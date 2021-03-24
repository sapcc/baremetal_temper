package provision

import (
	"github.com/sapcc/baremetal_temper/pkg/clients"
	"github.com/sapcc/baremetal_temper/pkg/config"
	"github.com/sapcc/baremetal_temper/pkg/model"

	log "github.com/sirupsen/logrus"
)

type Provisioner struct {
	Node            model.Node
	clientOpenstack *clients.Client
	clientRedfish   *clients.RedfishClient
	clientInspector *clients.InspectorClient
	clientNetbox    *clients.NetboxClient
}

func NewProvisioner(node model.Node, cfg config.Config) (*Provisioner, error) {
	ctxLogger := log.WithFields(log.Fields{
		"node": node.Name,
	})
	openstackClient, err := clients.NewClient(cfg, ctxLogger)
	if err != nil {
		return nil, err
	}
	clientRedfish := clients.NewRedfishClient(cfg, ctxLogger)
	clientInspector := clients.NewInspectorClient(cfg, ctxLogger)
	clientNetbox, err := clients.NewNetboxClient(cfg, ctxLogger)
	if err != nil {
		return nil, err
	}
	return &Provisioner{node, openstackClient, clientRedfish, clientInspector, clientNetbox}, nil
}
