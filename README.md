[![Community Project header](https://github.com/newrelic/opensource-website/raw/master/src/images/categories/Community_Project.png)](https://opensource.newrelic.com/oss-category/#community-project)

# New Relic VMware vSphere integration

[![Known Vulnerabilities](https://snyk.io/test/github/newrelic/nri-vsphere/badge.svg?targetFile=vendor/vendor.json)](https://snyk.io/test/github/newrelic/nri-vsphere?targetFile=vendor/vendor.json)

New Relic's VMware vSphere integration captures critical summary and performance metrics data by connecting to VMware vCenter or an ESXi Host. For more information, see [the official New Relic docs](https://docs.newrelic.com/docs/integrations/host-integrations/host-integrations-list/vmware-vsphere-monitoring-integration).

The integration collects data about data centers, clusters, virtual machines, hosts, datastores, resource pools, and networks. In addition to metrics, the integration can also capture vSphere events and VM snapshot information when enabled by the appropriate flags.

## Requirements

- [New Relic Infrastructure agent version 1.8.0 or higher](https://docs.newrelic.com/docs/infrastructure/install-configure-manage-infrastructure)
- New Relic Infrastructure Pro subscription or trial

## Installation

To install the integration, follow the official [documentation](https://docs.newrelic.com/docs/integrations/host-integrations/host-integrations-list/vmware-vsphere-monitoring-integration). We recommend using your operating system's package manager.

## Getting started

After you've [installed](#installation) the integration make sure that you have the required configuration for your environment.

To configure the integration go to `/etc/newrelic-infra/integrations.d/` (Linux) or `C:\Program Files\New Relic\newrelic-infra\integrations.d\` (Windows) and open the`vpshere-config.yml` configuration file.

Configure the `URL`, `user`, and `password` fields -- they are required to connect to your vCenter or ESXi host.

To select which performance metrics to capture, you must define them in the `vsphere-performance.metrics` file per each `performance level` you require.
You can find this file in `/etc/newrelic-infra/integrations.d/vsphere-performance.metrics` (Linux) or `C:\Program Files\New Relic\newrelic-infra\integrations.d\vsphere-performance.metrics` (Windows).
Use the flag `--perf_level` to select which level of **performance metrics** you want to capture.

Please note that the more performance metrics you enable the more load you add to your environment.

## Building

If you have downloaded the source code and installed the Go toolchain, you can build and run the vSphere integration locally.

vSphere data is obtained via a client provided by the `govmomi` package, the official Go wrapper for the VMware API. The same package provides a simulator for the virtual center, leveraged to run integration tests.

The output of the integration is determined by the `newrelic/infra-integrations-sdk` package and provided in JSON format. The [New Relic Infrastructure agent](https://github.com/newrelic/infrastructure-agent) collects such data and sends it to New Relic.

1. After cloning this repository, go to the vSphere integration directory and build the integration:

    ```bash
    make compile
    ```

    The command above executes the tests for the vSphere integration and builds an executable file named `nri-vsphere` under `bin/`.

2. Run the executable with the following arguments:

    ```bash
    ./bin/darwin/nri-vsphere --url 127.0.0.1:8989/sdk --user user --pass pass --validate_ssl false
    ```

To learn more about the usage of `./bin/darwin/nri-vsphere`, pass the `-help` argument.

```bash
./bin/darwin/nri-vsphere -help
```

External dependencies are managed via `govendor`.

## Testing

After cloning this repository, go to vSphere integration directory and build the integration:

```bash
make test
```

You need `docker-compose` to run the integration tests.

## Support

Should you need assistance with New Relic products, you are in good hands with several support channels.

**Support Channels**
* [New Relic Documentation](https://docs.newrelic.com/docs/integrations/host-integrations/host-integrations-list/vmware-vsphere-monitoring-integration): Comprehensive guidance for using our platform
* [New Relic Community](https://discuss.newrelic.com/t/new-relic-vmware-vsphere-integration/): The best place to engage in troubleshooting questions
* [New Relic Developer](https://developer.newrelic.com/): Resources for building a custom observability applications
* [New Relic University](https://learn.newrelic.com/): A range of online training for New Relic users of every level

## Privacy
At New Relic we take your privacy and the security of your information seriously, and are committed to protecting your information. We must emphasize the importance of not sharing personal data in public forums, and ask all users to scrub logs and diagnostic information for sensitive information, whether personal, proprietary, or otherwise.

We define “Personal Data” as any information relating to an identified or identifiable individual, including, for example, your name, phone number, post code or zip code, Device ID, IP address and email address.

Please review [New Relic’s General Data Privacy Notice](https://newrelic.com/termsandconditions/privacy) for more information.

## Contributing

We encourage your contributions to improve New Relic's VMware vSphere integration. Keep in mind when you submit your pull request, you'll need to sign the CLA via the click-through using CLA-Assistant. You only have to sign the CLA one time per project.
If you have any questions, or to execute our corporate CLA, required if your contribution is on behalf of a company,  please drop us an email at opensource@newrelic.com.

**A note about vulnerabilities**

As noted in our [security policy](/SECURITY.md), New Relic is committed to the privacy and security of our customers and their data. We believe that providing coordinated disclosure by security researchers and engaging with the security community are important means to achieve our security goals.

If you believe you have found a security vulnerability in this project or any of New Relic's products or websites, we welcome and greatly appreciate you reporting it to New Relic through [HackerOne](https://hackerone.com/newrelic).

If you would like to contribute to this project, please review [these guidelines](./CONTRIBUTING.md).

To [all contributors](https://github.com/newrelic/nri-vsphere/graphs/contributors), we thank you!  Without your contribution, this project would not be what it is today.  We also host a community project page dedicated to 
the [New Relic's VMware vSphere integration](https://opensource.newrelic.com/projects/newrelic/nri-vsphere).

## License

New Relic's VMware vSphere integration is licensed under the [Apache 2.0](http://apache.org/licenses/LICENSE-2.0.txt) License.

New Relic Infrastructure Integration for VMware vSphere also uses source code from third-party libraries. You can find full details on which libraries are used and the terms under which they are licensed in the [third-party notices document](https://github.com/newrelic/nri-vsphere/blob/master/THIRD_PARTY_NOTICES.md).
