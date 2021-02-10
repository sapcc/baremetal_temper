package diagnostics

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/sapcc/ironic_temper/pkg/config"
	"github.com/sapcc/ironic_temper/pkg/model"
	log "github.com/sirupsen/logrus"
)

type ACIClient struct {
	cfg config.Config
	log *log.Entry
}

func (a ACIClient) Run(n *model.IronicNode) (err error) {
	foundAllNeighbors := true
	for _, i := range n.Interfaces {
		cfg := a.cfg.AciAuth
		c := client.NewClient("https://"+i.ConnectionIP, cfg.User, client.Password(cfg.Password), client.Insecure(true))
		co, err := c.GetViaURL("/api/node/class/lldpIf.json?rsp-subtree=children&rsp-subtree-class=lldpAdjEp&rsp-subtree-include=required&order-by=lldpIf.id")
		/*
			contList := models.ListFromContainer(co, "lldpIf")
			for _, c := range contList {
				fmt.Println(c.Search("id"))
				fmt.Println(c.Search("mac"))
				fmt.Println(c.Search("portDesc"))
			}
		*/
		if err != nil {
			return err
		}
		l, _ := co.Search("imdata").Children()
		foundNeighbor := false
		for _, c := range l {
			fmt.Println(c.S("lldpIf").S("children").S("lldpAdjEp", "attributes").S("chassisIdV").String(), i.Mac)
			if c.S("lldpIf").S("children").S("lldpAdjEp", "attributes").S("chassisIdV").String() == i.Mac {
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
