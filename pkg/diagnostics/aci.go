package diagnostics

import (
	"net/http"
	"strings"
	"time"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/sapcc/baremetal_temper/pkg/config"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"
)

type ACIClient struct {
	Cfg config.Config
	Log *log.Entry
	c   map[string]*client.Client
	co  map[string]*container.Container
}

type Lldp struct {
	LldpIf LldpIf `json:"lldpIf"`
}

type LldpIf struct {
	LldpIfAttributes LldpIfAttributes `json:"attributes"`
	LldpIfChildren   []LldpIfChildren `json:"children"`
}

type LldpIfAttributes struct {
	ID       string `json:"id"`
	Mac      string `json:"mac"`
	PortMode string `json:"portMode"`
	Status   string `json:"status"`
}

type LldpIfChildren struct {
	LldpAdjEp LldpAdjEp `json:"lldpAdjEp"`
}

type LldpAdjEp struct {
	LldpAdjEpAttributes LldpAdjEpAttributes `json:"attributes"`
}

type LldpAdjEpAttributes struct {
	ChassisIdT string `json:"chassisIdT"`
	MgmtIp     string `json:"mgmtIp"`
	PortIdV    string `json:"portIdV"`
	PortVlan   string `json:"portVlan"`
	Status     string `json:"status"`
	SysDesc    string `json:"sysDesc"`
}

func NewACI(cfg config.Config, log *log.Entry) (c *ACIClient) {
	return &ACIClient{
		Cfg: cfg,
		Log: log,
		c:   make(map[string]*client.Client, 0),
		co:  make(map[string]*container.Container, 0),
	}
}

func (a ACIClient) GetClient(host string) (c *client.Client) {
	c, ok := a.c[host]
	if !ok {
		cfg := a.Cfg.Aci
		c = client.NewClient("https://"+host, cfg.User, client.Password(cfg.Password), client.Insecure(true), client.SkipLoggingPayload(true), client.ReqTimeout(20))
		a.c[host] = c
	}
	return
}

func (a ACIClient) GetContainer(host string) (co *container.Container, err error) {
	co, ok := a.co[host]
	if ok {
		return
	}
	c := a.GetClient(host)
	request := wait.ConditionFunc(func() (bool, error) {
		var resp *http.Response
		url := c.BaseURL.String() + "/api/node/class/lldpIf.json?rsp-subtree=children&rsp-subtree-class=lldpAdjEp&rsp-subtree-include=required&order-by=lldpIf.id"
		req, err := c.MakeRestRequest("GET", url, nil, true)
		if err != nil && strings.Contains(err.Error(), "invalid character '<'") {
			return false, nil
		}
		if err != nil {
			return true, err
		}
		co, resp, _ = c.Do(req)
		if resp.StatusCode == http.StatusServiceUnavailable {
			return false, nil
		}
		if resp.StatusCode == http.StatusOK {
			return true, nil
		}
		if err = client.CheckForErrors(co, "GET", true); err != nil {
			return false, nil
		}
		return true, err
	})

	if err = wait.Poll(5*time.Second, 2*time.Minute, request); err != nil {
		return
	}
	a.co[host] = co
	return
}
