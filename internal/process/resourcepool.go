// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package process

import (
	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/nri-vsphere/internal/config"
)

func createResourcePoolSamples(config *config.Config) {
	for _, dc := range config.Datacenters {
		for _, rp := range dc.ResourcePools {

			// filtering here will to avoid sending data to backend
			if config.TagFilteringEnabled() && !config.TagCollector.MatchObjectTags(rp.Self) {
				continue
			}

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

			e, ms, err := createNewEntityWithMetricSet(config, entityTypeResourcePool, entityName, entityName)
			if err != nil {
				config.Logrus.WithError(err).WithField("resourcePoolName", entityName).Error("failed to create metricSet")
				continue
			}

			checkError(config.Logrus, ms.SetMetric("resourcePoolName", resourcePoolName, metric.ATTRIBUTE))

			if config.Args.DatacenterLocation != "" {
				checkError(config.Logrus, ms.SetMetric("datacenterLocation", config.Args.DatacenterLocation, metric.ATTRIBUTE))
			}

			if config.IsVcenterAPIType {
				checkError(config.Logrus, ms.SetMetric("datacenterName", datacenterName, metric.ATTRIBUTE))
				if cluster, ok := dc.Clusters[rp.Owner]; ok {
					checkError(config.Logrus, ms.SetMetric("clusterName", cluster.Name, metric.ATTRIBUTE))
				}
			}

			memTotal := (rp.Runtime.Memory.ReservationUsed + rp.Runtime.Memory.UnreservedForPool) / (1 << 20)
			checkError(config.Logrus, ms.SetMetric("mem.size", memTotal, metric.GAUGE))

			summary := rp.Summary.GetResourcePoolSummary()
			// esxi api reports nil quickstats
			if summary.QuickStats != nil {
				checkError(config.Logrus, ms.SetMetric("mem.usage", summary.QuickStats.GuestMemoryUsage, metric.GAUGE))
				memFree := memTotal - summary.QuickStats.GuestMemoryUsage
				checkError(config.Logrus, ms.SetMetric("mem.free", memFree, metric.GAUGE))
				checkError(config.Logrus, ms.SetMetric("mem.ballooned", summary.QuickStats.BalloonedMemory, metric.GAUGE))
				checkError(config.Logrus, ms.SetMetric("mem.swapped", summary.QuickStats.SwappedMemory, metric.GAUGE))
				checkError(config.Logrus, ms.SetMetric("cpu.overallUsage", summary.QuickStats.OverallCpuUsage, metric.GAUGE))
			}
			cpuTotal := rp.Runtime.Cpu.ReservationUsed + rp.Runtime.Cpu.UnreservedForPool
			checkError(config.Logrus, ms.SetMetric("cpu.totalMHz", cpuTotal, metric.GAUGE))

			checkError(config.Logrus, ms.SetMetric("vmCount", len(rp.Vm), metric.GAUGE))

			checkError(config.Logrus, ms.SetMetric("overallStatus", string(rp.OverallStatus), metric.ATTRIBUTE))

			// Tags
			if config.TagCollectionEnabled() {
				tagsByCategory := config.TagCollector.GetTagsByCategories(rp.Self)
				for k, v := range tagsByCategory {
					checkError(config.Logrus, ms.SetMetric(tagsPrefix+k, v, metric.ATTRIBUTE))
					// add tags to inventory due to the inventory workaround
					checkError(config.Logrus, e.SetInventoryItem("tags", tagsPrefix+k, v))
				}
			}
			// Performance metrics
			if config.PerfMetricsCollectionEnabled() {
				perfMetrics := dc.GetPerfMetrics(rp.Self)
				for _, perfMetric := range perfMetrics {
					checkError(config.Logrus, ms.SetMetric(perfMetricPrefix+perfMetric.Counter, perfMetric.Value, metric.GAUGE))
				}
			}
		}
	}
}
