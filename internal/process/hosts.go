package process

import (
	"context"
	"fmt"
	"strings"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/nri-vmware-vsphere/internal/load"
	"github.com/vmware/govmomi/object"
)

func createHostSamples(timestamp int64) {
	ctx := context.Background()

	// create new entities for each host
	for _, dc := range load.Datacenters {
		for _, host := range dc.Hosts {
			hostConfigName := host.Summary.Config.Name
			entityName := hostConfigName + ":host"
			datacenterName := dc.Datacenter.Name
			clusterName := dc.Clusters[host.Parent.Reference()].Name
			if load.IsVcenterAPIType {
				// to avoid redundant names when the host doesn't belong to any cluster
				if clusterName == hostConfigName {
					entityName = datacenterName + ":" + entityName
				} else {
					entityName = datacenterName + ":" + clusterName + ":" + entityName
				}
			}

			if load.Args.DatacenterLocation != "" {
				entityName = load.Args.DatacenterLocation + ":" + entityName
			}
			entityName = strings.ToLower(entityName)
			entityName = strings.ReplaceAll(entityName, ".", "-")

			// bios uuid identifies the host unequivocally and is available from vcenter/host api
			// uuid := integration.IDAttribute{Key: "uuid", Value: host.Summary.Hardware.Uuid}
			workingEntity, err := load.Integration.Entity(entityName, "vsphere")
			if err != nil {
				load.Logrus.WithError(err).Error("failed to create entity")
			}

			workingEntity.SetInventoryItem("name", "value", fmt.Sprintf("%v:%d", entityName, timestamp))

			systemSampleMetricSet := workingEntity.NewMetricSet("VSphereHostSample")

			if load.IsVcenterAPIType {
				checkError(systemSampleMetricSet.SetMetric("datacenterName", datacenterName, metric.ATTRIBUTE))
				checkError(systemSampleMetricSet.SetMetric("clusterName", clusterName, metric.ATTRIBUTE))

				resourcePools := dc.FindResourcePool(host.Parent.Reference())
				resourcePoolList := ""
				for _, rp := range resourcePools {
					resourcePoolList += rp.Name + "|"
				}
				checkError(systemSampleMetricSet.SetMetric("resourcePoolNameList", resourcePoolList, metric.ATTRIBUTE))
			}
			datastoreList := ""
			for _, ds := range host.Datastore {
				datastoreList += dc.Datastores[ds].Name + "|"
			}
			checkError(systemSampleMetricSet.SetMetric("datastoreNameList", datastoreList, metric.ATTRIBUTE))
			if load.Args.DatacenterLocation != "" {
				checkError(systemSampleMetricSet.SetMetric("datacenterLocation", load.Args.DatacenterLocation, metric.ATTRIBUTE))
			}
			checkError(systemSampleMetricSet.SetMetric("hypervisorHostName", hostConfigName, metric.ATTRIBUTE))

			checkError(systemSampleMetricSet.SetMetric("vmCount", len(host.Vm), metric.GAUGE))

			//TODO change to list | and add map of networks
			for i, nw := range host.Network {

				network := object.NewNetwork(load.NetworkContainerView.Client(), nw)
				hostNetwork, err := network.ObjectName(ctx)
				if err == nil {
					checkError(systemSampleMetricSet.SetMetric(fmt.Sprintf("network.%d", i), hostNetwork, metric.ATTRIBUTE))
				}
			}
			// memory
			memoryTotalBytes := float64(host.Summary.Hardware.MemorySize)
			checkError(systemSampleMetricSet.SetMetric("mem.size", memoryTotalBytes, metric.GAUGE))

			memoryUsedBytes := float64(host.Summary.QuickStats.OverallMemoryUsage) * 1e+6
			checkError(systemSampleMetricSet.SetMetric("mem.usage", memoryUsedBytes, metric.GAUGE))

			memoryFreeBytes := memoryTotalBytes - memoryUsedBytes
			checkError(systemSampleMetricSet.SetMetric("mem.free", memoryFreeBytes, metric.GAUGE))

			// cpu
			CPUCores := host.Summary.Hardware.NumCpuCores
			checkError(systemSampleMetricSet.SetMetric("cpu.cores", CPUCores, metric.GAUGE))

			CPUThreads := host.Summary.Hardware.NumCpuThreads
			checkError(systemSampleMetricSet.SetMetric("cpu.threads", CPUThreads, metric.GAUGE))

			CPUMhz := host.Summary.Hardware.CpuMhz
			checkError(systemSampleMetricSet.SetMetric("cpu.coreMHz", CPUMhz, metric.GAUGE))

			TotalMHz := float64(CPUMhz) * float64(CPUCores)
			checkError(systemSampleMetricSet.SetMetric("cpu.totalMHz", TotalMHz, metric.GAUGE))

			cpuPercent := (float64(host.Summary.QuickStats.OverallCpuUsage) / TotalMHz) * 100
			checkError(systemSampleMetricSet.SetMetric("cpu.percent", cpuPercent, metric.GAUGE))
			checkError(systemSampleMetricSet.SetMetric("cpu.overallUsage", host.Summary.QuickStats.OverallCpuUsage, metric.GAUGE))

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
			checkError(systemSampleMetricSet.SetMetric("disk.TotalBytes", diskTotalBytes, metric.GAUGE))

		}
	}
}
