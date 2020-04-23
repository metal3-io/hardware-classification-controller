# hardware-classification
Controller for matching host hardware characteristics to expected values.

The HWC (Hardware Classification) implements Kubernetes API for labeling the valid hosts.
Implemented `hardware-classification` CRD expects the Profiles to be validated as yaml input.

Comparision and validation is done on baremetalhost list provided `BMO` against hardware profiles mentioned in metal3.io_v1alpha1_hardwareclassification.yaml.

More capabilities are being added regularly. See open issues and pull requests for more information on work in progress.

For more information about Metal³, the Hardware Classification, and other related components, see the [Metal³ docs](https://github.com/metal3-io/metal3-docs).

## Resources

* API documentation
* Setup Development Environment
* Configuration
