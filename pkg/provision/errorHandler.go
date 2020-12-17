package provision

import (
	"context"

	"github.com/sapcc/ironic_temper/pkg/model"
	log "github.com/sirupsen/logrus"
)

type SchedulerError struct {
	Err  string
	Node *model.IronicNode
}

func (n *SchedulerError) Error() string {
	return n.Err
}

type ErrorHandler struct {
	Errors chan error
	ctx    context.Context
	p      *Provisioner
}

func NewErrorHandler(ctx context.Context, p *Provisioner) (e ErrorHandler) {
	errors := make(chan error, 0)
	e.Errors = errors
	e.ctx = ctx
	e.p = p
	go e.initHandler()
	return e
}

func (e ErrorHandler) initHandler() {
	go func() {
		select {
		case err := <-e.Errors:
			if serr, ok := err.(*SchedulerError); ok {
				log.Errorf("error tempering node %s. err: %s", serr.Node.UUID, serr.Err)
				if serr.Node.InstanceUUID != "" {
					if err = e.p.clientOpenstack.DeleteTestInstance(serr.Node); err != nil {
						log.Error("cannot delete compute instance %s. err: %s", serr.Node.InstanceUUID, err.Error())
					}
				}
				if err = e.p.clientOpenstack.DeleteNode(serr.Node); err != nil {
					log.Errorf("cannot delete node %s. err: %s", serr.Node.Name, err.Error())
				}
				if err = e.p.clientNetbox.SetNodeStatusFailed(serr.Node); err != nil {
					log.Errorf("cannot set node %s status in netbox. err: %s", serr.Node.Name, err.Error())
				}
			} else {
				log.Error(err.Error())
			}
		case <-e.ctx.Done():
			return
		}
	}()
}
