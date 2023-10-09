package process

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-vsphere/internal/process/testdata"
)

func TestSnapshotsRealData(t *testing.T) {
	t.Parallel()

	vm := testdata.GetVMFromStaticData(t)
	// This dataAndDisk describes the following scenario
	// Taken directly marshaling the dataAndDisk into JSON from a real VCenter
	// [276.0KB]  Snap1 No Mem
	//   [63.5MB]  Snap2
	//     [1.1GB]  Snap3 No Mem
	//       [2.1GB]  Snap4
	//         [1.7GB]  Snap5
	//           [60.3MB]  Snap6 No Mem
	//             [59.5MB]  Snap7

	// Compared to the UI all values matches, but the one of the first snapshot
	// Since we decided not to count the disk of the machine in the first snapshot

	sp := newSnapshotProcessor(nil, &vm)

	sp.processSnapshotTree(nil, vm.Snapshot.RootSnapshotList)
	testRawResults(t, sp)

	// Verify against library
	assert.Equal(t, int64(object.SnapshotSize(getSnapshotReference(1), nil, vm.LayoutEx, false)), sp.results[getSnapshotReference(1)].totalDisk)
	parent := getSnapshotReference(2)
	assert.Equal(t, int64(object.SnapshotSize(getSnapshotReference(3), &parent, vm.LayoutEx, false)), sp.results[getSnapshotReference(3)].totalDisk)
	parent = getSnapshotReference(6)
	assert.Equal(t, int64(object.SnapshotSize(getSnapshotReference(7), &parent, vm.LayoutEx, true)), sp.results[getSnapshotReference(7)].totalDisk)

	i, _ := integration.New("test", "0.0.0")
	e := i.LocalEntity()
	sp.createSnapshotSamples(e, "vmName", vm.Snapshot.RootSnapshotList)
	testMetrics(t, e)
}

func testRawResults(t *testing.T, sp snapshotProcessor) {
	t.Helper()

	require.NotNil(t, sp.results)
	assert.Len(t, sp.results, 7)

	assert.Equal(t, int64(282644), sp.results[getSnapshotReference(1)].totalDisk)
	assert.Equal(t, int64(66575311), sp.results[getSnapshotReference(2)].totalDisk)
	assert.Equal(t, int64(1195659291), sp.results[getSnapshotReference(3)].totalDisk)
	assert.Equal(t, int64(2276978168), sp.results[getSnapshotReference(4)].totalDisk)
	assert.Equal(t, int64(1790451907), sp.results[getSnapshotReference(5)].totalDisk)
	assert.Equal(t, int64(63197211), sp.results[getSnapshotReference(6)].totalDisk)
	assert.Equal(t, int64(79175863), sp.results[getSnapshotReference(7)].totalDisk)

	assert.Equal(t, int64(0), sp.results[getSnapshotReference(1)].totalMemoryInDisk)
	assert.Equal(t, int64(2147483648), sp.results[getSnapshotReference(2)].totalMemoryInDisk)
	assert.Equal(t, int64(0), sp.results[getSnapshotReference(3)].totalMemoryInDisk)
	assert.Equal(t, int64(2147483648), sp.results[getSnapshotReference(4)].totalMemoryInDisk)
	assert.Equal(t, int64(2202009600), sp.results[getSnapshotReference(5)].totalUniqueMemoryInDisk)
	assert.Equal(t, int64(0), sp.results[getSnapshotReference(6)].totalUniqueMemoryInDisk)
	assert.Equal(t, int64(2202009600), sp.results[getSnapshotReference(7)].totalUniqueMemoryInDisk)

	assert.Contains(t, sp.results[getSnapshotReference(1)].datastorePathDisk, "test-snap-Snapshot3.vmsn")
	assert.Contains(t, sp.results[getSnapshotReference(2)].datastorePathDisk, "test-snap-Snapshot4.vmsn")
	assert.Contains(t, sp.results[getSnapshotReference(2)].datastorePathDisk, "test-snap-000001.vmdk")
	assert.Contains(t, sp.results[getSnapshotReference(3)].datastorePathMemory, "")
	assert.Contains(t, sp.results[getSnapshotReference(4)].datastorePathMemory, "test-snap-Snapshot6.vmem")
	assert.Contains(t, sp.results[getSnapshotReference(5)].datastorePathDisk, "test-snap-Snapshot7.vmsn")
	assert.Contains(t, sp.results[getSnapshotReference(5)].datastorePathDisk, "test-snap-000004.vmdk")
	assert.Contains(t, sp.results[getSnapshotReference(6)].datastorePathDisk, "test-snap-000005.vmdk")
}

func testMetrics(t *testing.T, e *integration.Entity) {
	t.Helper()

	assert.Len(t, e.Metrics, 7)
	assert.Equal(t, "3", e.Metrics[0].Metrics["snapshotId"])
	assert.Equal(t, "4", e.Metrics[1].Metrics["snapshotId"])
	assert.Equal(t, "5", e.Metrics[2].Metrics["snapshotId"])
	assert.Equal(t, "2023-10-05 11:17:05.637214 +0000 UTC", e.Metrics[1].Metrics["creationTime"])
	assert.Equal(t, "vmName:Snap1:Snap2", e.Metrics[1].Metrics["snapshotTreeInfo"])
	assert.Equal(t, float64(1), e.Metrics[0].Metrics["totalUniqueDisk"])
	assert.Equal(t, float64(1141), e.Metrics[2].Metrics["totalUniqueDisk"])
	assert.Equal(t, float64(2172), e.Metrics[3].Metrics["totalUniqueDisk"])
	assert.Equal(t, float64(61), e.Metrics[5].Metrics["totalUniqueDisk"])
	assert.Contains(t, e.Metrics[0].Metrics["datastorePathDisk"], "test-snap-Snapshot3.vmsn")
	assert.Equal(t, "false", e.Metrics[2].Metrics["replaySupported"])
	assert.Equal(t, float64(2172), e.Metrics[3].Metrics["totalDisk"])
	assert.Equal(t, float64(0), e.Metrics[5].Metrics["totalMemoryInDisk"])
	assert.Equal(t, float64(2048), e.Metrics[6].Metrics["totalMemoryInDisk"])
}

func getSnapshotReference(number int) types.ManagedObjectReference {
	return types.ManagedObjectReference{
		Type:  "VirtualMachineSnapshot",
		Value: fmt.Sprintf("snapshot-%d", number),
	}
}

// This test even if it is similar allows to control the vm structure.
func TestSnapshots(t *testing.T) {
	t.Parallel()

	snapshot := types.ManagedObjectReference{
		Type:  "snapshot",
		Value: "1",
	}

	vm := mo.VirtualMachine{
		LayoutEx: getLayout(),
		Snapshot: &types.VirtualMachineSnapshotInfo{
			CurrentSnapshot: &snapshot,
		},
	}
	sp := newSnapshotProcessor(nil, &vm)
	sp.processSnapshotTree(nil, getTree())

	require.NotNil(t, sp.results)
	assert.Len(t, sp.results, 3)
	s := sp.results[snapshot]
	assert.Equal(t, int64(4010), s.totalDisk)
	assert.Equal(t, int64(2005), s.totalUniqueDisk)
	assert.Equal(t, int64(50), s.totalUniqueMemoryInDisk)
	assert.Equal(t, int64(100), s.totalMemoryInDisk)
	assert.Contains(t, s.datastorePathDisk, "test-000001.vmdk")
	assert.Contains(t, s.datastorePathDisk, "testPath2")
	assert.Equal(t, "testPath3", s.datastorePathMemory)

	// We can check the totalSize against the govmomi library
	size := object.SnapshotSize(snapshot, nil, getLayout(), true)
	assert.Equal(t, s.totalDisk, int64(size))
}

var now = time.Now()

func getTree() []types.VirtualMachineSnapshotTree {
	return []types.VirtualMachineSnapshotTree{
		{
			Snapshot: types.ManagedObjectReference{
				Type:  "snapshot",
				Value: "1",
			},
			Name:           "snapshot1",
			Description:    "Description",
			Id:             15,
			CreateTime:     now,
			State:          "ready",
			Quiesced:       false,
			BackupManifest: "BackupManifest",
			ChildSnapshotList: append([]types.VirtualMachineSnapshotTree{},
				types.VirtualMachineSnapshotTree{
					Snapshot: types.ManagedObjectReference{
						Type:  "snapshot",
						Value: "2",
					},
					Name:           "snapshot2",
					Description:    "Description",
					Id:             16,
					CreateTime:     now,
					State:          "state",
					BackupManifest: "BackupManifest",
					ChildSnapshotList: append([]types.VirtualMachineSnapshotTree{},
						types.VirtualMachineSnapshotTree{
							Snapshot: types.ManagedObjectReference{
								Type:  "snapshot",
								Value: "3",
							},
							Name:           "snapshot3",
							Description:    "Description",
							Id:             17,
							CreateTime:     now,
							State:          "state",
							BackupManifest: "BackupManifest",
						}),
				}),
		},
	}
}

func getLayout() *types.VirtualMachineFileLayoutEx {
	return &types.VirtualMachineFileLayoutEx{
		Disk: []types.VirtualMachineFileLayoutExDiskLayout{
			{
				DynamicData: types.DynamicData{},
				Key:         200,
				Chain: []types.VirtualMachineFileLayoutExDiskUnit{
					{
						FileKey: []int32{
							5,
						},
					},
					{
						FileKey: []int32{
							6,
						},
					},
				},
			},
		},
		File: []types.VirtualMachineFileLayoutExFileInfo{
			{
				Key:        0,
				Name:       "testPath",
				Type:       "snapshotInfo",
				Size:       2,
				UniqueSize: 1,
			},
			{
				Key:        1,
				Name:       "testPath2",
				Type:       "snapshotInfo",
				Size:       10,
				UniqueSize: 5,
			},
			{
				Key:        2,
				Name:       "testPath3",
				Type:       "snapshotMemory",
				Size:       100,
				UniqueSize: 50,
			},
			{
				Key:        3,
				Name:       "testPath4",
				Type:       "suspendMemory",
				Size:       1000,
				UniqueSize: 500,
			},
			{
				Key:        4,
				Name:       "testPath4",
				Type:       "other",
				Size:       10000,
				UniqueSize: 5000,
			},
			{
				Key:        5,
				Name:       "test.vmdk",
				Type:       "diskDescriptor",
				Size:       10000,
				UniqueSize: 5000,
			},
			{
				Key:        6,
				Name:       "test-000001.vmdk",
				Type:       "diskDescriptor",
				Size:       4000,
				UniqueSize: 2000,
			},
		},

		Snapshot: []types.VirtualMachineFileLayoutExSnapshotLayout{
			{
				DataKey:   1,
				MemoryKey: 2,
				Disk: []types.VirtualMachineFileLayoutExDiskLayout{
					{
						DynamicData: types.DynamicData{},
						Key:         200,
						Chain: []types.VirtualMachineFileLayoutExDiskUnit{
							{
								FileKey: []int32{
									5,
								},
							},
						},
					},
				},
				Key: types.ManagedObjectReference{
					Type:  "snapshot",
					Value: "1",
				},
			},
			{
				DataKey:   4,
				MemoryKey: 2,
				Key: types.ManagedObjectReference{
					Type:  "snapshot",
					Value: "2",
				},
			},
			{
				DataKey:   -1,
				MemoryKey: -1,
				Key: types.ManagedObjectReference{
					Type:  "snapshot",
					Value: "3",
				},
			},
		},
	}
}
