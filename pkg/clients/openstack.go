package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"text/template"

	"github.com/sapcc/baremetal_temper/pkg/config"
	"github.com/sapcc/baremetal_temper/pkg/model"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/baremetal/apiversions"
	"github.com/gophercloud/gophercloud/openstack/baremetal/v1/nodes"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/services"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/images"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/dns/v2/recordsets"
	"github.com/gophercloud/gophercloud/openstack/dns/v2/zones"
	iDservices "github.com/gophercloud/gophercloud/openstack/identity/v3/services"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/pagination"
	log "github.com/sirupsen/logrus"
)

//Client is
type OpenstackClient struct {
	baremetalClient *gophercloud.ServiceClient
	dnsClient       *gophercloud.ServiceClient
	computeClient   *gophercloud.ServiceClient
	keystoneClient  *gophercloud.ServiceClient
	networkClient   *gophercloud.ServiceClient
	domain          string
	log             *log.Entry
	cfg             config.Config
}

//NodeNotFoundError error for missing node
type NodeNotFoundError struct {
	Err string
}

func (n *NodeNotFoundError) Error() string {
	return n.Err
}

//NewClient creates a new client containing different openstack-clients (baremetal, compute, dns)
func NewClient(cfg config.Config, ctxLogger *log.Entry) (*OpenstackClient, error) {
	provider, err := newProviderClient(cfg.Openstack)
	if err != nil {
		return nil, err
	}

	dc, err := openstack.NewDNSV2(provider, gophercloud.EndpointOpts{
		Region: cfg.Region,
	})
	if err != nil {
		return nil, err
	}

	cc, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: cfg.Region,
	})
	if err != nil {
		return nil, err
	}
	ic, err := openstack.NewIdentityV3(provider, gophercloud.EndpointOpts{
		Region: cfg.Region,
	})
	if err != nil {
		return nil, err
	}

	nc, err := openstack.NewNetworkV2(provider, gophercloud.EndpointOpts{
		Region: cfg.Region,
	})

	bc, err := openstack.NewBareMetalV1(provider, gophercloud.EndpointOpts{
		Region: cfg.Region,
	})
	if err == nil {
		version, err := apiversions.Get(bc, "v1").Extract()
		if err != nil {
			return nil, err
		}
		bc.Microversion = version.Version
	} else {
		ctxLogger.Infof("baremetal service error: %s", err.Error())
	}

	return &OpenstackClient{networkClient: nc, baremetalClient: bc, dnsClient: dc, computeClient: cc, keystoneClient: ic, domain: cfg.Domain, log: ctxLogger, cfg: cfg}, nil
}

func newProviderClient(i config.OpenstackAuth) (pc *gophercloud.ProviderClient, err error) {
	os.Setenv("OS_USERNAME", i.User)
	os.Setenv("OS_PASSWORD", i.Password)
	os.Setenv("OS_PROJECT_NAME", i.ProjectName)
	os.Setenv("OS_DOMAIN_NAME", i.DomainName)
	os.Setenv("OS_PROJECT_DOMAIN_NAME", i.ProjectDomainName)
	os.Setenv("OS_AUTH_URL", i.Url)
	opts, err := openstack.AuthOptionsFromEnv()
	opts.AllowReauth = true
	opts.Scope = &gophercloud.AuthScope{
		ProjectName: opts.TenantName,
		DomainName:  os.Getenv("OS_PROJECT_DOMAIN_NAME"),
	}

	pc, err = openstack.AuthenticatedClient(opts)
	if err != nil {
		return pc, err
	}

	pc.UseTokenLock()

	return pc, nil
}

func (c *OpenstackClient) ServiceEnabled(service string) (bool, error) {
	a, err := iDservices.List(c.keystoneClient, iDservices.ListOpts{Name: service}).AllPages()
	if err != nil {
		return false, err
	}
	s, err := iDservices.ExtractServices(a)
	if len(s) == 0 {
		return false, fmt.Errorf("no service found")
	}
	if s[0].Enabled {
		return true, nil
	}

	return false, nil
}

func (c *OpenstackClient) getAPIVersion() (*apiversions.APIVersion, error) {
	return apiversions.Get(c.baremetalClient, "v1").Extract()
}

//CreateDNSRecords For creates a dns record for the given node if not exists
func (c *OpenstackClient) CreateDNSRecords(n *model.Node) (err error) {
	c.log.Debug("creating dns record")
	opts := zones.ListOpts{
		Name: c.domain + ".",
	}
	allPages, err := zones.List(c.dnsClient, opts).AllPages()
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
		log.Debug("Create A recordset:  ", ip.String(), allZones[0].ID, a.DNSName)

		if err = c.createDNSRecord(ip.String(), allZones[0].ID, a.DNSName+".", "A"); err != nil {
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
		zoneID, err := c.createArpaZone(ip.String())
		if err != nil {
			return err
		}
		log.Debug("Create PTR recordset: ", a.DNSName, zoneID, arpa)
		if err = c.createDNSRecord(a.DNSName+".", zoneID, arpa, "PTR"); err != nil {
			return err
		}
	}

	return
}

func (c *OpenstackClient) createArpaZone(ip string) (zoneID string, err error) {
	arpaZone, err := reverseZone(ip)
	if err != nil {
		return
	}

	allPages, err := zones.List(c.dnsClient, zones.ListOpts{
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
		z, err := zones.Create(c.dnsClient, zones.CreateOpts{
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

func (c *OpenstackClient) createDNSRecord(ip, zoneID, recordName, rType string) (err error) {
	_, err = recordsets.Create(c.dnsClient, zoneID, recordsets.CreateOpts{
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

func (c *OpenstackClient) getImageID(name string) (id string, err error) {
	err = images.ListDetail(c.computeClient, images.ListOpts{Name: name}).EachPage(
		func(p pagination.Page) (bool, error) {
			is, err := images.ExtractImages(p)
			if err != nil {
				return false, err
			}
			for _, i := range is {
				if i.Name == name {
					id = i.ID
					return false, nil
				}
			}
			return true, nil
		},
	)
	return
}

func (c *OpenstackClient) getFlavorID(name string) (id string, err error) {
	err = flavors.ListDetail(c.computeClient, nil).EachPage(func(p pagination.Page) (bool, error) {
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

func (c *OpenstackClient) getMatchingFlavorFor(n *model.Node) (name string, err error) {
	mem := 0.1
	disk := 0.2
	cpu := 0.1
	var fl flavors.Flavor
	err = flavors.ListDetail(c.computeClient, nil).EachPage(func(p pagination.Page) (bool, error) {
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

func (c *OpenstackClient) getConductorZone(name string) (id string, err error) {
	err = services.List(c.computeClient, services.ListOpts{Host: name}).EachPage(
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

func (c *OpenstackClient) getRules(n *model.Node) (r config.Rule, err error) {
	var funcMap = template.FuncMap{
		"imageToID":            c.getImageID,
		"getMatchingFlavorFor": c.getMatchingFlavorFor,
	}

	tmpl := template.New("rules.json").Funcs(funcMap)
	t, err := tmpl.ParseFiles(c.cfg.RulesPath)
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

func (c *OpenstackClient) getNetwork(name string) (n servers.Network, err error) {
	pr, err := newProviderClient(c.cfg.Deployment.Openstack)
	if err != nil {
		return
	}
	nc, err := openstack.NewNetworkV2(pr, gophercloud.EndpointOpts{
		Region: c.cfg.Region,
	})
	p, err := networks.List(nc, networks.ListOpts{Name: name}).AllPages()
	if err != nil {
		return
	}
	ns, err := networks.ExtractNetworks(p)
	if err != nil || len(ns) != 1 {
		return
	}
	n.UUID = ns[0].ID
	return
}
