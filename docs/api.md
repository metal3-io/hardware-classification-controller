# API and Resource Definitions

## HardwareClassificationController

**Metal³** introduces the concept of **HardwareClassificationController** resource, which
defines expected hardware configurations. The **HardwareClassificationController** embeds
two well differentiated sections, the hardware classification controller specification
and its current status.

### HardwareClassificationController spec

The *HardwareClassificationController's* *spec* defines the desire state of the HardwareClassificationController. It contains mainly expected hardware configuration details.

#### Spec fields

* *namespace* -- The namespace name to fetch the BareMetalHost under metal3 namespace.

* *expectedValidationConfiguration* -- A list of multiple profiles containing expected minimum
hardware configuration for CPU, Storage, RAM, NICS, Firmware, SystemVendor.
  * *minimumCPU* -- Defined minimum CPU configuration. It is mandatory configuration.
    * *count* -- Count of CPU's.
  * *minimumDisk* -- Defines minimum configuration for storage. It is mandatory configuration.
    * *sizeBytesGB* -- size of disk in GB.
    * *numberOfDisks* -- number of Disks.
  * *minimumNICS* -- Defines minimum Configuration of NICS. It is mandatory configuration.
    * *numberOfNICS* -- Count of NICS.
  * *minimumRAM* -- Defines minimum RAM size in GB. It is optional configuration.
  * *systemVendor* -- Name of System vendor.
  * *firmware* -- Defines firmware Configuration. It is optional configuration.
    * *version* -- Version for any of RAID, BaseBandManagement, BIOS, IDRAC.


### HardwareClassificationController status

Moving onto the next block, the *HardwareClassificationController's* *status* which represents
the current state.

#### Status fields

* *errorMessage* -- Details of the last error if any

### HardwareClassificationController Example

The following is a complete example from a running cluster of a *HardwareClassificationController*
resource (in YAML), it includes its specification section. In expectedValidationConfiguration multipl

```yaml
apiVersion: metal3.io.sigs.k8s.io/v1alpha1
kind: HardwareClassificationController
metadata:
  name: hardwareclassificationcontroller-sample
spec:
  namespace: metal3
  expectedValidationConfiguration:
        - profileName: Profile 1
          minimumCPU:
                count: 4
          minimumDisk:
                sizeBytesGB: 25
                numberOfDisks: 2
          minimumNICS:
                numberOfNICS: 4
          minimumRAM: 8
          systemVendor:
                name: "Dell Inc"
          firmware:
                version:
                    RAID: "25.5.3.0005"
```