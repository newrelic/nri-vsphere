# New Relic Infrastructure Integration for VMware vSphere

New Relic Infrastructure Integration for VMware vShpere captures critical summary metrics and inventory data by connecting to the vCenter or to an ESXi Host. For more information, see [the official New Relic docs](https://docs.newrelic.com/docs/integrations/host-integrations/host-integrations-list/vmware-vsphere-monitoring-integration).

vSphere data is obtained via a client provided by the `govmomi` package, the official go wrapper for VMware API. The same package provides a simulator for the virtual center, leveraged to run integration tests. The output of the integration is shaped by the `newrelic/infra-integrations-sdk` package and provided in JSON format. The agent collects such data and sends it to the New Relic infrastructure.

## Requirements

- [New Relic Infrastructure agent version 1.8.0 or higher](https://docs.newrelic.com/docs/infrastructure/install-configure-manage-infrastructure)
- New Relic Infrastructure Pro subscription or trial

## Installation

To install the integration, follow the official [documentation](https://docs.newrelic.com/docs/integrations/host-integrations/host-integrations-list/vmware-vsphere-monitoring-integration). We recommend using a package manager.

## Development

If you have the source code and Go toolchain installed, you can build and run the vSphere integration locally.

1. After cloning this repository, go to vSphere integration directory and build the integration:
```bash
$ make test
$ make compile-all
```
The command above executes the tests for the vSphere integration and builds an executable file named `nri-vsphere` under `bin/{architecture}`. 

2. Run the executable with the following arguments:
```bash
$ ./bin/darwin/nri-vsphere --url 127.0.0.1:8989/sdk --user user --pass pass --validate_ssl false
```
To learn more about the usage of `./bin/darwin/nri-vsphere`, pass the `-help` argument.
```bash
$ ./bin/darwin/nri-vsphere -help
```

External dependencies are managed via govendor.

## Contributing code

We welcome code contributions (in the form of pull requests) from our user community. Before submitting a pull request, please review [these guidelines](https://github.com/newrelic/nri-vsphere/blob/master/CONTRIBUTING.md).

Following these helps us efficiently review and incorporate your contribution and avoid breaking your code with future changes to the agent.

## Custom integrations

To extend your monitoring solution with custom metrics, we offer the Integrations Golang SDK, which can be found on [GitHub](https://github.com/newrelic/infra-integrations-sdk).

Refer to [our docs site](https://docs.newrelic.com/docs/infrastructure/integrations-sdk/get-started/intro-infrastructure-integrations-sdk) to get help on how to build your custom integrations.

## Support

Need help? See our [troubleshooting page](troubleshooting.md). You can find more detailed documentation [on the New Relic docs site](http://newrelic.com/docs).

If you can't find what you're looking for there, reach out to us on our [support site](http://support.newrelic.com/) or our [community forum](http://forum.newrelic.com) and we'll be happy to help you.

Found a bug? Contact us at [support.newrelic.com](http://support.newrelic.com/)

### Community

New Relic hosts and moderates an online forum where customers can interact with New Relic employees as well as other customers to get help and share best practices. Like all official New Relic open source projects, there's a related Community topic in the New Relic Explorers Hub. You can find this project's topic/threads here:

https://discuss.newrelic.com/c/support-products-agents/new-relic-infrastructure

### Issues / Enhancement Requests

Issues and enhancement requests can be submitted in the [Issues tab of this repository](../../issues). Please search for and review the existing open issues before submitting a new issue.

## License

The project is released under version 2.0 of the [Apache license](http://www.apache.org/licenses/LICENSE-2.0).
