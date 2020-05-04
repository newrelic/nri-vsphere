// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package process

import (
	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/nri-vsphere/internal/load"
	"strings"
	"sync"
)

const (
	entityTypeDatacenter   = "Datacenter"
	entityTypeCluster      = "Cluster"
	entityTypeHost         = "Host"
	entityTypeResourcePool = "ResourcePool"
	entityTypeVm           = "Vm"
	entityTypeDatastore    = "Datastore"
)

// Run process samples
func Run(config *load.Config) {

	// create samples async
	var wg sync.WaitGroup
	wg.Add(6)
	go func() {
		defer wg.Done()
		createVirtualMachineSamples(config)
	}()
	go func() {
		defer wg.Done()
		createHostSamples(config)
	}()
	go func() {
		defer wg.Done()
		createDatastoreSamples(config)
	}()
	go func() {
		defer wg.Done()
		createDatacenterSamples(config)
	}()
	go func() {
		defer wg.Done()
		createClusterSamples(config)
	}()
	go func() {
		defer wg.Done()
		createResourcePoolSamples(config)
	}()
	wg.Wait()
}

// determineOS perform best effor to determine the operatingSystem
func determineOS(guestFullName string) string {

	if guestFullName != "" {
		gfnLower := strings.ToLower(guestFullName)

		linuxChecks := []string{"nix", "nux", "centos", "red", "suse"}
		for _, check := range linuxChecks {
			if strings.Contains(gfnLower, check) {
				return "linux"
			}
		}

		otherChecks := []string{"bsd", "solaris"}
		for _, check := range otherChecks {
			if strings.Contains(gfnLower, check) {
				return "unix"
			}
		}

		winChecks := []string{"win"}
		for _, check := range winChecks {
			if strings.Contains(gfnLower, check) {
				return "windows"
			}
		}

		macChecks := []string{"mac"}
		for _, check := range macChecks {
			if strings.Contains(gfnLower, check) {
				return "mac"
			}
		}

	}

	return "unknown"
}

func checkError(config *load.Config, err error) {
	if err != nil {
		config.Logrus.WithError(err).Error("failed to set")
	}
}

func sanitizeEntityName(config *load.Config, entityName string, datacenterName string) string {
	if config.IsVcenterAPIType && (datacenterName != "") {
		entityName = datacenterName + ":" + entityName
	}

	if config.Args.DatacenterLocation != "" {
		entityName = config.Args.DatacenterLocation + ":" + entityName
	}

	entityName = strings.ToLower(entityName)
	entityName = strings.Replace(entityName, ".", "-", -1)
	return entityName
}

func createNewEntityWithMetricSet(config *load.Config, typeEntity string, entityName string, uniqueIdentifier string) *metric.Set {
	workingEntity, err := config.Integration.Entity(uniqueIdentifier, "vsphere-"+strings.ToLower(typeEntity))
	if err != nil {
		config.Logrus.WithError(err).Error("failed to create entity")
	}

	// entity displayName
	workingEntity.SetInventoryItem("vsphere"+typeEntity, "name", entityName)
	ms := workingEntity.NewMetricSet("VSphere" + typeEntity + "Sample")
	return ms
}
