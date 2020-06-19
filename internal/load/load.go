// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package load

import (
	sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
	"github.com/newrelic/infra-integrations-sdk/integration"
	logrus "github.com/sirupsen/Logrus"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/view"
	"time"
)

const (
	minScrapeInterval = 20 * time.Second
	heartBeatPeriod   = 5 * time.Second // Period for the hard beat signal should be less than timeout
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
	EnableEvents       bool   `default:"false" help:"If true the integration will collect as well events"`
	ScrapeInterval     string `default:"60s" help:"Scrape interval of the integration"`
}

type Config struct {
	Args                 ArgumentList
	StartTime            int64                    // StartTime time Flex starts in Nanoseconds
	Integration          *integration.Integration // Integration Infrastructure SDK Integration
	Entity               *integration.Entity      // Entity Infrastructure SDK Entity
	Hostname             string                   // Hostname current host
	Logrus               *logrus.Logger           // Logrus create instance of the logger
	IntegrationName      string                   // IntegrationName name of integration
	IntegrationNameShort string                   // IntegrationNameShort Short Name
	IntegrationVersion   string                   // IntegrationVersion Version
	VMWareClient         *govmomi.Client          // VMWareClient Client
	ViewManager          *view.Manager            // ViewManager Client
	Datacenters          []Datacenter             // Datacenters VMWare
	IsVcenterAPIType     bool                     // IsVcenterAPIType true if connecting to vcenter
	Done                 chan struct{}            // Used to Gracefully end integration
	HeartBeatPeriod      time.Duration            //number of seconds to wait to send heartbeat
	ScrapeInterval       time.Duration            //interval of time to wait before scraping again data

}

func NewConfig() *Config {

	config := Config{
		Logrus:               logrus.New(),
		IntegrationName:      "com.newrelic.vsphere",
		IntegrationNameShort: "vsphere",
		IntegrationVersion:   "Unknown-SNAPSHOT",
		StartTime:            time.Now().UnixNano() / int64(time.Millisecond),
		IsVcenterAPIType:     false,
		Done:                 make(chan struct{}),
		HeartBeatPeriod:      heartBeatPeriod,
	}

	interval, err := time.ParseDuration(config.Args.ScrapeInterval)
	if err != nil {
		config.Logrus.Errorf("error parsing scrape interval:%s `%s`", err.Error(), config.Args.ScrapeInterval)
		interval = minScrapeInterval
	}
	if interval < minScrapeInterval {
		config.Logrus.Warn("scrap interval defined is less than 15s. Interval has set to 15s ")
		interval = minScrapeInterval
	}

	config.ScrapeInterval = interval

	return &config
}
