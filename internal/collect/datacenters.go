package collect

import (
	"context"

	"github.com/newrelic/nri-vmware-vsphere/internal/load"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/mo"
)

// Datacenters VMWare
func Datacenters(c *govmomi.Client) {
	ctx := context.Background()
	m := load.ViewManager

	cv, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"Datacenter"}, true)
	if err != nil {
		load.Logrus.WithError(err).Fatal("failed to create Datacenter container view")
	}

	defer cv.Destroy(ctx)

	var datacenters []mo.Datacenter
	err = cv.Retrieve(ctx, []string{"Datacenter"}, nil, &datacenters)
	if err != nil {
		load.Logrus.WithError(err).Fatal("failed to retrieve Datacenter")
	}

	for i := range datacenters {
		load.Datacenters = append(load.Datacenters, load.NewDatacenter(&datacenters[i]))
	}
}
