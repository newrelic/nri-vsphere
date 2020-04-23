// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package process

import (
	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/nri-vmware-vsphere/internal/load"
	"strconv"
)

func createHostSamples(config *load.Config, timestamp int64) {
	for _, dc := range config.Datacenters {
		for _, host := range dc.Hosts {
			hostConfigName := host.Summary.Config.Name
			entityName := hostConfigName
			datacenterName := dc.Datacenter.Name

			if cluster, ok := dc.Clusters[host.Parent.Reference()]; ok {
				entityName = cluster.Name + ":" + entityName
			}

			entityName = sanitizeEntityName(config, entityName, datacenterName)

			// bios uuid identifies the host unequivocally and is available from vcenter/host api
			uuid := host.Summary.Hardware.Uuid

			ms := createNewEntityWithMetricSet(config, entityTypeHost, entityName, uuid)

			if cluster, ok := dc.Clusters[host.Parent.Reference()]; ok {
				checkError(config, ms.SetMetric("clusterName", cluster.Name, metric.ATTRIBUTE))
			}
			checkError(config, ms.SetMetric("overallStatus", string(host.OverallStatus), metric.ATTRIBUTE))

			if config.IsVcenterAPIType {
				checkError(config, ms.SetMetric("datacenterName", datacenterName, metric.ATTRIBUTE))
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

			if host.Runtime.InQuarantineMode != nil {
				checkError(config, ms.SetMetric("inQuarantineMode", strconv.FormatBool(*host.Runtime.InQuarantineMode), metric.ATTRIBUTE))
			}
			if host.Runtime.BootTime != nil {
				checkError(config, ms.SetMetric("bootTime", host.Runtime.BootTime.String(), metric.ATTRIBUTE))
			}
			checkError(config, ms.SetMetric("connectionState", string(host.Runtime.ConnectionState), metric.ATTRIBUTE))
			checkError(config, ms.SetMetric("inMaintenanceMode", strconv.FormatBool(host.Runtime.InMaintenanceMode), metric.ATTRIBUTE))
			checkError(config, ms.SetMetric("powerState", string(host.Runtime.PowerState), metric.ATTRIBUTE))
			checkError(config, ms.SetMetric("standbyMode", host.Runtime.StandbyMode, metric.ATTRIBUTE))
			checkError(config, ms.SetMetric("cryptoState", host.Runtime.CryptoState, metric.ATTRIBUTE))

			networkList := ""
			for _, nw := range host.Network {
				networkList += dc.Networks[nw].Name + "|"
			}
			checkError(config, ms.SetMetric("networkNameList", networkList, metric.ATTRIBUTE))

			// memory
			memoryTotal := host.Summary.Hardware.MemorySize / (1 << 20)
			checkError(config, ms.SetMetric("mem.size", memoryTotal, metric.GAUGE))

			memoryUsed := host.Summary.QuickStats.OverallMemoryUsage
			checkError(config, ms.SetMetric("mem.usage", memoryUsed, metric.GAUGE))

			memoryFree := int32(memoryTotal) - memoryUsed
			checkError(config, ms.SetMetric("mem.free", memoryFree, metric.GAUGE))

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

			CPUAvailable := TotalMHz - float64(host.Summary.QuickStats.OverallCpuUsage)
			checkError(config, ms.SetMetric("cpu.available", CPUAvailable, metric.GAUGE))

			// disk
			diskTotalMiB := int64(0)
			if host.Config != nil {
				if host.Config.FileSystemVolume != nil {
					for _, mount := range host.Config.FileSystemVolume.MountInfo {
						capacity := mount.Volume.GetHostFileSystemVolume().Capacity
						diskTotalMiB += capacity / (1 << 20)
					}
				}
			}
			checkError(config, ms.SetMetric("disk.totalMiB", diskTotalMiB, metric.GAUGE))

		}
	}
}
