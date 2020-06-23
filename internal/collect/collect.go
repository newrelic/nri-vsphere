package collect

import (
	"github.com/newrelic/nri-vsphere/internal/load"
	"sync"
)

func CollectData(config *load.Config) {

	Datacenters(config)

	// fetch vmware data async
	var wg sync.WaitGroup
	wg.Add(6)
	go func() {
		defer wg.Done()
		VirtualMachines(config)
	}()
	go func() {
		defer wg.Done()
		Networks(config)

	}()
	go func() {
		defer wg.Done()
		Hosts(config)

	}()
	go func() {
		defer wg.Done()
		Datastores(config)

	}()
	go func() {
		defer wg.Done()
		Clusters(config)
	}()
	go func() {
		defer wg.Done()
		ResourcePools(config)
	}()
	wg.Wait()
}
