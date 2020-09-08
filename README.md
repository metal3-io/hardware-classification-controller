# Hardware classification controller

Controller for matching expected hardware characteristics to BareMetalHost's
hardware details.

The HWCC (Hardware Classification Controller) implements Kubernetes API for
labeling the matching hosts. Implemented `hardware-classification` CRD expects
profile having hardware configurations for CPU, RAM, Disk and NIC.

Comparision and validation is done on list of baremetalhosts in ready state
provided by `BMO` against hardware configuration mentioned in
metal3.io_v1alpha1_hardwareclassification.yaml.

More capabilities are being added regularly. See open issues and pull
requests for more information on work in progress.

For more information about Metal³, the Hardware Classification, and other
related components, see the [Metal³ docs](https://github.com/metal3-io/metal3-docs) and
[Kubebuilder Book](https://book.kubebuilder.io/quick-start.html#create-a-project).

## Resources

* [API documentation](docs/api.md)
* [Setup Development Environment](docs/dev-setup.md)
* [User guide](docs/user-guide.md)
