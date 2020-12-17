package provision

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"sync"
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

// Scheduler is ...
type Scheduler struct {
	cfg             config.Config
	provisoners     map[string]*Provisioner
	erroHandler     ErrorHandler
	ctx             context.Context
	nodesInProgress map[string]struct{}
	sync.RWMutex
}

// NewScheduler New Redfish Instance
func NewScheduler(ctx context.Context, cfg config.Config) (s Scheduler, err error) {
	p, err := NewProvisioner(model.IronicNode{}, cfg)
	if err != nil {
		return
	}
	s = Scheduler{
		cfg:             cfg,
		provisoners:     make(map[string]*Provisioner),
		erroHandler:     NewErrorHandler(ctx, p),
		ctx:             ctx,
		nodesInProgress: make(map[string]struct{}),
	}
	return
}

//Start ...
func (r *Scheduler) Start(d time.Duration) {
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

			r.Lock()
			if _, ok := r.nodesInProgress[node.Name]; ok {
				r.Unlock()
				log.Infof("node %s is already being tempered", node.Name)
				continue
			}
			r.nodesInProgress[node.Name] = struct{}{}
			r.Unlock()

			go r.run([]func(n *model.IronicNode) error{
				p.clientOpenstack.CreateDNSRecordFor,
				p.clientRedfish.LoadRedfishInfo,
				p.clientInspector.CreateIronicNode,
				p.clientOpenstack.CheckIronicNodeCreated,
				p.clientOpenstack.ApplyRules,
				p.clientOpenstack.ValidateNode,
				p.clientOpenstack.PowerNodeOn,
				p.clientOpenstack.ProvideNode,
				p.clientOpenstack.WaitForNovaPropagation,
				p.clientOpenstack.CreateTestInstance,
				p.clientOpenstack.DeleteTestInstance,
				p.clientOpenstack.PrepareNode,
				p.clientNetbox.SetNodeStatusActive,
			}, p)
		}
		select {
		case <-ticker.C:
			continue
		case <-r.ctx.Done():
			break loop
		}
	}
}

func (r *Scheduler) run(fns []func(n *model.IronicNode) error, p *Provisioner) (err error) {
	log.Infof("tempering node %s", p.ironicNode.Name)
	for _, fn := range fns {
		if err = fn(&p.ironicNode); err != nil {
			if _, ok := err.(*clients.NodeAlreadyExists); ok {
				log.Infof("Node %s already exists, nothing to temper", p.ironicNode.Name)
				break
			}
			r.erroHandler.Errors <- &SchedulerError{
				Err:  err.Error(),
				Node: &p.ironicNode,
			}
			break
		}
	}
	log.Infof("finished tempering node: %s", p.ironicNode.Name)
	r.Lock()
	delete(r.nodesInProgress, p.ironicNode.Name)
	defer r.Unlock()
	return
}

func (r *Scheduler) loadNodes() (nodes []model.IronicNode, err error) {
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

func (r *Scheduler) getProvisioner(node model.IronicNode) (p *Provisioner, err error) {
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
