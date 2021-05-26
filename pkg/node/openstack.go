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
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net"
	"strconv"

	"github.com/gophercloud/gophercloud/openstack/baremetal/v1/nodes"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/openstack/dns/v2/zones"
	"github.com/gophercloud/gophercloud/pagination"
	"github.com/sapcc/baremetal_temper/pkg/config"
)

//CreateDNSRecords For creates a dns record for the given node if not exists
func (n *Node) CreateDNSRecords() (err error) {
	c, err := n.oc.GetServiceClient(n.cfg, "dns")
	if err != nil {
		return
	}
	n.log.Debug("creating dns record")
	opts := zones.ListOpts{
		Name: n.cfg.Domain + ".",
	}
	allPages, err := zones.List(c, opts).AllPages()
	if err != nil {
		return
	}
	allZones, err := zones.ExtractZones(allPages)
	if err != nil || len(allZones) == 0 {
		return fmt.Errorf("wrong dns zone")
	}

	for _, a := range n.IpamAddresses {
		var ip net.IP
		ip, _, err = net.ParseCIDR(*a.Address)
		if err != nil {
			return
		}
		n.log.Debug("Create A recordset:  ", ip.String(), allZones[0].ID, a.DNSName)

		if err = n.oc.CreateDNSRecord(ip.String(), allZones[0].ID, a.DNSName+".", "A"); err != nil {
			return
		}

	}

	for _, a := range n.IpamAddresses {
		var arpa string
		var ip net.IP
		ip, _, err = net.ParseCIDR(*a.Address)
		if err != nil {
			return
		}
		arpa, err = reverseaddr(ip.String())
		if err != nil {
			return err
		}
		zoneID, err := n.oc.CreateArpaZone(ip.String())
		if err != nil {
			return err
		}
		n.log.Debug("Create PTR recordset: ", a.DNSName, zoneID, arpa)
		if err = n.oc.CreateDNSRecord(a.DNSName+".", zoneID, arpa, "PTR"); err != nil {
			return err
		}
	}

	return
}

func (n *Node) getRules() (r config.Rule, err error) {
	var funcMap = template.FuncMap{
		"imageToID":            n.oc.GetImageID,
		"getMatchingFlavorFor": n.getMatchingFlavorFor,
		"getRootDeviceSize":    n.getRootDeviceSize,
	}

	tmpl := template.New("rules.json").Funcs(funcMap)
	t, err := tmpl.ParseFiles(n.cfg.RulesPath)
	if err != nil {
		return r, fmt.Errorf("Error parsing rules: %s", err.Error())
	}

	out := new(bytes.Buffer)
	d := map[string]interface{}{
		"node": n,
	}
	err = t.Execute(out, d)
	if err != nil {

	}
	json.Unmarshal(out.Bytes(), &r)

	return
}

func (n *Node) getRootDeviceSize() (size int64, err error) {
	size = n.InspectionData.RootDisk.Size / 1024 / 1024 / 1024
	l := len(strconv.FormatInt(size, 10))
	switch l {
	case 3:
		size = 100 * ((size + 90) / 100)
	case 4:
		size = 1000 * ((size + 900) / 1000)
	}
	return
}

func (n *Node) getMatchingFlavorFor() (name string, err error) {
	c, err := n.oc.GetServiceClient(n.cfg, "compute")
	if err != nil {
		return
	}
	mem := 0.1
	disk := 0.2
	cpu := 0.1
	var fl flavors.Flavor
	err = flavors.ListDetail(c, nil).EachPage(func(p pagination.Page) (bool, error) {
		fs, err := flavors.ExtractFlavors(p)
		if err != nil {
			return false, err
		}
		for _, f := range fs {
			deltaMem := calcDelta(f.RAM, n.InspectionData.Inventory.Memory.PhysicalMb)
			deltaDisk := calcDelta(f.Disk, int(n.InspectionData.RootDisk.Size/1024/1024/1024))
			deltaCPU := calcDelta(f.VCPUs, n.InspectionData.Inventory.CPU.Count)
			if deltaMem <= mem && deltaDisk <= disk && deltaCPU <= cpu {
				mem = deltaMem
				disk = deltaDisk
				cpu = deltaCPU
				name = f.Name
				fl = f
			}
		}
		return true, nil
	})
	if name == "" {
		return name, fmt.Errorf("no matching flavor found for node")
	}
	n.InspectionData.Inventory.Memory.PhysicalMb = fl.RAM
	n.InspectionData.RootDisk.Size = int64(fl.Disk)
	n.InspectionData.Inventory.CPU.Count = fl.VCPUs
	updateNode := nodes.UpdateOpts{}
	updateNode = append(updateNode, nodes.UpdateOperation{
		Op:    nodes.ReplaceOp,
		Path:  "/properties/memory_mb",
		Value: fl.RAM,
	})
	updateNode = append(updateNode, nodes.UpdateOperation{
		Op:    nodes.ReplaceOp,
		Path:  "/properties/local_gb",
		Value: fl.Disk,
	})
	updateNode = append(updateNode, nodes.UpdateOperation{
		Op:    nodes.ReplaceOp,
		Path:  "/properties/cpus",
		Value: fl.VCPUs,
	})
	return
}
