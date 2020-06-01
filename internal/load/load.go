// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package load

import (
	"time"

	sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
	"github.com/newrelic/infra-integrations-sdk/integration"
	logrus "github.com/sirupsen/Logrus"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/view"
)

// ArgumentList Available Arguments
type ArgumentList struct {
	sdkArgs.DefaultArgumentList
	Local               bool   `default:"true" help:"Collect local entity info"`
	Entity              string `default:"" help:"Manually set a remote entity name"`
	URL                 string `default:"" help:"Required: ESXi or vCenter SDK URL eg. https://172.16.53.129/sdk"`
	User                string `default:"" help:"Required: Username"`
	Pass                string `default:"" help:"Required: Password"`
	DatacenterLocation  string `default:"" help:"Datacenter Location of your vCenter or ESXi Host eg. sydney-ultimo"`
	EventsPageSize      string `default:"100" help:"Number of events fetched from the vCenter for each page"`
	EnableVsphereEvents bool   `default:"false" help:"If set the integration will collect as well vSphere events at datacenter level"`
	AgentDir            string `default:"" help:"Agent Directory, injected by agent to save cache in Linux environments, es: /var/db/newrelic-infra" os:"linux"`
	AppDataDir          string `default:"" help:"Agent Data Directory, injected by agent to save cache in Windows environments, es: %PROGRAMDATA%\\New Relic\\newrelic-infra" os:"windows"`
	ValidateSSL         bool   `default:"false" help:"Validate SSL"`
	Version             bool   `default:"false" help:"If set prints version and exit"`
}

type Config struct {
	Args                 ArgumentList
	StartTime            int64                    // StartTime time Flex starts in Nanoseconds
	Integration          *integration.Integration // Integration Infrastructure SDK Integration
	Entity               *integration.Entity      // Entity Infrastructure SDK Entity
	Hostname             string                   // Hostname current host
	Logrus               *logrus.Logger           // Logrus create instance of the logger
	CachePath            string                   // Integration cache path
	IntegrationName      string                   // IntegrationName name of integration
	IntegrationNameShort string                   // IntegrationNameShort Short Name
	IntegrationVersion   string                   // IntegrationVersion Version
	VMWareClient         *govmomi.Client          // VMWareClient Client
	ViewManager          *view.Manager            // ViewManager Client
	Datacenters          []Datacenter             // Datacenters VMWare
	IsVcenterAPIType     bool                     // IsVcenterAPIType true if connecting to vcenter
}

func NewConfig(buildVersion string) *Config {
	return &Config{
		Logrus:               logrus.New(),
		IntegrationName:      "com.newrelic.vsphere",
		IntegrationNameShort: "vsphere",
		IntegrationVersion:   buildVersion,
		StartTime:            time.Now().UnixNano() / int64(time.Millisecond),
		IsVcenterAPIType:     false,
	}
}
