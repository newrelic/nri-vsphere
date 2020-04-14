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
func InfraIntegration(config *load.Config) error {
	var err error
	config.Hostname, err = os.Hostname() // set hostname
	if err != nil {
		config.Logrus.
			WithError(err).
			Debug("failed to get the hostname while creating integration")
	}

	config.Integration, err = Integration.New(config.IntegrationName, config.IntegrationVersion, Integration.Args(&config.Args))
	if err != nil {
		return fmt.Errorf("failed to create integration %v", err)
	}
	return nil
}

