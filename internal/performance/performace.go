// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package performance

import (
	"context"
	"encoding/json"
	"fmt"
	logrus "github.com/sirupsen/Logrus"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/performance"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/types"
	"os"
	"strconv"
)

const maxBatchSizePerfEntities = 100
const maxBatchSizePerfMetrics = 100

//As a general rule, specify between 10 and 50 entities in a single call to the QueryPerf method.
//This is a general recommendation because your system configuration may impose different
//constraints.
//https://vdc-download.vmware.com/vmwb-repository/dcr-public/cdbbd51c-4824-4a1b-ad43-45df55a76a76/8cb3ed93-cac2-46aa-b329-db5a096af5bc/vsphere-web-services-sdk-67-programming-guide.pdf

type PerfCollector struct {
	client           *govmomi.Client
	perfManager      *performance.Manager
	logger           *logrus.Logger
	MetricDefinition *perfMetricsIDs

	metricsAvaliableByID   map[int32]string
	metricsAvaliableByName map[string]int32
}

//this struct is not needed we can decide to pass more info and process it in the process, it would hide logic
type PerfMetric struct {
	Value   int64
	Counter string
}

func NewPerfCollector(client *govmomi.Client, logger *logrus.Logger, perfMetricFile string, logAvailableCounters bool) (*PerfCollector, error) {

	perfManager := performance.NewManager(client.Client)

	perfCollector := &PerfCollector{
		client:      client,
		perfManager: perfManager,
		logger:      logger,
	}

	err := perfCollector.retrieveCounterMetadata(logAvailableCounters)
	if err != nil {
		logger.WithError(err).Errorf("failed to fetch available metrics from perfManager")
		return nil, err
	}
	err = perfCollector.parseConfigFile(perfMetricFile)
	if err != nil {
		logger.WithError(err).Errorf("failed to fetch data from config file")
		return nil, err
	}

	return perfCollector, err
}

func (c *PerfCollector) Collect(mos []types.ManagedObjectReference, metrics []types.PerfMetricId, batchSizePerfEntitiesString string, batchSizePerfMetricsString string) map[types.ManagedObjectReference][]PerfMetric {
	ctx := context.Background()
	perfMetricsByRef := map[types.ManagedObjectReference][]PerfMetric{}

	batchSizePerfEntities, batchSizePerfMetrics, err := sanitizeArgs(batchSizePerfEntitiesString, c.logger, batchSizePerfMetricsString)
	if err != nil {
		return nil
	}

	query := types.QueryPerf{
		This:      c.perfManager.Reference(),
		QuerySpec: []types.PerfQuerySpec{},
	}

	for i := 0; i < len(mos); i += batchSizePerfEntities {
		for m := 0; m < len(metrics); m += batchSizePerfMetrics {

			chunkEntities := mos[i:min(i+batchSizePerfEntities, len(mos))]
			chunkMetrics := metrics[m:min(m+batchSizePerfMetrics, len(metrics))]

			for _, vm := range chunkEntities {
				querySpec := types.PerfQuerySpec{
					Entity:     vm.Reference(),
					MaxSample:  1,
					MetricId:   chunkMetrics,
					IntervalId: 20,
					//If the optional intervalId is omitted, the metrics are returned in their originally sampled interval.
					//When an intervalId is specified, the server tries to summarize the information for the specified intervalId.
					//However, if that interval does not exist or has no data, the server summarizes the information using the best interval available.
				}
				query.QuerySpec = append(query.QuerySpec, querySpec)
			}

			c.logger.WithField("number of entities", len(query.QuerySpec)).Debug("querying for perf metrics")
			retrievedStats, err := methods.QueryPerf(ctx, c.perfManager.Client(), &query)
			if err != nil {
				c.logger.Error(err)
				continue
			}

			for _, returnVal := range retrievedStats.Returnval {
				metricsValues, ok := returnVal.(*types.PerfEntityMetric) //TODO IT is guarantee
				if !ok || metricsValues == nil {
					continue
				}
				e := metricsValues.Entity
				c.processEntityMetrics(metricsValues, perfMetricsByRef, e)
			}
		}
	}
	return perfMetricsByRef
}

func (c *PerfCollector) processEntityMetrics(metricsValues *types.PerfEntityMetric, perfMetricsByRef map[types.ManagedObjectReference][]PerfMetric, e types.ManagedObjectReference) {
	for _, metricValue := range metricsValues.Value {
		metricValueSeries, ok2 := metricValue.(*types.PerfMetricIntSeries) //TODO IT is guarantee
		if !ok2 || metricValueSeries == nil {
			continue
		}
		name, ok := c.metricsAvaliableByID[metricValueSeries.Id.CounterId]
		if !ok {
			continue
		}

		if len(metricValueSeries.Value) != 1 {
			c.logger.Debug("The metrics is not containing one sample, this is not expected")
			continue
		}

		perfMetricsByRef[e] = append(perfMetricsByRef[e], PerfMetric{
			Counter: name,
			Value:   metricValueSeries.Value[0],
		})

	}
}

func (c *PerfCollector) retrieveCounterMetadata(logAvailableCounters bool) (err error) {
	ctx := context.Background()

	counters, err := c.perfManager.CounterInfo(ctx)
	c.metricsAvaliableByID = map[int32]string{}
	c.metricsAvaliableByName = map[string]int32{}

	if logAvailableCounters {
		c.logger.Info("LogAvailableCounters FLAG ON, printing all %d available counters", len(counters))
	}
	for _, perfCounter := range counters {
		groupInfo := perfCounter.GroupInfo.GetElementDescription()
		nameInfo := perfCounter.NameInfo.GetElementDescription()
		fullCounterName := groupInfo.Key + "." + nameInfo.Key + "." + fmt.Sprint(perfCounter.RollupType)

		c.metricsAvaliableByName[fullCounterName] = perfCounter.Key
		c.metricsAvaliableByID[perfCounter.Key] = fullCounterName

		if logAvailableCounters {
			c.logger.Info("\t %s [%d]\n", fullCounterName, perfCounter.Level)
		}
	}
	return nil
}

func (c *PerfCollector) parseConfigFile(fileName string) error {

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

	c.MetricDefinition = &perfMetricsIDs{
		VM:                     c.BuildPerMetricIdSlice(cf.VM),
		ClusterComputeResource: c.BuildPerMetricIdSlice(cf.ClusterComputeResource),
		ResourcePool:           c.BuildPerMetricIdSlice(cf.ResourcePool),
		Datastore:              c.BuildPerMetricIdSlice(cf.Datastore),
		Host:                   c.BuildPerMetricIdSlice(cf.Host),
	}

	return nil
}

func (c *PerfCollector) BuildPerMetricIdSlice(slice []string) []types.PerfMetricId {
	var tmp []types.PerfMetricId
	for _, metricName := range slice {
		if counterID, ok := c.metricsAvaliableByName[metricName]; ok {
			//““ – A string of length zero directs the vSphere Server to return only aggregated instance
			//data or rollup type data
			//https://vdc-download.vmware.com/vmwb-repository/dcr-public/cdbbd51c-4824-4a1b-ad43-45df55a76a76/8cb3ed93-cac2-46aa-b329-db5a096af5bc/vsphere-web-services-sdk-67-programming-guide.pdf
			pfi := types.PerfMetricId{CounterId: counterID, Instance: ""}

			tmp = append(tmp, pfi)
		} else {
			c.logger.WithField("metricName", metricName).Debug("metric not available")
		}
	}
	return tmp
}

type perfMetricsIDs struct {
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func sanitizeArgs(batchSizePerfEntitiesString string, logger *logrus.Logger, batchSizePerfMetricsString string) (int, int, error) {
	batchSizePerfEntities, err := strconv.Atoi(batchSizePerfEntitiesString)
	if err != nil {
		logger.WithError(err).Warn("Failed to parse batchSizePerf flag")
		return 0, 0, err
	}

	if batchSizePerfEntities > maxBatchSizePerfEntities {
		batchSizePerfEntities = maxBatchSizePerfEntities
		logger.WithField("maxBatchSizePerfEntities", maxBatchSizePerfEntities).Warn("maxBatchSizePerfEntities above the maximum, setting it to the maximum")
	} else if batchSizePerfEntities < 0 {
		batchSizePerfEntities = 1
		logger.WithField("min size", 1).Warn("batchSizePerf less then 0 no allowed, setting it to 1")
	}

	batchSizePerfMetrics, err := strconv.Atoi(batchSizePerfMetricsString)
	if err != nil {
		logger.WithError(err).Warn("Failed to parse batchSizePerf flag")
		return 0, 0, err
	}

	if batchSizePerfMetrics > maxBatchSizePerfMetrics {
		batchSizePerfMetrics = maxBatchSizePerfMetrics
		logger.WithField("maxBatchSizePerfMetrics", maxBatchSizePerfMetrics).Warn("maxBatchSizePerfMetrics above the maximum, setting it to the maximum")
	} else if batchSizePerfMetrics < 0 {
		batchSizePerfMetrics = 1
		logger.WithField("min size", 1).Warn("maxBatchSizePerfMetrics less then 0 no allowed, setting it to 1")
	}

	return batchSizePerfEntities, batchSizePerfMetrics, nil
}
