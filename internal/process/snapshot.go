// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package process

import (
	"strconv"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/vmware/govmomi/vim25/types"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/nri-vsphere/internal/load"
)

// It takes care of going through VirtualMachineFileLayoutEx building up infoSnapshot
func processLayoutEx(ex *types.VirtualMachineFileLayoutEx) (info map[types.ManagedObjectReference]*infoSnapshot, suspendMemory int64, suspendMemoryUnique int64) {

	info = map[types.ManagedObjectReference]*infoSnapshot{}
	if ex != nil {
		for _, exFile := range ex.File {
			if exFile.Type == "snapshotData" {
				findSnapshotAndUpdateInfoData(ex.Snapshot, exFile, info)
			} else if exFile.Type == "snapshotMemory" {
				findSnapshotAndUpdateInfoMemory(ex.Snapshot, exFile, info)
			} else if exFile.Type == "suspendMemory" {
				suspendMemory += exFile.Size
				suspendMemoryUnique += exFile.UniqueSize
			}
		}
	}
	return info, suspendMemory, suspendMemoryUnique
}

// It finds for the given datakey the snapshotRef and saves the data
func findSnapshotAndUpdateInfoData(snapshostFiles []types.VirtualMachineFileLayoutExSnapshotLayout, exFile types.VirtualMachineFileLayoutExFileInfo, info map[types.ManagedObjectReference]*infoSnapshot) {
	for _, exs := range snapshostFiles {
		if exs.DataKey == exFile.Key {
			i, ok := info[exs.Key]
			if !ok {
				i = &infoSnapshot{}
				info[exs.Key] = i
			}
			i.totalDisk += exFile.Size
			i.totalUniqueDisk += exFile.UniqueSize
			i.datastorePathDisk = i.datastorePathDisk + exFile.Name + "|"
		}
	}
}

// It finds for the given memorykey the snapshotRef and saves the data
func findSnapshotAndUpdateInfoMemory(snapshostFiles []types.VirtualMachineFileLayoutExSnapshotLayout, exFile types.VirtualMachineFileLayoutExFileInfo, info map[types.ManagedObjectReference]*infoSnapshot) {
	for _, exs := range snapshostFiles {
		if exs.MemoryKey == exFile.Key && exs.MemoryKey != -1 {
			i, ok := info[exs.Key]
			if !ok {
				i = &infoSnapshot{}
				info[exs.Key] = i
			}
			i.totalMemoryInDisk += exFile.Size
			i.totalUniqueMemoryInDisk += exFile.UniqueSize
			i.datastorePathMemory = i.datastorePathMemory + exFile.Name + "|"
		}
	}
}

// It adds a new sample for each snapshot following the whole tree in a recursive way
func traverseSnapshotList(e *integration.Entity, config *load.Config, tree types.VirtualMachineSnapshotTree, treeInfo string, info map[types.ManagedObjectReference]*infoSnapshot) {

	ms := e.NewMetricSet("VSphere" + sampleTypeSnapshotVm + "Sample")
	treeInfo = treeInfo + ":" + tree.Name
	createMetricsCurrentSnapshot(treeInfo, tree, config, ms, info)

	//A recursive function is needed since the actual size of the tree in unknown
	for _, s := range tree.ChildSnapshotList {
		traverseSnapshotList(e, config, s, treeInfo, info)
	}

}

func createMetricsCurrentSnapshot(treeInfo string, tree types.VirtualMachineSnapshotTree, config *load.Config, ms *metric.Set, info map[types.ManagedObjectReference]*infoSnapshot) {
	checkError(config, ms.SetMetric("snapshotTreeInfo", treeInfo, metric.ATTRIBUTE))
	checkError(config, ms.SetMetric("name", tree.Name, metric.ATTRIBUTE))
	checkError(config, ms.SetMetric("creationTime", tree.CreateTime.String(), metric.ATTRIBUTE))
	checkError(config, ms.SetMetric("powerState", string(tree.State), metric.ATTRIBUTE))
	checkError(config, ms.SetMetric("snapshotId", strconv.FormatInt(int64(tree.Id), 10), metric.ATTRIBUTE))
	checkError(config, ms.SetMetric("quiesced", strconv.FormatBool(tree.Quiesced), metric.ATTRIBUTE))
	if tree.BackupManifest != "" {
		checkError(config, ms.SetMetric("backupManifest", tree.BackupManifest, metric.ATTRIBUTE))
	}
	if tree.Description != "" {
		checkError(config, ms.SetMetric("description", tree.Description, metric.ATTRIBUTE))
	}
	if tree.ReplaySupported != nil {
		checkError(config, ms.SetMetric("replaySupported", strconv.FormatBool(*tree.ReplaySupported), metric.ATTRIBUTE))
	}

	if i, ok := info[tree.Snapshot]; ok {
		checkError(config, ms.SetMetric("totalMemoryInDisk", i.totalMemoryInDisk/(1<<20), metric.GAUGE))
		checkError(config, ms.SetMetric("totalUniqueMemoryInDisk", i.totalUniqueMemoryInDisk/(1<<20), metric.GAUGE))
		checkError(config, ms.SetMetric("totalDisk", i.totalDisk/(1<<20), metric.GAUGE))
		checkError(config, ms.SetMetric("totalUniqueDisk", i.totalUniqueDisk/(1<<20), metric.GAUGE))
		checkError(config, ms.SetMetric("datastorePathDisk", i.datastorePathDisk, metric.ATTRIBUTE))
		checkError(config, ms.SetMetric("datastorePathMemory", i.datastorePathMemory, metric.ATTRIBUTE))
	}
}

// This struct is used to save data before creating the metrics. Otherwise we would need to go thorugh the data structure many times
type infoSnapshot struct {
	totalMemoryInDisk int64
	totalDisk         int64
	//Size of the file in bytes corresponding to the file blocks that were allocated uniquely.
	//In other words, if the underlying storage supports sharing of file blocks across disk files,
	//the property corresponds to the size of the file blocks that were allocated only in context of this file,
	//i.e. it does not include shared blocks that were allocated in other files. This property will be unset if the
	//underlying implementation is unable to compute this information.
	//One example of this is when the file resides on a NAS datastore whose underlying storage doesn't support this
	//metric. In some cases the field might be set but the value could be over-estimated due to the inability of the NAS
	//based storage to provide an accurate value.
	totalUniqueMemoryInDisk int64
	totalUniqueDisk         int64
	datastorePathDisk       string
	datastorePathMemory     string
}
