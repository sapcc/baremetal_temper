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
	node         *model.Node
	log          *log.Entry
	cfg          config.Redfish
}

//NewRedfishClient creates redfish client
func NewRedfishClient(cfg config.Config, ctxLogger *log.Entry) *RedfishClient {
	return &RedfishClient{
		ClientConfig: &gofish.ClientConfig{
			Endpoint:  fmt.Sprintf("https://%s", "dummy.net"),
			Username:  cfg.Redfish.User,
			Password:  cfg.Redfish.Password,
			Insecure:  true,
			BasicAuth: true,
		},
		cfg: cfg.Redfish,
		log: ctxLogger,
	}
}

//SetEndpoint sets the redfish api endpoint
func (r *RedfishClient) SetEndpoint(n *model.Node) (err error) {
	r.ClientConfig.Endpoint = fmt.Sprintf("https://%s", n.RemoteIP)
	return
}

func (r *RedfishClient) BootImage(n *model.Node) (err error) {
	if err = r.connect(); err != nil {
		return
	}
	defer r.client.Logout()
	bootOverride := redfish.Boot{
		BootSourceOverrideTarget:  redfish.CdBootSourceOverrideTarget,
		BootSourceOverrideEnabled: redfish.OnceBootSourceOverrideEnabled,
	}
	if err = r.insertMedia(*r.cfg.BootImage); err != nil {
		return
	}
	return r.reboot(bootOverride)
}

func (r RedfishClient) insertMedia(image string) (err error) {
	m, err := r.client.Service.Managers()
	if err != nil {
		return
	}

	vms, err := m[0].VirtualMedia()
	if err != nil {
		return
	}
	var vm *redfish.VirtualMedia
	for _, v := range vms {
		for _, ty := range v.MediaTypes {
			if ty == redfish.CDMediaType || ty == redfish.DVDMediaType {
				vm = v
			}
		}
	}
	if vm.SupportsMediaInsert {
		if vm.Image != "" {
			err = vm.EjectMedia()
		}
		err = vm.InsertMedia(image, false, false)
	}

	return
}

func (r RedfishClient) CreateEventSubscription(n *model.Node) (err error) {
	es, err := r.client.Service.EventService()
	if err != nil {
		return
	}
	_, err = es.CreateEventSubscription(
		"https://baremetal_temper/events/"+n.Name,
		[]redfish.EventType{redfish.SupportedEventTypes["Alert"], redfish.SupportedEventTypes["StatusChange"]},
		nil,
		redfish.RedfishEventDestinationProtocol,
		"Public",
		nil,
	)
	return
}

func (r RedfishClient) DeleteEventSubscription(n *model.Node) (err error) {
	es, err := r.client.Service.EventService()
	if err != nil {
		r.log.Error(err)
		return nil
	}
	if err := es.DeleteEventSubscription("https://baremetal_temper/events/" + n.Name); err != nil {
		r.log.Error(err)
	}
	return nil
}

func (r RedfishClient) reboot(boot redfish.Boot) (err error) {
	type shareParameters struct {
		Target string
	}
	type temp struct {
		ShareParameters shareParameters
		ImportBuffer    string
	}

	bootOverride := redfish.Boot{
		BootSourceOverrideTarget:  redfish.CdBootSourceOverrideTarget,
		BootSourceOverrideEnabled: redfish.OnceBootSourceOverrideEnabled,
	}

	service := r.client.Service

	sys, err := service.Systems()
	if err != nil {
		return
	}
	var dellRe = regexp.MustCompile(`R640|R740|R840`)
	if dellRe.MatchString(sys[0].Model) {
		m, err := service.Managers()
		if err != nil {
			return err
		}
		t := temp{
			ShareParameters: shareParameters{Target: "ALL"},
			ImportBuffer:    "<SystemConfiguration><Component FQDD=\"iDRAC.Embedded.1\"><Attribute Name=\"ServerBoot.1#BootOnce\">Enabled</Attribute><Attribute Name=\"ServerBoot.1#FirstBootDevice\">VCD-DVD</Attribute></Component></SystemConfiguration>",
		}
		resp, err := m[0].Client.Post("/redfish/v1/Managers/iDRAC.Embedded.1/Actions/Oem/EID_674_Manager.ImportSystemConfiguration", t)
		fmt.Println(resp, err)
	} else {
		err = sys[0].SetBoot(bootOverride)
		if err != nil {
			return
		}
	}
	if sys[0].PowerState != redfish.OnPowerState {
		err = sys[0].Reset(redfish.OnResetType)
	} else {
		err = sys[0].Reset(redfish.ForceRestartResetType)
	}
	if err != nil {
		return
	}

	return
}

func (r *RedfishClient) connect() (err error) {
	client, err := gofish.Connect(*r.ClientConfig)
	if err != nil {
		return
	}
	r.client = client
	return
}

//LoadInventory loads the node's inventory via it's redfish api
func (r *RedfishClient) LoadInventory(n *model.Node) (err error) {
	r.log.Debug("calling redfish api to load node info")
	if err = r.connect(); err != nil {
		return
	}
	r.node = n
	defer r.client.Logout()
	ch, err := r.client.Service.Chassis()
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

	s, err := r.client.Service.Systems()
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
			mac, err = parseMac(mac, ':')
			if err != nil {
				log.Errorf("no mac address for port id: %s, name: %s. ignoring it", id, np.Name)
				continue
			}
			r.addBootInterface(id, np, n)
			//add baremetal ports (only link up and no integrated ports)
			if np.LinkStatus == redfish.UpPortLinkStatus && strings.Contains(id, "PCI") {
				n.InspectionData.Inventory.Interfaces = append(n.InspectionData.Inventory.Interfaces, model.Interface{
					Name:       strings.ToLower(id),
					MacAddress: strings.ToLower(mac),
					Vendor:     &a.Manufacturer,
					Product:    a.Model,
					HasCarrier: true,
				})
			}

			intfs[id] = model.NodeInterface{
				Mac:            mac,
				PortLinkStatus: np.LinkStatus,
			}
		}
	}
	r.node.Interfaces = intfs
	n.InspectionData.Inventory.Boot.CurrentBootMode = "uefi"
	return
}

func (r RedfishClient) addBootInterface(id string, np *redfish.NetworkPort, n *model.Node) {
	mac := np.AssociatedNetworkAddresses[0]
	if np.LinkStatus == redfish.UpPortLinkStatus && n.InspectionData.BootInterface == "" {
		mac, err := parseMac(mac, '-')
		if err != nil {
			log.Errorf("no mac address for port id: %s, name: %s", id, np.Name)
		} else {
			n.InspectionData.Inventory.Boot.PxeInterface = mac
			n.InspectionData.BootInterface = "01-" + strings.ToLower(mac)
		}
	}
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
