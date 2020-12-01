package redfish

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/sapcc/ironic_temper/pkg/config"
	"github.com/sapcc/ironic_temper/pkg/ironic"
)

// NetboxDiscovery is ...
type NetboxDiscovery struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}

// Redfish is ...
type Redfish struct {
	cfg         config.Config
	provisoners map[string]*Provisioner
}

// New Redfish Instance
func New(cfg config.Config) Redfish {
	r := Redfish{
		cfg:         cfg,
		provisoners: make(map[string]*Provisioner),
	}
	return r
}

// Start ...
func (r Redfish) Start(ctx context.Context, errors chan<- error) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

loop:
	for {
		nodes, err := r.loadNodes()
		if err != nil {
			errors <- err
			continue
		}
		for _, node := range nodes {
			// var p *Provisioner
			p, ok := r.provisoners[node.Name]
			if !ok {
				var err error
				p, err = NewProvisioner(node, r.cfg.IronicInspectorHost)
				if err != nil {
					// fail to create provisioner
					errors <- err
					continue
				}
				r.provisoners[node.Name] = p
			}
			bmc, err := node.LoadRedfishInfo()
			if err != nil {
				// fail to load data with redfish client
				errors <- err
				continue
			}
			if ok, err := p.CheckIronicNodeExists(); err != nil {
				// fail to check ironic node
				errors <- err
				continue
			} else {
				if ok {
					// ironic node exists
					errors <- fmt.Errorf("Node %s exist", node.Name)
					continue
				}
			}

			// create ironic node with insepctor
			if false {
				p.CreateIronicNodeWithInspector(&bmc)
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

func (r Redfish) loadNodes() (nodes []ironic.Node, err error) {
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
		node := ironic.Node{
			IP:             nodeIP,
			Name:           nodeName,
			Region:         r.cfg.OsRegion,
			IronicUser:     r.cfg.IronicUser,
			IronicPassword: r.cfg.IronicPassword,
		}
		nodes = append(nodes, node)
	}

	return
}

func (r Redfish) updateNetbox(i ironic.InspectorCallbackData) {
	// update provision_state
}
