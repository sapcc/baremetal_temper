package provision

import (
	"context"

	log "github.com/sirupsen/logrus"
)

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
			log.Error(err.Error())
		case <-e.ctx.Done():
			return
		}
	}()
}
