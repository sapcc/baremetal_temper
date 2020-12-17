package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/sapcc/ironic_temper/pkg/config"
	"github.com/sapcc/ironic_temper/pkg/model"

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
	"github.com/gophercloud/gophercloud/pagination"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"
)

type Client struct {
	IronicClient  *gophercloud.ServiceClient
	DNSClient     *gophercloud.ServiceClient
	ComputeClient *gophercloud.ServiceClient
	Domain        string
	log           *log.Entry
	cfg           config.Config
}

type NodeNotFoundError struct {
	Err string
}

func (n *NodeNotFoundError) Error() string {
	return n.Err
}

func NewClient(cfg config.Config, ctxLogger *log.Entry) (*Client, error) {
	provider, err := newProviderClient(cfg.IronicAuth)
	if err != nil {
		return nil, err
	}
	iclient, err := openstack.NewBareMetalV1(provider, gophercloud.EndpointOpts{
		Region: cfg.OsRegion,
	})

	dnsClient, err := openstack.NewDNSV2(provider, gophercloud.EndpointOpts{
		Region: cfg.OsRegion,
	})

	cclient, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: cfg.OsRegion,
	})

	if err != nil {
		return nil, err
	}
	version, err := apiversions.Get(iclient, "v1").Extract()
	if err != nil {
		return nil, err
	}
	iclient.Microversion = version.Version
	return &Client{IronicClient: iclient, DNSClient: dnsClient, ComputeClient: cclient, Domain: cfg.Domain, log: ctxLogger, cfg: cfg}, nil
}

func newProviderClient(i config.IronicAuth) (pc *gophercloud.ProviderClient, err error) {
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

func (c *Client) CheckIronicNodeCreated(n *model.IronicNode) error {
	c.log.Info("checking node creation")
	if n.UUID != "" {
		return nil
	}
	_, err := nodes.Get(c.IronicClient, n.UUID).Extract()
	if err != nil {
		return &NodeNotFoundError{
			Err: fmt.Sprintf("could not find node %s", n.UUID),
		}
	}
	return nil
}

func (c *Client) ApplyRules(n *model.IronicNode) (err error) {
	c.log.Info("updating node")
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

func (c *Client) updatePorts(opts ports.UpdateOpts, n *model.IronicNode) (err error) {
	listOpts := ports.ListOpts{
		NodeUUID: n.UUID,
	}

	l, err := ports.List(c.IronicClient, listOpts).AllPages()
	if err != nil {
		return
	}

	ps, err := ports.ExtractPorts(l)
	if err != nil {
		return
	}

	for _, p := range ps {
		cf := wait.ConditionFunc(func() (bool, error) {
			_, err = ports.Update(c.IronicClient, p.UUID, opts).Extract()
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

func (c *Client) updateNode(opts nodes.UpdateOpts, n *model.IronicNode) (err error) {
	cf := wait.ConditionFunc(func() (bool, error) {
		r := nodes.Update(c.IronicClient, n.UUID, opts)
		_, err = r.Extract()
		fmt.Println(err)
		if err != nil {
			return false, nil
		}
		return true, nil
	})
	return wait.Poll(5*time.Second, 60*time.Second, cf)
}

func (c *Client) getAPIVersion() (*apiversions.APIVersion, error) {
	return apiversions.Get(c.IronicClient, "v1").Extract()
}

func (c *Client) CreateDNSRecordFor(n *model.IronicNode) (err error) {
	c.log.Info("creating dns record")
	opts := zones.ListOpts{
		Name: c.Domain + ".",
	}
	allPages, err := zones.List(c.DNSClient, opts).AllPages()
	if err != nil {
		return
	}
	allZones, err := zones.ExtractZones(allPages)
	if err != nil || len(allZones) == 0 {
		return fmt.Errorf("wrong dns zone")
	}

	na := strings.Split(n.Name, "-")

	if len(na) < 1 {
		return fmt.Errorf("wrong node name")
	}

	name := fmt.Sprintf("%sr-%s", na[0], na[1])
	recordName := fmt.Sprintf("%s.%s.", name, c.Domain)
	n.Host = recordName

	_, err = recordsets.Create(c.DNSClient, allZones[0].ID, recordsets.CreateOpts{
		Name:    recordName,
		TTL:     3600,
		Type:    "A",
		Records: []string{n.IP},
	}).Extract()
	if httpStatus, ok := err.(gophercloud.ErrDefault409); ok {
		if httpStatus.Actual == 409 {
			// record already exists
			return nil
		}
	}

	return
}

func (c *Client) PowerNodeOn(n *model.IronicNode) (err error) {
	c.log.Info("powering on node")
	powerStateOpts := nodes.PowerStateOpts{
		Target: nodes.PowerOn,
	}
	r := nodes.ChangePowerState(c.IronicClient, n.UUID, powerStateOpts)

	if r.Err != nil {
		switch r.Err.(type) {
		case gophercloud.ErrDefault409:
			return fmt.Errorf("cannot power on node %s", n.UUID)
		default:
			return fmt.Errorf("cannot power on node %s", n.UUID)
		}
	}

	cf := wait.ConditionFunc(func() (bool, error) {
		r := nodes.Get(c.IronicClient, n.UUID)
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

func (c *Client) ValidateNode(n *model.IronicNode) (err error) {
	c.log.Info("validating node")
	v, err := nodes.Validate(c.IronicClient, n.UUID).Extract()
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

func (c *Client) WaitForNovaPropagation(n *model.IronicNode) (err error) {
	c.log.Info("waiting for nova propagation")
	cfp := wait.ConditionFunc(func() (bool, error) {
		p, err := hypervisors.List(c.ComputeClient).AllPages()
		if err != nil {
			return true, err
		}
		hys, err := hypervisors.ExtractHypervisors(p)
		if err != nil {
			return true, err
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

	return wait.Poll(10*time.Second, 600*time.Second, cfp)
}

func (c *Client) CreateTestInstance(n *model.IronicNode) (err error) {
	c.log.Info("creating node test deployment")
	fID, err := c.getFlavorID("inspection_test")
	iID, err := c.getImageID("ubuntu-20.04-amd64-baremetal")
	zID, err := c.getConductorZone("nova-compute-ironic-testing")
	if err != nil {
		return
	}

	opts := servers.CreateOpts{
		Name:             fmt.Sprintf("%s_inspector_test", time.Now().Format("2006-01-02T15:04:05")),
		FlavorRef:        fID,
		ImageRef:         iID,
		AvailabilityZone: fmt.Sprintf("%s::%s", zID, n.UUID),
	}
	r := servers.Create(c.ComputeClient, opts)
	s, err := r.Extract()
	n.InstanceUUID = s.ID
	if err != nil {
		return
	}

	return servers.WaitForStatus(c.ComputeClient, s.ID, "ACTIVE", 60)
}

func (c *Client) DeleteTestInstance(n *model.IronicNode) (err error) {
	if err = servers.ForceDelete(c.ComputeClient, n.InstanceUUID).ExtractErr(); err != nil {
		return
	}
	return servers.WaitForStatus(c.ComputeClient, n.InstanceUUID, "DELETED", 60)
}

func (c *Client) getImageID(name string) (id string, err error) {
	err = images.ListDetail(c.ComputeClient, images.ListOpts{Name: name}).EachPage(
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
	err = flavors.ListDetail(c.ComputeClient, nil).EachPage(func(p pagination.Page) (bool, error) {
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

func (c *Client) getConductorZone(name string) (id string, err error) {
	err = services.List(c.ComputeClient, services.ListOpts{Host: name}).EachPage(
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

func (c *Client) ProvideNode(n *model.IronicNode) (err error) {
	c.log.Info("providing node")
	cf := func(tp nodes.TargetProvisionState) wait.ConditionFunc {
		return wait.ConditionFunc(func() (bool, error) {
			if err = nodes.ChangeProvisionState(c.IronicClient, n.UUID, nodes.ProvisionStateOpts{
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
		n, err := nodes.Get(c.IronicClient, n.UUID).Extract()
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

func (c *Client) PrepareNode(n *model.IronicNode) (err error) {
	conductor := strings.Split(n.Name, "-")[1]
	opts := nodes.UpdateOpts{
		nodes.UpdateOperation{
			Op:   nodes.RemoveOp,
			Path: "/resource_class",
		},
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

func (c *Client) DeleteNode(n *model.IronicNode) (err error) {
	cfp := wait.ConditionFunc(func() (bool, error) {
		err = nodes.Delete(c.IronicClient, n.UUID).ExtractErr()
		if err != nil {
			return false, err
		}
		return true, nil
	})

	return wait.Poll(5*time.Second, 30*time.Second, cfp)
}

func (c *Client) getRules(n *model.IronicNode) (r config.Rule, err error) {
	var funcMap = template.FuncMap{
		"imageToID": c.getImageID,
	}

	tmpl := template.New("rules.json").Funcs(funcMap)
	t, err := tmpl.ParseFiles(c.cfg.RulesPath)
	if err != nil {
		return r, fmt.Errorf("Error parsing index: %s", err.Error())
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
