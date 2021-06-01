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

import "time"

func (n *Node) AddBaremetalCreateTasks() {
	n.Tasks[100] = &Task{Name: "create_ironic_node", Exec: n.Create}
	n.Tasks[90] = &Task{Name: "check_ironic_node_created", Exec: n.CheckCreated}
	n.Tasks[80] = &Task{Name: "apply_ironic_rules", Exec: n.ApplyRules}
	n.Tasks[70] = &Task{Name: "validate_ironic_node", Exec: n.Validate}
	n.Tasks[60] = &Task{Name: "power_on_ironic_node", Exec: n.PowerOn}
	n.Tasks[50] = &Task{Name: "provide_ironic_node", Exec: n.Provide}
	n.Tasks[40] = &Task{Name: "power_off_ironic_node", Exec: n.PowerOff}
}

func (n *Node) AddDeploymentTestTasks() {
	n.Tasks[100] = &Task{Name: "get_ironic_uuid", Exec: n.getUUID}
	n.Tasks[90] = &Task{Name: "wait_nova_propagation", Exec: n.WaitForNovaPropagation}
	n.Tasks[80] = &Task{Name: "deploy_test_instance", Exec: n.DeployTestInstance}
}

func (n *Node) AddTask(p int, na string) (t *Task) {
	if n.Tasks == nil {
		n.Tasks = make(map[int]*Task, 0)
	}
	t, ok := n.Tasks[p]
	if !ok {
		n.Tasks[p] = &Task{Name: na}
	}
	return n.Tasks[p]
}

func (n *Node) AddAllTemperTasks(diag bool, bm bool, events bool, image bool) {
	n.AddTask(0, "create_dns_records").Exec = n.CreateDNSRecords
	if events {
		n.AddTask(10, "create_event_sub").Exec = n.CreateEventSubscription
	}
	if diag {
		n.AddTask(20, "hardware_check").Exec = n.RunHardwareChecks

		if n.cfg.Redfish.BootImage != nil && image {
			n.AddTask(30, "boot_image").Exec = n.BootImage
			n.AddTask(40, "boot_image_wait").Exec = TimeoutTask(5 * time.Minute)
		}
		n.AddTask(50, "aci_check").Exec = n.RunACICheck
		n.AddTask(51, "arista_check").Exec = n.RunAristaCheck
	}
	if bm {
		n.AddBaremetalCreateTasks()

	}
	if events {
		n.AddTask(70, "delete_event_sub").Exec = n.DeleteEventSubscription
	}
	n.AddTask(100, "update_netbox").Exec = n.Update
	return
}

func TimeoutTask(d time.Duration) func() (err error) {
	return func() (err error) {
		time.Sleep(d)
		return
	}
}
