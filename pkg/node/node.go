package node

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/netbox-community/go-netbox/netbox/models"
	"github.com/sapcc/baremetal_temper/pkg/clients"
	"github.com/sapcc/baremetal_temper/pkg/config"
	"github.com/sapcc/baremetal_temper/pkg/netbox"
	_redfish "github.com/sapcc/baremetal_temper/pkg/redfish"
	log "github.com/sirupsen/logrus"
)

type Node struct {
	Name           string             `json:"name"`
	RemoteIP       string             `json:"remoteIP"`
	PrimaryIP      string             `json:"primaryIP"`
	UUID           string             `json:"uuid"`
	ProvisionState string             `json:"provisionState"`
	InstanceUUID   string             `json:"instanceUUID"`
	InstanceIPv4   string             `json:"instanceIP"`
	Host           string             `json:"host"`
	Tasks          []*netbox.Task     `json:"tasks"`
	Status         string             `json:"status"`
	PortGroupUUID  string             `json:"portGroupUUID"`
	ResourceClass  string             `json:"-"`
	IpamAddresses  []models.IPAddress `json:"-"`

	tasksExecs map[string]map[string][]*netbox.Exec `json:"-"`
	Updated    time.Time                            `json:"-"`

	Redfish _redfish.Redfish `json:"-"`
	Netbox  *netbox.Netbox   `json:"-"`

	log *log.Entry         `json:"-"`
	cfg config.Config      `json:"-"`
	oc  *clients.Openstack `json:"-"`
}

func New(name string, cfg config.Config) (n *Node, err error) {
	ctxLogger := log.WithFields(log.Fields{
		"node": name,
	})
	if len(strings.Split(name, "-")) != 2 {
		return n, fmt.Errorf("wrong node name format. e.g. node001-ap001")
	}
	n = &Node{
		Name:       name,
		Status:     "progress",
		cfg:        cfg,
		Tasks:      make([]*netbox.Task, 0),
		log:        ctxLogger,
		oc:         clients.NewClient(cfg, ctxLogger),
		tasksExecs: make(map[string]map[string][]*netbox.Exec),
	}
	if cfg.Netbox.Token == "" {
		return n, fmt.Errorf("missing netbox token")
	}
	n.Netbox, err = netbox.New(n.Name, cfg, ctxLogger)
	if err != nil {
		return n, fmt.Errorf("cannot create netbox client: %s", err.Error())
	}
	if err = n.createRedfishClient(); err != nil {
		err = fmt.Errorf("cannot create redfish client: %s", err.Error())
		n.Status = "failed"
		return
	}
	n.initTaskExecs()
	return
}

func (n *Node) Setup() (err error) {
	if err = n.Redfish.Power(false, false); err != nil {
		n.Status = "failed"
		err = fmt.Errorf("cannot power on node: %s", err.Error())
		return
	}
	if err = n.Redfish.WaitPowerStateOn(); err != nil {
		n.Status = "failed"
		err = fmt.Errorf("node does not power on: %s", err.Error())
		return
	}
	return n.mergeInterfaces()
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
	if err := n.Setup(); err != nil {
		n.log.Error(err)
		n.Status = "failed"
		return
	}
TasksLoop:
	for _, t := range n.Tasks {
		for i, exec := range t.Exec {
			if t.Status == "success" || t.Status == "done" {
				continue
			}
			n.log.Infof("executing temper task: %s", exec.Name)
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
					t.Status = "success"
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
		if err := n.Netbox.SetStatus(n.Status); err != nil {
			n.log.Errorf("cannot set node %s status in netbox. err: %s", n.Name, err.Error())
		}
	}

}

func recoverTaskExec(n *Node) {
	if r := recover(); r != nil {
		n.log.Error("recovered from ", r)
		n.Status = "failed"
	}
}

func (n *Node) createRedfishClient() (err error) {
	d, err := n.Netbox.GetData()
	if err != nil {
		return
	}

	lenovo := regexp.MustCompile(`(?i)SR950|SR650|SR850P`)
	dell := regexp.MustCompile(`(?i)R640|R730|R740|R840`)
	hpe := regexp.MustCompile(`(?i)DL560|DL360`)

	switch {
	case lenovo.MatchString(*d.Device.DeviceType.Model):
		n.log.Info("loading LENOVO redfish client")
		n.Redfish, err = _redfish.NewLenovo(d.RemoteIP, n.cfg, n.log)
	case dell.MatchString(*d.Device.DeviceType.Model):
		n.log.Info("loading DELL redfish client")
		n.Redfish, err = _redfish.NewDell(d.RemoteIP, n.cfg, n.log)
	case hpe.MatchString(*d.Device.DeviceType.Model):
		n.log.Info("loading HPE redfish client")
		n.Redfish, err = _redfish.NewHpe(d.RemoteIP, n.cfg, n.log)
	default:
		n.log.Info("loading DEFAULT redfish client")
		n.Redfish, err = _redfish.NewDefault(d.RemoteIP, n.cfg, n.log)
	}
	return
}

func (n *Node) mergeInterfaces() (err error) {
	nd, err := n.Netbox.GetData()
	if err != nil {
		return
	}
	rd, err := n.Redfish.GetData()
	if err != nil {
		return
	}
	sort.Slice(nd.Interfaces, func(i, j int) bool {
		nic1 := nd.Interfaces[i].Nic * 10
		port1 := nd.Interfaces[i].PortNumber
		nic2 := nd.Interfaces[j].Nic * 10
		port2 := nd.Interfaces[j].PortNumber
		return nic1+port1 < nic2+port2
	})
	rd.Inventory.BmcAddress = nd.DNSName
	sort.Slice(rd.Inventory.Interfaces, func(i, j int) bool {
		//make sure nic = 0 port = x is always smaller than a nic = 1, port = x
		nic1 := rd.Inventory.Interfaces[i].Nic * 10
		port1 := rd.Inventory.Interfaces[i].Port
		nic2 := rd.Inventory.Interfaces[j].Nic * 10
		port2 := rd.Inventory.Interfaces[j].Port
		return nic1+port1 < nic2+port2
	})
	interfaces := make([]netbox.NodeInterface, 0)
	i := 0
	for _, intf := range nd.Interfaces {
		redfishIntf := rd.Inventory.Interfaces[i]
		intf.Mac = redfishIntf.MacAddress
		intf.PortLinkStatus = redfishIntf.PortLinkStatus
		if intf.Nic == 0 {
			if redfishIntf.Nic != 0 {
				continue
			}
			intf.RedfishName = fmt.Sprintf("L%d", redfishIntf.Port)
			i++
			n.log.Debugf("found interface: %s, redfish: %s, mac: %s", intf.Name, intf.RedfishName, intf.Mac)
			interfaces = append(interfaces, intf)
			continue
		}
		if redfishIntf.Nic == 0 {
			i++
			continue
		}
		intf.RedfishName = fmt.Sprintf("PCI%d-port%d", redfishIntf.Nic, redfishIntf.Port)
		i++
		n.log.Debugf("found interface: %s, redfish: %s, mac: %s", intf.Name, intf.RedfishName, intf.Mac)
		interfaces = append(interfaces, intf)
	}
	nd.Interfaces = interfaces
	return
}
