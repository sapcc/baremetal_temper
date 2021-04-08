package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	OpenstackAuth   OpenstackAuth `yaml:"os_auth"`
	Inspector       Inspector     `yaml:"inspector"`
	Redfish         Redfish       `yaml:"redfish"`
	NetboxAuth      NetboxAuth    `yaml:"netbox_auth"`
	AristaAuth      AristaAuth    `yaml:"arista_auth"`
	AciAuth         AciAuth       `yaml:"aci_auth"`
	AwxAuth         AwxAuth       `yaml:"awx_auth"`
	NetboxNodesPath string        `yaml:"netbox_nodes_path"`
	RulesPath       string        `yaml:"rules_path"`
	Region          string        `yaml:"region"`
	NetboxQuery     *string       `yaml:"netbox_query"`
	Domain          string        `yaml:"domain"`
	NameSpace       string        `yaml:"namespace"`
	Deployment      Deployment    `yaml:"deployment"`
}

type Inspector struct {
	Host string `yaml:"host"`
}

type Redfish struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type OpenstackAuth struct {
	User              string `yaml:"user"`
	Password          string `yaml:"password"`
	AuthURL           string `yaml:"auth_url"`
	DomainName        string `yaml:"user_domain_name"`
	ProjectName       string `yaml:"project_name"`
	ProjectDomainName string `yaml:"domain_name"`
}

type NetboxAuth struct {
	Host  string `yaml:"host"`
	Token string `yaml:"token"`
}

type AristaAuth struct {
	Transport string `yaml:"transport"`
	Password  string `yaml:"password"`
	User      string `yaml:"user"`
	Port      int    `yaml:"port"`
}

type AciAuth struct {
	Password string `yaml:"password"`
	User     string `yaml:"user"`
}

type AwxAuth struct {
	Password string `yaml:"password"`
	User     string `yaml:"user"`
	Host     string `yaml:"host"`
}

type Deployment struct {
	Image         string `yaml:"image"`
	ConductorZone string `yaml:"conductor_zone"`
	Flavor        string `yaml:"flavor"`
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
