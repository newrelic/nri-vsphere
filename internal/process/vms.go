package process

import (
	"context"
	"fmt"
	"strings"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/nri-vmware-vsphere/internal/load"
	"github.com/vmware/govmomi/object"
)

func createVirtualMachineSamples(timestamp int64) {
	ctx := context.Background()

	// create new entities for each vm
	for _, vm := range load.VirtualMachines {
		// resolve hypervisor host
		vmHost := vm.Summary.Runtime.Host
		hostSystem := object.NewHostSystem(load.HostSystemContainerView.Client(), *vmHost)
		hypervisorHost, err := hostSystem.ObjectName(ctx)
		discoveredHost := findHost(hostSystem.Reference())
		entityName := ""

		if discoveredHost.Self.Value != "" && discoveredHost.Summary.Config.Name != "" {
			entityName = discoveredHost.Summary.Config.Name + ":" + vm.Summary.Config.Name + ":vm"
		} else {
			entityName = hypervisorHost + ":" + vm.Summary.Config.Name + ":vm"
		}

		if load.Args.DatacenterLocation != "" {
			entityName = load.Args.DatacenterLocation + ":" + entityName
		}

		entityName = strings.ToLower(entityName)
		entityName = strings.ReplaceAll(entityName, ".", "-")

		// not available if VM is offline
		vmHostname := vm.Summary.Guest.HostName

		workingEntity := setEntity(entityName, "vmware") // default type instance
		workingEntity.SetInventoryItem("name", "value", fmt.Sprintf("%v:%d", entityName, timestamp))

		// create SystemSample metric set
		systemSampleMetricSet := workingEntity.NewMetricSet("SystemSample")

		// defaults
		checkError(systemSampleMetricSet.SetMetric("integration_version", load.IntegrationVersion, metric.ATTRIBUTE))
		checkError(systemSampleMetricSet.SetMetric("integration_name", load.IntegrationName, metric.ATTRIBUTE))
		checkError(systemSampleMetricSet.SetMetric("timestamp", timestamp, metric.GAUGE))
		checkError(systemSampleMetricSet.SetMetric("instanceType", "vmware-guest", metric.ATTRIBUTE))
		checkError(systemSampleMetricSet.SetMetric("datacenterLocation", load.Args.DatacenterLocation, metric.ATTRIBUTE))
		checkError(systemSampleMetricSet.SetMetric("type", "virtualmachine", metric.ATTRIBUTE))
		checkError(systemSampleMetricSet.SetMetric("hostname", vm.Summary.Config.Name, metric.ATTRIBUTE))

		// host system
		checkError(systemSampleMetricSet.SetMetric("hostSystem", vm.Summary.Runtime.Host.Value, metric.ATTRIBUTE))
		if err == nil && hypervisorHost != "" {
			checkError(systemSampleMetricSet.SetMetric("hypervisorHostname", hypervisorHost, metric.ATTRIBUTE))
		}

		// vm
		checkError(systemSampleMetricSet.SetMetric("configName", vm.Summary.Config.Name, metric.ATTRIBUTE))
		checkError(systemSampleMetricSet.SetMetric("vmHostname", vmHostname, metric.ATTRIBUTE))
		checkError(systemSampleMetricSet.SetMetric("vmNumber", vm.Self.Value, metric.ATTRIBUTE))

		for i, nw := range vm.Network {
			network := object.NewNetwork(load.NetworkContainerView.Client(), nw)
			vmNetwork, err := network.ObjectName(ctx)
			if err == nil {
				checkError(systemSampleMetricSet.SetMetric(fmt.Sprintf("network.%d", i), vmNetwork, metric.ATTRIBUTE))
			}
		}

		operatingSystem := determineOS(vm.Summary.Config.GuestFullName)
		checkError(systemSampleMetricSet.SetMetric("operatingSystem", operatingSystem, metric.ATTRIBUTE))
		checkError(systemSampleMetricSet.SetMetric("guestFullName", vm.Summary.Config.GuestFullName, metric.ATTRIBUTE))

		// SystemSample metrics

		// memory
		memoryTotalBytes := float64(vm.Summary.Config.MemorySizeMB) * 1e+6
		checkError(systemSampleMetricSet.SetMetric("memoryTotalBytes", memoryTotalBytes, metric.GAUGE))
		checkError(systemSampleMetricSet.SetMetric("systemMemoryBytes", memoryTotalBytes, metric.GAUGE))
		memoryUsedBytes := float64(vm.Summary.QuickStats.GuestMemoryUsage) * 1e+6
		memoryFreeBytes := memoryTotalBytes - memoryUsedBytes
		checkError(systemSampleMetricSet.SetMetric("memoryUsedBytes", memoryUsedBytes, metric.GAUGE))
		checkError(systemSampleMetricSet.SetMetric("memoryFreeBytes", memoryFreeBytes, metric.GAUGE))

		// cpu
		checkError(systemSampleMetricSet.SetMetric("coreCount", vm.Summary.Config.NumCpu, metric.GAUGE))
		checkError(systemSampleMetricSet.SetMetric("overallCpuUsageMHz", vm.Summary.QuickStats.OverallCpuUsage, metric.GAUGE))

		cpuAllocationLimit := float64(0)
		if vm.Config.CpuAllocation.Limit != nil {
			cpuAllocationLimit = float64(*vm.Config.CpuAllocation.Limit)
		}

		checkError(systemSampleMetricSet.SetMetric("cpuAllocationLimit", cpuAllocationLimit, metric.GAUGE))

		if discoveredHost.Self.Value != "" {
			CPUMhz := discoveredHost.Summary.Hardware.CpuMhz
			CPUCores := discoveredHost.Summary.Hardware.NumCpuCores
			CPUThreads := discoveredHost.Summary.Hardware.NumCpuThreads
			TotalMHz := float64(CPUMhz) * float64(CPUCores)
			checkError(systemSampleMetricSet.SetMetric("hypervisorCpuThreads", CPUThreads, metric.GAUGE))
			checkError(systemSampleMetricSet.SetMetric("hypervisorCpuMhz", CPUMhz, metric.GAUGE))
			checkError(systemSampleMetricSet.SetMetric("hypervisorCpuCores", CPUCores, metric.GAUGE))
			checkError(systemSampleMetricSet.SetMetric("hypervisorTotalMHz", TotalMHz, metric.GAUGE))

			cpuPercent := float64(0)
			if cpuAllocationLimit > TotalMHz || cpuAllocationLimit < 0 {
				cpuPercent = (float64(vm.Summary.QuickStats.OverallCpuUsage) / TotalMHz) * 100
			} else {
				cpuPercent = (float64(vm.Summary.QuickStats.OverallCpuUsage) / cpuAllocationLimit) * 100
			}

			checkError(systemSampleMetricSet.SetMetric("cpuPercent", cpuPercent, metric.GAUGE))
			checkError(systemSampleMetricSet.SetMetric("hypervisorConfigName", discoveredHost.Name, metric.ATTRIBUTE))

		}

		// disk
		checkError(systemSampleMetricSet.SetMetric("diskTotalBytes", vm.Summary.Storage.Committed, metric.GAUGE))

		// network
		checkError(systemSampleMetricSet.SetMetric("ipAddress", vm.Guest.IpAddress, metric.ATTRIBUTE))

		// vm state
		checkError(systemSampleMetricSet.SetMetric("connectionState", fmt.Sprintf("%v", vm.Runtime.ConnectionState), metric.ATTRIBUTE))
		checkError(systemSampleMetricSet.SetMetric("powerState", fmt.Sprintf("%v", vm.Runtime.PowerState), metric.ATTRIBUTE))

	}
}
