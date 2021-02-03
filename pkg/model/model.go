package model

type IronicNode struct {
	Name           string
	IP             string
	UUID           string `json:"uuid"`
	InstanceUUID   string
	InstanceIPv4   string
	Host           string
	ResourceClass  string
	InspectionData InspectonData
	Connections    map[string]string
}

type InspectonData struct {
	RootDisk      RootDisk  `json:"root_disk"`
	BootInterface string    `json:"boot_interface"`
	Inventory     Inventory `json:"inventory"`

	Logs string `json:"logs"`
}

type RootDisk struct {
	Rotational bool   `json:"rotational"`
	Vendor     string `json:"vendor"`
	Name       string `json:"name"`
	Model      string `json:"model"`
	Serial     string `json:"serial"`
	Size       int64  `json:"size"`
}

type Inventory struct {
	BmcAddress   string       `json:"bmc_address"`
	SystemVendor SystemVendor `json:"system_vendor"`
	Boot         Boot         `json:"boot"`
	Interfaces   []Interface  `json:"interfaces"`
	Disks        []Disk       `json:"disks"`
	Memory       Memory       `json:"memory"`
	CPU          CPU          `json:"cpu"`
}

type Interface struct {
	Lldp       map[string]string `json:"lldp"`
	Product    string            `json:"product"`
	Vendor     *string           `json:"vendor"`
	Name       string            `json:"name"`
	HasCarrier bool              `json:"has_carrier"`
	IP4Address string            `json:"ipv4_address"`
	ClientID   *string           `json:"client_id"`
	MacAddress string            `json:"mac_address"`
}

type Boot struct {
	CurrentBootMode string `json:"current_boot_mode"`
	PxeInterface    string `json:"pxe_interface"`
}

type SystemVendor struct {
	SerialNumber string `json:"serial_number"`
	ProductName  string `json:"product_name"`
	Manufacturer string `json:"manufacturer"`
}

type Disk struct {
	Rotational         bool    `json:"rotational"`
	Vendor             string  `json:"vendor"`
	Name               string  `json:"name"`
	Hctl               *string `json:"hctl"`
	WwnVendorExtension *string `json:"wwn_vendor_extension"`
	WwnWithExtension   *string `json:"wwn_with_extension"`
	Model              string  `json:"model"`
	Wwn                *string `json:"wwn"`
	Serial             *string `json:"serial"`
	Size               int64   `json:"size"`
}

type Memory struct {
	PhysicalMb int     `json:"physical_mb"`
	Total      float32 `json:"total"`
}

type CPU struct {
	Count        int      `json:"count"`
	Frequency    string   `json:"frequency"`
	Flags        []string `json:"flags"`
	Architecture string   `json:"architecture"`
}
