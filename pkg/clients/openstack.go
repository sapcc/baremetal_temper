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
	"github.com/gophercloud/gophercloud/openstack/dns/v2/recordsets"
	"github.com/gophercloud/gophercloud/openstack/dns/v2/zones"
	"github.com/gophercloud/gophercloud/pagination"
	"k8s.io/apimachinery/pkg/util/wait"
)

type Client struct {
	IronicClient *gophercloud.ServiceClient
	DnsClient    *gophercloud.ServiceClient
	Domain       string
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

	if err != nil {
		return nil, err
	}
	version, err := apiversions.Get(iclient, "v1").Extract()
	if err != nil {
		return nil, err
	}
	iclient.Microversion = version.Version
	return &Client{IronicClient: iclient, DnsClient: dnsClient, Domain: domain}, nil
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

func (c *Client) UpdateNode(n model.IronicNode) (err error) {
	updateName := nodes.UpdateOpts{
		nodes.UpdateOperation{
			Op:    nodes.ReplaceOp,
			Path:  "/name",
			Value: n.Name,
		},
	}

	return c.updateNode(updateName, n)
}

func (c *Client) updateNode(opts nodes.UpdateOpts, n model.IronicNode) (err error) {
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

func (c *Client) PowerNodeOn(uuid string) (err error) {
	powerStateOpts := nodes.PowerStateOpts{
		Target: nodes.PowerOn,
	}
	r := nodes.ChangePowerState(c.IronicClient, uuid, powerStateOpts)

	if r.Err != nil {
		switch r.Err.(type) {
		case gophercloud.ErrDefault409:
			//p.log.Info("host is locked, trying again after delay", "delay", powerRequeueDelay)
			return fmt.Errorf("cannot power on node %s", uuid)
		default:
			return fmt.Errorf("cannot power on node %s", uuid)
		}
	}

	cf := wait.ConditionFunc(func() (bool, error) {
		r := nodes.Get(c.IronicClient, uuid)
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

func (c *Client) listNodes() (l []nodes.Node, err error) {
	pages := 0
	opts := nodes.ListOpts{}
	err = nodes.List(c.IronicClient, opts).EachPage(func(p pagination.Page) (bool, error) {
		pages++
		extracted, err := nodes.ExtractNodes(p)
		if err != nil {
			return false, err
		}
		l = append(l, extracted...)
		return true, nil
	})
	return
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
