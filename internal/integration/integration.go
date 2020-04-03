/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package integration

import (
	"os"

	"github.com/newrelic/nri-vmware-vsphere/internal/load"
	"github.com/sirupsen/logrus"
)


func SetupLogger(config *load.Config) {
	verboseLogging := os.Getenv("VERBOSE")
	if config.Args.Verbose || verboseLogging == "true" || verboseLogging == "1" {
		config.Logrus.SetLevel(logrus.TraceLevel)
	}
	config.Logrus.Out = os.Stderr
}

