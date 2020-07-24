// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package collect

import (
	"context"
	"github.com/newrelic/nri-vsphere/internal/config"
	"github.com/newrelic/nri-vsphere/internal/performance"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

// VirtualMachines vms
func VirtualMachines(config *config.Config) {
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

		logger.WithField("seconds", config.Uptime().Seconds()).Debug("before collecting vm data method.Retrieve")

		var vms []mo.VirtualMachine
		err = cv.Retrieve(ctx, []string{VIRTUAL_MACHINE}, propertiesToRetrieve, &vms)
		if err != nil {
			logger.WithError(err).WithField("datacenter", dc.Datacenter.Name).
				Error("failed to retrieve VM data for datacenter")
			continue
		}
		logger.WithField("seconds", config.Uptime().Seconds()).Debug("after collecting vm data method.Retrieve")

		if collectTags {
			_, err = config.TagCollector.FetchTagsForObjects(vms)
			if err != nil {
				logger.WithError(err).Warn("failed to retrieve tags for virtual machines")
			} else {
				logger.WithField("seconds", config.Uptime().Seconds()).Debug("vms tags collected")
			}
		}

		logger.WithField("seconds", config.Uptime().Seconds()).Debug("vm tags collected")

		var vmRefs []types.ManagedObjectReference
		for _, vm := range vms {
			if filterByTag && !config.TagCollector.MatchObjectTags(vm.Reference()) {
				logger.WithField("virtual machine", vm.Name).
					Debug("ignoring virtual machine since no tags matched the configured filters")
				continue
			}

			config.Datacenters[i].VirtualMachines[vm.Self] = &vm
			vmRefs = append(vmRefs, vm.Self)
		}

		if config.PerfMetricsCollectionEnabled() {
			metricsToCollect := config.PerfCollector.MetricDefinition.VM
			collectedData := config.PerfCollector.Collect(vmRefs, metricsToCollect, performance.RealTimeInterval)
			dc.AddPerfMetrics(collectedData)

			logger.WithField("seconds", config.Uptime().Seconds()).Debug("vms perf metrics collected")
		}
	}
}
