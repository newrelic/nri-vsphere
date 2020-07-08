package performance

import (
	"context"
	logrus "github.com/sirupsen/Logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/performance"
	"github.com/vmware/govmomi/simulator"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	"io/ioutil"
	"os"
	"testing"
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

func TestPerfCollector_NewCollector(t *testing.T) {

	content := []byte(`
host:
  level_1:
    - cpu.coreUtilization.average
    - metric.not.available
  level_2:
    - cpu.demand.average
vm:
  level_1:
    - metric.not.available1
    - cpu.demand.average
    - metric.not.available2
    #- commented
  level_3:
    - cpu.outoflevel
`)

	tmpfile, err := ioutil.TempFile("", "config")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name()) // clean up
	_, err = tmpfile.Write(content)
	require.NoError(t, err)

	_, err, c := startVcSim(t)

	pc, err := NewPerfCollector(c, logrus.New(), tmpfile.Name(), false, 2, "100", "50")
	assert.NoError(t, err)
	assert.Len(t, pc.MetricDefinition.Host, 2)
	assert.Len(t, pc.MetricDefinition.VM, 1)
	assert.Equal(t, 50, pc.batchSizePerfMetrics)
	assert.Equal(t, 100, pc.batchSizePerfEntities)

	ref := types.ManagedObjectReference{Type: "VirtualMachine", Value: "vm-87"}

	metrics := pc.Collect([]types.ManagedObjectReference{ref}, pc.MetricDefinition.VM)

	assert.Equal(t, 1, len(metrics), "we fetched events for 1 vm only")
	assert.Equal(t, 1, len(metrics[ref]), "we expect only one metric since only metrics with id 2 and 6 are defined for vms and only 2 is map in metricsAvaliableByID")
	assert.Greater(t, metrics[ref][0].Value, int64(0), "the value is not static, therefore we assume that a value grater then 0 is there")

}

func TestPerfMetricsEmptyPerfCollector(t *testing.T) {

	ctx, err, c := startVcSim(t)

	var vms []mo.VirtualMachine
	m := view.NewManager(c.Client)
	cv, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"VirtualMachine"}, true)
	assert.NoError(t, err)

	err = cv.Retrieve(ctx, []string{"VirtualMachine"}, []string{"name", "overallStatus"}, &vms)
	assert.NoError(t, err)
	var refSlice []types.ManagedObjectReference

	for _, vm := range vms {
		refSlice = append(refSlice, vm.Self)
	}

	p := PerfCollector{
		client:                 c,
		perfManager:            performance.NewManager(c.Client),
		logger:                 logrus.New(),
		MetricDefinition:       nil,
		metricsAvaliableByID:   nil,
		metricsAvaliableByName: nil,
		batchSizePerfEntities:  1,
		batchSizePerfMetrics:   1,
	}

	//no fail SEG/Fault expected
	metrics := p.Collect(refSlice, nil)
	assert.Equal(t, map[types.ManagedObjectReference][]PerfMetric{}, metrics)

	ms := []types.PerfMetricId{{CounterId: 1, Instance: ""}, {CounterId: 2, Instance: ""}, {CounterId: 3, Instance: ""}, {CounterId: 4, Instance: ""}}
	metrics = p.Collect(refSlice, ms)
	assert.Equal(t, map[types.ManagedObjectReference][]PerfMetric{}, metrics)
}

func startVcSim(t *testing.T) (context.Context, error, *govmomi.Client) {
	ctx := context.Background()

	//SettingUp Simulator
	model := simulator.VPX()
	model.Machine = 51
	err := model.Create()
	assert.NoError(t, err)

	s := model.Service.NewServer()

	c, _ := govmomi.NewClient(ctx, s.URL, true)
	return ctx, err, c
}

func TestPerfMetrics(t *testing.T) {

	ctx, err, c := startVcSim(t)

	var vms []mo.VirtualMachine
	m := view.NewManager(c.Client)
	cv, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"VirtualMachine"}, true)
	assert.NoError(t, err)

	err = cv.Retrieve(ctx, []string{"VirtualMachine"}, []string{"name", "overallStatus"}, &vms)
	assert.NoError(t, err)
	var refSlice []types.ManagedObjectReference

	for _, vm := range vms {
		refSlice = append(refSlice, vm.Self)
	}

	p := PerfCollector{
		client:                 c,
		perfManager:            performance.NewManager(c.Client),
		logger:                 logrus.New(),
		MetricDefinition:       nil,
		metricsAvaliableByID:   nil,
		metricsAvaliableByName: nil,
		batchSizePerfEntities:  1,
		batchSizePerfMetrics:   1,
	}

	//no fail SEG/Fault expected
	metrics := p.Collect(refSlice, nil)
	assert.Equal(t, map[types.ManagedObjectReference][]PerfMetric{}, metrics)

	//Please notice that only value for ID 2 and 6 is defined
	ms := []types.PerfMetricId{{CounterId: 1, Instance: ""}, {CounterId: 2, Instance: ""}, {CounterId: 5, Instance: ""}, {CounterId: 6, Instance: ""}}
	metrics = p.Collect(refSlice, ms)
	assert.Equal(t, map[types.ManagedObjectReference][]PerfMetric{}, metrics)

	p = PerfCollector{
		client:                 c,
		perfManager:            performance.NewManager(c.Client),
		logger:                 logrus.New(),
		MetricDefinition:       nil,
		metricsAvaliableByID:   map[int32]string{1: "test1", 2: "test2", 3: "test3"},
		metricsAvaliableByName: map[string]int32{"test1": 1, "test2": 2, "test3": 3},
		batchSizePerfEntities:  3,
		batchSizePerfMetrics:   3,
	}

	metrics = p.Collect(refSlice, ms)
	assert.Equal(t, len(refSlice), len(metrics), "we have 100 vm, all of them should be present in the map")
	assert.Equal(t, 1, len(metrics[refSlice[0]]), "we expect only one metric since only metrics with id 2 and 6 are defined for vms and only 2 is map in metricsAvaliableByID")

}

func TestSanitize(t *testing.T) {
	_, _, err := sanitizeArgs("1", "2")
	assert.NoError(t, err)
	_, _, err = sanitizeArgs("pg", "2")
	assert.Error(t, err)
	_, _, err = sanitizeArgs("-1", "2")
	assert.Error(t, err)
	_, _, err = sanitizeArgs("1", "0")
	assert.Error(t, err)
}
