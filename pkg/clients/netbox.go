package clients

import (
	"context"
	"fmt"
	"net"
	"strings"

	runtimeclient "github.com/go-openapi/runtime/client"
	netboxclient "github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/dcim"
	"github.com/netbox-community/go-netbox/netbox/client/ipam"
	"github.com/netbox-community/go-netbox/netbox/models"
	"github.com/sapcc/baremetal_temper/pkg/config"
	"github.com/sapcc/baremetal_temper/pkg/model"
	log "github.com/sirupsen/logrus"
)

//NetboxClient is ..
type NetboxClient struct {
	client *netboxclient.NetBoxAPI
	log    *log.Entry
}

//NewNetboxClient creates netbox client instance
func NewNetboxClient(cfg config.Config, ctxLogger *log.Entry) (n *NetboxClient, err error) {
	tlsClient, err := runtimeclient.TLSClient(runtimeclient.TLSClientOptions{InsecureSkipVerify: true})
	if err != nil {
		return
	}

	transport := runtimeclient.NewWithClient(cfg.Netbox.Host, netboxclient.DefaultBasePath, []string{"https"}, tlsClient)
	transport.DefaultAuthentication = runtimeclient.APIKeyAuth("Authorization", "header", fmt.Sprintf("Token %v", cfg.Netbox.Token))
	n = &NetboxClient{
		client: netboxclient.New(transport, nil),
		log:    ctxLogger,
	}
	return
}

//Update node serial and primaryIP. Does not return error to not trigger errorhandler and cleanup of node
func (n *NetboxClient) Update(i *model.Node) error {
	params := models.WritableDeviceWithConfigContext{
		Serial: i.InspectionData.Inventory.SystemVendor.SerialNumber,
	}
	d, err := n.getDevice(nil, &i.Name)
	if err != nil {
		log.Error(err)
		return nil
	}
	if d.PrimaryIP == nil {
		ips, err := n.client.Ipam.IpamIPAddressesList(&ipam.IpamIPAddressesListParams{
			Address: &i.PrimaryIP,
			Context: context.Background(),
		}, nil)
		if err != nil {
			log.Error(err)
			return nil
		}
		if len(ips.Payload.Results) > 1 || len(ips.Payload.Results) == 0 {
			log.Errorf("could not find ip %s", i.PrimaryIP)
			return nil
		}
		params.PrimaryIp4 = &ips.Payload.Results[0].ID
	}

	_, err = n.updateNodeInfo(i.Name, params)
	if err != nil {
		log.Error(err)
		return nil
	}

	if err = n.updateNodeInterfaces(i); err != nil {
		log.Error(err)
		return nil
	}

	return nil
}

//LoadIpamAddresses loads all ipam addresse of a node
func (n *NetboxClient) LoadIpamAddresses(i *model.Node) (err error) {
	n.log.Debug("calling netbox api to load ipam Addresses")
	block := strings.Split(i.Name, "-")[1]
	name := strings.Split(i.Name, "-")[0]
	var limit int64
	limit = 200
	al, err := n.client.Ipam.IpamIPAddressesList(&ipam.IpamIPAddressesListParams{
		Q:       &block,
		Context: context.Background(),
		Limit:   &limit,
	}, nil)
	if err != nil {
		return
	}
	addr := al.Payload.Results
	for _, a := range addr {
		if strings.Contains(a.DNSName, i.Name) {
			ip, _, err := net.ParseCIDR(*a.Address)
			if err != nil {
				return err
			}
			fmt.Println(a.DNSName, i.Name, ip.String())
			i.PrimaryIP = ip.String()
		}
		if strings.Contains(a.DNSName, name+"r") {
			ip, _, err := net.ParseCIDR(*a.Address)
			if err != nil {
				return err
			}
			i.RemoteIP = ip.String()
			i.InspectionData.Inventory.BmcAddress = a.DNSName
		}
		if strings.Contains(a.DNSName, name) {
			i.IpamAddresses = append(i.IpamAddresses, *a)
		}
	}
	return
}

//SetStatusStaged does not return error to not trigger errorhandler and cleanup of node
func (n *NetboxClient) SetStatusStaged(i *model.Node) error {
	p, err := n.updateNodeInfo(i.Name, models.WritableDeviceWithConfigContext{
		Status:   models.DeviceWithConfigContextStatusValueStaged,
		Comments: "temper successful",
	})
	if err != nil {
		log.Error(err)
		return nil
	}
	if *p.Payload.Status.Value != models.DeviceWithConfigContextStatusValueStaged {
		log.Errorf("cannot update node status in netbox")
	}
	return nil
}

//SetStatusFailed sets status to failed in netbox
func (n *NetboxClient) SetStatusFailed(i *model.Node, comments string) (err error) {
	comments = "temper failed: " + comments
	p, err := n.updateNodeInfo(i.Name, models.WritableDeviceWithConfigContext{
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
func (n *NetboxClient) LoadInterfaces(i *model.Node) (err error) {
	n.log.Debug("calling netbox api to load node interfaces")
	in, err := n.getInterfaces(i)
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
			log.Error(err)
			continue
		}

		intf, ok := i.Interfaces[*in.Name]
		if !ok {
			log.Infof("redfish missing interface: %s", *in.Name)
			intf = model.NodeInterface{}
		}

		intf.Connection = fmt.Sprintf("%v", device["name"])
		intf.ConnectionIP = ip.String()
		intf.Port = fmt.Sprintf("%v", conn["name"])
		i.Interfaces[*in.Name] = intf
		if in.MacAddress == nil {
			log.Infof("netbox interface %s no mac", *in.Name)
			continue
		}

		if *in.MacAddress != intf.Mac {
			log.Infof("netbox / redfish interface %s mac mismatch: %s, %s", *in.Name, intf.Mac, *in.MacAddress)
		}

	}
	return
}

//LoadPlannedNodes loads all nodes in status planned
func (n *NetboxClient) LoadPlannedNodes(cfg config.Config) (nodes []*models.DeviceWithConfigContext, err error) {
	role := "server"
	status := "planned"

	param := dcim.DcimDevicesListParams{
		Context: context.Background(),
		Status:  &status,
		Role:    &role,
		Region:  &cfg.Region,
		Q:       cfg.NetboxQuery,
	}
	l, err := n.client.Dcim.DcimDevicesList(&param, nil)
	if err != nil {
		return
	}
	nodes = l.Payload.Results
	return
}

func (n *NetboxClient) updateNodeInfo(name string, data models.WritableDeviceWithConfigContext) (p *dcim.DcimDevicesUpdateOK, err error) {
	node, err := n.getDevice(nil, &name)
	if err != nil {
		return
	}
	data.DeviceType = &node.DeviceType.ID
	data.DeviceRole = &node.DeviceRole.ID
	data.Site = &node.Site.ID
	p, err = n.client.Dcim.DcimDevicesUpdate(&dcim.DcimDevicesUpdateParams{
		ID:      node.ID,
		Data:    &data,
		Context: context.Background(),
	}, nil)

	return
}

func (n *NetboxClient) updateNodeInterfaces(i *model.Node) (err error) {
	intf, err := n.getInterfaces(i)
	if err != nil {
		return
	}
	for _, in := range intf {
		nIntf, ok := i.Interfaces[*in.Name]
		if ok {
			_, err = n.client.Dcim.DcimInterfacesUpdate(&dcim.DcimInterfacesUpdateParams{
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

func (n *NetboxClient) getInterfaces(i *model.Node) (in []*models.Interface, err error) {
	l, err := n.client.Dcim.DcimInterfacesList(&dcim.DcimInterfacesListParams{
		Device:  &i.Name,
		Context: context.Background(),
	}, nil)
	if err != nil {
		return
	}
	in = l.Payload.Results
	if len(in) == 0 {
		return in, fmt.Errorf("could not find interfaces for node with name %s", i.Name)
	}
	return
}

func (n *NetboxClient) getInterfaceIP(id string) (ip net.IP, err error) {
	d, err := n.getDevice(&id, nil)
	if d.PrimaryIp4 == nil {
		return ip, fmt.Errorf("no ip available for switch %s", id)
	}
	ip, _, err = net.ParseCIDR(*d.PrimaryIp4.Address)
	return
}

func (n *NetboxClient) getDevice(id, name *string) (node *models.DeviceWithConfigContext, err error) {
	param := dcim.DcimDevicesListParams{
		Context: context.Background(),
	}
	if id != nil {
		param.ID = id
	} else {
		param.Name = name
	}
	l, err := n.client.Dcim.DcimDevicesList(&param, nil)
	if err != nil {
		return
	}
	if len(l.Payload.Results) == 0 {
		return node, fmt.Errorf("no device found")
	}
	return l.Payload.Results[0], nil
}
