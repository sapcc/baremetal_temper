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
	"net"
	"regexp"
	"strconv"
	"text/template"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/aggregates"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/services"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/images"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/dns/v2/recordsets"
	"github.com/gophercloud/gophercloud/openstack/dns/v2/zones"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/pagination"
	"github.com/sapcc/baremetal_temper/pkg/clients"
	"github.com/sapcc/baremetal_temper/pkg/config"
	"k8s.io/apimachinery/pkg/util/wait"
)

// CreateDNSRecords For creates a dns record for the given node if not exists
func (n *Node) createDNSRecords() (err error) {
	c, err := n.oc.GetServiceClient("dns")
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

		if err = n.createDNSRecord(ip.String(), allZones[0].ID, a.DNSName+".", "A"); err != nil {
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
		zoneID, err := n.createArpaZone(ip.String())
		if err != nil {
			return err
		}
		n.log.Debug("Create PTR recordset: ", a.DNSName, zoneID, arpa)
		if err = n.createDNSRecord(a.DNSName+".", zoneID, arpa, "PTR"); err != nil {
			return err
		}
	}

	return
}

func (n *Node) getRules() (r config.Rule, err error) {
	var funcMap = template.FuncMap{
		"imageToID":                n.getImageID,
		"getMatchingFlavorForNode": n.getMatchingFlavorFor,
		"getRootDeviceSize":        n.getRootDeviceSize,
		"getPortGroupUUID":         n.createPortGroup,
	}

	tmpl := template.New("rules.json").Funcs(funcMap)
	t, err := tmpl.ParseFiles(n.cfg.RulesPath)
	if err != nil {
		return r, fmt.Errorf("Error parsing rules: %s", err.Error())
	}

	data, err := n.Redfish.GetData()
	out := new(bytes.Buffer)
	d := map[string]interface{}{
		"node":      n.Name,
		"inventory": data.Inventory,
	}
	err = t.Execute(out, d)
	if err != nil {
		return
	}
	json.Unmarshal(out.Bytes(), &r)
	return
}

func (n *Node) getRootDeviceSize() (size int64, err error) {
	data, err := n.Redfish.GetData()
	if err != nil {
		return
	}
	size = data.RootDisk.Size / 1024 / 1024 / 1024
	l := len(strconv.FormatInt(size, 10))
	switch l {
	case 3:
		size = 100 * ((size + 90) / 100)
	case 4:
		size = 1000 * ((size + 900) / 1000)
	}
	if size == 0 {
		return size, fmt.Errorf("could not get root disk size")
	}
	return
}

func (n *Node) getMatchingFlavorFor() (name string, err error) {
	c, err := n.oc.GetServiceClient("compute")
	if err != nil {
		return
	}
	data, err := n.Redfish.GetData()
	if err != nil {
		return
	}
	mem := 0.1
	//disk := 0.2
	cpu := 0.1
	var fl flavors.Flavor
	flavorNameRules := regexp.MustCompile(`^zh.+|^hv_.+`)
	err = flavors.ListDetail(c, flavors.ListOpts{AccessType: n.cfg.FlavorAccessType}).EachPage(func(p pagination.Page) (bool, error) {
		fs, err := flavors.ExtractFlavors(p)
		if err != nil {
			return false, err
		}
		for _, f := range fs {
			if !flavorNameRules.MatchString(f.Name) {
				continue
			}
			deltaMem := calcDelta(f.RAM, data.Inventory.Memory.PhysicalMb)
			//deltaDisk := calcDelta(f.Disk, int(data.RootDisk.Size/1024/1024/1024))
			deltaCPU := calcDelta(f.VCPUs, data.Inventory.CPU.Count)
			if deltaMem <= mem && deltaCPU <= cpu {
				mem = deltaMem
				//disk = deltaDisk
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
	data.Inventory.Memory.PhysicalMb = fl.RAM
	//data.RootDisk.Size = int64(fl.Disk)
	data.Inventory.CPU.Count = fl.VCPUs
	return
}

func (n *Node) enableComputeService(svc services.Service) (id string, err error) {
	cl, err := n.oc.GetServiceClient("compute")
	if err != nil {
		return
	}
	cl.Microversion = "2.53"
	if svc.Status == string(services.ServiceDisabled) {
		r := services.Update(cl, svc.ID, services.UpdateOpts{Status: services.ServiceEnabled})
		return svc.ID, r.Err
	}
	return
}

func (n *Node) waitComputeServiceCreated(host string) (svc services.Service, err error) {
	cl, err := n.oc.GetServiceClient("compute")
	if err != nil {
		return
	}
	cl.Microversion = "2.53"
	n.log.Debug("waiting compute service to be created")
	cf := wait.ConditionFunc(func() (bool, error) {
		ps := services.List(cl, services.ListOpts{Host: host})
		p, err := ps.AllPages()
		if err != nil {
			return false, nil
		}
		//var svc services.Service
		svcs, err := services.ExtractServices(p)
		if err != nil {
			return false, err
		}
		for _, s := range svcs {
			if s.Host == host {
				svc = s
				break
			}
		}
		if svc.Status == "" {
			n.log.Debugf("compute service %s does not exist yet", host)
			return false, nil
		}
		return true, nil
	})
	return svc, wait.Poll(20*time.Second, 11*time.Minute, cf)
}

func (n *Node) addHostToAggregate(host, az string) (err error) {
	cl, err := n.oc.GetServiceClient("compute")
	if err != nil {
		return
	}

	p := aggregates.List(cl)
	ps, err := p.AllPages()
	if err != nil {
		return
	}
	aggs, err := aggregates.ExtractAggregates(ps)
	if err != nil {
		return
	}
	var aggregate aggregates.Aggregate
	for _, a := range aggs {
		if a.AvailabilityZone == az && a.Name == az {
			aggregate = a
			break
		}
	}
	if aggregate.Name == "" {
		return fmt.Errorf("cannot find aggregate for az: %s", az)
	}

	foundHost := false
	for _, h := range aggregate.Hosts {
		if h == host {
			foundHost = true
			break
		}
	}
	if !foundHost {
		r := aggregates.AddHost(cl, aggregate.ID, aggregates.AddHostOpts{Host: host})
		if r.Header.Values("Status")[0] == "409" {
			return
		}
		return r.Err
	}
	return
}

func (n *Node) createArpaZone(ip string) (zoneID string, err error) {
	c, err := n.oc.GetServiceClient("dns")
	if err != nil {
		return
	}
	arpaZone, err := reverseZone(ip)
	if err != nil {
		return
	}

	allPages, err := zones.List(c, zones.ListOpts{
		Name: arpaZone,
	}).AllPages()
	if err != nil {
		return
	}
	allZones, err := zones.ExtractZones(allPages)
	if err != nil {
		return
	}

	if len(allZones) == 0 {
		z, err := zones.Create(c, zones.CreateOpts{
			Name:        arpaZone,
			TTL:         3600,
			Description: "An in-addr.arpa. zone for reverse lookups set up by baremetal temper",
			Email:       "stefan.hipfel@sap.com",
		}).Extract()
		if err != nil {
			return zoneID, err
		}
		zoneID = z.ID
	} else {
		zoneID = allZones[0].ID
	}
	return
}

func (n *Node) createDNSRecord(ip, zoneID, recordName, rType string) (err error) {
	c, err := n.oc.GetServiceClient("dns")
	if err != nil {
		return
	}
	_, err = recordsets.Create(c, zoneID, recordsets.CreateOpts{
		Name:    recordName,
		TTL:     3600,
		Type:    rType,
		Records: []string{ip},
	}).Extract()
	if httpStatus, ok := err.(gophercloud.ErrDefault409); ok {
		if httpStatus.Actual == 409 {
			// record already exists
			return nil
		}
	}
	return
}

func (n *Node) getImageID(name string) (id string, err error) {
	cl, err := n.oc.GetServiceClient("compute")
	if err != nil {
		return
	}
	err = images.ListDetail(cl, images.ListOpts{Name: name, Status: "active"}).EachPage(
		func(p pagination.Page) (bool, error) {
			is, err := images.ExtractImages(p)
			if err != nil {
				return false, err
			}
			var latest time.Time
			for _, i := range is {
				g, ok := i.Metadata["git_branch"]
				if !ok || fmt.Sprintf("%v", g) != "master" {
					continue
				}
				//2021-06-28T14:04:58
				if i.Name == name {
					ts, err := time.Parse(time.RFC3339, i.Created)
					if err != nil {
						continue
					}
					if ts.After(latest) {
						id = i.ID
						latest = ts
					}
					continue
				}
			}
			return true, nil
		},
	)
	return
}

func (n *Node) getFlavorID(name string) (id string, err error) {
	cl, err := n.oc.GetServiceClient("compute")
	if err != nil {
		return
	}
	err = flavors.ListDetail(cl, nil).EachPage(func(p pagination.Page) (bool, error) {
		fs, err := flavors.ExtractFlavors(p)
		if err != nil {
			return true, err
		}
		for _, f := range fs {
			if f.Name == name {
				id = f.ID
				return true, nil
			}
		}
		return false, nil
	})
	return
}

func (n *Node) getConductorZone(name string) (id string, err error) {
	cl, err := n.oc.GetServiceClient("compute")
	if err != nil {
		return
	}
	err = services.List(cl, services.ListOpts{Host: name}).EachPage(
		func(p pagination.Page) (bool, error) {
			svs, err := services.ExtractServices(p)
			if err != nil {
				return true, err
			}
			for _, sv := range svs {
				if sv.Host == name {
					id = sv.Zone
					return true, nil
				}
			}
			return false, nil
		})
	return
}

func (n *Node) getNetwork(name string) (net servers.Network, err error) {
	pr, err := clients.NewProviderClient(n.cfg.Deployment.Openstack)
	if err != nil {
		return
	}
	nc, err := openstack.NewNetworkV2(pr, gophercloud.EndpointOpts{
		Region: n.cfg.Region,
	})
	p, err := networks.List(nc, networks.ListOpts{Name: name}).AllPages()
	if err != nil {
		return
	}
	ns, err := networks.ExtractNetworks(p)
	if err != nil || len(ns) != 1 {
		return
	}
	net.UUID = ns[0].ID
	return
}
