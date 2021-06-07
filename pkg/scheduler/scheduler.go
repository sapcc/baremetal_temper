package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/sapcc/baremetal_temper/pkg/clients"
	"github.com/sapcc/baremetal_temper/pkg/config"
	"github.com/sapcc/baremetal_temper/pkg/node"
	"github.com/sapcc/baremetal_temper/pkg/server"
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
	nc              *clients.Netbox
	sync.RWMutex
}

// New Scheduler Instance
func New(ctx context.Context, cfg config.Config, opts config.Options) (s Scheduler, err error) {
	ctxLogger := log.WithFields(log.Fields{
		"temper": "scheduler",
	})
	n, err := clients.NewNetbox(cfg, ctxLogger)
	if err != nil {
		return
	}
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: false,
	})

	s = Scheduler{
		cfg:             cfg,
		ctx:             ctx,
		nodesInProgress: make(map[string]struct{}),
		log:             ctxLogger,
		opts:            opts,
		nc:              n,
	}
	return
}

//Start ...
func (r *Scheduler) Start(d time.Duration) {
	ticker := time.NewTicker(d)
	defer ticker.Stop()

	if r.opts.RedfishEvents {
		mux := http.NewServeMux()
		r.server = server.New(r.cfg, r.log, nil)
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
	var wg sync.WaitGroup
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
	ni.AddTask("temper_dns")
	ni.Temper(true, &wg)
	r.log.Infof("finished tempering node: %s", n)
	r.Lock()
	delete(r.nodesInProgress, n)
	defer r.Unlock()
}

func (r *Scheduler) loadNodes() (nodes []string, err error) {
	targets := make([]NetboxDiscovery, 0)
	nodes = make([]string, 0)
	if r.cfg.NetboxNodesPath == "" {
		nodes, err = r.nc.LoadNodes(r.cfg.NetboxQuery, nil, &r.cfg.Region)
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
