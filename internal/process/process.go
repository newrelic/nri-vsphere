package process

import (
	// "fmt"

	"strings"
	"sync"
	"time"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-vmware-vsphere/internal/load"
)

// Run process samples
func Run(config *load.Config) {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)

	// create samples async
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		createVirtualMachineSamples(config, timestamp)
	}()
	go func() {
		defer wg.Done()
		createHostSamples(config, timestamp)
	}()
	go func() {
		defer wg.Done()
		createDatastoreSamples(config, timestamp)
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

// setEntity sets the entity to be used for the configured API
// defaults the type aka namespace to instance
func setEntity(config *load.Config, entity string, customNamespace string) *integration.Entity {
	if entity != "" {
		if customNamespace == "" {
			customNamespace = "instance"
		}
		workingEntity, err := config.Integration.Entity(entity, customNamespace)
		if err == nil {
			return workingEntity
		}
	}
	return config.Entity
}

func checkError(config *load.Config, err error) {
	if err != nil {
		config.Logrus.WithError(err).Error("failed to set")
	}
}
