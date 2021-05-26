package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/sapcc/baremetal_temper/pkg/config"
	"github.com/sapcc/baremetal_temper/pkg/node"
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
	ctx             context.Context
	nodesInProgress map[string]struct{}
	log             *log.Entry
	server          *server.Handler
	tp              *temper.Temper
	sync.RWMutex
}

// New Scheduler Instance
func New(ctx context.Context, cfg config.Config, opts config.Options) (s Scheduler, err error) {
	t := temper.New(cfg, ctx, true)
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: false,
	})
	ctxLogger := log.WithFields(log.Fields{
		"scheduler": "temper",
	})

	s = Scheduler{
		cfg:             cfg,
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
			r.log.Error(err)
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

func (r *Scheduler) temper(n string) {
	var err error
	r.log.Infof("tempering node %s", n)
	r.Lock()
	if _, ok := r.nodesInProgress[n]; ok {
		r.Unlock()
		r.log.Infof("node %s is already being tempered", n)
		return
	}
	r.nodesInProgress[n] = struct{}{}
	r.Unlock()
	ni, err := node.New(n, r.cfg)
	if err != nil {
		r.log.Error(err)
		return
	}
	r.tp.SetAllTemperTasks(ni, r.opts.Diagnostics, r.opts.Baremetal, r.opts.RedfishEvents, true)
	r.tp.TemperNode(ni, true)
	r.log.Infof("finished tempering node: %s", n)
	r.Lock()
	delete(r.nodesInProgress, n)
	defer r.Unlock()
}

func (r *Scheduler) loadNodes() (nodes []string, err error) {
	targets := make([]NetboxDiscovery, 0)
	nodes = make([]string, 0)
	if r.cfg.NetboxNodesPath == "" {
		nodes, err = r.tp.LoadPlannedNodes(nil)
		if err != nil {
			return
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
		nodes = append(nodes, t.Labels["server_name"])
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
