package config

import (
	"fmt"
	"io/ioutil"

	"github.com/gophercloud/gophercloud/openstack/baremetal/v1/nodes"
	"github.com/gophercloud/gophercloud/openstack/baremetal/v1/ports"
)

type Rule struct {
	Properties Property `json:"properties"`
}

type Property struct {
	Node      []NodeConfig `json:"node"`
	Port      []PortConfig `json:"port"`
	PortGroup []PortConfig `json:"port_group"`
}

type NodeConfig struct {
	Op    nodes.UpdateOp `json:"op"`
	Path  string         `json:"path"`
	Value interface{}    `json:"value"`
}

type PortConfig struct {
	Op    ports.UpdateOp `json:"op"`
	Path  string         `json:"path"`
	Value interface{}    `json:"value"`
}

func GetRules(path string) (jsonBytes []byte, err error) {
	if path == "" {
		return
	}
	jsonBytes, err = ioutil.ReadFile(path)
	if err != nil {
		return jsonBytes, fmt.Errorf("read file file: %s", err.Error())
	}
	return
}
