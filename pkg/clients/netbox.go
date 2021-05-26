package clients

import (
	"context"
	"fmt"

	runtimeclient "github.com/go-openapi/runtime/client"
	netboxclient "github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/dcim"
	"github.com/netbox-community/go-netbox/netbox/models"
	"github.com/sapcc/baremetal_temper/pkg/config"
	log "github.com/sirupsen/logrus"
)

//NetboxClient is ..
type NetboxClient struct {
	Client *netboxclient.NetBoxAPI
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
		Client: netboxclient.New(transport, nil),
		log:    ctxLogger,
	}
	return
}

//LoadPlannedNodes loads all nodes in status planned
func (n *NetboxClient) LoadPlannedNodes(query *string, region *string) (nodes []*models.DeviceWithConfigContext, err error) {
	role := "server"
	status := "planned"

	param := dcim.DcimDevicesListParams{
		Context: context.Background(),
		Status:  &status,
		Role:    &role,
		Region:  region,
		Q:       query,
	}
	l, err := n.Client.Dcim.DcimDevicesList(&param, nil)
	if err != nil {
		return
	}
	nodes = l.Payload.Results
	return
}
