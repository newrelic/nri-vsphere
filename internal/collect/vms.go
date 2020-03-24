package collect

import (
	"context"

	"github.com/newrelic/nri-vmware-vsphere/internal/load"
	"github.com/vmware/govmomi"
)

// VirtualMachines vms
func VirtualMachines(c *govmomi.Client) {
	ctx := context.Background()
	var err error
	m := load.ViewManager

	load.VirutalMachineContainerView, err = m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"VirtualMachine"}, true)
	if err != nil {
		load.Logrus.WithError(err).Fatal("failed to create VirtualMachine container view")
	}
	defer load.VirutalMachineContainerView.Destroy(ctx)

	// Retrieve summary property for all machines
	// Reference: http://pubs.vmware.com/vsphere-60/topic/com.vmware.wssdk.apiref.doc/vim.VirtualMachine.html
	err = load.VirutalMachineContainerView.Retrieve(ctx, []string{"VirtualMachine"}, []string{"summary", "network", "config", "guest", "runtime"}, &load.VirtualMachines)
	if err != nil {
		load.Logrus.WithError(err).Fatal("failed to retrieve VM Summaries")
	}
}
