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
	"github.com/sapcc/baremetal_temper/pkg/node"
)

type dispatcher struct {
	Workers  []*Worker
	WorkChan JobChannel
	Queue    JobQueue
	quit     chan struct{}
}

func NewDispatcher(num int) *dispatcher {
	return &dispatcher{
		Workers:  make([]*Worker, num),
		WorkChan: make(JobChannel),
		Queue:    make(JobQueue),
		quit:     make(chan struct{}),
	}
}

func (d *dispatcher) Start() *dispatcher {
	for i := 0; i < len(d.Workers); i++ {
		w := Worker{make(JobChannel), d.Queue, d.quit}
		w.Start()
		d.Workers[i] = &w
	}
	go d.run()
	return d
}

func (d *dispatcher) Stop() *dispatcher {
	for i := 0; i < len(d.Workers); i++ {
		<-d.Queue
		d.quit <- struct{}{}
	}
	return d
}

func (d *dispatcher) Dispatch(job *node.Node) {
	job.Status = "progress"
	d.WorkChan <- job
}

func (d *dispatcher) run() {
	for {
		select {
		case job := <-d.WorkChan:
			jobChan := <-d.Queue
			jobChan <- job
		}
	}
}
