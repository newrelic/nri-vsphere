package collect

import (
	"context"

	"github.com/newrelic/nri-vmware-vsphere/internal/load"
	"github.com/vmware/govmomi/vim25/mo"
)

// Datacenters VMWare
func Datacenters(config *load.Config) {
	ctx := context.Background()
	m := config.ViewManager

	cv, err := m.CreateContainerView(ctx, config.VMWareClient.ServiceContent.RootFolder, []string{"Datacenter"}, true)
	if err != nil {
		config.Logrus.WithError(err).Fatal("failed to create Datacenter container view")
	}

	defer cv.Destroy(ctx)

	var datacenters []mo.Datacenter
	err = cv.Retrieve(ctx, []string{"Datacenter"}, nil, &datacenters)
	if err != nil {
		config.Logrus.WithError(err).Fatal("failed to retrieve Datacenter")
	}

	for i := range datacenters {
		config.Datacenters = append(config.Datacenters, load.NewDatacenter(&datacenters[i]))
	}
}
