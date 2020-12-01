package ironic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
)

type InspectorCallbackData struct {
	RootDisk      RootDisk    `json:"root_disk"`
	BootInterface string      `json:"boot_interface"`
	Inventory     Inventory   `json:"inventory"`
	Interfaces    []Interface `json:"interfaces"`
	Disks         []Disk      `json:"disks"`
	Memory        Memory      `json:"memory"`
	CPU           CPU         `json:"cpu"`
	Logs          string      `json:"logs"`
}

type RootDisk struct {
	Rotational bool   `json:"rotational"`
	Vendor     string `json:"vendor"`
	Name       string `json:"name"`
	Model      string `json:"model"`
	Serial     string `json:"serial"`
	Size       int    `json:"size"`
}

type Inventory struct {
	BmcAddress string `json:"bmc_address"`
}

type Interface struct {
	Lldp       *string `json:"lldp"`
	Product    string  `json:"product"`
	Vendor     *string `json:"vendor"`
	Name       string  `json:"name"`
	HasCarrier bool    `json:"has_carrier"`
	IP4Address string  `json:"ipv4_address"`
	ClientID   *string `json:"client_id"`
	MacAddress string  `json:"mac_address"`
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
	Size               int     `json:"size"`
}

type Memory struct {
	PhysicalMb int `json:"physical_mb"`
	Total      int `json:"total"`
}

type CPU struct {
	Count        int      `json:"count"`
	Frequency    string   `json:"frequency"`
	Flags        []string `json:"flags"`
	Architecture string   `json:"architecture"`
}

func CreateNodeWithInspector(d *InspectorCallbackData, host string) (err error) {
	client := &http.Client{}
	u, err := url.Parse(fmt.Sprintf("https://%s", host))
	if err != nil {
		return
	}
	u.Path = path.Join(u.Path, "/v1/continue")
	db, err := json.Marshal(d)
	if err != nil {
		return
	}
	req, err := http.NewRequest(http.MethodPatch, u.String(), bytes.NewBuffer(db))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	_, err = client.Do(req)
	if err != nil {
		return
	}
	return
}
