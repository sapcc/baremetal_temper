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
	"github.com/sapcc/baremetal_temper/pkg/clients"
	"github.com/sapcc/baremetal_temper/pkg/config"
	log "github.com/sirupsen/logrus"
	"github.com/stmcginnis/gofish/redfish"
)

type Dell struct {
	Default
}

func NewDell(remoteIP string, cfg config.Config, ctxLogger *log.Entry) (Redfish, error) {
	c := clients.NewRedfish(cfg, ctxLogger)
	c.SetEndpoint(remoteIP)
	r := &Dell{Default: Default{client: c, cfg: cfg, log: ctxLogger}}
	return r, r.check()
}

func (d *Dell) GetData() (*Data, error) {
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

func (d *Dell) getVendorData() (err error) {
	d.log.Debug("calling redfish api to load node info")
	ch, err := d.client.Client.Service.Chassis()
	if err != nil || len(ch) == 0 {
		return
	}

	s, err := d.client.Client.Service.Systems()
	if err != nil || len(s) == 0 {
		return
	}

	d.Data.Inventory.SystemVendor.Manufacturer = ch[0].Manufacturer
	d.Data.Inventory.SystemVendor.SerialNumber = ch[0].SKU
	d.Data.Inventory.SystemVendor.Model = s[0].Model

	d.Data.Inventory.SystemVendor.ProductName = ch[0].Model
	return
}

func (d *Dell) rebootFromVirtualMedia(boot redfish.Boot) (err error) {
	d.log.Debug("boot from virtual media")
	type shareParameters struct {
		Target string
	}
	type temp struct {
		ShareParameters shareParameters
		ImportBuffer    string
	}

	service := d.client.Client.Service

	m, err := service.Managers()
	if err != nil {
		return err
	}
	t := temp{
		ShareParameters: shareParameters{Target: "ALL"},
		ImportBuffer:    "<SystemConfiguration><Component FQDD=\"iDRAC.Embedded.1\"><Attribute Name=\"ServerBoot.1#BootOnce\">Enabled</Attribute><Attribute Name=\"ServerBoot.1#FirstBootDevice\">VCD-DVD</Attribute></Component></SystemConfiguration>",
	}
	_, err = m[0].Client.Post("/redfish/v1/Managers/iDRAC.Embedded.1/Actions/Oem/EID_674_Manager.ImportSystemConfiguration", t)
	if err != nil {
		return
	}

	return d.Power(false, true)
}
