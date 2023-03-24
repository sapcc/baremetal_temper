package clients

import (
	"fmt"
	"os"

	"github.com/sapcc/baremetal_temper/pkg/config"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/baremetal/apiversions"
	iDservices "github.com/gophercloud/gophercloud/openstack/identity/v3/services"
	log "github.com/sirupsen/logrus"
)

// Openstack is
type Openstack struct {
	Clients map[string]*gophercloud.ServiceClient
	log     *log.Entry
	cfg     config.Config
}

type PortGroup struct {
	UUID                     string                 `json:"uuid,omitempty"`
	NodeUUID                 string                 `json:"node_uuid"`
	Name                     string                 `json:"name,omitempty"`
	Address                  string                 `json:"address,omitempty"`
	StandalonePortsSupported bool                   `json:"standalone_ports_supported,omitempty"`
	Mode                     string                 `json:"mode,omitempty"`
	Properties               map[string]interface{} `json:"properties,omitempty"`
}

type Console struct {
	Enabled bool `json:"enabled"`
}

// ToPortCreateMap assembles a request body based on the contents of a CreateOpts.
func (opts PortGroup) ToPortCreateMap() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return body, nil
}

// NodeNotFoundError error for missing node
type NodeNotFoundError struct {
	Err string
}

func (n *NodeNotFoundError) Error() string {
	return n.Err
}

// NewClient creates a new client containing different openstack-clients (baremetal, compute, dns)
func NewClient(cfg config.Config, ctxLogger *log.Entry) *Openstack {
	return &Openstack{cfg: cfg, log: ctxLogger, Clients: make(map[string]*gophercloud.ServiceClient, 0)}
}

func (oc *Openstack) GetServiceClient(client string) (c *gophercloud.ServiceClient, err error) {
	c, ok := oc.Clients[client]
	if ok {
		return
	}
	provider, err := NewProviderClient(oc.cfg.Openstack)
	if err != nil {
		return nil, err
	}
	switch client {
	case "compute":
		c, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
			Region: oc.cfg.Region,
		})
		if err != nil {
			return nil, err
		}
		oc.Clients["compute"] = c
		return c, err
	case "dns":
		c, err := openstack.NewDNSV2(provider, gophercloud.EndpointOpts{
			Region: oc.cfg.Region,
		})
		if err != nil {
			return nil, err
		}
		oc.Clients["dns"] = c
		return c, err
	case "identity":
		c, err := openstack.NewIdentityV3(provider, gophercloud.EndpointOpts{
			Region: oc.cfg.Region,
		})
		if err != nil {
			return nil, err
		}
		oc.Clients["identity"] = c
		return c, err
	case "network":
		c, err := openstack.NewNetworkV2(provider, gophercloud.EndpointOpts{
			Region: oc.cfg.Region,
		})
		if err != nil {
			return nil, err
		}
		oc.Clients["network"] = c
		return c, err
	case "object":
		c, err := openstack.NewObjectStorageV1(provider, gophercloud.EndpointOpts{
			Region: oc.cfg.Region,
		})
		if err != nil {
			return nil, err
		}
		oc.Clients["object"] = c
		return c, err
	case "baremetal":
		c, err := openstack.NewBareMetalV1(provider, gophercloud.EndpointOpts{
			Region: oc.cfg.Region,
		})
		if err == nil {
			version, err := apiversions.Get(c, "v1").Extract()
			if err != nil {
				return nil, err
			}
			if version.Version == "" {
				c.Microversion = "1.38"
			} else {
				c.Microversion = version.Version
			}
			oc.Clients["baremetal"] = c
			return c, err
		}
		oc.log.Infof("baremetal service error: %s", err.Error())
		return c, err
	}
	return
}

func NewProviderClient(i config.OpenstackAuth) (pc *gophercloud.ProviderClient, err error) {
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

func (c *Openstack) ServiceEnabled(service string) (bool, error) {
	a, err := iDservices.List(c.Clients["identity"], iDservices.ListOpts{Name: service}).AllPages()
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

func (c *Openstack) getAPIVersion() (*apiversions.APIVersion, error) {
	return apiversions.Get(c.Clients["baremetal"], "v1").Extract()
}
