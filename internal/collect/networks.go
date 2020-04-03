package collect

import (
	"context"

	"github.com/newrelic/nri-vmware-vsphere/internal/load"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/mo"
)

// Networks ESXi
func Networks(c *govmomi.Client) {
	ctx := context.Background()
	m := load.ViewManager

	for i, dc := range load.Datacenters {
		cv, err := m.CreateContainerView(ctx, dc.Datacenter.Reference(), []string{"Network"}, true)
		if err != nil {
			load.Logrus.WithError(err).Fatal("failed to create Network container view")
		}
		defer cv.Destroy(ctx)

		var networks []mo.Network
		// Reference: http://pubs.vmware.com/vsphere-60/topic/com.vmware.wssdk.apiref.doc/vim.Network.html
		err = cv.Retrieve(ctx, []string{"Network"}, nil, &networks)
		if err != nil {
			load.Logrus.WithError(err).Fatal("failed to retrieve Networks")
		}
		for j := 0; j < len(networks); j++ {
			load.Datacenters[i].Networks[networks[j].Self] = &networks[j]
		}
	}
}
