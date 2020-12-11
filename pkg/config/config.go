package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	IronicAuth      IronicAuth `yaml:"ironic_auth"`
	Inspector       Inspector  `yaml:"inspector"`
	Redfish         Redfish    `yaml:"redfish"`
	NetboxNodesPath string     `yaml:"netbox_nodes_path"`
	OsRegion        string     `yaml:"os_region"`
	Domain          string     `yaml:"domain"`
	NameSpace       string     `yaml:"namespace"`
}

type Inspector struct {
	Host     string `yaml:"host"`
	Callback string `yaml:"callback"`
}

type Redfish struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type IronicAuth struct {
	User              string `yaml:"user"`
	Password          string `yaml:"password"`
	AuthURL           string `yaml:"auth_url"`
	DomainName        string `yaml:"user_domain_name"`
	ProjectName       string `yaml:"project_name"`
	ProjectDomainName string `yaml:"domain_name"`
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
