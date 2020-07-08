// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package process

import (
	"fmt"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/nri-vsphere/internal/load"
	"github.com/vmware/govmomi/vim25/types"
)

func createDatastoreSamples(config *load.Config) {
	for _, dc := range config.Datacenters {
		for _, ds := range dc.Datastores {
			datacenterName := dc.Datacenter.Name

			entityName := sanitizeEntityName(config, ds.Summary.Name, datacenterName)

			dataStoreID := ds.Summary.Url

			e, ms, err := createNewEntityWithMetricSet(config, entityTypeDatastore, entityName, dataStoreID)
			if err != nil {
				config.Logrus.WithError(err).WithField("datastoreName", entityName).WithField("dataStoreID", dataStoreID).Error("failed to create metricSet")
				continue
			}

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
			checkError(config, ms.SetMetric("hostCount", len(ds.Host), metric.GAUGE))
			checkError(config, ms.SetMetric("url", ds.Summary.Url, metric.ATTRIBUTE))
			checkError(config, ms.SetMetric("capacity", float64(ds.Summary.Capacity)/(1<<30), metric.GAUGE))
			checkError(config, ms.SetMetric("freeSpace", float64(ds.Summary.FreeSpace)/(1<<30), metric.GAUGE))
			checkError(config, ms.SetMetric("uncommitted", float64(ds.Summary.Uncommitted)/(1<<30), metric.GAUGE))

			switch info := ds.Info.(type) {
			case *types.NasDatastoreInfo:
				if info.Nas != nil {
					checkError(config, ms.SetMetric("nas.remoteHost", info.Nas.RemoteHost, metric.ATTRIBUTE))
					checkError(config, ms.SetMetric("nas.remotePath", info.Nas.RemotePath, metric.ATTRIBUTE))

				}
			}

			// Tags
			tagsByCategory := dc.GetTagsByCategories(ds.Self)
			for k, v := range tagsByCategory {
				checkError(config, ms.SetMetric(tagsPrefix+k, v, metric.ATTRIBUTE))
				// add tags to inventory due to the inventory workaround
				checkError(config, e.SetInventoryItem("tags", tagsPrefix+k, v))
			}

		}
	}
}
