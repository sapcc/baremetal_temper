package clients

import (
	"context"
	"fmt"
	"net"

	runtimeclient "github.com/go-openapi/runtime/client"
	netboxclient "github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/dcim"
	"github.com/netbox-community/go-netbox/netbox/models"
	"github.com/sapcc/ironic_temper/pkg/config"
	"github.com/sapcc/ironic_temper/pkg/model"
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

	transport := runtimeclient.NewWithClient(cfg.NetboxAuth.Host, netboxclient.DefaultBasePath, []string{"https"}, tlsClient)
	transport.DefaultAuthentication = runtimeclient.APIKeyAuth("Authorization", "header", fmt.Sprintf("Token %v", cfg.NetboxAuth.Token))
	n = &NetboxClient{
		client: netboxclient.New(transport, nil),
		log:    ctxLogger,
	}
	return
}

//Activate does not return error to not trigger errorhandler and cleanup of node
func (n *NetboxClient) Activate(i *model.IronicNode) error {
	p, err := n.updateNodeByName(i.Name, models.WritableDeviceWithConfigContext{
		Status: models.DeviceWithConfigContextStatusValueActive,
	})
	if err != nil {
		log.Error(err)
	}
	if *p.Payload.Status.Value != models.DeviceWithConfigContextStatusValueActive {
		log.Errorf("cannot update node status in netbox")
	}
	return nil
}

//SetNodeStatusFailed sets status to failed in netbox
func (n *NetboxClient) SetNodeStatusFailed(i *model.IronicNode) (err error) {
	p, err := n.updateNodeByName(i.Name, models.WritableDeviceWithConfigContext{
		Status: models.DeviceWithConfigContextStatusValueFailed,
	})
	if err != nil {
		return
	}
	if *p.Payload.Status.Value != models.DeviceWithConfigContextStatusValueFailed {
		return fmt.Errorf("cannot update node status in netbox")
	}
	return
}

func (n *NetboxClient) updateNodeByName(name string, data models.WritableDeviceWithConfigContext) (p *dcim.DcimDevicesUpdateOK, err error) {
	l, err := n.client.Dcim.DcimDevicesList(&dcim.DcimDevicesListParams{
		Name:    &name,
		Context: context.Background(),
	}, nil)
	if err != nil {
		return
	}
	if len(l.Payload.Results) > 1 || len(l.Payload.Results) == 0 {
		return p, fmt.Errorf("could not find node with name %s", name)
	}
	node := l.Payload.Results[0]
	p, err = n.client.Dcim.DcimDevicesUpdate(&dcim.DcimDevicesUpdateParams{
		ID:      node.ID,
		Data:    &data,
		Context: context.Background(),
	}, nil)

	return
}

func (n *NetboxClient) LoadInterfaces(i *model.IronicNode) (err error) {
	n.log.Debug("calling netbox api to load node interface")
	l, err := n.client.Dcim.DcimInterfacesList(&dcim.DcimInterfacesListParams{
		Device:  &i.Name,
		Context: context.Background(),
	}, nil)
	if err != nil {
		return
	}
	if len(l.Payload.Results) == 0 {
		return fmt.Errorf("could not find interfaces for node with name %s", i.Name)
	}

	intfs := make([]model.IronicInterface, 0)

	for _, i := range l.Payload.Results {
		if i.MacAddress == nil {
			continue
		}

		conn, ok := i.ConnectedEndpoint.(map[string]interface{})
		if !ok {
			return fmt.Errorf("no device connection info")
		}
		device, ok := conn["device"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("no device connection info")
		}

		ip, err := n.getInterfaceIP(fmt.Sprintf("%v", device["id"]))
		if err != nil {
			return err
		}
		fmt.Println(ip.String())

		intf := model.IronicInterface{
			Connection:   fmt.Sprintf("%v", device["name"]),
			ConnectionIP: ip.String(),
			Mac:          *i.MacAddress,
		}
		intfs = append(intfs, intf)
	}
	i.Interfaces = append(i.Interfaces, intfs...)
	return
}

func (n *NetboxClient) getInterfaceIP(id string) (ip net.IP, err error) {
	l, err := n.client.Dcim.DcimDevicesList(&dcim.DcimDevicesListParams{
		ID:      &id,
		Context: context.Background(),
	}, nil)
	if err != nil {
		return
	}
	if len(l.Payload.Results) == 0 {
		return ip, fmt.Errorf("no device found")
	}
	ip, _, err = net.ParseCIDR(*l.Payload.Results[0].PrimaryIp4.Address)
	return
}
