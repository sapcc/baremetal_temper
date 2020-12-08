package provision

import (
	"fmt"

	"github.com/sapcc/ironic_temper/pkg/clients"
	"github.com/sapcc/ironic_temper/pkg/config"
	"github.com/sapcc/ironic_temper/pkg/model"
)

type Provisioner struct {
	ironicNode      model.IronicNode
	clientIronic    *clients.Client
	clientRedfish   clients.RedfishClient
	clientInspector clients.InspectorClient
}

type NodeNotFoundError struct {
	Err string
}

func (n *NodeNotFoundError) Error() string {
	return n.Err
}

func NewProvisioner(node model.IronicNode, cfg config.Config) (*Provisioner, error) {
	clientIronic, err := clients.NewClient(node.Region, cfg.IronicAuth)
	if err != nil {
		return nil, err
	}
	clientRedfish := clients.RedfishClient{Host: node.IP, User: cfg.Redfish.User, Password: cfg.Redfish.Password}
	clientInspector := clients.InspectorClient{Host: cfg.Inspector.Host}
	return &Provisioner{node, clientIronic, clientRedfish, clientInspector}, nil
}

func (p *Provisioner) CheckIronicNodeExists() error {
	if p.ironicNode.UUID != "" {
		return nil
	}
	_, err := p.clientIronic.GetNodeByUUID(p.ironicNode.UUID)
	if err != nil {
		return &NodeNotFoundError{
			Err: fmt.Sprintf("could not find node %s", p.ironicNode.UUID),
		}
	}
	return nil
}
