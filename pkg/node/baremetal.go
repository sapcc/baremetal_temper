/**
 * Copyright 2021 SAP SE
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package node

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
	"github.com/sapcc/baremetal_temper/pkg/clients"
	"k8s.io/apimachinery/pkg/util/wait"
)

//NodeAlreadyExists custom error
type AlreadyExists struct {
	Err string
}

func (n *AlreadyExists) Error() string {
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

//Create creates a new ironic node based on the provided ironic model
func (n *Node) create() (err error) {
	n.log.Debug("calling inspector api for node creation")
	if len(n.InspectionData.Inventory.Interfaces) == 0 {
		panic("no interfaces with linkStatus up found. cannot create ironic node")
	}
	client := &http.Client{Timeout: 90 * time.Second}
	u, err := url.Parse(fmt.Sprintf("http://%s", n.cfg.Inspector.Host))
	if err != nil {
		panic("could not create ironic node: " + err.Error())
	}
	u.Path = path.Join(u.Path, "/v1/continue")
	db, err := json.Marshal(n.InspectionData)
	if err != nil {
		panic("could not create ironic node: " + err.Error())
	}
	n.log.Debugf("calling (%s) with data: %s", u.String(), string(db))
	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(db))
	if err != nil {
		panic("could not create ironic node: " + err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		panic("could not create ironic node: " + err.Error())
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic("could not create ironic node: " + err.Error())
	}

	if res.StatusCode != http.StatusOK {
		ierr := &InspectorErr{}
		if err = json.Unmarshal(bodyBytes, ierr); err != nil {
			panic("could not create ironic node: " + err.Error())
		}
		if strings.Contains(ierr.Error.Message, "already exists, uuid") {
			return &AlreadyExists{}
		}
		panic("could not create ironic node: " + ierr.Error.Message)
	}
	if err = json.Unmarshal(bodyBytes, n); err != nil {
		panic("could not create ironic node: " + err.Error())
	}
	return
}

//DeleteNode deletes a node via the baremetal api
func (n *Node) DeleteNode() (err error) {
	if n.UUID == "" {
		return
	}
	c, err := n.oc.GetServiceClient("baremetal")
	if err != nil {
		return
	}
	n.log.Debug("deleting node")
	cfp := wait.ConditionFunc(func() (bool, error) {
		err = nodes.Delete(c, n.UUID).ExtractErr()
		if err != nil {
			return false, err
		}
		return true, nil
	})

	return wait.Poll(5*time.Second, 30*time.Second, cfp)
}

//CheckCreated checks if node was created
func (n *Node) checkCreated() (err error) {
	if n.UUID == "" {
		return
	}
	c, err := n.oc.GetServiceClient("baremetal")
	if err != nil {
		return
	}
	n.log.Debug("checking node creation")
	r, err := nodes.Get(c, n.UUID).Extract()
	if err != nil {
		return &clients.NodeNotFoundError{
			Err: fmt.Sprintf("could not find node %s", n.UUID),
		}
	}
	n.ResourceClass = r.ResourceClass
	return
}

//Prepare prepares the node for customers.
//Removes resource_class, sets the rightful conductor and maintenance to true
func (n *Node) prepare() (err error) {
	if err = n.getUUID(); err != nil {
		return
	}
	n.log.Debug("preparing node")
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
	return n.updateNode(opts)
}

//PowerOn powers on the node
func (n *Node) changePowerState(powerState nodes.TargetPowerState) (err error) {
	c, err := n.oc.GetServiceClient("baremetal")
	if err != nil {
		return
	}
	n.log.Debug("powering on node")
	powerStateOpts := nodes.PowerStateOpts{
		Target: powerState,
	}
	r := nodes.ChangePowerState(c, n.UUID, powerStateOpts)

	if r.Err != nil {
		switch r.Err.(type) {
		case gophercloud.ErrDefault409:
			return fmt.Errorf("cannot power on node %s", n.UUID)
		default:
			return fmt.Errorf("cannot power on node %s", n.UUID)
		}
	}

	cf := wait.ConditionFunc(func() (bool, error) {
		r := nodes.Get(c, n.UUID)
		n, err := r.Extract()
		if err != nil {
			return false, fmt.Errorf("cannot power on node")
		}
		if n.PowerState != string(nodes.PowerOn) {
			return false, nil
		}
		return true, nil
	})
	return wait.Poll(5*time.Second, 120*time.Second, cf)
}

//PowerOn powers on the node
func (n *Node) powerOn() (err error) {
	return n.changePowerState(nodes.PowerOn)
}

//PowerOff node off
func (n *Node) powerOff() (err error) {
	return n.changePowerState(nodes.PowerOff)
}

//Validate calls the baremetal validate api
func (n *Node) validate() (err error) {
	c, err := n.oc.GetServiceClient("baremetal")
	if err != nil {
		return
	}
	n.log.Debug("validating node")
	if err = n.getUUID(); err != nil {
		return
	}
	v, err := nodes.Validate(c, n.UUID).Extract()
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

//DeleteTestInstance deletes the test instance via the nova api
func (n *Node) DeleteTestInstance() (err error) {
	c, err := n.oc.GetServiceClient("compute")
	if err != nil {
		return
	}
	n.log.Debug("deleting instance on node")
	if err = servers.ForceDelete(c, n.InstanceUUID).ExtractErr(); err != nil {
		return
	}
	return servers.WaitForStatus(c, n.InstanceUUID, "DELETED", 60)
}

//Provide sets node provisionstate to provided (available).
//Needed to deploy a test instance on this node
func (n *Node) provide() (err error) {
	c, err := n.oc.GetServiceClient("baremetal")
	if err != nil {
		return
	}
	n.log.Debug("providing node")
	cf := func(tp nodes.TargetProvisionState) wait.ConditionFunc {
		return wait.ConditionFunc(func() (bool, error) {
			if err = nodes.ChangeProvisionState(c, n.UUID, nodes.ProvisionStateOpts{
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
		n, err := nodes.Get(c, n.UUID).Extract()
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

func (n *Node) getUUID() (err error) {
	if n.UUID != "" {
		return
	}
	c, err := n.oc.GetServiceClient("baremetal")
	if err != nil {
		return
	}
	n.log.Debug("get node uuid")
	p, err := nodes.List(c, nodes.ListOpts{}).AllPages()
	if err != nil {
		return &clients.NodeNotFoundError{
			Err: fmt.Sprintf("could not find node %s uuid", n.Name),
		}
	}
	nodes, err := nodes.ExtractNodes(p)
	if err != nil {
		return
	}
	for _, no := range nodes {
		if no.Name == n.Name {
			n.UUID = no.UUID
			break
		}
	}
	if n.UUID == "" {
		return &clients.NodeNotFoundError{
			Err: fmt.Sprintf("could not find node %s uuid", n.Name),
		}
	}
	return
}

//WaitForNovaPropagation calls the hypervisor api to check if new node has been
//propagated to nova
func (n *Node) waitForNovaPropagation() (err error) {
	c, err := n.oc.GetServiceClient("compute")
	if err != nil {
		return
	}
	if err = n.getUUID(); err != nil {
		return
	}
	n.log.Debug("waiting for nova propagation")
	cfp := wait.ConditionFunc(func() (bool, error) {
		p, err := hypervisors.List(c).AllPages()
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
func (n *Node) applyRules() (err error) {
	n.log.Debug("applying rules on node")
	rules, err := n.getRules()
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
	if err = n.updatePorts(updatePorts); err != nil {
		return
	}

	return n.updateNode(updateNode)
}

//DeployTestInstance creates a new test instance on the newly created node
func (n *Node) deployTestInstance() (err error) {
	c, err := n.oc.GetServiceClient("compute")
	if err != nil {
		return
	}
	n.log.Debug("creating test instance on node")
	iID, err := n.oc.GetImageID(n.cfg.Deployment.Image)
	if err != nil {
		return
	}
	zID, err := n.oc.GetConductorZone(n.cfg.Deployment.ConductorZone)
	if err != nil {
		return
	}
	fID, err := n.oc.GetFlavorID(n.cfg.Deployment.Flavor)
	if err != nil {
		return
	}

	net, err := n.oc.GetNetwork(n.cfg.Deployment.Network)
	if err != nil {
		return
	}
	nets := make([]servers.Network, 0, 2)
	nets = append(nets, net, net, net)

	pr, err := clients.NewProviderClient(n.cfg.Deployment.Openstack)
	if err != nil {
		return
	}
	cc, err := openstack.NewComputeV2(pr, gophercloud.EndpointOpts{
		Region: n.cfg.Region,
	})

	opts := servers.CreateOpts{
		Name:             fmt.Sprintf("%s_node_test", time.Now().Format("2006-01-02T15:04:05")),
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
	n.log.Debugf("waiting test instance %s to be created", s.ID)
	instError := true
	if err := servers.WaitForStatus(c, s.ID, "ERROR", 60); err != nil {
		if err.Error() == "A timeout occurred" {
			instError = false
		}
	}
	if instError {
		return fmt.Errorf("create test instance %s failed", n.InstanceUUID)
	}
	n.log.Debugf("waiting test instance %s to be active", s.ID)
	if err = servers.WaitForStatus(c, s.ID, "ACTIVE", 1200); err != nil {
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

func (n *Node) createPortGroup(name string) (id string, err error) {
	if n.PortGroupUUID != "" {
		return n.PortGroupUUID, err
	}
	pg := clients.PortGroup{
		NodeUUID:                 n.UUID,
		StandalonePortsSupported: true,
		Name:                     name,
		Address:                  n.InspectionData.Inventory.Interfaces[0].MacAddress, // use the MAC of the first interface
		Mode:                     "802.3ad",
		Properties:               map[string]interface{}{"miimon": 100},
	}
	id, err = n.oc.CreatePortGroup(pg)
	if err != nil {
		return
	}
	n.PortGroupUUID = id
	return
}

func (n *Node) updatePorts(opts ports.UpdateOpts) (err error) {
	c, err := n.oc.GetServiceClient("baremetal")
	if err != nil {
		return
	}
	listOpts := ports.ListOpts{
		NodeUUID: n.UUID,
	}

	l, err := ports.List(c, listOpts).AllPages()
	if err != nil {
		return
	}

	ps, err := ports.ExtractPorts(l)
	if err != nil {
		return
	}

	for _, p := range ps {
		cf := wait.ConditionFunc(func() (bool, error) {
			_, err = ports.Update(c, p.UUID, opts).Extract()
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

func (n *Node) updateNode(opts nodes.UpdateOpts) (err error) {
	c, err := n.oc.GetServiceClient("baremetal")
	if err != nil {
		return
	}
	cf := wait.ConditionFunc(func() (bool, error) {
		r := nodes.Update(c, n.UUID, opts)
		_, err = r.Extract()
		if err != nil {
			return false, nil
		}
		return true, nil
	})
	return wait.Poll(5*time.Second, 60*time.Second, cf)
}
