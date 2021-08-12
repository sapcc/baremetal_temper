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

package redfish

type Lenovo struct {
	Default
}

func (p *Lenovo) insertMedia(image string) (err error) {
	//n.log.Debug("insert virtual media")
	vm, err := p.getDVDMediaType()
	if err != nil {
		return
	}
	if vm.SupportsMediaInsert {
		if vm.Image != "" {
			err = vm.EjectMedia()
		}
		type temp struct {
			Image          string
			Inserted       bool `json:"Inserted"`
			WriteProtected bool `json:"WriteProtected"`
		}
		t := temp{
			Image:          image,
			Inserted:       true,
			WriteProtected: false,
		}
		return vm.InsertMedia(image, "PATCH", t)
	}
	return
}

func (p *Lenovo) ejectMedia() (err error) {
	//n.log.Debug("eject media image")
	vm, err := p.getDVDMediaType()
	if err != nil {
		return
	}
	if vm.SupportsMediaInsert {
		if vm.Image != "" {
			return vm.InsertMedia("", "PATCH", nil)
		}
	}
	return
}
