package process

import (
	"fmt"
	"strings"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/nri-vmware-vsphere/internal/load"
)

func createHostSamples(config *load.Config, timestamp int64) {
	for _, dc := range config.Datacenters {
		for _, host := range dc.Hosts {
			hostConfigName := host.Summary.Config.Name
			entityName := hostConfigName + ":host"
			datacenterName := dc.Datacenter.Name
			clusterName := dc.Clusters[host.Parent.Reference()].Name
			if config.IsVcenterAPIType {
				// to avoid redundant names when the host doesn't belong to any cluster
				if clusterName == hostConfigName {
					entityName = datacenterName + ":" + entityName
				} else {
					entityName = datacenterName + ":" + clusterName + ":" + entityName
				}
			}

			if config.Args.DatacenterLocation != "" {
				entityName = config.Args.DatacenterLocation + ":" + entityName
			}
			entityName = strings.ToLower(entityName)
			entityName = strings.ReplaceAll(entityName, ".", "-")

			// bios uuid identifies the host unequivocally and is available from vcenter/host api
			// uuid := integration.IDAttribute{Key: "uuid", Value: host.Summary.Hardware.Uuid}
			workingEntity, err := config.Integration.Entity(entityName, "vsphere")
			if err != nil {
				config.Logrus.WithError(err).Error("failed to create entity")
			}

			workingEntity.SetInventoryItem("name", "value", fmt.Sprintf("%v:%d", entityName, timestamp))

			ms := workingEntity.NewMetricSet("VSphereHostSample")

			if config.IsVcenterAPIType {
				checkError(config, ms.SetMetric("datacenterName", datacenterName, metric.ATTRIBUTE))
				checkError(config, ms.SetMetric("clusterName", clusterName, metric.ATTRIBUTE))
			}
			resourcePools := dc.FindResourcePool(host.Parent.Reference())
			resourcePoolList := ""
			for _, rp := range resourcePools {
				resourcePoolList += rp.Name + "|"
			}
			checkError(config, ms.SetMetric("resourcePoolNameList", resourcePoolList, metric.ATTRIBUTE))

			datastoreList := ""
			for _, ds := range host.Datastore {
				datastoreList += dc.Datastores[ds].Name + "|"
			}
			checkError(config, ms.SetMetric("datastoreNameList", datastoreList, metric.ATTRIBUTE))
			if config.Args.DatacenterLocation != "" {
				checkError(config, ms.SetMetric("datacenterLocation", config.Args.DatacenterLocation, metric.ATTRIBUTE))
			}
			checkError(config, ms.SetMetric("hypervisorHostName", hostConfigName, metric.ATTRIBUTE))
			checkError(config, ms.SetMetric("uuid", host.Summary.Hardware.Uuid, metric.ATTRIBUTE))

			checkError(config, ms.SetMetric("vmCount", len(host.Vm), metric.GAUGE))

			networkList := ""
			for _, nw := range host.Network {
				networkList += dc.Networks[nw].Name + "|"
			}
			checkError(config, ms.SetMetric("networkNameList", networkList, metric.ATTRIBUTE))

			// memory
			memoryTotalBytes := float64(host.Summary.Hardware.MemorySize)
			checkError(config, ms.SetMetric("mem.size", memoryTotalBytes, metric.GAUGE))

			memoryUsedBytes := float64(host.Summary.QuickStats.OverallMemoryUsage) * 1e+6
			checkError(config, ms.SetMetric("mem.usage", memoryUsedBytes, metric.GAUGE))

			memoryFreeBytes := memoryTotalBytes - memoryUsedBytes
			checkError(config, ms.SetMetric("mem.free", memoryFreeBytes, metric.GAUGE))

			// cpu
			CPUCores := host.Summary.Hardware.NumCpuCores
			checkError(config, ms.SetMetric("cpu.cores", CPUCores, metric.GAUGE))

			CPUThreads := host.Summary.Hardware.NumCpuThreads
			checkError(config, ms.SetMetric("cpu.threads", CPUThreads, metric.GAUGE))

			CPUMhz := host.Summary.Hardware.CpuMhz
			checkError(config, ms.SetMetric("cpu.coreMHz", CPUMhz, metric.GAUGE))

			TotalMHz := float64(CPUMhz) * float64(CPUCores)
			checkError(config, ms.SetMetric("cpu.totalMHz", TotalMHz, metric.GAUGE))

			cpuPercent := (float64(host.Summary.QuickStats.OverallCpuUsage) / TotalMHz) * 100
			checkError(config, ms.SetMetric("cpu.percent", cpuPercent, metric.GAUGE))
			checkError(config, ms.SetMetric("cpu.overallUsage", host.Summary.QuickStats.OverallCpuUsage, metric.GAUGE))

			// disk
			diskTotalBytes := int64(0)
			if host.Config != nil {
				if host.Config.FileSystemVolume != nil {
					for _, mount := range host.Config.FileSystemVolume.MountInfo {
						capacity := mount.Volume.GetHostFileSystemVolume().Capacity
						diskTotalBytes += capacity
					}
				}
			}
			checkError(config, ms.SetMetric("disk.TotalBytes", diskTotalBytes, metric.GAUGE))

		}
	}
}
