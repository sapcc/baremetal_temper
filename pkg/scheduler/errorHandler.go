package scheduler

import (
	"context"

	"github.com/sapcc/baremetal_temper/pkg/model"
	"github.com/sapcc/baremetal_temper/pkg/temper"
	log "github.com/sirupsen/logrus"
)

type SchedulerError struct {
	Err  string
	Node *model.Node
}

func (n *SchedulerError) Error() string {
	return n.Err
}

type ErrorHandler struct {
	Errors chan error
	ctx    context.Context
	tp     *temper.Temper
}

func NewErrorHandler(ctx context.Context, p *temper.Temper) (e ErrorHandler) {
	errors := make(chan error)
	e.Errors = errors
	e.ctx = ctx
	e.tp = p
	go e.initHandler()
	return e
}

func (e ErrorHandler) initHandler() {
	for {
		select {
		case err := <-e.Errors:
			if serr, ok := err.(*SchedulerError); ok {
				log.Errorf("error tempering node %s. err: %s", serr.Node.Name, serr.Err)
				c, err := e.tp.GetClients(serr.Node.Name)
				if serr.Node.InstanceUUID != "" {
					if err = c.Openstack.DeleteTestInstance(serr.Node); err != nil {
						log.Error("cannot delete compute instance %s. err: %s", serr.Node.InstanceUUID, err.Error())
					}
				}
				if err = c.Openstack.DeleteNode(serr.Node); err != nil {
					log.Errorf("cannot delete node %s. err: %s", serr.Node.Name, err.Error())
				}
				if err = c.Netbox.SetStatusFailed(serr.Node, serr.Err); err != nil {
					log.Errorf("cannot set node %s status in netbox. err: %s", serr.Node.Name, err.Error())
				}
			} else {
				log.Error(err.Error())
			}
		case <-e.ctx.Done():
			return
		}
	}
}
