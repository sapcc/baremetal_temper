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
	var dellRe = regexp.MustCompile(`R640|R740|R840`)

	if dellRe.MatchString(n.InspectionData.Inventory.SystemVendor.Model) {
		c := diagnostics.NewDellClient(*n.Clients.Redfish.ClientConfig, n.log)
		return c.Run()
	}

	return
}

func (n *Node) runACICheck() (err error) {
	n.log.Debug("calling aci api for node cable check")
	aci := diagnostics.NewACI(n.cfg, n.log)
	noLldp := make([]string, 0)

	defer func() {
		if len(noLldp) > 0 {
			err = fmt.Errorf("cable check not successful for: %s", noLldp)
		}
	}()

	for iName, i := range n.Interfaces {
		n.log.Debug("checking interface: %s --> %s", iName, i.Connection)
		if !strings.Contains(i.Connection, "aci") {
			continue
		}
		if i.PortLinkStatus == redfish.DownPortLinkStatus {
			noLldp = append(noLldp, iName+"(interface_down)")
			continue
		}
		var co *container.Container
		if i.ConnectionIP == "" {
			noLldp = append(noLldp, iName+"(no_aci_ip)")
			continue
		}
		co, err = aci.GetContainer(i.ConnectionIP)
		if err != nil {
			noLldp = append(noLldp, iName+"("+err.Error()+")")
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
				n.log.Debugf("aci lldp: %s / node %s", prepareMac(ch.LldpAdjEp.LldpAdjEpAttributes.PortIdV), prepareMac(i.Mac))
				if prepareMac(i.Mac) == prepareMac(ch.LldpAdjEp.LldpAdjEpAttributes.PortIdV) {
					if l.LldpIf.LldpIfAttributes.ID != i.Port {
						errMsg := fmt.Sprintf("%s(wrong switch port: %s)", iName, l.LldpIf.LldpIfAttributes.ID)
						n.log.Debugf("%s(wrong switch port: %s)", iName, l.LldpIf.LldpIfAttributes.ID)
						noLldp = append(noLldp, errMsg)
						break aciPortLoop
					}
					n.log.Debugf("found aci lldp neighbor: %s", i.Mac)
					foundNeighbor = true
					break aciPortLoop
				}
			}
		}
		if !foundNeighbor {
			noLldp = append(noLldp, iName+"(lldp_missing)")
		}
	}

	return
}

func (n *Node) runAristaCheck() (err error) {
	noLldp := make([]string, 0)
	cfg := n.cfg.Arista
	defer func() {
		if len(noLldp) > 0 {
			err = fmt.Errorf("cable check not successful for: %s", noLldp)
		}
	}()
	for iName, i := range n.Interfaces {
		if !strings.Contains(i.Connection, "sw") {
			continue
		}
		n.log.Debug("calling arista api for node cable check")
		host := fmt.Sprintf("%s.%s", i.Connection, n.cfg.Domain)
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
			if strings.ToLower(strings.ReplaceAll(i.Mac, ":", "")) == strings.ReplaceAll(ln.NeighborPort, ".", "") {
				n.log.Debugf("found aci lldp neighbor: %s", i.Mac)
				foundNeighbor = true
				break
			}
		}
		if !foundNeighbor {
			noLldp = append(noLldp, iName+"(lldp_missing)")
		}
	}

	return
}
func prepareMac(m string) string {
	return strings.ToLower(strings.ReplaceAll(m, ":", ""))
}
