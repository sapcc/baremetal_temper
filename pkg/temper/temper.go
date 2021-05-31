package temper

import (
	"context"
	"sync"

	"github.com/sapcc/baremetal_temper/pkg/config"
	"github.com/sapcc/baremetal_temper/pkg/node"
)

type Temper struct {
	nodes map[string]*node.Node
	cfg   config.Config
	opts  config.Options
	ctx   context.Context
	disp  *dispatcher
	sync.RWMutex
}

func New(numWorkers int) *Temper {
	t := &Temper{
		nodes: make(map[string]*node.Node, 0),
		disp:  NewDispatcher(numWorkers),
	}
	t.disp.Start()
	return t
}

func (t *Temper) AddNodes(nodes []*node.Node) {
	t.Lock()
	for _, n := range nodes {
		t.nodes[n.Name] = n
		t.disp.Dispatch(n)
	}
	t.Unlock()
	go t.cleanup()
}

func (t *Temper) Stop() {
	t.disp.Stop()
}

func (t *Temper) cleanup() {
	t.Lock()
	defer t.Unlock()
	for i, n := range t.nodes {
		if n.Status != "staged" {
			delete(t.nodes, i)
		}
	}
}

func (t *Temper) GetNodes() map[string]*node.Node {
	t.RLock()
	defer t.RUnlock()
	return t.nodes
}
