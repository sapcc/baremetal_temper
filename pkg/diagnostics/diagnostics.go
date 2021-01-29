package diagnostics

import (
	"regexp"

	log "github.com/sirupsen/logrus"
	"github.com/stmcginnis/gofish"
)

func GetRemoteDiagnostics(c *gofish.APIClient, l *log.Entry) (Diagnostics, error) {
	var dellRe = regexp.MustCompile(`R640|R740|R840`)
	s, err := c.Service.Systems()
	if err != nil {
		return nil, err
	}
	if dellRe.MatchString(s[0].Model) {
		return NewDell(c, l), nil
	}

	return nil, err
}
