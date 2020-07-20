package tag

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"

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

var tagByIDCache TagsByID

// TagsByID stores tags per object
type TagsByObject = map[mor][]Tag

var tagsByObjectCache TagsByObject

var tagMux sync.Mutex

var filterTags []Tag

// separate in chunks of 2000 objects following performance recommendation
// https://www.vmware.com/content/dam/digitalmarketing/vmware/en/pdf/techpaper/performance/tagging-vsphere67-perf.pdf
const maxBatchSize = 2000

func init() {
	tagByIDCache = TagsByID{}
	tagsByObjectCache = TagsByObject{}
	filterTags = []Tag{}
}

// BuildTagCache caches all tag and categories from vCenter and stores them for future reference
func BuildTagCache(tm *tags.Manager) error {
	ctx := context.Background()

	categories, err := tm.GetCategories(ctx)
	if err != nil {
		return err
	}
	categoriesByID := make(map[string]string)
	for _, c := range categories {
		categoriesByID[c.ID] = c.Name
	}

	ts, err := tm.GetTags(ctx)
	if err != nil {
		return err
	}
	for _, tag := range ts {
		if category, ok := categoriesByID[tag.CategoryID]; ok {
			tagByIDCache[tag.ID] = Tag{Name: tag.Name, Category: category}
		}
	}
	return nil
}

// GetTagById gets a tag by it's id
func GetTagByID(id string) Tag {
	return tagByIDCache[id]
}

// GetTagsByCategories return a map of tags categories and the corresponding tags associated to the object
func GetTagsByCategories(ref mor) map[string]string {
	tagsByCategory := make(map[string]string)

	if ts, ok := tagsByObjectCache[ref]; ok {
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
func GetTagsForObject(tm *tags.Manager, or mor) []Tag {
	return tagsByObjectCache[or]
}

func FetchTagsForObjects(tm *tags.Manager, objectsSlice interface{}) (TagsByObject, error) {
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

	tagsByObject, err := getTags(ref, tm)
	if err != nil {
		return nil, fmt.Errorf("failed to collect tags:%v", err)
	}
	cacheTags(tagsByObject)

	return tagsByObjectCache, nil
}

// cache tags grouped by object reference
func cacheTags(tagsByObject TagsByObject) {
	tagMux.Lock()
	defer tagMux.Unlock()
	for or, ts := range tagsByObject {
		tagsByObjectCache[or] = append(tagsByObjectCache[or], ts...)
	}
}

// return all tags attached to objects in ref grouped by the object reference
func getTags(ref []mo.Reference, tm *tags.Manager) (TagsByObject, error) {
	ctx := context.Background()

	var attachedTags []tags.AttachedTags
	for i := 0; i < len(ref); i += maxBatchSize {
		batch := ref[i:min(i+maxBatchSize, len(ref))]

		result, err := tm.ListAttachedTagsOnObjects(ctx, batch)
		if err != nil {
			return nil, fmt.Errorf("fail to get tags:%v", err)
		}
		attachedTags = append(attachedTags, result...)
	}

	tagsByObject := make(map[mor][]Tag)
	for _, tag := range attachedTags {
		or := tag.ObjectID.Reference()
		for _, tagID := range tag.TagIDs {
			if tag, ok := tagByIDCache[tagID]; ok {
				tagsByObject[or] = append(tagsByObject[or], tag)
			}
		}
	}
	return tagsByObject, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ParseFilterTagExpression converts a filter tag expression into a slice of Tag
// example: tag=value tag1=value1 --> t:{c: tag, n: value}, t1:{c: tag1, n: value1}
func ParseFilterTagExpression(tagExpression string) {
	if len(tagExpression) == 0 {
		return
	}

	fields := strings.Fields(tagExpression)
	for _, t := range fields {
		kv := strings.Split(t, "=")
		if len(kv) != 2 {
			continue
		}

		filterTags = append(filterTags, Tag{
			Category: kv[0],
			Name:     kv[1],
		})
	}
}

// MatchObjectTags checks if any tag in 'objectTags' matches any of the 'filterTags'
func MatchObjectTags(objectTags []Tag) bool {
	if len(objectTags) == 0 {
		return false
	}
	for _, ft := range filterTags {
		for _, ot := range objectTags {
			if ot.Category == ft.Category && ot.Name == ft.Name {
				return true
			}
		}
	}
	return false
}
