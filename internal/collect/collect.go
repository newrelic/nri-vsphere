package collect

import (
	"errors"
	"github.com/newrelic/nri-vsphere/internal/config"
	"sync"
)

const (
	DATACENTER      = "Datacenter"
	VIRTUAL_MACHINE = "VirtualMachine"
	DATASTORE       = "Datastore"
	HOST            = "HostSystem"
	RESOURCE_POOL   = "ResourcePool"
	NETWORK         = "Network"
	CLUSTER         = "ClusterComputeResource"
)

func CollectData(config *config.Config) error {

	if config.TagCollectionEnabled() {
		err := config.TagCollector.BuildTagCache()
		if err != nil {
			config.Logrus.WithError(err).Error("failed to build tag cache")
		}
	}
	config.Logrus.WithField("seconds", config.Uptime()).Debug("after collecting tags")

	err := Datacenters(config)
	if err != nil {
		return err
	}
	config.Logrus.WithField("seconds", config.Uptime()).Debug("after collecting dc data")

	if len(config.Datacenters) == 0 {
		return errors.New("no datacenter was collected. this is most likely an error in your filter")
	}

	// fetch vmware data async
	var wg sync.WaitGroup
	wg.Add(6)
	go func() {
		defer wg.Done()
		VirtualMachines(config)
		config.Logrus.WithField("seconds", config.Uptime()).Debug("after collecting vms data")
	}()
	go func() {
		defer wg.Done()
		Networks(config)
		config.Logrus.WithField("seconds", config.Uptime()).Debug("after collecting network data")

	}()
	go func() {
		defer wg.Done()
		Hosts(config)
		config.Logrus.WithField("seconds", config.Uptime()).Debug("after collecting hosts data")
	}()
	go func() {
		defer wg.Done()
		Datastores(config)
		config.Logrus.WithField("seconds", config.Uptime()).Debug("after collecting datastores data")

	}()
	go func() {
		defer wg.Done()
		Clusters(config)
		config.Logrus.WithField("seconds", config.Uptime()).Debug("after collecting clusters data")

	}()
	go func() {
		defer wg.Done()
		ResourcePools(config)
		config.Logrus.WithField("seconds", config.Uptime()).Debug("after collecting resourcepools data")

	}()
	wg.Wait()

	return nil
}
