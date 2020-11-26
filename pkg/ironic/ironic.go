package ironic

import "github.com/gophercloud/gophercloud/openstack"

type Ironic struct {
}

func NewIronic() (*Ironic, error) {
	authOpts, err := openstack.AuthOptionsFromEnv()
	if err != nil {
		return nil, err
	}
	provider, err := openstack.AuthenticatedClient(authOpts)
	if err != nil {
		return nil, err
	}

}
