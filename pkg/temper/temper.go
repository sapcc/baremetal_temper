package temper

import (
	"context"
	"fmt"
	"sync"
	"time"

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

func (t *Temper) AddNode(node *node.Node) {
	t.Lock()
	n, ok := t.nodes[node.Name]
	if ok {
		if n.Status == "progress" {
			fmt.Println("node: " + node.Name + "already being tempered")
			return
		}
	}
	t.nodes[node.Name] = node
	t.disp.Dispatch(node)
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
		if n.Status == "failed" {
			if n.Updated.After(time.Now().Add(30 * time.Minute)) {
				t.disp.Dispatch(n)
			}
		}
		if n.Status == "success" {
			delete(t.nodes, i)
		}
	}
}

func (t *Temper) GetNodes() map[string]*node.Node {
	t.RLock()
	defer t.RUnlock()
	return t.nodes
}
