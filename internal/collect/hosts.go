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

// Hosts VMWare
func Hosts(config *config.Config) {
	now := time.Now()

	ctx := context.Background()
	m := config.ViewManager

	collectTags := config.TagCollectionEnabled()
	filterByTag := config.TagFilteringEnabled()

	// Reference: http://pubs.vmware.com/vsphere-60/topic/com.vmware.wssdk.apiref.doc/vim.HostSystem.html
	propertiesToRetrieve := []string{"summary", "overallStatus", "config", "network", "vm", "runtime", "parent", "datastore"}
	for i, dc := range config.Datacenters {
		logger := config.Logrus.WithField("datacenter", dc.Datacenter.Name)

		cv, err := m.CreateContainerView(ctx, dc.Datacenter.Reference(), []string{HOST}, true)
		if err != nil {
			logger.WithError(err).Error("failed to create HostSystem container view")
			continue
		}

		defer func() {
			err := cv.Destroy(ctx)
			if err != nil {
				config.Logrus.WithError(err).Error("error while cleaning up host container view")
			}
		}()

		var hosts []mo.HostSystem
		err = cv.Retrieve(ctx, []string{HOST}, propertiesToRetrieve, &hosts)
		if err != nil {
			logger.WithError(err).Error("failed to retrieve HostSystems")
			continue
		}

		var objectTags tag.TagsByObject
		if collectTags {
			objectTags, err = config.TagCollector.FetchTagsForObjects(hosts)
			if err != nil {
				logger.WithError(err).Warn("failed to retrieve tags for hosts", err)
			} else {
				logger.WithField("seconds", time.Since(now).Seconds()).Debug("hosts tags collected")
			}
		}

		var hostsRefs []types.ManagedObjectReference
		for _, host := range hosts {
			if filterByTag && len(objectTags) == 0 {
				logger.WithField("host", host.Name).
					Debug("ignoring host since no tags were collected and we have filters configured")
				continue
			}
			// if object has no tags attached or no tag matches any of the tag filters, object will be ignored
			if filterByTag && !config.TagCollector.MatchObjectTags(objectTags[host.Reference()]) {
				logger.WithField("host", host.Name).
					Debug("ignoring host since it does not match any configured tag")
				continue
			}

			config.Datacenters[i].Hosts[host.Self] = &host
			hostsRefs = append(hostsRefs, host.Self)
		}

		if config.Args.EnableVspherePerfMetrics && dc.PerfCollector != nil {
			collectedData := dc.PerfCollector.Collect(hostsRefs, dc.PerfCollector.MetricDefinition.Host, performance.RealTimeInterval)
			dc.AddPerfMetrics(collectedData)
		}

		logger.WithField("seconds", time.Since(now).Seconds()).Debug("hosts perf metrics collected")
	}
}
