/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package load

import (
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

type mor = types.ManagedObjectReference

// Datacenter struct
type Datacenter struct {
	Datacenter      *mo.Datacenter
	Hosts           map[mor]*mo.HostSystem
	Clusters        map[mor]*mo.ClusterComputeResource
	ResourcePools   map[mor]*mo.ResourcePool
	Datastores      map[mor]*mo.Datastore
	Networks        map[mor]*mo.Network
	VirtualMachines map[mor]*mo.VirtualMachine
}

// NewDatacenter Initialize datacenter struct
func NewDatacenter(datacenter *mo.Datacenter) Datacenter {
	return Datacenter{
		Datacenter:      datacenter,
		Hosts:           make(map[mor]*mo.HostSystem),
		Clusters:        make(map[mor]*mo.ClusterComputeResource),
		ResourcePools:   make(map[mor]*mo.ResourcePool),
		Datastores:      make(map[mor]*mo.Datastore),
		Networks:        make(map[mor]*mo.Network),
		VirtualMachines: make(map[mor]*mo.VirtualMachine),
	}
}

// FindResourcePool finds the ResourcePool associated to a Cluster except for the default resource pool
func (dc *Datacenter) FindResourcePool(clusterReference mor) (rp []*mo.ResourcePool) {
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
