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
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/netbox-community/go-netbox/netbox/client/dcim"
	"github.com/netbox-community/go-netbox/netbox/client/ipam"
	"github.com/netbox-community/go-netbox/netbox/models"
)

//Update node serial and primaryIP. Does not return error to not trigger errorhandler and cleanup of node
func (n *Node) Update() error {
	params := models.WritableDeviceWithConfigContext{
		Serial: n.InspectionData.Inventory.SystemVendor.SerialNumber,
	}
	d, err := n.getDevice(nil, &n.Name)
	if err != nil {
		n.log.Error(err)
		return nil
	}
	if d.PrimaryIP == nil {
		ips, err := n.Clients.Netbox.Client.Ipam.IpamIPAddressesList(&ipam.IpamIPAddressesListParams{
			Address: &n.PrimaryIP,
			Context: context.Background(),
		}, nil)
		if err != nil {
			n.log.Error(err)
			return nil
		}
		if len(ips.Payload.Results) > 1 || len(ips.Payload.Results) == 0 {
			n.log.Errorf("could not find ip %s", n.PrimaryIP)
			return nil
		}
		params.PrimaryIp4 = &ips.Payload.Results[0].ID
	}

	_, err = n.updateNodeInfo(n.Name, params)
	if err != nil {
		n.log.Error(err)
		return nil
	}

	if err = n.updateNodeInterfaces(); err != nil {
		n.log.Error(err)
		return nil
	}

	return nil
}

//LoadIpamAddresses loads all ipam addresse of a node
func (n *Node) loadIpamAddresses() (err error) {
	n.log.Debug("calling netbox api to load ipam Addresses")
	split := strings.Split(n.Name, "-")
	if len(split) == 1 {
		return fmt.Errorf("wrong node name format: node[001]-[block_name]")
	}
	block := strings.Split(n.Name, "-")[1]
	name := strings.Split(n.Name, "-")[0]
	var limit int64
	limit = 200
	al, err := n.Clients.Netbox.Client.Ipam.IpamIPAddressesList(&ipam.IpamIPAddressesListParams{
		Q:       &block,
		Context: context.Background(),
		Limit:   &limit,
	}, nil)
	if err != nil {
		return
	}
	addr := al.Payload.Results
	for _, a := range addr {
		if strings.Contains(a.DNSName, n.Name) {
			ip, _, err := net.ParseCIDR(*a.Address)
			if err != nil {
				return err
			}
			n.PrimaryIP = ip.String()
		}

		if strings.Contains(a.DNSName, name+"r") {
			ip, _, err := net.ParseCIDR(*a.Address)
			if err != nil {
				return err
			}
			n.RemoteIP = ip.String()
			n.InspectionData.Inventory.BmcAddress = a.DNSName
		}
		if strings.Contains(a.DNSName, name) {
			n.IpamAddresses = append(n.IpamAddresses, *a)
		}
	}
	return
}

//SetStatusStaged does not return error to not trigger errorhandler and cleanup of node
func (n *Node) SetStatus() error {
	errors := make([]string, 0)
	for _, t := range n.Tasks {
		if t.Error != "" {
			m := fmt.Sprintf("%s failed: %s", t.Name, t.Error)
			errors = append(errors, m)
		}
	}

	if len(errors) != 0 {
		return n.setStatusFailed(strings.Join(errors, " "))
	}
	n.setStatusStaged()
	return nil
}

//SetStatusStaged does not return error to not trigger errorhandler and cleanup of node
func (n *Node) setStatusStaged() error {
	p, err := n.updateNodeInfo(n.Name, models.WritableDeviceWithConfigContext{
		Status:   models.DeviceWithConfigContextStatusValueStaged,
		Comments: "temper successful",
	})
	if err != nil {
		n.log.Error(err)
		return nil
	}
	if *p.Payload.Status.Value != models.DeviceWithConfigContextStatusValueStaged {
		n.log.Errorf("cannot update node status in netbox")
	}
	return nil
}

//SetStatusFailed sets status to failed in netbox
func (n *Node) setStatusFailed(comments string) (err error) {
	comments = "temper failed: " + comments
	p, err := n.updateNodeInfo(n.Name, models.WritableDeviceWithConfigContext{
		Status:   models.DeviceWithConfigContextStatusValueFailed,
		Comments: comments,
	})
	if err != nil {
		return
	}
	if *p.Payload.Status.Value != models.DeviceWithConfigContextStatusValueFailed {
		return fmt.Errorf("cannot update node status in netbox")
	}
	return
}

//LoadInterfaces loads additional node interface info
func (n *Node) loadInterfaces() (err error) {
	n.log.Debug("calling netbox api to load node interfaces")
	in, err := n.getInterfaces()
	if err != nil {
		return
	}
	for _, in := range in {
		if in.ConnectionStatus == nil || !*in.ConnectionStatus.Value {
			continue
		}
		if !strings.Contains(*in.Name, "PCI") && !strings.Contains(*in.Name, "LCI") {
			continue
		}

		conn, ok := in.ConnectedEndpoint.(map[string]interface{})
		if !ok {
			return fmt.Errorf("no device connection info")
		}
		device, ok := conn["device"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("no device connection info")
		}

		ip, err := n.getInterfaceIP(fmt.Sprintf("%v", device["id"]))
		if err != nil {
			n.log.Error(err)
			continue
		}

		intf, ok := n.Interfaces[*in.Name]
		if !ok {
			n.log.Infof("redfish missing interface: %s", *in.Name)
			intf = NodeInterface{}
		}

		intf.Connection = fmt.Sprintf("%v", device["name"])
		intf.ConnectionIP = ip.String()
		intf.Port = fmt.Sprintf("%v", conn["name"])
		n.Interfaces[*in.Name] = intf
		if in.MacAddress == nil {
			n.log.Infof("netbox interface %s no mac", *in.Name)
			continue
		}

		if *in.MacAddress != intf.Mac {
			n.log.Infof("netbox / redfish interface %s mac mismatch: %s, %s", *in.Name, intf.Mac, *in.MacAddress)
		}

	}
	return
}

func (n *Node) updateNodeInfo(name string, data models.WritableDeviceWithConfigContext) (p *dcim.DcimDevicesUpdateOK, err error) {
	node, err := n.getDevice(nil, &name)
	if err != nil {
		return
	}
	data.DeviceType = &node.DeviceType.ID
	data.DeviceRole = &node.DeviceRole.ID
	data.Site = &node.Site.ID
	p, err = n.Clients.Netbox.Client.Dcim.DcimDevicesUpdate(&dcim.DcimDevicesUpdateParams{
		ID:      node.ID,
		Data:    &data,
		Context: context.Background(),
	}, nil)

	return
}

func (n *Node) updateNodeInterfaces() (err error) {
	intf, err := n.getInterfaces()
	if err != nil {
		return
	}
	for _, in := range intf {
		nIntf, ok := n.Interfaces[*in.Name]
		if ok {
			_, err = n.Clients.Netbox.Client.Dcim.DcimInterfacesUpdate(&dcim.DcimInterfacesUpdateParams{
				ID: in.ID,
				Data: &models.WritableInterface{
					MacAddress:  &nIntf.Mac,
					Name:        in.Name,
					Type:        in.Type.Value,
					TaggedVlans: []int64{},
					Device:      &in.Device.ID,
				},
				Context: context.Background(),
			}, nil)
		}
	}
	return
}

func (n *Node) getInterfaces() (in []*models.Interface, err error) {
	l, err := n.Clients.Netbox.Client.Dcim.DcimInterfacesList(&dcim.DcimInterfacesListParams{
		Device:  &n.Name,
		Context: context.Background(),
	}, nil)
	if err != nil {
		return
	}
	in = l.Payload.Results
	if len(in) == 0 {
		return in, fmt.Errorf("could not find interfaces for node with name %s", n.Name)
	}
	return
}

func (n *Node) getInterfaceIP(id string) (ip net.IP, err error) {
	d, err := n.getDevice(&id, nil)
	if d.PrimaryIp4 == nil {
		return ip, fmt.Errorf("no ip available for switch %s", id)
	}
	ip, _, err = net.ParseCIDR(*d.PrimaryIp4.Address)
	return
}

func (n *Node) getDevice(id, name *string) (node *models.DeviceWithConfigContext, err error) {
	param := dcim.DcimDevicesListParams{
		Context: context.Background(),
	}
	if id != nil {
		param.ID = id
	} else {
		param.Name = name
	}
	l, err := n.Clients.Netbox.Client.Dcim.DcimDevicesList(&param, nil)
	if err != nil {
		return
	}
	if len(l.Payload.Results) == 0 {
		return node, fmt.Errorf("no device found")
	}
	return l.Payload.Results[0], nil
}
