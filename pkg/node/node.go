package node

import (
	"fmt"
	"sync"

	"github.com/netbox-community/go-netbox/netbox/models"
	"github.com/sapcc/baremetal_temper/pkg/clients"
	"github.com/sapcc/baremetal_temper/pkg/config"
	"github.com/sapcc/baremetal_temper/pkg/task"
	"github.com/stmcginnis/gofish/redfish"

	log "github.com/sirupsen/logrus"
)

type Node struct {
	Name           string                   `json:"name"`
	RemoteIP       string                   `json:"remoteIP"`
	PrimaryIP      string                   `json:"primaryIP"`
	UUID           string                   `json:"uuid"`
	ProvisionState string                   `json:"provisionState"`
	InstanceUUID   string                   `json:"instanceUUID"`
	InstanceIPv4   string                   `json:"instanceIP"`
	Host           string                   `json:"host"`
	Tasks          []*task.Task             `json:"tasks"`
	Status         string                   `json:"status"`
	Clients        ApiClients               `json:"-"`
	PortGroupUUID  string                   `json:"portGroupUUID"`
	ResourceClass  string                   `json:"-"`
	InspectionData InspectonData            `json:"-"`
	Interfaces     map[string]NodeInterface `json:"-"`
	IpamAddresses  []models.IPAddress       `json:"-"`

	DeviceConfig *models.DeviceWithConfigContext    `json:"-"`
	tasksExecs   map[string]map[string][]*task.Exec `json:"-"`

	log *log.Entry         `json:"-"`
	cfg config.Config      `json:"-"`
	oc  *clients.Openstack `json:"-"`
}

type ApiClients struct {
	Redfish *clients.Redfish
	Netbox  *clients.Netbox
}

type NodeInterface struct {
	Connection     string
	ConnectionIP   string
	Port           string
	Mac            string
	PortLinkStatus redfish.PortLinkStatus
}

type InspectonData struct {
	RootDisk      RootDisk  `json:"root_disk"`
	BootInterface string    `json:"boot_interface"`
	Inventory     Inventory `json:"inventory"`
	Logs          string    `json:"logs"`
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
	Interfaces   []Interface  `json:"interfaces"`
	Boot         Boot         `json:"boot"`
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
	Model        string
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

func New(name string, cfg config.Config) (n *Node, err error) {
	ctxLogger := log.WithFields(log.Fields{
		"node": name,
	})
	n = &Node{
		Name:       name,
		Status:     "planned",
		cfg:        cfg,
		Tasks:      make([]*task.Task, 0),
		log:        ctxLogger,
		oc:         clients.NewClient(cfg, ctxLogger),
		tasksExecs: make(map[string]map[string][]*task.Exec),
	}
	if cfg.Redfish.User != "" {
		n.Clients.Redfish = clients.NewRedfish(cfg, ctxLogger)
	}
	if cfg.Netbox.Token == "" {
		return n, fmt.Errorf("missing netbox token")
	}
	n.Clients.Netbox, err = clients.NewNetbox(cfg, ctxLogger)
	if err != nil {
		return
	}
	n.initTaskExecs()
	return
}

func (n *Node) Temper(netboxSts bool, wg *sync.WaitGroup) {
	defer func() {
		if r := recover(); r != nil {
			n.log.Errorf("aborting node temper: error  %s", r)
			n.Status = "failed"
		}
		n.cleanupHandler(netboxSts)
		wg.Done()
	}()

TasksLoop:
	for _, t := range n.Tasks {
		n.log.Infof("executing temper task: %s.%s", t.Service, t.Task)
		for i, exec := range t.Exec {
			if t.Status == "done" {
				continue
			}
			if err := exec.Fn(); err != nil {
				if _, ok := err.(*AlreadyExists); ok {
					if err := n.loadBaremetalNodeInfo(); err != nil {
						break TasksLoop
					}
					if n.ProvisionState != "enroll" {
						n.log.Infof("node %s already exists, nothing to temper", n.Name)
						break TasksLoop
					}
					n.log.Info("found existing node in enroll state. ")
				} else {
					t.Error = err.Error()
					t.Status = "failed"
					n.Status = "failed"
				}

			} else {
				if i == len(t.Exec)-1 {
					t.Status = "done"
				}
			}
		}

	}
	if n.Status != "failed" {
		n.Status = "staged"
	}
	return
}

func (n *Node) cleanupHandler(netboxSts bool) {
	n.log.Debugf("calling cleanupHandler, node status: %s", n.Status)
	for _, t := range n.Tasks {
		if t.Error != "" {
			n.log.Errorf("error tempering node %s. task: %s err: %s", n.Name, t.Task, t.Error)
		}
	}
	if n.InstanceUUID != "" {
		if err := n.DeleteTestInstance(); err != nil {
			n.log.Error("cannot delete compute instance %s. err: %s", n.InstanceUUID, err.Error())
		}
	}
	if n.Status == "failed" {
		if err := n.DeleteNode(); err != nil {
			n.log.Errorf("cannot delete node %s. err: %s", n.Name, err.Error())
		}
	}
	if netboxSts {
		if err := n.SetStatus(); err != nil {
			n.log.Errorf("cannot set node %s status in netbox. err: %s", n.Name, err.Error())
		}
	}

}

func recoverTaskExec(n *Node) {
	if r := recover(); r != nil {
		fmt.Println("recovered from ", r)
		n.Status = "failed"
	}
}

func (n *Node) loadNetboxInfos() (err error) {
	if err = n.LoadDeviceConfig(); err != nil {
		return
	}
	if err = n.loadIpamAddresses(); err != nil {
		return
	}
	if err = n.Clients.Redfish.SetEndpoint(n.RemoteIP); err != nil {
		return
	}
	return
}

func (n *Node) loadRedfishInfos() (err error) {
	if err = n.loadInventory(); err != nil {
		return
	}
	if err = n.loadInterfaces(); err != nil {
		return
	}
	return
}
