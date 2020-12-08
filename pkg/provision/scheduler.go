package provision

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/sapcc/ironic_temper/pkg/clients"
	"github.com/sapcc/ironic_temper/pkg/config"
	"github.com/sapcc/ironic_temper/pkg/model"
)

// NetboxDiscovery is ...
type NetboxDiscovery struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}

// Redfish is ...
type Scheduler struct {
	cfg         config.Config
	provisoners map[string]*Provisioner
}

// New Redfish Instance
func NewScheduler(cfg config.Config) Scheduler {
	r := Scheduler{
		cfg:         cfg,
		provisoners: make(map[string]*Provisioner),
	}
	return r
}

// Start ...
func (r Scheduler) Start(ctx context.Context, errors chan<- error) {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

loop:
	for {
		nodes, err := r.loadNodes()
		if err != nil {
			errors <- err
			continue
		}
		for _, node := range nodes {
			p, err := r.getProvisioner(node)
			if err != nil {
				errors <- err
				continue
			}
			bmc, err := p.clientRedfish.LoadRedfishInfo(node.IP)
			if err != nil {
				// fail to load data with redfish client
				errors <- err
				continue
			}
			// create ironic node with insepctor
			err = p.clientInspector.CreateIronicNode(bmc, &p.ironicNode)
			if err != nil {
				// fail to create ironic node
				errors <- err
				continue
			}
			if err = p.CheckIronicNodeExists(); err != nil {
				// fail to create ironic node
				errors <- err
				continue
			}
			if err = p.clientIronic.SetNodeName(p.ironicNode.UUID, p.ironicNode.Name); err != nil {
				// fail to update ironic node name
				errors <- err
				continue
			}
			if err = p.clientIronic.PowerNodeOn(p.ironicNode.UUID); err != nil {
				// fail to power on ironic node
				errors <- err
				continue
			}
		}
		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			break loop
		}
	}
}

func (r Scheduler) loadNodes() (nodes []model.IronicNode, err error) {
	d, err := ioutil.ReadFile(r.cfg.NetboxNodesPath)
	if err != nil {
		return
	}

	targets := make([]NetboxDiscovery, 0)
	if err = json.Unmarshal(d, &targets); err != nil {
		return
	}

	for _, t := range targets {
		nodeIP := t.Targets[0]
		nodeName := t.Labels["server_name"]
		node := model.IronicNode{
			IP:     nodeIP,
			Name:   nodeName,
			Region: r.cfg.OsRegion,
		}
		nodes = append(nodes, node)
	}

	return
}

func (r Scheduler) getProvisioner(node model.IronicNode) (p *Provisioner, err error) {
	p, ok := r.provisoners[node.Name]
	if ok {
		return
	}
	p, err = NewProvisioner(node, r.cfg)
	if err == nil {
		r.provisoners[node.Name] = p
	}
	return
}

func (r Scheduler) updateNetbox(d clients.InspectorCallbackData) {
	// update provision_state
}
