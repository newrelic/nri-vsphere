// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package load

import (
	"time"

	sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
	"github.com/newrelic/infra-integrations-sdk/integration"
	logrus "github.com/sirupsen/logrus"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vapi/rest"
	"github.com/vmware/govmomi/vapi/tags"
	"github.com/vmware/govmomi/view"
)

var Now = time.Now()

// ArgumentList Available Arguments
type ArgumentList struct {
	sdkArgs.DefaultArgumentList
	Local              bool   `default:"true" help:"Collect local entity info"`
	Entity             string `default:"" help:"Manually set a remote entity name"`
	URL                string `default:"" help:"Required: ESXi or vCenter SDK URL eg. https://172.16.53.129/sdk"`
	User               string `default:"" help:"Required: Username"`
	Pass               string `default:"" help:"Required: Password"`
	DatacenterLocation string `default:"" help:"Datacenter Location of your vCenter or ESXi Host eg. sydney-ultimo"`

	EnableVsphereEvents bool   `default:"false" help:"If set the integration will collect as well vSphere events at datacenter level"`
	EventsPageSize      string `default:"100" help:"Number of events fetched from the vCenter for each page"`
	AgentDir            string `default:"" help:"Agent Directory, injected by agent to save cache in Linux environments, es: /var/db/newrelic-infra" os:"linux"`
	AppDataDir          string `default:"" help:"Agent Data Directory, injected by agent to save cache in Windows environments, es: %PROGRAMDATA%\\New Relic\\newrelic-infra" os:"windows"`

	EnableVspherePerfMetrics bool   `default:"false" help:"If set the integration will collect as well vSphere performance metrics"`
	PerfLevel                int    `default:"1" help:"Performance counter level that will be collected"`
	LogAvailableCounters     bool   `default:"false" help:"Print available performance metrics"`
	PerfMetricFile           string `default:"" help:"location of the configuration file containing perfMetrics to be retrieved"`

	//As a general rule, specify between 10 and 50 entities in a single call to the QueryPerf method.
	//This is a general recommendation because your system configuration may impose different
	//constraints.
	//https://vdc-download.vmware.com/vmwb-repository/dcr-public/cdbbd51c-4824-4a1b-ad43-45df55a76a76/8cb3ed93-cac2-46aa-b329-db5a096af5bc/vsphere-web-services-sdk-67-programming-guide.pdf
	BatchSizePerfEntities string `default:"50" help:"Number of entities requested at the same time when querying perf metrics"`
	BatchSizePerfMetrics  string `default:"50" help:"Number of metrics requested at the same time when querying perf metrics"`

	EnableVsphereTags      bool `default:"false" help:"If true tags will be collected. Tags are available when connecting to vcenter"`
	EnableVsphereSnapshots bool `default:"false" help:"If set to true integration will collect, process and send as well data regarding vm Snapshots"`

	ValidateSSL bool `default:"false" help:"Validate SSL"`
	Version     bool `default:"false" help:"If set prints version and exit"`
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
	VMWareClientRest     *rest.Client             // VMWareClientRest Client
	ViewManager          *view.Manager            // ViewManager Client
	TagsManager          *tags.Manager            // TagsManager Client
	Datacenters          []*Datacenter            // Datacenters VMWare
	TagsByID             TagsByID                 // Lists of tags by id
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
		TagsByID:             make(map[string]Tag),
	}
}
