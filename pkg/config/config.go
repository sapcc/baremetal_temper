package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Openstack       OpenstackAuth `yaml:"openstack"`
	Inspector       Inspector     `yaml:"inspector"`
	Redfish         Redfish       `yaml:"wathever"`
	Netbox          NetboxAuth    `yaml:"netbox"`
	Arista          AristaAuth    `yaml:"arista"`
	AciAuth         AciAuth       `yaml:"aciAuth"`
	AwxAuth         AwxAuth       `yaml:"awxAuth"`
	NetboxNodesPath string        `yaml:"netboxNodesPath"`
	RulesPath       string        `yaml:"rulesPath"`
	Region          string        `yaml:"region"`
	NetboxQuery     *string       `yaml:"netboxQuery"`
	Domain          string        `yaml:"domain"`
	NameSpace       string        `yaml:"namespace"`
	Deployment      Deployment    `yaml:"deployment"`
}

type Inspector struct {
	Host string `yaml:"host"`
}

type Redfish struct {
	User      string  `yaml:"user"`
	Password  string  `yaml:"password"`
	BootImage *string `yaml:"bootImage"`
}

type OpenstackAuth struct {
	User              string `yaml:"user"`
	Password          string `yaml:"password"`
	Url               string `yaml:"url"`
	DomainName        string `yaml:"domainName"`
	ProjectName       string `yaml:"projectName"`
	ProjectDomainName string `yaml:"projectDomainName"`
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
	Image         string        `yaml:"image"`
	ConductorZone string        `yaml:"conductorZone"`
	Flavor        string        `yaml:"flavor"`
	Network       string        `yaml:"network"`
	OpenstackAuth OpenstackAuth `yaml:"osAuth"`
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
