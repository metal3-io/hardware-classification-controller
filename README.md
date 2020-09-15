# Hardware classification controller

Hardware Classification Controller (HWCC) goal is to identify right matched
host for user provided hardware configurations.

User provides workload profile which has set of hardware configuration
parameters (CPU, RAM, DISK and NIC). HWCC fetches the hosts from BMO.
It then filters the ready state hosts and compares those with the expected
configurations provided by user in profile. The hosts matched to any of the
configurations are then labelled accordingly.

More capabilities are being added regularly. See open issues and pull
requests for more information on work in progress.

For more information about Metal³, the Hardware Classification, and other
related components, see the
[Metal³ docs](https://github.com/metal3-io/metal3-docs)and
[Kubebuilder Book](https://book.kubebuilder.io/quick-start.html#create-a-project).

## Resources

* [API documentation](docs/api.md)
* [Setup Development Environment](docs/dev-setup.md)
* [User guide](docs/user-guide.md)
