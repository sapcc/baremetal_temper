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

package task

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
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

func UnmarshalConfigContext(data interface{}) (cfgCtx ConfigContext, err error) {
	cfgCtx = ConfigContext{}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err = enc.Encode(data); err != nil {
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &cfgCtx); err != nil {
		return
	}
	return
}
