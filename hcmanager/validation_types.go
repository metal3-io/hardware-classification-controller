package hcmanager

//RAM contains ram details fetched from the introspection data
type RAM struct {
	RAMGb int `json:"ramMebibytes"`
}

//HardwareSystemVendor contains hardware manufacturer details fetched from the introspection data
type HardwareSystemVendor struct {
	Manufacturer string `json:"manufacturer"`
}

//NIC contains the nic details fetched from the introspection data
type NIC struct {
	Name  string `json:"name"`
	PXE   bool   `json:"pxe"`
	Count int    `json:"count"`
}

//Storage contains disk details fetched from the introspection data
type Storage struct {
	Count int    `json:"count"`
	Disk  []Disk `json:"disk"`
}

//Disk contains disk size fetched from the introspection data
type Disk struct {
	Name   string `json:"name"`
	SizeGb int64  `json:"sizeBytes"`
}

//CPU contains the clockspeed and count details fetched from the introspection data
type CPU struct {
	Count      int     `json:"count"`
	ClockSpeed float64 `json:"clockspeed"`
}
