// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package collect

import (
	"context"
	eventSDK "github.com/newrelic/infra-integrations-sdk/data/event"
	"github.com/newrelic/infra-integrations-sdk/integration"
	logrus "github.com/sirupsen/Logrus"
	"github.com/vmware/govmomi/event"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/types"
)

type EventDispacher struct {
	ErrorEvent chan error // Used to rerun event dispatcher if needed
	PageSize   int32      //size of the page
}

func NewEventDispacher() *EventDispacher {
	ed := EventDispacher{
		ErrorEvent: make(chan error),
		PageSize:   100,
	}
	return &ed
}

// Clusters VMWare
func (ed EventDispacher) Events(client *vim25.Client, integration *integration.Integration) {
	go func() {

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		managedTypes := []types.ManagedObjectReference{client.ServiceContent.RootFolder}

		manager := event.NewManager(client)
		//this function should never return
		err := manager.Events(ctx, managedTypes, ed.PageSize, true, false, process(integration))
		ed.ErrorEvent <- err //informing of the event the caller
	}()
}

// process is the name of the function
// `func(types.ManagedObjectReference, []types.BaseEvent) error` this is the return value
// the whole purpose of this is to pass an extra argument
func process(i *integration.Integration) func(types.ManagedObjectReference, []types.BaseEvent) error {
	return func(moref types.ManagedObjectReference, baseEvent []types.BaseEvent) error {
		return processEvent(i, baseEvent)
	}
}

func processEvent(i *integration.Integration, eSlice []types.BaseEvent) error {

	for _, e := range eSlice {
		ev := &eventSDK.Event{
			Summary:  e.GetEvent().FullFormattedMessage,
			Category: "vSphereEvent",
			Attributes: map[string]interface{}{
				"vSphereEvent.userName": e.GetEvent().UserName,
				"vSphereEvent.tag":      e.GetEvent().ChangeTag,
				"vSphereEvent.date":     e.GetEvent().CreatedTime.String(),
			},
		}
		if e.GetEvent().Vm != nil {
			ev.Attributes["vSphereEvent.vm"] = e.GetEvent().Vm.Name
		}
		if e.GetEvent().Host != nil {
			ev.Attributes["vSphereEvent.host"] = e.GetEvent().Host.Name
		}
		if e.GetEvent().Datacenter != nil {
			ev.Attributes["vSphereEvent.datacenter"] = e.GetEvent().Datacenter.Name
		}
		if e.GetEvent().ComputeResource != nil {
			ev.Attributes["vSphereEvent.computeResource"] = e.GetEvent().ComputeResource.Name
		}
		if e.GetEvent().Ds != nil {
			ev.Attributes["vSphereEvent.datastore"] = e.GetEvent().Ds.Name
		}
		err := i.LocalEntity().AddEvent(ev)
		if err != nil {
			logrus.Error()
		}
	}
	return nil
}
