package diagnostics

import (
	"regexp"

	"github.com/sapcc/ironic_temper/pkg/config"
	"github.com/sapcc/ironic_temper/pkg/model"
	log "github.com/sirupsen/logrus"
	"github.com/stmcginnis/gofish"
)

func GetDiagnosticTasks(n model.Node, gc gofish.ClientConfig, cfg config.Config, l *log.Entry) (d []func(n *model.Node) error, err error) {
	d = make([]func(n *model.Node) error, 0)
	c, err := gofish.Connect(gc)
	if err != nil {
		return
	}
	defer c.Logout()
	d = append(d, ACIClient{
		cfg: cfg,
		log: l,
	}.Run, AristaClient{cfg, l}.Run)

	var dellRe = regexp.MustCompile(`R640|R740|R840`)
	s, err := c.Service.Systems()
	if err != nil {
		return nil, err
	}
	if dellRe.MatchString(s[0].Model) {
		d = append(d, DellClient{gCfg: gc, log: l}.Run)
	}

	return
}
