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
	"strings"
	"sync"

	"github.com/sapcc/baremetal_temper/pkg/node"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var tasks []string
var workers int

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run tasks via -t 'service.task'. If all tasks from a service should be run use: 'service.all'",
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
			wg.Add(1)
			for _, t := range tasks {
				s := strings.Split(t, ".")
				if len(s) != 2 {
					log.Error("wrong task format. It should be [service].[task]")
					continue
				}
				if err = n.AddTask(s[0], s[1]); err != nil {
					log.Error(err)
					continue
				}
			}
			go n.Temper(netboxStatus, &wg, limiter)
		}
		log.Info("number of go-routines: ", runtime.NumGoroutine())
		wg.Wait()
		log.Info("command complete")
	},
}

func init() {
	runCmd.PersistentFlags().StringArrayVarP(&tasks, "tasks", "t", []string{}, "array of tasks to run e.g. 'ironic.create'")
	runCmd.PersistentFlags().IntVarP(&workers, "workers", "w", 0, "number of max worker to execute tasks concurrently. Default is 0: infinite.")
	rootCmd.AddCommand(runCmd)
}
