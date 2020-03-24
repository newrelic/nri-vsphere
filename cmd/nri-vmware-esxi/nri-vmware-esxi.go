/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"strings"
	"sync"

	"github.com/kav91/nri-vmware-esxi/internal/client"
	"github.com/kav91/nri-vmware-esxi/internal/collect"
	"github.com/kav91/nri-vmware-esxi/internal/integration"
	"github.com/kav91/nri-vmware-esxi/internal/load"
	"github.com/kav91/nri-vmware-esxi/internal/outputs"
	"github.com/kav91/nri-vmware-esxi/internal/process"
	"github.com/vmware/govmomi/view"
)

func main() {
	load.StartTime = load.MakeTimestamp()

	err := outputs.InfraIntegration()
	if err != nil {
		load.Logrus.WithError(err).Fatal("failed to initialize integration")
	}

	if load.Args.URL == "" || load.Args.User == "" || load.Args.Pass == "" {
		load.Logrus.Fatal("missing argument, please check if URL, User, Pass has been supplied")
	}
	load.Args.DatacenterLocation = strings.ToLower(load.Args.DatacenterLocation)

	integration.SetDefaults()
	load.VMWareClient, err = client.New(load.Args.URL, load.Args.User, load.Args.Pass, load.Args.ValidateSSL)
	if err != nil {
		load.Logrus.WithError(err).Fatal("failed to create client")
	}

	load.ViewManager = view.NewManager(load.VMWareClient.Client)

	// fetch vmware data async
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		collect.VirtualMachines(load.VMWareClient)
	}()
	go func() {
		defer wg.Done()
		collect.Networks(load.VMWareClient)

	}()
	go func() {
		defer wg.Done()
		collect.Hosts(load.VMWareClient)

	}()
	wg.Wait()

	process.Run()

	err = load.Integration.Publish()
	if err != nil {
		load.Logrus.WithError(err).Fatal("failed to publish")
	}
}
