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

package temper

import (
	"testing"
	"time"

	"github.com/netbox-community/go-netbox/netbox/models"
	"github.com/sapcc/baremetal_temper/pkg/config"
	"github.com/sapcc/baremetal_temper/pkg/node"
	"github.com/stretchr/testify/assert"
)

func TestTemper(t *testing.T) {
	tp := New(1)
	assert.Equal(t, 0, len(tp.GetNodes()), "expects node list to be 0")
	n, _ := node.New("test1", config.Config{})
	n.DeviceConfig = &models.DeviceWithConfigContext{}
	tp.AddNode(n)
	assert.Equal(t, 1, len(tp.GetNodes()), "expects node list to be 1")
	n, _ = node.New("test1", config.Config{})
	n.DeviceConfig = &models.DeviceWithConfigContext{}
	n.Status = "planned"
	tp.AddNode(n)
	time.Sleep(1 * time.Millisecond)
	assert.Equal(t, "failed", n.Status, "expects node status to be failed")
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, 0, len(tp.GetNodes()), "expects node list to be 0")
}
