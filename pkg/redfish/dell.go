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

import "github.com/stmcginnis/gofish/redfish"

type Dell struct {
	Default
}

func (d Dell) getVendorData() {
	//node.log.Debug("calling redfish api to load node info")
	defer d.client.Logout()
	ch, err := d.client.Client.Service.Chassis()
	if err != nil || len(ch) == 0 {
		return
	}

	s, err := d.client.Client.Service.Systems()
	if err != nil || len(s) == 0 {
		return
	}

	d.node.InspectionData.Inventory.SystemVendor.Manufacturer = ch[0].Manufacturer
	d.node.InspectionData.Inventory.SystemVendor.SerialNumber = ch[0].SKU
	d.node.InspectionData.Inventory.SystemVendor.Model = s[0].Model

	d.node.InspectionData.Inventory.SystemVendor.ProductName = ch[0].Model
}

func (d *Dell) rebootFromVirtualMedia(boot redfish.Boot) (err error) {
	//n.log.Debug("boot from virtual media")
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

	return d.power(false, true)
}
