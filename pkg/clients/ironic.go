package clients

import (
	"fmt"
	"os"

	"github.com/sapcc/ironic_temper/pkg/config"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/baremetal/apiversions"
	"github.com/gophercloud/gophercloud/openstack/baremetal/v1/nodes"
	"github.com/gophercloud/gophercloud/pagination"
)

type Client struct {
	Client *gophercloud.ServiceClient
}

func NewClient(region string, i config.IronicAuth) (*Client, error) {
	provider, err := newProviderClient(i)
	if err != nil {
		return nil, err
	}
	client, err := openstack.NewBareMetalV1(provider, gophercloud.EndpointOpts{
		Region: region,
	})
	if err != nil {
		return nil, err
	}
	version, err := apiversions.Get(client, "v1").Extract()
	if err != nil {
		return nil, err
	}
	client.Microversion = version.Version
	return &Client{client}, nil
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
	r := nodes.Get(c.Client, uuid)
	return r.Extract()
}

func (c *Client) SetNodeName(uuid, name string) (err error) {
	updateOpts := nodes.UpdateOpts{
		nodes.UpdateOperation{
			Op:    nodes.ReplaceOp,
			Path:  "/name",
			Value: name,
		},
	}
	r := nodes.Update(c.Client, uuid, updateOpts)
	_, err = r.Extract()
	return
}

func (c *Client) PowerNodeOn(uuid string) (err error) {
	powerStateOpts := nodes.PowerStateOpts{
		Target: nodes.PowerOn,
	}
	r := nodes.ChangePowerState(c.Client, uuid, powerStateOpts)
	switch r.Err.(type) {
	case nil:
		return
	case gophercloud.ErrDefault409:
		//p.log.Info("host is locked, trying again after delay", "delay", powerRequeueDelay)
		return fmt.Errorf("cannot power on node %s", uuid)
	default:
		return fmt.Errorf("cannot power on node %s", uuid)
	}
}

func (c *Client) listNodes() (l []nodes.Node, err error) {
	pages := 0
	opts := nodes.ListOpts{}
	err = nodes.List(c.Client, opts).EachPage(func(p pagination.Page) (bool, error) {
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
	return nodes.Get(c.Client, id).Extract()
}

func (c *Client) getAPIVersion() (*apiversions.APIVersion, error) {
	return apiversions.Get(c.Client, "v1").Extract()
}
