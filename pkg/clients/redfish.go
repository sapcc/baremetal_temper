package clients

import (
	"fmt"
	"net"
	"strings"

	"github.com/sapcc/ironic_temper/pkg/config"
	"github.com/sapcc/ironic_temper/pkg/diagnostics"
	"github.com/sapcc/ironic_temper/pkg/model"
	log "github.com/sirupsen/logrus"
	"github.com/stmcginnis/gofish"
	"github.com/stmcginnis/gofish/redfish"
)

type RedfishClient struct {
	gCfg    gofish.ClientConfig
	client  *gofish.APIClient
	service *gofish.Service
	data    *model.InspectonData
	node    *model.IronicNode
	log     *log.Entry
}

//NewRedfishClient creates redfish client
func NewRedfishClient(cfg config.Config, host string, ctxLogger *log.Entry) *RedfishClient {
	return &RedfishClient{
		gCfg: gofish.ClientConfig{
			Endpoint:  fmt.Sprintf("https://%s", host),
			Username:  cfg.Redfish.User,
			Password:  cfg.Redfish.Password,
			Insecure:  true,
			BasicAuth: false,
		},
		log: ctxLogger,
	}
}

//LoadInventory loads the node's inventory via it's redfish api
func (r RedfishClient) LoadInventory(n *model.IronicNode) (err error) {
	r.log.Debug("calling redfish api to load node info")
	client, err := gofish.Connect(r.gCfg)
	if err != nil {
		return
	}
	r.node = n
	defer client.Logout()
	r.client = client
	r.data = &model.InspectonData{}
	r.service = client.Service
	if err = r.setBMCAddress(); err != nil {
		return
	}
	if err = r.setInventory(); err != nil {
		return
	}
	n.InspectionData = *r.data
	return
}

func (r RedfishClient) setBMCAddress() (err error) {
	m, err := r.service.Managers()
	if err != nil && len(m) == 0 {
		return fmt.Errorf("cannot set bmc address")
	}
	in, err := m[0].EthernetInterfaces()
	if err != nil || len(in) == 0 {
		return
	}
	addr, err := net.LookupAddr(in[0].IPv4Addresses[0].Address)
	if err != nil || len(addr) == 0 {
		return
	}
	if r.node.Host == addr[0] {
		r.data.Inventory.BmcAddress = addr[0]
		return
	}
	return fmt.Errorf("dns record %s does not map to ip: %s", addr[0], in[0].IPv4Addresses[0].Address)
}

func (r RedfishClient) setInventory() (err error) {
	ch, err := r.service.Chassis()
	if err != nil || len(ch) == 0 {
		return
	}

	r.data.Inventory.SystemVendor.Manufacturer = ch[0].Manufacturer
	r.data.Inventory.SystemVendor.SerialNumber = ch[0].SerialNumber
	r.data.Inventory.SystemVendor.ProductName = ch[0].Model

	s, err := r.service.Systems()
	if err != nil || len(s) == 0 {
		return
	}
	if err = r.setMemory(s[0]); err != nil {
		return
	}
	if err = r.setDisks(s[0]); err != nil {
		return
	}
	if err = r.setCPUs(s[0]); err != nil {
		return
	}
	if err = r.setNetworkDevicesData(s[0]); err != nil {
		return
	}
	return
}

func (r RedfishClient) setMemory(s *redfish.ComputerSystem) (err error) {
	mem, err := s.Memory()
	if err != nil {
		return
	}
	r.data.Inventory.Memory.PhysicalMb = calcTotalMemory(mem)
	return
}

func (r RedfishClient) setDisks(s *redfish.ComputerSystem) (err error) {
	st, err := s.Storage()
	rootDisk := model.RootDisk{
		Rotational: true,
	}
	r.data.Inventory.Disks = make([]model.Disk, 0)
	for _, s := range st {
		ds, err := s.Drives()
		if err != nil {
			continue
		}
		for _, s := range ds {
			rotational := true
			if s.RotationSpeedRPM == 0 {
				rotational = false
			}
			disk := model.Disk{
				Name:   s.Name,
				Model:  s.Model,
				Vendor: s.Manufacturer,
				//inspector converts bytes to gibibyte
				Size:       int64(float64(s.CapacityBytes) * 1.074),
				Rotational: rotational,
			}

			if s.CapacityBytes > rootDisk.Size {
				rootDisk.Size = int64(float64(s.CapacityBytes) * 1.074)
				rootDisk.Name = s.Name
				rootDisk.Model = s.Model
				rootDisk.Vendor = s.Manufacturer
				if s.RotationSpeedRPM == 0 {
					rootDisk.Rotational = rotational
				}
			}
			r.data.Inventory.Disks = append(r.data.Inventory.Disks, disk)
		}
	}

	r.data.RootDisk = rootDisk
	return
}

func (r RedfishClient) setCPUs(s *redfish.ComputerSystem) (err error) {
	cpu, err := s.Processors()
	if err != nil || len(cpu) == 0 {
		return
	}
	r.data.Inventory.CPU.Count = s.ProcessorSummary.LogicalProcessorCount / s.ProcessorSummary.Count
	r.data.Inventory.CPU.Architecture = strings.Replace(string(cpu[0].InstructionSet), "-", "_", 1)
	return
}

func (r RedfishClient) setNetworkDevicesData(s *redfish.ComputerSystem) (err error) {
	ethInt, err := s.EthernetInterfaces()
	if err != nil || len(ethInt) == 0 {
		return
	}
	r.data.Inventory.Boot.PxeInterface = ethInt[0].MACAddress
	r.data.BootInterface = "01-" + strings.ReplaceAll(ethInt[0].MACAddress, ":", "-")
	r.data.Inventory.Boot.CurrentBootMode = "bios"
	r.data.Inventory.Interfaces = make([]model.Interface, len(ethInt))
	for i, e := range ethInt {
		r.data.Inventory.Interfaces[i].MacAddress = e.MACAddress
		r.data.Inventory.Interfaces[i].Name = e.ID
		//r.data.Inventory.Interfaces[i].HasCarrier = true
	}

	return
}

func (r RedfishClient) RunRemoteDiagnostics(n *model.IronicNode) (err error) {
	r.log.Debug("calling redfish api to run remote diagnostics")
	client, err := gofish.Connect(r.gCfg)
	if err != nil {
		return
	}
	defer client.Logout()
	d, err := diagnostics.GetRemoteDiagnostics(client, r.log)
	if err != nil {
		return
	}

	if d == nil {
		return
	}

	return d.Run()
}
