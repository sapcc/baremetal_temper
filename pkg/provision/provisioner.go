package provision

import (
	"github.com/sapcc/ironic_temper/pkg/clients"
	"github.com/sapcc/ironic_temper/pkg/config"
)

type Provisioner struct {
	ironicNode      IronicNode
	clientIronic    *clients.Client
	clientRedfish   clients.RedfishClient
	clientInspector clients.InspectorClient
}

type IronicNode struct {
	Name   string
	IP     string
	UUID   string
	Region string
}

func NewProvisioner(node IronicNode, cfg config.Config) (*Provisioner, error) {
	clientIronic, err := clients.NewClient(node.Region)
	if err != nil {
		return nil, err
	}
	clientRedfish := clients.RedfishClient{Host: node.IP, User: cfg.IronicUser, Password: cfg.IronicPassword}
	clientInspector := clients.InspectorClient{Host: cfg.IronicInspectorHost}
	return &Provisioner{node, clientIronic, clientRedfish, clientInspector}, nil
}

func (p *Provisioner) CheckIronicNodeExists() (bool, error) {
	if p.ironicNode.UUID != "" {
		return true, nil
	}
	uuid, err := p.clientIronic.GetNodeUUIDByName(p.ironicNode.Name)
	if err != nil {
		return false, err
	}
	p.ironicNode.UUID = uuid
	return true, nil
}
