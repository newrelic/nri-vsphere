package tag

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/vmware/govmomi/vapi/tags"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

type mor = types.ManagedObjectReference

type Tag struct {
	Name     string
	Category string
}

// TagsByID stores tags per tag id
type TagsByID map[string]Tag

// TagsByID stores tags per object
type TagsByObject = map[mor][]Tag

// separate in chunks of 2000 objects following performance recommendation
// https://www.vmware.com/content/dam/digitalmarketing/vmware/en/pdf/techpaper/performance/tagging-vsphere67-perf.pdf
const maxBatchSize = 2000

type Collector struct {
	tm     *tags.Manager
	logger *logrus.Logger

	tagByIDCache      TagsByID
	tagsByObjectCache TagsByObject
	filterTags        []Tag
	mutex             *sync.Mutex
}

// ParseFilterTagExpression converts a filter tag expression into a slice of Tag
// example: tag=value tag1=value1 --> t:{c: tag, n: value}, t1:{c: tag1, n: value1}
// each invocation of this function resets any previously created filter
func (c *Collector) ParseFilterTagExpression(tagFilterExpression string) {
	if len(tagFilterExpression) == 0 {
		return
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()
	// clear previous filter tags if any
	c.filterTags = nil

	fields := strings.Fields(tagFilterExpression)
	for _, t := range fields {
		kv := strings.Split(t, "=")
		if len(kv) != 2 {
			c.logger.WithField("tag", t).Warn("invalid tag definition")
			continue
		}

		c.filterTags = append(c.filterTags, Tag{
			Category: kv[0],
			Name:     kv[1],
		})
	}
}

// BuildTagCache caches all tag and categories from vCenter and stores them for future reference
// each invocation of this func will clear any previously cached values
func (c *Collector) BuildTagCache() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	//clear previous cache if any
	c.tagByIDCache = TagsByID{}

	ctx := context.Background()

	categories, err := c.tm.GetCategories(ctx)
	if err != nil {
		return err
	}
	categoriesByID := make(map[string]string)
	for _, c := range categories {
		categoriesByID[c.ID] = c.Name
	}

	ts, err := c.tm.GetTags(ctx)
	if err != nil {
		return err
	}
	for _, t := range ts {
		if category, ok := categoriesByID[t.CategoryID]; ok {
			c.tagByIDCache[t.ID] = Tag{Name: t.Name, Category: category}
		}
	}
	return nil
}

// GetTagById gets a tag by it's id
func (c *Collector) GetTagByID(id string) Tag {
	return c.tagByIDCache[id]
}

// GetTagsByCategories return a map of tags categories and the corresponding tags associated to the object
func (c *Collector) GetTagsByCategories(ref mor) map[string]string {
	tagsByCategory := make(map[string]string)

	if c.tagsByObjectCache == nil {
		c.logger.Fatal(">>>> tagsByObjectCache is nill")
	}
	if ts, ok := c.tagsByObjectCache[ref]; ok {
		sort.Slice(ts, func(i, j int) bool {
			return ts[i].Name < ts[j].Name
		})
		for _, t := range ts {
			if _, ok := tagsByCategory[t.Category]; ok {
				tagsByCategory[t.Category] = tagsByCategory[t.Category] + "|" + t.Name
			} else {
				tagsByCategory[t.Category] = t.Name
			}
		}
	}
	return tagsByCategory
}

// GetTagsForObject gets all tags for a object
func (c *Collector) GetTagsForObject(or mor) []Tag {
	return c.tagsByObjectCache[or]
}

func (c *Collector) FetchTagsForObjects(objectsSlice interface{}) (TagsByObject, error) {
	var ref []mo.Reference

	switch obs := objectsSlice.(type) {
	case []mo.VirtualMachine:
		for _, o := range obs {
			ref = append(ref, o.Self)
		}
	case []mo.Datacenter:
		for _, o := range obs {
			ref = append(ref, o.Self)
		}
	case []mo.ClusterComputeResource:
		for _, o := range obs {
			ref = append(ref, o.Self)
		}
	case []mo.Datastore:
		for _, o := range obs {
			ref = append(ref, o.Self)
		}
	case []mo.HostSystem:
		for _, o := range obs {
			ref = append(ref, o.Self)
		}
	case []mo.ResourcePool:
		for _, o := range obs {
			ref = append(ref, o.Self)
		}
	case []mo.Network:
		for _, o := range obs {
			ref = append(ref, o.Self)
		}
	default:
		return nil, fmt.Errorf("type unknown")
	}

	if len(ref) < 1 {
		return nil, nil
	}

	tagsByObject, err := c.getTags(ref)
	if err != nil {
		return nil, fmt.Errorf("failed to collect tags:%v", err)
	}
	c.cacheTags(tagsByObject)

	return tagsByObject, nil
}

// MatchObjectTags checks if any tag in the resource tags matches any of the 'filterTags'
func (c *Collector) MatchObjectTags(resource mor) bool {
	objectTags := c.tagsByObjectCache[resource.Reference()]
	if len(objectTags) == 0 {
		return false
	}

	return c.matchTags(objectTags)
}

func (c *Collector) matchTags(objectTags []Tag) bool {
	for _, ft := range c.filterTags {
		for _, ot := range objectTags {
			if ot.Category == ft.Category && ot.Name == ft.Name {
				return true
			}
		}
	}
	return false
}

// cache tags grouped by object reference
func (c *Collector) cacheTags(tagsByObject TagsByObject) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for or, ts := range tagsByObject {
		c.tagsByObjectCache[or] = append(c.tagsByObjectCache[or], ts...)
	}
}

// return all tags attached to objects in ref grouped by the object reference
func (c *Collector) getTags(ref []mo.Reference) (TagsByObject, error) {
	ctx := context.Background()

	var attachedTags []tags.AttachedTags
	for i := 0; i < len(ref); i += maxBatchSize {
		batch := ref[i:min(i+maxBatchSize, len(ref))]

		result, err := c.tm.ListAttachedTagsOnObjects(ctx, batch)
		if err != nil {
			return nil, fmt.Errorf("fail to get tags:%v", err)
		}
		attachedTags = append(attachedTags, result...)
	}

	tagsByObject := make(map[mor][]Tag)
	for _, tag := range attachedTags {
		or := tag.ObjectID.Reference()
		for _, tagID := range tag.TagIDs {
			if tag, ok := c.tagByIDCache[tagID]; ok {
				tagsByObject[or] = append(tagsByObject[or], tag)
			}
		}
	}
	return tagsByObject, nil
}

func NewCollector(tagManager *tags.Manager, logger *logrus.Logger) *Collector {
	return &Collector{
		tm:                tagManager,
		logger:            logger,
		tagsByObjectCache: TagsByObject{},
		tagByIDCache:      TagsByID{},
		filterTags:        []Tag{},
		mutex:             &sync.Mutex{},
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
