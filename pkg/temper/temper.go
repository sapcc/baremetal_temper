package temper

import (
	"context"
	"sync"
	"time"

	"github.com/sapcc/baremetal_temper/pkg/clients"
	"github.com/sapcc/baremetal_temper/pkg/config"
	"github.com/sapcc/baremetal_temper/pkg/diagnostics"
	"github.com/sapcc/baremetal_temper/pkg/model"
	log "github.com/sirupsen/logrus"
)

type Temper struct {
	cfg          config.Config
	clients      map[string]apiClients
	Errors       chan error
	ctx          context.Context
	netboxStatus bool
	sync.RWMutex
}

type TemperError struct {
	Err  string
	Node *model.Node
}

func (n *TemperError) Error() string {
	return n.Err
}

type apiClients struct {
	Openstack *clients.OpenstackClient
	Redfish   *clients.RedfishClient
	Netbox    *clients.NetboxClient
}

func New(cfg config.Config, ctx context.Context, setNetboxStatus bool) (t *Temper) {
	t = &Temper{cfg: cfg, clients: make(map[string]apiClients), ctx: ctx, netboxStatus: setNetboxStatus}
	go t.initErrorHandler()
	return
}

func (t *Temper) GetClients(node string) (c apiClients, err error) {
	c, ok := t.clients[node]
	if ok {
		return
	}
	c, err = t.createClients(node)
	if err != nil {
		return
	}
	t.Lock()
	t.clients[node] = c
	t.Unlock()

	return
}

func (t *Temper) createClients(node string) (c apiClients, err error) {
	ctxLogger := log.WithFields(log.Fields{
		"node": node,
	})
	if t.cfg.Openstack.Url != "" {
		c.Openstack, err = clients.NewClient(t.cfg, ctxLogger)
		if err != nil {
			return
		}
	}
	if t.cfg.Redfish.User != "" {
		c.Redfish = clients.NewRedfishClient(t.cfg, ctxLogger)
	}
	if t.cfg.Netbox.Token != "" {
		c.Netbox, err = clients.NewNetboxClient(t.cfg, ctxLogger)
		if err != nil {
			return
		}
	}
	return
}

func (t *Temper) LoadNodeInfos(node string) (n model.Node, err error) {
	c, err := t.GetClients(node)
	if err != nil {
		return
	}
	n.Name = node
	if err = c.Netbox.LoadIpamAddresses(&n); err != nil {
		return
	}
	c.Redfish.SetEndpoint(&n)
	if err = c.Redfish.LoadInventory(&n); err != nil {
		return
	}
	if err = c.Netbox.LoadInterfaces(&n); err != nil {
		return
	}
	return
}

func (t *Temper) TemperNode(n *model.Node, tasks []func(n *model.Node) error) (err error) {
	for _, task := range tasks {
		if err = task(n); err != nil {
			if _, ok := err.(*clients.NodeAlreadyExists); ok {
				log.Infof("Node %s already exists, nothing to temper", n.Name)
				break
			}
			t.Errors <- &TemperError{
				Err:  err.Error(),
				Node: n,
			}
			return err
		}
		if t.netboxStatus {
			c, err := t.GetClients(n.Name)
			if err != nil {
				log.Error(err)
			}
			if err = c.Netbox.SetStatusStaged(n); err != nil {
				log.Error(err)
			}
		}
	}
	return
}

func (t *Temper) GetAllTemperTasks(node string, diag bool, bm bool, events bool, image bool) (tasks []func(n *model.Node) error, err error) {
	c, err := t.GetClients(node)
	if err != nil {
		return
	}
	ctxLogger := log.WithFields(log.Fields{
		"node": node,
	})
	tasks = make([]func(n *model.Node) error, 0)
	tasks = append(tasks, c.Openstack.CreateDNSRecords)
	if diag {
		d, err := diagnostics.GetHardwareCheckTasks(*c.Redfish.ClientConfig, t.cfg, ctxLogger)
		if err != nil {
			return d, err
		}
		tasks = append(tasks, d...)
		if t.cfg.Redfish.BootImage != nil && image {
			tasks = append(tasks, c.Redfish.BootImage, t.GetTimeoutTask(10*time.Minute))
		}
		tasks = append(tasks, diagnostics.GetCableCheckTasks(t.cfg, ctxLogger)...)
	}
	if bm {
		if baremetal, err := c.Openstack.ServiceEnabled("ironic"); err == nil && baremetal {
			tasks = append(tasks, c.Openstack.Create()...)
			tasks = append(tasks, c.Openstack.DeploymentTest()...)
			tasks = append(tasks, c.Openstack.Prepare)
		}
	}
	if events {
		tasks = append(tasks, c.Redfish.DeleteEventSubscription)
	}

	tasks = append(tasks, c.Netbox.Update)
	return
}

func (t *Temper) GetTimeoutTask(d time.Duration) (task func(n *model.Node) error) {
	task = func(n *model.Node) (err error) {
		time.Sleep(d)
		return
	}
	return
}

func (t *Temper) initErrorHandler() {
	for {
		select {
		case err := <-t.Errors:
			if serr, ok := err.(*TemperError); ok {
				log.Errorf("error tempering node %s. err: %s", serr.Node.Name, serr.Err)
				c, err := t.GetClients(serr.Node.Name)
				if serr.Node.InstanceUUID != "" {
					if err = c.Openstack.DeleteTestInstance(serr.Node); err != nil {
						log.Error("cannot delete compute instance %s. err: %s", serr.Node.InstanceUUID, err.Error())
					}
				}
				if err = c.Openstack.DeleteNode(serr.Node); err != nil {
					log.Errorf("cannot delete node %s. err: %s", serr.Node.Name, err.Error())
				}
				if t.netboxStatus {
					if err = c.Netbox.SetStatusFailed(serr.Node, serr.Err); err != nil {
						log.Errorf("cannot set node %s status in netbox. err: %s", serr.Node.Name, err.Error())
					}
				}

			} else {
				log.Error(err.Error())
			}
		case <-t.ctx.Done():
			return
		}
	}
}
