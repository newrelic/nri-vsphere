// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package process

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/nri-vsphere/internal/config"
	"github.com/vmware/govmomi/vim25/mo"
)

func createVirtualMachineSamples(config *config.Config) {
	for _, dc := range config.Datacenters {
		for _, vm := range dc.VirtualMachines {

			// filtering here will to avoid sending data to backend
			if config.TagFilteringEnabled() && !config.TagCollector.MatchObjectTags(vm.Self) {
				continue
			}

			// The virtual machine configuration is not guaranteed to be available. For example, the configuration
			// information would be unavailable if the server is unable to access the virtual machine files on disk,
			// and is often also unavailable during the initial phases of virtual machine creation.
			if vm.Config == nil {
				config.Logrus.WithField("vmMOR", vm.Self.String()).Debug("vm.Config is nil for this vm")
				continue
			}

			// resourcePool Returns null if the virtual machine is a template or the session has no access to the resource pool.
			if vm.ResourcePool == nil {
				config.Logrus.WithField("vmName", vm.Config.Name).Debug("vm.ResourcePool is nil for this vm")
				continue
			}

			// This property is null if the virtual machine is not running and is not assigned to run on a particular host.
			if vm.Summary.Runtime.Host == nil {
				config.Logrus.WithField("vmName", vm.Config.Name).Debug("vm.Summary.Runtime.Host is nil for this vm")
				continue
			}

			// we need the host and it's parent
			if h, ok := dc.Hosts[vm.Summary.Runtime.Host.Reference()]; !ok || ok && h.Parent == nil {
				config.Logrus.WithField("vmName", vm.Config.Name).Debug("host not found for this vm")
				continue
			}

			vmHost := dc.Hosts[*vm.Summary.Runtime.Host]
			vmHostParent := *vmHost.Parent
			hostConfigName := vmHost.Summary.Config.Name
			vmConfigName := vm.Summary.Config.Name
			datacenterName := dc.Datacenter.Name

			entityName := hostConfigName + ":" + vmConfigName
			if c, ok := dc.Clusters[vmHostParent]; ok {
				entityName = c.Name + ":" + entityName
			}
			entityName = sanitizeEntityName(config, entityName, datacenterName)

			// Unique identifier for the vm entity
			instanceUuid := vm.Config.InstanceUuid

			e, ms, err := createNewEntityWithMetricSet(config, entityTypeVm, entityName, instanceUuid)
			if err != nil {
				config.Logrus.WithError(err).
					WithField("vmName", entityName).
					WithField("instanceUuid", instanceUuid).
					Error("failed to create metricSet")
				continue
			}

			checkError(config.Logrus, ms.SetMetric("overallStatus", string(vm.OverallStatus), metric.ATTRIBUTE))

			if config.Args.DatacenterLocation != "" {
				checkError(config.Logrus, ms.SetMetric("datacenterLocation", config.Args.DatacenterLocation, metric.ATTRIBUTE))
			}

			if config.IsVcenterAPIType {
				checkError(config.Logrus, ms.SetMetric("datacenterName", datacenterName, metric.ATTRIBUTE))
			}

			if c, ok := dc.Clusters[vmHostParent]; ok {
				checkError(config.Logrus, ms.SetMetric("clusterName", c.Name, metric.ATTRIBUTE))
			}

			checkError(config.Logrus, ms.SetMetric("hypervisorHostname", hostConfigName, metric.ATTRIBUTE))

			// vm
			checkError(config.Logrus, ms.SetMetric("vmConfigName", vmConfigName, metric.ATTRIBUTE))
			checkError(config.Logrus, ms.SetMetric("instanceUuid", instanceUuid, metric.ATTRIBUTE))

			// not available if VM is offline
			if vm.Summary.Guest != nil {
				checkError(config.Logrus, ms.SetMetric("vmHostname", vm.Summary.Guest.HostName, metric.ATTRIBUTE))
			}

			vmResourcePool := *vm.ResourcePool
			if resourcePool, ok := dc.GetResourcePool(vmResourcePool); ok {
				checkError(config.Logrus, ms.SetMetric("resourcePoolName", resourcePool.Name, metric.ATTRIBUTE))
			}

			datastoreList := ""
			for _, ds := range vm.Datastore {
				if d, ok := dc.Datastores[ds]; ok {
					datastoreList += d.Name + "|"
				}
			}
			datastoreList = strings.TrimSuffix(datastoreList, "|")
			checkError(config.Logrus, ms.SetMetric("datastoreNameList", datastoreList, metric.ATTRIBUTE))

			networkList := ""
			for _, nw := range vm.Network {
				if n, ok := dc.Networks[nw]; ok {
					networkList += n.Name + "|"
				}
			}
			networkList = strings.TrimSuffix(networkList, "|")
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
				var ipAddresses strings.Builder
				for _, nic := range vm.Guest.Net {
					// available in api v5
					if nic.IpConfig != nil {
						for _, addr := range nic.IpConfig.IpAddress {
							ipAddresses.WriteString(addr.IpAddress)
							ipAddresses.WriteRune('|')
						}
					} else {
						for _, ip := range nic.IpAddress {
							ipAddresses.WriteString(ip)
							ipAddresses.WriteRune('|')
						}
					}
				}

				if fqdn := computeFullHostname(vm); fqdn != "" {
					checkError(config.Logrus, ms.SetMetric("vmFullHostname", fqdn, metric.ATTRIBUTE))
				}

				ipAddressesTrimmed := strings.TrimSuffix(ipAddresses.String(), "|")
				// it might be empty, but we still add the attribute for consistency
				checkError(config.Logrus, ms.SetMetric("ipAddresses", ipAddressesTrimmed, metric.ATTRIBUTE))
			}

			// vm state
			checkError(config.Logrus, ms.SetMetric("connectionState", fmt.Sprintf("%v", vm.Runtime.ConnectionState), metric.ATTRIBUTE))
			checkError(config.Logrus, ms.SetMetric("powerState", fmt.Sprintf("%v", vm.Runtime.PowerState), metric.ATTRIBUTE))

			// Tags
			if config.TagCollectionEnabled() {
				tagsByCategory := config.TagCollector.GetTagsByCategories(vm.Self)
				for k, v := range tagsByCategory {
					checkError(config.Logrus, ms.SetMetric(tagsPrefix+k, v, metric.ATTRIBUTE))
					// add tags to inventory due to the inventory workaround
					addTagsToInventory(config, e, k, v)
				}
			}

			// Performance metrics
			if config.PerfMetricsCollectionEnabled() {
				perfMetrics := dc.GetPerfMetrics(vm.Self)
				for _, perfMetric := range perfMetrics {
					checkError(config.Logrus, ms.SetMetric(perfMetricPrefix+perfMetric.Counter, perfMetric.Value, metric.GAUGE))
				}
			}

			// Snapshots
			if vm.Snapshot != nil && vm.LayoutEx != nil && config.Args.EnableVsphereSnapshots {
				sp := newSnapshotProcessor(config.Logrus, vm)
				sp.processSnapshotTree(nil, vm.Snapshot.RootSnapshotList)
				sp.createSnapshotSamples(e, entityName, vm.Snapshot.RootSnapshotList)
			}

			// suspendMemory
			if vm.LayoutEx != nil {
				var suspendMemory, suspendMemoryUnique int64
				for _, exFile := range vm.LayoutEx.File {
					if exFile.Type == "suspendMemory" {
						suspendMemory += exFile.Size
						suspendMemoryUnique += exFile.UniqueSize
					}
				}
				checkError(config.Logrus, ms.SetMetric("disk.suspendMemory", strconv.FormatInt(suspendMemory, 10), metric.GAUGE))
				checkError(config.Logrus, ms.SetMetric("disk.suspendMemoryUnique", strconv.FormatInt(suspendMemoryUnique, 10), metric.GAUGE))
			}
		}
	}
}

// computeFullHostname joins hostname and domain for each VM
// These data depends on the vmwareTools, therefore they need to be installed,
// and we depend on how such tool is collecting the value.
// Moreover, notice that we are returning the first domain contained in IpStack array.
func computeFullHostname(vm *mo.VirtualMachine) string {
	if vm.Guest == nil {
		return ""
	}
	if vm.Summary.Guest == nil {
		return ""
	}

	for _, is := range vm.Guest.IpStack {
		if is.DnsConfig == nil {
			continue
		}

		var domain = is.DnsConfig.DomainName
		var hostname = is.DnsConfig.HostName
		if domain == "" || hostname == "" {
			continue
		}
		if hostname != vm.Summary.Guest.HostName {
			continue
		}

		// we noticed that the hostname is sometimes the short one and sometimes the fqdn
		hostname = strings.TrimSuffix(hostname, domain)
		hostname = strings.TrimSuffix(hostname, ".")

		var fullHostname = hostname + "." + domain
		fullHostname = strings.TrimSuffix(fullHostname, ".")

		return fullHostname
	}
	return ""
}
