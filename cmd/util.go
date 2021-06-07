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
	"fmt"

	"github.com/sapcc/baremetal_temper/pkg/clients"

	log "github.com/sirupsen/logrus"
)

func loadNodes() (err error) {
	ctxLogger := log.WithFields(log.Fields{
		"cmd": "temper",
	})
	n, err := clients.NewNetbox(cfg, ctxLogger)
	if err != nil {
		return
	}
	if nodeQuery != "" {
		nodes, err = n.LoadNodes(&nodeQuery, &nodeStatus, &cfg.Region)
		if err != nil {
			log.Errorf("error loading nodes: %s", err.Error())
			return
		}
	}
	if len(nodes) == 0 {
		return fmt.Errorf("no nodes provided")
	}
	return
}
