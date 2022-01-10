# Change Log
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

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
