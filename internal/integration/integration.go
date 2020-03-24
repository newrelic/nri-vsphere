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

var log = load.Logrus

func setupLogger() {
	verboseLogging := os.Getenv("VERBOSE")
	if load.Args.Verbose || verboseLogging == "true" || verboseLogging == "1" {
		log.SetLevel(logrus.TraceLevel)
	}

	// if load.Args.StructuredLogs {
	// 	log.SetFormatter(&logrus.JSONFormatter{})
	// }
}

// SetDefaults set defaults
func SetDefaults() {
	log.Out = os.Stderr
	// load.FlexStatusCounter.M = make(map[string]int)
	// load.FlexStatusCounter.M["EventCount"] = 0
	// load.FlexStatusCounter.M["EventDropCount"] = 0
	// load.FlexStatusCounter.M["ConfigsProcessed"] = 0
}
