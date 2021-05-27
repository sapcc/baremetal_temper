package clients

import (
	"fmt"

	"github.com/sapcc/baremetal_temper/pkg/config"
	log "github.com/sirupsen/logrus"
	"github.com/stmcginnis/gofish"
)

type Redfish struct {
	ClientConfig *gofish.ClientConfig
	Client       *gofish.APIClient
	log          *log.Entry
	cfg          config.Redfish
}

//NewRedfishClient creates redfish client
func NewRedfish(cfg config.Config, ctxLogger *log.Entry) *Redfish {
	return &Redfish{
		ClientConfig: &gofish.ClientConfig{
			Endpoint:  fmt.Sprintf("https://%s", "dummy.net"),
			Username:  cfg.Redfish.User,
			Password:  cfg.Redfish.Password,
			Insecure:  true,
			BasicAuth: true,
		},
		cfg: cfg.Redfish,
		log: ctxLogger,
	}
}

//SetEndpoint sets the redfish api endpoint
func (r *Redfish) SetEndpoint(remoteIP string) (err error) {
	if remoteIP == "" {
		return fmt.Errorf("no remote ip address set")
	}
	r.ClientConfig.Endpoint = fmt.Sprintf("https://%s", remoteIP)
	return
}

func (r *Redfish) Connect() (err error) {
	client, err := gofish.Connect(*r.ClientConfig)
	if err != nil {
		return
	}
	r.Client = client
	return
}

func (r *Redfish) Logout() {
	r.Client.Logout()
}
