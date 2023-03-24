<a href="https://opensource.newrelic.com/oss-category/#community-plus"><picture><source media="(prefers-color-scheme: dark)" srcset="https://github.com/newrelic/opensource-website/raw/main/src/images/categories/dark/Community_Plus.png"><source media="(prefers-color-scheme: light)" srcset="https://github.com/newrelic/opensource-website/raw/main/src/images/categories/Community_Plus.png"><img alt="New Relic Open Source community plus project banner." src="https://github.com/newrelic/opensource-website/raw/main/src/images/categories/Community_Plus.png"></picture></a>


# New Relic VMware vSphere integration

New Relic's VMware vSphere integration captures critical summary and performance metrics data by connecting to VMware vCenter or an ESXi Host. For more information, see [the official New Relic docs](https://docs.newrelic.com/docs/integrations/host-integrations/host-integrations-list/vmware-vsphere-monitoring-integration).

The integration collects data about data centers, clusters, virtual machines, hosts, datastores, resource pools, and networks. In addition to metrics, the integration can also capture vSphere events and VM snapshot information when enabled by the appropriate flags.

## Requirements

- [New Relic infrastructure agent version 1.8.0 or higher](https://docs.newrelic.com/docs/infrastructure/install-configure-manage-infrastructure)

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

Notice that the integration fetches multiple values for a single performance metrics related to different "instances" 
belonging to a single object, but only the average value is stored.

For example, the counter `cpu.usage.average` returns multiple values: one for each CPU core of an host.
The integration uses these values to compute the average, that is then included in the `VSphereHostSample` sample.

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

## Running the integration locally with a vCenter simulator

For testing purposes there is the possibility to run a docker compose cluster with 2 containers:
- One NewRelic Infrastructure agent with the integration installed.
- One container with [VCSIM](https://github.com/vmware/govmomi/tree/master/vcsim) running with the defaults to simulate a VCenter in port 8989.

This will emit VSphere metrics to NRONE based on the license key provided. 

A License Key (NRIA_LICENSE_KEY) env var must be provided.

You can also set the following env vars:
- STAGING   # To use a staging account 
- VS_HOSTS  # To set the number of simulated hosts per cluster (default: 1)
- VS_VMS    # To set the number of simulated Virtual Machines per resource pool (default: 10)
- VS_DS     # To set the number of simulated Data stores (default: 4)

Example:

```bash
NRIA_LICENSE_KEY=xxx make tools-vcsim-run
```

Example with hosts, vm, and ds custom:

```bash
NRIA_LICENSE_KEY=xxx VS_HOSTS=2 VS_VMS=20 VS_DS=10 make tools-vcsim-run
```

In order to stop the cluster you will need to run the following command:

```bash
make tools-vcsim-stop
```

## Support

Should you need assistance with New Relic products, you are in good hands with several support channels.


>New Relic offers NRDiag, [a client-side diagnostic utility](https://docs.newrelic.com/docs/using-new-relic/cross-product-functions/troubleshooting/new-relic-diagnostics) that automatically detects common problems with New Relic agents. If NRDiag detects a problem, it suggests troubleshooting steps. NRDiag can also automatically attach troubleshooting data to a New Relic Support ticket. Remove this section if it doesn't apply.

If the issue has been confirmed as a bug or is a feature request, file a GitHub issue.

**Support Channels**
* [New Relic Documentation](https://docs.newrelic.com/docs/integrations/host-integrations/host-integrations-list/vmware-vsphere-monitoring-integration): Comprehensive guidance for using our platform
* [New Relic Community](https://discuss.newrelic.com/t/new-relic-vmware-vsphere-integration/): The best place to engage in troubleshooting questions
* [New Relic Developer](https://developer.newrelic.com/): Resources for building a custom observability applications
* [New Relic University](https://learn.newrelic.com/): A range of online training for New Relic users of every level
* [New Relic Technical Support](https://support.newrelic.com/) 24/7/365 ticketed support. Read more about our [Technical Support Offerings](https://docs.newrelic.com/docs/licenses/license-information/general-usage-licenses/support-plan).

## Privacy

At New Relic we take your privacy and the security of your information seriously, and are committed to protecting your information. We must emphasize the importance of not sharing personal data in public forums, and ask all users to scrub logs and diagnostic information for sensitive information, whether personal, proprietary, or otherwise.

We define “Personal Data” as any information relating to an identified or identifiable individual, including, for example, your name, phone number, post code or zip code, Device ID, IP address, and email address.

For more information, review [New Relic’s General Data Privacy Notice](https://newrelic.com/termsandconditions/privacy).

## Contribute

We encourage your contributions to improve this project! Keep in mind that when you submit your pull request, you'll need to sign the CLA via the click-through using CLA-Assistant. You only have to sign the CLA one time per project.

If you have any questions, or to execute our corporate CLA (which is required if your contribution is on behalf of a company), drop us an email at opensource@newrelic.com.

**A note about vulnerabilities**

As noted in our [security policy](../../security/policy), New Relic is committed to the privacy and security of our customers and their data. We believe that providing coordinated disclosure by security researchers and engaging with the security community are important means to achieve our security goals.

If you believe you have found a security vulnerability in this project or any of New Relic's products or websites, we welcome and greatly appreciate you reporting it to New Relic through [HackerOne](https://hackerone.com/newrelic).

If you would like to contribute to this project, review [these guidelines](./CONTRIBUTING.md).

To [all contributors](https://github.com/newrelic/nri-vsphere/graphs/contributors), we thank you!  Without your contribution, this project would not be what it is today.  We also host a community project page dedicated to 
the [New Relic's VMware vSphere integration](https://opensource.newrelic.com/projects/newrelic/nri-vsphere).

## License

nri-vsphere is licensed under the [Apache 2.0](http://apache.org/licenses/LICENSE-2.0.txt) License.

New Relic Infrastructure Integration for VMware vSphere also uses source code from third-party libraries. You can find full details on which libraries are used and the terms under which they are licensed in the [third-party notices document](https://github.com/newrelic/nri-vsphere/blob/master/THIRD_PARTY_NOTICES.md).
