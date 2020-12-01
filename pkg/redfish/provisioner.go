package redfish

import (
	"github.com/sapcc/ironic_temper/pkg/ironic"
)

type Provisioner struct {
	ironicNode    ironic.Node
	ironicClient  *ironic.Client
	inspectorHost string
}

func NewProvisioner(node ironic.Node, inspectorHost string) (*Provisioner, error) {
	client, err := ironic.NewClient(node.Region)
	if err != nil {
		return nil, err
	}
	return &Provisioner{node, client, inspectorHost}, nil
}

func (p *Provisioner) CheckIronicNodeExists() (bool, error) {
	if p.ironicNode.UUID != "" {
		return true, nil
	}
	uuid, err := p.ironicClient.GetNodeUUIDByName(p.ironicNode.Name)
	if err != nil {
		return false, err
	}
	p.ironicNode.UUID = uuid
	return true, nil
}

func (p *Provisioner) CreateIronicNodeWithInspector(d *ironic.InspectorCallbackData) error {
	return ironic.CreateNodeWithInspector(d, p.inspectorHost)
}
