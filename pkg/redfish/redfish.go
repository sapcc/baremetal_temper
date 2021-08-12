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

package redfish

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/sapcc/baremetal_temper/pkg/clients"
	"github.com/sapcc/baremetal_temper/pkg/node"
	"github.com/stmcginnis/gofish/redfish"
	"k8s.io/apimachinery/pkg/util/wait"
)

type Redfish interface {
	getVendorData() (err error)
	getMemory() (err error)
	getCPU() (err error)
	getDisks() (err error)
	getNetworkDevices() (err error)

	power(forceOff bool, restart bool) (err error)
	waitPowerStateOn() (err error)
	bootFromImage() (err error)
}

type Default struct {
	client *clients.Redfish
	node   *node.Node
}

func (p Default) init() (err error) {
	ch, err := p.client.Client.Service.Chassis()
	if err != nil || len(ch) == 0 {
		return fmt.Errorf("redfish chassis != 1")
	}

	s, err := p.client.Client.Service.Systems()
	if err != nil || len(s) == 0 {
		return fmt.Errorf("redfish systems != 1")
	}
	return
}

func (p Default) power(forceOff bool, restart bool) (err error) {
	s, err := p.client.Client.Service.Systems()
	if err != nil {
		return
	}
	if forceOff {
		return s[0].Reset(redfish.ForceOffResetType)
	}
	if restart {
		return s[0].Reset(redfish.ForceRestartResetType)
	}
	//n.log.Debugf("node power state: %s", s[0].PowerState)
	if s[0].PowerState != redfish.OnPowerState {
		//n.log.Debug("node power on")
		err = s[0].Reset(redfish.OnResetType)
		// lets give the server some time to fully boot,
		// otherwise redfish resources may not be ready (e.g. ports)
		time.Sleep(5 * time.Minute)
	}
	return
}

func (p *Default) waitPowerStateOn() (err error) {
	//n.log.Infof("waiting for node to power on")
	cf := wait.ConditionFunc(func() (bool, error) {
		sys, err := p.client.Client.Service.Systems()
		if err != nil {
			return false, fmt.Errorf("cannot power on node")
		}
		p := sys[0].PowerState
		//n.log.Debugf("node power state: %s", p)
		if p != redfish.OnPowerState {
			return false, nil
		}
		return true, nil
	})
	return wait.Poll(10*time.Second, 5*time.Minute, cf)
}

func (p Default) getVendorData() (err error) {
	ch, err := p.client.Client.Service.Chassis()
	if err != nil {
		return
	}
	s, err := p.client.Client.Service.Systems()
	if err != nil {
		return
	}
	p.node.InspectionData.Inventory.SystemVendor.Manufacturer = ch[0].Manufacturer
	p.node.InspectionData.Inventory.SystemVendor.SerialNumber = ch[0].SerialNumber
	p.node.InspectionData.Inventory.SystemVendor.Model = s[0].Model
	p.node.InspectionData.Inventory.SystemVendor.ProductName = ch[0].Model
	return
}
func (p Default) getMemory() (err error) {
	s, err := p.client.Client.Service.Systems()
	if err != nil {
		return
	}
	st, err := s[0].Storage()
	rootDisk := RootDisk{
		Rotational: true,
	}
	p.node.InspectionData.Inventory.Disks = make([]node.Disk, 0)
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

			p.node.InspectionData.Inventory.Disks = append(p.node.InspectionData.Inventory.Disks, disk)
		}
	}

	p.node.InspectionData.RootDisk = rootDisk
	return
}
func (p Default) setCPUs() (err error) {
	s, err := p.client.Client.Service.Systems()
	if err != nil {
		return
	}
	cpu, err := s[0].Processors()
	if err != nil || len(cpu) == 0 {
		return
	}
	p.node.InspectionData.Inventory.CPU.Count = s[0].ProcessorSummary.LogicalProcessorCount / 2 // threads
	p.node.InspectionData.Inventory.CPU.Architecture = strings.Replace(string(cpu[0].InstructionSet), "-", "_", 1)
	return
}

func (p Default) setNetworkDevicesData() (err error) {
	ch, err := p.client.Client.Service.Chassis()
	if err != nil {
		return
	}
	pciHigh, pciLow := 0, 0
	intfs := make(map[string]node.NodeInterface, 0)
	p.node.InspectionData.Inventory.Interfaces = make([]node.Interface, 0)
	na, err := ch[0].NetworkAdapters()
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
				//n.log.Errorf("no mac address for port id: %s, name: %s. ignoring it", id, np.Name)
				continue
			}
			p.addBootInterface(id, np)
			//add baremetal ports (only link up and no integrated ports)
			if np.LinkStatus == redfish.UpPortLinkStatus && strings.Contains(id, "PCI") {
				p.node.InspectionData.Inventory.Interfaces = append(p.node.InspectionData.Inventory.Interfaces, node.Interface{
					Name:       strings.ToLower(id),
					MacAddress: strings.ToLower(mac),
					Vendor:     &a.Manufacturer,
					Product:    a.Model,
					HasCarrier: true,
				})
			}

			intfs[id] = node.NodeInterface{
				Mac:            mac,
				PortLinkStatus: np.LinkStatus,
			}
		}
	}
	p.node.Interfaces = intfs
	p.node.InspectionData.Inventory.Boot.CurrentBootMode = "uefi"
	return
}

func (p *Default) bootFromImage(path string) (err error) {
	//n.log.Debugf("booting image for cable check: %s", *n.cfg.Redfish.BootImage)
	bootOverride := redfish.Boot{
		BootSourceOverrideTarget:  redfish.CdBootSourceOverrideTarget,
		BootSourceOverrideEnabled: redfish.OnceBootSourceOverrideEnabled,
	}
	if err = p.insertMedia(path); err != nil {
		return
	}
	return p.rebootFromVirtualMedia(bootOverride)
}

func (p *Default) rebootFromVirtualMedia(boot redfish.Boot) (err error) {
	//n.log.Debug("boot from virtual media")
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

	service := p.client.Client.Service

	sys, err := service.Systems()
	if err != nil {
		return
	}

	err = sys[0].SetBoot(bootOverride)
	if err != nil {
		return
	}

	return p.power(false, true)
}

func (p *Default) insertMedia(image string) (err error) {
	//n.log.Debug("insert virtual media")
	vm, err := p.getDVDMediaType()
	if err != nil {
		return
	}
	if vm.SupportsMediaInsert {
		if vm.Image != "" {
			err = vm.EjectMedia()
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

func (p *Default) ejectMedia() (err error) {
	//n.log.Debug("eject media image")
	vm, err := p.getDVDMediaType()
	if err != nil {
		return
	}
	if vm.SupportsMediaInsert {
		if vm.Image != "" {
			return vm.EjectMedia()
		}
	}
	return
}

func (p *Default) getDVDMediaType() (vm *redfish.VirtualMedia, err error) {
	m, err := p.client.Client.Service.Managers()
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

func (p *Default) addBootInterface(id string, np *redfish.NetworkPort) {
	mac := np.AssociatedNetworkAddresses[0]
	if np.LinkStatus == redfish.UpPortLinkStatus && p.node.InspectionData.BootInterface == "" {
		mac, err := parseMac(mac, '-')
		if err != nil {
			//n.log.Errorf("no mac address for port id: %s, name: %s", id, np.Name)
		} else {
			p.node.InspectionData.Inventory.Boot.PxeInterface = mac
			p.node.InspectionData.BootInterface = "01-" + strings.ToLower(mac)
		}
	}
}
