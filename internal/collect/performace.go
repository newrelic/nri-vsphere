// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package collect

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/performance"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/types"
	"os"

	"github.com/newrelic/nri-vsphere/internal/load"
)

const maxBatchSizePerf = 30

type perfCollector struct {
	client      *govmomi.Client
	entity      *integration.Entity
	finder      *find.Finder
	perfManager *performance.Manager

	metricDefinition *perfMetrics

	//TODO verify which on is actually needed
	metricsAvaliableByID   map[int32]string
	metricsAvaliableByName map[string]int32
}

//this struct is not needed we can decide to pass more info and process it in the process, it would hide logic
type perfMetric struct {
	value   int64
	counter string
}

func newPerfCollector(config *load.Config) (*perfCollector, error) {

	finder := find.NewFinder(config.VMWareClient.Client, true)

	perfManager := performance.NewManager(config.VMWareClient.Client)

	perfCollector := &perfCollector{
		client:      config.VMWareClient,
		finder:      finder,
		perfManager: perfManager,
	}

	err := perfCollector.retrieveCounterMetadata(config)
	if err != nil {
		config.Logrus.WithError(err).Errorf("failed to fetch available metrics from perfManager")
		return nil, err
	}
	err = perfCollector.parseConfigFile(config.Args.PerfMetricFile)
	if err != nil {
		config.Logrus.WithError(err).Errorf("failed to fetch data from config file")
		return nil, err
	}

	return perfCollector, err
}

func (c *perfCollector) collect(config *load.Config, mos []types.ManagedObjectReference, metrics []types.PerfMetricId) map[types.ManagedObjectReference][]perfMetric {
	ctx := context.Background()
	perfMetricsByRef := map[types.ManagedObjectReference][]perfMetric{}

	query := types.QueryPerf{
		This:      c.perfManager.Reference(),
		QuerySpec: []types.PerfQuerySpec{},
	}

	for i := 0; i < len(mos); i += maxBatchSizePerf {

		chunk := mos[i:min(i+maxBatchSize, len(mos))]

		for _, vm := range chunk {
			querySpec := types.PerfQuerySpec{
				Entity:    vm.Reference(),
				MaxSample: 1,
				MetricId:  metrics,
				//IntervalId: 20,
				//If the optional intervalId is omitted, the metrics are returned in their originally sampled interval.
				//When an intervalId is specified, the server tries to summarize the information for the specified intervalId.
				//However, if that interval does not exist or has no data, the server summarizes the information using the best interval available.
			}
			query.QuerySpec = append(query.QuerySpec, querySpec)
		}

		retrievedStats, err := methods.QueryPerf(ctx, c.perfManager.Client(), &query)
		if err != nil {
			config.Logrus.Error(err)
			continue
		}

		for _, returnVal := range retrievedStats.Returnval {
			metricsValues, ok := returnVal.(*types.PerfEntityMetric)
			if !ok {
				continue
			}
			e := metricsValues.Entity
			c.processEntityMetrics(metricsValues, perfMetricsByRef, e)

		}

	}
	return perfMetricsByRef
}

func (c *perfCollector) processEntityMetrics(metricsValues *types.PerfEntityMetric, perfMetricsByRef map[types.ManagedObjectReference][]perfMetric, e types.ManagedObjectReference) {
	for _, metricValue := range metricsValues.Value {
		metricValueSeries, ok2 := metricValue.(*types.PerfMetricIntSeries)
		if !ok2 {
			continue
		}
		name, ok := c.metricsAvaliableByID[metricValueSeries.Id.CounterId]
		if !ok {
			continue
		}

		perfMetricsByRef[e] = append(perfMetricsByRef[e], perfMetric{
			counter: name,
			value:   metricValueSeries.Value[0], //TODO IT is guarantee only one sample but w should not check this like this
		})

	}
}

func (c *perfCollector) retrieveCounterMetadata(config *load.Config) (err error) {
	ctx := context.Background()

	counters, err := c.perfManager.CounterInfo(ctx)
	c.metricsAvaliableByID = map[int32]string{}
	c.metricsAvaliableByName = map[string]int32{}

	if config.Args.LogAvailableCounters {
		config.Logrus.Info("LogAvailableCounters FLAG ON, printing all %d available counters", len(counters))
	}
	for _, perfCounter := range counters {
		groupInfo := perfCounter.GroupInfo.GetElementDescription()
		nameInfo := perfCounter.NameInfo.GetElementDescription()
		fullCounterName := groupInfo.Key + "." + nameInfo.Key + "." + fmt.Sprint(perfCounter.RollupType)

		c.metricsAvaliableByName[fullCounterName] = perfCounter.Key
		c.metricsAvaliableByID[perfCounter.Key] = fullCounterName

		if config.Args.LogAvailableCounters {
			config.Logrus.Info("\t %s [%d]\n", fullCounterName, perfCounter.Level)
		}
	}
	return nil
}

func (c *perfCollector) parseConfigFile(fileName string) error {

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return fmt.Errorf("error loading configuration from file. Configuration file does not exist")

	}
	var cf configFile

	configFile, err := os.Open(fileName)
	defer configFile.Close()
	if err != nil {
		return err
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&cf)
	if err != nil {
		return err
	}

	c.metricDefinition = &perfMetrics{
		VM:                     c.BuildPerMetricIdSlice(cf.VM),
		ClusterComputeResource: c.BuildPerMetricIdSlice(cf.ClusterComputeResource),
		ResourcePool:           c.BuildPerMetricIdSlice(cf.ResourcePool),
		Datastore:              c.BuildPerMetricIdSlice(cf.Datastore),
		Host:                   c.BuildPerMetricIdSlice(cf.Host),
	}

	return nil
}

func (c *perfCollector) BuildPerMetricIdSlice(slice []string) []types.PerfMetricId {
	var tmp []types.PerfMetricId
	for _, metricName := range slice {
		if counterID, ok := c.metricsAvaliableByName[metricName]; ok {
			pfi := types.PerfMetricId{CounterId: counterID, Instance: "*"}
			tmp = append(tmp, pfi)
		} //todo what we should do if a metric is not available?
	}
	return tmp
}

type perfMetrics struct {
	Host                   []types.PerfMetricId
	VM                     []types.PerfMetricId
	ResourcePool           []types.PerfMetricId
	ClusterComputeResource []types.PerfMetricId
	Datastore              []types.PerfMetricId
}

//This struct is used to parse the config file
type configFile struct {
	Host                   []string `json:"host"`
	VM                     []string `json:"vm"`
	ResourcePool           []string `json:"resourcePool"`
	ClusterComputeResource []string `json:"clusterComputeResource"`
	Datastore              []string `json:"datastore"`
}
