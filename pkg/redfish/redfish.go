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
	"strconv"
	"strings"
	"time"

	"github.com/sapcc/baremetal_temper/pkg/clients"
	"github.com/sapcc/baremetal_temper/pkg/config"
	"github.com/stmcginnis/gofish"
	"github.com/stmcginnis/gofish/redfish"
	"k8s.io/apimachinery/pkg/util/wait"

	log "github.com/sirupsen/logrus"
)

type Data struct {
	RootDisk      RootDisk  `json:"root_disk"`
	BootInterface string    `json:"boot_interface"`
	Inventory     Inventory `json:"inventory"`
	Logs          string    `json:"logs"`
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
	Lldp           map[string]string      `json:"lldp"`
	Product        string                 `json:"product"`
	Vendor         *string                `json:"vendor"`
	Name           string                 `json:"name"`
	HasCarrier     bool                   `json:"has_carrier"`
	IP4Address     string                 `json:"ipv4_address"`
	ClientID       *string                `json:"client_id"`
	MacAddress     string                 `json:"mac_address"`
	PortLinkStatus redfish.PortLinkStatus `json:"-"`
	Nic            int                    `json:"-"`
	Port           int                    `json:"-"`
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
	ID                 string  `json:"-"`
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

type RootDisk struct {
	Rotational bool   `json:"rotational"`
	Vendor     string `json:"vendor"`
	Name       string `json:"name"`
	Model      string `json:"model"`
	Serial     string `json:"serial"`
	Size       int64  `json:"size"`
}

type Redfish interface {
	GetData() (*Data, error)
	GetClientConfig() *gofish.ClientConfig
	getVendorData() (err error)
	getMemory() (err error)
	getCPUs() (err error)
	getDisks() (err error)
	getNetworkDevices() (err error)
	rebootFromVirtualMedia(boot redfish.Boot) (err error)
	mapInterfaceToNetbox(id string, slot int) (name string, port, nic int)

	Power(forceOff bool, restart bool) (err error)
	WaitPowerStateOn() (err error)
	BootFromImage(path string) (err error)
	EjectMedia() (err error)
	InsertMedia(image string) (err error)
}

type Default struct {
	client *clients.Redfish
	log    *log.Entry
	cfg    config.Config
	Data   *Data
}

func NewDefault(remoteIP string, cfg config.Config, ctxLogger *log.Entry) (Redfish, error) {
	c := clients.NewRedfish(cfg, ctxLogger)
	c.SetEndpoint(remoteIP)
	r := &Default{client: c, log: ctxLogger, cfg: cfg}
	return r, r.check()
}

func (p Default) check() (err error) {
	if err = p.client.Connect(); err != nil {
		return err
	}
	defer p.client.Logout()
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

func (p *Default) GetClientConfig() *gofish.ClientConfig {
	return p.client.ClientConfig
}

func (p *Default) GetData() (*Data, error) {
	if p.Data != nil {
		return p.Data, nil
	}
	if err := p.client.Connect(); err != nil {
		return p.Data, err
	}
	defer p.client.Logout()
	p.Data = &Data{}
	if err := p.getVendorData(); err != nil {
		return p.Data, err
	}
	if err := p.getDisks(); err != nil {
		return p.Data, err
	}
	if err := p.getCPUs(); err != nil {
		return p.Data, err
	}
	if err := p.getMemory(); err != nil {
		return p.Data, err
	}
	if err := p.getNetworkDevices(); err != nil {
		return p.Data, err
	}
	return p.Data, nil
}

func (p *Default) Power(forceOff bool, restart bool) (err error) {
	defer p.client.Logout()
	if err = p.client.Connect(); err != nil {
		return
	}
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
	p.log.Debugf("node power state: %s", s[0].PowerState)
	if s[0].PowerState == redfish.OffPowerState {
		p.log.Debug("node power on")
		err = s[0].Reset(redfish.OnResetType)
	}
	return
}

func (p *Default) WaitPowerStateOn() (err error) {
	defer p.client.Logout()
	if err = p.client.Connect(); err != nil {
		return
	}
	cf := wait.ConditionFunc(func() (bool, error) {
		sys, err := p.client.Client.Service.Systems()
		if err != nil {
			return false, fmt.Errorf("cannot power on node")
		}
		power := sys[0].PowerState
		p.log.Debugf("waiting for node to power on: current state: %s", power)
		if power != redfish.OnPowerState {
			return false, nil
		}
		return true, nil
	})

	r, err := cf()
	if r && err == nil {
		return
	}
	if err = wait.Poll(10*time.Second, 5*time.Minute, cf); err != nil {
		return
	}
	// lets give the server some time to fully boot. powerON state is not sufficient in most cases
	// otherwise redfish resources may not be ready (e.g. ports)
	time.Sleep(5 * time.Minute)
	return
}

func (p *Default) getVendorData() (err error) {
	ch, err := p.client.Client.Service.Chassis()
	if err != nil {
		return
	}
	s, err := p.client.Client.Service.Systems()
	if err != nil {
		return
	}
	p.Data.Inventory.SystemVendor.Manufacturer = ch[0].Manufacturer
	p.Data.Inventory.SystemVendor.SerialNumber = ch[0].SerialNumber
	p.Data.Inventory.SystemVendor.Model = s[0].Model
	p.Data.Inventory.SystemVendor.ProductName = ch[0].Model
	return
}

func (p *Default) getDisks() (err error) {
	s, err := p.client.Client.Service.Systems()
	if err != nil {
		return
	}
	st, err := s[0].Storage()
	rootDisk := RootDisk{
		Rotational: true,
	}
	p.Data.Inventory.Disks = make([]Disk, 0)
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
			p.Data.Inventory.Disks = append(p.Data.Inventory.Disks, disk)
		}
	}
	p.Data.RootDisk = rootDisk
	return
}

func (p *Default) getCPUs() (err error) {
	s, err := p.client.Client.Service.Systems()
	if err != nil {
		return
	}
	cpu, err := s[0].Processors()
	if err != nil || len(cpu) == 0 {
		return
	}
	p.Data.Inventory.CPU.Count = s[0].ProcessorSummary.LogicalProcessorCount / 2 // threads
	p.Data.Inventory.CPU.Architecture = strings.Replace(string(cpu[0].InstructionSet), "-", "_", 1)
	return
}

func (p *Default) getNetworkDevices() (err error) {
	ch, err := p.client.Client.Service.Chassis()
	if err != nil {
		return
	}

	p.Data.Inventory.Interfaces = make([]Interface, 0)
	na, err := ch[0].NetworkAdapters()
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
			name, port, nic := p.mapInterfaceToNetbox(np.ID, slot)
			mac, err = parseMac(mac, ':')
			if err != nil {
				p.log.Errorf("no mac address for port id: %s, name: %s. ignoring it", name, np.Name)
				continue
			}
			p.addBootInterface(name, np)
			p.Data.Inventory.Interfaces = append(p.Data.Inventory.Interfaces, Interface{
				Name:           strings.ToLower(name),
				MacAddress:     strings.ToLower(mac),
				Vendor:         &a.Manufacturer,
				Product:        a.Model,
				HasCarrier:     true,
				Nic:            nic,
				Port:           port,
				PortLinkStatus: np.LinkStatus,
			})
		}
	}
	p.Data.Inventory.Boot.CurrentBootMode = "uefi"
	return
}

func (p *Default) getMemory() (err error) {
	s, err := p.client.Client.Service.Systems()
	if err != nil {
		return
	}
	mem, err := s[0].Memory()
	if err != nil {
		return
	}
	p.Data.Inventory.Memory.PhysicalMb = calcTotalMemory(mem)
	return
}

func (p *Default) BootFromImage(path string) (err error) {
	defer p.client.Logout()
	if err = p.client.Connect(); err != nil {
		return
	}
	p.log.Debugf("booting image for cable check: %s", *p.cfg.Redfish.BootImage)
	bootOverride := redfish.Boot{
		BootSourceOverrideTarget:  redfish.CdBootSourceOverrideTarget,
		BootSourceOverrideEnabled: redfish.OnceBootSourceOverrideEnabled,
	}
	if err = p.InsertMedia(path); err != nil {
		return
	}
	return p.rebootFromVirtualMedia(bootOverride)
}

func (p *Default) rebootFromVirtualMedia(boot redfish.Boot) (err error) {
	p.log.Debug("boot from virtual media")
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

	return p.Power(false, true)
}

func (p *Default) InsertMedia(image string) (err error) {
	defer p.client.Logout()
	if err = p.client.Connect(); err != nil {
		return
	}
	p.log.Debug("insert virtual media")
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

func (p *Default) EjectMedia() (err error) {
	defer p.client.Logout()
	if err = p.client.Connect(); err != nil {
		return
	}
	p.log.Debug("eject media image")
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

func (p *Default) addBootInterface(name string, np *redfish.NetworkPort) {
	mac := np.AssociatedNetworkAddresses[0]
	if strings.HasPrefix(name, "L") {
		return
	}
	if np.LinkStatus == redfish.UpPortLinkStatus && p.Data.BootInterface == "" {
		mac, err := parseMac(mac, '-')
		if err != nil {
			p.log.Errorf("no mac address for port id: %s, name: %s", name, np.Name)
		} else {
			p.Data.Inventory.Boot.PxeInterface = mac
			p.Data.BootInterface = "01-" + strings.ToLower(mac)
		}
	}
}

func (p *Default) mapInterfaceToNetbox(id string, slot int) (name string, port, nic int) {
	//netbox example: NIC1-port1 (vmnic4)
	n := strings.Split(id, ".")
	if len(n) <= 1 {
		port, _ = strconv.Atoi(id)
		return fmt.Sprintf("NIC%d-port%s", slot, id), port, slot
	}
	//NIC.Integrated.1-1-1 => L1
	if n[1] == "Integrated" {
		nr := strings.Split(n[2], "-")
		port, _ = strconv.Atoi(nr[1])
		return "L" + nr[1], port, 0
	}
	//NIC.Slot.3-2-1 => PCI3-P2
	if n[1] == "Slot" {
		nr := strings.Split(n[2], "-")
		nic, _ = strconv.Atoi(nr[0])
		port, _ = strconv.Atoi(nr[1])
		return fmt.Sprintf("NIC%s-port%s", nr[0], nr[1]), port, nic
	}
	return
}

func handleError(err error) {
	if strings.Contains(err.Error(), "ResourceNotReadyRetry") {
	}
}
