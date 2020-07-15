// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-vsphere/internal/client"
	"github.com/newrelic/nri-vsphere/internal/collect"
	"github.com/newrelic/nri-vsphere/internal/load"
	"github.com/newrelic/nri-vsphere/internal/process"
	logrus "github.com/sirupsen/logrus"
	"github.com/vmware/govmomi/vapi/tags"
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
	if config.Args.Version {
		config.Logrus.Infof("integration version: %s", buildVersion)
		return
	}
	config.Logrus.Debugf("integration version: %s", buildVersion)

	checkAndSanitizeConfig(config)

	config.Logrus.WithField("seconds", time.Since(load.Now).Seconds()).Debug("before creating client")

	config.VMWareClient, err = client.New(config.Args.URL, config.Args.User, config.Args.Pass, config.Args.ValidateSSL)
	if err != nil {
		config.Logrus.WithError(err).Fatal("failed to create client")
	}
	defer func() {
		err := client.Logout(config.VMWareClient)
		if err != nil {
			config.Logrus.WithError(err).Error("error while logging out client")
		}
	}()

	if config.VMWareClient.ServiceContent.About.ApiType == "VirtualCenter" {
		config.IsVcenterAPIType = true
	}

	if !config.IsVcenterAPIType && config.Args.EnableVsphereEvents {
		config.Logrus.Warn("It is not possible to fetch events from the vCenter if the integration is pointing to an host")
	}

	if !config.IsVcenterAPIType && config.Args.EnableVsphereTags {
		config.Logrus.Warn("It is not possible to fetch Tags from the vCenter if the integration is pointing to an host")
	}

	if config.IsVcenterAPIType && config.Args.EnableVsphereTags {
		config.VMWareClientRest, err = client.NewRest(config.VMWareClient, config.Args.User, config.Args.Pass)
		if err != nil {
			config.Logrus.WithError(err).Fatal("failed to create client rest")
		}

		defer func() {
			err := client.LogoutRest(config.VMWareClientRest)
			if err != nil {
				config.Logrus.WithError(err).Error("error while logging out RestClient")
			}
		}()

		config.TagsManager = tags.NewManager(config.VMWareClientRest)
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

	if config.Args.EnableVspherePerfMetrics && config.Args.PerfMetricFile ==""{
		if runtime.GOOS == "windows" {
			config.Args.PerfMetricFile = "C:\\Program Files\\New Relic\\newrelic-infra\\integrations.d\\vsphere-performance.metrics"
		} else {
			config.Args.PerfMetricFile = "/etc/newrelic-infra/integrations.d/vsphere-performance.metrics"
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

	config.Logrus.WithField("seconds", time.Since(load.Now).Seconds()).Debug("before collecting data")
	collect.CollectData(config)
	config.Logrus.WithField("seconds", time.Since(load.Now).Seconds()).Debug("before processing data")
	process.ProcessData(config)
	config.Logrus.WithField("seconds", time.Since(load.Now).Seconds()).Debug("after processing data")

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
