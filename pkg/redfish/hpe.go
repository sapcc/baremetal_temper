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
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/sapcc/baremetal_temper/pkg/clients"
	"github.com/sapcc/baremetal_temper/pkg/config"
	log "github.com/sirupsen/logrus"
	"github.com/stmcginnis/gofish/common"
	"github.com/stmcginnis/gofish/redfish"
)

type Hpe struct {
	Default
}

func NewHpe(remoteIP string, cfg config.Config, ctxLogger *log.Entry) (Redfish, error) {
	c := clients.NewRedfish(cfg, ctxLogger)
	c.SetEndpoint(remoteIP)
	r := &Hpe{Default: Default{client: c, cfg: cfg, log: ctxLogger}}
	return r, r.check()
}

func (d *Hpe) GetData() (*Data, error) {
	defer d.client.Logout()
	if err := d.client.Connect(); err != nil {
		return d.Data, err
	}
	if d.Data != nil {
		return d.Data, nil
	}
	if err := d.client.Connect(); err != nil {
		return d.Data, err
	}
	defer d.client.Logout()
	d.Data = &Data{}
	if err := d.getVendorData(); err != nil {
		return d.Data, err
	}
	if err := d.getDisks(); err != nil {
		return d.Data, err
	}
	if err := d.getCPUs(); err != nil {
		return d.Data, err
	}
	if err := d.getMemory(); err != nil {
		return d.Data, err
	}
	if err := d.getNetworkDevices(); err != nil {
		return d.Data, err
	}
	return d.Data, nil
}

func (p *Hpe) getNetworkDevices() (err error) {
	resp, err := p.client.Client.Get("/redfish/v1/Systems/1/BaseNetworkAdapters/")
	if err != nil {
		return
	}
	var result common.Collection
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return
	}
	p.Data.Inventory.Interfaces = make([]Interface, 0)
	for _, l := range result.ItemLinks {
		resp, err := p.client.Client.Get(l)
		if err != nil {
			return err
		}
		var r HpeBaseNetworkAdapter
		err = json.NewDecoder(resp.Body).Decode(&r)
		if err != nil {
			return err
		}

		// sort hex mac address
		sort.Slice(r.PhysicalPorts, func(i, j int) bool {
			s1 := strings.ReplaceAll(r.PhysicalPorts[i].MACAddress, ":", "")
			s2 := strings.ReplaceAll(r.PhysicalPorts[j].MACAddress, ":", "")
			mac1, _ := strconv.ParseInt(s1, 16, 64)
			mac2, _ := strconv.ParseInt(s2, 16, 64)
			return mac1 < mac2
		})
		for i, e := range r.PhysicalPorts {
			name, port, nic := p.mapInterfaceToNetbox(r.StructuredName, i+1)
			mac, err := parseMac(e.MACAddress, ':')
			if err != nil {
				continue
			}
			var ls redfish.PortLinkStatus
			if redfish.LinkUpLinkStatus == e.LinkStatus {
				ls = redfish.UpPortLinkStatus
			} else {
				ls = redfish.DownPortLinkStatus
			}
			p.addBootInterface(name, e)
			p.Data.Inventory.Interfaces = append(p.Data.Inventory.Interfaces, Interface{
				Name:           strings.ToLower(name),
				MacAddress:     strings.ToLower(mac),
				Vendor:         &r.Name,
				Product:        r.Name,
				HasCarrier:     true,
				Nic:            nic,
				Port:           port,
				PortLinkStatus: ls,
			})
		}
	}
	p.Data.Inventory.Boot.CurrentBootMode = "uefi"
	return
}

func (p *Hpe) addBootInterface(name string, np redfish.EthernetInterface) {
	if strings.HasPrefix(name, "L") {
		return
	}
	if np.LinkStatus == redfish.LinkUpLinkStatus && p.Data.BootInterface == "" {
		mac, err := parseMac(np.MACAddress, '-')
		if err != nil {
			p.log.Errorf("no mac address for port id: %s, name: %s", name, np.Name)
		} else {
			p.Data.Inventory.Boot.PxeInterface = mac
			p.Data.BootInterface = "01-" + strings.ToLower(mac)
		}
	}
}

func (p *Hpe) mapInterfaceToNetbox(id string, slot int) (name string, port, nic int) {
	//NIC.Slot.5.1
	n := strings.Split(id, ".")
	if len(n) < 4 {
		return
	}
	//NIC.FlexLOM.1.1 => L1
	if n[1] == "FlexLOM" {
		return "L" + strconv.Itoa(slot), slot, 0
	}
	//NIC.Slot.5.1  => PCI3-P2
	if n[1] == "Slot" {
		nic, _ = strconv.Atoi(n[2])
		return fmt.Sprintf("NIC%s-port%d", n[2], slot), slot, nic
	}
	return
}

func (p *Hpe) getCPUs() (err error) {
	s, err := p.client.Client.Service.Systems()
	if err != nil {
		return
	}
	cpus, err := s[0].Processors()
	if err != nil || len(cpus) == 0 {
		return
	}
	cores := 0
	for _, c := range cpus {
		cores = cores + c.TotalCores
	}
	p.Data.Inventory.CPU.Count = cores
	p.Data.Inventory.CPU.Architecture = strings.Replace(string(cpus[0].InstructionSet), "-", "_", 1)
	return
}

func (p *Hpe) getDisks() (err error) {
	s, err := p.client.Client.Service.Systems()
	if err != nil {
		return
	}
	st, err := s[0].Storage()
	p.Data.Inventory.Disks = make([]Disk, 0)
	for _, s := range st {
		vol, err := s.Volumes()
		if err != nil {
			fmt.Println(err, "volumes")
			continue
		}
		if len(vol) == 1 {
			dr, err := vol[0].Drives()
			if err != nil {
				fmt.Println(err)
				continue
			}
			if len(dr) == 2 {
				p.addDisk(dr[0], s.Name)
				continue
			}
		}
		ds, errD := s.Drives()
		if errD != nil {
			continue
		}
		for _, d := range ds {
			p.addDisk(d, s.Name)
		}
	}
	return
}

func (p *Hpe) addDisk(d *redfish.Drive, storageName string) {
	rotational := true
	if d.RotationSpeedRPM == 0 {
		rotational = false
	}
	disk := Disk{
		ID:     d.ID,
		Name:   d.Name,
		Model:  d.Model,
		Vendor: d.Manufacturer,
		//inspector converts bytes to gibibyte
		Size:       int64(float64(d.CapacityBytes) * 1.074),
		Rotational: rotational,
	}
	p.Data.Inventory.Disks = append(p.Data.Inventory.Disks, disk)
	if strings.Contains(storageName, "Boot Controller") {
		p.Data.RootDisk = RootDisk{
			Size:       int64(float64(d.CapacityBytes) * 1.074),
			Name:       d.Name,
			Model:      d.Model,
			Vendor:     d.Manufacturer,
			Rotational: rotational,
		}
	}
}
