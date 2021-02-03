package provision

import (
	"github.com/sapcc/ironic_temper/pkg/clients"
	"github.com/sapcc/ironic_temper/pkg/config"
	"github.com/sapcc/ironic_temper/pkg/model"

	log "github.com/sirupsen/logrus"
)

type Provisioner struct {
	ironicNode      model.IronicNode
	clientOpenstack *clients.Client
	clientRedfish   *clients.RedfishClient
	clientInspector *clients.InspectorClient
	clientNetbox    *clients.NetboxClient
	clientArista    *clients.AristaClient
}

func NewProvisioner(node model.IronicNode, cfg config.Config) (*Provisioner, error) {
	ctxLogger := log.WithFields(log.Fields{
		"node": node.Name,
	})
	openstackClient, err := clients.NewClient(cfg, ctxLogger)
	if err != nil {
		return nil, err
	}
	clientRedfish := clients.NewRedfishClient(cfg, node.IP, ctxLogger)
	clientInspector := clients.NewInspectorClient(cfg, ctxLogger)
	clientNetbox, err := clients.NewNetboxClient(cfg, ctxLogger)
	clientArista, err := clients.NewAristaClient(cfg, ctxLogger)
	if err != nil {
		return nil, err
	}
	return &Provisioner{node, openstackClient, clientRedfish, clientInspector, clientNetbox, clientArista}, nil
}
