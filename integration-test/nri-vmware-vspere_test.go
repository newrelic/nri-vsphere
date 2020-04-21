// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

// +build integration

package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os/exec"
	"reflect"
	"regexp"
	"testing"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/stretchr/testify/assert"
)

func TestMainFunction(t *testing.T) {

	expectedOutput, err := ioutil.ReadFile("./expectedOutput/default")
	if err != nil {
		panic(err)
	}

	actualOutput, _ := exectuteIntegration()
	actual, expected := transformAndSanitizeOutput(string(actualOutput), string(expectedOutput))

	assert.Equal(t, len(actual.Entities), len(expected.Entities), "The number of entities is different to the one expected")

	entitiesNotMatchig := []string{}
	//We cannot trust the order of the slice, I think that it is caused by a map transformed to a slice
	for _, entityActual := range expected.Entities {
		isTheEntityPresentInTheSlice := false
		for _, entity := range actual.Entities {
			test := reflect.DeepEqual(entityActual, entity)
			if test {
				isTheEntityPresentInTheSlice = true
				break
			}
		}
		if isTheEntityPresentInTheSlice == false {
			entitiesNotMatchig = append(entitiesNotMatchig, entityActual.Metadata.Namespace+"    "+entityActual.Metadata.Name)
		}
	}
	assert.Equal(t, []string{}, entitiesNotMatchig, "Some entities are not matching with the mock:\n\n"+string(actualOutput))
}

func exectuteIntegration() ([]byte, []byte) {
	var cmdLine []string
	cmdLine = append(cmdLine, "exec", "-i")
	cmdLine = append(cmdLine, "vmware-integration-with-mock")
	cmdLine = append(cmdLine, "/go/src/github.com/newrelic/nri-vmware-vsphere/bin/linux/nri-vmware-vsphere",
		"-user", "user",
		"-pass", "pass",
		"-url", "127.0.0.1:8989/sdk",
		"-validate_ssl", "false")

	cmd := exec.Command("docker", cmdLine...)

	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	output := outbuf.Bytes()
	error := errbuf.Bytes()

	return output, error
}

func transformAndSanitizeOutput(expectedOutput string, actualOutput string) (integration.Integration, integration.Integration) {
	var expected integration.Integration
	var actual integration.Integration

	re := regexp.MustCompile(Myregex)
	actualOutput = re.ReplaceAllString(actualOutput, "")
	expectedOutput = re.ReplaceAllString(expectedOutput, "")

	err := json.Unmarshal([]byte(actualOutput), &actual)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(expectedOutput), &expected)
	if err != nil {
		panic(err)
	}
	return actual, expected
}

var Myregex = `("timestamp":[0-9]*,|,"timestamp":[0-9]*)`
