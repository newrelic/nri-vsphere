// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"github.com/newrelic/infra-integrations-sdk/log"
	logrus "github.com/sirupsen/Logrus"
	"os"
	"strings"
	"time"

	"github.com/newrelic/nri-vsphere/internal/client"
	"github.com/newrelic/nri-vsphere/internal/collect"
	"github.com/newrelic/nri-vsphere/internal/load"
	"github.com/newrelic/nri-vsphere/internal/outputs"
	"github.com/newrelic/nri-vsphere/internal/process"
	"github.com/vmware/govmomi/view"
)

func main() {
	config := load.NewConfig()

	err := outputs.InfraIntegration(config)
	if err != nil {
		config.Logrus.WithError(err).Fatal("failed to initialize integration")
	}

	if config.Args.URL == "" || config.Args.User == "" || config.Args.Pass == "" {
		config.Logrus.Fatal("missing argument, please check if URL, User, Pass has been supplied")
	}
	config.Args.DatacenterLocation = strings.ToLower(config.Args.DatacenterLocation)

	setupLogger(config)

	config.VMWareClient, err = client.New(config.Args.URL, config.Args.User, config.Args.Pass, config.Args.ValidateSSL)
	if err != nil {
		config.Logrus.WithError(err).Fatal("failed to create client")
	}

	if config.VMWareClient.ServiceContent.About.ApiType == "VirtualCenter" {
		config.IsVcenterAPIType = true
	}

	config.ViewManager = view.NewManager(config.VMWareClient.Client)

	runIntegration(config)
}

func setupLogger(config *load.Config) {
	verboseLogging := os.Getenv("VERBOSE")
	if config.Args.Verbose || verboseLogging == "true" || verboseLogging == "1" {
		config.Logrus.SetLevel(logrus.TraceLevel)
	}
	config.Logrus.Out = os.Stderr
}

func runIntegration(config *load.Config) {

	heartBeat := time.NewTicker(config.HeartBeatPeriod)
	metricInterval := time.NewTicker(config.ScrapeInterval)
	eventDispacher := collect.NewEventDispacher()

	eventDispacher.Events(config.VMWareClient.Client, config.Integration) //this never ends if no error or connection issue occurs
	for {
		select {
		case <-heartBeat.C:
			log.Debug("Sending heartBeat")
			// hart beat signal for long running integrations
			// https://docs.newrelic.com/docs/integrations/integrations-sdk/file-specifications/host-integrations-newer-configuration-format#timeout
			fmt.Println("{}")

		case <-metricInterval.C:
			collect.CollectData(config)
			process.ProcessData(config)

			err := config.Integration.Publish()
			if err != nil {
				config.Logrus.WithError(err).Fatal("failed to publish")
			}

		case err := <-eventDispacher.ErrorEvent:
			if err != nil {
				config.Logrus.WithError(err).Error("The event dispatcher returned")
			}
			eventDispacher.Events(config.VMWareClient.Client, config.Integration) //this never ends if no error or connection issue occurs

		}
	}
}
