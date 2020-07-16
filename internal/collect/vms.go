// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package collect

import (
	"context"
	"time"

	"github.com/newrelic/nri-vsphere/internal/performance"
	"github.com/vmware/govmomi/vim25/types"

	"github.com/newrelic/nri-vsphere/internal/load"
	"github.com/vmware/govmomi/vim25/mo"
)

// VirtualMachines vms
func VirtualMachines(config *load.Config) {
	ctx := context.Background()
	m := config.ViewManager

	for i, dc := range config.Datacenters {
		cv, err := m.CreateContainerView(ctx, dc.Datacenter.Reference(), []string{"VirtualMachine"}, true)
		if err != nil {
			config.Logrus.WithError(err).Error("failed to create VirtualMachine container view")
			continue
		}

		defer func() {
			err := cv.Destroy(ctx)
			if err != nil {
				config.Logrus.WithError(err).Error("error while cleaning up virtual machines container view")
			}
		}()

		var vms []mo.VirtualMachine

		propertiesToRetrieve := []string{"summary", "network", "config", "guest", "runtime", "resourcePool", "datastore", "overallStatus"}
		if config.Args.EnableVsphereSnapshots {
			config.Logrus.Debug("collecting as well snapshot and layoutEx properties")
			propertiesToRetrieve = append(propertiesToRetrieve, "snapshot", "layoutEx.file", "layoutEx.snapshot")
		}

		config.Logrus.WithField("seconds", time.Since(load.Now).Seconds()).Debug("before collecting vm data method.Retrieve")

		// Reference: http://pubs.vmware.com/vsphere-60/topic/com.vmware.wssdk.apiref.doc/vim.VirtualMachine.html
		err = cv.Retrieve(ctx, []string{"VirtualMachine"}, propertiesToRetrieve, &vms)
		if err != nil {
			config.Logrus.WithError(err).Error("failed to retrieve VM Summaries")
			continue
		}
		config.Logrus.WithField("seconds", time.Since(load.Now).Seconds()).Debug("after collecting vm data method.Retrieve")

		if err := collectTags(config, vms, config.Datacenters[i]); err != nil {
			config.Logrus.WithError(err).Errorf("failed to retrieve tags:%v", err)
		} else {
			config.Logrus.WithField("seconds", time.Since(load.Now).Seconds()).Debug("vm tags collected")
		}

		var refSlice []types.ManagedObjectReference
		for j := 0; j < len(vms); j++ {
			config.Datacenters[i].VirtualMachines[vms[j].Self] = &vms[j]
			refSlice = append(refSlice, vms[j].Self)
		}

		if config.Args.EnableVspherePerfMetrics && dc.PerfCollector != nil {
			collectedData := dc.PerfCollector.Collect(refSlice, dc.PerfCollector.MetricDefinition.VM, performance.RealTimeInterval)
			dc.AddPerfMetrics(collectedData)
		}
		config.Logrus.WithField("seconds", time.Since(load.Now).Seconds()).Debug("vm perf metrics collected")

	}
}
