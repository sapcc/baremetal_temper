/**
 * Copyright 2021 SAP SE
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package node

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/aristanetworks/goeapi"
	"github.com/aristanetworks/goeapi/module"
	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/sapcc/baremetal_temper/pkg/diagnostics"
	"github.com/stmcginnis/gofish/redfish"
)

func (n *Node) runHardwareChecks() (err error) {
	var dellRe = regexp.MustCompile(`R640|R740|R760|R840`)

	d, err := n.Redfish.GetData()
	if err != nil {
		return
	}
	if dellRe.MatchString(d.Inventory.SystemVendor.Model) {
		c := diagnostics.NewDellClient(*n.Redfish.GetClientConfig(), n.log)
		return c.Run()
	}

	return
}

func (n *Node) runACICheck() (err error) {
	n.log.Debug("calling aci api for node cable check")
	aci := diagnostics.NewACI(n.cfg, n.log)
	noLldp := make([]string, 0)
	d, err := n.Netbox.GetData()
	if err != nil {
		return
	}
	defer func() {
		if len(noLldp) > 0 {
			err = fmt.Errorf("cable check not successful for: %s", noLldp)
		}
	}()

	for _, intf := range d.Interfaces {
		if !strings.Contains(intf.Connection, "aci") {
			continue
		}
		if intf.PortLinkStatus == redfish.DownPortLinkStatus {
			noLldp = append(noLldp, intf.Name+"(interface_down)")
			continue
		}
		var co *container.Container
		if intf.ConnectionIP == "" {
			noLldp = append(noLldp, intf.Name+"(no_aci_ip)")
			continue
		}
		n.log.Debugf("checking interface: %s --> %s", intf.Name, intf.Connection)
		co, err = aci.GetContainer(intf.ConnectionIP)
		if err != nil {
			noLldp = append(noLldp, intf.Name+"("+err.Error()+")")
			continue
		}

		l, _ := co.Search("imdata").Children()
		foundNeighbor := false

	aciPortLoop:
		for _, c := range l {
			var l diagnostics.Lldp
			if err = json.Unmarshal(c.Bytes(), &l); err != nil {
				n.log.Errorf("cannot unmarshal aci lldp: %s", err.Error())
				continue
			}
			for _, ch := range l.LldpIf.LldpIfChildren {
				if ch.LldpAdjEp.LldpAdjEpAttributes.SysDesc != "" {
					interCon := strings.Split(ch.LldpAdjEp.LldpAdjEpAttributes.SysDesc, "/")
					if len(interCon) == 3 {
						n.log.Debugf("intra aci: aci-%s", interCon[2])
					}
				}
				n.log.Debugf("aci lldp: %s / node %s", prepareMac(ch.LldpAdjEp.LldpAdjEpAttributes.PortIdV), prepareMac(intf.Mac))
				if prepareMac(intf.Mac) == prepareMac(ch.LldpAdjEp.LldpAdjEpAttributes.PortIdV) {
					if l.LldpIf.LldpIfAttributes.ID != intf.Port {
						errMsg := fmt.Sprintf("%s(wrong switch port: %s)", intf.Name, l.LldpIf.LldpIfAttributes.ID)
						n.log.Debugf("%s(wrong switch port: %s)", intf.Name, l.LldpIf.LldpIfAttributes.ID)
						noLldp = append(noLldp, errMsg)
						break aciPortLoop
					}
					n.log.Debugf("found aci lldp neighbor: %s", intf.Mac)
					foundNeighbor = true
					break aciPortLoop
				}
			}
		}
		if !foundNeighbor {
			noLldp = append(noLldp, intf.Name+"(lldp_missing)")
		}
	}

	return
}

func (n *Node) runAristaCheck() (err error) {
	noLldp := make([]string, 0)
	cfg := n.cfg.Arista
	netboxData, err := n.Netbox.GetData()
	if err != nil {
		return
	}
	defer func() {
		if len(noLldp) > 0 {
			err = fmt.Errorf("cable check not successful for: %s", noLldp)
		}
	}()
	for _, intf := range netboxData.Interfaces {
		if intf.Connection == "" || !strings.Contains(intf.Connection, "sw") {
			continue
		}
		n.log.Debug("calling arista api for node cable check")
		host := fmt.Sprintf("%s.%s", intf.Connection, n.cfg.Domain)
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
			if strings.ToLower(strings.ReplaceAll(intf.Mac, ":", "")) == strings.ReplaceAll(ln.NeighborPort, ".", "") {
				n.log.Debugf("found aci lldp neighbor: %s", intf.Mac)
				foundNeighbor = true
				break
			}
		}
		if !foundNeighbor {
			noLldp = append(noLldp, intf.Name+"(lldp_missing)")
		}
	}

	return
}
func prepareMac(m string) string {
	return strings.ToLower(strings.ReplaceAll(m, ":", ""))
}
