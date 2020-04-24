// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package process

import (
	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/nri-vsphere/internal/load"
)

func createResourcePoolSamples(config *load.Config, timestamp int64) {
	for _, dc := range config.Datacenters {
		for _, rp := range dc.ResourcePools {
			// Skip root default ResourcePool (not created by user)
			if dc.IsDefaultResourcePool(rp.Reference()) {
				continue
			}
			resourcePoolName := rp.Name
			datacenterName := dc.Datacenter.Name

			// Resource Pool could be owned by Cluster or a Host
			ownerName := ""
			if cluster, ok := dc.Clusters[rp.Owner]; ok {
				ownerName = cluster.Name
			} else if host := dc.FindHost(rp.Owner); host != nil {
				ownerName = host.Summary.Config.Name
			}
			entityName := ownerName + ":" + resourcePoolName

			entityName = sanitizeEntityName(config, entityName, datacenterName)

			ms := createNewEntityWithMetricSet(config, entityTypeResourcePool, entityName, entityName)

			checkError(config, ms.SetMetric("resourcePoolName", resourcePoolName, metric.ATTRIBUTE))
			if config.Args.DatacenterLocation != "" {
				checkError(config, ms.SetMetric("datacenterLocation", config.Args.DatacenterLocation, metric.ATTRIBUTE))
			}
			if config.IsVcenterAPIType {
				checkError(config, ms.SetMetric("datacenterName", datacenterName, metric.ATTRIBUTE))
				if cluster, ok := dc.Clusters[rp.Owner]; ok {
					checkError(config, ms.SetMetric("clusterName", cluster.Name, metric.ATTRIBUTE))
				}
			}

			memTotal := (rp.Runtime.Memory.ReservationUsed + rp.Runtime.Memory.UnreservedForPool) / (1e6)
			checkError(config, ms.SetMetric("mem.size", memTotal, metric.GAUGE))

			summary := rp.Summary.GetResourcePoolSummary()
			// esxi api reports nil quickstats
			if summary.QuickStats != nil {
				checkError(config, ms.SetMetric("mem.usage", summary.QuickStats.GuestMemoryUsage, metric.GAUGE))
				memFree := memTotal - summary.QuickStats.GuestMemoryUsage
				checkError(config, ms.SetMetric("mem.free", memFree, metric.GAUGE))
				checkError(config, ms.SetMetric("mem.ballooned", summary.QuickStats.BalloonedMemory, metric.GAUGE))
				checkError(config, ms.SetMetric("mem.swapped", summary.QuickStats.SwappedMemory, metric.GAUGE))
				checkError(config, ms.SetMetric("cpu.overallUsage", summary.QuickStats.OverallCpuUsage, metric.GAUGE))
			}
			cpuTotal := rp.Runtime.Cpu.ReservationUsed + rp.Runtime.Cpu.UnreservedForPool
			checkError(config, ms.SetMetric("cpu.totalMHz", cpuTotal, metric.GAUGE))

			checkError(config, ms.SetMetric("vmCount", len(rp.Vm), metric.GAUGE))

			checkError(config, ms.SetMetric("overallStatus", string(rp.OverallStatus), metric.ATTRIBUTE))

		}
	}
}
