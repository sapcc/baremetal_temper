package clients

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/sapcc/baremetal_temper/pkg/config"
	"github.com/sapcc/baremetal_temper/pkg/model"
	log "github.com/sirupsen/logrus"
	"github.com/stmcginnis/gofish"
	"github.com/stmcginnis/gofish/redfish"
)

type RedfishClient struct {
	ClientConfig *gofish.ClientConfig
	client       *gofish.APIClient
	service      *gofish.Service
	node         *model.Node
	log          *log.Entry
}

//NewRedfishClient creates redfish client
func NewRedfishClient(cfg config.Config, ctxLogger *log.Entry) *RedfishClient {
	return &RedfishClient{
		ClientConfig: &gofish.ClientConfig{
			Endpoint:  fmt.Sprintf("https://%s", "dummy.net"),
			Username:  cfg.Redfish.User,
			Password:  cfg.Redfish.Password,
			Insecure:  true,
			BasicAuth: false,
		},
		log: ctxLogger,
	}
}

//SetEndpoint sets the redfish api endpoint
func (r RedfishClient) SetEndpoint(n *model.Node) (err error) {
	r.ClientConfig.Endpoint = fmt.Sprintf("https://%s", n.RemoteIP)
	return
}

//LoadInventory loads the node's inventory via it's redfish api
func (r RedfishClient) LoadInventory(n *model.Node) (err error) {
	r.log.Debug("calling redfish api to load node info")
	client, err := gofish.Connect(*r.ClientConfig)
	if err != nil {
		return
	}
	r.node = n
	defer client.Logout()
	r.client = client
	r.service = client.Service
	if err = r.setInventory(n); err != nil {
		return
	}
	return
}

func (r RedfishClient) setInventory(n *model.Node) (err error) {
	ch, err := r.service.Chassis()
	if err != nil || len(ch) == 0 {
		return
	}

	n.InspectionData.Inventory.SystemVendor.Manufacturer = ch[0].Manufacturer
	n.InspectionData.Inventory.SystemVendor.SerialNumber = ch[0].SerialNumber

	// not performant string comparison due to toLower
	if strings.Contains(strings.ToLower(ch[0].Manufacturer), "dell") {
		n.InspectionData.Inventory.SystemVendor.SerialNumber = ch[0].SKU
	}
	n.InspectionData.Inventory.SystemVendor.ProductName = ch[0].Model

	s, err := r.service.Systems()
	if err != nil || len(s) == 0 {
		return
	}
	if err = r.setMemory(s[0], n); err != nil {
		return
	}
	if err = r.setDisks(s[0], n); err != nil {
		return
	}
	if err = r.setCPUs(s[0], n); err != nil {
		return
	}
	if err = r.setNetworkDevicesData(ch[0], n); err != nil {
		return
	}
	return
}

func (r RedfishClient) setMemory(s *redfish.ComputerSystem, n *model.Node) (err error) {
	mem, err := s.Memory()
	if err != nil {
		return
	}
	n.InspectionData.Inventory.Memory.PhysicalMb = calcTotalMemory(mem)
	return
}

func (r RedfishClient) setDisks(s *redfish.ComputerSystem, n *model.Node) (err error) {
	st, err := s.Storage()
	rootDisk := model.RootDisk{
		Rotational: true,
	}
	n.InspectionData.Inventory.Disks = make([]model.Disk, 0)
	re := regexp.MustCompile(`^(?i)(ssd|hdd)\s*(\d+)$`)
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

			//"SSD 1" or "HDD 2"
			match := re.FindStringSubmatch(s.Name)
			if match != nil {
				rootDisk.Size = int64(float64(s.CapacityBytes) * 1.074)
				rootDisk.Name = s.Name
				rootDisk.Model = s.Model
				rootDisk.Vendor = s.Manufacturer
				if s.RotationSpeedRPM == 0 {
					rootDisk.Rotational = rotational
				}
			}

			n.InspectionData.Inventory.Disks = append(n.InspectionData.Inventory.Disks, disk)
		}
	}

	n.InspectionData.RootDisk = rootDisk
	return
}
func (r RedfishClient) setCPUs(s *redfish.ComputerSystem, n *model.Node) (err error) {
	cpu, err := s.Processors()
	if err != nil || len(cpu) == 0 {
		return
	}
	//n.InspectionData.Inventory.CPU.Count = s.ProcessorSummary.LogicalProcessorCount / s.ProcessorSummary.Count
	n.InspectionData.Inventory.CPU.Count = s.ProcessorSummary.LogicalProcessorCount / 2 // threads
	n.InspectionData.Inventory.CPU.Architecture = strings.Replace(string(cpu[0].InstructionSet), "-", "_", 1)
	return
}

func (r RedfishClient) setNetworkDevicesData(c *redfish.Chassis, n *model.Node) (err error) {
	intfs := make(map[string]model.NodeInterface, 0)
	n.InspectionData.Inventory.Interfaces = make([]model.Interface, 0)
	na, err := c.NetworkAdapters()
	if err != nil {
		return
	}
	for _, a := range na {
		slot := a.Controllers[0].Location.PartLocation.LocationOrdinalValue
		nps, err := a.NetworkPorts()
		if err != nil {
			return err
		}
		for _, np := range nps {
			mac := np.AssociatedNetworkAddresses[0]
			id := mapInterfaceToNetbox(np.ID, slot)
			if np.LinkStatus == redfish.UpPortLinkStatus && n.InspectionData.BootInterface == "" {
				mac, err = parseMac(mac, '-')
				if err != nil {
					return err
				}
				n.InspectionData.Inventory.Boot.PxeInterface = mac
				n.InspectionData.BootInterface = "01-" + strings.ToLower(mac)
			}
			mac, err = parseMac(mac, ':')
			if err != nil {
				return err
			}
			//add baremetal ports
			if np.LinkStatus == redfish.UpPortLinkStatus {
				n.InspectionData.Inventory.Interfaces = append(n.InspectionData.Inventory.Interfaces, model.Interface{
					Name:       strings.ToLower(id),
					MacAddress: strings.ToLower(mac),
					Vendor:     &a.Manufacturer,
					Product:    a.Model,
					HasCarrier: true,
				})
			}

			intfs[id] = model.NodeInterface{
				Connection:     "",
				ConnectionIP:   "",
				Mac:            mac,
				PortLinkStatus: np.LinkStatus,
			}
		}
	}
	r.node.Interfaces = intfs
	n.InspectionData.Inventory.Boot.CurrentBootMode = "uefi"
	return
}

func parseMac(s string, sep rune) (string, error) {
	if len(s) < 12 {
		return "", fmt.Errorf("invalid MAC address: %s", s)
	}
	s = strings.ReplaceAll(s, ":", "")
	s = strings.ReplaceAll(s, "-", "")
	var buf bytes.Buffer
	for i, char := range s {
		buf.WriteRune(char)
		if i%2 == 1 && i != len(s)-1 {
			buf.WriteRune(sep)
		}

	}

	return buf.String(), nil
}

func mapInterfaceToNetbox(id string, slot int) (intf string) {
	p := strings.Split(id, ".")
	if len(p) <= 1 {
		return fmt.Sprintf("PCI%d-P%s", slot, id)
	}
	//NIC.Integrated.1-1-1 => L1
	if p[1] == "Integrated" {
		nr := strings.Split(p[2], "-")
		intf = "L" + nr[1]
	}
	//NIC.Slot.3-2-1 => PCI3-P2
	if p[1] == "Slot" {
		nr := strings.Split(p[2], "-")
		intf = fmt.Sprintf("PCI%s-P%s", nr[0], nr[1])
	}
	return
}
