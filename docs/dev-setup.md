# Development Environment Setup for HWCC

## Prerequisites

* metal3
* go version v1.13+.
* docker version 17.03+.
* kubectl version v1.11.3+.
* kustomize v3.1.0+
* Access to a Kubernetes v1.11.3+ cluster.

Please follow metal3 dev guide for setting up above prerequisites -
<https://github.com/metal3-io/metal3-dev-env/blob/master/README.md>

## Running Controller within Pod

1. Clone Hardware Classification Controller

    ```bash
   git clone https://github.com/metal3-io/hardware-classification-controller
   cd hardware-classification-controller/
   ```

2. Deploy Hardware Classification Controller to the cluster with image
specified.

    ```bash
   make deploy IMG=quay.io/metal3-io/hardware-classification-controller:latest
   ```

3. Verify successfull deployment

    ```bash
   kubectl get pods -A
   ```

## Running Controller locally

1. Clone Hardware Classification Controller

    ```bash
   git clone https://github.com/metal3-io/hardware-classification-controller
   cd hardware-classification-controller/
   ```

2. Install the CRDs

    ```bash
   make install
   ```

3. Run controller (this will run in the foreground, so switch to a new
terminal if you want controller running).

    ```bash
   make run
   ```

Note: Setup is completed here. To use HWCC follow HWCC user-guide
(Link will be shared soon)
