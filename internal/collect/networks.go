// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package collect

import (
	"context"

	"github.com/newrelic/nri-vsphere/internal/load"
	"github.com/vmware/govmomi/vim25/mo"
)

// Networks ESXi
func Networks(config *load.Config) {
	ctx := context.Background()
	m := config.ViewManager

	for i, dc := range config.Datacenters {
		cv, err := m.CreateContainerView(ctx, dc.Datacenter.Reference(), []string{"Network"}, true)
		if err != nil {
			config.Logrus.WithError(err).Error("failed to create Network container view")
			continue
		}
		defer cv.Destroy(ctx)

		var networks []mo.Network
		// Reference: http://pubs.vmware.com/vsphere-60/topic/com.vmware.wssdk.apiref.doc/vim.Network.html
		err = cv.Retrieve(ctx, []string{"Network"}, []string{"name"}, &networks)
		if err != nil {
			config.Logrus.WithError(err).Error("failed to retrieve Networks")
			continue
		}
		for j := 0; j < len(networks); j++ {
			config.Datacenters[i].Networks[networks[j].Self] = &networks[j]
		}
	}
}
