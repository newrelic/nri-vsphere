// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package process

import (
	"strconv"

	"github.com/newrelic/nri-vsphere/internal/config"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
)

func createClusterSamples(config *config.Config) {
	for _, dc := range config.Datacenters {
		for _, cluster := range dc.Clusters {
			// // resolve hypervisor host
			datacenterName := dc.Datacenter.Name

			//Retrieving the list of host belonging to the cluster
			hostList := ""
			for _, hr := range cluster.Host {
				if h := dc.GetHost(hr); h != nil {
					hostList += h.Summary.Config.Name + "|"
				}
			}

			//Retrieving the list of networks attached to the cluster
			networkList := ""
			for _, nr := range cluster.Network {
				if n := dc.GetNetwork(nr); n != nil {
					networkList += n.Name + "|"
				}
			}

			//Retrieving the list of datastores attached to the cluster
			datastoreList := ""
			for _, dr := range cluster.Datastore {
				if ds := dc.GetDatastore(dr); ds != nil {
					datastoreList += ds.Name + "|"
				}
			}

			entityName := sanitizeEntityName(config, cluster.Name, datacenterName)

			e, ms, err := createNewEntityWithMetricSet(config, entityTypeCluster, entityName, entityName)
			if err != nil {
				config.Logrus.WithError(err).WithField("clusterName", entityName).Error("failed to create metricSet")
				continue
			}

			if config.Args.DatacenterLocation != "" {
				checkError(config.Logrus, ms.SetMetric("datacenterLocation", config.Args.DatacenterLocation, metric.ATTRIBUTE))
			}

			if config.IsVcenterAPIType {
				checkError(config.Logrus, ms.SetMetric("datacenterName", datacenterName, metric.ATTRIBUTE))
			}

			checkError(config.Logrus, ms.SetMetric("networkList", networkList, metric.ATTRIBUTE))
			checkError(config.Logrus, ms.SetMetric("hostList", hostList, metric.ATTRIBUTE))
			checkError(config.Logrus, ms.SetMetric("datastoreList", datastoreList, metric.ATTRIBUTE))

			summary := cluster.Summary.GetComputeResourceSummary()
			if summary != nil {
				checkError(config.Logrus, ms.SetMetric("overallStatus", string(summary.OverallStatus), metric.ATTRIBUTE))
				checkError(config.Logrus, ms.SetMetric("cpu.cores", summary.NumCpuCores, metric.GAUGE))
				checkError(config.Logrus, ms.SetMetric("cpu.threads", summary.NumCpuThreads, metric.GAUGE))
				checkError(config.Logrus, ms.SetMetric("cpu.totalEffectiveMHz", summary.EffectiveCpu, metric.GAUGE))
				checkError(config.Logrus, ms.SetMetric("cpu.totalMHz", summary.TotalCpu, metric.GAUGE))
				checkError(config.Logrus, ms.SetMetric("mem.size", summary.TotalMemory/(1<<20), metric.GAUGE))
				checkError(config.Logrus, ms.SetMetric("mem.effectiveSize", summary.EffectiveMemory, metric.GAUGE))
				checkError(config.Logrus, ms.SetMetric("effectiveHosts", summary.NumEffectiveHosts, metric.GAUGE))
				checkError(config.Logrus, ms.SetMetric("hosts", summary.NumHosts, metric.GAUGE))
			}

			//DRS metrics
			if cluster.Configuration.DrsConfig.Enabled != nil {
				checkError(config.Logrus, ms.SetMetric("drsConfig.enabled", strconv.FormatBool(*cluster.Configuration.DrsConfig.Enabled), metric.ATTRIBUTE))
			}
			if cluster.Configuration.DrsConfig.EnableVmBehaviorOverrides != nil {
				checkError(config.Logrus, ms.SetMetric("drsConfig.enableVmBehaviorOverrides", strconv.FormatBool(*cluster.Configuration.DrsConfig.EnableVmBehaviorOverrides), metric.ATTRIBUTE))
			}
			checkError(config.Logrus, ms.SetMetric("drsConfig.vmotionRate", cluster.Configuration.DrsConfig.VmotionRate, metric.GAUGE))
			checkError(config.Logrus, ms.SetMetric("drsConfig.defaultVmBehavior", string(cluster.Configuration.DrsConfig.DefaultVmBehavior), metric.ATTRIBUTE))

			//DAS metrics
			if cluster.Configuration.DasConfig.Enabled != nil {
				checkError(config.Logrus, ms.SetMetric("dasConfig.enabled", strconv.FormatBool(*cluster.Configuration.DasConfig.Enabled), metric.ATTRIBUTE))
			}
			if cluster.Configuration.DasConfig.AdmissionControlEnabled != nil {
				checkError(config.Logrus, ms.SetMetric("dasConfig.admissionControlEnabled", strconv.FormatBool(*cluster.Configuration.DasConfig.AdmissionControlEnabled), metric.ATTRIBUTE))
			}
			if cluster.Configuration.DasConfig.DefaultVmSettings != nil {
				checkError(config.Logrus, ms.SetMetric("dasConfig.isolationResponse", cluster.Configuration.DasConfig.DefaultVmSettings.IsolationResponse, metric.ATTRIBUTE))
				checkError(config.Logrus, ms.SetMetric("dasConfig.restartPriority", cluster.Configuration.DasConfig.DefaultVmSettings.RestartPriority, metric.ATTRIBUTE))
				checkError(config.Logrus, ms.SetMetric("dasConfig.restartPriorityTimeout", cluster.Configuration.DasConfig.DefaultVmSettings.RestartPriorityTimeout, metric.GAUGE))
			}
			checkError(config.Logrus, ms.SetMetric("dasConfig.hostMonitoring", cluster.Configuration.DasConfig.HostMonitoring, metric.ATTRIBUTE))
			checkError(config.Logrus, ms.SetMetric("dasConfig.vmMonitoring", cluster.Configuration.DasConfig.VmMonitoring, metric.ATTRIBUTE))
			checkError(config.Logrus, ms.SetMetric("dasConfig.vmComponentProtecting", cluster.Configuration.DasConfig.VmComponentProtecting, metric.ATTRIBUTE))
			checkError(config.Logrus, ms.SetMetric("dasConfig.hbDatastoreCandidatePolicy", cluster.Configuration.DasConfig.HBDatastoreCandidatePolicy, metric.ATTRIBUTE))

			// Tags
			if config.TagCollectionEnabled() {
				tagsByCategory := config.TagCollector.GetTagsByCategories(cluster.Self)
				for k, v := range tagsByCategory {
					checkError(config.Logrus, ms.SetMetric(tagsPrefix+k, v, metric.ATTRIBUTE))
					// add tags to inventory due to the inventory workaround
					checkError(config.Logrus, e.SetInventoryItem("tags", tagsPrefix+k, v))
				}
			}
			// Performance metrics
			if config.PerfMetricsCollectionEnabled() {
				perfMetrics := dc.GetPerfMetrics(cluster.Self)
				for _, perfMetric := range perfMetrics {
					checkError(config.Logrus, ms.SetMetric(perfMetricPrefix+perfMetric.Counter, perfMetric.Value, metric.GAUGE))
				}
			}
		}
	}
}
