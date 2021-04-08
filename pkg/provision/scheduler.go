package provision

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"sync"
	"time"

	"github.com/netbox-community/go-netbox/netbox/models"
	"github.com/sapcc/baremetal_temper/pkg/clients"
	"github.com/sapcc/baremetal_temper/pkg/config"
	"github.com/sapcc/baremetal_temper/pkg/diagnostics"
	"github.com/sapcc/baremetal_temper/pkg/model"
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
	opts            config.Options
	provisoners     map[string]*Provisioner
	erroHandler     ErrorHandler
	ctx             context.Context
	nodesInProgress map[string]struct{}
	log             *log.Entry
	sync.RWMutex
}

// NewScheduler New Redfish Instance
func NewScheduler(ctx context.Context, cfg config.Config, opts config.Options) (s Scheduler, err error) {
	p, err := NewProvisioner(model.Node{}, cfg)
	if err != nil {
		return
	}
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: false,
	})
	ctxLogger := log.WithFields(log.Fields{
		"scheduler": "temper",
	})

	s = Scheduler{
		cfg:             cfg,
		provisoners:     make(map[string]*Provisioner),
		erroHandler:     NewErrorHandler(ctx, p),
		ctx:             ctx,
		nodesInProgress: make(map[string]struct{}),
		log:             ctxLogger,
		opts:            opts,
	}
	return
}

//Start ...
func (r *Scheduler) Start(d time.Duration) {
	ticker := time.NewTicker(d)
	defer ticker.Stop()

loop:
	for {
		r.log.Debug("starting temper loop...")
		nodes, err := r.loadNodes()
		if err != nil {
			r.erroHandler.Errors <- err
			continue
		}
		for _, node := range nodes {
			go r.temper(node)
		}
		select {
		case <-ticker.C:
			continue
		case <-r.ctx.Done():
			break loop
		}
	}
}

func (r *Scheduler) temper(node model.Node) {
	p, err := r.getProvisioner(node)
	if err != nil {
		r.erroHandler.Errors <- err
		return
	}
	r.log.Infof("tempering node %s", p.Node.Name)
	r.Lock()
	if _, ok := r.nodesInProgress[node.Name]; ok {
		r.Unlock()
		r.log.Infof("node %s is already being tempered", node.Name)
		return
	}
	r.nodesInProgress[node.Name] = struct{}{}
	r.Unlock()
	if err = p.clientNetbox.LoadIpamAddresses(&node); err != nil {
		r.erroHandler.Errors <- err
		return
	}
	p.clientRedfish.SetEndpoint(&node)
	t, err := r.getTasks(node)
	if err != nil {
		r.erroHandler.Errors <- err
		return
	}
	r.execTasks(t, &node)
}

func (r *Scheduler) execTasks(fns []func(n *model.Node) error, n *model.Node) (err error) {
	p, err := r.getProvisioner(*n)
	for _, fn := range fns {
		if err = fn(n); err != nil {
			if _, ok := err.(*clients.NodeAlreadyExists); ok {
				r.log.Infof("Node %s already exists, nothing to temper", p.Node.Name)
				break
			}
			r.erroHandler.Errors <- &SchedulerError{
				Err:  err.Error(),
				Node: n,
			}
			break
		}
	}
	r.log.Infof("finished tempering node: %s", p.Node.Name)
	r.Lock()
	delete(r.nodesInProgress, p.Node.Name)
	defer r.Unlock()
	return
}

func (r *Scheduler) loadNodes() (nodes []model.Node, err error) {
	targets := make([]NetboxDiscovery, 0)

	if r.cfg.NetboxNodesPath == "" {
		var n *clients.NetboxClient
		n, err = clients.NewNetboxClient(r.cfg, r.log)
		if err != nil {
			return
		}
		var pNodes []*models.DeviceWithConfigContext
		pNodes, err = n.LoadPlannedNodes(r.cfg)
		if err != nil {
			return
		}

		for _, n := range pNodes {
			nodes = append(nodes, model.Node{
				Name: *n.Name,
			})
		}
		return
	}
	d, err := ioutil.ReadFile(r.cfg.NetboxNodesPath)
	if err != nil {
		return
	}
	if err = json.Unmarshal(d, &targets); err != nil {
		return
	}

	for _, t := range targets {
		nodes = append(nodes, model.Node{
			Name: t.Labels["server_name"],
		})
	}

	return
}

func (r *Scheduler) getProvisioner(node model.Node) (p *Provisioner, err error) {
	p, ok := r.provisoners[node.Name]
	if ok {
		return
	}
	p, err = NewProvisioner(node, r.cfg)
	if err == nil {
		r.Lock()
		r.provisoners[node.Name] = p
		r.Unlock()
	}
	return
}

func (r *Scheduler) getTasks(n model.Node) (tasks []func(n *model.Node) error, err error) {

	ctxLogger := log.WithFields(log.Fields{
		"node": n.Name,
	})

	p, err := r.getProvisioner(n)
	if err != nil {
		return
	}
	tasks = make([]func(n *model.Node) error, 0)
	tasks = append(tasks,
		p.clientOpenstack.CreateDNSRecords,
		p.clientRedfish.LoadInventory,
		p.clientNetbox.LoadInterfaces,
	)
	if r.opts.Baremetal {
		if baremetal, err := p.clientOpenstack.ServiceEnabled("ironic"); err == nil && baremetal {
			tasks = append(tasks, p.clientOpenstack.GetBaremetalTasks()...)
		}
	}
	if r.opts.Diagnostics {
		d, err := diagnostics.GetTasks(n, *p.clientRedfish.ClientConfig, r.cfg, ctxLogger)
		if err != nil {
			return d, err
		}
		tasks = append(tasks, d...)
	}

	return
}
