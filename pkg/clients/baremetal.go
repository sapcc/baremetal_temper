package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/go-ping/ping"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/baremetal/v1/nodes"
	"github.com/gophercloud/gophercloud/openstack/baremetal/v1/ports"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/hypervisors"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/sapcc/baremetal_temper/pkg/model"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"
)

//NodeAlreadyExists custom error
type NodeAlreadyExists struct {
	Err string
}

func (n *NodeAlreadyExists) Error() string {
	return n.Err
}

//InspectorErr custom error struct for inspector callback errors
type InspectorErr struct {
	Error ErrorMessage `json:"error"`
}

//ErrorMessage message struct for InspectorErr
type ErrorMessage struct {
	Message string `json:"message"`
}

func (c *OpenstackClient) Create() (d []func(n *model.Node) error) {
	d = make([]func(n *model.Node) error, 0)
	d = append(d,
		c.create,
		c.checkCreated,
		c.applyRules,
		c.Validate,
		c.powerOn,
		c.provide,
	)
	return
}

func (c *OpenstackClient) TestAndPrepare() (d []func(n *model.Node) error) {
	d = make([]func(n *model.Node) error, 0)
	d = append(d,
		c.waitForNovaPropagation,
		c.deployTestInstance,
		c.DeleteTestInstance,
		c.prepare,
	)
	return
}

//Create creates a new ironic node based on the provided ironic model
func (c *OpenstackClient) create(in *model.Node) (err error) {
	c.log.Debug("calling inspector api for node creation")
	client := &http.Client{}
	u, err := url.Parse(fmt.Sprintf("http://%s", c.cfg.Inspector.Host))
	if err != nil {
		return
	}
	u.Path = path.Join(u.Path, "/v1/continue")
	db, err := json.Marshal(in.InspectionData)
	if err != nil {
		return
	}
	log.Debugf("calling (%s) with data: %s", u.String(), string(db))
	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(db))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	if res.StatusCode != http.StatusOK {
		ierr := &InspectorErr{}
		if err = json.Unmarshal(bodyBytes, ierr); err != nil {
			return fmt.Errorf("could not create node")
		}
		if strings.Contains(ierr.Error.Message, "already exists, uuid") {
			return &NodeAlreadyExists{}
		}
		return fmt.Errorf(ierr.Error.Message)
	}
	if err = json.Unmarshal(bodyBytes, in); err != nil {
		return
	}
	return
}

//CheckCreated checks if node was created
func (c *OpenstackClient) checkCreated(n *model.Node) error {
	if n.UUID == "" {
		return nil
	}
	c.log.Debug("checking node creation")
	r, err := nodes.Get(c.baremetalClient, n.UUID).Extract()
	if err != nil {
		return &NodeNotFoundError{
			Err: fmt.Sprintf("could not find node %s", n.UUID),
		}
	}
	n.ResourceClass = r.ResourceClass
	return nil
}

//PowerOn powers on the node
func (c *OpenstackClient) powerOn(n *model.Node) (err error) {
	c.log.Debug("powering on node")
	powerStateOpts := nodes.PowerStateOpts{
		Target: nodes.PowerOn,
	}
	r := nodes.ChangePowerState(c.baremetalClient, n.UUID, powerStateOpts)

	if r.Err != nil {
		switch r.Err.(type) {
		case gophercloud.ErrDefault409:
			return fmt.Errorf("cannot power on node %s", n.UUID)
		default:
			return fmt.Errorf("cannot power on node %s", n.UUID)
		}
	}

	cf := wait.ConditionFunc(func() (bool, error) {
		r := nodes.Get(c.baremetalClient, n.UUID)
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

//Validate calls the baremetal validate api
func (c *OpenstackClient) Validate(n *model.Node) (err error) {
	c.log.Debug("validating node")
	v, err := nodes.Validate(c.baremetalClient, n.UUID).Extract()
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

//WaitForNovaPropagation calls the hypervisor api to check if new node has been
//propagated to nova
func (c *OpenstackClient) waitForNovaPropagation(n *model.Node) (err error) {
	c.log.Debug("waiting for nova propagation")
	cfp := wait.ConditionFunc(func() (bool, error) {
		p, err := hypervisors.List(c.computeClient).AllPages()
		if err != nil {
			return false, err
		}
		hys, err := hypervisors.ExtractHypervisors(p)
		if err != nil {
			fmt.Println(err)
			return false, err
		}
		for _, hv := range hys {
			if hv.HypervisorHostname == n.UUID && hv.State == "up" {
				if hv.LocalGB > 0 && hv.MemoryMB > 0 {
					return true, nil
				}
			}
		}
		return false, nil
	})

	return wait.Poll(10*time.Second, 20*time.Minute, cfp)
}

//ApplyRules applies rules from a json file
func (c *OpenstackClient) applyRules(n *model.Node) (err error) {
	c.log.Debug("applying rules on node")
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

//DeployTestInstance creates a new test instance on the newly created node
func (c *OpenstackClient) deployTestInstance(n *model.Node) (err error) {
	c.log.Debug("creating test instance on node")
	iID, err := c.getImageID(c.cfg.Deployment.Image)
	zID, err := c.getConductorZone(c.cfg.Deployment.ConductorZone)
	fID, err := c.getFlavorID(c.cfg.Deployment.Flavor)
	if err != nil {
		return
	}

	net, err := c.getNetwork(c.cfg.Deployment.Network)
	if err != nil {
		return
	}
	nets := make([]servers.Network, 2)
	nets = append(nets, net, net, net)

	pr, err := newProviderClient(c.cfg.Deployment.Openstack)
	if err != nil {
		return
	}
	cc, err := openstack.NewComputeV2(pr, gophercloud.EndpointOpts{
		Region: c.cfg.Region,
	})

	opts := servers.CreateOpts{
		Name:             fmt.Sprintf("%s_temper_test", time.Now().Format("2006-01-02T15:04:05")),
		FlavorRef:        fID,
		ImageRef:         iID,
		AvailabilityZone: fmt.Sprintf("%s::%s", zID, n.UUID),
		Networks:         nets,
	}
	r := servers.Create(cc, opts)
	s, err := r.Extract()
	if err != nil {
		return
	}
	n.InstanceUUID = s.ID
	c.log.Debugf("waiting test instance %s to be created", s.ID)
	instError := true
	if err := servers.WaitForStatus(c.computeClient, s.ID, "ERROR", 60); err != nil {
		if err.Error() == "A timeout occurred" {
			instError = false
		}
	}
	if instError {
		return fmt.Errorf("create test instance %s failed", n.InstanceUUID)
	}
	c.log.Debugf("waiting test instance %s to be active", s.ID)
	if err = servers.WaitForStatus(c.computeClient, s.ID, "ACTIVE", 1200); err != nil {
		return
	}
	n.InstanceIPv4 = s.AccessIPv4
	pinger, err := ping.NewPinger(n.InstanceIPv4)
	if err != nil {
		return
	}
	pinger.Timeout = 1 * time.Minute
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		return
	}
	return
}

//DeleteTestInstance deletes the test instance via the nova api
func (c *OpenstackClient) DeleteTestInstance(n *model.Node) (err error) {
	c.log.Debug("deleting instance on node")
	if err = servers.ForceDelete(c.computeClient, n.InstanceUUID).ExtractErr(); err != nil {
		return
	}
	return servers.WaitForStatus(c.computeClient, n.InstanceUUID, "DELETED", 60)
}

//Provide sets node provisionstate to provided (available).
//Needed to deploy a test instance on this node
func (c *OpenstackClient) provide(n *model.Node) (err error) {
	c.log.Debug("providing node")
	cf := func(tp nodes.TargetProvisionState) wait.ConditionFunc {
		return wait.ConditionFunc(func() (bool, error) {
			if err = nodes.ChangeProvisionState(c.baremetalClient, n.UUID, nodes.ProvisionStateOpts{
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
		n, err := nodes.Get(c.baremetalClient, n.UUID).Extract()
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

func (c *OpenstackClient) CreatePortGroup(n *model.Node) (err error) {
	//TODO: create port group
	return
}

//Prepare prepares the node for customers.
//Removes resource_class, sets the rightful conductor and maintenance to true
func (c *OpenstackClient) prepare(n *model.Node) (err error) {
	c.log.Debug("preparing node")
	conductor := strings.Split(n.Name, "-")[1]
	opts := nodes.UpdateOpts{
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

//DeleteNode deletes a node via the baremetal api
func (c *OpenstackClient) DeleteNode(n *model.Node) (err error) {
	if n.UUID == "" {
		return
	}
	c.log.Debug("deleting node")
	cfp := wait.ConditionFunc(func() (bool, error) {
		err = nodes.Delete(c.baremetalClient, n.UUID).ExtractErr()
		if err != nil {
			return false, err
		}
		return true, nil
	})

	return wait.Poll(5*time.Second, 30*time.Second, cfp)
}

func (c *OpenstackClient) updatePorts(opts ports.UpdateOpts, n *model.Node) (err error) {
	listOpts := ports.ListOpts{
		NodeUUID: n.UUID,
	}

	l, err := ports.List(c.baremetalClient, listOpts).AllPages()
	if err != nil {
		return
	}

	ps, err := ports.ExtractPorts(l)
	if err != nil {
		return
	}

	for _, p := range ps {
		cf := wait.ConditionFunc(func() (bool, error) {
			_, err = ports.Update(c.baremetalClient, p.UUID, opts).Extract()
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

func (c *OpenstackClient) updateNode(opts nodes.UpdateOpts, n *model.Node) (err error) {
	cf := wait.ConditionFunc(func() (bool, error) {
		r := nodes.Update(c.baremetalClient, n.UUID, opts)
		_, err = r.Extract()
		if err != nil {
			return false, nil
		}
		return true, nil
	})
	return wait.Poll(5*time.Second, 60*time.Second, cf)
}
