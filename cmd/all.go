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

	"github.com/sapcc/baremetal_temper/pkg/temper"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	baremetal     bool
	diag          bool
	redfishEvents bool
	bootImg       bool
)

var complete = &cobra.Command{
	Use:   "complete",
	Short: "tempers a node",
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		t := temper.New(cfg)
		if len(nodes) > 0 {
			for _, n := range nodes {
				wg.Add(1)
				go execComplete(t, n, &wg)
			}
		}
		wg.Wait()
		log.Info("command complete")
	},
}

func execComplete(t *temper.Temper, node string, wg *sync.WaitGroup) {
	defer wg.Done()
	c, err := t.GetClients(node)
	if err != nil {
		log.Errorf("error node %s: %s", node, err.Error())
		return
	}
	n, err := t.LoadNodeInfos(node)
	if err != nil {
		log.Errorf("error node %s: %s", node, err.Error())
		return
	}
	tasks, err := t.GetAllTemperTasks(n.Name, diag, baremetal, redfishEvents, bootImg)
	if err != nil {
		log.Errorf("error node %s: %s", node, err.Error())
		return
	}
	if err = t.TemperNode(&n, tasks); err != nil {
		log.Errorf("error node %s: %s", node, err.Error())
		return
	}

	if err = c.Netbox.Update(&n); err != nil {
		log.Errorf("error node %s: %s", node, err.Error())
		return
	}
	log.Infof("node %s done", node)
}

func init() {
	complete.PersistentFlags().BoolVar(&baremetal, "baremetal", false, "run baremetal tasks")
	complete.PersistentFlags().BoolVar(&diag, "diagnostics", true, "run diagnostics tasks")
	complete.PersistentFlags().BoolVar(&redfishEvents, "redfishEvents", false, "use redfish events")
	complete.PersistentFlags().BoolVar(&bootImg, "bootImage", false, "boots an image before running cablecheck")

	rootCmd.AddCommand(complete)
}
