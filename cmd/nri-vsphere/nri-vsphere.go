// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-vsphere/internal/client"
	"github.com/newrelic/nri-vsphere/internal/collect"
	"github.com/newrelic/nri-vsphere/internal/config"
	"github.com/newrelic/nri-vsphere/internal/model/tag"
	"github.com/newrelic/nri-vsphere/internal/process"

	"github.com/sirupsen/logrus"
	"github.com/vmware/govmomi/vapi/tags"
	"github.com/vmware/govmomi/view"
)

var (
	buildVersion = "0.0.0" // set by -ldflags on build
)

func main() {

	cfg := config.New(buildVersion)

	err := infraIntegration(cfg)
	if err != nil {
		cfg.Logrus.WithError(err).Fatal("failed to initialize integration")
	}
	setupLogger(cfg)

	// print integration version and exit
	if cfg.Args.Version {
		cfg.Logrus.Infof("integration version: %s", buildVersion)
		return
	}
	cfg.Logrus.Debugf("integration version: %s", buildVersion)

	checkAndSanitizeConfig(cfg)

	cfg.VMWareClient, err = client.New(cfg.Args.URL, cfg.Args.User, cfg.Args.Pass, cfg.Args.ValidateSSL)
	if err != nil {
		cfg.Logrus.WithError(err).Fatal("failed to create client")
	}
	defer func() {
		err := client.Logout(cfg.VMWareClient)
		if err != nil {
			cfg.Logrus.WithError(err).Error("error while logging out client")
		}
	}()

	if cfg.VMWareClient.ServiceContent.About.ApiType == "VirtualCenter" {
		cfg.IsVcenterAPIType = true
	}

	if cfg.IsVcenterAPIType && cfg.Args.EnableVsphereTags {
		cfg.VMWareClientRest, err = client.NewRest(cfg.VMWareClient, cfg.Args.User, cfg.Args.Pass)
		if err != nil {
			cfg.Logrus.WithError(err).Fatal("failed to create client rest")
		}

		defer func() {
			err := client.LogoutRest(cfg.VMWareClientRest)
			if err != nil {
				cfg.Logrus.WithError(err).Error("error while logging out RestClient")
			}
		}()

		cfg.TagsManager = tags.NewManager(cfg.VMWareClientRest)
		if len(cfg.Args.IncludeTags) > 0 {
			tag.ParseFilterTagExpression(cfg.Args.IncludeTags)
		}
	}

	cfg.ViewManager = view.NewManager(cfg.VMWareClient.Client)

	runIntegration(cfg)

}

func checkAndSanitizeConfig(cfg *config.Config) {
	if cfg.Args.URL == "" {
		cfg.Logrus.Fatal("missing argument `url`, please check if URL has been supplied in the config file")
	}
	if cfg.Args.User == "" {
		cfg.Logrus.Fatal("missing argument `user`, please check if username has been supplied in the config file")
	}
	if cfg.Args.Pass == "" {
		cfg.Logrus.Fatal("missing argument `pass`, please check if password has been supplied")
	}

	if !cfg.IsVcenterAPIType && cfg.Args.EnableVsphereEvents {
		cfg.Logrus.Warn("It is not possible to fetch events from the vCenter if the integration is pointing to an host")
	}
	if !cfg.IsVcenterAPIType && cfg.Args.EnableVsphereTags {
		cfg.Logrus.Warn("It is not possible to fetch Tags from the vCenter if the integration is pointing to an host")
	}

	if cfg.Args.EnableVspherePerfMetrics && cfg.Args.PerfMetricFile == "" {
		var err error
		if runtime.GOOS == "windows" {
			cfg.Args.PerfMetricFile, err = filepath.Abs(config.WindowsPerfMetricFile)
		} else {
			cfg.Args.PerfMetricFile, err = filepath.Abs(config.LinuxDefaultPerfMetricFile)
		}
		if err != nil {
			cfg.Logrus.Fatal("error while setting default path for performance metrics configuration file")
		}
	}

	cfg.Args.DatacenterLocation = strings.ToLower(cfg.Args.DatacenterLocation)
}

func setupLogger(config *config.Config) {
	verboseLogging := os.Getenv("VERBOSE")
	if config.Args.Verbose || verboseLogging == "true" || verboseLogging == "1" {
		config.Logrus.SetLevel(logrus.TraceLevel)
	}
	config.Logrus.Out = os.Stderr
}

func runIntegration(config *config.Config) {
	now := time.Now()

	config.Logrus.WithField("seconds", time.Since(now).Seconds()).Debug("before collecting data")
	collect.CollectData(config)
	config.Logrus.WithField("seconds", time.Since(now).Seconds()).Debug("before processing data")
	process.ProcessData(config)
	config.Logrus.WithField("seconds", time.Since(now).Seconds()).Debug("after processing data")

	err := config.Integration.Publish()
	if err != nil {
		config.Logrus.WithError(err).Fatal("failed to publish")
	}

}

func infraIntegration(config *config.Config) error {
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
