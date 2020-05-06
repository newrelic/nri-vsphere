// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package process

import (
	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/nri-vsphere/internal/load"
	"strconv"
)

func createClusterSamples(config *load.Config) {
	for _, dc := range config.Datacenters {
		for _, cluster := range dc.Clusters {
			// // resolve hypervisor host
			summary := cluster.Summary.GetComputeResourceSummary()
			datacenterName := dc.Datacenter.Name

			//Retrieving the list of host belonging to the cluster
			hostList := ""
			for _, hostReference := range cluster.Host {
				if host, ok := dc.Hosts[hostReference.Reference()]; ok {
					hostList = hostList + host.Summary.Config.Name + "|"
				}
			}

			//Retrieving the list of networks attached to the cluster
			networkList := ""
			for _, networkReference := range cluster.Network {
				if network, ok := dc.Networks[networkReference]; ok {
					networkList = networkList + network.Name + "|"
				}
			}

			//Retrieving the list of datastores attached to the cluster
			datastoreList := ""
			for _, datastoreReference := range cluster.Datastore {
				if datastore, ok := dc.Datastores[datastoreReference]; ok {
					datastoreList = datastoreList + datastore.Name + "|"
				}
			}

			entityName := sanitizeEntityName(config, cluster.Name, datacenterName)

			ms := createNewEntityWithMetricSet(config, entityTypeCluster, entityName, entityName)

			if config.Args.DatacenterLocation != "" {
				checkError(config, ms.SetMetric("datacenterLocation", config.Args.DatacenterLocation, metric.ATTRIBUTE))
			}

			if config.IsVcenterAPIType {
				checkError(config, ms.SetMetric("datacenterName", datacenterName, metric.ATTRIBUTE))
			}

			checkError(config, ms.SetMetric("networkList", networkList, metric.ATTRIBUTE))
			checkError(config, ms.SetMetric("hostList", hostList, metric.ATTRIBUTE))
			checkError(config, ms.SetMetric("datastoreList", datastoreList, metric.ATTRIBUTE))

			checkError(config, ms.SetMetric("overallStatus", string(summary.OverallStatus), metric.ATTRIBUTE))

			checkError(config, ms.SetMetric("cpu.cores", summary.NumCpuCores, metric.GAUGE))
			checkError(config, ms.SetMetric("cpu.threads", summary.NumCpuThreads, metric.GAUGE))
			checkError(config, ms.SetMetric("cpu.totalEffectiveMHz", summary.EffectiveCpu, metric.GAUGE))
			checkError(config, ms.SetMetric("cpu.totalMHz", summary.TotalCpu, metric.GAUGE))
			checkError(config, ms.SetMetric("mem.size", summary.TotalMemory/(1<<20), metric.GAUGE))
			checkError(config, ms.SetMetric("mem.effectiveSize", summary.EffectiveMemory, metric.GAUGE))
			checkError(config, ms.SetMetric("effectiveHosts", summary.NumEffectiveHosts, metric.GAUGE))
			checkError(config, ms.SetMetric("hosts", summary.NumHosts, metric.GAUGE))

			//DRS metrics
			if cluster.Configuration.DrsConfig.Enabled != nil {
				checkError(config, ms.SetMetric("drsConfig.Enabled", strconv.FormatBool(*cluster.Configuration.DrsConfig.Enabled), metric.ATTRIBUTE))
			}
			if cluster.Configuration.DrsConfig.EnableVmBehaviorOverrides != nil {
				checkError(config, ms.SetMetric("drsConfig.EnableVmBehaviorOverrides", strconv.FormatBool(*cluster.Configuration.DrsConfig.EnableVmBehaviorOverrides), metric.ATTRIBUTE))
			}
			checkError(config, ms.SetMetric("drsConfig.VmotionRate", cluster.Configuration.DrsConfig.VmotionRate, metric.GAUGE))
			checkError(config, ms.SetMetric("drsConfig.DefaultVmBehavior", string(cluster.Configuration.DrsConfig.DefaultVmBehavior), metric.ATTRIBUTE))

			//DAS metrics
			if cluster.Configuration.DasConfig.Enabled != nil {
				checkError(config, ms.SetMetric("dasConfig.Enabled", strconv.FormatBool(*cluster.Configuration.DasConfig.Enabled), metric.ATTRIBUTE))
			}
			if cluster.Configuration.DasConfig.AdmissionControlEnabled != nil {
				checkError(config, ms.SetMetric("dasConfig.AdmissionControlEnabled", strconv.FormatBool(*cluster.Configuration.DasConfig.AdmissionControlEnabled), metric.ATTRIBUTE))
			}
			if cluster.Configuration.DasConfig.DefaultVmSettings != nil {
				checkError(config, ms.SetMetric("dasConfig.IsolationResponse", cluster.Configuration.DasConfig.DefaultVmSettings.IsolationResponse, metric.ATTRIBUTE))
				checkError(config, ms.SetMetric("dasConfig.RestartPriority", cluster.Configuration.DasConfig.DefaultVmSettings.RestartPriority, metric.ATTRIBUTE))
				checkError(config, ms.SetMetric("dasConfig.RestartPriorityTimeout", cluster.Configuration.DasConfig.DefaultVmSettings.RestartPriorityTimeout, metric.GAUGE))
			}
			checkError(config, ms.SetMetric("dasConfig.HostMonitoring", cluster.Configuration.DasConfig.HostMonitoring, metric.ATTRIBUTE))
			checkError(config, ms.SetMetric("dasConfig.VmMonitoring", cluster.Configuration.DasConfig.VmMonitoring, metric.ATTRIBUTE))
			checkError(config, ms.SetMetric("dasConfig.VmComponentProtecting", cluster.Configuration.DasConfig.VmComponentProtecting, metric.ATTRIBUTE))
			checkError(config, ms.SetMetric("dasConfig.HBDatastoreCandidatePolicy", cluster.Configuration.DasConfig.HBDatastoreCandidatePolicy, metric.ATTRIBUTE))

		}
	}
}
