# API and Resource Definitions

## HardwareClassificationController

**Metal³** introduces the concept of **HardwareClassificationController** resource, which defines expected hardware configurations to classify BareMetalHosts. The HardwareClassificationController embeds two well differentiated sections, the hardware classification controller specification and its current status.


### HardwareClassificationController metadata

* name -- name of profile
* namespace -- namespace from which BareMetalHosts to be fetched
* labels -- Label is a key-value pair where key should be 'profile-name' and value can be anything. If not provided by user in yaml **default** label is `hardwareclassification.metal3.io/<profile-name> : matches`. This label is set on BaremetalHosts matching to expected hardware configurations provided by user in YAML.

      To check labels assigned on BaremetalHosts:

      $ kubectl get bmh -n <namespace> --show-labels

### HardwareClassificationController spec

The *HardwareClassificationController's* *spec* contains mainly expected hardware configuration details.

#### Spec fields

* *hardwareCharacteristics* -- HardwareCharacteristics defines expected hardware configurations for CPU, DISK, NIC and RAM.
  * *cpu* -- Expected CPU configurations:
    * minimumCount -- minimum cpu count
    * maximumCount -- maximum cpu count
    * minimumSpeedMHz -- minimum speed in MHz
    * maximumSpeedMHz -- maximum speed in MHz
  * *disk* -- Expected DISK configurations:
    * minimumCount -- minimum disk count
    * maximumCount -- maximum disk count
    * minimumIndividualSizeGB -- minimum individual disk size in GB
    * maximumIndividualSizeGB -- maximum individual disk size in GB
  * *ram* -- Expected RAM configurations:
    * minimumSizeGB -- minimum ram size in GB
    * maximumSizeGB -- maximum ram size in GB
  * *nic* -- Expected NIC configurations:
    * minimumCount -- minimum nic count
    * maximumCount -- maximum nic count


### HardwareClassificationController status

The *HardwareClassificationController's* *status* which represents the observed state of HardwareClassification. 

#### Status fields

* *errorType* -- errorType indicates the type of failure encountered
  * LabelUpdateFailure -- LabelUpdateFailure is an error condition occurring when the controller is unable to update label of BareMetalHost.
  * LabelDeleteFailure -- LabelDeleteFailure is an error condition occurring when the controller is unable to delete label of BareMetalHost.
  * FetchBMHListFailure -- FetchBMHListFailure is an error condition occurring when the controller is unable to fetch BareMetalHost from BMO.
  * ProfileMisConfigured -- ProfileMisConfigured is an error condition occurring when the extracted profile is misconfigured.

* *profileMatchStatus* -- profileMatchStatus indicates whether expected hardwareCharacteristics matches to any of BareMetalHost or not.
  * ProfileMatchStatusEmpty -- default is empty
  * ProfileMatchStatusMatched -- profileMatchStatusMatched is the status value when the profile matches to one of the BareMetalHost.
  * ProfileMatchStatusUnMatched -- profileMatchStatusUnMatched is the status value when the profile does not matches to any of the BareMetalHost.

* *errorMessage* -- Details of the last error reported by the hardwareclassification system.

### HardwareClassificationController Example

The following is a sample CRD of a HardwareClassificationController resource (in YAML), it includes its metadata and specification section.

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
