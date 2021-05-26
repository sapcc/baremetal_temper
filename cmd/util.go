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
package cmd

import (
	"sync"

	"github.com/sapcc/baremetal_temper/pkg/node"
	"github.com/sapcc/baremetal_temper/pkg/temper"

	log "github.com/sirupsen/logrus"
)

func temperNode(t *temper.Temper, n *node.Node, wg *sync.WaitGroup) {
	defer wg.Done()
	if err := t.TemperNode(n, netboxStatus); err != nil {
		log.Errorf("error node %s: %s", n.Name, err.Error())
	}
}

func loadNodes(t *temper.Temper) (err error) {
	if nodeQuery != "" {
		nodes, err = t.LoadPlannedNodes(&nodeQuery)
		if err != nil {
			log.Errorf("error loading nodes: %s", err.Error())
			return
		}
	}
	return
}
