package collect

import (
	"context"
	"github.com/newrelic/nri-vsphere/internal/load"
	"github.com/newrelic/nri-vsphere/internal/outputs"
	"github.com/stretchr/testify/assert"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/simulator"
	"testing"
	"time"
)

func TestEvents(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())

	config := load.NewConfig()
	config.Args = load.ArgumentList{}
	outputs.InfraIntegration(config)

	model := simulator.VPX()
	model.Machine = 5
	defer model.Remove()
	err := model.Create()
	if err != nil {
		config.Logrus.Fatal(err)
	}

	s := model.Service.NewServer()

	c, _ := govmomi.NewClient(ctx, s.URL, true)
	config.VMWareClient = c
	NewEventDispacher().Events(c.Client, config.Integration)

	time.Sleep(5 * time.Second)
	assert.Equal(t, 62, len(config.Integration.LocalEntity().Events), "We were expecting 10 events")

	vm, _ := find.NewFinder(c.Client).VirtualMachine(ctx, "DC0_H0_VM0")
	task, _ := vm.CreateSnapshot(ctx, "backup", "Backup", false, false)
	task.Wait(ctx)
	task, _ = vm.PowerOff(ctx)
	task.Wait(ctx)
	task, _ = vm.PowerOn(ctx)
	task.Wait(ctx)

	time.Sleep(5 * time.Second)
	assert.Equal(t, 66, len(config.Integration.LocalEntity().Events), "We were expecting 14 events")

	//TODO we should add test to check attributes

	//s.Close() We cannot close properly since the the connection is never closed and it is removed only if considered idle
	cancel()

}
