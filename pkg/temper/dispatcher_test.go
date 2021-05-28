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

	"github.com/sapcc/baremetal_temper/pkg/config"
	"github.com/sapcc/baremetal_temper/pkg/node"
	"github.com/stretchr/testify/assert"
)

func TestDispatcher(t *testing.T) {
	d := NewDispatcher(10)
	d.Start()

	assert.Equal(t, 10, len(d.Workers), "expects worker count to be 10")

	n, _ := node.New("test", config.Config{})
	d.Dispatch(n)

	assert.Equal(t, n.Status, "progress", "expects node job status to be 'progress'")

	time.Sleep(1 * time.Millisecond)

	assert.Equal(t, n.Status, "staged", "expects node job status to be 'staged'")

	d.Stop()
	time.Sleep(1 * time.Millisecond)
	for _, w := range d.Workers {
		checkChannelClosed(t, w.JobChan)
	}
}

func checkChannelClosed(t *testing.T, JobChan JobChannel) {
	select {
	case <-JobChan:
	default:
		t.Errorf("expects worker channel to be closed")
	}
}
