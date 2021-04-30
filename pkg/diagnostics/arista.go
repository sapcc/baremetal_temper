package diagnostics

import (
	"fmt"
	"strings"

	"github.com/aristanetworks/goeapi"
	"github.com/aristanetworks/goeapi/module"
	"github.com/sapcc/baremetal_temper/pkg/config"
	"github.com/sapcc/baremetal_temper/pkg/model"
	log "github.com/sirupsen/logrus"
)

type AristaClient struct {
	cfg config.Config
	log *log.Entry
}

func (a AristaClient) Run(n *model.Node) (err error) {
	foundAllNeighbors := true
	cfg := a.cfg.Arista
	for _, i := range n.Interfaces {
		if !strings.Contains(i.Connection, "asw") {
			continue
		}
		a.log.Debug("calling arista api for node cable check")
		host := fmt.Sprintf("%s.%s", i.Connection, a.cfg.Domain)
		c, err := goeapi.Connect(cfg.Transport, host, cfg.User, cfg.Password, cfg.Port)
		if err != nil {
			return err
		}
		s := module.Show(c)
		lldp := s.ShowLLDPNeighbors()
		foundNeighbor := false
		for _, ln := range lldp.LLDPNeighbors {
			//244a.979a.b76b
			//24:4a:97:9a:b7:6b
			if strings.ToLower(strings.ReplaceAll(i.Mac, ":", "")) == strings.ReplaceAll(ln.NeighborPort, ".", "") {
				foundNeighbor = true
				break
			}
		}
		if !foundNeighbor {
			foundAllNeighbors = false
		}
	}

	if !foundAllNeighbors {
		return fmt.Errorf("cable check not successful")
	}

	return
}
