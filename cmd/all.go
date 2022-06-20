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
	"runtime"
	"sync"

	"github.com/sapcc/baremetal_temper/pkg/node"
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
		err := loadNodes()
		if err != nil {
			log.Errorf("error loading nodes: %s", err.Error())
			return
		}
		limiter := make(chan bool, workers)
		if workers == 0 {
			limiter = nil
		}
		for _, na := range nodes {
			n, err := node.New(na, cfg)
			if err != nil {
				log.Errorf("error node %s: %s", na, err.Error())
				continue
			}
			n.AddTask("dns", "create")
			if diag {
				n.AddTask("diagnostics", "cablecheck")
				n.AddTask("diagnostics", "hardwarecheck")
			}
			if baremetal {
				n.AddTask("ironic", "create")
				n.AddTask("ironic", "validate")
				n.AddTask("ironic", "test")
				n.AddTask("ironic", "prepare")
			}
			if netboxStatus {
				n.AddTask("netbox", "sync")
			}
			wg.Add(1)
			go n.Temper(netboxStatus, &wg, limiter)
			log.Info("number of go-routines: ", runtime.NumGoroutine())
		}
		wg.Wait()
		log.Info("command complete")
	},
}

func init() {
	complete.PersistentFlags().BoolVar(&baremetal, "baremetal", false, "run baremetal tasks")
	complete.PersistentFlags().BoolVar(&diag, "diagnostics", true, "run diagnostics tasks")
	complete.PersistentFlags().BoolVar(&redfishEvents, "redfishEvents", false, "use redfish events")
	complete.PersistentFlags().BoolVar(&bootImg, "bootImage", false, "boots an image before running cablecheck")
	complete.PersistentFlags().IntVarP(&workers, "workers", "w", 0, "number of max worker to execute tasks concurrently. Default is 0: infinite")

	rootCmd.AddCommand(complete)
}
