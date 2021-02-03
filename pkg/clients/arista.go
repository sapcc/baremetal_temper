package clients

import (
	"fmt"
	"strings"

	"github.com/aristanetworks/goeapi"
	"github.com/aristanetworks/goeapi/module"
	"github.com/sapcc/ironic_temper/pkg/config"
	"github.com/sapcc/ironic_temper/pkg/model"
	log "github.com/sirupsen/logrus"
)

type AristaClient struct {
	cfg config.Config
	log *log.Entry
}

func NewAristaClient(cfg config.Config, ctxLogger *log.Entry) (*AristaClient, error) {
	return &AristaClient{
			cfg: cfg,
			log: ctxLogger,
		},
		nil
}

func (a AristaClient) RunCableCheck(n *model.IronicNode) (err error) {
	a.log.Debug("calling arista api for node cable check")
	foundAllNeighbors := true
	cfg := a.cfg.AristaAuth
	for mac, sw := range n.Connections {
		host := fmt.Sprintf("%s.%s", sw, a.cfg.Domain)
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
			if strings.ToLower(strings.ReplaceAll(mac, ":", "")) == strings.ReplaceAll(ln.NeighborPort, ".", "") {
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
