// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package collect

import (
	"context"

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

		defer cv.Destroy(ctx)

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
		for j := 0; j < len(hosts); j++ {
			config.Datacenters[i].Hosts[hosts[j].Self] = &hosts[j]
		}
	}
}
