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

// ResourcePools VMWare
func ResourcePools(config *config.Config) {
	now := time.Now()

	ctx := context.Background()
	m := config.ViewManager

	collectTags := config.TagCollectionEnabled()
	filterByTag := config.TagFilteringEnabled()

	propertiesToRetrieve := []string{"summary", "owner", "parent", "runtime", "name", "overallStatus", "vm", "resourcePool"}
	for i, dc := range config.Datacenters {
		logger := config.Logrus.WithField("datacenter", dc.Datacenter.Name)

		cv, err := m.CreateContainerView(ctx, dc.Datacenter.Reference(), []string{RESOURCE_POOL}, true)
		if err != nil {
			logger.WithError(err).Error("failed to create ResourcePool container view")
			continue
		}

		defer func() {
			err := cv.Destroy(ctx)
			if err != nil {
				config.Logrus.WithError(err).Error("error while cleaning up resourcePools container view")
			}
		}()

		var resourcePools []mo.ResourcePool
		err = cv.Retrieve(ctx, []string{RESOURCE_POOL}, propertiesToRetrieve, &resourcePools)
		if err != nil {
			logger.WithError(err).Error("failed to retrieve ResourcePools")
			continue
		}

		var objectTags tag.TagsByObject
		if collectTags {
			objectTags, err = config.TagCollector.FetchTagsForObjects(resourcePools)
			if err != nil {
				logger.WithError(err).Warn("failed to retrieve tags for resourcePools", err)
			} else {
				logger.WithField("seconds", time.Since(now).Seconds()).Debug("resourcePools tags collected")
			}
		}

		var rpRefs []types.ManagedObjectReference
		for _, rp := range resourcePools {
			if filterByTag && len(objectTags) == 0 {
				logger.WithField("resource pool", rp.Name).
					Debug("ignoring resource pool since no tags were collected and we have filters configured")
				continue
			}
			// if object has no tags attached or no tag matches any of the tag filters, object will be ignored
			if filterByTag && !config.TagCollector.MatchObjectTags(objectTags[rp.Reference()]) {
				logger.WithField("resource pool", rp.Name).
					Debug("ignoring resource pool since it does not match any configured tag")
				continue
			}

			config.Datacenters[i].ResourcePools[rp.Self] = &rp
			rpRefs = append(rpRefs, rp.Self)
		}

		if config.Args.EnableVspherePerfMetrics && dc.PerfCollector != nil {
			collectedData := dc.PerfCollector.Collect(rpRefs, dc.PerfCollector.MetricDefinition.ResourcePool, performance.FiveMinutesInterval)
			dc.AddPerfMetrics(collectedData)
		}

		logger.WithField("seconds", time.Since(now).Seconds()).Debug("resource pools perf metrics collected")
	}
}
