// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package collect

import (
	"context"
	"github.com/newrelic/nri-vsphere/internal/load"
	"github.com/vmware/govmomi/vim25/mo"
)

// Clusters VMWare
func Clusters(config *load.Config) {
	ctx := context.Background()
	m := config.ViewManager

	for i, dc := range config.Datacenters {
		cv, err := m.CreateContainerView(ctx, dc.Datacenter.Reference(), []string{"ComputeResource"}, true)
		if err != nil {
			config.Logrus.WithError(err).Error("failed to create ComputeResource container view")
			continue
		}
		defer cv.Destroy(ctx)
		var clusters []mo.ClusterComputeResource
		// Reference: https://code.vmware.com/apis/704/vsphere/vim.ClusterComputeResource.html
		err = cv.Retrieve(
			ctx,
			[]string{"ClusterComputeResource"},
			[]string{"summary", "host", "datastore", "name", "network", "configuration"},
			&clusters)
		if err != nil {
			config.Logrus.WithError(err).Error("failed to retrieve ClusterComputeResource")
			continue
		}
		for j := 0; j < len(clusters); j++ {
			config.Datacenters[i].Clusters[clusters[j].Self] = &clusters[j]
		}
	}
}
