package temper

import (
	"context"
	"sync"

	"github.com/sapcc/baremetal_temper/pkg/clients"
	"github.com/sapcc/baremetal_temper/pkg/config"
)

type Temper struct {
	cfg          config.Config
	ctx          context.Context
	netbox       *clients.Netbox
	netboxStatus bool
	sync.RWMutex
}
