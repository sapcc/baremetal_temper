package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/netbox-community/go-netbox/netbox/models"
	"github.com/sapcc/baremetal_temper/pkg/clients"
	"github.com/sapcc/baremetal_temper/pkg/config"
	"github.com/sapcc/baremetal_temper/pkg/model"
	"github.com/sapcc/baremetal_temper/pkg/server"
	"github.com/sapcc/baremetal_temper/pkg/temper"
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
	erroHandler     ErrorHandler
	ctx             context.Context
	nodesInProgress map[string]struct{}
	log             *log.Entry
	server          *server.Handler
	tp              *temper.Temper
	sync.RWMutex
}

// New Scheduler Instance
func New(ctx context.Context, cfg config.Config, opts config.Options) (s Scheduler, err error) {
	t := temper.New(cfg)
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: false,
	})
	ctxLogger := log.WithFields(log.Fields{
		"scheduler": "temper",
	})

	s = Scheduler{
		cfg:             cfg,
		erroHandler:     NewErrorHandler(ctx, t),
		ctx:             ctx,
		nodesInProgress: make(map[string]struct{}),
		log:             ctxLogger,
		opts:            opts,
		tp:              t,
	}
	return
}

//Start ...
func (r *Scheduler) Start(d time.Duration) {
	ticker := time.NewTicker(d)
	defer ticker.Stop()

	if r.opts.RedfishEvents {
		mux := http.NewServeMux()
		r.server = server.New(r.cfg, r.log)
		go http.ListenAndServe(":9090", mux)
		go r.eventLoop()
	}

loop:
	for {
		r.log.Debug("scheduling temper...")
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
	var err error
	r.log.Infof("tempering node %s", node.Name)
	r.Lock()
	if _, ok := r.nodesInProgress[node.Name]; ok {
		r.Unlock()
		r.log.Infof("node %s is already being tempered", node.Name)
		return
	}
	r.nodesInProgress[node.Name] = struct{}{}
	r.Unlock()
	node, err = r.tp.LoadNodeInfos(node.Name)
	if err != nil {
		r.erroHandler.Errors <- err
		return
	}
	t, err := r.tp.GetAllTemperTasks(node.Name, r.opts.Diagnostics, r.opts.Baremetal, r.opts.RedfishEvents, true)
	if err != nil {
		r.erroHandler.Errors <- err
		return
	}
	if err = r.tp.TemperNode(&node, t); err != nil {
		r.erroHandler.Errors <- &SchedulerError{
			Err:  err.Error(),
			Node: &node,
		}
	}
	r.log.Infof("finished tempering node: %s", node.Name)
	r.Lock()
	delete(r.nodesInProgress, node.Name)
	defer r.Unlock()
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

func (r *Scheduler) eventLoop() {
	for {
		select {
		case e := <-r.server.Events:
			fmt.Println(e.Name)
		case <-r.ctx.Done():
			return
		}
	}
}
