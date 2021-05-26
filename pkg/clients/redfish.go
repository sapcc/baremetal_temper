package clients

import (
	"fmt"

	"github.com/sapcc/baremetal_temper/pkg/config"
	log "github.com/sirupsen/logrus"
	"github.com/stmcginnis/gofish"
)

type RedfishClient struct {
	ClientConfig *gofish.ClientConfig
	Client       *gofish.APIClient
	log          *log.Entry
	cfg          config.Redfish
}

//NewRedfishClient creates redfish client
func NewRedfishClient(cfg config.Config, ctxLogger *log.Entry) *RedfishClient {
	return &RedfishClient{
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
func (r *RedfishClient) SetEndpoint(remoteIP string) (err error) {
	if remoteIP == "" {
		return fmt.Errorf("no remote ip address set")
	}
	r.ClientConfig.Endpoint = fmt.Sprintf("https://%s", remoteIP)
	return
}

func (r *RedfishClient) Connect() (err error) {
	client, err := gofish.Connect(*r.ClientConfig)
	if err != nil {
		return
	}
	r.Client = client
	return
}

func (r *RedfishClient) Logout() {
	r.Client.Logout()
}
