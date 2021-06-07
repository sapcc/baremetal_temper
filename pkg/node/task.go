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
	"fmt"
	"time"
)

func (n *Node) initTasks() {
	n.taskList["temper_dns"] = []*Task{
		{Exec: n.createDNSRecords, Name: "create_dns_records"},
	}
	n.taskList["temper_cable-check"] = []*Task{
		{Exec: n.bootImage, Name: "boot_image"},
		{Exec: TimeoutTask(5 * time.Minute), Name: "boot_image_wait"},
		{Exec: n.runACICheck, Name: "aci_cable_check"},
		{Exec: n.runAristaCheck, Name: "arista_cable_check"},
	}
	n.taskList["temper_import-ironic"] = []*Task{
		{Exec: n.create, Name: "create_ironic_node"},
		{Exec: n.checkCreated, Name: "check_ironic_node_created"},
		{Exec: n.applyRules, Name: "apply_ironic_rules"},
		{Exec: n.validate, Name: "validate_ironic_node"},
		{Exec: n.powerOn, Name: "power_on_ironic_node"},
		{Exec: n.provide, Name: "provide_ironic_node"},
	}
	n.taskList["temper_ironic-test-deployment"] = []*Task{
		{Exec: n.waitForNovaPropagation, Name: "wait_nova_propagation"},
		{Exec: n.deployTestInstance, Name: "deploy_test_instance"},
	}
	n.taskList["temper_hardware-check"] = []*Task{
		{Exec: n.runHardwareChecks, Name: "hardware_checks"},
	}
	n.taskList["temper_sync-netbox"] = []*Task{
		{Exec: n.update, Name: "sync_netbox"},
	}
	n.taskList["temper_fw-upgrade"] = []*Task{}
	n.taskList["temper_bmc-settings"] = []*Task{}
}

func (n *Node) AddTask(name string) error {
	if n.Tasks == nil {
		n.Tasks = make([]*Task, 0)
	}
	t, ok := n.taskList[name]
	if !ok {
		return fmt.Errorf("unknown task")
	}
	n.Tasks = append(n.Tasks, t...)
	return nil
}

func TimeoutTask(d time.Duration) func() (err error) {
	return func() (err error) {
		time.Sleep(d)
		return
	}
}
