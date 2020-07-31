package collect

import (
	"context"
	"testing"

	"github.com/newrelic/nri-vsphere/internal/config"

	logrus "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/govmomi"

	"github.com/vmware/govmomi/simulator"
	"github.com/vmware/govmomi/view"
)

func TestCollectData(t *testing.T) {
	c := &config.Config{
		Logrus: logrus.New(),
	}

	ctx := context.Background()

	//SettingUp Simulator
	model := simulator.VPX()
	defer model.Remove()
	require.NoError(t, model.Create())

	s := model.Service.NewServer()
	var err error
	c.VMWareClient, err = govmomi.NewClient(ctx, s.URL, true)
	require.NoError(t, err)

	c.ViewManager = view.NewManager(c.VMWareClient.Client)
	_ = CollectData(c)

	assert.Len(t, c.Datacenters, model.Datacenter)
	assert.Len(t, c.Datacenters[0].Datastores, model.Datastore)
	assert.Len(t, c.Datacenters[0].Hosts, model.Host+model.ClusterHost)
	// by default 2 RP are created for the 2 groups of vm
	assert.Len(t, c.Datacenters[0].ResourcePools, 2)
	assert.Len(t, c.Datacenters[0].Clusters, model.Cluster)
	assert.Len(t, c.Datacenters[0].VirtualMachines, (model.Machine*model.Host)+(model.Machine*model.Cluster))

	// Folder structure of the generated vcenter
	// /DC0
	// /DC0/vm
	// /DC0/vm/DC0_H0_VM0
	// /DC0/vm/DC0_H0_VM1
	// /DC0/vm/DC0_C0_RP0_VM0
	// /DC0/vm/DC0_C0_RP0_VM1
	// /DC0/host
	// /DC0/host/DC0_H0
	// /DC0/host/DC0_H0/DC0_H0
	// /DC0/host/DC0_H0/Resources
	// /DC0/host/DC0_C0
	// /DC0/host/DC0_C0/DC0_C0_H0
	// /DC0/host/DC0_C0/DC0_C0_H1
	// /DC0/host/DC0_C0/DC0_C0_H2
	// /DC0/host/DC0_C0/Resources
	// /DC0/datastore
	// /DC0/datastore/LocalDS_0
	// /DC0/network
	// /DC0/network/VM Network
	// /DC0/network/DVS0
	// /DC0/network/DVS0-DVUplinks-9
	// /DC0/network/DC0_DVPG0
}
