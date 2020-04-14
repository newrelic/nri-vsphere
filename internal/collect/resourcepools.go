package collect

import (
	"context"

	"github.com/newrelic/nri-vmware-vsphere/internal/load"
	"github.com/vmware/govmomi/vim25/mo"
)

// ResourcePools VMWare
func ResourcePools(config *load.Config) {
	ctx := context.Background()
	m := config.ViewManager

	for i, dc := range config.Datacenters {
		cv, err := m.CreateContainerView(ctx, dc.Datacenter.Reference(), []string{"ResourcePool"}, true)
		if err != nil {
			config.Logrus.WithError(err).Fatal("failed to create ResourcePool container view")
		}
		defer cv.Destroy(ctx)
		var resourcePools []mo.ResourcePool
		err = cv.Retrieve(
			ctx,
			[]string{"ResourcePool"},
			[]string{"summary", "owner", "parent", "runtime", "name", "overallStatus", "vm"},
			&resourcePools)
		if err != nil {
			config.Logrus.WithError(err).Fatal("failed to retrieve ResourcePools")
		}
		for j := 0; j < len(resourcePools); j++ {
			config.Datacenters[i].ResourcePools[resourcePools[j].Self] = &resourcePools[j]
		}
	}
}
