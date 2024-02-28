package testdata

import (
	"encoding/json"
	"github.com/vmware/govmomi/vim25/mo"
	"testing"
)

func GetVMFromStaticData(t *testing.T) (test mo.VirtualMachine) {
	t.Helper()

	err := json.Unmarshal([]byte(vmData), &test)
	if err != nil {
		t.Fatal()
	}
	return test
}

const vmData = `
{
  "Self": {
    "Type": "VirtualMachine",
    "Value": "vm-3"
  },
  "Value": null,
  "AvailableField": null,
  "Parent": null,
  "CustomValue": null,
  "OverallStatus": "green",
  "ConfigStatus": "",
  "ConfigIssue": null,
  "EffectiveRole": null,
  "Permission": null,
  "Name": "test-snap",
  "DisabledMethod": null,
  "RecentTask": null,
  "DeclaredAlarmState": null,
  "TriggeredAlarmState": null,
  "AlarmActionsEnabled": null,
  "Tag": null,
  "Capability": {
    "SnapshotOperationsSupported": false,
    "MultipleSnapshotsSupported": false,
    "SnapshotConfigSupported": false,
    "PoweredOffSnapshotsSupported": false,
    "MemorySnapshotsSupported": false,
    "RevertToSnapshotSupported": false,
    "QuiescedSnapshotsSupported": false,
    "DisableSnapshotsSupported": false,
    "LockSnapshotsSupported": false,
    "ConsolePreferencesSupported": false,
    "CpuFeatureMaskSupported": false,
    "S1AcpiManagementSupported": false,
    "SettingScreenResolutionSupported": false,
    "ToolsAutoUpdateSupported": false,
    "VmNpivWwnSupported": false,
    "NpivWwnOnNonRdmVmSupported": false,
    "VmNpivWwnDisableSupported": null,
    "VmNpivWwnUpdateSupported": null,
    "SwapPlacementSupported": false,
    "ToolsSyncTimeSupported": false,
    "VirtualMmuUsageSupported": false,
    "DiskSharesSupported": false,
    "BootOptionsSupported": false,
    "BootRetryOptionsSupported": null,
    "SettingVideoRamSizeSupported": false,
    "SettingDisplayTopologySupported": null,
    "RecordReplaySupported": null,
    "ChangeTrackingSupported": null,
    "MultipleCoresPerSocketSupported": null,
    "HostBasedReplicationSupported": null,
    "GuestAutoLockSupported": null,
    "MemoryReservationLockSupported": null,
    "FeatureRequirementSupported": null,
    "PoweredOnMonitorTypeChangeSupported": null,
    "SeSparseDiskSupported": null,
    "NestedHVSupported": null,
    "VPMCSupported": null,
    "SecureBootSupported": null,
    "PerVmEvcSupported": null,
    "VirtualMmuUsageIgnored": null,
    "VirtualExecUsageIgnored": null,
    "DiskOnlySnapshotOnSuspendedVMSupported": null,
    "SuspendToMemorySupported": null,
    "ToolsSyncTimeAllowSupported": null,
    "SevSupported": null,
    "PmemFailoverSupported": null,
    "RequireSgxAttestationSupported": null,
    "ChangeModeDisksSupported": null
  },
  "Config": {
    "ChangeVersion": "2023-10-05T10:09:31.249952Z",
    "Modified": "1970-01-01T00:00:00Z",
    "Name": "test-snap",
    "GuestFullName": "CentOS 8 (64-bit)",
    "Version": "vmx-14",
    "Uuid": "420fe454-848e-4d0b-a08f-c2423606b46c",
    "CreateDate": "2023-10-05T08:33:27.781429Z",
    "InstanceUuid": "500f7e58-8c6b-2aae-3ffb-ea87f750a55d",
    "NpivNodeWorldWideName": null,
    "NpivPortWorldWideName": null,
    "NpivWorldWideNameType": "",
    "NpivDesiredNodeWwns": 0,
    "NpivDesiredPortWwns": 0,
    "NpivTemporaryDisabled": true,
    "NpivOnNonRdmDisks": null,
    "LocationId": "564debde-dd20-acc1-24b0-d17d1d6a37b0",
    "Template": false,
    "GuestId": "centos8_64Guest",
    "AlternateGuestName": "",
    "Annotation": "",
    "Files": {
      "VmPathName": "[WorkloadDatastore] d7741e65-148e-b5ff-2458-068cdd1d8254/test-snap.vmx",
      "SnapshotDirectory": "[WorkloadDatastore] d7741e65-148e-b5ff-2458-068cdd1d8254/",
      "SuspendDirectory": "[WorkloadDatastore] d7741e65-148e-b5ff-2458-068cdd1d8254/",
      "LogDirectory": "[WorkloadDatastore] d7741e65-148e-b5ff-2458-068cdd1d8254/",
      "FtMetadataDirectory": ""
    },
    "Tools": {
      "ToolsVersion": 0,
      "ToolsInstallType": "",
      "AfterPowerOn": true,
      "AfterResume": true,
      "BeforeGuestStandby": true,
      "BeforeGuestShutdown": true,
      "BeforeGuestReboot": null,
      "ToolsUpgradePolicy": "manual",
      "PendingCustomization": "",
      "CustomizationKeyId": null,
      "SyncTimeWithHostAllowed": true,
      "SyncTimeWithHost": false,
      "LastInstallInfo": {
        "Counter": 0,
        "Fault": null
      }
    },
    "Flags": {
      "DisableAcceleration": null,
      "EnableLogging": true,
      "UseToe": false,
      "RunWithDebugInfo": false,
      "MonitorType": "release",
      "HtSharing": "any",
      "SnapshotDisabled": false,
      "SnapshotLocked": false,
      "DiskUuidEnabled": false,
      "VirtualMmuUsage": "",
      "VirtualExecUsage": "",
      "SnapshotPowerOffBehavior": "powerOff",
      "RecordReplayEnabled": false,
      "FaultToleranceType": "unset",
      "CbrcCacheEnabled": false,
      "VvtdEnabled": false,
      "VbsEnabled": false
    },
    "ConsolePreferences": null,
    "DefaultPowerOps": {
      "PowerOffType": "soft",
      "SuspendType": "hard",
      "ResetType": "soft",
      "DefaultPowerOffType": "soft",
      "DefaultSuspendType": "hard",
      "DefaultResetType": "soft",
      "StandbyAction": "checkpoint"
    },
    "RebootPowerOff": false,
    "VcpuConfig": null,
    "CpuAllocation": {
      "Reservation": 0,
      "ExpandableReservation": false,
      "Limit": -1,
      "Shares": {
        "Shares": 1000,
        "Level": "normal"
      },
      "OverheadLimit": null
    },
    "MemoryAllocation": {
      "Reservation": 0,
      "ExpandableReservation": false,
      "Limit": -1,
      "Shares": {
        "Shares": 20480,
        "Level": "normal"
      },
      "OverheadLimit": 73
    },
    "LatencySensitivity": {
      "Level": "normal",
      "Sensitivity": 0
    },
    "MemoryHotAddEnabled": false,
    "CpuHotAddEnabled": false,
    "CpuHotRemoveEnabled": false,
    "HotPlugMemoryLimit": 2048,
    "HotPlugMemoryIncrementSize": 0,
    "CpuAffinity": null,
    "MemoryAffinity": null,
    "NetworkShaper": null,
    "CpuFeatureMask": null,
    "DatastoreUrl": [
      {
        "Name": "WorkloadDatastore",
        "Url": "/vmfs/volumes/vsan:daa0538afb0a4be7-953fb0344f2d1123"
      }
    ],
    "SwapPlacement": "inherit",
    "BootOptions": {
      "BootDelay": 0,
      "EnterBIOSSetup": false,
      "EfiSecureBootEnabled": true,
      "BootRetryEnabled": false,
      "BootRetryDelay": 10000,
      "BootOrder": null,
      "NetworkBootProtocol": "ipv4"
    },
    "FtInfo": null,
    "RepConfig": null,
    "VAppConfig": null,
    "VAssertsEnabled": false,
    "ChangeTrackingEnabled": false,
    "Firmware": "efi",
    "MaxMksConnections": -1,
    "GuestAutoLockEnabled": true,
    "ManagedBy": null,
    "MemoryReservationLockedToMax": false,
    "InitialOverhead": {
      "InitialMemoryReservation": 70643712,
      "InitialSwapReservation": 920489984
    },
    "NestedHVEnabled": false,
    "VPMCEnabled": false,
    "ScheduledHardwareUpgradeInfo": {
      "UpgradePolicy": "never",
      "VersionKey": "",
      "ScheduledHardwareUpgradeStatus": "none",
      "Fault": null
    },
    "ForkConfigInfo": null,
    "VFlashCacheReservation": 0,
    "VmxConfigChecksum": "R2JvQmhFT2dsZW1FL1BwN05KK0p6YklhUGg4PQ==",
    "MessageBusTunnelEnabled": false,
    "VmStorageObjectId": "d7741e65-148e-b5ff-2458-068cdd1d8254",
    "SwapStorageObjectId": "dd741e65-d34b-edea-9e04-068cdd1d8254",
    "KeyId": null,
    "GuestIntegrityInfo": {
      "Enabled": false
    },
    "MigrateEncryption": "opportunistic",
    "SgxInfo": {
      "EpcSize": 0,
      "FlcMode": "unlocked",
      "LePubKeyHash": "",
      "RequireAttestation": false
    },
    "ContentLibItemInfo": null,
    "FtEncryptionMode": "ftEncryptionOpportunistic",
    "GuestMonitoringModeInfo": {
      "GmmFile": "",
      "GmmAppliance": ""
    },
    "SevEnabled": false,
    "NumaInfo": {
      "CoresPerNumaNode": 1,
      "AutoCoresPerNumaNode": true,
      "VnumaOnCpuHotaddExposed": false
    },
    "PmemFailoverEnabled": false,
    "VmxStatsCollectionEnabled": true,
    "VmOpNotificationToAppEnabled": false,
    "VmOpNotificationTimeout": -1,
    "DeviceSwap": {
      "LsiToPvscsi": {
        "Enabled": true,
        "Applicable": false,
        "Status": "none"
      }
    },
    "Pmem": null,
    "DeviceGroups": {
      "DeviceGroup": null
    }
  },
  "Layout": null,
  "LayoutEx": {
    "File": [
      {
        "Key": 0,
        "Name": "[WorkloadDatastore] d7741e65-148e-b5ff-2458-068cdd1d8254/test-snap.vmx",
        "Type": "config",
        "Size": 0,
        "UniqueSize": 0,
        "BackingObjectId": "d7741e65-148e-b5ff-2458-068cdd1d8254",
        "Accessible": true
      },
      {
        "Key": 3,
        "Name": "[WorkloadDatastore] d7741e65-148e-b5ff-2458-068cdd1d8254/test-snap.nvram",
        "Type": "nvram",
        "Size": 270840,
        "UniqueSize": 270840,
        "BackingObjectId": "",
        "Accessible": true
      },
      {
        "Key": 1,
        "Name": "[WorkloadDatastore] d7741e65-148e-b5ff-2458-068cdd1d8254/test-snap.vmsd",
        "Type": "snapshotList",
        "Size": 0,
        "UniqueSize": 0,
        "BackingObjectId": "",
        "Accessible": true
      },
      {
        "Key": 2,
        "Name": "[WorkloadDatastore] d7741e65-148e-b5ff-2458-068cdd1d8254/test-snap.vmdk",
        "Type": "diskDescriptor",
        "Size": 2218786816,
        "UniqueSize": 2218786816,
        "BackingObjectId": "d8741e65-6877-4539-7350-068cdd1d8254",
        "Accessible": true
      },
      {
        "Key": 6,
        "Name": "[WorkloadDatastore] d7741e65-148e-b5ff-2458-068cdd1d8254/test-snap-000001.vmdk",
        "Type": "diskDescriptor",
        "Size": 54525952,
        "UniqueSize": 54525952,
        "BackingObjectId": "1c8e1e65-2fae-3944-7cd8-068cdd1d8254",
        "Accessible": true
      },
      {
        "Key": 8,
        "Name": "[WorkloadDatastore] d7741e65-148e-b5ff-2458-068cdd1d8254/test-snap-000002.vmdk",
        "Type": "diskDescriptor",
        "Size": 1195376640,
        "UniqueSize": 1195376640,
        "BackingObjectId": "319b1e65-6854-53c2-88f9-068cdd1d8254",
        "Accessible": true
      },
      {
        "Key": 12,
        "Name": "[WorkloadDatastore] d7741e65-148e-b5ff-2458-068cdd1d8254/test-snap-000003.vmdk",
        "Type": "diskDescriptor",
        "Size": 2264924160,
        "UniqueSize": 2264924160,
        "BackingObjectId": "20b51e65-8f9b-5fca-d03f-068cdd1d8254",
        "Accessible": true
      },
      {
        "Key": 14,
        "Name": "[WorkloadDatastore] d7741e65-148e-b5ff-2458-068cdd1d8254/test-snap-000004.vmdk",
        "Type": "diskDescriptor",
        "Size": 1778384896,
        "UniqueSize": 1778384896,
        "BackingObjectId": "66bb1e65-a6f9-6127-6954-068cdd1d8254",
        "Accessible": true
      },
      {
        "Key": 17,
        "Name": "[WorkloadDatastore] d7741e65-148e-b5ff-2458-068cdd1d8254/test-snap-000005.vmdk",
        "Type": "diskDescriptor",
        "Size": 62914560,
        "UniqueSize": 62914560,
        "BackingObjectId": "dddb1f65-42e3-0da1-f5b4-068cdd1d8254",
        "Accessible": true
      },
      {
        "Key": 22,
        "Name": "[WorkloadDatastore] d7741e65-148e-b5ff-2458-068cdd1d8254/test-snap-000007.vmdk",
        "Type": "diskDescriptor",
        "Size": 54525952,
        "UniqueSize": 54525952,
        "BackingObjectId": "c7fa1f65-3d92-3edd-2c35-068cdd1d8254",
        "Accessible": true
      },
      {
        "Key": 7,
        "Name": "[WorkloadDatastore] d7741e65-148e-b5ff-2458-068cdd1d8254/test-snap-Snapshot3.vmsn",
        "Type": "snapshotData",
        "Size": 282644,
        "UniqueSize": 282644,
        "BackingObjectId": "",
        "Accessible": true
      },
      {
        "Key": 20,
        "Name": "[WorkloadDatastore] d7741e65-148e-b5ff-2458-068cdd1d8254/test-snap-000006.vmdk",
        "Type": "diskDescriptor",
        "Size": 12582912,
        "UniqueSize": 12582912,
        "BackingObjectId": "b8fa1f65-2dac-5ca3-dffd-068cdd1d8254",
        "Accessible": true
      },
      {
        "Key": 10,
        "Name": "[WorkloadDatastore] d7741e65-148e-b5ff-2458-068cdd1d8254/test-snap-Snapshot4.vmem",
        "Type": "snapshotMemory",
        "Size": 2147483648,
        "UniqueSize": 2202009600,
        "BackingObjectId": "319b1e65-62eb-39b1-b4b4-068cdd1d8254",
        "Accessible": true
      },
      {
        "Key": 11,
        "Name": "[WorkloadDatastore] d7741e65-148e-b5ff-2458-068cdd1d8254/test-snap-Snapshot4.vmsn",
        "Type": "snapshotData",
        "Size": 12049359,
        "UniqueSize": 12049359,
        "BackingObjectId": "",
        "Accessible": true
      },
      {
        "Key": 13,
        "Name": "[WorkloadDatastore] d7741e65-148e-b5ff-2458-068cdd1d8254/test-snap-Snapshot5.vmsn",
        "Type": "snapshotData",
        "Size": 282651,
        "UniqueSize": 282651,
        "BackingObjectId": "",
        "Accessible": true
      },
      {
        "Key": 15,
        "Name": "[WorkloadDatastore] d7741e65-148e-b5ff-2458-068cdd1d8254/test-snap-Snapshot6.vmem",
        "Type": "snapshotMemory",
        "Size": 2147483648,
        "UniqueSize": 2202009600,
        "BackingObjectId": "66bb1e65-9e8b-3d10-8ca8-068cdd1d8254",
        "Accessible": true
      },
      {
        "Key": 16,
        "Name": "[WorkloadDatastore] d7741e65-148e-b5ff-2458-068cdd1d8254/test-snap-Snapshot6.vmsn",
        "Type": "snapshotData",
        "Size": 12054008,
        "UniqueSize": 12054008,
        "BackingObjectId": "",
        "Accessible": true
      },
      {
        "Key": 18,
        "Name": "[WorkloadDatastore] d7741e65-148e-b5ff-2458-068cdd1d8254/test-snap-Snapshot7.vmem",
        "Type": "snapshotMemory",
        "Size": 2147483648,
        "UniqueSize": 2202009600,
        "BackingObjectId": "dddb1f65-c58a-c485-98fd-068cdd1d8254",
        "Accessible": true
      },
      {
        "Key": 19,
        "Name": "[WorkloadDatastore] d7741e65-148e-b5ff-2458-068cdd1d8254/test-snap-Snapshot7.vmsn",
        "Type": "snapshotData",
        "Size": 12067011,
        "UniqueSize": 12067011,
        "BackingObjectId": "",
        "Accessible": true
      },
      {
        "Key": 21,
        "Name": "[WorkloadDatastore] d7741e65-148e-b5ff-2458-068cdd1d8254/test-snap-Snapshot8.vmsn",
        "Type": "snapshotData",
        "Size": 282651,
        "UniqueSize": 282651,
        "BackingObjectId": "",
        "Accessible": true
      },
      {
        "Key": 23,
        "Name": "[WorkloadDatastore] d7741e65-148e-b5ff-2458-068cdd1d8254/test-snap-Snapshot9.vmem",
        "Type": "snapshotMemory",
        "Size": 2147483648,
        "UniqueSize": 2202009600,
        "BackingObjectId": "c7fa1f65-1683-08bc-e4a9-068cdd1d8254",
        "Accessible": true
      },
      {
        "Key": 24,
        "Name": "[WorkloadDatastore] d7741e65-148e-b5ff-2458-068cdd1d8254/test-snap-Snapshot9.vmsn",
        "Type": "snapshotData",
        "Size": 12066999,
        "UniqueSize": 12066999,
        "BackingObjectId": "",
        "Accessible": true
      },
      {
        "Key": 4,
        "Name": "[WorkloadDatastore] d7741e65-148e-b5ff-2458-068cdd1d8254/test-snap-b52f7a84.vswp",
        "Type": "swap",
        "Size": 2147483648,
        "UniqueSize": 12582912,
        "BackingObjectId": "dd741e65-d34b-edea-9e04-068cdd1d8254",
        "Accessible": true
      },
      {
        "Key": 5,
        "Name": "[WorkloadDatastore] d7741e65-148e-b5ff-2458-068cdd1d8254/vmx-test-snap-5cb76317af187947bb1301e1bba4d8dc4171e03e74df4379920246ccedf02462-1.vswp",
        "Type": "uwswap",
        "Size": 83886080,
        "UniqueSize": 83886080,
        "BackingObjectId": "",
        "Accessible": true
      },
      {
        "Key": 9,
        "Name": "[WorkloadDatastore] d7741e65-148e-b5ff-2458-068cdd1d8254/vmware.log",
        "Type": "log",
        "Size": 395279,
        "UniqueSize": 395279,
        "BackingObjectId": "",
        "Accessible": true
      }
    ],
    "Disk": [
      {
        "Key": 2000,
        "Chain": [
          {
            "FileKey": [
              2
            ]
          },
          {
            "FileKey": [
              6
            ]
          },
          {
            "FileKey": [
              8
            ]
          },
          {
            "FileKey": [
              12
            ]
          },
          {
            "FileKey": [
              14
            ]
          },
          {
            "FileKey": [
              17
            ]
          },
          {
            "FileKey": [
              20
            ]
          },
          {
            "FileKey": [
              22
            ]
          }
        ]
      }
    ],
    "Snapshot": [
      {
        "Key": {
          "Type": "VirtualMachineSnapshot",
          "Value": "snapshot-1"
        },
        "DataKey": 7,
        "MemoryKey": -1,
        "Disk": [
          {
            "Key": 2000,
            "Chain": [
              {
                "FileKey": [
                  2
                ]
              }
            ]
          }
        ]
      },
      {
        "Key": {
          "Type": "VirtualMachineSnapshot",
          "Value": "snapshot-2"
        },
        "DataKey": 11,
        "MemoryKey": 10,
        "Disk": [
          {
            "Key": 2000,
            "Chain": [
              {
                "FileKey": [
                  2
                ]
              },
              {
                "FileKey": [
                  6
                ]
              }
            ]
          }
        ]
      },
      {
        "Key": {
          "Type": "VirtualMachineSnapshot",
          "Value": "snapshot-3"
        },
        "DataKey": 13,
        "MemoryKey": -1,
        "Disk": [
          {
            "Key": 2000,
            "Chain": [
              {
                "FileKey": [
                  2
                ]
              },
              {
                "FileKey": [
                  6
                ]
              },
              {
                "FileKey": [
                  8
                ]
              }
            ]
          }
        ]
      },
      {
        "Key": {
          "Type": "VirtualMachineSnapshot",
          "Value": "snapshot-4"
        },
        "DataKey": 16,
        "MemoryKey": 15,
        "Disk": [
          {
            "Key": 2000,
            "Chain": [
              {
                "FileKey": [
                  2
                ]
              },
              {
                "FileKey": [
                  6
                ]
              },
              {
                "FileKey": [
                  8
                ]
              },
              {
                "FileKey": [
                  12
                ]
              }
            ]
          }
        ]
      },
      {
        "Key": {
          "Type": "VirtualMachineSnapshot",
          "Value": "snapshot-5"
        },
        "DataKey": 19,
        "MemoryKey": 18,
        "Disk": [
          {
            "Key": 2000,
            "Chain": [
              {
                "FileKey": [
                  2
                ]
              },
              {
                "FileKey": [
                  6
                ]
              },
              {
                "FileKey": [
                  8
                ]
              },
              {
                "FileKey": [
                  12
                ]
              },
              {
                "FileKey": [
                  14
                ]
              }
            ]
          }
        ]
      },
      {
        "Key": {
          "Type": "VirtualMachineSnapshot",
          "Value": "snapshot-6"
        },
        "DataKey": 21,
        "MemoryKey": -1,
        "Disk": [
          {
            "Key": 2000,
            "Chain": [
              {
                "FileKey": [
                  2
                ]
              },
              {
                "FileKey": [
                  6
                ]
              },
              {
                "FileKey": [
                  8
                ]
              },
              {
                "FileKey": [
                  12
                ]
              },
              {
                "FileKey": [
                  14
                ]
              },
              {
                "FileKey": [
                  17
                ]
              }
            ]
          }
        ]
      },
      {
        "Key": {
          "Type": "VirtualMachineSnapshot",
          "Value": "snapshot-7"
        },
        "DataKey": 24,
        "MemoryKey": 23,
        "Disk": [
          {
            "Key": 2000,
            "Chain": [
              {
                "FileKey": [
                  2
                ]
              },
              {
                "FileKey": [
                  6
                ]
              },
              {
                "FileKey": [
                  8
                ]
              },
              {
                "FileKey": [
                  12
                ]
              },
              {
                "FileKey": [
                  14
                ]
              },
              {
                "FileKey": [
                  17
                ]
              },
              {
                "FileKey": [
                  20
                ]
              }
            ]
          }
        ]
      }
    ],
    "Timestamp": "0001-01-01T00:00:00Z"
  },
  "Storage": null,
  "EnvironmentBrowser": {
    "Type": "",
    "Value": ""
  },
  "ResourcePool": {
    "Type": "ResourcePool",
    "Value": "resgroup-41"
  },
  "ParentVApp": null,
  "ResourceConfig": null,
  "Runtime": {
    "Host": {
      "Type": "HostSystem",
      "Value": "host-14"
    },
    "ConnectionState": "connected",
    "PowerState": "poweredOn",
    "VmFailoverInProgress": null,
    "FaultToleranceState": "notConfigured",
    "DasVmProtection": {
      "DasProtected": true
    },
    "ToolsInstallerMounted": false,
    "SuspendTime": null,
    "BootTime": "2023-10-05T10:20:12.996473Z",
    "SuspendInterval": 0,
    "Question": null,
    "MemoryOverhead": 0,
    "MaxCpuUsage": 2299,
    "MaxMemoryUsage": 2048,
    "NumMksConnections": 0,
    "RecordReplayState": "inactive",
    "CleanPowerOff": null,
    "NeedSecondaryReason": "",
    "OnlineStandby": false,
    "MinRequiredEVCModeKey": "intel-broadwell",
    "ConsolidationNeeded": false,
    "OfflineFeatureRequirement": [
      {
        "Key": "cpuid.lm",
        "FeatureName": "cpuid.lm",
        "Value": "Num:Min:1"
      }
    ],
    "FeatureRequirement": [
      {
        "Key": "cpuid.3dnprefetch",
        "FeatureName": "cpuid.3dnprefetch",
        "Value": "Bool:Min:1"
      },
      {
        "Key": "cpuid.abm",
        "FeatureName": "cpuid.abm",
        "Value": "Bool:Min:1"
      },
      {
        "Key": "cpuid.adx",
        "FeatureName": "cpuid.adx",
        "Value": "Bool:Min:1"
      },
      {
        "Key": "cpuid.aes",
        "FeatureName": "cpuid.aes",
        "Value": "Bool:Min:1"
      },
      {
        "Key": "cpuid.avx",
        "FeatureName": "cpuid.avx",
        "Value": "Bool:Min:1"
      },
      {
        "Key": "cpuid.avx2",
        "FeatureName": "cpuid.avx2",
        "Value": "Bool:Min:1"
      },
      {
        "Key": "cpuid.bmi1",
        "FeatureName": "cpuid.bmi1",
        "Value": "Bool:Min:1"
      },
      {
        "Key": "cpuid.bmi2",
        "FeatureName": "cpuid.bmi2",
        "Value": "Bool:Min:1"
      },
      {
        "Key": "cpuid.cmpxchg16b",
        "FeatureName": "cpuid.cmpxchg16b",
        "Value": "Bool:Min:1"
      },
      {
        "Key": "cpuid.enfstrg",
        "FeatureName": "cpuid.enfstrg",
        "Value": "Bool:Min:1"
      },
      {
        "Key": "cpuid.f16c",
        "FeatureName": "cpuid.f16c",
        "Value": "Bool:Min:1"
      },
      {
        "Key": "cpuid.fcmd",
        "FeatureName": "cpuid.fcmd",
        "Value": "Bool:Min:1"
      },
      {
        "Key": "cpuid.fma",
        "FeatureName": "cpuid.fma",
        "Value": "Bool:Min:1"
      },
      {
        "Key": "cpuid.fsgsbase",
        "FeatureName": "cpuid.fsgsbase",
        "Value": "Bool:Min:1"
      },
      {
        "Key": "cpuid.ibpb",
        "FeatureName": "cpuid.ibpb",
        "Value": "Bool:Min:1"
      },
      {
        "Key": "cpuid.ibrs",
        "FeatureName": "cpuid.ibrs",
        "Value": "Bool:Min:1"
      },
      {
        "Key": "cpuid.intel",
        "FeatureName": "cpuid.intel",
        "Value": "Bool:Min:1"
      },
      {
        "Key": "cpuid.invpcid",
        "FeatureName": "cpuid.invpcid",
        "Value": "Bool:Min:1"
      },
      {
        "Key": "cpuid.lahf64",
        "FeatureName": "cpuid.lahf64",
        "Value": "Bool:Min:1"
      },
      {
        "Key": "cpuid.lm",
        "FeatureName": "cpuid.lm",
        "Value": "Bool:Min:1"
      },
      {
        "Key": "cpuid.mdclear",
        "FeatureName": "cpuid.mdclear",
        "Value": "Bool:Min:1"
      },
      {
        "Key": "cpuid.movbe",
        "FeatureName": "cpuid.movbe",
        "Value": "Bool:Min:1"
      },
      {
        "Key": "cpuid.nx",
        "FeatureName": "cpuid.nx",
        "Value": "Bool:Min:1"
      },
      {
        "Key": "cpuid.pcid",
        "FeatureName": "cpuid.pcid",
        "Value": "Bool:Min:1"
      },
      {
        "Key": "cpuid.pclmulqdq",
        "FeatureName": "cpuid.pclmulqdq",
        "Value": "Bool:Min:1"
      },
      {
        "Key": "cpuid.pdpe1gb",
        "FeatureName": "cpuid.pdpe1gb",
        "Value": "Bool:Min:1"
      },
      {
        "Key": "cpuid.popcnt",
        "FeatureName": "cpuid.popcnt",
        "Value": "Bool:Min:1"
      },
      {
        "Key": "cpuid.rdrand",
        "FeatureName": "cpuid.rdrand",
        "Value": "Bool:Min:1"
      },
      {
        "Key": "cpuid.rdseed",
        "FeatureName": "cpuid.rdseed",
        "Value": "Bool:Min:1"
      },
      {
        "Key": "cpuid.rdtscp",
        "FeatureName": "cpuid.rdtscp",
        "Value": "Bool:Min:1"
      },
      {
        "Key": "cpuid.smap",
        "FeatureName": "cpuid.smap",
        "Value": "Bool:Min:1"
      },
      {
        "Key": "cpuid.smep",
        "FeatureName": "cpuid.smep",
        "Value": "Bool:Min:1"
      },
      {
        "Key": "cpuid.ss",
        "FeatureName": "cpuid.ss",
        "Value": "Bool:Min:1"
      },
      {
        "Key": "cpuid.ssbd",
        "FeatureName": "cpuid.ssbd",
        "Value": "Bool:Min:1"
      },
      {
        "Key": "cpuid.sse3",
        "FeatureName": "cpuid.sse3",
        "Value": "Bool:Min:1"
      },
      {
        "Key": "cpuid.sse41",
        "FeatureName": "cpuid.sse41",
        "Value": "Bool:Min:1"
      },
      {
        "Key": "cpuid.sse42",
        "FeatureName": "cpuid.sse42",
        "Value": "Bool:Min:1"
      },
      {
        "Key": "cpuid.ssse3",
        "FeatureName": "cpuid.ssse3",
        "Value": "Bool:Min:1"
      },
      {
        "Key": "cpuid.stibp",
        "FeatureName": "cpuid.stibp",
        "Value": "Bool:Min:1"
      },
      {
        "Key": "cpuid.xcr0_master_sse",
        "FeatureName": "cpuid.xcr0_master_sse",
        "Value": "Bool:Min:1"
      },
      {
        "Key": "cpuid.xcr0_master_ymm_h",
        "FeatureName": "cpuid.xcr0_master_ymm_h",
        "Value": "Bool:Min:1"
      },
      {
        "Key": "cpuid.xsave",
        "FeatureName": "cpuid.xsave",
        "Value": "Bool:Min:1"
      },
      {
        "Key": "cpuid.xsaveopt",
        "FeatureName": "cpuid.xsaveopt",
        "Value": "Bool:Min:1"
      }
    ],
    "FeatureMask": null,
    "VFlashCacheAllocation": 0,
    "Paused": false,
    "SnapshotInBackground": false,
    "QuiescedForkParent": null,
    "InstantCloneFrozen": false,
    "CryptoState": "",
    "SuspendedToMemory": null,
    "OpNotificationTimeout": 0
  },
  "Guest": {
    "ToolsStatus": "toolsNotInstalled",
    "ToolsVersionStatus": "guestToolsNotInstalled",
    "ToolsVersionStatus2": "guestToolsNotInstalled",
    "ToolsRunningStatus": "guestToolsNotRunning",
    "ToolsVersion": "0",
    "ToolsInstallType": "guestToolsTypeUnknown",
    "GuestId": "",
    "GuestFamily": "",
    "GuestFullName": "",
    "HostName": "",
    "IpAddress": "",
    "Net": null,
    "IpStack": null,
    "Disk": null,
    "Screen": {
      "Width": 1024,
      "Height": 768
    },
    "GuestState": "notRunning",
    "AppHeartbeatStatus": "appStatusGray",
    "GuestKernelCrashed": false,
    "AppState": "none",
    "GuestOperationsReady": false,
    "InteractiveGuestOperationsReady": false,
    "GuestStateChangeSupported": false,
    "GenerationInfo": null,
    "HwVersion": "vmx-14",
    "CustomizationInfo": {
      "CustomizationStatus": "TOOLSDEPLOYPKG_IDLE",
      "StartTime": null,
      "EndTime": null,
      "ErrorMsg": ""
    }
  },
  "Summary": {
    "Vm": {
      "Type": "VirtualMachine",
      "Value": "vm-3"
    },
    "Runtime": {
      "Host": {
        "Type": "HostSystem",
        "Value": "host-14"
      },
      "ConnectionState": "connected",
      "PowerState": "poweredOn",
      "VmFailoverInProgress": null,
      "FaultToleranceState": "notConfigured",
      "DasVmProtection": {
        "DasProtected": true
      },
      "ToolsInstallerMounted": false,
      "SuspendTime": null,
      "BootTime": "2023-10-05T10:20:12.996473Z",
      "SuspendInterval": 0,
      "Question": null,
      "MemoryOverhead": 0,
      "MaxCpuUsage": 2299,
      "MaxMemoryUsage": 2048,
      "NumMksConnections": 0,
      "RecordReplayState": "inactive",
      "CleanPowerOff": null,
      "NeedSecondaryReason": "",
      "OnlineStandby": false,
      "MinRequiredEVCModeKey": "intel-broadwell",
      "ConsolidationNeeded": false,
      "OfflineFeatureRequirement": [
        {
          "Key": "cpuid.lm",
          "FeatureName": "cpuid.lm",
          "Value": "Num:Min:1"
        }
      ],
      "FeatureRequirement": [
        {
          "Key": "cpuid.3dnprefetch",
          "FeatureName": "cpuid.3dnprefetch",
          "Value": "Bool:Min:1"
        },
        {
          "Key": "cpuid.abm",
          "FeatureName": "cpuid.abm",
          "Value": "Bool:Min:1"
        },
        {
          "Key": "cpuid.adx",
          "FeatureName": "cpuid.adx",
          "Value": "Bool:Min:1"
        },
        {
          "Key": "cpuid.aes",
          "FeatureName": "cpuid.aes",
          "Value": "Bool:Min:1"
        },
        {
          "Key": "cpuid.avx",
          "FeatureName": "cpuid.avx",
          "Value": "Bool:Min:1"
        },
        {
          "Key": "cpuid.avx2",
          "FeatureName": "cpuid.avx2",
          "Value": "Bool:Min:1"
        },
        {
          "Key": "cpuid.bmi1",
          "FeatureName": "cpuid.bmi1",
          "Value": "Bool:Min:1"
        },
        {
          "Key": "cpuid.bmi2",
          "FeatureName": "cpuid.bmi2",
          "Value": "Bool:Min:1"
        },
        {
          "Key": "cpuid.cmpxchg16b",
          "FeatureName": "cpuid.cmpxchg16b",
          "Value": "Bool:Min:1"
        },
        {
          "Key": "cpuid.enfstrg",
          "FeatureName": "cpuid.enfstrg",
          "Value": "Bool:Min:1"
        },
        {
          "Key": "cpuid.f16c",
          "FeatureName": "cpuid.f16c",
          "Value": "Bool:Min:1"
        },
        {
          "Key": "cpuid.fcmd",
          "FeatureName": "cpuid.fcmd",
          "Value": "Bool:Min:1"
        },
        {
          "Key": "cpuid.fma",
          "FeatureName": "cpuid.fma",
          "Value": "Bool:Min:1"
        },
        {
          "Key": "cpuid.fsgsbase",
          "FeatureName": "cpuid.fsgsbase",
          "Value": "Bool:Min:1"
        },
        {
          "Key": "cpuid.ibpb",
          "FeatureName": "cpuid.ibpb",
          "Value": "Bool:Min:1"
        },
        {
          "Key": "cpuid.ibrs",
          "FeatureName": "cpuid.ibrs",
          "Value": "Bool:Min:1"
        },
        {
          "Key": "cpuid.intel",
          "FeatureName": "cpuid.intel",
          "Value": "Bool:Min:1"
        },
        {
          "Key": "cpuid.invpcid",
          "FeatureName": "cpuid.invpcid",
          "Value": "Bool:Min:1"
        },
        {
          "Key": "cpuid.lahf64",
          "FeatureName": "cpuid.lahf64",
          "Value": "Bool:Min:1"
        },
        {
          "Key": "cpuid.lm",
          "FeatureName": "cpuid.lm",
          "Value": "Bool:Min:1"
        },
        {
          "Key": "cpuid.mdclear",
          "FeatureName": "cpuid.mdclear",
          "Value": "Bool:Min:1"
        },
        {
          "Key": "cpuid.movbe",
          "FeatureName": "cpuid.movbe",
          "Value": "Bool:Min:1"
        },
        {
          "Key": "cpuid.nx",
          "FeatureName": "cpuid.nx",
          "Value": "Bool:Min:1"
        },
        {
          "Key": "cpuid.pcid",
          "FeatureName": "cpuid.pcid",
          "Value": "Bool:Min:1"
        },
        {
          "Key": "cpuid.pclmulqdq",
          "FeatureName": "cpuid.pclmulqdq",
          "Value": "Bool:Min:1"
        },
        {
          "Key": "cpuid.pdpe1gb",
          "FeatureName": "cpuid.pdpe1gb",
          "Value": "Bool:Min:1"
        },
        {
          "Key": "cpuid.popcnt",
          "FeatureName": "cpuid.popcnt",
          "Value": "Bool:Min:1"
        },
        {
          "Key": "cpuid.rdrand",
          "FeatureName": "cpuid.rdrand",
          "Value": "Bool:Min:1"
        },
        {
          "Key": "cpuid.rdseed",
          "FeatureName": "cpuid.rdseed",
          "Value": "Bool:Min:1"
        },
        {
          "Key": "cpuid.rdtscp",
          "FeatureName": "cpuid.rdtscp",
          "Value": "Bool:Min:1"
        },
        {
          "Key": "cpuid.smap",
          "FeatureName": "cpuid.smap",
          "Value": "Bool:Min:1"
        },
        {
          "Key": "cpuid.smep",
          "FeatureName": "cpuid.smep",
          "Value": "Bool:Min:1"
        },
        {
          "Key": "cpuid.ss",
          "FeatureName": "cpuid.ss",
          "Value": "Bool:Min:1"
        },
        {
          "Key": "cpuid.ssbd",
          "FeatureName": "cpuid.ssbd",
          "Value": "Bool:Min:1"
        },
        {
          "Key": "cpuid.sse3",
          "FeatureName": "cpuid.sse3",
          "Value": "Bool:Min:1"
        },
        {
          "Key": "cpuid.sse41",
          "FeatureName": "cpuid.sse41",
          "Value": "Bool:Min:1"
        },
        {
          "Key": "cpuid.sse42",
          "FeatureName": "cpuid.sse42",
          "Value": "Bool:Min:1"
        },
        {
          "Key": "cpuid.ssse3",
          "FeatureName": "cpuid.ssse3",
          "Value": "Bool:Min:1"
        },
        {
          "Key": "cpuid.stibp",
          "FeatureName": "cpuid.stibp",
          "Value": "Bool:Min:1"
        },
        {
          "Key": "cpuid.xcr0_master_sse",
          "FeatureName": "cpuid.xcr0_master_sse",
          "Value": "Bool:Min:1"
        },
        {
          "Key": "cpuid.xcr0_master_ymm_h",
          "FeatureName": "cpuid.xcr0_master_ymm_h",
          "Value": "Bool:Min:1"
        },
        {
          "Key": "cpuid.xsave",
          "FeatureName": "cpuid.xsave",
          "Value": "Bool:Min:1"
        },
        {
          "Key": "cpuid.xsaveopt",
          "FeatureName": "cpuid.xsaveopt",
          "Value": "Bool:Min:1"
        }
      ],
      "FeatureMask": null,
      "VFlashCacheAllocation": 0,
      "Paused": false,
      "SnapshotInBackground": false,
      "QuiescedForkParent": null,
      "InstantCloneFrozen": false,
      "CryptoState": "",
      "SuspendedToMemory": null,
      "OpNotificationTimeout": 0
    },
    "Guest": {
      "GuestId": "",
      "GuestFullName": "",
      "ToolsStatus": "toolsNotInstalled",
      "ToolsVersionStatus": "guestToolsNotInstalled",
      "ToolsVersionStatus2": "guestToolsNotInstalled",
      "ToolsRunningStatus": "guestToolsNotRunning",
      "HostName": "",
      "IpAddress": "",
      "HwVersion": "vmx-14"
    },
    "Config": {
      "Name": "test-snap",
      "Template": false,
      "VmPathName": "[WorkloadDatastore] d7741e65-148e-b5ff-2458-068cdd1d8254/test-snap.vmx",
      "MemorySizeMB": 2048,
      "CpuReservation": 0,
      "MemoryReservation": 0,
      "NumCpu": 1,
      "NumEthernetCards": 1,
      "NumVirtualDisks": 1,
      "Uuid": "420fe454-848e-4d0b-a08f-c2423606b46c",
      "InstanceUuid": "500f7e58-8c6b-2aae-3ffb-ea87f750a55d",
      "GuestId": "centos8_64Guest",
      "GuestFullName": "CentOS 8 (64-bit)",
      "Annotation": "",
      "Product": null,
      "InstallBootRequired": false,
      "FtInfo": null,
      "ManagedBy": null,
      "TpmPresent": false,
      "NumVmiopBackings": 0,
      "HwVersion": "vmx-14"
    },
    "Storage": {
      "Committed": 16684941312,
      "Uncommitted": 131931832320,
      "Unshared": 7642021888,
      "Timestamp": "2023-10-06T13:15:48.131079Z"
    },
    "QuickStats": {
      "OverallCpuUsage": 0,
      "OverallCpuDemand": 0,
      "OverallCpuReadiness": 0,
      "GuestMemoryUsage": 20,
      "HostMemoryUsage": 2087,
      "GuestHeartbeatStatus": "gray",
      "DistributedCpuEntitlement": 0,
      "DistributedMemoryEntitlement": 1479,
      "StaticCpuEntitlement": 2299,
      "StaticMemoryEntitlement": 2999,
      "GrantedMemory": 2042,
      "PrivateMemory": 2042,
      "SharedMemory": 0,
      "SwappedMemory": 0,
      "BalloonedMemory": 0,
      "ConsumedOverheadMemory": 45,
      "FtLogBandwidth": -1,
      "FtSecondaryLatency": -1,
      "FtLatencyStatus": "gray",
      "CompressedMemory": 0,
      "UptimeSeconds": 103491,
      "SsdSwappedMemory": 0,
      "ActiveMemory": 20,
      "MemoryTierStats": null
    },
    "OverallStatus": "green",
    "CustomValue": null
  },
  "Datastore": [
    {
      "Type": "Datastore",
      "Value": "datastore-40"
    }
  ],
  "Network": [
    {
      "Type": "DistributedVirtualPortgroup",
      "Value": "dvportgroup-46"
    }
  ],
  "Snapshot": {
    "CurrentSnapshot": {
      "Type": "VirtualMachineSnapshot",
      "Value": "snapshot-7"
    },
    "RootSnapshotList": [
      {
        "Snapshot": {
          "Type": "VirtualMachineSnapshot",
          "Value": "snapshot-1"
        },
        "Vm": {
          "Type": "VirtualMachine",
          "Value": "vm-3"
        },
        "Name": "Snap1",
        "Description": "",
        "Id": 3,
        "CreateTime": "2023-10-05T10:21:16.766578Z",
        "State": "poweredOff",
        "Quiesced": false,
        "BackupManifest": "",
        "ChildSnapshotList": [
          {
            "Snapshot": {
              "Type": "VirtualMachineSnapshot",
              "Value": "snapshot-2"
            },
            "Vm": {
              "Type": "VirtualMachine",
              "Value": "vm-3"
            },
            "Name": "Snap2",
            "Description": "",
            "Id": 4,
            "CreateTime": "2023-10-05T11:17:05.637214Z",
            "State": "poweredOn",
            "Quiesced": false,
            "BackupManifest": "",
            "ChildSnapshotList": [
              {
                "Snapshot": {
                  "Type": "VirtualMachineSnapshot",
                  "Value": "snapshot-3"
                },
                "Vm": {
                  "Type": "VirtualMachine",
                  "Value": "vm-3"
                },
                "Name": "Snap3",
                "Description": "",
                "Id": 5,
                "CreateTime": "2023-10-05T13:07:44.044097Z",
                "State": "poweredOff",
                "Quiesced": false,
                "BackupManifest": "",
                "ChildSnapshotList": [
                  {
                    "Snapshot": {
                      "Type": "VirtualMachineSnapshot",
                      "Value": "snapshot-4"
                    },
                    "Vm": {
                      "Type": "VirtualMachine",
                      "Value": "vm-3"
                    },
                    "Name": "Snap4",
                    "Description": "",
                    "Id": 6,
                    "CreateTime": "2023-10-05T13:34:30.441561Z",
                    "State": "poweredOn",
                    "Quiesced": false,
                    "BackupManifest": "",
                    "ChildSnapshotList": [
                      {
                        "Snapshot": {
                          "Type": "VirtualMachineSnapshot",
                          "Value": "snapshot-5"
                        },
                        "Vm": {
                          "Type": "VirtualMachine",
                          "Value": "vm-3"
                        },
                        "Name": "Snap5",
                        "Description": "",
                        "Id": 7,
                        "CreateTime": "2023-10-06T10:05:17.201285Z",
                        "State": "poweredOn",
                        "Quiesced": false,
                        "BackupManifest": "",
                        "ChildSnapshotList": [
                          {
                            "Snapshot": {
                              "Type": "VirtualMachineSnapshot",
                              "Value": "snapshot-6"
                            },
                            "Vm": {
                              "Type": "VirtualMachine",
                              "Value": "vm-3"
                            },
                            "Name": "Snap6",
                            "Description": "",
                            "Id": 8,
                            "CreateTime": "2023-10-06T12:16:56.053819Z",
                            "State": "poweredOff",
                            "Quiesced": false,
                            "BackupManifest": "",
                            "ChildSnapshotList": [
                              {
                                "Snapshot": {
                                  "Type": "VirtualMachineSnapshot",
                                  "Value": "snapshot-7"
                                },
                                "Vm": {
                                  "Type": "VirtualMachine",
                                  "Value": "vm-3"
                                },
                                "Name": "Snap7",
                                "Description": "",
                                "Id": 9,
                                "CreateTime": "2023-10-06T12:17:11.193681Z",
                                "State": "poweredOn",
                                "Quiesced": false,
                                "BackupManifest": "",
                                "ChildSnapshotList": null,
                                "ReplaySupported": false
                              }
                            ],
                            "ReplaySupported": false
                          }
                        ],
                        "ReplaySupported": false
                      }
                    ],
                    "ReplaySupported": false
                  }
                ],
                "ReplaySupported": false
              }
            ],
            "ReplaySupported": false
          }
        ],
        "ReplaySupported": false
      }
    ]
  },
  "RootSnapshot": null,
  "GuestHeartbeatStatus": ""
}
`
