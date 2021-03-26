package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/sapcc/baremetal_temper/pkg/config"
	"github.com/sapcc/baremetal_temper/pkg/model"

	"github.com/go-ping/ping"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/baremetal/apiversions"
	"github.com/gophercloud/gophercloud/openstack/baremetal/v1/nodes"
	"github.com/gophercloud/gophercloud/openstack/baremetal/v1/ports"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/hypervisors"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/services"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/images"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/dns/v2/recordsets"
	"github.com/gophercloud/gophercloud/openstack/dns/v2/zones"
	iDservices "github.com/gophercloud/gophercloud/openstack/identity/v3/services"
	"github.com/gophercloud/gophercloud/pagination"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"
)

//Client is
type Client struct {
	baremetalClient *gophercloud.ServiceClient
	dnsClient       *gophercloud.ServiceClient
	computeClient   *gophercloud.ServiceClient
	keystoneClient  *gophercloud.ServiceClient
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
func NewClient(cfg config.Config, ctxLogger *log.Entry) (*Client, error) {
	provider, err := newProviderClient(cfg.OpenstackAuth)
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
		ctxLogger.Error("baremetal service error", err.Error())
	}

	return &Client{baremetalClient: bc, dnsClient: dc, computeClient: cc, keystoneClient: ic, domain: cfg.Domain, log: ctxLogger, cfg: cfg}, nil
}

func newProviderClient(i config.OpenstackAuth) (pc *gophercloud.ProviderClient, err error) {
	os.Setenv("OS_USERNAME", i.User)
	os.Setenv("OS_PASSWORD", i.Password)
	os.Setenv("OS_PROJECT_NAME", i.ProjectName)
	os.Setenv("OS_DOMAIN_NAME", i.DomainName)
	os.Setenv("OS_PROJECT_DOMAIN_NAME", i.ProjectDomainName)
	os.Setenv("OS_AUTH_URL", i.AuthURL)
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

func (c *Client) ServiceEnabled(service string) (bool, error) {
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

//CheckCreated checks if node was created
func (c *Client) CheckCreated(n *model.Node) error {
	if !n.Baremetal {
		return nil
	}
	c.log.Debug("checking node creation")
	if n.UUID != "" {
		return nil
	}
	r, err := nodes.Get(c.baremetalClient, n.UUID).Extract()
	if err != nil {
		return &NodeNotFoundError{
			Err: fmt.Sprintf("could not find node %s", n.UUID),
		}
	}
	n.ResourceClass = r.ResourceClass
	return nil
}

//ApplyRules applies rules from a json file
func (c *Client) ApplyRules(n *model.Node) (err error) {
	if !n.Baremetal {
		return
	}
	c.log.Debug("applying rules on node")
	rules, err := c.getRules(n)
	if err != nil {
		return
	}
	updateNode := nodes.UpdateOpts{}
	updatePorts := ports.UpdateOpts{}

	for _, n := range rules.Properties.Node {
		updateNode = append(updateNode, nodes.UpdateOperation{
			Op:    n.Op,
			Path:  n.Path,
			Value: n.Value,
		})
	}
	for _, p := range rules.Properties.Port {
		updatePorts = append(updatePorts, ports.UpdateOperation{
			Op:    p.Op,
			Path:  p.Path,
			Value: p.Value,
		})
	}
	if err = c.updatePorts(updatePorts, n); err != nil {
		return
	}

	return c.updateNode(updateNode, n)
}

func (c *Client) updatePorts(opts ports.UpdateOpts, n *model.Node) (err error) {
	listOpts := ports.ListOpts{
		NodeUUID: n.UUID,
	}

	l, err := ports.List(c.baremetalClient, listOpts).AllPages()
	if err != nil {
		return
	}

	ps, err := ports.ExtractPorts(l)
	if err != nil {
		return
	}

	for _, p := range ps {
		cf := wait.ConditionFunc(func() (bool, error) {
			_, err = ports.Update(c.baremetalClient, p.UUID, opts).Extract()
			if err != nil {
				switch err.(type) {
				case gophercloud.ErrDefault409:
					//node is locked
					return false, nil
				}
				return true, err
			}
			return true, nil
		})
		if err = wait.Poll(5*time.Second, 60*time.Second, cf); err != nil {
			return
		}
	}

	return
}

func (c *Client) updateNode(opts nodes.UpdateOpts, n *model.Node) (err error) {
	cf := wait.ConditionFunc(func() (bool, error) {
		r := nodes.Update(c.baremetalClient, n.UUID, opts)
		_, err = r.Extract()
		if err != nil {
			return false, nil
		}
		return true, nil
	})
	return wait.Poll(5*time.Second, 60*time.Second, cf)
}

func (c *Client) getAPIVersion() (*apiversions.APIVersion, error) {
	return apiversions.Get(c.baremetalClient, "v1").Extract()
}

//CreateDNSRecords For creates a dns record for the given node if not exists
func (c *Client) CreateDNSRecords(n *model.Node) (err error) {
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

func (c *Client) createArpaZone(ip string) (zoneID string, err error) {
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

func (c *Client) createDNSRecord(ip, zoneID, recordName, rType string) (err error) {
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

//PowerOn powers on the node
func (c *Client) PowerOn(n *model.Node) (err error) {
	if !n.Baremetal {
		return
	}
	c.log.Debug("powering on node")
	powerStateOpts := nodes.PowerStateOpts{
		Target: nodes.PowerOn,
	}
	r := nodes.ChangePowerState(c.baremetalClient, n.UUID, powerStateOpts)

	if r.Err != nil {
		switch r.Err.(type) {
		case gophercloud.ErrDefault409:
			return fmt.Errorf("cannot power on node %s", n.UUID)
		default:
			return fmt.Errorf("cannot power on node %s", n.UUID)
		}
	}

	cf := wait.ConditionFunc(func() (bool, error) {
		r := nodes.Get(c.baremetalClient, n.UUID)
		n, err := r.Extract()
		if err != nil {
			return false, fmt.Errorf("cannot power on node")
		}
		if n.PowerState != string(nodes.PowerOn) {
			return false, nil
		}
		return true, nil
	})
	return wait.Poll(5*time.Second, 30*time.Second, cf)
}

//Validate calls the baremetal validate api
func (c *Client) Validate(n *model.Node) (err error) {
	if !n.Baremetal {
		return
	}
	c.log.Debug("validating node")
	v, err := nodes.Validate(c.baremetalClient, n.UUID).Extract()
	if err != nil {
		return
	}
	if !v.Inspect.Result {
		return fmt.Errorf(v.Inspect.Reason)
	}
	if !v.Power.Result {
		return fmt.Errorf(v.Power.Reason)
	}

	if !v.Management.Result {
		return fmt.Errorf(v.Management.Reason)
	}

	if !v.Network.Result {
		return fmt.Errorf(v.Network.Reason)
	}
	return
}

//WaitForNovaPropagation calls the hypervisor api to check if new node has been
//propagated to nova
func (c *Client) WaitForNovaPropagation(n *model.Node) (err error) {
	if !n.Baremetal {
		return
	}
	c.log.Debug("waiting for nova propagation")
	cfp := wait.ConditionFunc(func() (bool, error) {
		p, err := hypervisors.List(c.computeClient).AllPages()
		if err != nil {
			return false, err
		}
		hys, err := hypervisors.ExtractHypervisors(p)
		if err != nil {
			return false, err
		}
		for _, hv := range hys {
			if hv.HypervisorHostname == n.UUID {
				if hv.LocalGB > 0 && hv.MemoryMB > 0 {
					return true, nil
				}
			}
		}
		return false, nil
	})

	return wait.Poll(10*time.Second, 20*time.Minute, cfp)
}

//DeployTestInstance creates a new test instance on the newly created node
func (c *Client) DeployTestInstance(n *model.Node) (err error) {
	if !n.Baremetal {
		return
	}
	c.log.Debug("creating test instance on node")
	iID, err := c.getImageID(c.cfg.Deployment.Image)
	zID, err := c.getConductorZone(c.cfg.Deployment.ConductorZone)
	if err != nil {
		return
	}

	opts := servers.CreateOpts{
		Name:             fmt.Sprintf("%s_inspector_test", time.Now().Format("2006-01-02T15:04:05")),
		FlavorRef:        n.ResourceClass,
		ImageRef:         iID,
		AvailabilityZone: fmt.Sprintf("%s::%s", zID, n.UUID),
	}
	r := servers.Create(c.computeClient, opts)
	s, err := r.Extract()
	if err != nil {
		return
	}
	n.InstanceUUID = s.ID

	if err = servers.WaitForStatus(c.computeClient, s.ID, "ACTIVE", 60); err != nil {
		return
	}
	n.InstanceIPv4 = s.AccessIPv4
	pinger, err := ping.NewPinger(n.InstanceIPv4)
	if err != nil {
		return
	}
	pinger.Count = 3
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		return
	}
	return
}

//DeleteTestInstance deletes the test instance via the nova api
func (c *Client) DeleteTestInstance(n *model.Node) (err error) {
	if !n.Baremetal {
		return
	}
	c.log.Debug("deleting instance on node")
	if err = servers.ForceDelete(c.computeClient, n.InstanceUUID).ExtractErr(); err != nil {
		return
	}
	return servers.WaitForStatus(c.computeClient, n.InstanceUUID, "DELETED", 60)
}

func (c *Client) getImageID(name string) (id string, err error) {
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

func (c *Client) getFlavorID(name string) (id string, err error) {
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

func (c *Client) getMatchingFlavorFor(n *model.Node) (name string, err error) {
	err = flavors.ListDetail(c.computeClient, nil).EachPage(func(p pagination.Page) (bool, error) {
		fs, err := flavors.ExtractFlavors(p)
		if err != nil {
			return true, err
		}
		ram := 0.1
		disk := 0.2
		cpu := 0.1
		for _, f := range fs {
			delta := calcDelta(f.RAM, n.InspectionData.Inventory.Memory.PhysicalMb)
			if delta > ram {
				continue
			}
			ram = delta
			delta = calcDelta(f.Disk, int(n.InspectionData.RootDisk.Size))
			if delta > disk {
				continue
			}
			disk = delta
			delta = calcDelta(f.VCPUs, n.InspectionData.Inventory.CPU.Count)
			if delta > cpu {
				continue
			}
			cpu = delta
			name = f.Name
			n.ResourceClass = f.Name
		}
		return false, nil
	})
	return
}

func (c *Client) getConductorZone(name string) (id string, err error) {
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

//Provide sets node provisionstate to provided (available).
//Needed to deploy a test instance on this node
func (c *Client) Provide(n *model.Node) (err error) {
	if !n.Baremetal {
		return
	}
	c.log.Debug("providing node")
	cf := func(tp nodes.TargetProvisionState) wait.ConditionFunc {
		return wait.ConditionFunc(func() (bool, error) {
			if err = nodes.ChangeProvisionState(c.baremetalClient, n.UUID, nodes.ProvisionStateOpts{
				Target: tp,
			}).ExtractErr(); err != nil {
				switch err.(type) {
				case gophercloud.ErrDefault409:
					//node is locked
					return false, nil
				}
				return true, err
			}
			return true, nil
		})
	}
	if err = wait.Poll(5*time.Second, 30*time.Second, cf(nodes.TargetManage)); err != nil {
		return
	}
	if err = wait.Poll(5*time.Second, 30*time.Second, cf(nodes.TargetProvide)); err != nil {
		return
	}

	cfp := wait.ConditionFunc(func() (bool, error) {
		n, err := nodes.Get(c.baremetalClient, n.UUID).Extract()
		if err != nil {
			return true, err
		}

		if n.ProvisionState != "available" {
			return false, nil
		}
		return true, nil
	})

	return wait.Poll(5*time.Second, 30*time.Second, cfp)
}

//Prepare prepares the node for customers.
//Removes resource_class, sets the rightful conductor and maintenance to true
func (c *Client) Prepare(n *model.Node) (err error) {
	if !n.Baremetal {
		return
	}
	c.log.Debug("preparing node")
	conductor := strings.Split(n.Name, "-")[1]
	opts := nodes.UpdateOpts{
		nodes.UpdateOperation{
			Op:    nodes.ReplaceOp,
			Path:  "/conductor_group",
			Value: conductor,
		},
		nodes.UpdateOperation{
			Op:    nodes.ReplaceOp,
			Path:  "/maintenance",
			Value: true,
		},
	}
	return c.updateNode(opts, n)
}

//DeleteNode deletes a node via the baremetal api
func (c *Client) DeleteNode(n *model.Node) (err error) {
	if n.UUID == "" {
		return
	}
	c.log.Debug("deleting node")
	cfp := wait.ConditionFunc(func() (bool, error) {
		err = nodes.Delete(c.baremetalClient, n.UUID).ExtractErr()
		if err != nil {
			return false, err
		}
		return true, nil
	})

	return wait.Poll(5*time.Second, 30*time.Second, cfp)
}

func (c *Client) getRules(n *model.Node) (r config.Rule, err error) {
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
