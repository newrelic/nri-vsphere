package collect

import (
	"github.com/newrelic/nri-vsphere/internal/load"
	"sync"
	"time"
)

func CollectData(config *load.Config) {

	Datacenters(config)
	config.Logrus.WithField("seconds", time.Since(load.Now).Seconds()).Debug("after collecting dc data")

	// fetch vmware data async
	var wg sync.WaitGroup
	wg.Add(6)
	go func() {
		defer wg.Done()
		VirtualMachines(config)
		config.Logrus.WithField("seconds", time.Since(load.Now).Seconds()).Debug("after collecting vms data")
	}()
	go func() {
		defer wg.Done()
		Networks(config)
		config.Logrus.WithField("seconds", time.Since(load.Now).Seconds()).Debug("after collecting network data")

	}()
	go func() {
		defer wg.Done()
		Hosts(config)
		config.Logrus.WithField("seconds", time.Since(load.Now).Seconds()).Debug("after collecting hosts data")
	}()
	go func() {
		defer wg.Done()
		Datastores(config)
		config.Logrus.WithField("seconds", time.Since(load.Now).Seconds()).Debug("after collecting datastores data")

	}()
	go func() {
		defer wg.Done()
		Clusters(config)
		config.Logrus.WithField("seconds", time.Since(load.Now).Seconds()).Debug("after collecting clusters data")

	}()
	go func() {
		defer wg.Done()
		ResourcePools(config)
		config.Logrus.WithField("seconds", time.Since(load.Now).Seconds()).Debug("after collecting resourcepools data")

	}()
	wg.Wait()
}
