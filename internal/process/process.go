// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package process

import (
	"strings"
	"sync"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-vsphere/internal/config"

	logrus "github.com/sirupsen/logrus"
)

const (
	entityTypeDatacenter   = "Datacenter"
	entityTypeCluster      = "Cluster"
	entityTypeHost         = "Host"
	entityTypeResourcePool = "ResourcePool"
	entityTypeVm           = "Vm"
	entityTypeDatastore    = "Datastore"
	//The sampleTypeSnapshotVm is used to create a sample, however it does not have a corresponding entity
	//sampleTypeSnapshotVm is attached to a vm entity.
	sampleTypeSnapshotVm = "SnapshotVm"

	tagsPrefix       = "label."
	perfMetricPrefix = "perf."
)

// Run process samples
func ProcessData(config *config.Config) {
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

func checkError(logger *logrus.Logger, err error) {
	if err != nil {
		logger.WithError(err).Error("failed to set")
	}
}

func sanitizeEntityName(config *config.Config, entityName string, datacenterName string) string {
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

func createNewEntityWithMetricSet(config *config.Config, typeEntity string, entityName string, uniqueIdentifier string) (*integration.Entity, *metric.Set, error) {
	workingEntity, err := config.Integration.Entity(uniqueIdentifier, "vsphere-"+strings.ToLower(typeEntity))
	if err != nil {
		config.Logrus.WithError(err).Error("failed to create entity")
		return nil, nil, err
	}

	// entity displayName
	checkError(config.Logrus, workingEntity.SetInventoryItem("vsphere"+typeEntity, "name", entityName))
	ms := workingEntity.NewMetricSet("VSphere" + typeEntity + "Sample")
	return workingEntity, ms, nil
}
