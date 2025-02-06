# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

Unreleased section should follow [Release Toolkit](https://github.com/newrelic/release-toolkit#render-markdown-and-update-markdown)

## Unreleased

## v1.6.3 - 2025-02-06

### ⛓️ Dependencies
- Updated golang patch version to v1.23.6

## v1.6.2 - 2025-01-30

### ⛓️ Dependencies
- Updated golang patch version to v1.23.5

## v1.6.1 - 2024-12-05

### ⛓️ Dependencies
- Updated golang patch version to v1.23.4

## v1.6.0 - 2024-10-10

### dependency
- Upgrade go to 1.23.2

### 🚀 Enhancements
- Upgrade integrations SDK so the interval is variable and allows intervals up to 5 minutes

## v1.5.4 - 2024-09-12

### ⛓️ Dependencies
- Updated golang version to v1.23.1

## v1.5.3 - 2024-07-04

### ⛓️ Dependencies
- Updated golang version to v1.22.5

## v1.5.2 - 2024-05-09

### ⛓️ Dependencies
- Updated golang version to v1.22.3

## v1.5.1 - 2024-04-11

### ⛓️ Dependencies
- Updated golang version

## v1.5.0 - 2024-04-09

### 🛡️ Security notices
- Updated dependencies

## v1.4.1 (2024-02-28)
### Fixed
 - Reverted a new metric that was reporting a FQDN different from the one reported by the infra-agent

## v1.3.0 (2023-10-10)
### Changed
 - Dependencies have been updated: testify, logrus, govmomi
 - The way snapshot sizes are computed has been refactored in order to take into account delta disks


## v1.2.6 (2023-01-05)
### Changed
- Bump dependencies
- Bump go version 1.19
### Fix
- Event timestamp is now set explicitly to avoid the agent to set it to time.Now() instead of e.CreatedTime

## v1.2.5 (2022-06-23)

### Changed
 - Bump dependencies
### Added
 - Added support for more distributions:
    RHEL(EL) 9
    Ubuntu 22.04

## v1.2.4 (2022-05-05)

### Fix
* Performance metrics were not taking into consideration instanceless values by @paologallinaharbur in https://github.com/newrelic/nri-vsphere/pull/120
* DatacenterName attribute was missing from DatacenterSample by @paologallinaharbur in https://github.com/newrelic/nri-vsphere/pull/120

### Dependencies
* bump github.com/newrelic/infra-integrations-sdk from 3.7.0+incompatible to 3.7.1+incompatible by @dependabot in https://github.com/newrelic/nri-vsphere/pull/110
* bump github.com/newrelic/infra-integrations-sdk from 3.7.1+incompatible to 3.7.2+incompatible by @dependabot in https://github.com/newrelic/nri-vsphere/pull/114
* bump github.com/vmware/govmomi from 0.27.2 to 0.27.3 by @dependabot in https://github.com/newrelic/nri-vsphere/pull/112
* bump github.com/vmware/govmomi from 0.27.3 to 0.27.4 by @dependabot in https://github.com/newrelic/nri-vsphere/pull/113
* bump github.com/vmware/govmomi from 0.27.4 to 0.28.0 by @dependabot in https://github.com/newrelic/nri-vsphere/pull/118
* bump github.com/stretchr/testify from 1.7.0 to 1.7.1 by @dependabot in https://github.com/newrelic/nri-vsphere/pull/115
* bump goversion from 1.14 to 1.18  by @paologallinaharbur in https://github.com/newrelic/nri-vsphere/pull/119


**Full Changelog**: https://github.com/newrelic/nri-vsphere/compare/v1.2.3...v1.2.4

## v1.2.3 (2022-01-10)

### Changed

- bump github.com/vmware/govmomi from 0.27.1 to 0.27.2.

## v1.2.2 (2021-11-16)

### Changed

- Upgrade dependency version for govmomi, testify, yaml.v2.

## v1.2.1 (2021-03-24)

### Changed

- Added arm packages and binaries
- fix: List separator missing (#85)

## v1.2.0 (2020-09-26)

### Changed

- Performance metrics are now retrieved taking in consideration all instances
- In case multiple instances are mapped to a single newrelic DataSample the average is computed (ES: one metric per eare CPU core of an host)

### Others

- Now in case the instance name is specified as "*" when retrieving performance metrics.
These change allows retrieving the performance metrics that were not available yet.
