package diagnostics

import (
	"regexp"

	"github.com/sapcc/ironic_temper/pkg/config"
	"github.com/sapcc/ironic_temper/pkg/model"
	log "github.com/sirupsen/logrus"
	"github.com/stmcginnis/gofish"
)

func GetRemoteDiagnostics(gc gofish.ClientConfig, cfg config.Config, l *log.Entry) (d []func(n *model.IronicNode) error, err error) {
	d = make([]func(n *model.IronicNode) error, 0)
	c, err := gofish.Connect(gc)
	defer c.Logout()
	if err != nil {
		return
	}
	var dellRe = regexp.MustCompile(`R640|R740|R840`)
	s, err := c.Service.Systems()
	if err != nil {
		return nil, err
	}
	if dellRe.MatchString(s[0].Model) {
		d = append(d, DellClient{gCfg: gc, log: l}.Run)
	}

	//TODO: distinguish between bpod and vpod
	if "vpod" == "vpod" {
		d = append(d, ACIClient{cfg, l}.Run)
	} else {
		d = append(d, AristaClient{cfg, l}.Run)
	}

	return
}
