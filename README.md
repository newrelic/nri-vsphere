[![Community Project header](https://github.com/newrelic/opensource-website/raw/master/src/images/categories/Community_Project.png)](https://opensource.newrelic.com/oss-category/#community-project)

# New Relic VMware vShpere integration [![Known Vulnerabilities](https://snyk.io/test/github/newrelic/nri-vsphere/badge.svg?targetFile=vendor/vendor.json)](https://snyk.io/test/github/newrelic/nri-vsphere?targetFile=vendor/vendor.json)

New Relic's VMware vShpere integration captures critical summary and performance metrics data by connecting to VMware vCenter or an ESXi Host. For more information, see [the official New Relic docs](https://docs.newrelic.com/docs/integrations/host-integrations/host-integrations-list/vmware-vsphere-monitoring-integration).

The integration captures data about Datacenters, Clusters, Virtual Machines, Hosts, Datastores, ResourcePools and Networks. Apart from metrics, the integration also can capture vSphere events and VM snapshot information when enabled by the appropriate flags.

## Requirements

- [New Relic Infrastructure agent version 1.8.0 or higher](https://docs.newrelic.com/docs/infrastructure/install-configure-manage-infrastructure)
- New Relic Infrastructure Pro subscription or trial

## Installation

To install the integration, follow the official [documentation](https://docs.newrelic.com/docs/integrations/host-integrations/host-integrations-list/vmware-vsphere-monitoring-integration). We recommend using your operating system native package manager.

## Building

If you have the source code and Go toolchain installed, you can build and run the vSphere integration locally.

vSphere data is obtained via a client provided by the `govmomi` package, the official go wrapper for VMware API. The same package provides a simulator for the virtual center, leveraged to run integration tests.

The output of the integration is shaped by the `newrelic/infra-integrations-sdk` package and provided in JSON format. The New Relic Infra agent collects such data and sends it to New Relic.

1. After cloning this repository, go to vSphere integration directory and build the integration:

    ```bash
    make test
    make compile-all
    ```

    The command above executes the tests for the vSphere integration and builds an executable file named `nri-vsphere` under `bin/{architecture}`.

    You will also need `docker` to run the integration tests.

2. Run the executable with the following arguments:

    ```bash
    ./bin/darwin/nri-vsphere --url 127.0.0.1:8989/sdk --user user --pass pass --validate_ssl false
    ```

To learn more about the usage of `./bin/darwin/nri-vsphere`, pass the `-help` argument.

```bash
./bin/darwin/nri-vsphere -help
```

External dependencies are managed via govendor.

## Support

New Relic hosts and moderates an online forum where customers can interact with New Relic employees as well as other customers to get help and share best practices. Like all official New Relic open source projects, there's a related Community topic in the New Relic Explorers Hub. You can find this project's topic/threads here:

[Community support thread](https://discuss.newrelic.com/t/new-relic-vmware-vsphere-integration/)

## Contributing

We encourage your contributions to improve New Relic's VMware vShpere integration. Keep in mind when you submit your pull request, you'll need to sign the CLA via the click-through using CLA-Assistant. You only have to sign the CLA one time per project.
If you have any questions, or to execute our corporate CLA, required if your contribution is on behalf of a company,  please drop us an email at opensource@newrelic.com.

Before submitting a pull request, please review [these guidelines](https://github.com/newrelic/nri-vsphere/blob/master/CONTRIBUTING.md).

## License

New Relic's VMware vShpere integration is licensed under the [Apache 2.0](http://apache.org/licenses/LICENSE-2.0.txt) License.

New Relic Infrastructure Integration for VMware vSphere also uses source code from third-party libraries. You can find full details on which libraries are used and the terms under which they are licensed in the [third-party notices document](https://github.com/newrelic/nri-vsphere/blob/master/THIRD_PARTY_NOTICES.md).
