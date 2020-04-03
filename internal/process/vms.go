package process

import (
	"fmt"
	"strings"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/nri-vmware-vsphere/internal/load"
)

func createVirtualMachineSamples(config *load.Config, timestamp int64) {
	for _, dc := range config.Datacenters {
		for _, vm := range dc.VirtualMachines {
			// // resolve hypervisor host
			vmHost := dc.Hosts[*vm.Summary.Runtime.Host]
			hostConfigName := vmHost.Summary.Config.Name
			vmConfigName := vm.Summary.Config.Name
			// hostSystem := object.NewHostSystem(load.HostSystemContainerView.Client(), *vmHost)
			// hypervisorHost, err := hostSystem.ObjectName(ctx)
			// vmHost := findHost(hostSystem.Reference())
			datacenterName := dc.Datacenter.Name
			clusterName := dc.Clusters[*vmHost.Parent].Name

			entityName := hostConfigName + ":" + vmConfigName + ":vm"
			if config.IsVcenterAPIType {
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

			// Unique identifier for the vm entity
			// uuid := integration.IDAttribute{Key: "uuid", Value: vm.Config.instanceUuid}
			workingEntity, err := config.Integration.Entity(entityName, "vsphere")
			if err != nil {
				config.Logrus.WithError(err).Error("failed to create entity")
			}

			workingEntity.SetInventoryItem("name", "value", fmt.Sprintf("%v:%d", entityName, timestamp))

			ms := workingEntity.NewMetricSet("VSphereVmSample")

			if config.Args.DatacenterLocation != "" {
				checkError(config, ms.SetMetric("datacenterLocation", config.Args.DatacenterLocation, metric.ATTRIBUTE))
			}
			if config.IsVcenterAPIType {
				checkError(config, ms.SetMetric("datacenterName", datacenterName, metric.ATTRIBUTE))
				checkError(config, ms.SetMetric("clusterName", clusterName, metric.ATTRIBUTE))
			}
			checkError(config, ms.SetMetric("hypervisorHostname", hostConfigName, metric.ATTRIBUTE))

			resourcePoolName := dc.ResourcePools[*vm.ResourcePool].Name
			checkError(config, ms.SetMetric("resourcePoolName", resourcePoolName, metric.ATTRIBUTE))

			datastoreList := ""
			for _, ds := range vm.Datastore {
				datastoreList += dc.Datastores[ds].Name + "|"
			}
			checkError(config, ms.SetMetric("datastoreNameList", datastoreList, metric.ATTRIBUTE))
			// vm
			// not available if VM is offline
			vmHostname := vm.Summary.Guest.HostName
			checkError(config, ms.SetMetric("vmConfigName", vmConfigName, metric.ATTRIBUTE))
			checkError(config, ms.SetMetric("vmHostname", vmHostname, metric.ATTRIBUTE))
			checkError(config, ms.SetMetric("instanceUuid", vm.Config.InstanceUuid, metric.ATTRIBUTE))

			networkList := ""
			for _, nw := range vm.Network {
				networkList += dc.Networks[nw].Name + "|"
			}
			checkError(config, ms.SetMetric("networkNameList", networkList, metric.ATTRIBUTE))

			operatingSystem := determineOS(vm.Summary.Config.GuestFullName)
			checkError(config, ms.SetMetric("operatingSystem", operatingSystem, metric.ATTRIBUTE))
			checkError(config, ms.SetMetric("guestFullName", vm.Summary.Config.GuestFullName, metric.ATTRIBUTE))

			// SystemSample metrics

			// memory
			memoryTotalBytes := float64(vm.Summary.Config.MemorySizeMB) * 1e+6
			checkError(config, ms.SetMetric("memoryTotalBytes", memoryTotalBytes, metric.GAUGE))
			checkError(config, ms.SetMetric("systemMemoryBytes", memoryTotalBytes, metric.GAUGE))
			memoryUsedBytes := float64(vm.Summary.QuickStats.GuestMemoryUsage) * 1e+6
			memoryFreeBytes := memoryTotalBytes - memoryUsedBytes
			checkError(config, ms.SetMetric("memoryUsedBytes", memoryUsedBytes, metric.GAUGE))
			checkError(config, ms.SetMetric("memoryFreeBytes", memoryFreeBytes, metric.GAUGE))

			// cpu
			checkError(config, ms.SetMetric("coreCount", vm.Summary.Config.NumCpu, metric.GAUGE))
			checkError(config, ms.SetMetric("overallCpuUsageMHz", vm.Summary.QuickStats.OverallCpuUsage, metric.GAUGE))

			cpuAllocationLimit := float64(0)
			if vm.Config.CpuAllocation.Limit != nil {
				cpuAllocationLimit = float64(*vm.Config.CpuAllocation.Limit)
			}

			checkError(config, ms.SetMetric("cpuAllocationLimit", cpuAllocationLimit, metric.GAUGE))

			if vmHost.Self.Value != "" {
				CPUMhz := vmHost.Summary.Hardware.CpuMhz
				CPUCores := vmHost.Summary.Hardware.NumCpuCores
				CPUThreads := vmHost.Summary.Hardware.NumCpuThreads
				TotalMHz := float64(CPUMhz) * float64(CPUCores)
				checkError(config, ms.SetMetric("hypervisorCpuThreads", CPUThreads, metric.GAUGE))
				checkError(config, ms.SetMetric("hypervisorCpuMhz", CPUMhz, metric.GAUGE))
				checkError(config, ms.SetMetric("hypervisorCpuCores", CPUCores, metric.GAUGE))
				checkError(config, ms.SetMetric("hypervisorTotalMHz", TotalMHz, metric.GAUGE))

				cpuPercent := float64(0)
				if cpuAllocationLimit > TotalMHz || cpuAllocationLimit < 0 {
					cpuPercent = (float64(vm.Summary.QuickStats.OverallCpuUsage) / TotalMHz) * 100
				} else {
					cpuPercent = (float64(vm.Summary.QuickStats.OverallCpuUsage) / cpuAllocationLimit) * 100
				}

				checkError(config, ms.SetMetric("cpuPercent", cpuPercent, metric.GAUGE))
				checkError(config, ms.SetMetric("hypervisorConfigName", vmHost.Name, metric.ATTRIBUTE))

			}

			// disk
			checkError(config, ms.SetMetric("diskTotalBytes", vm.Summary.Storage.Committed, metric.GAUGE))

			// network
			checkError(config, ms.SetMetric("ipAddress", vm.Guest.IpAddress, metric.ATTRIBUTE))

			// vm state
			checkError(config, ms.SetMetric("connectionState", fmt.Sprintf("%v", vm.Runtime.ConnectionState), metric.ATTRIBUTE))
			checkError(config, ms.SetMetric("powerState", fmt.Sprintf("%v", vm.Runtime.PowerState), metric.ATTRIBUTE))

		}
	}
}
