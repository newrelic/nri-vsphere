package collect

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/simulator"
	"github.com/vmware/govmomi/vapi/rest"
	"github.com/vmware/govmomi/vapi/tags"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"

	_ "github.com/vmware/govmomi/vapi/simulator"
)

func TestTagCategories(t *testing.T) {
	simulator.Run(func(ctx context.Context, vc *vim25.Client) error {
		c := rest.NewClient(vc)
		_ = c.Login(ctx, simulator.DefaultLogin)

		m := tags.NewManager(c)

		var tc = make(map[string]string)

		assert.NoError(t, collectTagCategories(tc, m))
		assert.Len(t, tc, 0)

		categoryName := "my-category"
		categoryID, err := m.CreateCategory(ctx, &tags.Category{
			AssociableTypes: []string{"VirtualMachine"},
			Cardinality:     "SINGLE",
			Name:            categoryName,
		})
		if err != nil {
			return err
		}

		assert.NoError(t, collectTagCategories(tc, m))
		assert.Equal(t, categoryName, tc[categoryID])

		return nil
	})
}

func TestGetTags(t *testing.T) {
	simulator.Test(func(ctx context.Context, vc *vim25.Client) {
		c := rest.NewClient(vc)
		_ = c.Login(ctx, simulator.DefaultLogin)

		m := tags.NewManager(c)

		var tc = make(map[string]string)

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

		assert.NoError(t, collectTagCategories(tc, m))
		assert.Equal(t, categoryName, tc[categoryID])

		vm, err := find.NewFinder(vc).VirtualMachine(ctx, "DC0_H0_VM0")
		assert.NoError(t, err)
		err = m.AttachTag(ctx, tagID, vm.Reference())
		assert.NoError(t, err)

		vms := []mo.Reference{vm.Reference()}
		tagsByCategory, _ := getTags(vms, m, tc)
		assert.NotEmpty(t, tagsByCategory[vm.Reference()][0])
		assert.Equal(t, tagName, tagsByCategory[vm.Reference()][0].Name)

	})
}
