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
			datacenterName := dc.Datacenter.Name
			entityName := hostConfigName + ":" + vmConfigName + ":vm"

			if cluster, ok := dc.Clusters[*vmHost.Parent]; ok {
				entityName =  cluster.Name + ":" + entityName
			}
			if config.IsVcenterAPIType {
				entityName = datacenterName + ":" + entityName
			}

			if config.Args.DatacenterLocation != "" {
				entityName = config.Args.DatacenterLocation + ":" + entityName
			}

			entityName = strings.ToLower(entityName)
			entityName = strings.ReplaceAll(entityName, ".", "-")

			// Unique identifier for the vm entity
			instanceUuid := vm.Config.InstanceUuid
			workingEntity, err := config.Integration.Entity(instanceUuid, "vsphere-vm")
			if err != nil {
				config.Logrus.WithError(err).Error("failed to create entity")
			}

			// entity displayName
			workingEntity.SetInventoryItem("vsphereVm", "name", entityName)

			ms := workingEntity.NewMetricSet("VSphereVmSample")

			if config.Args.DatacenterLocation != "" {
				checkError(config, ms.SetMetric("datacenterLocation", config.Args.DatacenterLocation, metric.ATTRIBUTE))
			}

			if cluster, ok := dc.Clusters[*vmHost.Parent]; ok {
				checkError(config, ms.SetMetric("clusterName", cluster.Name, metric.ATTRIBUTE))
			}

			if config.IsVcenterAPIType {
				checkError(config, ms.SetMetric("datacenterName", datacenterName, metric.ATTRIBUTE))
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
			checkError(config, ms.SetMetric("instanceUuid", instanceUuid, metric.ATTRIBUTE))

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
			memorySize := vm.Summary.Config.MemorySizeMB
			checkError(config, ms.SetMetric("mem.size", memorySize, metric.GAUGE))
			memoryUsed := vm.Summary.QuickStats.GuestMemoryUsage
			checkError(config, ms.SetMetric("mem.usage", memoryUsed, metric.GAUGE))
			memoryFree := memorySize - memoryUsed
			checkError(config, ms.SetMetric("mem.free", memoryFree, metric.GAUGE))
			checkError(config, ms.SetMetric("mem.balloned", vm.Summary.QuickStats.BalloonedMemory, metric.GAUGE))
			checkError(config, ms.SetMetric("mem.swapped", vm.Summary.QuickStats.SwappedMemory, metric.GAUGE))
			swappedSsd := float64(vm.Summary.QuickStats.SsdSwappedMemory) / 1e+3
			checkError(config, ms.SetMetric("mem.swappedSsd", swappedSsd, metric.GAUGE))

			// cpu
			checkError(config, ms.SetMetric("cpu.cores", vm.Summary.Config.NumCpu, metric.GAUGE))
			checkError(config, ms.SetMetric("cpu.overallUsage", vm.Summary.QuickStats.OverallCpuUsage, metric.GAUGE))

			cpuAllocationLimit := float64(0)
			if vm.Config.CpuAllocation.Limit != nil {
				cpuAllocationLimit = float64(*vm.Config.CpuAllocation.Limit)
			}
			checkError(config, ms.SetMetric("cpu.allocationLimit", cpuAllocationLimit, metric.GAUGE))

			CPUMhz := vmHost.Summary.Hardware.CpuMhz
			CPUCores := vmHost.Summary.Hardware.NumCpuCores
			TotalMHz := float64(CPUMhz) * float64(CPUCores)

			cpuPercent := float64(0)
			if cpuAllocationLimit > TotalMHz || cpuAllocationLimit < 0 {
				cpuPercent = (float64(vm.Summary.QuickStats.OverallCpuUsage) / TotalMHz) * 100
			} else {
				cpuPercent = (float64(vm.Summary.QuickStats.OverallCpuUsage) / cpuAllocationLimit) * 100
			}
			checkError(config, ms.SetMetric("cpu.hostUsagePercent", cpuPercent, metric.GAUGE))

			// disk
			diskTotal := vm.Summary.Storage.Committed / 1e+6
			checkError(config, ms.SetMetric("disk.totalMB", diskTotal, metric.GAUGE))

			// network
			checkError(config, ms.SetMetric("ipAddress", vm.Guest.IpAddress, metric.ATTRIBUTE))

			// vm state
			checkError(config, ms.SetMetric("connectionState", fmt.Sprintf("%v", vm.Runtime.ConnectionState), metric.ATTRIBUTE))
			checkError(config, ms.SetMetric("powerState", fmt.Sprintf("%v", vm.Runtime.PowerState), metric.ATTRIBUTE))

		}
	}
}
