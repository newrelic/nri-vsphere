package collect

import (
	"context"

	"github.com/newrelic/nri-vmware-vsphere/internal/load"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/mo"
)

// Hosts VMWare
func Hosts(c *govmomi.Client) {
	ctx := context.Background()
	var err error
	m := load.ViewManager

	for i, dc := range load.Datacenters {

		load.HostSystemContainerView, err = m.CreateContainerView(ctx, dc.Datacenter.Reference(), []string{"HostSystem"}, true)
		if err != nil {
			load.Logrus.WithError(err).Fatal("failed to create HostSystem container view")
		}

		defer load.HostSystemContainerView.Destroy(ctx)

		var hosts []mo.HostSystem
		// Reference: http://pubs.vmware.com/vsphere-60/topic/com.vmware.wssdk.apiref.doc/vim.HostSystem.html
		err = load.HostSystemContainerView.Retrieve(
			ctx,
			[]string{"HostSystem"},
			[]string{"summary", "config", "network", "vm", "parent", "datastore"},
			&hosts)
		if err != nil {
			load.Logrus.WithError(err).Fatal("failed to retrieve HostSystems")
		}
		for j := 0; j < len(hosts); j++ {
			load.Datacenters[i].Hosts[hosts[j].Self] = &hosts[j]
		}
	}
}
