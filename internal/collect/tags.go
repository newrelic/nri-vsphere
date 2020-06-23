package collect

import (
	"context"
	"fmt"

	"github.com/vmware/govmomi/vapi/tags"

	"github.com/newrelic/nri-vsphere/internal/load"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

// separate in chunks of 2000 objects following performance recommendation
// https://www.vmware.com/content/dam/digitalmarketing/vmware/en/pdf/techpaper/performance/tagging-vsphere67-perf.pdf
const maxBatchSize = 2000

// TagCategories retreive all tag categories from vcenter and store them for future use
func collectTagCategories(tc load.TagCategories, tm *tags.Manager) error {
	ctx := context.Background()
	categories, err := tm.GetCategories(ctx)
	if err != nil {
		return err
	}
	for _, c := range categories {
		tc[c.ID] = c.Name
	}
	return nil
}

func collectTags(config *load.Config, managedObjectsSlice interface{}, dc *load.Datacenter) error {
	var ref []mo.Reference

	if !(config.Args.EnableVsphereTags && config.IsVcenterAPIType) {
		return nil
	}

	switch obs := managedObjectsSlice.(type) {
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
		return fmt.Errorf("type unknown")
	}

	if len(ref) < 1 {
		return nil
	}

	tagsByObject, err := getTags(ref, config.TagsManager, config.TagCategories)
	if err != nil {
		return fmt.Errorf("failed to collect tags:%v", err)
	}

	dc.AddTags(tagsByObject)

	return nil
}

func getTags(ref []mo.Reference, tm *tags.Manager, tc load.TagCategories) (map[types.ManagedObjectReference][]load.Tag, error) {
	ctx := context.Background()

	var attachedTags []tags.AttachedTags
	for i := 0; i < len(ref); i += maxBatchSize {
		batch := ref[i:min(i+maxBatchSize, len(ref))]

		tags, err := tm.GetAttachedTagsOnObjects(ctx, batch)
		if err != nil {
			return nil, fmt.Errorf("fail to get tags:%v", err)
		}

		attachedTags = append(attachedTags, tags...)
	}

	tagsByObject := make(map[types.ManagedObjectReference][]load.Tag)
	for _, object := range attachedTags {
		r := object.ObjectID.Reference()
		for _, tag := range object.Tags {
			if category, ok := tc[tag.CategoryID]; ok {
				t := load.Tag{Name: tag.Name, Category: category}
				tagsByObject[r] = append(tagsByObject[r], t)
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
