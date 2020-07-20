package process

import (
	"testing"
	"time"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-vsphere/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/vmware/govmomi/vim25/types"
)

func TestSnapshots(t *testing.T) {

	snapshot := types.ManagedObjectReference{
		Type:  "snapshot",
		Value: "1",
	}

	info, total, unique := processLayoutEx(getLayout())

	assert.NotNil(t, info)
	assert.Equal(t, int64(12*(1<<20)), info[snapshot].totalDisk)
	assert.Equal(t, int64(6*(1<<20)), info[snapshot].totalUniqueDisk)
	assert.Equal(t, int64(50*(1<<20)), info[snapshot].totalUniqueMemoryInDisk)
	assert.Equal(t, int64(100*(1<<20)), info[snapshot].totalMemoryInDisk)
	assert.Equal(t, "testPath|testPath2|", info[snapshot].datastorePathDisk)
	assert.Equal(t, "testPath3|", info[snapshot].datastorePathMemory)
	assert.Equal(t, int64(1000*(1<<20)), total)
	assert.Equal(t, int64(500*(1<<20)), unique)
}

func TestTRaverseSnapshots(t *testing.T) {

	info, _, _ := processLayoutEx(getLayout())
	config := config.Config{}
	i, _ := integration.New("test", "0.0.0")
	e := i.LocalEntity()
	traverseSnapshotList(e, &config, getTree(), "vmName", info)

	assert.Len(t, e.Metrics, 3)
	assert.Equal(t, "15", e.Metrics[0].Metrics["snapshotId"])
	assert.Equal(t, "16", e.Metrics[1].Metrics["snapshotId"])
	assert.Equal(t, "17", e.Metrics[2].Metrics["snapshotId"])
	assert.Equal(t, now.String(), e.Metrics[1].Metrics["creationTime"])
	assert.Equal(t, "vmName:snapshot1:snapshot2", e.Metrics[1].Metrics["snapshotTreeInfo"])
	assert.Equal(t, float64(6), e.Metrics[0].Metrics["totalUniqueDisk"])

}

var now = time.Now()

func getTree() types.VirtualMachineSnapshotTree {
	return types.VirtualMachineSnapshotTree{
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
	}
}

func getLayout() *types.VirtualMachineFileLayoutEx {
	return &types.VirtualMachineFileLayoutEx{
		File: append([]types.VirtualMachineFileLayoutExFileInfo{},
			types.VirtualMachineFileLayoutExFileInfo{
				Key:        0,
				Name:       "testPath",
				Type:       "snapshotData",
				Size:       2 * (1 << 20),
				UniqueSize: 1 * (1 << 20),
			},
			types.VirtualMachineFileLayoutExFileInfo{
				Key:        1,
				Name:       "testPath2",
				Type:       "snapshotData",
				Size:       10 * (1 << 20),
				UniqueSize: 5 * (1 << 20),
			},
			types.VirtualMachineFileLayoutExFileInfo{
				Key:        2,
				Name:       "testPath3",
				Type:       "snapshotMemory",
				Size:       100 * (1 << 20),
				UniqueSize: 50 * (1 << 20),
			},
			types.VirtualMachineFileLayoutExFileInfo{
				Key:        3,
				Name:       "testPath4",
				Type:       "suspendMemory",
				Size:       1000 * (1 << 20),
				UniqueSize: 500 * (1 << 20),
			},
			types.VirtualMachineFileLayoutExFileInfo{
				Key:        3,
				Name:       "testPath4",
				Type:       "other",
				Size:       10000 * (1 << 20),
				UniqueSize: 5000 * (1 << 20),
			}),

		Snapshot: append([]types.VirtualMachineFileLayoutExSnapshotLayout{},
			types.VirtualMachineFileLayoutExSnapshotLayout{
				DataKey:   0,
				MemoryKey: -1,
				Key: types.ManagedObjectReference{
					Type:  "snapshot",
					Value: "1",
				},
			},
			types.VirtualMachineFileLayoutExSnapshotLayout{
				DataKey:   1,
				MemoryKey: -1,
				Key: types.ManagedObjectReference{
					Type:  "snapshot",
					Value: "1",
				},
			},
			types.VirtualMachineFileLayoutExSnapshotLayout{
				DataKey:   -1,
				MemoryKey: 2,
				Key: types.ManagedObjectReference{
					Type:  "snapshot",
					Value: "1",
				},
			}),
	}
}
