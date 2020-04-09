package process

import (
	"fmt"
	"strings"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/nri-vmware-vsphere/internal/load"
	"github.com/vmware/govmomi/vim25/types"
)

func createDatastoreSamples(config *load.Config, timestamp int64) {
	for _, dc := range config.Datacenters {
		for _, ds := range dc.Datastores {
			entityName := ds.Summary.Name + ":datastore"
			datacenterName := dc.Datacenter.Name
			if config.IsVcenterAPIType {
				entityName = datacenterName + ":" + entityName
			}

			if config.Args.DatacenterLocation != "" {
				entityName = config.Args.DatacenterLocation + ":" + entityName
			}

			entityName = strings.ToLower(entityName)
			entityName = strings.ReplaceAll(entityName, ".", "-")

			dataStoreID := ds.Summary.Url

			workingEntity, err := config.Integration.Entity(dataStoreID, "vsphere-datastore")
			if err != nil {
				config.Logrus.WithError(err).Error("failed to create entity")
			}

			// entity displayName
			err = workingEntity.SetInventoryItem("vsphereDatastore", "name", entityName)
			if err != nil {
				config.Logrus.WithError(err).Error("failed to set datastore inventory")
			}

			ms := workingEntity.NewMetricSet("VSphereDatastoreSample")

			if config.Args.DatacenterLocation != "" {
				checkError(config, ms.SetMetric("datacenterLocation", config.Args.DatacenterLocation, metric.ATTRIBUTE))
			}
			if config.IsVcenterAPIType {
				checkError(config, ms.SetMetric("datacenterName", datacenterName, metric.ATTRIBUTE))
			}

			checkError(config, ms.SetMetric("name", ds.Summary.Name, metric.ATTRIBUTE))
			checkError(config, ms.SetMetric("fileSystemType", ds.Summary.Type, metric.ATTRIBUTE))
			checkError(config, ms.SetMetric("overallStatus", string(ds.OverallStatus), metric.ATTRIBUTE))
			checkError(config, ms.SetMetric("accessible", fmt.Sprintf("%t", ds.Summary.Accessible), metric.ATTRIBUTE))
			checkError(config, ms.SetMetric("vmCount", len(ds.Vm), metric.GAUGE))
			checkError(config, ms.SetMetric("hostsCount", len(ds.Host), metric.GAUGE))
			checkError(config, ms.SetMetric("url", ds.Summary.Url, metric.ATTRIBUTE))
			checkError(config, ms.SetMetric("capacity", float64(ds.Summary.Capacity)/(1<<30), metric.GAUGE))
			checkError(config, ms.SetMetric("freeSpace", float64(ds.Summary.FreeSpace)/(1<<30), metric.GAUGE))
			checkError(config, ms.SetMetric("uncommitted", float64(ds.Summary.Uncommitted)/(1<<30), metric.GAUGE))

			switch info := ds.Info.(type) {
			case *types.NasDatastoreInfo:
				checkError(config, ms.SetMetric("nas.remoteHost", info.Nas.RemoteHost, metric.ATTRIBUTE))
				checkError(config, ms.SetMetric("nas.remotePath", info.Nas.RemotePath, metric.ATTRIBUTE))
			}
		}
	}
}
