package collect

import (
	"context"
	"github.com/newrelic/nri-vsphere/internal/model"
	"github.com/vmware/govmomi/vim25/mo"
	"testing"

	"github.com/newrelic/nri-vsphere/internal/client"
	"github.com/newrelic/nri-vsphere/internal/config"
	"github.com/newrelic/nri-vsphere/internal/tag"

	logrus "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/simulator"
	"github.com/vmware/govmomi/vapi/rest"
	_ "github.com/vmware/govmomi/vapi/simulator"
	"github.com/vmware/govmomi/vapi/tags"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25"
)

func getDatacenter(ctx context.Context, vm *view.Manager) *model.Datacenter {
	cv, _ := vm.CreateContainerView(ctx, vm.Client().ServiceContent.RootFolder, []string{DATACENTER}, false)

	var datacenters []mo.Datacenter
	_ = cv.Retrieve(ctx, []string{DATACENTER}, []string{"name"}, &datacenters)
	return model.NewDatacenter(datacenters[0])
}

func Test_ListDatacenters_WithEmptyFilter_ReturnsAllDatacenters(t *testing.T) {

	simulator.Run(func(ctx context.Context, vc *vim25.Client) error {
		vmClient, err := client.New(vc.URL().String(), "user", "pass", false)
		assert.NoError(t, err)
		vm := view.NewManager(vc)
		assert.NotNil(t, vm)

		// given
		config := &config.Config{VMWareClient: vmClient, ViewManager: vm, Logrus: logrus.StandardLogger()}

		// when
		err = Datacenters(config)
		assert.NoError(t, err)

		// then
		assert.Len(t, config.Datacenters, 1)

		return nil
	})
}

func Test_ListDatacenters_WithNonEmptyFilter(t *testing.T) {
	simulator.Run(func(ctx context.Context, vc *vim25.Client) error {
		vmClient, err := client.New(vc.URL().String(), "user", "pass", false)
		assert.NoError(t, err)

		c := rest.NewClient(vc)
		err = c.Login(ctx, simulator.DefaultLogin)
		assert.NoError(t, err)

		m := tags.NewManager(c)
		addTag(ctx, m, vc, "region", "eu")
		addTag(ctx, m, vc, "env", "test")

		// given
		collector := tag.NewCollector(m, logrus.StandardLogger())
		_ = collector.BuildTagCache()

		cfg := &config.Config{
			Args: config.ArgumentList{
				EnableVsphereTags: true,
			},
			IsVcenterAPIType: true,
			VMWareClient:     vmClient,
			ViewManager:      view.NewManager(vc),
			TagCollector:     collector,
			Logrus:           logrus.StandardLogger(),
		}

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
				cfg.Datacenters = nil
				// make sure tag filtering is enabled
				cfg.Args.IncludeTags = tt.args
				collector.ParseFilterTagExpression(tt.args)

				// when
				_ = Datacenters(cfg)

				// then
				assert.Equal(t, tt.want, len(cfg.Datacenters))
			})
		}
		return nil
	})
}

func addTag(ctx context.Context, m *tags.Manager, vc *vim25.Client, category string, value string) {
	categoryID, _ := m.CreateCategory(ctx, &tags.Category{
		AssociableTypes: []string{DATACENTER},
		Cardinality:     "SINGLE",
		Name:            category,
	})
	tagID, _ := m.CreateTag(ctx, &tags.Tag{CategoryID: categoryID, Name: value})
	vm, _ := find.NewFinder(vc).Datacenter(ctx, "DC0")
	_ = m.AttachTag(ctx, tagID, vm.Reference())
}
