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

//SetNodeStatusActive does not return error to not trigger errorhandler and cleanup of node
func (n *NetboxClient) SetNodeStatusActive(i *model.IronicNode) error {
	id, err := strconv.ParseInt(i.UUID, 10, 64)
	if err != nil {
		log.Error(err)
	}

	p, err := n.updateNode(id, models.WritableDeviceWithConfigContext{
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
	id, err := strconv.ParseInt(i.UUID, 10, 64)
	if err != nil {
		return
	}
	p, err := n.updateNode(id, models.WritableDeviceWithConfigContext{
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

func (n *NetboxClient) updateNode(id int64, data models.WritableDeviceWithConfigContext) (p *dcim.DcimDevicesUpdateOK, err error) {
	p, err = n.client.Dcim.DcimDevicesUpdate(&dcim.DcimDevicesUpdateParams{
		ID:   id,
		Data: &data,
	}, nil)

	return
}
