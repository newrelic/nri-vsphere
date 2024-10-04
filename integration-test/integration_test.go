//go:build integration

/*
* Copyright 2020 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/newrelic/nri-vsphere/integration-test/jsonschema"
	"github.com/stretchr/testify/require"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/simulator"
	"github.com/vmware/govmomi/vapi/rest"
	_ "github.com/vmware/govmomi/vapi/simulator"
	"github.com/vmware/govmomi/vapi/tags"
	"github.com/vmware/govmomi/vim25"
)

// This can set when running the test as -ldflags "-X github.com/newrelic/nri-winservices/test.integrationPath="
var (
	integrationPath = "../bin/nri-vsphere"
)

func TestIntegration(t *testing.T) {
	simulator.Test(func(ctx context.Context, vc *vim25.Client) {
		// Add tag to a vm in the simulator
		require.NoError(t, addTag(ctx, vc))

		stdout, stderr, err := runIntegration([]string{
			"-url", vc.URL().String(),
			"-enable_vsphere_tags",
			"-enable_vsphere_events",
		})
		//Notice that stdErr contains as well normal logs of the integration
		require.NotNil(t, stderr, "unexpected stderr")
		require.NoError(t, err, "Unexpected error")

		schemaPath := filepath.Join("json-schema-files", "vsphere-schema.json")
		err = jsonschema.Validate(schemaPath, stdout)
		require.NoError(t, err, "The output of vsphere integration doesn't have expected format")
		require.Less(t, 0, len(stdout), "The output should be longer than 0")

	})
}
func TestIntegrationPerformanceMetrics(t *testing.T) {
	ctx := context.Background()

	model := simulator.VPX()
	// adding a resource pool to the model, default is 0
	model.Pool = 1
	require.NoError(t, model.Create())

	s := model.Service.NewServer()

	vc, err := govmomi.NewClient(ctx, s.URL, true)
	require.NoError(t, err)

	stdout, stderr, err := runIntegration([]string{
		"-url", vc.URL().String(),
		"-enable_vsphere_perf_metrics",
		"-perf_metric_file", "../vsphere-performance.metrics",
		"-perf_level", "4",
	})
	//Notice that stdErr contains as well normal logs of the integration
	require.NotNil(t, stderr, "unexpected stderr")
	require.NoError(t, err, "Unexpected error")

	schemaPath := filepath.Join("json-schema-files", "vsphere-perf-schema.json")
	err = jsonschema.Validate(schemaPath, stdout)
	require.NoError(t, err, "The output of vsphere integration doesn't have expected format")
	require.Less(t, 0, len(stdout), "The output should be longer than 0")

}

func runIntegration(args []string) (string, string, error) {
	defaultArgs := []string{"-user", "user", "-pass", "pass"}
	cmdArgs := append(defaultArgs, args...)
	cmd := exec.Command(
		integrationPath,
		cmdArgs...,
	)

	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	if err != nil {
		path, _ := os.Getwd()
		return "", "", fmt.Errorf("fail to start cmd: %v, currentPath: %s", err, path)
	}
	stdout := outbuf.String()
	stderr := errbuf.String()

	return stdout, stderr, nil
}

func addTag(ctx context.Context, vc *vim25.Client) error {
	c := rest.NewClient(vc)
	_ = c.Login(ctx, simulator.DefaultLogin)

	m := tags.NewManager(c)

	categoryName := "my-category"
	categoryID, err := m.CreateCategory(ctx, &tags.Category{
		AssociableTypes: []string{"VirtualMachine"},
		Cardinality:     "SINGLE",
		Name:            categoryName,
	})
	if err != nil {
		return err
	}
	tagName := "vm-tag"
	tagID, err := m.CreateTag(ctx, &tags.Tag{CategoryID: categoryID, Name: tagName})
	if err != nil {
		return err
	}
	// "DC0_H0_VM0" is the name of a default vm added by the vcsim
	vm, err := find.NewFinder(vc).VirtualMachine(ctx, "DC0_H0_VM0")
	if err != nil {
		return err
	}
	err = m.AttachTag(ctx, tagID, vm.Reference())
	if err != nil {
		return err
	}
	return nil
}
