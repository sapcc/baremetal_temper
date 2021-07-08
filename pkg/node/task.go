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

	"github.com/sapcc/baremetal_temper/pkg/task"
)

func (n *Node) initTaskExecs() {
	n.tasksExecs["dns"] = map[string][]*task.Exec{
		"create": {
			{Fn: n.createDNSRecords, Name: "dns.create"},
		},
	}
	if *n.cfg.Redfish.BootImage == "" {
		n.log.Warning("did not find boot image for cable check. run check without it")
		n.tasksExecs["diagnostics"] = map[string][]*task.Exec{
			"cablecheck": {
				{Fn: n.runACICheck, Name: "diagnostics.cablecheck.aci"},
				//{Exec: n.runAristaCheck, Name: "arista_cable_check"},
			},
			"hardwarecheck": {
				{Fn: n.runHardwareChecks, Name: "diagnostics.hardwarecheck"},
			},
			"bootimage": {
				{Fn: n.bootImage, Name: "diagnostics.bootimage"},
			},
		}
	} else {
		n.tasksExecs["diagnostics"] = map[string][]*task.Exec{
			"cablecheck": {
				{Fn: n.bootImage, Name: "diagnostics.cablecheck.bootimage"},
				{Fn: TimeoutTask(10 * time.Minute), Name: "diagnostics.cablecheck.bootimage.wait"},
				{Fn: n.runACICheck, Name: "diagnostics.cablecheck.aci"},
				//{Exec: n.runAristaCheck, Name: "arista_cable_check"},
				{Fn: n.ejectMedia, Name: "diagnostics.cablecheck.bootimage.eject"},
				{Fn: func() error { return n.power(false) }, Name: "diagnostics.cablecheck.reboot"},
			},
			"hardwarecheck": {
				{Fn: n.runHardwareChecks, Name: "diagnostics.cablecheck.hardwarecheck"},
			},
		}
	}
	n.tasksExecs["ironic"] = map[string][]*task.Exec{
		"create": {
			{Fn: n.create, Name: "ironic.create"},
			{Fn: n.checkCreated, Name: "ironic.create.check"},
			{Fn: TimeoutTask(10 * time.Second), Name: "ironic.create.wait"},
			{Fn: n.applyRules, Name: "ironic.create.applyRules"},
			{Fn: n.powerOn, Name: "ironic.create.powerOn"},
			{Fn: n.provide, Name: "ironic.create.provide"},
		},
		"validate": {
			{Fn: n.validate, Name: "ironic.validate"},
		},
		"test": {
			{Fn: n.waitForNovaPropagation, Name: "ironic.test.waitForNovaPropagation"},
			{Fn: n.deployTestInstance, Name: "ironic.test.deploy"},
		},
		"prepare": {
			{Fn: n.prepare, Name: "ironic.prepare"},
		},
	}
	n.tasksExecs["netbox"] = map[string][]*task.Exec{
		"sync": {
			{Fn: n.update, Name: "netbox.sync"},
		},
	}
	n.tasksExecs["firmware"] = map[string][]*task.Exec{
		"profile": {},
		"update":  {},
	}
	n.tasksExecs["bios"] = map[string][]*task.Exec{
		"profile": {},
		"update":  {},
	}
}

func (n *Node) AddTask(service, taskName string) error {
	if n.Tasks == nil {
		n.Tasks = make([]*task.Task, 0)
	}
	if taskName == "all" {
		for t, e := range n.tasksExecs[service] {
			t := &task.Task{
				Service: service,
				Task:    t,
				Exec:    e,
			}
			n.Tasks = append(n.Tasks, t)
		}
		return nil
	}
	execs, ok := n.tasksExecs[service][taskName]
	if !ok {
		return fmt.Errorf("unknown task")
	}

	t := &task.Task{
		Service: service,
		Task:    taskName,
		Exec:    execs,
	}
	n.Tasks = append(n.Tasks, t)
	return nil
}

func (n *Node) MergeTaskWithContext(cfgCtx task.ConfigContext) error {
	if n.Tasks == nil {
		n.Tasks = make([]*task.Task, 0)
	}
	for _, t := range cfgCtx.Baremetal.Temper.Tasks {
		execs, ok := n.tasksExecs[t.Service][t.Task]
		if !ok {
			return fmt.Errorf("unknown task")
		}
		t.Exec = execs
		n.Tasks = append(n.Tasks, t)
	}
	return nil
}

func TimeoutTask(d time.Duration) func() (err error) {
	return func() (err error) {
		time.Sleep(d)
		return
	}
}
