/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package outputs

import (
	"fmt"
	"os"

	Integration "github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-vmware-vsphere/internal/load"
)

// InfraIntegration Creates Infrastructure SDK Integration
func InfraIntegration() error {
	var err error
	load.Hostname, err = os.Hostname() // set hostname
	if err != nil {
		load.Logrus.
			WithError(err).
			Debug("failed to get the hostname while creating integration")
	}

	load.Integration, err = Integration.New(load.IntegrationName, load.IntegrationVersion, Integration.Args(&load.Args))
	if err != nil {
		return fmt.Errorf("failed to create integration %v", err)
	}

	load.Entity, err = createEntity(load.Args.Local, load.Args.Entity)
	if err != nil {
		return fmt.Errorf("failed create entity: %v", err)
	}
	return nil
}

func createEntity(isLocalEntity bool, entityName string) (*Integration.Entity, error) {
	if isLocalEntity {
		return load.Integration.LocalEntity(), nil
	}

	if entityName == "" {
		entityName = load.Hostname // default hostname
	}

	return load.Integration.Entity(entityName, load.IntegrationNameShort)
}
