package collect

import (
	"context"
	"github.com/newrelic/nri-vmware-vsphere/internal/load"
	"github.com/vmware/govmomi/vim25/mo"
)

// Clusters VMWare
func Clusters(config *load.Config) {
	ctx := context.Background()
	m := config.ViewManager

	for i, dc := range config.Datacenters {
		cv, err := m.CreateContainerView(ctx, dc.Datacenter.Reference(), []string{"ComputeResource"}, true)
		if err != nil {
			config.Logrus.WithError(err).Fatal("failed to create ComputeResource container view")
		}
		defer cv.Destroy(ctx)
		var clusters []mo.ClusterComputeResource
		// Reference: https://code.vmware.com/apis/704/vsphere/vim.ClusterComputeResource.html
		err = cv.Retrieve(
			ctx,
			[]string{"ClusterComputeResource"},
			[]string{"summary", "host", "datastore", "name", "network"},
			&clusters)
		if err != nil {
			config.Logrus.WithError(err).Fatal("failed to retrieve ClusterComputeResource")
		}
		for j := 0; j < len(clusters); j++ {
			config.Datacenters[i].Clusters[clusters[j].Self] = &clusters[j]
		}
	}
}
