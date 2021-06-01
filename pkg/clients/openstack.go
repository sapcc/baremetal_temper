package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/sapcc/baremetal_temper/pkg/config"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/baremetal/apiversions"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/services"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/images"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/dns/v2/recordsets"
	"github.com/gophercloud/gophercloud/openstack/dns/v2/zones"
	iDservices "github.com/gophercloud/gophercloud/openstack/identity/v3/services"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/pagination"
	log "github.com/sirupsen/logrus"
)

//Openstack is
type Openstack struct {
	Clients map[string]*gophercloud.ServiceClient
	log     *log.Entry
	cfg     config.Config
}

type PortGroup struct {
	UUID                     string                 `json:"uuid"`
	NodeUUID                 string                 `json:"node_uuid"`
	Name                     string                 `json:"name,omitempty"`
	Address                  string                 `json:"address,omitempty"`
	StandalonePortsSupported bool                   `json:"standalone_ports_supported,omitempty"`
	Mode                     string                 `json:"mode,omitempty"`
	Properties               map[string]interface{} `json:"properties,omitempty"`
}

// ToPortCreateMap assembles a request body based on the contents of a CreateOpts.
func (opts PortGroup) ToPortCreateMap() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return body, nil
}

//NodeNotFoundError error for missing node
type NodeNotFoundError struct {
	Err string
}

func (n *NodeNotFoundError) Error() string {
	return n.Err
}

//NewClient creates a new client containing different openstack-clients (baremetal, compute, dns)
func NewClient(cfg config.Config, ctxLogger *log.Entry) *Openstack {
	return &Openstack{cfg: cfg, log: ctxLogger, Clients: make(map[string]*gophercloud.ServiceClient, 0)}
}

func (oc *Openstack) GetServiceClient(cfg config.Config, client string) (c *gophercloud.ServiceClient, err error) {
	c, ok := oc.Clients[client]
	if ok {
		return
	}
	provider, err := NewProviderClient(cfg.Openstack)
	if err != nil {
		return nil, err
	}
	switch client {
	case "compute":
		c, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
			Region: cfg.Region,
		})
		if err != nil {
			return nil, err
		}
		oc.Clients["compute"] = c
		return c, err
	case "dns":
		c, err := openstack.NewDNSV2(provider, gophercloud.EndpointOpts{
			Region: cfg.Region,
		})
		if err != nil {
			return nil, err
		}
		oc.Clients["dns"] = c
		return c, err
	case "identity":
		c, err := openstack.NewIdentityV3(provider, gophercloud.EndpointOpts{
			Region: cfg.Region,
		})
		if err != nil {
			return nil, err
		}
		oc.Clients["identity"] = c
		return c, err
	case "network":
		c, err := openstack.NewNetworkV2(provider, gophercloud.EndpointOpts{
			Region: cfg.Region,
		})
		if err != nil {
			return nil, err
		}
		oc.Clients["network"] = c
		return c, err
	case "baremetal":
		c, err := openstack.NewBareMetalV1(provider, gophercloud.EndpointOpts{
			Region: cfg.Region,
		})
		if err == nil {
			version, err := apiversions.Get(c, "v1").Extract()
			if err != nil {
				return nil, err
			}
			c.Microversion = version.Version
			oc.Clients["baremetal"] = c
			return c, err
		} else {
			oc.log.Infof("baremetal service error: %s", err.Error())
			return c, err
		}
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

func (c *Openstack) CreateArpaZone(ip string) (zoneID string, err error) {
	arpaZone, err := reverseZone(ip)
	if err != nil {
		return
	}

	allPages, err := zones.List(c.Clients["dns"], zones.ListOpts{
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
		z, err := zones.Create(c.Clients["dns"], zones.CreateOpts{
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

func (c *Openstack) CreateDNSRecord(ip, zoneID, recordName, rType string) (err error) {
	_, err = recordsets.Create(c.Clients["dns"], zoneID, recordsets.CreateOpts{
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

func (c *Openstack) GetImageID(name string) (id string, err error) {
	err = images.ListDetail(c.Clients["compute"], images.ListOpts{Name: name}).EachPage(
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

func (c *Openstack) GetFlavorID(name string) (id string, err error) {
	err = flavors.ListDetail(c.Clients["compute"], nil).EachPage(func(p pagination.Page) (bool, error) {
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

func (c *Openstack) GetConductorZone(name string) (id string, err error) {
	err = services.List(c.Clients["compute"], services.ListOpts{Host: name}).EachPage(
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

func (c *Openstack) GetNetwork(name string) (n servers.Network, err error) {
	pr, err := NewProviderClient(c.cfg.Deployment.Openstack)
	if err != nil {
		return
	}
	nc, err := openstack.NewNetworkV2(pr, gophercloud.EndpointOpts{
		Region: c.cfg.Region,
	})
	p, err := networks.List(nc, networks.ListOpts{Name: name}).AllPages()
	if err != nil {
		return
	}
	ns, err := networks.ExtractNetworks(p)
	if err != nil || len(ns) != 1 {
		return
	}
	n.UUID = ns[0].ID
	return
}

func (c *Openstack) CreatePortGroup(pg PortGroup) (uuid string, err error) {
	u := c.Clients["baremetal"].ServiceURL("ports")
	reqBody, err := pg.ToPortCreateMap()
	if err != nil {
		return
	}
	resp, err := c.Clients["baremetal"].Post(u, reqBody, nil, nil)
	if resp.StatusCode != http.StatusCreated {
		return uuid, fmt.Errorf("error creating port group: %s", err.Error())
	}
	r := PortGroup{}
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return
	}
	return r.UUID, err
}
