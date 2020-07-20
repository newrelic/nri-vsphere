// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package collect

import (
	"context"
	"time"

	"github.com/newrelic/nri-vsphere/internal/config"
	"github.com/newrelic/nri-vsphere/internal/performance"
	"github.com/newrelic/nri-vsphere/internal/tag"
	"github.com/vmware/govmomi/vim25/types"

	"github.com/vmware/govmomi/vim25/mo"
)

// Clusters VMWare
func Clusters(config *config.Config) {
	now := time.Now()

	ctx := context.Background()
	m := config.ViewManager

	collectTags := config.TagCollectionEnabled()
	filterByTag := config.TagFilteringEnabled()

	propertiesToRetrieve := []string{"summary", "host", "datastore", "name", "network", "configuration"}
	for i, dc := range config.Datacenters {
		logger := config.Logrus.WithField("datacenter", dc.Datacenter.Name)

		cv, err := m.CreateContainerView(ctx, dc.Datacenter.Reference(), []string{CLUSTER}, true)
		if err != nil {
			logger.WithError(err).Error("failed to create ComputeResource container view")
			continue
		}
		defer func() {
			err := cv.Destroy(ctx)
			if err != nil {
				config.Logrus.WithError(err).Error("error while cleaning up cluster container view")
			}
		}()

		var clusters []mo.ClusterComputeResource
		// Reference: https://code.vmware.com/apis/704/vsphere/vim.ClusterComputeResource.html
		err = cv.Retrieve(ctx, []string{CLUSTER}, propertiesToRetrieve, &clusters)
		if err != nil {
			logger.WithError(err).Error("failed to retrieve ClusterComputeResource")
			continue
		}

		var objectTags tag.TagsByObject
		if collectTags {
			objectTags, err = config.TagCollector.FetchTagsForObjects(clusters)
			if err != nil {
				logger.WithError(err).Warn("failed to retrieve tags for clusters", err)
			} else {
				logger.WithField("seconds", time.Since(now).Seconds()).Debug("clusters tags collected")
			}
		}

		var clusterRefs []types.ManagedObjectReference
		for _, cluster := range clusters {
			if filterByTag && len(objectTags) == 0 {
				logger.WithField("cluster", cluster.Name).
					Debugf("ignoring cluster since not tags were collected and we have filters configured")
				continue
			}
			// if object has no tags attached or no tag matches any of the tag filters, object will be ignored
			if filterByTag && !config.TagCollector.MatchObjectTags(objectTags[cluster.Reference()]) {
				logger.WithField("cluster", cluster.Name).
					Debugf("ignoring cluster since it does not match any configured tag")
				continue
			}

			config.Datacenters[i].Clusters[cluster.Self] = &cluster
			clusterRefs = append(clusterRefs, cluster.Self)
		}

		if config.PerfMetricsCollectionEnabled() {
			metricsToCollect := config.PerfCollector.MetricDefinition.ClusterComputeResource
			collectedData := config.PerfCollector.Collect(clusterRefs, metricsToCollect, performance.FiveMinutesInterval)
			dc.AddPerfMetrics(collectedData)
		}

		logger.WithField("seconds", time.Since(now).Seconds()).Debug("clusters perf metrics collected")
	}
}
