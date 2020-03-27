package collect

import (
	"context"

	"github.com/newrelic/nri-vmware-vsphere/internal/load"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/mo"
)

// ResourcePools VMWare
func ResourcePools(c *govmomi.Client) {
	ctx := context.Background()
	m := load.ViewManager

	for i, dc := range load.Datacenters {
		cv, err := m.CreateContainerView(ctx, dc.Datacenter.Reference(), []string{"ResourcePool"}, true)
		if err != nil {
			load.Logrus.WithError(err).Fatal("failed to create ResourcePool container view")
		}
		defer cv.Destroy(ctx)
		var resourcePools []mo.ResourcePool
		err = cv.Retrieve(
			ctx,
			[]string{"ResourcePool"},
			[]string{"summary", "owner", "runtime", "name", "resourcePool"},
			&resourcePools)
		if err != nil {
			load.Logrus.WithError(err).Fatal("failed to retrieve ResourcePools")
		}
		for j := 0; j < len(resourcePools); j++ {
			load.Datacenters[i].ResourcePools[resourcePools[j].Self] = &resourcePools[j]
		}
	}
}
