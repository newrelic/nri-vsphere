/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"strings"
	"sync"

	"github.com/newrelic/nri-vmware-vsphere/internal/client"
	"github.com/newrelic/nri-vmware-vsphere/internal/collect"
	"github.com/newrelic/nri-vmware-vsphere/internal/integration"
	"github.com/newrelic/nri-vmware-vsphere/internal/load"
	"github.com/newrelic/nri-vmware-vsphere/internal/outputs"
	"github.com/newrelic/nri-vmware-vsphere/internal/process"
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
	integration.SetupLogger(config)

	config.VMWareClient, err = client.New(config.Args.URL, config.Args.User, config.Args.Pass, config.Args.ValidateSSL)
	if err != nil {
		config.Logrus.WithError(err).Fatal("failed to create client")
	}

	if config.VMWareClient.ServiceContent.About.ApiType == "VirtualCenter" {
		config.IsVcenterAPIType = true
	}

	config.ViewManager = view.NewManager(config.VMWareClient.Client)

	collect.Datacenters(config)

	// fetch vmware data async
	var wg sync.WaitGroup
	wg.Add(6)
	go func() {
		defer wg.Done()
		collect.VirtualMachines(config)
	}()
	go func() {
		defer wg.Done()
		collect.Networks(config)

	}()
	go func() {
		defer wg.Done()
		collect.Hosts(config)

	}()
	go func() {
		defer wg.Done()
		collect.Datastores(config)

	}()
	go func() {
		defer wg.Done()
		collect.Clusters(config)
	}()
	go func() {
		defer wg.Done()
		collect.ResourcePools(config)
	}()
	wg.Wait()

	process.Run(config)

	err = config.Integration.Publish()
	if err != nil {
		config.Logrus.WithError(err).Fatal("failed to publish")
	}
}
