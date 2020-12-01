package ironic

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/baremetal/v1/nodes"
	"github.com/gophercloud/gophercloud/pagination"
)

type Client struct {
	Client *gophercloud.ServiceClient
}

type Node struct {
	Name           string
	IP             string
	UUID           string
	Region         string
	IronicUser     string
	IronicPassword string
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
	return &Client{client}, nil
}

// GetNodeUUIDByName gets node's uuid by node name
func (c *Client) GetNodeUUIDByName(name string) (nodeUUID string, err error) {
	pages := 0
	err = nodes.List(c.Client, nodes.ListOpts{}).EachPage(func(p pagination.Page) (bool, error) {
		pages++
		extracted, err := nodes.ExtractNodes(p)
		if err != nil {
			return false, err
		}
		for _, n := range extracted {
			if n.Name == name {
				nodeUUID = n.UUID
				return false, nil
			}
		}
		return true, nil
	})
	return
}
