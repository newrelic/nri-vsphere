package process

import (
	"fmt"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-vmware-vsphere/internal/load"
	"github.com/vmware/govmomi/vim25/types"
)

func createDatastoreSamples(config *load.Config, timestamp int64) {
	// ctx := context.Background()

	// create new entities for each host
	for _, dc := range config.Datacenters {
		for _, ds := range dc.Datastores {
			// entityName := ds.Summary.Name + ":ds"

			// if load.Args.DatacenterLocation != "" {
			// 	entityName = load.Args.DatacenterLocation + ":" + entityName
			// }

			// entityName = strings.ToLower(entityName)
			// entityName = strings.ReplaceAll(entityName, ".", "-")

			// workingEntity := setEntity(entityName, "vmware") // default type instance
			// workingEntity.SetInventoryItem("name", "value", fmt.Sprintf("%v:%d", entityName, timestamp))

			id := integration.IDAttribute{Key: "id", Value: ds.Summary.Datastore.Value}
			workingEntity, err := config.Integration.Entity(ds.Summary.Name, "datastore", id)
			if err != nil {
				config.Logrus.WithError(err).Error("failed to create entity")
			}

			// create SystemSample metric set
			systemSampleMetricSet := workingEntity.NewMetricSet("VSphereDatastoreSample")

			// defaults
			checkError(config, systemSampleMetricSet.SetMetric("integration_version", config.IntegrationVersion, metric.ATTRIBUTE))
			checkError(config, systemSampleMetricSet.SetMetric("integration_name", config.IntegrationName, metric.ATTRIBUTE))
			checkError(config, systemSampleMetricSet.SetMetric("timestamp", timestamp, metric.GAUGE))
			checkError(config, systemSampleMetricSet.SetMetric("instanceType", "vmware-datastore", metric.ATTRIBUTE))
			checkError(config, systemSampleMetricSet.SetMetric("datacenterLocation", config.Args.DatacenterLocation, metric.ATTRIBUTE))
			checkError(config, systemSampleMetricSet.SetMetric("type", "datastore", metric.ATTRIBUTE))

			// Ds metrics
			checkError(config, systemSampleMetricSet.SetMetric("overallStatus", string(ds.OverallStatus), metric.ATTRIBUTE))
			checkError(config, systemSampleMetricSet.SetMetric("accessible", fmt.Sprintf("%t", ds.Summary.Accessible), metric.ATTRIBUTE))
			// Properties not valid if accessible is false
			if ds.Summary.Accessible {
				checkError(config, systemSampleMetricSet.SetMetric("url", ds.Summary.Url, metric.ATTRIBUTE))
				checkError(config, systemSampleMetricSet.SetMetric("capacity", float64(ds.Summary.Capacity)/(1<<30), metric.GAUGE))
				checkError(config, systemSampleMetricSet.SetMetric("freespace", float64(ds.Summary.FreeSpace)/(1<<30), metric.GAUGE))
				checkError(config, systemSampleMetricSet.SetMetric("uncommitted", float64(ds.Summary.Uncommitted)/(1<<30), metric.GAUGE))

				switch info := ds.Info.(type) {
				case *types.NasDatastoreInfo:
					checkError(config, systemSampleMetricSet.SetMetric("nas.remoteHost", info.Nas.RemoteHost, metric.ATTRIBUTE))
					checkError(config, systemSampleMetricSet.SetMetric("nas.remotePath", info.Nas.RemotePath, metric.ATTRIBUTE))
				}
			}
		}
	}
}
