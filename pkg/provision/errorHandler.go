package provision

import (
	"context"

	log "github.com/sirupsen/logrus"
)

type SchedulerError struct {
	Err  string
	Node string
}

func (n *SchedulerError) Error() string {
	return n.Err
}

type ErrorHandler struct {
	Errors chan error
	ctx    context.Context
}

func NewErrorHandler(ctx context.Context) (e ErrorHandler) {
	errors := make(chan error, 0)
	e.Errors = errors
	e.ctx = ctx
	go e.initHandler()
	return e
}

func (e ErrorHandler) initHandler() {
	go func() {
		select {
		case err := <-e.Errors:
			if serr, ok := err.(*SchedulerError); ok {
				log.Infof("error tempering node %s", serr.Node)
			} else {
				log.Error(err.Error())
			}
		case <-e.ctx.Done():
			return
		}
	}()
}
