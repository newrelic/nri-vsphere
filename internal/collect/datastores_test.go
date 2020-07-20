package collect

import (
	"context"

	"github.com/newrelic/nri-vsphere/internal/client"
	"github.com/newrelic/nri-vsphere/internal/config"
	"github.com/newrelic/nri-vsphere/internal/model/tag"

	"testing"

	. "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/simulator"
	"github.com/vmware/govmomi/vapi/rest"
	_ "github.com/vmware/govmomi/vapi/simulator"
	"github.com/vmware/govmomi/vapi/tags"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

func Test_ListDatastores_WithEmptyFilter_ReturnsAllDatastores(t *testing.T) {
	simulator.Run(func(ctx context.Context, vc *vim25.Client) error {
		vmClient, err := client.New(vc.URL().String(), "user", "pass", false)
		assert.NoError(t, err)
		vm := view.NewManager(vc)
		assert.NotNil(t, vm)
		// given
		cfg := &config.Config{VMWareClient: vmClient, ViewManager: vm, Logrus: StandardLogger()}
		cfg.Datacenters = append(cfg.Datacenters, getDatacenter(ctx, vm))

		// when
		Datastores(cfg)

		// then
		assert.True(t, len(cfg.Datacenters[0].Datastores) > 0)

		return nil
	})
}

func Test_ListDatastoress_WithNonEmptyFilter(t *testing.T) {
	simulator.Run(func(ctx context.Context, vc *vim25.Client) error {
		vmClient, err := client.New(vc.URL().String(), "user", "pass", false)
		assert.NoError(t, err)

		c := rest.NewClient(vc)
		err = c.Login(ctx, simulator.DefaultLogin)
		assert.NoError(t, err)

		m := tags.NewManager(c)
		addDsTag(ctx, m, vc, "region", "eu")
		addDsTag(ctx, m, vc, "env", "test")

		// given
		// given
		cfg := &config.Config{
			Args: config.ArgumentList{
				EnableVsphereTags: true,
			},
			IsVcenterAPIType: true,
			VMWareClient:     vmClient,
			ViewManager:      view.NewManager(vc),
			TagsManager:      m,
			Logrus:           StandardLogger(),
			TagsByID:         map[string]tag.Tag{},
		}

		_ = tag.BuildTagCache(cfg.TagsManager)

		tests := []struct {
			name string
			args string
			want int
		}{
			{
				name: "ByNonExistingTag",
				args: "key=value",
				want: 0,
			},
			{
				name: "ByExistingTag",
				args: "region=eu",
				want: 1,
			},
			{
				name: "ByMultipleMixedTags",
				args: "key=value env=test",
				want: 1,
			},
			{
				name: "ByMultipleExistingTags",
				args: "region=eu env=test",
				want: 1,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// we need to get the datacenter first
				cfg.Datacenters = append(cfg.Datacenters, getDatacenter(ctx, cfg.ViewManager))
				cfg.Datacenters[0].Datastores = make(map[types.ManagedObjectReference]*mo.Datastore)

				// when
				cfg.Args.IncludeTags = tt.args
				Datastores(cfg)

				// then
				assert.Equal(t, tt.want, len(cfg.Datacenters[0].Datastores))
			})
		}

		return nil
	})
}

func addDsTag(ctx context.Context, m *tags.Manager, vc *vim25.Client, category string, value string) {
	categoryID, _ := m.CreateCategory(ctx, &tags.Category{
		AssociableTypes: []string{DATASTORE},
		Cardinality:     "SINGLE",
		Name:            category,
	})
	tagID, _ := m.CreateTag(ctx, &tags.Tag{CategoryID: categoryID, Name: value})
	finder := find.NewFinder(vc)
	dss, _ := finder.DatastoreList(ctx, "/DC0/...")
	for _, ds := range dss {
		_ = m.AttachTag(ctx, tagID, ds.Reference())
	}
}
