package collect

import (
	"context"

	"github.com/newrelic/nri-vmware-vsphere/internal/load"
	"github.com/vmware/govmomi"
)

// Networks ESXi
func Networks(c *govmomi.Client) {
	ctx := context.Background()
	var err error
	m := load.ViewManager

	load.NetworkContainerView, err = m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"Network"}, true)
	if err != nil {
		load.Logrus.WithError(err).Fatal("failed to create Network container view")
	}
	defer load.NetworkContainerView.Destroy(ctx)

	// Reference: http://pubs.vmware.com/vsphere-60/topic/com.vmware.wssdk.apiref.doc/vim.Network.html
	err = load.NetworkContainerView.Retrieve(ctx, []string{"Network"}, nil, &load.Networks)
	if err != nil {
		load.Logrus.WithError(err).Fatal("failed to retrieve Networks")
	}
}
