package diagnostics

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/sapcc/ironic_temper/pkg/config"
	"github.com/sapcc/ironic_temper/pkg/model"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"
)

type ACIClient struct {
	cfg config.Config
	log *log.Entry
	c   map[string]*client.Client
}

func (a ACIClient) Run(n *model.Node) (err error) {
	a.log.Debug("calling aci api for node cable check")
	a.c = make(map[string]*client.Client)
	noLldp := make([]string, 0)
	for iName, i := range n.Interfaces {
		if !strings.Contains(i.Connection, "aci") {
			continue
		}
		var co *container.Container
		c, err := a.getClient(i.ConnectionIP)
		if err != nil {
			return err
		}
		request := wait.ConditionFunc(func() (bool, error) {
			var resp *http.Response
			url := c.BaseURL.String() + "/api/node/class/lldpIf.json?rsp-subtree=children&rsp-subtree-class=lldpAdjEp&rsp-subtree-include=required&order-by=lldpIf.id"
			req, err := c.MakeRestRequest("GET", url, nil, true)
			if err != nil && strings.Contains(err.Error(), "invalid character '<'") {
				return false, nil
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

		if err := wait.Poll(5*time.Second, 2*time.Minute, request); err != nil {
			return err
		}

		l, _ := co.Search("imdata").Children()
		foundNeighbor := false

		for _, c := range l {
			rgx := regexp.MustCompile(`\["(.*?)"\]`)
			m := c.S("lldpIf").S("children").S("lldpAdjEp", "attributes").S("portIdV").String()
			mac := rgx.FindStringSubmatch(m)
			if strings.ToLower(strings.ReplaceAll(i.Mac, ":", "")) == strings.ReplaceAll(mac[1], ":", "") {
				a.log.Debugf("found aci lldap neighbor: %s", m)
				foundNeighbor = true
				break
			}
		}
		if !foundNeighbor {
			noLldp = append(noLldp, iName)
		}
	}

	if len(noLldp) > 0 {
		return fmt.Errorf("cable check not successful for: %s", noLldp)
	}

	return
}

func (a ACIClient) getClient(host string) (c *client.Client, err error) {
	c, ok := a.c[host]
	if !ok {
		cfg := a.cfg.AciAuth
		c = client.NewClient("https://"+host, cfg.User, client.Password(cfg.Password), client.Insecure(true), client.SkipLoggingPayload(true), client.ReqTimeout(20))
		a.c[host] = c
	}
	return
}
