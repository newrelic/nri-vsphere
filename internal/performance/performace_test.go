// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package performance

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	logrus "github.com/sirupsen/Logrus"
	"github.com/stretchr/testify/require"
)

func TestPerfCollector_parseConfigFile(t *testing.T) {
	c := PerfCollector{
		logger:                 logrus.New(),
		collectionLevel:        2,
		metricsAvaliableByID:   make(map[int32]string),
		metricsAvaliableByName: make(map[string]int32),
	}

	c.metricsAvaliableByName["cpu.coreUtilization.average"] = 1
	c.metricsAvaliableByName["cpu.demand.average"] = 2
	c.metricsAvaliableByName["cpu.outoflevel"] = 3
	content := []byte(`
host:
  level_1:
    - cpu.coreUtilization.average
    - not.considered
  level_2:
    - cpu.demand.average
vm:
  level_1:
    - cpu.demand.average
  level_3:
    - cpu.outoflevel
`)

	tmpfile, err := ioutil.TempFile("", "config")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name()) // clean up
	_, err = tmpfile.Write(content)
	require.NoError(t, err)

	err = c.parseConfigFile(tmpfile.Name())
	require.NoError(t, err)

	// - cpu.costop.summation is discarded since is not in c.metricsAvaliableByName
	assert.Len(t, c.MetricDefinition.Host, 2)

	assert.Len(t, c.MetricDefinition.VM, 1)
}
