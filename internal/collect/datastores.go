// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package collect

import (
	"context"

	"github.com/newrelic/nri-vmware-vsphere/internal/load"
	"github.com/vmware/govmomi/vim25/mo"
)

// Datastores collects data of all datastores
func Datastores(config *load.Config) {
	ctx := context.Background()
	m := config.ViewManager

	for i, dc := range config.Datacenters {
		cv, err := m.CreateContainerView(ctx, dc.Datacenter.Reference(), []string{"Datastore"}, true)
		if err != nil {
			config.Logrus.WithError(err).Fatal("failed to create Datastore container view")
		}
		defer cv.Destroy(ctx)

		var datastores []mo.Datastore
		// Reference: https://code.vmware.com/apis/42/vsphere/doc/vim.Datastore.html
		err = cv.Retrieve(ctx, []string{"Datastore"}, nil, &datastores)
		if err != nil {
			config.Logrus.WithError(err).Fatal("failed to retrieve Datastore")
		}
		for j := 0; j < len(datastores); j++ {
			config.Datacenters[i].Datastores[datastores[j].Self] = &datastores[j]
		}
	}
}
