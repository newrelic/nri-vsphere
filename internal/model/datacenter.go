// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package model

import (
	"sync"

	"github.com/newrelic/nri-vsphere/internal/events"
	"github.com/newrelic/nri-vsphere/internal/performance"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

type mor = types.ManagedObjectReference

// Datacenter struct
type Datacenter struct {
	Datacenter      *mo.Datacenter
	EventDispacher  *events.EventDispacher
	Hosts           map[mor]*mo.HostSystem
	Clusters        map[mor]*mo.ClusterComputeResource
	ResourcePools   map[mor]*mo.ResourcePool
	Datastores      map[mor]*mo.Datastore
	Networks        map[mor]*mo.Network
	VirtualMachines map[mor]*mo.VirtualMachine
	PerfMetrics     map[mor][]performance.PerfMetric
	PerfMetricsMux  sync.Mutex
}

// NewDatacenter Initialize datacenter struct
func NewDatacenter(datacenter mo.Datacenter) *Datacenter {
	return &Datacenter{
		Datacenter:      &datacenter,
		Hosts:           make(map[mor]*mo.HostSystem),
		Clusters:        make(map[mor]*mo.ClusterComputeResource),
		ResourcePools:   make(map[mor]*mo.ResourcePool),
		Datastores:      make(map[mor]*mo.Datastore),
		Networks:        make(map[mor]*mo.Network),
		VirtualMachines: make(map[mor]*mo.VirtualMachine),
		PerfMetrics:     make(map[mor][]performance.PerfMetric),
	}
}

// FindResourcePools finds the ResourcePool associated to a Cluster except for the default resource pool
func (dc *Datacenter) FindResourcePools(clusterReference mor) (rp []*mo.ResourcePool) {
	for _, resourcePool := range dc.ResourcePools {
		// Default ResourcePool is the root, the rest should be listed as child
		if (resourcePool.Owner == clusterReference) && (len(resourcePool.ResourcePool) > 0) {
			for _, rpChild := range resourcePool.ResourcePool {
				rp = append(rp, dc.ResourcePools[rpChild])
			}
		}
	}
	return
}

// FindHost returns the child Host for a computeResource
func (dc *Datacenter) FindHost(computeResourceReference mor) *mo.HostSystem {
	for _, host := range dc.Hosts {
		if host.Parent.Reference() == computeResourceReference {
			return host
		}
	}
	return nil
}

// GetResourcePool returns the name of the Resource Pool if is not the default
func (dc *Datacenter) GetResourcePool(resourcePoolReference mor) *mo.ResourcePool {
	if !dc.IsDefaultResourcePool(resourcePoolReference) {
		return dc.ResourcePools[resourcePoolReference]
	}
	return nil
}

// IsDefaultResourcePool returns true if the resource pool is the default
func (dc *Datacenter) IsDefaultResourcePool(resourcePoolReference mor) bool {
	if rp, ok := dc.ResourcePools[resourcePoolReference]; ok {
		if rp.Parent.Type != "ResourcePool" {
			return true
		}
	}
	return false
}

// AddTags appends a tag batch to dc Tags map
func (dc *Datacenter) AddPerfMetrics(data map[types.ManagedObjectReference][]performance.PerfMetric) {
	dc.PerfMetricsMux.Lock()
	defer dc.PerfMetricsMux.Unlock()
	for m, value := range data {
		dc.PerfMetrics[m] = append(dc.PerfMetrics[m], value...)
	}
}

// GetPerfMetrics returns the slice of Perf metrics for the given object reference
func (dc *Datacenter) GetPerfMetrics(ref mor) []performance.PerfMetric {
	if perfMetrics, ok := dc.PerfMetrics[ref]; ok {
		return perfMetrics
	}
	return nil
}

func (dc *Datacenter) GetDatastore(ds mor) *mo.Datastore {
	return dc.Datastores[ds]
}

func (dc *Datacenter) GetNetwork(n mor) *mo.Network {
	return dc.Networks[n]
}

func (dc *Datacenter) GetHost(h mor) *mo.HostSystem {
	return dc.Hosts[h]
}

func (dc *Datacenter) GetCluster(c mor) *mo.ClusterComputeResource {
	return dc.Clusters[c]
}
