package process

import (
	"strings"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/nri-vmware-vsphere/internal/load"
)

func createClusterSamples(config *load.Config, timestamp int64) {
	for _, dc := range config.Datacenters {
		for _, cluster := range dc.Clusters {
			// // resolve hypervisor host
			summary := cluster.Summary.GetComputeResourceSummary()
			datacenterName := dc.Datacenter.Name

			//Retrieving the list of host belonging to the cluster
			hostList:=""
			for _, hostReference:= range cluster.Host{
				if host, ok := dc.Hosts[hostReference.Reference()]; ok {
					hostList = hostList+ host.Summary.Config.Name+"|"
				}
			}

			//Retrieving the list of networks attached to the cluster
			networkList:=""
			for _, networkReference:= range cluster.Network{
				if network, ok := dc.Networks[networkReference]; ok {
					networkList = networkList+ network.Name+"|"
				}
			}

			//Retrieving the list of datastores attached to the cluster
			datastoreList:=""
			for _, datastoreReference:= range cluster.Datastore{
				if datastore, ok := dc.Datastores[datastoreReference]; ok {
					datastoreList = datastoreList+ datastore.Name+"|"
				}
			}

			//Creating entity name
			entityName := cluster.Name + ":cluster"
			if config.IsVcenterAPIType {
				entityName = datacenterName + ":" + entityName
			}
			if config.Args.DatacenterLocation != "" {
				entityName = config.Args.DatacenterLocation + ":" + entityName
			}
			entityName = strings.ToLower(entityName)
			entityName = strings.ReplaceAll(entityName, ".", "-")

			// Identifier for cluster entity
			workingEntity, err := config.Integration.Entity(entityName, "vsphere-cluster")
			if err != nil {
				config.Logrus.WithError(err).Error("failed to create entity")
			}

			// entity displayName
			workingEntity.SetInventoryItem("vsphereCluster", "name", entityName)
			ms := workingEntity.NewMetricSet("VSphereClusterSample")

			if config.Args.DatacenterLocation != "" {
				checkError(config, ms.SetMetric("datacenterLocation", config.Args.DatacenterLocation, metric.ATTRIBUTE))
			}

			if config.IsVcenterAPIType {
				checkError(config, ms.SetMetric("datacenterName", datacenterName, metric.ATTRIBUTE))
			}

			checkError(config, ms.SetMetric("NetworkList", networkList, metric.ATTRIBUTE))
			checkError(config, ms.SetMetric("HostList", hostList, metric.ATTRIBUTE))
			checkError(config, ms.SetMetric("DatastoreList", datastoreList, metric.ATTRIBUTE))

			checkError(config, ms.SetMetric("OverallStatus", string(summary.OverallStatus), metric.ATTRIBUTE))

			checkError(config, ms.SetMetric("cpu.cores", summary.NumCpuCores, metric.GAUGE))
			checkError(config, ms.SetMetric("cpu.threads", summary.NumCpuThreads, metric.GAUGE))
			checkError(config, ms.SetMetric("cpu.totalEffectiveMHz", summary.EffectiveCpu, metric.GAUGE))
			checkError(config, ms.SetMetric("cpu.totalMHz", summary.TotalCpu, metric.GAUGE))
			checkError(config, ms.SetMetric("mem.size", summary.TotalMemory / 1e+6, metric.GAUGE))
			checkError(config, ms.SetMetric("mem.effectiveSize", summary.EffectiveMemory, metric.GAUGE))
			checkError(config, ms.SetMetric("effectiveHosts", summary.NumEffectiveHosts, metric.GAUGE))
			checkError(config, ms.SetMetric("hosts", summary.NumHosts, metric.GAUGE))

		}
	}
}
