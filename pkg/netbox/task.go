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

package netbox

import (
	"encoding/json"
	"fmt"

	"github.com/netbox-community/go-netbox/v3/netbox/models"
)

type ConfigContext struct {
	Baremetal TemperContext `json:"baremetal"`
}

type TemperContext struct {
	Temper TaskContext `json:"temper"`
}
type TaskContext struct {
	Tasks []*Task `json:"tasks"`
}

type Task struct {
	Service string  `json:"service"`
	Task    string  `json:"task"`
	Exec    []*Exec `json:"-"`
	Error   string  `json:"error,omitempty"`
	Status  string  `json:"status"`
}

type Exec struct {
	Fn   func() error
	Name string
}

func (n *Netbox) GetTemperConfigContext() (temperCtx ConfigContext, err error) {
	d, err := n.GetData()
	if err != nil {
		return
	}
	ctx, ok := d.Device.ConfigContext.(map[string]interface{})
	if !ok {
		return temperCtx, fmt.Errorf("cannot cast interface to netbox ConfigContext")
	}
	bm, ok := ctx["baremetal"].(map[string]interface{})
	if !ok {
		return temperCtx, fmt.Errorf("cannot cast interface to netbox ConfigContext")
	}
	b, err := json.Marshal(bm["temper"])
	taskCtx := TaskContext{}
	json.Unmarshal(b, &taskCtx)
	for _, t := range taskCtx.Tasks {
		//remove old error logs
		t.Error = ""
	}

	temperCtx.Baremetal.Temper = taskCtx
	return
}

func (n *Netbox) WriteLocalContextData(t []*Task) (err error) {
	d, err := n.GetData()
	if err != nil {
		return
	}
	ctx, ok := d.Device.ConfigContext.(map[string]interface{})
	if !ok {
		return fmt.Errorf("cannot cast interface to netbox ConfigContext")
	}
	ctx["baremetal"] = TemperContext{
		Temper: TaskContext{
			Tasks: t,
		},
	}
	if err != nil {
		return err
	}

	_, err = n.updateNodeInfo(
		models.WritableDeviceWithConfigContext{
			LocalContextData: ctx,
		},
	)
	return
}
