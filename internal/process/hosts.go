package process

import (
	"context"
	"fmt"
	"strings"

	"github.com/kav91/nri-vmware-esxi/internal/load"
	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/vmware/govmomi/object"
)

func createHostSamples(timestamp int64) {
	ctx := context.Background()

	// create new entities for each host
	for _, host := range load.Hosts {
		entityName := host.Summary.Config.Name + ":hypervisor"

		if load.Args.DatacenterLocation != "" {
			entityName = load.Args.DatacenterLocation + ":" + entityName
		}

		entityName = strings.ToLower(entityName)
		entityName = strings.ReplaceAll(entityName, ".", "-")

		workingEntity := setEntity(entityName, "vmware") // default type instance
		workingEntity.SetInventoryItem("name", "value", fmt.Sprintf("%v:%d", entityName, timestamp))

		// create SystemSample metric set
		systemSampleMetricSet := workingEntity.NewMetricSet("SystemSample")

		// defaults
		checkError(systemSampleMetricSet.SetMetric("integration_version", load.IntegrationVersion, metric.ATTRIBUTE))
		checkError(systemSampleMetricSet.SetMetric("integration_name", load.IntegrationName, metric.ATTRIBUTE))
		checkError(systemSampleMetricSet.SetMetric("timestamp", timestamp, metric.GAUGE))
		checkError(systemSampleMetricSet.SetMetric("instanceType", "vmware-hypervisor", metric.ATTRIBUTE))
		checkError(systemSampleMetricSet.SetMetric("datacenterLocation", load.Args.DatacenterLocation, metric.ATTRIBUTE))
		checkError(systemSampleMetricSet.SetMetric("type", "hypervisor", metric.ATTRIBUTE))
		checkError(systemSampleMetricSet.SetMetric("hostname", entityName, metric.ATTRIBUTE))

		// host system
		checkError(systemSampleMetricSet.SetMetric("hostSystem", host.Self.Value, metric.ATTRIBUTE))
		checkError(systemSampleMetricSet.SetMetric("hypervisorHostname", entityName, metric.ATTRIBUTE))

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
