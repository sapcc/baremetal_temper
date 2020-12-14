package provision

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/sapcc/ironic_temper/pkg/clients"
	"github.com/sapcc/ironic_temper/pkg/config"
	"github.com/sapcc/ironic_temper/pkg/model"
	log "github.com/sirupsen/logrus"
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
	erroHandler ErrorHandler
	ctx         context.Context
}

// New Redfish Instance
func NewScheduler(ctx context.Context, cfg config.Config) Scheduler {
	r := Scheduler{
		cfg:         cfg,
		provisoners: make(map[string]*Provisioner),
		erroHandler: NewErrorHandler(ctx),
		ctx:         ctx,
	}
	return r
}

// Start ...
func (r Scheduler) Start(d time.Duration) {
	ticker := time.NewTicker(d)
	defer ticker.Stop()

loop:
	for {
		log.Debug("starting temper loop...")
		nodes, err := r.loadNodes()
		if err != nil {
			r.erroHandler.Errors <- err
			continue
		}
		for _, node := range nodes {
			p, err := r.getProvisioner(node)
			if err != nil {
				r.erroHandler.Errors <- err
				continue
			}
			if err = p.clientOpenstack.CreateDNSRecordFor(&p.ironicNode); err != nil {
				// fail to load data with redfish client
				r.erroHandler.Errors <- &SchedulerError{
					Err:  err.Error(),
					Node: p.ironicNode.Name,
				}
				continue
			}
			bmc, err := p.clientRedfish.LoadRedfishInfo(p.ironicNode)
			if err != nil {
				// fail to load data with redfish client
				r.erroHandler.Errors <- &SchedulerError{
					Err:  err.Error(),
					Node: p.ironicNode.UUID,
				}
				continue
			}
			// create ironic node with insepctor
			err = p.clientInspector.CreateIronicNode(bmc, &p.ironicNode)
			if err != nil {
				if _, ok := err.(*clients.NodeAlreadyExists); ok {
					log.Infof("Node %s already exists, continue with the next node", p.ironicNode.Name)
					continue
				}
				// fail to create ironic node
				r.erroHandler.Errors <- &SchedulerError{
					Err:  err.Error(),
					Node: p.ironicNode.UUID,
				}
				continue
			}
			//p.ironicNode.UUID = "e847cdbd-2d63-4145-81a3-cef227fcb313"
			if err = p.CheckIronicNodeCreated(); err != nil {
				// fail check if ironic node was created
				r.erroHandler.Errors <- &SchedulerError{
					Err:  err.Error(),
					Node: p.ironicNode.UUID,
				}
				continue
			}
			if err = p.clientOpenstack.UpdateNode(p.ironicNode); err != nil {
				// fail to update ironic node name
				r.erroHandler.Errors <- &SchedulerError{
					Err:  err.Error(),
					Node: p.ironicNode.UUID,
				}
				continue
			}
			log.Debug("powering on node")
			if err = p.clientOpenstack.PowerNodeOn(p.ironicNode.UUID); err != nil {
				// fail to power on ironic node
				r.erroHandler.Errors <- &SchedulerError{
					Err:  err.Error(),
					Node: p.ironicNode.UUID,
				}
				continue
			}
			log.Debug("creating test node deployment")
			if err = p.clientOpenstack.CreateNodeTestDeployment(&p.ironicNode); err != nil {
				r.erroHandler.Errors <- &SchedulerError{
					Err:  err.Error(),
					Node: p.ironicNode.UUID,
				}
				continue
			}
			log.Infof("finished tempering node: %s", p.ironicNode.Name)
		}
		select {
		case <-ticker.C:
			continue
		case <-r.ctx.Done():
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
