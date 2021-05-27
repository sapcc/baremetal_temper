package clients

import (
	"context"
	"fmt"

	runtimeclient "github.com/go-openapi/runtime/client"
	netboxclient "github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/dcim"
	"github.com/sapcc/baremetal_temper/pkg/config"
	log "github.com/sirupsen/logrus"
)

//NetboxClient is ..
type Netbox struct {
	Client *netboxclient.NetBoxAPI
	log    *log.Entry
}

//NewNetboxClient creates netbox client instance
func NewNetbox(cfg config.Config, ctxLogger *log.Entry) (n *Netbox, err error) {
	tlsClient, err := runtimeclient.TLSClient(runtimeclient.TLSClientOptions{InsecureSkipVerify: true})
	if err != nil {
		return
	}

	transport := runtimeclient.NewWithClient(cfg.Netbox.Host, netboxclient.DefaultBasePath, []string{"https"}, tlsClient)
	transport.DefaultAuthentication = runtimeclient.APIKeyAuth("Authorization", "header", fmt.Sprintf("Token %v", cfg.Netbox.Token))
	n = &Netbox{
		Client: netboxclient.New(transport, nil),
		log:    ctxLogger,
	}
	return
}

//LoadPlannedNodes loads all nodes in status planned
func (n *Netbox) LoadPlannedNodes(query *string, region *string) (nodes []string, err error) {
	nodes = make([]string, 0)
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
	for _, n := range l.Payload.Results {
		nodes = append(nodes, *n.Name)
	}
	return
}
