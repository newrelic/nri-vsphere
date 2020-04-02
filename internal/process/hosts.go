package process

import (
	"context"
	"fmt"
	"strings"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-vmware-vsphere/internal/load"
	"github.com/vmware/govmomi/object"
)

func createHostSamples(timestamp int64) {
	ctx := context.Background()

	// create new entities for each host
	for _, dc := range load.Datacenters {
		for _, host := range dc.Hosts {
			// entityName := host.Summary.Config.Name + ":hypervisor"

			// if load.Args.DatacenterLocation != "" {
			// 	entityName = load.Args.DatacenterLocation + ":" + entityName
			// }

			// entityName = strings.ToLower(entityName)
			// entityName = strings.ReplaceAll(entityName, ".", "-")

			// workingEntity := setEntity(entityName, "vmware") // default type instance
			// workingEntity.SetInventoryItem("name", "value", fmt.Sprintf("%v:%d", entityName, timestamp))

			// // create SystemSample metric set
			// systemSampleMetricSet := workingEntity.NewMetricSet("SystemSample")

			// bios uuid identifies the host unequivocally and is available from vcenter/host api
			uuid := integration.IDAttribute{Key: "uuid", Value: host.Summary.Hardware.Uuid}
			workingEntity, err := load.Integration.Entity(host.Summary.Config.Name, "host", uuid)
			if err != nil {
				load.Logrus.WithError(err).Error("failed to create entity")
			}

			hostName := strings.ToLower(host.Summary.Config.Name)

			workingEntity.SetInventoryItem("name", "value", fmt.Sprintf("%v:%d", hostName, timestamp))

			systemSampleMetricSet := workingEntity.NewMetricSet("VSphereHostSample")

			checkError(systemSampleMetricSet.SetMetric("datacenterName", dc.Datacenter.ManagedEntity.Name, metric.ATTRIBUTE))

			// cluster := dc.FindCluster(host.Reference())
			cluster := dc.Clusters[host.Parent.Reference()]
			checkError(systemSampleMetricSet.SetMetric("clusterName", cluster.Name, metric.ATTRIBUTE))

			resourcePools := dc.FindResourcePool(cluster.Reference())
			resourcePoolList := ""
			for _, rp := range resourcePools {
				resourcePoolList += rp.Name + ","
			}
			checkError(systemSampleMetricSet.SetMetric("resourcePoolNameList", resourcePoolList, metric.ATTRIBUTE))

			datastoreList := ""
			for _, ds := range host.Datastore {
				datastoreList += dc.Datastores[ds].Name + ","
			}
			checkError(systemSampleMetricSet.SetMetric("datastoreNameList", datastoreList, metric.ATTRIBUTE))

			// defaults
			checkError(systemSampleMetricSet.SetMetric("integration_version", load.IntegrationVersion, metric.ATTRIBUTE))
			checkError(systemSampleMetricSet.SetMetric("integration_name", load.IntegrationName, metric.ATTRIBUTE))
			checkError(systemSampleMetricSet.SetMetric("timestamp", timestamp, metric.GAUGE))
			checkError(systemSampleMetricSet.SetMetric("instanceType", "vmware-hypervisor", metric.ATTRIBUTE))
			checkError(systemSampleMetricSet.SetMetric("datacenterLocation", load.Args.DatacenterLocation, metric.ATTRIBUTE))
			checkError(systemSampleMetricSet.SetMetric("type", "hypervisor", metric.ATTRIBUTE))
			checkError(systemSampleMetricSet.SetMetric("hostname", hostName, metric.ATTRIBUTE))

			// host system
			checkError(systemSampleMetricSet.SetMetric("hostSystem", host.Self.Value, metric.ATTRIBUTE))
			checkError(systemSampleMetricSet.SetMetric("hypervisorHostname", hostName, metric.ATTRIBUTE))

			// vms
			checkError(systemSampleMetricSet.SetMetric("configName", host.Summary.Config.Name, metric.ATTRIBUTE))
			checkError(systemSampleMetricSet.SetMetric("vmCount", len(host.Vm), metric.GAUGE))

			for i, nw := range host.Network {

				network := object.NewNetwork(load.NetworkContainerView.Client(), nw)
				hostNetwork, err := network.ObjectName(ctx)
				if err == nil {
					checkError(systemSampleMetricSet.SetMetric(fmt.Sprintf("network.%d", i), hostNetwork, metric.ATTRIBUTE))
				}
			}

			// SystemSample metrics

			// // memory
			memoryTotalBytes := float64(host.Summary.Hardware.MemorySize)
			checkError(systemSampleMetricSet.SetMetric("memoryTotalBytes", memoryTotalBytes, metric.GAUGE))
			checkError(systemSampleMetricSet.SetMetric("systemMemoryBytes", memoryTotalBytes, metric.GAUGE))
			memoryUsedBytes := float64(host.Summary.QuickStats.OverallMemoryUsage) * 1e+6
			memoryFreeBytes := memoryTotalBytes - memoryUsedBytes
			checkError(systemSampleMetricSet.SetMetric("memoryUsedBytes", memoryUsedBytes, metric.GAUGE))
			checkError(systemSampleMetricSet.SetMetric("memoryFreeBytes", memoryFreeBytes, metric.GAUGE))

			// // cpu
			CPUCores := host.Summary.Hardware.NumCpuCores
			checkError(systemSampleMetricSet.SetMetric("coreCount", CPUCores, metric.GAUGE))

			CPUThreads := host.Summary.Hardware.NumCpuThreads
			checkError(systemSampleMetricSet.SetMetric("cpuThreads", CPUThreads, metric.GAUGE))

			checkError(systemSampleMetricSet.SetMetric("overallCpuUsageMHz", host.Summary.QuickStats.OverallCpuUsage, metric.GAUGE))

			CPUMhz := host.Summary.Hardware.CpuMhz
			checkError(systemSampleMetricSet.SetMetric("cpuMHz", CPUMhz, metric.GAUGE))

			TotalMHz := float64(CPUMhz) * float64(CPUCores)
			checkError(systemSampleMetricSet.SetMetric("totalMHz", TotalMHz, metric.GAUGE))

			cpuPercent := (float64(host.Summary.QuickStats.OverallCpuUsage) / TotalMHz) * 100
			checkError(systemSampleMetricSet.SetMetric("cpuPercent", cpuPercent, metric.GAUGE))

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

			checkError(systemSampleMetricSet.SetMetric("diskTotalBytes", diskTotalBytes, metric.GAUGE))

		}
	}
}
