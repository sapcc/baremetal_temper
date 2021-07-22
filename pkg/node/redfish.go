/**
 * Copyright 2021 SAP SE
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package node

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/stmcginnis/gofish/redfish"
	"k8s.io/apimachinery/pkg/util/wait"
)

//LoadInventory loads the node's inventory via it's redfish api
func (n *Node) loadInventory() (err error) {
	n.log.Debug("calling redfish api to load node info")
	if err = n.Clients.Redfish.Connect(); err != nil {
		return
	}
	defer n.Clients.Redfish.Logout()
	ch, err := n.Clients.Redfish.Client.Service.Chassis()
	if err != nil || len(ch) == 0 {
		return
	}

	s, err := n.Clients.Redfish.Client.Service.Systems()
	if err != nil || len(s) == 0 {
		return
	}

	n.InspectionData.Inventory.SystemVendor.Manufacturer = ch[0].Manufacturer
	n.InspectionData.Inventory.SystemVendor.SerialNumber = ch[0].SerialNumber
	n.InspectionData.Inventory.SystemVendor.Model = s[0].Model

	// not performant string comparison due to toLower
	if strings.Contains(strings.ToLower(ch[0].Manufacturer), "dell") {
		n.InspectionData.Inventory.SystemVendor.SerialNumber = ch[0].SKU
	}
	n.InspectionData.Inventory.SystemVendor.ProductName = ch[0].Model

	if err = n.setMemory(s[0]); err != nil {
		return
	}
	if err = n.setDisks(s[0]); err != nil {
		return
	}
	if err = n.setCPUs(s[0]); err != nil {
		return
	}
	if err = n.setNetworkDevicesData(ch[0]); err != nil {
		return
	}
	return
}

func (n *Node) setMemory(s *redfish.ComputerSystem) (err error) {
	mem, err := s.Memory()
	if err != nil {
		return
	}
	n.InspectionData.Inventory.Memory.PhysicalMb = calcTotalMemory(mem)
	return
}

func (n *Node) setDisks(s *redfish.ComputerSystem) (err error) {
	st, err := s.Storage()
	rootDisk := RootDisk{
		Rotational: true,
	}
	n.InspectionData.Inventory.Disks = make([]Disk, 0)
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
			disk := Disk{
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
func (n *Node) setCPUs(s *redfish.ComputerSystem) (err error) {
	cpu, err := s.Processors()
	if err != nil || len(cpu) == 0 {
		return
	}
	n.InspectionData.Inventory.CPU.Count = s.ProcessorSummary.LogicalProcessorCount / 2 // threads
	n.InspectionData.Inventory.CPU.Architecture = strings.Replace(string(cpu[0].InstructionSet), "-", "_", 1)
	return
}

func (n *Node) setNetworkDevicesData(c *redfish.Chassis) (err error) {
	pciHigh, pciLow := 0, 0
	intfs := make(map[string]NodeInterface, 0)
	n.InspectionData.Inventory.Interfaces = make([]Interface, 0)
	na, err := c.NetworkAdapters()
	if err != nil {
		return
	}

	for _, a := range na {
		slot := a.Controllers[0].Location.PartLocation.LocationOrdinalValue
		if slot > pciHigh {
			pciHigh = slot
		} else if slot < pciLow {
			pciLow = slot
		}
		nps, err := a.NetworkPorts()
		if err != nil {
			return err
		}
		for _, np := range nps {
			mac := np.AssociatedNetworkAddresses[0]
			id := mapInterfaceToNetbox(np.ID, slot)
			mac, err = parseMac(mac, ':')
			if err != nil {
				n.log.Errorf("no mac address for port id: %s, name: %s. ignoring it", id, np.Name)
				continue
			}
			n.addBootInterface(id, np)
			//add baremetal ports (only link up and no integrated ports)
			if np.LinkStatus == redfish.UpPortLinkStatus && strings.Contains(id, "PCI") {
				n.InspectionData.Inventory.Interfaces = append(n.InspectionData.Inventory.Interfaces, Interface{
					Name:       strings.ToLower(id),
					MacAddress: strings.ToLower(mac),
					Vendor:     &a.Manufacturer,
					Product:    a.Model,
					HasCarrier: true,
				})
			}

			intfs[id] = NodeInterface{
				Mac:            mac,
				PortLinkStatus: np.LinkStatus,
			}
		}
	}
	n.Interfaces = intfs
	n.InspectionData.Inventory.Boot.CurrentBootMode = "uefi"
	return
}

func (n *Node) addBootInterface(id string, np *redfish.NetworkPort) {
	mac := np.AssociatedNetworkAddresses[0]
	if np.LinkStatus == redfish.UpPortLinkStatus && n.InspectionData.BootInterface == "" {
		mac, err := parseMac(mac, '-')
		if err != nil {
			n.log.Errorf("no mac address for port id: %s, name: %s", id, np.Name)
		} else {
			n.InspectionData.Inventory.Boot.PxeInterface = mac
			n.InspectionData.BootInterface = "01-" + strings.ToLower(mac)
		}
	}
}

func (n *Node) bootImage() (err error) {
	n.log.Debugf("booting image for cable check: %s", *n.cfg.Redfish.BootImage)
	if err = n.Clients.Redfish.Connect(); err != nil {
		return
	}
	defer n.Clients.Redfish.Logout()
	bootOverride := redfish.Boot{
		BootSourceOverrideTarget:  redfish.CdBootSourceOverrideTarget,
		BootSourceOverrideEnabled: redfish.OnceBootSourceOverrideEnabled,
	}
	if err = n.insertMedia(*n.cfg.Redfish.BootImage); err != nil {
		return
	}
	return n.rebootFromVirtualMedia(bootOverride)
}

func (n *Node) insertMedia(image string) (err error) {
	n.log.Debug("insert virtual media")
	vm, err := n.getMedia()
	if err != nil {
		return
	}
	if vm.SupportsMediaInsert {
		if vm.Image != "" {
			err = vm.EjectMedia()
		}
		var patchRe = regexp.MustCompile(`SR950`)
		if patchRe.MatchString(n.InspectionData.Inventory.SystemVendor.Model) {
			type temp struct {
				Image          string
				Inserted       bool `json:"Inserted"`
				WriteProtected bool `json:"WriteProtected"`
			}
			t := temp{
				Image:          image,
				Inserted:       true,
				WriteProtected: false,
			}
			return vm.InsertMedia(image, "PATCH", t)
		}
		type temp struct {
			Image          string
			Inserted       bool `json:"-"`
			WriteProtected bool `json:"-"`
		}
		t := temp{
			Image:          image,
			Inserted:       true,
			WriteProtected: false,
		}
		return vm.InsertMedia(image, "POST", t)
	}
	return
}

func (n *Node) ejectMedia() (err error) {
	n.log.Debug("eject media image")
	vm, err := n.getMedia()
	if err != nil {
		return
	}
	if vm.SupportsMediaInsert {
		if vm.Image != "" {
			var patchRe = regexp.MustCompile(`SR950`)
			if patchRe.MatchString(n.InspectionData.Inventory.SystemVendor.Model) {
				return vm.InsertMedia("", "PATCH", nil)
			}
			return vm.EjectMedia()
		}
	}
	return
}

func (n *Node) getMedia() (vm *redfish.VirtualMedia, err error) {
	m, err := n.Clients.Redfish.Client.Service.Managers()
	if err != nil {
		return
	}
	vms, err := m[0].VirtualMedia()
	if err != nil {
		return
	}
vmLoop:
	for _, v := range vms {
		for _, ty := range v.MediaTypes {
			if ty == redfish.CDMediaType || ty == redfish.DVDMediaType {
				vm = v
				break vmLoop
			}
		}
	}
	return
}

func (n *Node) CreateEventSubscription() (err error) {
	es, err := n.Clients.Redfish.Client.Service.EventService()
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

func (n *Node) DeleteEventSubscription() (err error) {
	es, err := n.Clients.Redfish.Client.Service.EventService()
	if err != nil {
		n.log.Error(err)
		return nil
	}
	if err := es.DeleteEventSubscription("https://baremetal_temper/events/" + n.Name); err != nil {
		n.log.Error(err)
	}
	return nil
}

func (n *Node) rebootFromVirtualMedia(boot redfish.Boot) (err error) {
	n.log.Debug("boot from virtual media")
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

	service := n.Clients.Redfish.Client.Service

	sys, err := service.Systems()
	if err != nil {
		return
	}
	var dellRe = regexp.MustCompile(`R640|R740|R840`)
	if dellRe.MatchString(n.InspectionData.Inventory.SystemVendor.Model) {
		m, err := service.Managers()
		if err != nil {
			return err
		}
		t := temp{
			ShareParameters: shareParameters{Target: "ALL"},
			ImportBuffer:    "<SystemConfiguration><Component FQDD=\"iDRAC.Embedded.1\"><Attribute Name=\"ServerBoot.1#BootOnce\">Enabled</Attribute><Attribute Name=\"ServerBoot.1#FirstBootDevice\">VCD-DVD</Attribute></Component></SystemConfiguration>",
		}
		_, err = m[0].Client.Post("/redfish/v1/Managers/iDRAC.Embedded.1/Actions/Oem/EID_674_Manager.ImportSystemConfiguration", t)
	} else {
		err = sys[0].SetBoot(bootOverride)
		if err != nil {
			return
		}
	}
	return n.power(false, true)
}

func (n *Node) power(forceOff bool, restart bool) (err error) {
	if err = n.Clients.Redfish.Connect(); err != nil {
		return
	}
	sys, err := n.Clients.Redfish.Client.Service.Systems()
	if err != nil {
		return
	}
	if forceOff {
		return sys[0].Reset(redfish.ForceOffResetType)
	}
	if restart {
		return sys[0].Reset(redfish.ForceRestartResetType)
	}
	n.log.Debugf("node power state: %s", sys[0].PowerState)
	if sys[0].PowerState != redfish.OnPowerState {
		n.log.Debug("node power on")
		err = sys[0].Reset(redfish.OnResetType)
		// lets give the server some time to fully boot,
		// otherwise redfish resources may not be ready (e.g. ports)
		TimeoutTask(1 * time.Minute)()
	}

	if err != nil {
		return
	}
	return
}

func (n *Node) waitPowerStateOn() (err error) {
	n.log.Infof("waiting for node to power on")
	cf := wait.ConditionFunc(func() (bool, error) {
		sys, err := n.Clients.Redfish.Client.Service.Systems()
		if err != nil {
			return false, fmt.Errorf("cannot power on node")
		}
		p := sys[0].PowerState
		n.log.Debugf("node power state: %s", p)
		if p != redfish.OnPowerState {
			return false, nil
		}
		return true, nil
	})
	return wait.Poll(10*time.Second, 5*time.Minute, cf)
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
