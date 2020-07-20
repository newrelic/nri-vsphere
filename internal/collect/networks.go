// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package collect

import (
	"context"

	"github.com/newrelic/nri-vsphere/internal/config"

	"github.com/vmware/govmomi/vim25/mo"
)

// Networks ESXi
func Networks(config *config.Config) {
	ctx := context.Background()
	m := config.ViewManager

	// Reference: http://pubs.vmware.com/vsphere-60/topic/com.vmware.wssdk.apiref.doc/vim.Network.html
	propertiesToRetrieve := []string{"name"}
	for i, dc := range config.Datacenters {
		logger := config.Logrus.WithField("datacenter", dc.Datacenter.Name)

		cv, err := m.CreateContainerView(ctx, dc.Datacenter.Reference(), []string{NETWORK}, true)
		if err != nil {
			logger.WithError(err).Error("failed to create Network container view")
			continue
		}
		defer func() {
			err := cv.Destroy(ctx)
			if err != nil {
				config.Logrus.WithError(err).Error("error while cleaning up network container view")
			}
		}()

		var networks []mo.Network
		err = cv.Retrieve(ctx, []string{NETWORK}, propertiesToRetrieve, &networks)
		if err != nil {
			logger.WithError(err).Error("failed to retrieve Networks")
			continue
		}
		for j := 0; j < len(networks); j++ {
			config.Datacenters[i].Networks[networks[j].Self] = &networks[j]
		}
	}
}
