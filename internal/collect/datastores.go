// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package collect

import (
	"context"
	"time"

	"github.com/newrelic/nri-vsphere/internal/config"
	"github.com/newrelic/nri-vsphere/internal/performance"
	"github.com/newrelic/nri-vsphere/internal/tag"

	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

// Datastores collects data of all datastores
func Datastores(config *config.Config) {
	now := time.Now()

	ctx := context.Background()
	m := config.ViewManager

	collectTags := config.TagCollectionEnabled()
	filterByTag := config.TagFilteringEnabled()

	// Reference: https://code.vmware.com/apis/42/vsphere/doc/vim.Datastore.html
	propertiesToRetrieve := []string{"name", "summary", "overallStatus", "vm", "host", "info"}
	for i, dc := range config.Datacenters {
		logger := config.Logrus.WithField("datacenter", dc.Datacenter.Name)

		cv, err := m.CreateContainerView(ctx, dc.Datacenter.Reference(), []string{DATASTORE}, true)
		if err != nil {
			logger.WithError(err).Error("failed to create Datastore container view")
			continue
		}
		defer func() {
			err := cv.Destroy(ctx)
			if err != nil {
				config.Logrus.WithError(err).Error("error while cleaning up datastores container view")
			}
		}()

		var datastores []mo.Datastore
		err = cv.Retrieve(ctx, []string{DATASTORE}, propertiesToRetrieve, &datastores)
		if err != nil {
			logger.WithError(err).Error("failed to retrieve Datastore")
			continue
		}

		var objectTags tag.TagsByObject
		if collectTags {
			objectTags, err = config.TagCollector.FetchTagsForObjects(datastores)
			if err != nil {
				logger.WithError(err).Warn("failed to retrieve tags for datastores", err)
			} else {
				logger.WithField("seconds", time.Since(now).Seconds()).Debug("datastores tags collected")
			}
		}

		var dsRefs []types.ManagedObjectReference
		for _, ds := range datastores {
			if filterByTag && len(objectTags) == 0 {
				logger.WithField("datastore", ds.Name).
					Debug("ignoring datastore since not tags were collected and we have filters configured")
				continue
			}
			// if object has no tags attached or no tag matches any of the tag filters, object will be ignored
			if filterByTag && !config.TagCollector.MatchObjectTags(objectTags[ds.Reference()]) {
				logger.WithField("datastore", ds.Name).
					Debug("ignoring datastore since it does not match any configured tag")
				continue
			}

			config.Datacenters[i].Datastores[ds.Self] = &ds
			dsRefs = append(dsRefs, ds.Self)
		}

		if config.PerfMetricsCollectionEnabled() {
			metricsToCollect := config.PerfCollector.MetricDefinition.Datastore
			collectedData := config.PerfCollector.Collect(dsRefs, metricsToCollect, performance.FiveMinutesInterval)
			dc.AddPerfMetrics(collectedData)
		}

		logger.WithField("seconds", time.Since(now).Seconds()).Debug("datastores perf metrics collected")
	}
}
