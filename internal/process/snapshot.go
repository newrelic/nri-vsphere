// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package process

import (
	"math"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"

	"github.com/newrelic/infra-integrations-sdk/v3/data/metric"
	"github.com/newrelic/infra-integrations-sdk/v3/integration"
)

const invalidFile = -1

type snapshotProcessor struct {
	vmLayoutEx      *types.VirtualMachineFileLayoutEx
	currentSnapshot *types.ManagedObjectReference
	logger          *logrus.Logger

	// These structures are needed just to speedUpComputation.
	filesInfoByKey     map[int32]types.VirtualMachineFileLayoutExFileInfo
	snapshotsInfoByKey map[types.ManagedObjectReference]snapshotInfo
	allSnapshotsFiles  []int32
	allFiles           []int32

	results map[types.ManagedObjectReference]*infoSnapshot
}

// snapshotInfo holds raw data that needs to be computed.
type snapshotInfo struct {
	allFiles   []int32
	memoryFile int32
	dataFile   int32
}

func newSnapshotProcessor(logger *logrus.Logger, vm *mo.VirtualMachine) snapshotProcessor {
	snapshotsInfoByKey, allSnapshotFiles := buildSnapshotMapAndSlice(vm.LayoutEx.Snapshot)
	sp := snapshotProcessor{
		vmLayoutEx:         vm.LayoutEx,
		currentSnapshot:    vm.Snapshot.CurrentSnapshot,
		results:            map[types.ManagedObjectReference]*infoSnapshot{},
		logger:             logger,
		snapshotsInfoByKey: snapshotsInfoByKey,
		allSnapshotsFiles:  allSnapshotFiles,
		allFiles:           extractDiskLayoutFiles(vm.LayoutEx.Disk),
		filesInfoByKey:     buildFileKeyMap(vm.LayoutEx.File),
	}

	return sp
}

func buildFileKeyMap(filesInfo []types.VirtualMachineFileLayoutExFileInfo) map[int32]types.VirtualMachineFileLayoutExFileInfo {
	fileKeyMap := map[int32]types.VirtualMachineFileLayoutExFileInfo{}
	for _, file := range filesInfo {
		fileKeyMap[file.Key] = file
	}
	return fileKeyMap
}

func buildSnapshotMapAndSlice(layout []types.VirtualMachineFileLayoutExSnapshotLayout) (map[types.ManagedObjectReference]snapshotInfo, []int32) {
	var allSnapshotsFiles []int32
	snapshotFileKeyMap := map[types.ManagedObjectReference]snapshotInfo{}
	for _, snapLayout := range layout { // Extracting the list of files of the current snapshot of the loop.
		diskFiles := extractDiskLayoutFiles(snapLayout.Disk)
		// We create the list of all the files of all the snapshots.
		allSnapshotsFiles = append(allSnapshotsFiles, diskFiles...)
		snapshotFileKeyMap[snapLayout.Key] = snapshotInfo{
			allFiles:   diskFiles,
			memoryFile: snapLayout.MemoryKey,
			dataFile:   snapLayout.DataKey,
		}
	}
	return snapshotFileKeyMap, allSnapshotsFiles
}

// It takes care of going through VirtualMachineFileLayoutEx building up infoSnapshot.
func (sp snapshotProcessor) processSnapshotTree(parentSnapshot *types.ManagedObjectReference, snapshotTrees []types.VirtualMachineSnapshotTree) {
	for _, st := range snapshotTrees {
		st := st
		sp.snapshotSize(parentSnapshot, st.Snapshot)
		sp.processSnapshotTree(&st.Snapshot, st.ChildSnapshotList)
	}
}

// Function logic taken from the govc implementation
//
// SnapshotSize calculates the size of a given snapshot in bytes.
// List of snapshot files https://docs.vmware.com/en/VMware-vSphere/7.0/com.vmware.vsphere.hostclient.doc/GUID-38F4D574-ADE7-4B80-AEAB-7EC502A379F4.html.
func (sp snapshotProcessor) snapshotSize(parentSnapshot *types.ManagedObjectReference, snapshotRef types.ManagedObjectReference) {
	isCurrent := sp.isCurrent(snapshotRef)

	// Creating the computedSnapshotFiles needed to compute the snapshotSize.
	files := sp.buildFileStructure(parentSnapshot, snapshotRef, isCurrent)
	sp.computeDiskSizes(parentSnapshot, snapshotRef, files, isCurrent)

	if files.memory != 0 {
		if file, ok := sp.filesInfoByKey[files.memory]; ok {
			sp.results[snapshotRef].totalUniqueMemoryInDisk = file.UniqueSize
			sp.results[snapshotRef].totalMemoryInDisk = file.Size
			sp.results[snapshotRef].datastorePathMemory = file.Name
		}
	}
}

func (sp snapshotProcessor) isCurrent(snapshot types.ManagedObjectReference) bool {
	if sp.currentSnapshot != nil {
		if snapshot == *sp.currentSnapshot {
			return true
		}
	}
	return false
}

// This structure is needed just to return data in an easier way.
type computedSnapshotFiles struct {
	// memory points to the memory key if set.
	memory int32
	// dataAndDisk is the list of all files deltas and data if the snapshot has a parent
	// otherwise, it is just the snapshot data (to avoid considering the vm Disk).
	dataAndDisk []int32
	// deltaOfCurrent are the delta disks of the "Current" (from a Vsphere point of view) snapshot.
	deltaOfCurrent []int32
}

// Vmware consider that the snapshot files are the ones attached to the vm at the point of taking the snapshot.
// For example: the files of "snapshot 3" will be: delta of snapshot 2 (containing the delta of 1) + initial VM disk.
// Therefore, to compute the size of Snapshot 3 you need to do "Size of snapshot 3 files" - Size of 2 (parent) + Current delta (if snapshot if Current)
// In particular:
//
//	If the snapshot has a parent -> "diskFiles += snapshotFiles - parent snapshot files"
//	If the snapshot is the "current" (from a Vsphere point of view) one -> "diskFiles += allFiles - allSnapshotFiles"
//	To these files it is always added the "data file" that is quite small.
func (sp snapshotProcessor) computeDiskSizes(parentSnapshot *types.ManagedObjectReference, snapshotRef types.ManagedObjectReference, files computedSnapshotFiles, isCurrent bool) {
	var datastorePathDisk string
	var diskSize int64
	var uniqueDiskSize int64

	for _, fileKey := range files.dataAndDisk {
		if file, ok := sp.filesInfoByKey[fileKey]; ok {
			// if the parent is nil, we do not consider files belonging to the VM itself
			if parentSnapshot != nil || isMetadata(file) {
				diskSize += file.Size
				uniqueDiskSize += file.UniqueSize
				datastorePathDisk = datastorePathDisk + file.Name + "|"
			}
		}
	}

	if isCurrent {
		for _, diskFile := range files.deltaOfCurrent {
			if file, ok := sp.filesInfoByKey[diskFile]; ok {
				diskSize += file.Size
				uniqueDiskSize += file.UniqueSize
				datastorePathDisk = datastorePathDisk + file.Name + "|"
			}
		}
	}

	sp.results[snapshotRef] = &infoSnapshot{
		totalDisk:         diskSize,
		totalUniqueDisk:   uniqueDiskSize,
		datastorePathDisk: datastorePathDisk,
	}
}

func isMetadata(file types.VirtualMachineFileLayoutExFileInfo) bool {
	return file.Type != string(types.VirtualMachineFileLayoutExFileTypeDiskDescriptor) && file.Type != string(types.VirtualMachineFileLayoutExFileTypeDiskExtent)
}

func (sp snapshotProcessor) buildFileStructure(parentSnapshot *types.ManagedObjectReference, snapshotRef types.ManagedObjectReference, isCurrent bool) computedSnapshotFiles {
	var memoryFile int32
	var parentFiles []int32
	var dataAndDiskFiles []int32

	if data, ok := sp.snapshotsInfoByKey[snapshotRef]; ok {
		// Adding the .vmdk files of the snapshot we are interested into.
		dataAndDiskFiles = append(dataAndDiskFiles, data.allFiles...)
		// Adding the .vmsn file of the snapshot we are interested into.
		dataAndDiskFiles = append(dataAndDiskFiles, data.dataFile)
		memoryFile = data.memoryFile
	}

	if parentSnapshot != nil {
		if data, ok := sp.snapshotsInfoByKey[*parentSnapshot]; ok {
			parentFiles = append(parentFiles, data.allFiles...)
		}
	}

	// We do not consider any file belonging to a parent
	for _, parentFile := range parentFiles {
		removeKey(&dataAndDiskFiles, parentFile)
	}
	removeKey(&dataAndDiskFiles, invalidFile)

	var deltaOfCurrent []int32
	if isCurrent {
		// Extracting the list of all files related to a virtualMachine.
		// Then we remove all snapshots files that are already considered by parent snapshots
		// Remaining files are counted if the Snapshot is the "Current" one (from a Vsphere point of view).
		deltaOfCurrent = sp.allFiles
		for _, file := range sp.allSnapshotsFiles {
			removeKey(&deltaOfCurrent, file)
		}

		removeKey(&deltaOfCurrent, invalidFile)
	}

	return computedSnapshotFiles{
		memory:         memoryFile,
		dataAndDisk:    dataAndDiskFiles,
		deltaOfCurrent: deltaOfCurrent,
	}
}

// extractDiskLayoutFiles is a helper function used to extract file keys for
// all disk files attached to the virtual machine at the current point of running.
func extractDiskLayoutFiles(diskLayoutList []types.VirtualMachineFileLayoutExDiskLayout) []int32 {
	var result []int32

	for _, layoutExDisk := range diskLayoutList {
		for _, link := range layoutExDisk.Chain {
			result = append(result, link.FileKey...)
		}
	}

	return result
}

// removeKey is a helper function for removing a specific file key from a list
// of keys associated with disks attached to a virtual machine.
func removeKey(l *[]int32, key int32) {
	for i, k := range *l {
		if k == key {
			*l = append((*l)[:i], (*l)[i+1:]...)
			break
		}
	}
}

// It adds a new sample for each snapshot following the whole tree in a recursive way.
func (sp snapshotProcessor) createSnapshotSamples(e *integration.Entity, treeInfo string, snapshotTrees []types.VirtualMachineSnapshotTree) {
	for _, st := range snapshotTrees {
		ms := e.NewMetricSet("VSphere" + sampleTypeSnapshotVm + "Sample")
		treeInfo = treeInfo + ":" + st.Name

		sp.createMetricsCurrentSnapshot(treeInfo, st, ms)
		sp.createSnapshotSamples(e, treeInfo, st.ChildSnapshotList)
	}
}

func (sp snapshotProcessor) createMetricsCurrentSnapshot(treeInfo string, tree types.VirtualMachineSnapshotTree, ms *metric.Set) {
	checkError(sp.logger, ms.SetMetric("snapshotTreeInfo", treeInfo, metric.ATTRIBUTE))
	checkError(sp.logger, ms.SetMetric("name", tree.Name, metric.ATTRIBUTE))
	checkError(sp.logger, ms.SetMetric("creationTime", tree.CreateTime.String(), metric.ATTRIBUTE))
	checkError(sp.logger, ms.SetMetric("powerState", string(tree.State), metric.ATTRIBUTE))
	checkError(sp.logger, ms.SetMetric("snapshotId", strconv.FormatInt(int64(tree.Id), 10), metric.ATTRIBUTE))
	checkError(sp.logger, ms.SetMetric("quiesced", strconv.FormatBool(tree.Quiesced), metric.ATTRIBUTE))
	if tree.BackupManifest != "" {
		checkError(sp.logger, ms.SetMetric("backupManifest", tree.BackupManifest, metric.ATTRIBUTE))
	}
	if tree.Description != "" {
		checkError(sp.logger, ms.SetMetric("description", tree.Description, metric.ATTRIBUTE))
	}
	if tree.ReplaySupported != nil {
		checkError(sp.logger, ms.SetMetric("replaySupported", strconv.FormatBool(*tree.ReplaySupported), metric.ATTRIBUTE))
	}

	if i, ok := sp.results[tree.Snapshot]; ok {
		checkError(sp.logger, ms.SetMetric("totalMemoryInDisk", math.Ceil(float64(i.totalMemoryInDisk)/(1<<20)), metric.GAUGE))
		checkError(sp.logger, ms.SetMetric("totalUniqueMemoryInDisk", math.Ceil(float64(i.totalUniqueMemoryInDisk)/(1<<20)), metric.GAUGE))
		checkError(sp.logger, ms.SetMetric("totalDisk", math.Ceil(float64(i.totalDisk)/(1<<20)), metric.GAUGE))
		checkError(sp.logger, ms.SetMetric("totalUniqueDisk", math.Ceil(float64(i.totalUniqueDisk)/(1<<20)), metric.GAUGE))
		checkError(sp.logger, ms.SetMetric("datastorePathDisk", i.datastorePathDisk, metric.ATTRIBUTE))
		checkError(sp.logger, ms.SetMetric("datastorePathMemory", i.datastorePathMemory, metric.ATTRIBUTE))
	}
}

// This struct is used to save dataAndDisk before creating the metrics. Otherwise, we would need to go thorugh the dataAndDisk structure many times.
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
