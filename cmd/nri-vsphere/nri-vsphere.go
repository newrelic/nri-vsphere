// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-vsphere/internal/client"
	"github.com/newrelic/nri-vsphere/internal/collect"
	"github.com/newrelic/nri-vsphere/internal/load"
	"github.com/newrelic/nri-vsphere/internal/process"
	logrus "github.com/sirupsen/Logrus"
	"github.com/vmware/govmomi/view"
)

var (
	buildVersion = "0.0.0" // set by -ldflags on build
)

func main() {

	config := load.NewConfig(buildVersion)

	err := infraIntegration(config)
	if err != nil {
		config.Logrus.WithError(err).Fatal("failed to initialize integration")
	}
	setupLogger(config)
	if config.Args.Version == true {
		config.Logrus.Infof("integration version: %s", buildVersion)
		return
	}
	config.Logrus.Debugf("integration version: %s", buildVersion)

	checkAndSanitizeConfig(config)

	config.VMWareClient, err = client.New(config.Args.URL, config.Args.User, config.Args.Pass, config.Args.ValidateSSL)
	if err != nil {
		config.Logrus.WithError(err).Fatal("failed to create client")
	}

	if config.VMWareClient.ServiceContent.About.ApiType == "VirtualCenter" {
		config.IsVcenterAPIType = true
	}

	if !config.IsVcenterAPIType && config.Args.EnableVsphereEvents {
		config.Logrus.Fatal("It is not possible to fetch events from the vCenter if the integration is pointing to an host")
	}

	config.ViewManager = view.NewManager(config.VMWareClient.Client)

	runIntegration(config)
}

func checkAndSanitizeConfig(config *load.Config) {
	if config.Args.URL == "" {
		config.Logrus.Fatal("missing argument `url`, please check if URL has been supplied in the config file")
	}
	if config.Args.User == "" {
		config.Logrus.Fatal("missing argument `user`, please check if username has been supplied in the config file")
	}
	if config.Args.Pass == "" {
		config.Logrus.Fatal("missing argument `pass`, please check if password has been supplied")
	}

	if config.Args.EnableVsphereEvents {
		if config.Args.AppDataDir == "" && runtime.GOOS == "windows" {
			config.Logrus.Fatal("missing argument `app_data_dir`, in newer version of the Agent it is injected automatically, please update or specify argument in integration it in config file")
		}

		if config.Args.AgentDir == "" && runtime.GOOS != "windows" {
			config.Logrus.Fatal("missing argument `agent_dir`, in newer version of the Agent it is injected automatically, please update or specify argument in integration config file")
		}
		if runtime.GOOS == "windows" {
			config.CachePath = filepath.Join(config.Args.AppDataDir, "/data/integration/events-cache")
		} else {
			//to test locally in darwin systems you can pass as argument agetn_dir=./ and create te folder "data/integration/events-cache"
			config.CachePath = filepath.Join(config.Args.AgentDir, "/data/integration/events-cache")
		}
	}
	config.Args.DatacenterLocation = strings.ToLower(config.Args.DatacenterLocation)
}

func setupLogger(config *load.Config) {
	verboseLogging := os.Getenv("VERBOSE")
	if config.Args.Verbose || verboseLogging == "true" || verboseLogging == "1" {
		config.Logrus.SetLevel(logrus.TraceLevel)
	}
	config.Logrus.Out = os.Stderr
}

func runIntegration(config *load.Config) {

	collect.CollectData(config)
	process.ProcessData(config)

	err := config.Integration.Publish()
	if err != nil {
		config.Logrus.WithError(err).Fatal("failed to publish")
	}

}

func infraIntegration(config *load.Config) error {
	var err error
	config.Hostname, err = os.Hostname() // set hostname
	if err != nil {
		config.Logrus.WithError(err).Debug("failed to get the hostname while creating integration")
	}

	config.Integration, err = integration.New(config.IntegrationName, config.IntegrationVersion, integration.Args(&config.Args))
	if err != nil {
		return fmt.Errorf("failed to create integration %v", err)
	}
	return nil
}
