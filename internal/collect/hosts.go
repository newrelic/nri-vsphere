package collect

import (
	"context"

	"github.com/kav91/nri-vmware-esxi/internal/load"
	"github.com/vmware/govmomi"
)

// Hosts VMWare
func Hosts(c *govmomi.Client) {
	ctx := context.Background()
	var err error
	m := load.ViewManager

	load.HostSystemContainerView, err = m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"HostSystem"}, true)
	if err != nil {
		load.Logrus.WithError(err).Fatal("failed to create HostSystem container view")
	}

	defer load.HostSystemContainerView.Destroy(ctx)

	// Retrieve summary property for all hosts
	// Reference: http://pubs.vmware.com/vsphere-60/topic/com.vmware.wssdk.apiref.doc/vim.HostSystem.html
	err = load.HostSystemContainerView.Retrieve(ctx, []string{"HostSystem"}, []string{"summary", "config", "network", "vm"}, &load.Hosts)
	if err != nil {
		load.Logrus.WithError(err).Fatal("failed to retrieve HostSystems")
	}
}
