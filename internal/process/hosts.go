// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package process

import (
	"strconv"
	"strings"

	"github.com/newrelic/nri-vsphere/internal/config"

	"github.com/newrelic/infra-integrations-sdk/v3/data/metric"
)

func createHostSamples(config *config.Config) {
	for _, dc := range config.Datacenters {
		for _, host := range dc.Hosts {

			// filtering here will to avoid sending data to backend
			if config.TagFilteringEnabled() && !config.TagCollector.MatchObjectTags(host.Self) {
				continue
			}

			if host.Summary.Hardware == nil {
				config.Logrus.WithField("hostMOR", host.Self.String()).Debug("host.Summary.Hardware is nil for this host")
				continue
			}
			// bios uuid identifies the host unequivocally and is available from vcenter/host api
			uuid := host.Summary.Hardware.Uuid

			hostConfigName := host.Summary.Config.Name
			entityName := hostConfigName
			datacenterName := dc.Datacenter.Name

			if cluster, ok := dc.Clusters[host.Parent.Reference()]; ok {
				entityName = cluster.Name + ":" + entityName
			}

			entityName = sanitizeEntityName(config, entityName, datacenterName)

			e, ms, err := createNewEntityWithMetricSet(config, entityTypeHost, entityName, uuid)
			if err != nil {
				config.Logrus.WithError(err).WithField("hostName", entityName).WithField("uuid", uuid).Error("failed to create metricSet")
				continue
			}

			if config.IsVcenterAPIType {
				checkError(config.Logrus, ms.SetMetric("datacenterName", datacenterName, metric.ATTRIBUTE))
			}

			if cluster, ok := dc.Clusters[host.Parent.Reference()]; ok {
				checkError(config.Logrus, ms.SetMetric("clusterName", cluster.Name, metric.ATTRIBUTE))
			}

			checkError(config.Logrus, ms.SetMetric("overallStatus", string(host.OverallStatus), metric.ATTRIBUTE))

			if config.Args.DatacenterLocation != "" {
				checkError(config.Logrus, ms.SetMetric("datacenterLocation", config.Args.DatacenterLocation, metric.ATTRIBUTE))
			}
			checkError(config.Logrus, ms.SetMetric("hypervisorHostname", hostConfigName, metric.ATTRIBUTE))

			checkError(config.Logrus, ms.SetMetric("vmCount", len(host.Vm), metric.GAUGE))

			if host.Runtime.InQuarantineMode != nil {
				checkError(config.Logrus, ms.SetMetric("inQuarantineMode", strconv.FormatBool(*host.Runtime.InQuarantineMode), metric.ATTRIBUTE))
			}

			if host.Runtime.BootTime != nil {
				checkError(config.Logrus, ms.SetMetric("bootTime", host.Runtime.BootTime.String(), metric.ATTRIBUTE))
			}

			checkError(config.Logrus, ms.SetMetric("connectionState", string(host.Runtime.ConnectionState), metric.ATTRIBUTE))
			checkError(config.Logrus, ms.SetMetric("inMaintenanceMode", strconv.FormatBool(host.Runtime.InMaintenanceMode), metric.ATTRIBUTE))
			checkError(config.Logrus, ms.SetMetric("powerState", string(host.Runtime.PowerState), metric.ATTRIBUTE))
			checkError(config.Logrus, ms.SetMetric("standbyMode", host.Runtime.StandbyMode, metric.ATTRIBUTE))
			checkError(config.Logrus, ms.SetMetric("cryptoState", host.Runtime.CryptoState, metric.ATTRIBUTE))

			resourcePools := dc.FindResourcePools(host.Parent.Reference())
			resourcePoolList := ""
			for _, rp := range resourcePools {
				resourcePoolList += rp.Name + "|"
			}
			resourcePoolList = strings.TrimSuffix(resourcePoolList, "|")
			checkError(config.Logrus, ms.SetMetric("resourcePoolNameList", resourcePoolList, metric.ATTRIBUTE))

			datastoreList := ""
			for _, ds := range host.Datastore {
				if d, ok := dc.Datastores[ds]; ok {
					datastoreList += d.Name + "|"
				}
			}
			datastoreList = strings.TrimSuffix(datastoreList, "|")
			checkError(config.Logrus, ms.SetMetric("datastoreNameList", datastoreList, metric.ATTRIBUTE))

			networkList := ""
			for _, nw := range host.Network {
				if n, ok := dc.Networks[nw]; ok {
					networkList += n.Name + "|"
				}
			}
			networkList = strings.TrimSuffix(networkList, "|")
			checkError(config.Logrus, ms.SetMetric("networkNameList", networkList, metric.ATTRIBUTE))

			checkError(config.Logrus, ms.SetMetric("uuid", host.Summary.Hardware.Uuid, metric.ATTRIBUTE))

			// memory
			memoryTotal := host.Summary.Hardware.MemorySize / (1 << 20)
			checkError(config.Logrus, ms.SetMetric("mem.size", memoryTotal, metric.GAUGE))

			memoryUsed := host.Summary.QuickStats.OverallMemoryUsage
			checkError(config.Logrus, ms.SetMetric("mem.usage", memoryUsed, metric.GAUGE))

			memoryFree := int32(memoryTotal) - memoryUsed
			checkError(config.Logrus, ms.SetMetric("mem.free", memoryFree, metric.GAUGE))

			// cpu
			CPUCores := host.Summary.Hardware.NumCpuCores
			checkError(config.Logrus, ms.SetMetric("cpu.cores", CPUCores, metric.GAUGE))

			CPUThreads := host.Summary.Hardware.NumCpuThreads
			checkError(config.Logrus, ms.SetMetric("cpu.threads", CPUThreads, metric.GAUGE))

			CPUMhz := host.Summary.Hardware.CpuMhz
			checkError(config.Logrus, ms.SetMetric("cpu.coreMHz", CPUMhz, metric.GAUGE))

			TotalMHz := float64(CPUMhz) * float64(CPUCores)
			checkError(config.Logrus, ms.SetMetric("cpu.totalMHz", TotalMHz, metric.GAUGE))

			if TotalMHz != 0 {
				cpuPercent := (float64(host.Summary.QuickStats.OverallCpuUsage) / TotalMHz) * 100
				checkError(config.Logrus, ms.SetMetric("cpu.percent", cpuPercent, metric.GAUGE))
			}

			checkError(config.Logrus, ms.SetMetric("cpu.overallUsage", host.Summary.QuickStats.OverallCpuUsage, metric.GAUGE))

			CPUAvailable := TotalMHz - float64(host.Summary.QuickStats.OverallCpuUsage)
			checkError(config.Logrus, ms.SetMetric("cpu.available", CPUAvailable, metric.GAUGE))

			// disk
			diskTotalMiB := int64(0)
			if host.Config != nil {
				if host.Config.FileSystemVolume != nil {
					for _, mount := range host.Config.FileSystemVolume.MountInfo {
						hostFileSystemVolume := mount.Volume.GetHostFileSystemVolume()
						if hostFileSystemVolume != nil {
							diskTotalMiB += hostFileSystemVolume.Capacity / (1 << 20)
						}
					}
				}
			}
			checkError(config.Logrus, ms.SetMetric("disk.totalMiB", diskTotalMiB, metric.GAUGE))

			// Tags
			if config.TagCollectionEnabled() {
				tagsByCategory := config.TagCollector.GetTagsByCategories(host.Self)
				for k, v := range tagsByCategory {
					checkError(config.Logrus, ms.SetMetric(tagsPrefix+k, v, metric.ATTRIBUTE))
					// add tags to inventory due to the inventory workaround
					addTagsToInventory(config, e, k, v)
				}
			}
			// Performance metrics
			if config.PerfMetricsCollectionEnabled() {
				perfMetrics := dc.GetPerfMetrics(host.Self)
				for _, perfMetric := range perfMetrics {
					checkError(config.Logrus, ms.SetMetric(perfMetricPrefix+perfMetric.Counter, perfMetric.Value, metric.GAUGE))
				}
			}

		}
	}
}
