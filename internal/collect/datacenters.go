// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package collect

import (
	"context"
	"time"

	"github.com/newrelic/infra-integrations-sdk/persist"
	"github.com/newrelic/nri-vsphere/internal/cache"
	"github.com/newrelic/nri-vsphere/internal/config"
	"github.com/newrelic/nri-vsphere/internal/events"
	"github.com/newrelic/nri-vsphere/internal/model"
	"github.com/vmware/govmomi/vim25/mo"
)

// Datacenters VMWare
func Datacenters(config *config.Config) error {
	ctx := context.Background()
	m := config.ViewManager

	cv, err := m.CreateContainerView(ctx, config.VMWareClient.ServiceContent.RootFolder, []string{DATACENTER}, true)
	if err != nil {
		config.Logrus.WithError(err).Fatal("failed to create Datacenter container view")
	}

	defer func() {
		err := cv.Destroy(ctx)
		if err != nil {
			config.Logrus.WithError(err).Error("error while cleaning up datacenter container view")
		}
	}()

	var datacenters []mo.Datacenter
	err = cv.Retrieve(ctx, []string{DATACENTER}, []string{"name", "overallStatus"}, &datacenters)
	if err != nil {
		config.Logrus.WithError(err).Error("failed to retrieve Datacenters")
		return err
	}

	if config.TagCollectionEnabled() {
		_, err = config.TagCollector.FetchTagsForObjects(datacenters)
		if err != nil {
			config.Logrus.WithError(err).Warn("failed to retrieve tags for datacenters")
		}
	}

	// cache store for events
	cs, err := newCacheStore(config)
	if err != nil {
		config.Logrus.WithError(err).Warn("could not create cache for vsphere events. all events will be returned")
	}

	for _, d := range datacenters {
		if config.TagFilteringEnabled() && !config.TagCollector.MatchObjectTags(d.Reference()) {
			config.Logrus.WithField("datacenter", d.Name).
				Debug("ignoring datacenter since no tags matched the configured filters")
			continue
		}

		newDatacenter := model.NewDatacenter(d)

		if config.EventCollectionEnabled() {
			c := cache.NewCache(d.Name, cs)
			collectEvents(config, d, newDatacenter, c)
		}

		config.Datacenters = append(config.Datacenters, newDatacenter)
	}

	return nil
}

func collectEvents(config *config.Config, d mo.Datacenter, newDatacenter *model.Datacenter, c *cache.Cache) {
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

func newCacheStore(config *config.Config) (persist.Storer, error) {
	// we have to set a distinct default path otherwise it gets overwritten by the default Infra SDK store
	path := persist.DefaultPath(config.IntegrationName + "_timestamps")
	store, err := persist.NewFileStore(path, config.Logrus, time.Hour*24)
	if err != nil {
		store = persist.NewInMemoryStore()
	}
	return store, err
}
