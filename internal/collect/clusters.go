package collect

import (
	"context"

	"github.com/newrelic/nri-vmware-vsphere/internal/load"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/mo"
)

// Clusters VMWare
func Clusters(c *govmomi.Client) {
	ctx := context.Background()
	m := load.ViewManager

	for i, dc := range load.Datacenters {
		cv, err := m.CreateContainerView(ctx, dc.Datacenter.Reference(), []string{"ComputeResource"}, true)
		if err != nil {
			load.Logrus.WithError(err).Fatal("failed to create ComputeResource container view")
		}
		defer cv.Destroy(ctx)
		var clusters []mo.ComputeResource
		// Reference: https://code.vmware.com/apis/704/vsphere/vim.ComputeResource.html
		err = cv.Retrieve(
			ctx,
			[]string{"ComputeResource"},
			[]string{"summary", "host", "resourcePool", "name"},
			&clusters)
		if err != nil {
			load.Logrus.WithError(err).Fatal("failed to retrieve ComputeResource")
		}
		for j := 0; j < len(clusters); j++ {
			load.Datacenters[i].Clusters[clusters[j].Self] = &clusters[j]
		}
	}
}
