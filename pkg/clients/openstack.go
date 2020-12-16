package clients

import (
	"fmt"
	"os"
	"strings"
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
	"k8s.io/apimachinery/pkg/util/wait"
)

type Client struct {
	IronicClient  *gophercloud.ServiceClient
	DnsClient     *gophercloud.ServiceClient
	ComputeClient *gophercloud.ServiceClient
	Domain        string
}

type NodeNotFoundError struct {
	Err string
}

func (n *NodeNotFoundError) Error() string {
	return n.Err
}

func NewClient(region string, i config.IronicAuth, domain string) (*Client, error) {
	provider, err := newProviderClient(i)
	if err != nil {
		return nil, err
	}
	iclient, err := openstack.NewBareMetalV1(provider, gophercloud.EndpointOpts{
		Region: region,
	})

	dnsClient, err := openstack.NewDNSV2(provider, gophercloud.EndpointOpts{
		Region: region,
	})

	cclient, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: region,
	})

	if err != nil {
		return nil, err
	}
	version, err := apiversions.Get(iclient, "v1").Extract()
	if err != nil {
		return nil, err
	}
	iclient.Microversion = version.Version
	return &Client{IronicClient: iclient, DnsClient: dnsClient, ComputeClient: cclient, Domain: domain}, nil
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

// GetNodeUUIDByName gets node's uuid by node name
func (c *Client) GetNodeByUUID(uuid string) (node *nodes.Node, err error) {
	r := nodes.Get(c.IronicClient, uuid)
	return r.Extract()
}

func (c *Client) CheckIronicNodeCreated(n *model.IronicNode) error {
	if n.UUID != "" {
		return nil
	}
	_, err := c.GetNodeByUUID(n.UUID)
	if err != nil {
		return &NodeNotFoundError{
			Err: fmt.Sprintf("could not find node %s", n.UUID),
		}
	}
	return nil
}

func (c *Client) UpdateNode(n *model.IronicNode) (err error) {
	kernel, err := c.getImageID("")
	ramdisk, err := c.getImageID("")
	if err != nil {
		return
	}
	updateNode := nodes.UpdateOpts{
		nodes.UpdateOperation{
			Op:    nodes.ReplaceOp,
			Path:  "/name",
			Value: n.Name,
		},
		nodes.UpdateOperation{
			Op:    nodes.ReplaceOp,
			Path:  "/resource_class",
			Value: "inspection_test",
		},
		nodes.UpdateOperation{
			Op:    nodes.AddOp,
			Path:  "/driver_info/deploy_kernel",
			Value: kernel,
		},
		nodes.UpdateOperation{
			Op:    nodes.AddOp,
			Path:  "/driver_info/deploy_ramdisk",
			Value: ramdisk,
		},
		nodes.UpdateOperation{
			Op:    nodes.AddOp,
			Path:  "/properties/manufacturer",
			Value: n.InspectionData.Inventory.SystemVendor.Manufacturer,
		},
		nodes.UpdateOperation{
			Op:    nodes.AddOp,
			Path:  "/properties/serial",
			Value: n.InspectionData.Inventory.SystemVendor.SerialNumber,
		},
		nodes.UpdateOperation{
			Op:    nodes.AddOp,
			Path:  "/properties/model",
			Value: n.InspectionData.Inventory.SystemVendor.ProductName,
		},
	}
	updatePorts := ports.UpdateOpts{
		ports.UpdateOperation{
			Op:    ports.AddOp,
			Path:  "/local_link_connection/switch_id",
			Value: "aa:bb:cc:dd:ee:ff",
		},
		ports.UpdateOperation{
			Op:    ports.AddOp,
			Path:  "/local_link_connection/port_id",
			Value: "Etherneth1/15",
		},
		ports.UpdateOperation{
			Op:    ports.AddOp,
			Path:  "/local_link_connection/switch_info",
			Value: n.Name,
		},
		ports.UpdateOperation{
			Op:    ports.ReplaceOp,
			Path:  "/pxe_enabled",
			Value: true,
		},
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
		if err != nil {
			return false, nil
		}
		return true, nil
	})
	return wait.Poll(5*time.Second, 60*time.Second, cf)
}

func (c *Client) PowerNodeOn(n *model.IronicNode) (err error) {
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

func (c *Client) getNodeByID(id string) (n *nodes.Node, err error) {
	return nodes.Get(c.IronicClient, id).Extract()
}

func (c *Client) getAPIVersion() (*apiversions.APIVersion, error) {
	return apiversions.Get(c.IronicClient, "v1").Extract()
}

func (c *Client) CreateDNSRecordFor(n *model.IronicNode) (err error) {
	opts := zones.ListOpts{
		Name: c.Domain + ".",
	}
	allPages, err := zones.List(c.DnsClient, opts).AllPages()
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

	_, err = recordsets.Create(c.DnsClient, allZones[0].ID, recordsets.CreateOpts{
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

func (c *Client) ValidateNode(n *model.IronicNode) (err error) {
	if err = c.provideNode(n.UUID); err != nil {
		return
	}
	v, err := nodes.Validate(c.IronicClient, n.UUID).Extract()
	if !v.Inspect.Result {
		return fmt.Errorf(v.Inspect.Reason)
	}
	if !v.Power.Result {
		return fmt.Errorf(v.Power.Reason)
	}
	return
}

func (c *Client) WaitForNovaPropagation(n *model.IronicNode) (err error) {
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

func (c *Client) CreateNodeTestDeployment(n *model.IronicNode) (err error) {
	fID, err := c.getFlavorID("")
	iID, err := c.getImageID("")
	zID, err := c.getConductorZone("")
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

	return servers.WaitForStatus(c.ComputeClient, s.ID, "ACTIVE", 600)
}

func (c *Client) DeleteNodeTestDeployment(n *model.IronicNode) (err error) {
	return servers.ForceDelete(c.ComputeClient, n.InstanceUUID).ExtractErr()
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

func (c *Client) provideNode(id string) (err error) {
	cf := func(tp nodes.TargetProvisionState) wait.ConditionFunc {
		return wait.ConditionFunc(func() (bool, error) {
			if err = nodes.ChangeProvisionState(c.IronicClient, id, nodes.ProvisionStateOpts{
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
		n, err := nodes.Get(c.IronicClient, id).Extract()
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
