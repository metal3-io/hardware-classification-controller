package classifier

import (
	bmh "github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"

	hwcc "github.com/metal3-io/hardware-classification-controller/api/v1alpha1"
)

func checkCPU(profile *hwcc.HardwareClassification, host *bmh.BareMetalHost) bool {
	cpuDetails := profile.Spec.HardwareCharacteristics.Cpu
	if cpuDetails == nil {
		return true
	}

	ok := checkRangeInt(
		cpuDetails.MinimumCount,
		cpuDetails.MaximumCount,
		host.Status.HardwareDetails.CPU.Count)
	log.Info("CPU",
		"host", host.Name,
		"profile", profile.Name,
		"namespace", host.Namespace,
		"minCount", cpuDetails.MinimumCount,
		"maxCount", cpuDetails.MaximumCount,
		"actualCount", host.Status.HardwareDetails.CPU.Count,
		"ok", ok,
	)
	if !ok {
		return false
	}

	ok = checkRangeClockSpeed(
		bmh.ClockSpeed(cpuDetails.MinimumSpeedMHz),
		bmh.ClockSpeed(cpuDetails.MaximumSpeedMHz),
		host.Status.HardwareDetails.CPU.ClockMegahertz)
	log.Info("CPU",
		"host", host.Name,
		"profile", profile.Name,
		"namespace", host.Namespace,
		"minSpeed", cpuDetails.MinimumSpeedMHz,
		"maxSpeed", cpuDetails.MaximumSpeedMHz,
		"actualSpeed", host.Status.HardwareDetails.CPU.ClockMegahertz,
		"ok", ok,
	)
	if !ok {
		return false
	}

	return true
}

func checkRangeClockSpeed(min, max, count bmh.ClockSpeed) bool {
	if min > 0 && count < min {
		return false
	}
	if max > 0 && count > max {
		return false
	}
	return true
}
