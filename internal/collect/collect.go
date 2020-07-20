package collect

import (
	"sync"
	"time"

	"github.com/newrelic/nri-vsphere/internal/config"
	"github.com/newrelic/nri-vsphere/internal/model/tag"
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

func CollectData(config *config.Config) {
	now := time.Now()

	if config.Args.EnableVsphereTags && config.IsVcenterAPIType {
		err := tag.BuildTagCache(config.TagsManager)
		if err != nil {
			config.Logrus.WithError(err).Error("failed to build tag cache")
		}
	}

	err := Datacenters(config)
	if err != nil {
		return
	}
	config.Logrus.WithField("seconds", time.Since(now).Seconds()).Debug("after collecting dc data")

	// fetch vmware data async
	var wg sync.WaitGroup
	wg.Add(6)
	go func() {
		defer wg.Done()
		VirtualMachines(config)
		config.Logrus.WithField("seconds", time.Since(now).Seconds()).Debug("after collecting vms data")
	}()
	go func() {
		defer wg.Done()
		Networks(config)
		config.Logrus.WithField("seconds", time.Since(now).Seconds()).Debug("after collecting network data")

	}()
	go func() {
		defer wg.Done()
		Hosts(config)
		config.Logrus.WithField("seconds", time.Since(now).Seconds()).Debug("after collecting hosts data")
	}()
	go func() {
		defer wg.Done()
		Datastores(config)
		config.Logrus.WithField("seconds", time.Since(now).Seconds()).Debug("after collecting datastores data")

	}()
	go func() {
		defer wg.Done()
		Clusters(config)
		config.Logrus.WithField("seconds", time.Since(now).Seconds()).Debug("after collecting clusters data")

	}()
	go func() {
		defer wg.Done()
		ResourcePools(config)
		config.Logrus.WithField("seconds", time.Since(now).Seconds()).Debug("after collecting resourcepools data")

	}()
	wg.Wait()
}
