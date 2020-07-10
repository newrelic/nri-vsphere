// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package collect

import (
	"context"
	"github.com/newrelic/infra-integrations-sdk/persist"
	"github.com/newrelic/nri-vsphere/internal/cache"
	"github.com/newrelic/nri-vsphere/internal/events"
	"github.com/newrelic/nri-vsphere/internal/load"
	"github.com/newrelic/nri-vsphere/internal/performance"
	"github.com/vmware/govmomi/vim25/mo"
	"time"
)

// Datacenters VMWare
func Datacenters(config *load.Config) {
	ctx := context.Background()
	m := config.ViewManager

	cv, err := m.CreateContainerView(ctx, config.VMWareClient.ServiceContent.RootFolder, []string{"Datacenter"}, true)
	if err != nil {
		config.Logrus.WithError(err).Fatal("failed to create Datacenter container view")
	}

	defer cv.Destroy(ctx)

	var datacenters []mo.Datacenter
	err = cv.Retrieve(ctx, []string{"Datacenter"}, []string{"name", "overallStatus"}, &datacenters)
	if err != nil {
		config.Logrus.WithError(err).Fatal("failed to retrieve Datacenter")
	}

	if config.Args.EnableVsphereTags && config.IsVcenterAPIType {
		err := collectTagsByID(config.TagsByID, config.TagsManager)
		if err != nil {
			config.Logrus.WithError(err).Error("failed to collect tags by id")
		}
	}

	cs := newCacheStore(config)
	for i, d := range datacenters {
		newDatacenter := load.NewDatacenter(&datacenters[i])
		if config.IsVcenterAPIType && config.Args.EnableVsphereEvents {
			c := cache.NewCache(d.Name, cs)
			collectEvents(config, d, newDatacenter, c)
		}

		if config.Args.EnableVspherePerfMetrics {
			newDatacenter.PerfCollector, err = performance.NewPerfCollector(config.VMWareClient, config.Logrus, config.Args.PerfMetricFile, config.Args.LogAvailableCounters, config.Args.PerfLevel, config.Args.BatchSizePerfEntities, config.Args.BatchSizePerfEntities)
			if err != nil {
				config.Logrus.Fatal(err)
			}
		}
		config.Datacenters = append(config.Datacenters, newDatacenter)

		// create a slice in order to collect tags just for the dc that will be used to store the tags
		dc := []mo.Datacenter{datacenters[i]}
		if err := collectTags(config, dc, config.Datacenters[i]); err != nil {
			config.Logrus.WithError(err).Errorf("failed to retrieve tags:%v", err)
		}

	}
}

func collectEvents(config *load.Config, d mo.Datacenter, newDatacenter *load.Datacenter, c *cache.Cache) {
	//https://pubs.vmware.com/vsphere-51/index.jsp?topic=%2Fcom.vmware.wssdk.apiref.doc%2Fvim.HistoryCollector.html
	ed, err := events.NewEventDispacher(config.VMWareClient.Client, d.Self, config.Logrus, c)
	if err != nil {
		config.Logrus.WithError(err).Error("error while creating event Dispatcher")
		return
	}
	defer ed.Cancel()

	newDatacenter.EventDispacher = ed
	ed.CollectEvents(config.Args.EventsPageSize)
}

func newCacheStore(config *load.Config) persist.Storer {
	path := persist.DefaultPath(config.IntegrationName)
	store, err := persist.NewFileStore(path, config.Logrus, time.Hour*24)
	// if for some reason we couldn't create a file store, use an in-memory store.
	if err != nil {
		config.Logrus.WithError(err).Warn("could not create file based cache for events. using in-memory store")
		store = persist.NewInMemoryStore()
	}
	return store
}
