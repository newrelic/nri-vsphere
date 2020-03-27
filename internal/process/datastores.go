package process

import (
	"fmt"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-vmware-vsphere/internal/load"
	"github.com/vmware/govmomi/vim25/types"
)

func createDatastoreSamples(timestamp int64) {
	// ctx := context.Background()

	// create new entities for each host
	for _, dc := range load.Datacenters {
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
			workingEntity, err := load.Integration.Entity(ds.Summary.Name, "datastore", id)
			if err != nil {
				load.Logrus.WithError(err).Error("failed to create entity")
			}

			// create SystemSample metric set
			systemSampleMetricSet := workingEntity.NewMetricSet("VSphereDatastoreSample")

			// defaults
			checkError(systemSampleMetricSet.SetMetric("integration_version", load.IntegrationVersion, metric.ATTRIBUTE))
			checkError(systemSampleMetricSet.SetMetric("integration_name", load.IntegrationName, metric.ATTRIBUTE))
			checkError(systemSampleMetricSet.SetMetric("timestamp", timestamp, metric.GAUGE))
			checkError(systemSampleMetricSet.SetMetric("instanceType", "vmware-datastore", metric.ATTRIBUTE))
			checkError(systemSampleMetricSet.SetMetric("datacenterLocation", load.Args.DatacenterLocation, metric.ATTRIBUTE))
			checkError(systemSampleMetricSet.SetMetric("type", "datastore", metric.ATTRIBUTE))

			// Ds metrics
			checkError(systemSampleMetricSet.SetMetric("overallStatus", string(ds.OverallStatus), metric.ATTRIBUTE))
			checkError(systemSampleMetricSet.SetMetric("accessible", fmt.Sprintf("%t", ds.Summary.Accessible), metric.ATTRIBUTE))
			// Properties not valid if accessible is false
			if ds.Summary.Accessible {
				checkError(systemSampleMetricSet.SetMetric("url", ds.Summary.Url, metric.ATTRIBUTE))
				checkError(systemSampleMetricSet.SetMetric("capacity", float64(ds.Summary.Capacity)/(1<<30), metric.GAUGE))
				checkError(systemSampleMetricSet.SetMetric("freespace", float64(ds.Summary.FreeSpace)/(1<<30), metric.GAUGE))
				checkError(systemSampleMetricSet.SetMetric("uncommitted", float64(ds.Summary.Uncommitted)/(1<<30), metric.GAUGE))

				switch info := ds.Info.(type) {
				case *types.NasDatastoreInfo:
					checkError(systemSampleMetricSet.SetMetric("nas.remoteHost", info.Nas.RemoteHost, metric.ATTRIBUTE))
					checkError(systemSampleMetricSet.SetMetric("nas.remotePath", info.Nas.RemotePath, metric.ATTRIBUTE))
				}
			}
		}
	}
}
