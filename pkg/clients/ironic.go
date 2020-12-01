package clients

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/baremetal/v1/nodes"
	"github.com/gophercloud/gophercloud/pagination"
)

type Client struct {
	Client *gophercloud.ServiceClient
}

func NewClient(region string) (*Client, error) {
	authOpts, err := openstack.AuthOptionsFromEnv()
	if err != nil {
		return nil, err
	}
	provider, err := openstack.AuthenticatedClient(authOpts)
	if err != nil {
		return nil, err
	}
	client, err := openstack.NewBareMetalV1(provider, gophercloud.EndpointOpts{
		Region: region,
	})
	if err != nil {
		return nil, err
	}
	return &Client{client}, nil
}

// GetNodeUUIDByName gets node's uuid by node name
func (c *Client) GetNodeUUIDByName(name string) (nodeUUID string, err error) {
	nodeList, err := c.listNodes()
	for _, _node := range nodeList {
		node, err := c.getNodeByID(_node.UUID)
		if err != nil {
			return "", err
		}
		fmt.Printf("%v\n\n", node)
	}
	return
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
