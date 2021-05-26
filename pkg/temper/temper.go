package temper

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/sapcc/baremetal_temper/pkg/clients"
	"github.com/sapcc/baremetal_temper/pkg/config"
	"github.com/sapcc/baremetal_temper/pkg/node"
	log "github.com/sirupsen/logrus"
)

type Temper struct {
	cfg          config.Config
	ctx          context.Context
	netbox       *clients.NetboxClient
	netboxStatus bool
	sync.RWMutex
}

type TemperError struct {
	Err  string
	Node *node.Node
}

func (n *TemperError) Error() string {
	return n.Err
}

func New(cfg config.Config, ctx context.Context, setNetboxStatus bool) (t *Temper) {
	ctxLogger := log.WithFields(log.Fields{
		"temper": "temper",
	})
	n, _ := clients.NewNetboxClient(cfg, ctxLogger)
	t = &Temper{
		cfg: cfg,
		ctx: ctx, netboxStatus: setNetboxStatus,
		netbox: n,
	}
	return
}

func (t *Temper) TemperNode(n *node.Node, netboxSts bool) (err error) {
	prios := make([]int, 0, len(n.Tasks))
	for k := range n.Tasks {
		prios = append(prios, k)
	}
	sort.Ints(prios)
	for _, k := range prios {
		fmt.Println(n.Tasks[k].Name)
	}
	for _, k := range prios {
		if err = n.Tasks[k].Exec(); err != nil {
			if _, ok := err.(*node.AlreadyExists); ok {
				log.Infof("Node %s already exists, nothing to temper", n.Name)
				break
			}
			n.Tasks[k].Error = err
			n.Status = "failed"
		}
	}
	t.cleanupHandler(n, netboxSts)
	return
}

func (t *Temper) SetAllTemperTasks(n *node.Node, diag bool, bm bool, events bool, image bool) {
	n.GetOrCreateTask(0, "create_dns_records").Exec = n.CreateDNSRecords
	if events {
		n.GetOrCreateTask(10, "create_event_sub").Exec = n.CreateEventSubscription
	}
	if diag {
		n.GetOrCreateTask(20, "hardware_check").Exec = n.RunHardwareChecks

		if t.cfg.Redfish.BootImage != nil && image {
			n.GetOrCreateTask(30, "boot_image").Exec = n.BootImage
			n.GetOrCreateTask(40, "boot_image_wait").Exec = TimeoutTask(5 * time.Minute)
		}
		n.GetOrCreateTask(50, "aci_check").Exec = n.RunACICheck
		n.GetOrCreateTask(51, "arista_check").Exec = n.RunAristaCheck
	}
	if bm {
		n.GetOrCreateTask(60, "aci_check").Exec = n.Create
		n.GetOrCreateTask(61, "aci_check").Exec = n.ApplyRules
		n.GetOrCreateTask(62, "aci_check").Exec = n.Validate
		n.GetOrCreateTask(63, "aci_check").Exec = n.Prepare
		n.GetOrCreateTask(64, "aci_check").Exec = n.WaitForNovaPropagation
		n.GetOrCreateTask(65, "aci_check").Exec = n.DeployTestInstance

	}
	if events {
		n.GetOrCreateTask(70, "delete_event_sub").Exec = n.DeleteEventSubscription
	}
	n.GetOrCreateTask(100, "update_netbox").Exec = n.Update
	return
}

func (t *Temper) LoadPlannedNodes(query *string) (nodes []string, err error) {
	nodes = make([]string, 0)
	pNodes, err := t.netbox.LoadPlannedNodes(query, &t.cfg.Region)
	if err != nil {
		return
	}
	for _, n := range pNodes {
		nodes = append(nodes, *n.Name)
	}
	return
}

func TimeoutTask(d time.Duration) func() (err error) {
	return func() (err error) {
		time.Sleep(d)
		return
	}
}

func (t *Temper) cleanupHandler(n *node.Node, netboxSts bool) {
	for _, t := range n.Tasks {
		if t.Error != nil {
			log.Errorf("error tempering node %s. task: %s err: %s", n.Name, t.Name, t.Error.Error())
		}
	}
	if n.InstanceUUID != "" {
		if err := n.DeleteTestInstance(); err != nil {
			log.Error("cannot delete compute instance %s. err: %s", n.InstanceUUID, err.Error())
		}
	}
	if err := n.DeleteNode(); err != nil {
		log.Errorf("cannot delete node %s. err: %s", n.Name, err.Error())
	}
	if netboxSts {
		if err := n.SetStatus(); err != nil {
			log.Errorf("cannot set node %s status in netbox. err: %s", n.Name, err.Error())
		}
	}

}
