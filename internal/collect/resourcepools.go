// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package collect

import (
	"context"

	"github.com/newrelic/nri-vsphere/internal/load"
	"github.com/vmware/govmomi/vim25/mo"
)

// ResourcePools VMWare
func ResourcePools(config *load.Config) {
	ctx := context.Background()
	m := config.ViewManager

	for i, dc := range config.Datacenters {
		cv, err := m.CreateContainerView(ctx, dc.Datacenter.Reference(), []string{"ResourcePool"}, true)
		if err != nil {
			config.Logrus.WithError(err).Error("failed to create ResourcePool container view")
			continue
		}
		defer cv.Destroy(ctx)
		var resourcePools []mo.ResourcePool
		err = cv.Retrieve(
			ctx,
			[]string{"ResourcePool"},
			[]string{"summary", "owner", "parent", "runtime", "name", "overallStatus", "vm", "resourcePool"},
			&resourcePools)
		if err != nil {
			config.Logrus.WithError(err).Error("failed to retrieve ResourcePools")
			continue
		}
		if err := collectTags(config, resourcePools, config.Datacenters[i]); err != nil {
			config.Logrus.WithError(err).Errorf("failed to retrieve tags:%v", err)
		}

		for j := 0; j < len(resourcePools); j++ {
			config.Datacenters[i].ResourcePools[resourcePools[j].Self] = &resourcePools[j]
		}

	}
}
