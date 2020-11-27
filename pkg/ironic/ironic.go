package ironic

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
)

type Ironic struct {
}

func NewIronic() (client *gophercloud.ProviderClient, err error) {
	authOpts, err := openstack.AuthOptionsFromEnv()
	if err != nil {
		return nil, err
	}
	client, err = openstack.AuthenticatedClient(authOpts)
	if err != nil {
		return nil, err
	}
	return
}
