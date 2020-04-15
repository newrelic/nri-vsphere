package process

import (
	"strings"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/nri-vmware-vsphere/internal/load"
)

func createResourcePoolSamples(config *load.Config, timestamp int64) {
	for _, dc := range config.Datacenters {
		for _, rp := range dc.ResourcePools {
			// Skip root default ResourcePool (not created by user)
			if rp.Parent.Type != "ResourcePool" {
				continue
			}
			resourcePoolName := rp.Name
			datacenterName := dc.Datacenter.Name
			ownerName := ""
			if cluster, ok := dc.Clusters[rp.Owner]; ok {
				ownerName = cluster.Name
			} else if host := dc.FindHost(rp.Owner); host != nil {
				ownerName = host.Summary.Config.Name
			}
			entityName := ownerName + ":" + resourcePoolName + ":resourcePool"
			if config.IsVcenterAPIType {
				entityName = datacenterName + ":" + entityName
			}
			if config.Args.DatacenterLocation != "" {
				entityName = config.Args.DatacenterLocation + ":" + entityName
			}

			entityName = strings.ToLower(entityName)
			entityName = strings.ReplaceAll(entityName, ".", "-")

			workingEntity, err := config.Integration.Entity(entityName, "vsphere-resourcepool")
			if err != nil {
				config.Logrus.WithError(err).Error("failed to create entity")
			}
			// entity displayName
			workingEntity.SetInventoryItem("vsphereResourcePool", "name", entityName)

			ms := workingEntity.NewMetricSet("VSphereResourcePoolSample")

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
