/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package load

import (
	"time"

	sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
	"github.com/newrelic/infra-integrations-sdk/integration"
	logrus "github.com/sirupsen/logrus"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
)

// ArgumentList Available Arguments
type ArgumentList struct {
	sdkArgs.DefaultArgumentList
	Local              bool   `default:"true" help:"Collect local entity info"`
	Entity             string `default:"" help:"Manually set a remote entity name"`
	URL                string `default:"" help:"Required: ESXi or vCenter SDK URL eg. https://172.16.53.129/sdk"`
	User               string `default:"" help:"Required: Username"`
	Pass               string `default:"" help:"Required: Password"`
	DatacenterLocation string `default:"false" help:"Required: Datacenter Location of your vCenter or ESXi Host eg. sydney-ultimo"`
	ValidateSSL        bool   `default:"false" help:"Validate SSL"`
}

// Args Infrastructure SDK Arguments List
var Args ArgumentList

// StartTime time Flex starts in Nanoseconds
var StartTime int64

// Integration Infrastructure SDK Integration
var Integration *integration.Integration

// Entity Infrastructure SDK Entity
var Entity *integration.Entity

// Hostname current host
var Hostname string

// Logrus create instance of the logger
var Logrus = logrus.New()

// IntegrationName name of integration
var IntegrationName = "com.newrelic.nri-vmware-vsphere"

// IntegrationNameShort Short Name
var IntegrationNameShort = "nri-vmware-vsphere"

// IntegrationVersion Version
var IntegrationVersion = "Unknown-SNAPSHOT"

// VMWareClient Client
var VMWareClient *govmomi.Client

// ViewManager Client
var ViewManager *view.Manager

// HostSystemContainerView x
var HostSystemContainerView *view.ContainerView

// NetworkContainerView x
var NetworkContainerView *view.ContainerView

// Networks VMWare
var Networks []mo.Network

// Hosts VMWare
var Hosts []mo.HostSystem

// Datacenters VMWare
var Datacenters []Datacenter

// MakeTimestamp creates timestamp in milliseconds
func MakeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
