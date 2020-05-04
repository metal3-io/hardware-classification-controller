# Hardware classification controller
Controller for matching host hardware characteristics to expected values.

The HWCC (Hardware Classification Controller) implements Kubernetes API for labeling the valid hosts.
Implemented `hardware-classification` CRD expects the Profiles to be validated as yaml input.

Comparision and validation is done on baremetalhost list provided `BMO` against hardware profile mentioned in metal3.io_v1alpha1_hardwareclassification.yaml.

More capabilities are being added regularly. See open issues and pull requests for more information on work in progress.

For more information about Metal³, the Hardware Classification, and other related components, see the [Metal³ docs](https://github.com/metal3-io/metal3-docs).

## Setup Development Environment

### Prerequisites
* metal3
* go version v1.13+.
* docker version 17.03+.
* kubectl version v1.11.3+.
* kustomize v3.1.0+
* kubebuilder v2.2.0+
* Access to a Kubernetes v1.11.3+ cluster.

### Install metal3 dev-env

- Refer [Metal³ dev env setup](https://github.com/metal3-io/metal3-dev-env/blob/master/README.md) for metal3 installation.

### Install hardware-classification-controller

- After successfully installation of metal3-dev-env, under /go/src directory pull [HWCC](https://github.com/metal3-io/hardware-classification-controller.git).

- Go under directory ./hardware-classification-controller.

    $ cd hardware-classification-controller

- Install HWCC

    $ make install

- Run HWCC

    $ make run

    Now keep the controller running.

- Apply CR

    Check kubectl is in PATH. Open new terminal and apply CR.

    $ kubectl apply -f metal3.io_v1alpha1_hardwareclassification.yaml

## Example Usage

User can validate and classify the hosts based on hardware requirement. User will get to know how many hosts matched to user profile.
User can select any of matched host and go for provisioning.

This feature helps user to avoid runtime failures and also increases the performances for workload deployments.

## Resources

* API documentation
* Setup Development Environment
* Configuration
