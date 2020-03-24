package process

import (
	// "fmt"

	"strings"
	"sync"

	"github.com/kav91/nri-vmware-esxi/internal/load"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

// Run process samples
func Run() {
	timestamp := load.MakeTimestamp()

	// create samples async
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		createVirtualMachineSamples(timestamp)
	}()
	go func() {
		defer wg.Done()
		createHostSamples(timestamp)
	}()
	wg.Wait()
}

func findHost(hostReference types.ManagedObjectReference) mo.HostSystem {
	host := mo.HostSystem{}
	for _, h := range load.Hosts {
		if h.Reference() == hostReference {
			return h
		}
	}
	return host
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
func setEntity(entity string, customNamespace string) *integration.Entity {
	if entity != "" {
		if customNamespace == "" {
			customNamespace = "instance"
		}
		workingEntity, err := load.Integration.Entity(entity, customNamespace)
		if err == nil {
			return workingEntity
		}
	}
	return load.Entity
}

func checkError(err error) {
	if err != nil {
		load.Logrus.WithError(err).Error("failed to set")
	}
}
