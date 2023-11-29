package process

import (
	"context"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-vsphere/internal/client"
	"github.com/newrelic/nri-vsphere/internal/collect"
	"github.com/newrelic/nri-vsphere/internal/config"
	"github.com/newrelic/nri-vsphere/internal/model"
	"github.com/newrelic/nri-vsphere/internal/process/testdata"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/vmware/govmomi/simulator"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	"testing"
)

func Test_createVirtualMachineSamples_HasIpAddresses(t *testing.T) {
	simulator.Run(func(ctx context.Context, vc *vim25.Client) error {
		vmClient, err := client.New(vc.URL().String(), "user", "pass", false)
		assert.NoError(t, err)
		vm := view.NewManager(vc)
		assert.NotNil(t, vm)
		// given
		cfg := &config.Config{VMWareClient: vmClient, ViewManager: vm, Logrus: logrus.StandardLogger()}
		cfg.Integration, _ = integration.New("test", "dev")
		cfg.Datacenters = append(cfg.Datacenters, getDatacenter(ctx, vm))

		// when
		// we need the host to create the vms
		collect.Hosts(cfg)
		collect.VirtualMachines(cfg)

		createVirtualMachineSamples(cfg)
		// then
		assert.True(t, len(cfg.Datacenters[0].VirtualMachines) > 0)
		for _, e := range cfg.Integration.Entities {
			for _, ms := range e.Metrics {
				// we just chek the presence because it might have no 'extra' ip addresses but we still add the attribute
				assert.Contains(t, ms.Metrics, "ipAddresses")
			}
		}
		return nil
	})
}

func getDatacenter(ctx context.Context, vm *view.Manager) *model.Datacenter {
	cv, err := vm.CreateContainerView(ctx, vm.Client().ServiceContent.RootFolder, []string{"Datacenter"}, false)
	if err != nil {
		logrus.Fatal("failed to get container view")
	}
	var datacenters []mo.Datacenter
	_ = cv.Retrieve(ctx, []string{"Datacenter"}, []string{"name"}, &datacenters)
	return model.NewDatacenter(&datacenters[0])
}

const hostname = "test"
const fullHostname = "test.this.com"
const domain = "this.com"

func TestComputeFullHostname(t *testing.T) {
	var vm = &mo.VirtualMachine{}
	assert.Equal(t, "", computeFullHostname(vm))

	vm = &mo.VirtualMachine{
		Guest: &types.GuestInfo{
			IpStack: []types.GuestStackInfo{
				{
					DnsConfig: &types.NetDnsConfigInfo{
						HostName:   hostname + "different",
						DomainName: domain,
					},
				},
				{
					DnsConfig: &types.NetDnsConfigInfo{
						HostName: hostname,
					},
				},
				{
					DnsConfig: &types.NetDnsConfigInfo{
						DomainName: domain,
					},
				},
				{ // This is the only entry that the implementation should consider
					DnsConfig: &types.NetDnsConfigInfo{
						HostName:   hostname,
						DomainName: domain,
					},
				},
			},
		},
		Summary: types.VirtualMachineSummary{
			Guest: &types.VirtualMachineGuestSummary{
				HostName: hostname,
			},
		},
	}
	assert.Equal(t, fullHostname, computeFullHostname(vm))

	// No matter if in hostname there is the fqdn, we do not place the suffix twice
	vm.Guest.IpStack[3].DnsConfig.HostName = fullHostname
	vm.Summary.Guest.HostName = fullHostname
	assert.Equal(t, fullHostname, computeFullHostname(vm))

	// No matter if in hostname there is the fqdn, we do not place the suffix twice
	vm.Guest.IpStack[3].DnsConfig.HostName = hostname + "."
	vm.Summary.Guest.HostName = hostname + "."
	assert.Equal(t, fullHostname, computeFullHostname(vm))

	// if the hostname is different in the summary we avoid computing the fqdn
	vm.Guest.IpStack[3].DnsConfig.HostName = hostname + "different"
	assert.Equal(t, "", computeFullHostname(vm))

	// Testing it with mock from real data
	realVm := testdata.GetVMFromStaticData(t)
	assert.Equal(t, "vm-3.test.com", computeFullHostname(&realVm))
}
