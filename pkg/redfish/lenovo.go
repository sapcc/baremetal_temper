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
	"strconv"
	"strings"

	"github.com/sapcc/baremetal_temper/pkg/clients"
	"github.com/sapcc/baremetal_temper/pkg/config"
	log "github.com/sirupsen/logrus"
)

type Lenovo struct {
	Default
}

func NewLenovo(remoteIP string, cfg config.Config, ctxLogger *log.Entry) (Redfish, error) {
	c := clients.NewRedfish(cfg, ctxLogger)
	c.SetEndpoint(remoteIP)
	r := &Lenovo{Default: Default{client: c, cfg: cfg, log: ctxLogger}}
	return r, r.check()
}

func (d *Lenovo) GetData() (*Data, error) {
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

func (p *Lenovo) InsertMedia(image string) (err error) {
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
	return
}

func (p *Lenovo) EjectMedia() (err error) {
	p.log.Debug("eject media image")
	vm, err := p.getDVDMediaType()
	if err != nil {
		return
	}
	if vm.SupportsMediaInsert {
		if vm.Image != "" {
			return vm.InsertMedia("", "PATCH", nil)
		}
	}
	return
}

func (p *Lenovo) getDisks() (err error) {
	s, err := p.client.Client.Service.Systems()
	if err != nil {
		return
	}
	st, err := s[0].Storage()
	rootDisk := RootDisk{
		Rotational: true,
	}
	p.Data.Inventory.Disks = make([]Disk, 0)
	for _, s := range st {
		vols, err := s.Volumes()
		if err != nil {
			continue
		}
		// no raid
		if len(vols) == 0 {
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
				p.Data.Inventory.Disks = append(p.Data.Inventory.Disks, disk)
			}
		} else {
			for _, v := range vols {
				ds, err := v.Drives()
				if err != nil {
					continue
				}
				rootDisk.Size = int64(float64(v.CapacityBytes) * 1.074)
				rootDisk.Name = ds[0].Name
				rootDisk.Model = ds[0].Model
				rootDisk.Vendor = ds[0].Manufacturer
				if ds[0].RotationSpeedRPM == 0 {
					rootDisk.Rotational = false
				} else {
					rootDisk.Rotational = true
				}
			}
		}
	}
	if rootDisk.Size == 0 {
		panic("unable to detect root disk")
	}
	p.Data.RootDisk = rootDisk
	fmt.Println(p.Data.RootDisk)
	fmt.Println(p.Data.Inventory.Disks)
	return
}

func (p *Lenovo) getNetworkDevices() (err error) {
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
			var (
				name string
				port int
				nic  int
			)
			fmt.Println(strings.ToLower(a.Manufacturer), np.ID, slot)
			if strings.Contains(strings.ToLower(a.Manufacturer), "intel") {
				name, port, nic = p.mapInterfaceToNetbox(np.ID, 0)
			} else {
				name, port, nic = p.mapInterfaceToNetbox(np.ID, slot)
			}
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

func (p *Lenovo) mapInterfaceToNetbox(id string, slot int) (name string, port, nic int) {
	port, _ = strconv.Atoi(id)
	if slot == 0 {
		return fmt.Sprintf("L%s", id), port, 0
	}
	return fmt.Sprintf("NIC%d-port%s", slot, id), port, slot
}
