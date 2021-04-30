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
	"github.com/sapcc/baremetal_temper/pkg/temper"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	baremetal     bool
	diag          bool
	redfishEvents bool
)

var complete = &cobra.Command{
	Use:   "complete",
	Short: "tempers a node",
	Run: func(cmd *cobra.Command, args []string) {
		t := temper.New(cfg)
		c, err := t.GetClients(node)
		if err != nil {
			log.Fatal(err)
		}
		n, err := t.LoadNodeInfos(node)
		if err != nil {
			log.Fatal(err)
		}
		tasks, err := t.GetAllTemperTasks(n.Name, diag, baremetal, redfishEvents)
		if err != nil {
			log.Fatal(err)
		}
		if err = t.TemperNode(&n, tasks); err != nil {
			log.Fatal(err)
		}

		if err = c.Netbox.Update(&n); err != nil {
			log.Fatal(err)
		}
		log.Info("node synced")
	},
}

func init() {
	complete.PersistentFlags().BoolVar(&baremetal, "baremetal", true, "run baremetal tasks")
	complete.PersistentFlags().BoolVar(&diag, "diagnostics", true, "run diagnostics tasks")
	complete.PersistentFlags().BoolVar(&redfishEvents, "redfishEvents", true, "use redfish events")

	rootCmd.AddCommand(complete)
}
