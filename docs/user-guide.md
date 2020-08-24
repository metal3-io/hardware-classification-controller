# User guide for Hardware Classification Controller

## Hardware Classification Controller

User can validate and classify the hosts based on hardware requirement.
User will get to know how many hosts matched to user profile.
User can select any of matched host and go for provisioning.

This feature helps user to avoid runtime failures.

## Prerequisites

1. HWCC setup

    Follow setup documentation guide for setup.
[Setup documentation](docs/dev-setup.md)

## HardwareClassificationController Example

The following is a sample CRD of a `HardwareClassificationController` resource
([YAML PATH](config/samples/metal3.io_v1alpha1_hardwareclassification.yaml)),
it includes its metadata and specification section.
User can create multiple profiles for different workloads as per the sample
shown below. Alternatively user can simply modify values of the parameters
shown in sample file as per application.

```yaml
apiVersion: metal3.io/v1alpha1
kind: HardwareClassification
metadata:
  name: hardwareclassification-sample
  namespace: metal3
  labels:
    hardwareclassification-sample: matches
spec:
  hardwareCharacteristics:
      cpu:
         minimumCount: 48
         maximumCount: 72
         minimumSpeedMHz: 2600
         maximumSpeedMHz: 3600
      disk:
         minimumCount: 1
         maximumCount: 8
         minimumIndividualSizeGB: 200
         maximumIndividualSizeGB: 3000
      ram:
         minimumSizeGB: 6
         maximumSizeGB: 180
      nic:
         minimumCount: 1
         maximumCount: 7
```

Note: Minimum 1 field under `hardwareCharacteristics` is mandatory.
If no field provided under `hardwareCharacteristics`, user willl get error.

e.g.

```yaml
    disk:
       minimumCount: 1
```

## Commands

User requires to use following commands for applying workload profiles
and get status for classified hosts.

### *Run*

Apply profile using `kubectl apply` command.

```yaml
    $ kubectl apply -f <path-to-hwcc.yaml>
```

### *Status*

Check labels on baremetal hosts

```yaml
    $ kubectl get bmh -n <namespace> --show-labels
```

Check status of profile by checking hardware-classification status.

```yaml
    $ kubectl get hardware-classification -n <namespace>
```

Note : Instead of hardware-classification shortform hwc or hc can be used.

### *Delete*

#### Deleting profile

To delete profile which is applied, user have two option.
    Delete using profile name:

```yaml
    $ kubectl delete -f <path-to-hwcc.yaml>
```

or delete using resource name:

```yaml
    $ kubectl delete hwc <profile-name> -n <namespace>
```

#### Deleting setup

1. Delete cluster based setup.

    To delete whole setup, delete deployment and namespace of hwcc.

    ```yaml
        $ kubectl delete deployement <deployment-name> -n <namespace>

        $ kubectl delete namespace <namespace>
    ```

1. Delete local setup.
    To delete local setup

    ```yaml
        $ make uninstall
    ```
