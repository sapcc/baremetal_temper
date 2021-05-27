package node

import (
	"sort"
	"sync"

	"github.com/netbox-community/go-netbox/netbox/models"
	"github.com/sapcc/baremetal_temper/pkg/clients"
	"github.com/sapcc/baremetal_temper/pkg/config"
	"github.com/stmcginnis/gofish/redfish"

	log "github.com/sirupsen/logrus"
)

type Node struct {
	Name         string `json:"name"`
	RemoteIP     string
	PrimaryIP    string
	UUID         string `json:"uuid"`
	InstanceUUID string
	InstanceIPv4 string
	Host         string
	Tasks        map[int]*Task
	Status       string `json:"status"`
	Clients      ApiClients

	ResourceClass  string
	InspectionData InspectonData
	Interfaces     map[string]NodeInterface
	IpamAddresses  []models.IPAddress

	log *log.Entry
	cfg config.Config
	oc  *clients.OpenstackClient
}

type Task struct {
	Exec     func() error
	Name     string
	Priority int
	Error    error
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
		Name:  name,
		cfg:   cfg,
		Tasks: make(map[int]*Task, 0),
		log:   ctxLogger,
		oc:    clients.NewClient(cfg, ctxLogger),
	}
	if cfg.Redfish.User != "" {
		n.Clients.Redfish = clients.NewRedfish(cfg, ctxLogger)
	}
	if cfg.Netbox.Token != "" {
		n.Clients.Netbox, err = clients.NewNetbox(cfg, ctxLogger)
		if err != nil {
			return
		}
	}
	return
}

func (n *Node) Temper(netboxSts bool, wg *sync.WaitGroup) {
	defer wg.Done()
	prios := make([]int, 0, len(n.Tasks))
	for k := range n.Tasks {
		prios = append(prios, k)
	}
	sort.Ints(prios)
	if err := n.loadInfos(); err != nil {
		n.Status = "failed"
		log.Errorf("failed to load %s info. err: %s", n.Name, err.Error())
		return
	}
	for _, k := range prios {
		if err := n.Tasks[k].Exec(); err != nil {
			if _, ok := err.(*AlreadyExists); ok {
				log.Infof("Node %s already exists, nothing to temper", n.Name)
				break
			}
			n.Tasks[k].Error = err
			n.Status = "failed"
		}
	}
	n.cleanupHandler(netboxSts)
	return
}

func (n *Node) cleanupHandler(netboxSts bool) {
	for _, t := range n.Tasks {
		if t.Error != nil {
			log.Errorf("error tempering node %s. task: %s err: %s", n.Name, t.Name, t.Error.Error())
		}
	}
	if n.InstanceUUID != "" {
		if err := n.DeleteTestInstance(); err != nil {
			log.Error("cannot delete compute instance %s. err: %s", n.InstanceUUID, err.Error())
		}
	}
	if err := n.DeleteNode(); err != nil {
		log.Errorf("cannot delete node %s. err: %s", n.Name, err.Error())
	}
	if netboxSts {
		if err := n.SetStatus(); err != nil {
			log.Errorf("cannot set node %s status in netbox. err: %s", n.Name, err.Error())
		}
	}

}

func (n *Node) loadInfos() (err error) {
	if err = n.LoadIpamAddresses(); err != nil {
		return
	}
	n.Clients.Redfish.SetEndpoint(n.RemoteIP)
	if err = n.LoadInventory(); err != nil {
		return
	}
	if err = n.LoadInterfaces(); err != nil {
		return
	}
	return
}
