package clients

import (
	"fmt"
	"strconv"

	runtimeclient "github.com/go-openapi/runtime/client"
	netboxclient "github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/dcim"
	"github.com/netbox-community/go-netbox/netbox/models"
	"github.com/sapcc/ironic_temper/pkg/config"
	"github.com/sapcc/ironic_temper/pkg/model"
	log "github.com/sirupsen/logrus"
)

type NetboxClient struct {
	client *netboxclient.NetBoxAPI
}

func NewNetboxClient(cfg config.Config, ctxLogger *log.Entry) (n *NetboxClient, err error) {
	tlsClient, err := runtimeclient.TLSClient(runtimeclient.TLSClientOptions{InsecureSkipVerify: true})
	if err != nil {
		return
	}

	transport := runtimeclient.NewWithClient(cfg.Netbox.Host, netboxclient.DefaultBasePath, []string{"https"}, tlsClient)
	transport.DefaultAuthentication = runtimeclient.APIKeyAuth("Authorization", "header", fmt.Sprintf("Token %v", cfg.Netbox.Token))
	n.client = netboxclient.New(transport, nil)
	return
}

func (n *NetboxClient) SetNodeStatusActive(i *model.IronicNode) (err error) {
	id, err := strconv.ParseInt(i.UUID, 10, 64)
	if err != nil {
		return
	}
	return n.setNodeStatus(id, models.DeviceWithConfigContextStatusValueActive)
}

func (n *NetboxClient) SetNodeStatusFailed(i *model.IronicNode) (err error) {
	id, err := strconv.ParseInt(i.UUID, 10, 64)
	if err != nil {
		return
	}
	return n.setNodeStatus(id, models.DeviceWithConfigContextStatusValueFailed)
}

func (n *NetboxClient) setNodeStatus(id int64, status string) (err error) {
	u, err := n.client.Dcim.DcimDevicesUpdate(&dcim.DcimDevicesUpdateParams{
		ID: id,
		Data: &models.WritableDeviceWithConfigContext{
			Status: status,
		},
	}, nil)
	if *u.Payload.Status.Value != status {
		return fmt.Errorf("cannot update node status in netbox")
	}
	return
}
