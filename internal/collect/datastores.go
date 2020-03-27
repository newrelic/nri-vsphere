package collect

import (
	"context"

	"github.com/newrelic/nri-vmware-vsphere/internal/load"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/mo"
)

// Datastores collects data of all datastores
func Datastores(c *govmomi.Client) {
	ctx := context.Background()
	m := load.ViewManager

	for i, dc := range load.Datacenters {
		cv, err := m.CreateContainerView(ctx, dc.Datacenter.Reference(), []string{"Datastore"}, true)
		if err != nil {
			load.Logrus.WithError(err).Fatal("failed to create Datastore container view")
		}
		defer cv.Destroy(ctx)

		var datastores []mo.Datastore
		// Reference: https://code.vmware.com/apis/42/vsphere/doc/vim.Datastore.html
		err = cv.Retrieve(ctx, []string{"Datastore"}, nil, &datastores)
		if err != nil {
			load.Logrus.WithError(err).Fatal("failed to retrieve Datastore")
		}
		for j := 0; j < len(datastores); j++ {
			load.Datacenters[i].Datastores[datastores[j].Self] = &datastores[j]
		}
	}
}
