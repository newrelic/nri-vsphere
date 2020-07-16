// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package collect

import (
	"context"

	"github.com/newrelic/nri-vsphere/internal/performance"
	"github.com/vmware/govmomi/vim25/types"

	"github.com/newrelic/nri-vsphere/internal/load"
	"github.com/vmware/govmomi/vim25/mo"
)

// Hosts VMWare
func Hosts(config *load.Config) {
	ctx := context.Background()
	m := config.ViewManager

	for i, dc := range config.Datacenters {

		cv, err := m.CreateContainerView(ctx, dc.Datacenter.Reference(), []string{"HostSystem"}, true)
		if err != nil {
			config.Logrus.WithError(err).Error("failed to create HostSystem container view")
			continue
		}

		defer func() {
			err := cv.Destroy(ctx)
			if err != nil {
				config.Logrus.WithError(err).Error("error while cleaning up host container view")
			}
		}()

		var hosts []mo.HostSystem
		// Reference: http://pubs.vmware.com/vsphere-60/topic/com.vmware.wssdk.apiref.doc/vim.HostSystem.html
		err = cv.Retrieve(
			ctx,
			[]string{"HostSystem"},
			[]string{"summary", "overallStatus", "config", "network", "vm", "runtime", "parent", "datastore"},
			&hosts)
		if err != nil {
			config.Logrus.WithError(err).Error("failed to retrieve HostSystems")
			continue
		}
		if err := collectTags(config, hosts, config.Datacenters[i]); err != nil {
			config.Logrus.WithError(err).Errorf("failed to retrieve tags:%v", err)
		}

		var refSlice []types.ManagedObjectReference
		for j := 0; j < len(hosts); j++ {
			config.Datacenters[i].Hosts[hosts[j].Self] = &hosts[j]
			refSlice = append(refSlice, hosts[j].Self)
		}

		if config.Args.EnableVspherePerfMetrics && dc.PerfCollector != nil {
			collectedData := dc.PerfCollector.Collect(refSlice, dc.PerfCollector.MetricDefinition.Host, performance.RealTimeInterval)
			dc.AddPerfMetrics(collectedData)
		}
	}
}
