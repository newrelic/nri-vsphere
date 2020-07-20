// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package process

import (
	"fmt"
	"github.com/newrelic/nri-vsphere/internal/model/tag"
	"strconv"

	"github.com/newrelic/nri-vsphere/internal/config"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
)

func createVirtualMachineSamples(config *config.Config) {
	for _, dc := range config.Datacenters {
		for _, vm := range dc.VirtualMachines {

			if vm.Config == nil {
				continue // The virtual machine configuration is not guaranteed to be available. For example, the configuration information would be unavailable if the server is unable to access the virtual machine files on disk, and is often also unavailable during the initial phases of virtual machine creation.
			}
			if vm.ResourcePool == nil {
				continue // resourcePool Returns null if the virtual machine is a template or the session has no access to the resource pool.
			}
			if vm.Summary.Runtime.Host == nil {
				continue // This property is null if the virtual machine is not running and is not assigned to run on a particular host.
			}
			if _, ok := dc.Hosts[*vm.Summary.Runtime.Host]; !ok {
				continue
			}
			if dc.Hosts[*vm.Summary.Runtime.Host].Parent == nil {
				continue
			}
			vmHost := dc.Hosts[*vm.Summary.Runtime.Host]
			vmHostParent := *vmHost.Parent
			vmResourcePool := *vm.ResourcePool
			hostConfigName := vmHost.Summary.Config.Name
			vmConfigName := vm.Summary.Config.Name
			datacenterName := dc.Datacenter.Name
			entityName := hostConfigName + ":" + vmConfigName

			if cluster, ok := dc.Clusters[vmHostParent]; ok {
				entityName = cluster.Name + ":" + entityName
			}

			entityName = sanitizeEntityName(config, entityName, datacenterName)

			// Unique identifier for the vm entity
			instanceUuid := vm.Config.InstanceUuid

			e, ms, err := createNewEntityWithMetricSet(config, entityTypeVm, entityName, instanceUuid)
			if err != nil {
				config.Logrus.WithError(err).WithField("vmName", entityName).WithField("instanceUuid", instanceUuid).Error("failed to create metricSet")
				continue
			}

			checkError(config.Logrus, ms.SetMetric("overallStatus", string(vm.OverallStatus), metric.ATTRIBUTE))

			if config.Args.DatacenterLocation != "" {
				checkError(config.Logrus, ms.SetMetric("datacenterLocation", config.Args.DatacenterLocation, metric.ATTRIBUTE))
			}

			if cluster, ok := dc.Clusters[vmHostParent]; ok {
				checkError(config.Logrus, ms.SetMetric("clusterName", cluster.Name, metric.ATTRIBUTE))
			}

			if config.IsVcenterAPIType {
				checkError(config.Logrus, ms.SetMetric("datacenterName", datacenterName, metric.ATTRIBUTE))
			}
			checkError(config.Logrus, ms.SetMetric("hypervisorHostname", hostConfigName, metric.ATTRIBUTE))

			resourcePoolName := dc.GetResourcePoolName(vmResourcePool)
			checkError(config.Logrus, ms.SetMetric("resourcePoolName", resourcePoolName, metric.ATTRIBUTE))

			datastoreList := ""
			for _, ds := range vm.Datastore {
				datastoreList += dc.Datastores[ds].Name + "|"
			}
			checkError(config.Logrus, ms.SetMetric("datastoreNameList", datastoreList, metric.ATTRIBUTE))
			// vm
			// not available if VM is offline
			if vm.Summary.Guest != nil {
				checkError(config.Logrus, ms.SetMetric("vmHostname", vm.Summary.Guest.HostName, metric.ATTRIBUTE))
			}
			checkError(config.Logrus, ms.SetMetric("vmConfigName", vmConfigName, metric.ATTRIBUTE))
			checkError(config.Logrus, ms.SetMetric("instanceUuid", instanceUuid, metric.ATTRIBUTE))

			networkList := ""
			for _, nw := range vm.Network {
				networkList += dc.Networks[nw].Name + "|"
			}
			checkError(config.Logrus, ms.SetMetric("networkNameList", networkList, metric.ATTRIBUTE))

			operatingSystem := determineOS(vm.Summary.Config.GuestFullName)
			checkError(config.Logrus, ms.SetMetric("operatingSystem", operatingSystem, metric.ATTRIBUTE))
			checkError(config.Logrus, ms.SetMetric("guestFullName", vm.Summary.Config.GuestFullName, metric.ATTRIBUTE))

			// SystemSample metrics

			// memory
			memorySize := vm.Summary.Config.MemorySizeMB
			checkError(config.Logrus, ms.SetMetric("mem.size", memorySize, metric.GAUGE))
			memoryUsed := vm.Summary.QuickStats.GuestMemoryUsage
			checkError(config.Logrus, ms.SetMetric("mem.usage", memoryUsed, metric.GAUGE))
			memoryFree := memorySize - memoryUsed
			checkError(config.Logrus, ms.SetMetric("mem.free", memoryFree, metric.GAUGE))
			checkError(config.Logrus, ms.SetMetric("mem.hostUsage", vm.Summary.QuickStats.HostMemoryUsage, metric.GAUGE))

			checkError(config.Logrus, ms.SetMetric("mem.balloned", vm.Summary.QuickStats.BalloonedMemory, metric.GAUGE))
			checkError(config.Logrus, ms.SetMetric("mem.swapped", vm.Summary.QuickStats.SwappedMemory, metric.GAUGE))
			swappedSsd := float64(vm.Summary.QuickStats.SsdSwappedMemory) / (1 << 10)
			checkError(config.Logrus, ms.SetMetric("mem.swappedSsd", swappedSsd, metric.GAUGE))

			// cpu
			checkError(config.Logrus, ms.SetMetric("cpu.cores", vm.Summary.Config.NumCpu, metric.GAUGE))
			checkError(config.Logrus, ms.SetMetric("cpu.overallUsage", vm.Summary.QuickStats.OverallCpuUsage, metric.GAUGE))

			var cpuAllocationLimit float64
			if vm.Config.CpuAllocation != nil {
				if vm.Config.CpuAllocation.Limit != nil {
					cpuAllocationLimit = float64(*vm.Config.CpuAllocation.Limit)
				}
			}
			checkError(config.Logrus, ms.SetMetric("cpu.allocationLimit", cpuAllocationLimit, metric.GAUGE))

			if vmHost.Summary.Hardware != nil {
				CPUMhz := vmHost.Summary.Hardware.CpuMhz
				CPUCores := vmHost.Summary.Hardware.NumCpuCores
				OverallCpuUsage := vm.Summary.QuickStats.OverallCpuUsage
				var cpuPercent float64

				TotalMHz := float64(CPUMhz) * float64(CPUCores)
				if (cpuAllocationLimit > TotalMHz || cpuAllocationLimit < 0) && TotalMHz != 0 {
					cpuPercent = float64(OverallCpuUsage) / TotalMHz * 100
				} else if cpuAllocationLimit != 0 {
					cpuPercent = float64(OverallCpuUsage) / cpuAllocationLimit * 100
				}
				checkError(config.Logrus, ms.SetMetric("cpu.hostUsagePercent", cpuPercent, metric.GAUGE))
			}

			// disk
			if vm.Summary.Storage != nil {
				checkError(config.Logrus, ms.SetMetric("disk.totalUncommittedMiB", vm.Summary.Storage.Uncommitted/(1<<20), metric.GAUGE))
				checkError(config.Logrus, ms.SetMetric("disk.totalMiB", vm.Summary.Storage.Committed/(1<<20), metric.GAUGE))
				checkError(config.Logrus, ms.SetMetric("disk.totalUnsharedMiB", vm.Summary.Storage.Unshared/(1<<20), metric.GAUGE))
			}

			// network
			if vm.Guest != nil {
				checkError(config.Logrus, ms.SetMetric("ipAddress", vm.Guest.IpAddress, metric.ATTRIBUTE))
			}
			// vm state
			checkError(config.Logrus, ms.SetMetric("connectionState", fmt.Sprintf("%v", vm.Runtime.ConnectionState), metric.ATTRIBUTE))
			checkError(config.Logrus, ms.SetMetric("powerState", fmt.Sprintf("%v", vm.Runtime.PowerState), metric.ATTRIBUTE))

			//Tags
			tagsByCategory := tag.GetTagsByCategories(vm.Self)
			for k, v := range tagsByCategory {
				checkError(config.Logrus, ms.SetMetric(tagsPrefix+k, v, metric.ATTRIBUTE))
				// add tags to inventory due to the inventory workaround
				checkError(config.Logrus, e.SetInventoryItem("tags", tagsPrefix+k, v))
			}
			// Performance metrics
			perfMetrics := dc.GetPerfMetrics(vm.Self)
			for _, perfMetric := range perfMetrics {
				checkError(config.Logrus, ms.SetMetric(perfMetricPrefix+perfMetric.Counter, perfMetric.Value, metric.GAUGE))
			}

			if vm.Snapshot != nil && config.Args.EnableVsphereSnapshots {
				infoSnapshot, suspendMemory, suspendMemoryUnique := processLayoutEx(vm.LayoutEx)
				checkError(config.Logrus, ms.SetMetric("disk.suspendMemory", strconv.FormatInt(suspendMemory, 10), metric.GAUGE))
				checkError(config.Logrus, ms.SetMetric("disk.suspendMemoryUnique", strconv.FormatInt(suspendMemoryUnique, 10), metric.GAUGE))

				for _, t := range vm.Snapshot.RootSnapshotList {
					traverseSnapshotList(e, config, t, entityName, infoSnapshot)
				}
			}
		}
	}
}
