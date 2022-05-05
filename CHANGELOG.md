# Change Log
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

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
