package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	IronicUser              string `yaml:"ironic_user"`
	IronicPassword          string `yaml:"ironic_password"`
	IronicInspectorHost     string `yaml:"ironic_inspector_host"`
	IronicInspectorCallback string `yaml:"ironic_inspector_callback"`
	NetboxNodesPath         string `yaml:"netbox_nodes_path"`
	OsRegion                string `yaml:"os_region"`
	NameSpace				string `yaml:"namespace"`
}

func GetConfig(opts Options) (cfg Config, err error) {
	if opts.ConfigFilePath == "" {
		return cfg, nil
	}
	yamlBytes, err := ioutil.ReadFile(opts.ConfigFilePath)
	if err != nil {
		return cfg, fmt.Errorf("read file file: %s", err.Error())
	}
	err = yaml.Unmarshal(yamlBytes, &cfg)
	if err != nil {
		return cfg, fmt.Errorf("parse config file: %s", err.Error())
	}

	return cfg, nil
}
