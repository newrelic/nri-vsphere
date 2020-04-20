# New Relic Infrastructure Integration for VMware vSphere

New Relic Infrastructure Integration for VMware vShpere captures critical summary
metrics and inventory data connecting to the vCenter or to an Esxi Host.

Data is obtained making use of a client provided by the `govmomi` package, the official go wrapper for VMware API.

The very same package provides as well a simulator for the virtual center leveraged to run integration tests.

The output of the integration is shaped by `newrelic/infra-integrations-sdk` package and provided in json format. The agent collects such data and ultimately sends it to the New Relic infrastructure.

See our [documentation web site](https://docs.newrelic.com/docs/integrations/host-integrations/host-integrations-list/vmware-vsphere-monitoring-integration) for more details.

## Installation

It is required to have the agent (v1.8.0+) installed 
(see
[agent installation](https://docs.newrelic.com/docs/infrastructure/new-relic-infrastructure/installation/install-infrastructure-linux)).

The preferable way to perform the installation fot this integration is to follow the official [documentation web site](https://docs.newrelic.com/docs/integrations/host-integrations/host-integrations-list/vmware-vsphere-monitoring-integration) and make use of the packet managers.

Tt is also possible, but not recommended, to perform a manual installation.
In the Linux environment, for example, in order to install the vSphere Integration it is required to configure
`./config/nri-vmware-config.yml` file and place it under 
`/etc/newrelic-infra/integrations.d/vsphere-config.yml`. 

The executable built running `make compile-linux` should be placed under: `/var/db/newrelic-infra/newrelic-integrations/bin/nri-vmware-vsphere` 

Once the file are placed in the correct folder and configured you can restart the agent:
`sudo service newrelic-infra restart`

## Integration development usage

Assuming that you have the source code and Go tool installed you can build and run the vSphere Integration locally.
* After cloning this repository, go to the directory of the vSphere Integration and build it
```bash
$ make test
$ make compile-all
```
* The command above will execute the tests for the vSphere Integration and build executable files called `nri-vmware-vsphere` under `bin/{architecture}` directory. 

Run the executable passing the required arguments:
```bash
$ ./bin/darwin/nri-vmware-vsphere --url 127.0.0.1:8989/sdk --user user --pass pass --validate_ssl false
```
* If you want to know more about usage of `./bin/darwin/nri-vmware-vsphere` check
```bash
$ ./bin/darwin/nri-vmware-vsphere -help
```

For managing external dependencies go modules are used. It is required to lock all external dependencies to specific version (if possible) into vendor directory.

## Contributing Code

We welcome code contributions (in the form of pull requests) from our user
community. Before submitting a pull request please review [these guidelines](https://github.com/newrelic/nri-vmware-vsphere/blob/master/CONTRIBUTING.md).

Following these helps us efficiently review and incorporate your contribution
and avoid breaking your code with future changes to the agent.

## Custom Integrations

To extend your monitoring solution with custom metrics, we offer the Integrations
Golang SDK which can be found on [github](https://github.com/newrelic/infra-integrations-sdk).

Refer to [our docs site](https://docs.newrelic.com/docs/infrastructure/integrations-sdk/get-started/intro-infrastructure-integrations-sdk)
to get help on how to build your custom integrations.

## Support

You can find more detailed documentation [on our website](http://newrelic.com/docs),
and specifically in the [Infrastructure category](https://docs.newrelic.com/docs/infrastructure).

If you can't find what you're looking for there, reach out to us on our [support
site](http://support.newrelic.com/) or our [community forum](http://forum.newrelic.com)
and we'll be happy to help you.

Find a bug? Contact us via [support.newrelic.com](http://support.newrelic.com/),
or email support@newrelic.com.

New Relic, Inc.
