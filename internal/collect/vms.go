// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package collect

import (
	"context"
	"time"

	"github.com/newrelic/nri-vsphere/internal/config"
	"github.com/newrelic/nri-vsphere/internal/model/tag"
	"github.com/newrelic/nri-vsphere/internal/performance"

	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

// VirtualMachines vms
func VirtualMachines(config *config.Config) {
	now := time.Now()

	ctx := context.Background()
	m := config.ViewManager

	collectTags := config.TagCollectionEnabled()
	filterByTag := config.TagFilteringEnabled()

	// Reference: http://pubs.vmware.com/vsphere-60/topic/com.vmware.wssdk.apiref.doc/vim.VirtualMachine.html
	propertiesToRetrieve := []string{"name", "summary", "network", "config", "guest", "runtime", "resourcePool", "datastore", "overallStatus"}
	if config.Args.EnableVsphereSnapshots {
		config.Logrus.Debug("collecting as well snapshot and layoutEx properties")
		propertiesToRetrieve = append(propertiesToRetrieve, "snapshot", "layoutEx.file", "layoutEx.snapshot")
	}

	for i, dc := range config.Datacenters {
		logger := config.Logrus.WithField("datacenter", dc.Datacenter.Name)

		cv, err := m.CreateContainerView(ctx, dc.Datacenter.Reference(), []string{VIRTUAL_MACHINE}, true)
		if err != nil {
			logger.WithError(err).Error("failed to create VirtualMachine container view")
			continue
		}

		defer func() {
			err := cv.Destroy(ctx)
			if err != nil {
				config.Logrus.WithError(err).Error("error while cleaning up virtual machines container view")
			}
		}()

		logger.WithField("seconds", time.Since(now).Seconds()).Debug("before collecting vm data method.Retrieve")

		var vms []mo.VirtualMachine
		err = cv.Retrieve(ctx, []string{VIRTUAL_MACHINE}, propertiesToRetrieve, &vms)
		if err != nil {
			logger.WithError(err).WithField("datacenter", dc.Datacenter.Name).
				Error("failed to retrieve VM data for datacenter")
			continue
		}
		logger.WithField("seconds", time.Since(now).Seconds()).Debug("after collecting vm data method.Retrieve")

		var objectTags tag.TagsByObject
		if collectTags {
			objectTags, err = tag.FetchTagsForObjects(config.TagsManager, vms)
			if err != nil {
				logger.WithError(err).Warn("failed to retrieve tags for virtual machines")
			} else {
				logger.WithField("seconds", time.Since(now).Seconds()).Debug("vms tags collected")
			}
		}

		logger.WithField("seconds", time.Since(now).Seconds()).Debug("vm tags collected")

		var vmRefs []types.ManagedObjectReference
		for _, vm := range vms {
			if filterByTag && len(objectTags) == 0 {
				logger.WithField("virtual machine", vm.Name).
					Debug("ignoring virtual machine since not tags were collected and we have filters configured")
				continue
			}
			// if object has no tags attached or no tag matches any of the tag filters, object will be ignored
			if filterByTag && !tag.MatchObjectTags(objectTags[vm.Reference()]) {
				logger.WithField("virtual machine", vm.Name).
					Debug("ignoring virtual machine since it does not match any configured tag")
				continue
			}

			config.Datacenters[i].VirtualMachines[vm.Self] = &vm
			vmRefs = append(vmRefs, vm.Self)
		}

		if config.Args.EnableVspherePerfMetrics && dc.PerfCollector != nil {
			collectedData := dc.PerfCollector.Collect(vmRefs, dc.PerfCollector.MetricDefinition.VM, performance.RealTimeInterval)
			dc.AddPerfMetrics(collectedData)
		}

		config.Logrus.
			WithField("seconds", time.Since(now).Seconds()).
			Debug("vms perf metrics collected")
	}
}
