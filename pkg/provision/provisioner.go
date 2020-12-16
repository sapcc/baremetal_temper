package provision

import (
	"github.com/sapcc/ironic_temper/pkg/clients"
	"github.com/sapcc/ironic_temper/pkg/config"
	"github.com/sapcc/ironic_temper/pkg/model"
)

type Provisioner struct {
	ironicNode      model.IronicNode
	clientOpenstack *clients.Client
	clientRedfish   clients.RedfishClient
	clientInspector clients.InspectorClient
}

func NewProvisioner(node model.IronicNode, cfg config.Config) (*Provisioner, error) {
	clientIronic, err := clients.NewClient(node, cfg)
	if err != nil {
		return nil, err
	}
	clientRedfish := clients.RedfishClient{User: cfg.Redfish.User, Password: cfg.Redfish.Password}
	clientInspector := clients.InspectorClient{Host: cfg.Inspector.Host}
	return &Provisioner{node, clientIronic, clientRedfish, clientInspector}, nil
}
