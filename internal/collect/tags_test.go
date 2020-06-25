package collect

import (
	"context"
	"testing"

	"github.com/newrelic/nri-vsphere/internal/load"

	"github.com/stretchr/testify/assert"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/simulator"
	"github.com/vmware/govmomi/vapi/rest"
	"github.com/vmware/govmomi/vapi/tags"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"

	_ "github.com/vmware/govmomi/vapi/simulator"
)

func TestCollectTagsByID(t *testing.T) {
	simulator.Run(func(ctx context.Context, vc *vim25.Client) error {
		c := rest.NewClient(vc)
		_ = c.Login(ctx, simulator.DefaultLogin)

		m := tags.NewManager(c)

		tagsByID := make(map[string]load.Tag)

		assert.NoError(t, collectTagsByID(tagsByID, m))
		assert.Len(t, tagsByID, 0)

		categoryName := "my-category"
		categoryID, err := m.CreateCategory(ctx, &tags.Category{Name: categoryName})
		assert.NoError(t, err)

		tagName := "vm-tag"
		tagID, err := m.CreateTag(ctx, &tags.Tag{CategoryID: categoryID, Name: tagName})
		assert.NoError(t, err)

		assert.NoError(t, collectTagsByID(tagsByID, m))
		assert.Equal(t, categoryName, tagsByID[tagID].Category)
		assert.Equal(t, tagName, tagsByID[tagID].Name)

		return nil
	})
}

func TestGetTags(t *testing.T) {
	simulator.Test(func(ctx context.Context, vc *vim25.Client) {
		c := rest.NewClient(vc)
		_ = c.Login(ctx, simulator.DefaultLogin)

		m := tags.NewManager(c)

		tagsByID := make(map[string]load.Tag)

		categoryName := "my-category"
		categoryID, err := m.CreateCategory(ctx, &tags.Category{
			AssociableTypes: []string{"VirtualMachine"},
			Cardinality:     "SINGLE",
			Name:            categoryName,
		})
		assert.NoError(t, err)
		tagName := "vm-tag"
		tagID, err := m.CreateTag(ctx, &tags.Tag{CategoryID: categoryID, Name: tagName})
		assert.NoError(t, err)

		assert.NoError(t, collectTagsByID(tagsByID, m))

		vm, err := find.NewFinder(vc).VirtualMachine(ctx, "DC0_H0_VM0")
		assert.NoError(t, err)
		err = m.AttachTag(ctx, tagID, vm.Reference())
		assert.NoError(t, err)

		vms := []mo.Reference{vm.Reference()}
		tagsByCategory, _ := getTags(vms, m, tagsByID)
		assert.Len(t, tagsByCategory, 1)
		assert.NotEmpty(t, tagsByCategory[vm.Reference()][0])
		assert.Equal(t, tagName, tagsByCategory[vm.Reference()][0].Name)

	})
}
