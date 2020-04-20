// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

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
	DatacenterLocation string `default:"" help:"Datacenter Location of your vCenter or ESXi Host eg. sydney-ultimo"`
	ValidateSSL        bool   `default:"false" help:"Validate SSL"`
}

type Config struct {
	Args                        ArgumentList
	StartTime                   int64                    // StartTime time Flex starts in Nanoseconds
	Integration                 *integration.Integration // Integration Infrastructure SDK Integration
	Entity                      *integration.Entity      // Entity Infrastructure SDK Entity
	Hostname                    string                   // Hostname current host
	Logrus                      *logrus.Logger           // Logrus create instance of the logger
	IntegrationName             string                   // IntegrationName name of integration
	IntegrationNameShort        string                   // IntegrationNameShort Short Name
	IntegrationVersion          string                   // IntegrationVersion Version
	VMWareClient                *govmomi.Client          // VMWareClient Client
	ViewManager                 *view.Manager            // ViewManager Client
	HostSystemContainerView     *view.ContainerView      // HostSystemContainerView x
	VirutalMachineContainerView *view.ContainerView      // VirutalMachineContainerView x
	NetworkContainerView        *view.ContainerView      // NetworkContainerView x
	VirtualMachines             []mo.VirtualMachine      // VirtualMachines VMWare
	Networks                    []mo.Network             // Networks VMWare
	Hosts                       []mo.HostSystem          // Hosts VMWare
	Datacenters                 []Datacenter             // Datacenters VMWare
	IsVcenterAPIType            bool                     // IsVcenterAPIType true if connecting to vcenter
}

func NewConfig() *Config {
	return &Config{
		Logrus:               logrus.New(),
		IntegrationName:      "com.newrelic.vmware-vsphere",
		IntegrationNameShort: "vmware-vsphere",
		IntegrationVersion:   "Unknown-SNAPSHOT",
		StartTime:            time.Now().UnixNano() / int64(time.Millisecond),
		IsVcenterAPIType:     false,
	}
}
