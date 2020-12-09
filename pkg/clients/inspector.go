package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/sapcc/ironic_temper/pkg/model"
)

type InspectorClient struct {
	Host string
}

type InspectorCallbackData struct {
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
	Lldp       *string `json:"lldp"`
	Product    string  `json:"product"`
	Vendor     *string `json:"vendor"`
	Name       string  `json:"name"`
	HasCarrier bool    `json:"has_carrier"`
	IP4Address string  `json:"ipv4_address"`
	ClientID   *string `json:"client_id"`
	MacAddress string  `json:"mac_address"`
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
	PhysicalMb float32 `json:"physical_mb"`
	Total      float32 `json:"total"`
}

type CPU struct {
	Count        int      `json:"count"`
	Frequency    string   `json:"frequency"`
	Flags        []string `json:"flags"`
	Architecture string   `json:"architecture"`
}

type InspectorErr struct {
	Error ErrorMessage `json:"error"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}

func (i InspectorClient) CreateIronicNode(d *InspectorCallbackData, in *model.IronicNode) (err error) {
	client := &http.Client{}
	u, err := url.Parse(fmt.Sprintf("http://%s", i.Host))
	if err != nil {
		return
	}
	u.Path = path.Join(u.Path, "/v1/continue")
	db, err := json.Marshal(d)
	if err != nil {
		return
	}
	fmt.Println(string(db))
	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(db))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	if res.StatusCode != http.StatusOK {
		ierr := &InspectorErr{}
		if err = json.Unmarshal(bodyBytes, ierr); err != nil {
			return fmt.Errorf("could not create node")
		}
		return fmt.Errorf(ierr.Error.Message)
	}

	if err = json.Unmarshal(bodyBytes, in); err != nil {
		return
	}
	name := strings.Split(d.Inventory.BmcAddress, ".")
	node := strings.Split(name[0], "-")
	nodeName := strings.Replace(node[0], "r", "", 1)
	in.Name = fmt.Sprintf("%s-%s", nodeName, node[1])
	return
}
