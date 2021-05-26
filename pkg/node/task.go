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

package node

func (n *Node) AddBaremetalCreateTasks() {
	n.Tasks[100] = &Task{Name: "create_ironic_node", Exec: n.Create}
	n.Tasks[90] = &Task{Name: "check_ironic_node_created", Exec: n.CheckCreated}
	n.Tasks[80] = &Task{Name: "apply_ironic_rules", Exec: n.ApplyRules}
	n.Tasks[70] = &Task{Name: "validate_ironic_node", Exec: n.Validate}
	n.Tasks[60] = &Task{Name: "power_on_ironic_node", Exec: n.PowerOn}
	n.Tasks[50] = &Task{Name: "provide_ironic_node", Exec: n.Provide}
}

func (n *Node) AddDeploymentTestTasks() {
	n.Tasks[100] = &Task{Name: "get_ironic_uuid", Exec: n.getUUID}
	n.Tasks[90] = &Task{Name: "wait_nova_propagation", Exec: n.WaitForNovaPropagation}
	n.Tasks[80] = &Task{Name: "deploy_test_instance", Exec: n.DeployTestInstance}
}
