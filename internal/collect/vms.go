// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package collect

import (
	"context"
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
		defer cv.Destroy(ctx)

		var vms []mo.VirtualMachine
		// Reference: http://pubs.vmware.com/vsphere-60/topic/com.vmware.wssdk.apiref.doc/vim.VirtualMachine.html
		err = cv.Retrieve(
			ctx,
			[]string{"VirtualMachine"},
			[]string{"summary", "network", "config", "guest", "runtime", "resourcePool", "datastore", "overallStatus"},
			&vms)
		if err != nil {
			config.Logrus.WithError(err).Error("failed to retrieve VM Summaries")
			continue
		}

		tagsCategories2, err := config.TagsManager.ListCategories(ctx)
		if err != nil {
			config.Logrus.WithError(err).Error(tagsCategories2)
		}
		config.Logrus.Error("ListOfCategories!")
		config.Logrus.Error(tagsCategories2)

		tagsList2, err := config.TagsManager.ListTags(ctx)
		if err != nil {
			config.Logrus.WithError(err).Error(tagsList2)
		}
		config.Logrus.Error("ListOfTags!")
		config.Logrus.Error(tagsList2)

		config.Logrus.Error("ListOfTags attached to VMs")
		for j := 0; j < len(vms); j++ {
			tagsList, err := config.TagsManager.ListAttachedTags(ctx, vms[j])
			if err != nil {
				config.Logrus.WithError(err).Error(tagsList2)
			}

			for _, t := range tagsList {
				config.Logrus.Error(vms[j].Config.Name)
				tag, err := config.TagsManager.GetTag(ctx, t)
				if err != nil {
					config.Logrus.WithError(err).Error("something weird")
				}
				config.Logrus.Error(tag.Name)
				config.Logrus.Error(tag.CategoryID)
				config.Logrus.Error(tag.Description)
				config.Logrus.Error(tag.UsedBy)
			}

			config.Datacenters[i].VirtualMachines[vms[j].Self] = &vms[j]
		}

	}
}
