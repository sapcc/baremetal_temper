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

package netbox

import (
	"context"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"

	"github.com/netbox-community/go-netbox/netbox/client/dcim"
	"github.com/netbox-community/go-netbox/netbox/client/ipam"
	"github.com/netbox-community/go-netbox/netbox/models"
	"github.com/sapcc/baremetal_temper/pkg/clients"
	"github.com/sapcc/baremetal_temper/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/stmcginnis/gofish/redfish"
)

type Netbox struct {
	Data   *Data
	client *clients.Netbox
	node   string
}

type Data struct {
	Device        *models.DeviceWithConfigContext
	RemoteIP      string
	PrimaryIP     string
	DNSName       string
	IpamAddresses []models.IPAddress
	Interfaces    []NodeInterface `json:"-"`
}

type NodeInterface struct {
	Name           string
	RedfishName    string
	Connection     string
	ConnectionIP   string
	Port           string
	Mac            string
	PortLinkStatus redfish.PortLinkStatus
	PortNumber     int
	Nic            int
}

func New(node string, cfg config.Config, ctxLogger *logrus.Entry) (*Netbox, error) {
	c, err := clients.NewNetbox(cfg, ctxLogger)
	if err != nil {
		return nil, err
	}
	return &Netbox{client: c, node: node}, nil
}

func (n *Netbox) GetData() (*Data, error) {
	if n.Data != nil {
		return n.Data, nil
	}
	d, err := n.getDeviceConfig(nil, &n.node)
	n.Data = &Data{}
	n.Data.Device = d
	if err != nil {
		return n.Data, err
	}
	if err := n.loadIpamAddresses(n.node); err != nil {
		return n.Data, nil
	}
	if err := n.loadInterfaces(); err != nil {
		return n.Data, nil
	}
	return n.Data, nil
}

//LoadIpamAddresses loads all ipam addresse of a node
func (n *Netbox) loadIpamAddresses(node string) (err error) {
	//n.log.Debug("calling netbox api to load ipam Addresses")
	split := strings.Split(node, "-")
	if len(split) == 1 {
		return fmt.Errorf("wrong node name format: node[001]-[block_name]")
	}
	block := strings.Split(node, "-")[1]
	name := strings.Split(node, "-")[0]
	var limit int64
	limit = 200
	al, err := n.client.Client.Ipam.IpamIPAddressesList(&ipam.IpamIPAddressesListParams{
		Q:       &block,
		Context: context.Background(),
		Limit:   &limit,
	}, nil)
	if err != nil {
		return
	}
	addr := al.Payload.Results
	for _, a := range addr {
		if strings.Contains(a.DNSName, node) {
			ip, _, err := net.ParseCIDR(*a.Address)
			if err != nil {
				return err
			}
			n.Data.PrimaryIP = ip.String()
		}

		if strings.Contains(a.DNSName, name+"r") {
			ip, _, err := net.ParseCIDR(*a.Address)
			if err != nil {
				return err
			}
			n.Data.RemoteIP = ip.String()
			n.Data.DNSName = a.DNSName
		}
		if strings.Contains(a.DNSName, name) {
			n.Data.IpamAddresses = append(n.Data.IpamAddresses, *a)
		}
	}
	return
}

func (n *Netbox) getDeviceConfig(id, name *string) (d *models.DeviceWithConfigContext, err error) {
	param := dcim.DcimDevicesListParams{
		Context: context.Background(),
	}
	if id != nil {
		param.ID = id
	} else {
		param.Name = name
	}
	l, err := n.client.Client.Dcim.DcimDevicesList(&param, nil)
	if err != nil {
		return
	}
	if len(l.Payload.Results) == 0 {
		return d, fmt.Errorf("no device found")
	}
	return l.Payload.Results[0], err
}

//Update node serial and primaryIP. Does not return error to not trigger errorhandler and cleanup of node
func (n *Netbox) Update(serialNumber string) error {
	params := models.WritableDeviceWithConfigContext{
		Serial: serialNumber,
	}
	if n.Data.PrimaryIP == "" {
		ips, err := n.client.Client.Ipam.IpamIPAddressesList(&ipam.IpamIPAddressesListParams{
			Address: &n.Data.PrimaryIP,
			Context: context.Background(),
		}, nil)
		if err != nil {
			//n.log.Error(err)
			return nil
		}
		if len(ips.Payload.Results) > 1 || len(ips.Payload.Results) == 0 {
			//n.log.Errorf("could not find ip %s", n.PrimaryIP)
			return nil
		}
		params.PrimaryIp4 = &ips.Payload.Results[0].ID
	}

	_, err := n.updateNodeInfo(params)
	if err != nil {
		//n.log.Error(err)
		return nil
	}

	if err = n.updateNodeInterfaces(); err != nil {
		//n.log.Error(err)
		return nil
	}

	return nil
}

//SetStatus does not return error to not trigger errorhandler and cleanup of node
func (n *Netbox) SetStatus(status string) (err error) {
	if status == "failed" {
		err = n.setStatusFailed()
	} else {
		err = n.setStatusStaged()
	}
	return
}

//SetStatusStaged does not return error to not trigger errorhandler and cleanup of node
func (n *Netbox) setStatusStaged() error {
	p, err := n.updateNodeInfo(models.WritableDeviceWithConfigContext{
		Status:   models.DeviceWithConfigContextStatusValueStaged,
		Comments: "temper successful",
	})
	if err != nil {
		//n.log.Error(err)
		return nil
	}
	if *p.Payload.Status.Value != models.DeviceWithConfigContextStatusValueStaged {
		//n.log.Errorf("cannot update node status in netbox")
	}
	return nil
}

//SetStatusFailed sets status to failed in netbox
func (n *Netbox) setStatusFailed() (err error) {
	p, err := n.updateNodeInfo(models.WritableDeviceWithConfigContext{
		Status:   models.DeviceWithConfigContextStatusValueFailed,
		Comments: "temper failed: check config context",
	})
	if err != nil {
		//n.log.Error(err)
		return nil
	}
	if *p.Payload.Status.Value != models.DeviceWithConfigContextStatusValueFailed {
		return fmt.Errorf("cannot update node status in netbox")
	}
	return
}

//LoadInterfaces loads additional node interface info
func (n *Netbox) loadInterfaces() (err error) {
	//n.log.Debug("calling netbox api to load node interfaces")
	n.Data.Interfaces = make([]NodeInterface, 0)
	in, err := n.getInterfaces()
	if err != nil {
		return
	}
	for _, in := range in {
		var (
			ip     net.IP
			device map[string]interface{}
			conn   map[string]interface{}
		)
		if !strings.Contains(*in.Name, "NIC") && !strings.Contains(*in.Name, "L") || strings.Contains(*in.Name, "LAG") {
			continue
		}

		if in.ConnectionStatus != nil && *in.ConnectionStatus.Value {
			conn = in.ConnectedEndpoint.(map[string]interface{})
			device = conn["device"].(map[string]interface{})
			ip, err = n.getInterfaceIP(fmt.Sprintf("%v", device["id"]))
			if err != nil {
				//n.log.Error(err)
				continue
			}
		}
		re := regexp.MustCompile("[0-9]+")
		nicPort := re.FindAllString(*in.Name, -1)
		nic := 0
		port := 0
		if len(nicPort) == 2 {
			nic, _ = strconv.Atoi(nicPort[0])
			port, _ = strconv.Atoi(nicPort[1])
		} else {
			port, _ = strconv.Atoi(nicPort[0])
		}

		intf := NodeInterface{}
		intf.Nic = nic
		intf.PortNumber = port
		if ip != nil {
			intf.Connection = fmt.Sprintf("%v", device["name"])
			intf.ConnectionIP = ip.String()
			intf.Port = fmt.Sprintf("%v", conn["name"])
		}
		intf.Name = *in.Name
		n.Data.Interfaces = append(n.Data.Interfaces, intf)
	}
	return
}

func (n *Netbox) updateNodeInfo(data models.WritableDeviceWithConfigContext) (p *dcim.DcimDevicesUpdateOK, err error) {
	data.DeviceType = &n.Data.Device.DeviceType.ID
	data.DeviceRole = &n.Data.Device.DeviceRole.ID
	data.Site = &n.Data.Device.Site.ID
	p, err = n.client.Client.Dcim.DcimDevicesUpdate(&dcim.DcimDevicesUpdateParams{
		ID:      n.Data.Device.ID,
		Data:    &data,
		Context: context.Background(),
	}, nil)

	return
}

func (n *Netbox) updateNodeInterfaces() (err error) {
	intf, err := n.getInterfaces()
	if err != nil {
		return
	}
	for _, in := range intf {
		for _, nIntf := range n.Data.Interfaces {
			if nIntf.Name == *in.Name {
				_, err = n.client.Client.Dcim.DcimInterfacesUpdate(&dcim.DcimInterfacesUpdateParams{
					ID: in.ID,
					Data: &models.WritableInterface{
						Description: nIntf.RedfishName,
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
	}
	return
}

func (n *Netbox) getInterfaces() (in []*models.Interface, err error) {
	l, err := n.client.Client.Dcim.DcimInterfacesList(&dcim.DcimInterfacesListParams{
		Device:  &n.Data.Device.DisplayName,
		Context: context.Background(),
	}, nil)
	if err != nil {
		return
	}
	in = l.Payload.Results
	if len(in) == 0 {
		return in, fmt.Errorf("could not find interfaces for node with name %s", n.Data.Device.DisplayName)
	}
	return
}

func (n *Netbox) getInterfaceIP(id string) (ip net.IP, err error) {
	d, err := n.getDeviceConfig(&id, nil)
	if d.PrimaryIp4 == nil {
		return ip, fmt.Errorf("no ip available for switch %s", id)
	}
	ip, _, err = net.ParseCIDR(*d.PrimaryIp4.Address)
	return
}
