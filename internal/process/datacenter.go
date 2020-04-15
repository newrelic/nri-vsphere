package process

import (
	"strings"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/nri-vmware-vsphere/internal/load"
)

func createClusterSamples(config *load.Config, timestamp int64) {
	for _, dc := range config.Datacenters {

		//Creating entity name
		entityName := dc.Datacenter.Name + ":datacenter"

		if config.Args.DatacenterLocation != "" {
			entityName = config.Args.DatacenterLocation + ":" + entityName
		}
		entityName = strings.ToLower(entityName)
		entityName = strings.ReplaceAll(entityName, ".", "-")


		// Identifier for cluster entity
		workingEntity, err := config.Integration.Entity(entityName, "vsphere-datacenter")
		if err != nil {
			config.Logrus.WithError(err).Error("failed to create entity")
		}

		// entity displayName
		workingEntity.SetInventoryItem("vspheredatacenter", "name", entityName)
		ms := workingEntity.NewMetricSet("VSphereDatacenterSample")

		checkError(config, ms.SetMetric("overallStatus", string(dc.Datacenter.OverallStatus), metric.ATTRIBUTE))


		checkError(config, ms.SetMetric("datastores", len(dc.Datastores), metric.GAUGE))
		checkError(config, ms.SetMetric("hosts", len(dc.Hosts), metric.GAUGE))
		checkError(config, ms.SetMetric("vms", len(dc.VirtualMachines), metric.GAUGE))
		checkError(config, ms.SetMetric("networks", len(dc.Networks), metric.GAUGE))
		checkError(config, ms.SetMetric("resourcePools", len(dc.ResourcePools), metric.GAUGE))

		checkError(config, ms.SetMetric("datastores", len(dc.Datastores), metric.GAUGE))
		var totalDatastoreCapacity int64
		var totalDatastoreFreeSpace int64
		for _, datastore:=range dc.Datastores{
			totalDatastoreCapacity = totalDatastoreCapacity + datastore.Summary.FreeSpace
			totalDatastoreFreeSpace = totalDatastoreFreeSpace + datastore.Summary.Capacity
		}
		checkError(config, ms.SetMetric("totalDatastoreCapacity", totalDatastoreCapacity/(1<<30), metric.GAUGE))
		checkError(config, ms.SetMetric("totalDatastoreFreeSpace", totalDatastoreFreeSpace/(1<<30), metric.GAUGE))
		checkError(config, ms.SetMetric("totalDatastoreUsedSpace", (totalDatastoreFreeSpace - totalDatastoreCapacity)/(1<<30), metric.GAUGE))

		var totalMemoryHost int64
		var totalMemoryUsedHost int32
		var totalCpuHost int16
		var totalMHz float64
		var totalOverallCpuUsage float64

		for _, host:=range dc.Hosts{
			totalMHz = totalMHz + (float64(host.Summary.Hardware.CpuMhz) * float64(host.Summary.Hardware.NumCpuCores))
			totalOverallCpuUsage = totalOverallCpuUsage + float64(host.Summary.QuickStats.OverallCpuUsage)
			totalCpuHost = totalCpuHost + host.Summary.Hardware.NumCpuCores
			totalMemoryHost = totalMemoryHost + host.Summary.Hardware.MemorySize / 1e+6
			totalMemoryUsedHost = totalMemoryUsedHost + host.Summary.QuickStats.OverallMemoryUsage

		}
		cpuPercentHost := totalOverallCpuUsage / totalMHz * 100
		checkError(config, ms.SetMetric("cpuPercentHost", cpuPercentHost, metric.GAUGE))
		checkError(config, ms.SetMetric("totalCpuHost", totalCpuHost, metric.GAUGE))
		checkError(config, ms.SetMetric("cpu.PercentHost", cpuPercentHost, metric.GAUGE))
		checkError(config, ms.SetMetric("cpu.totalOverallUsage", totalOverallCpuUsage, metric.GAUGE))
		checkError(config, ms.SetMetric("cpu.totalMHz", totalMHz, metric.GAUGE))

		memoryPercentHost := float64(totalMemoryUsedHost) / float64(totalMemoryHost) * 100
		checkError(config, ms.SetMetric("totalMemoryUsedHost", totalMemoryUsedHost, metric.GAUGE))
		checkError(config, ms.SetMetric("totalMemoryHost", totalMemoryHost, metric.GAUGE))
		checkError(config, ms.SetMetric("memoryPercentHost", memoryPercentHost, metric.GAUGE))

	}
}
