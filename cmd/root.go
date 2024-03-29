package cmd

import (
	"fmt"
	"os"

	"github.com/sapcc/baremetal_temper/pkg/config"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	logLevel     int
	nodes        []string
	nodeQuery    string
	netboxStatus bool
	nodeStatus   string
	cfgFile      string
	cfg          config.Config
)

var rootCmd = &cobra.Command{
	Use:   "temper",
	Short: "temper manages your node instances",
	Long:  `temper manages your node instances`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cobra.OnInitialize(InitConfig)
	rootCmd.PersistentFlags().IntVarP(&logLevel, "log-level", "l", 5, "temper log level")
	rootCmd.PersistentFlags().StringArrayVarP(&nodes, "nodes", "n", []string{}, "array of nodes")
	rootCmd.PersistentFlags().StringVarP(&nodeQuery, "node-query", "q", "", "query to load nodes via netbox")
	rootCmd.PersistentFlags().StringVar(&nodeStatus, "node-status", "planned", "node status to load via netbox")
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/config.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&netboxStatus, "netbox-status", "s", false, "set to true if the node's netbox status should be updated")

	switch logLevel {
	case 1:
		log.SetLevel(log.FatalLevel)
	case 2:
		log.SetLevel(log.ErrorLevel)
	case 3:
		log.SetLevel(log.WarnLevel)
	case 4:
		log.SetLevel(log.InfoLevel)
	case 5:
		log.SetLevel(log.DebugLevel)
	}

}

// Execute executes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}

// InitConfig reads in config file and ENV variables if set.
func InitConfig() {
	viper.SetDefault("region", "")
	viper.BindEnv("region", "region")

	viper.SetDefault("inspector.host", "")
	viper.BindEnv("inspector.host", "inspector_host")

	viper.SetDefault("redfish.user", "")
	viper.BindEnv("redfish.user", "redfish_user")
	viper.SetDefault("redfish.password", "")
	viper.BindEnv("redfish.password", "redfish_password")
	viper.SetDefault("redfish.bootImage", "")
	viper.BindEnv("redfish.bootImage", "redfish_bootImage")

	viper.SetDefault("netbox.token", "")
	viper.BindEnv("netbox.token", "netbox_token")
	viper.SetDefault("netbox.host", "")
	viper.BindEnv("netbox.host", "netbox_host")

	viper.SetDefault("openstack.user", "")
	viper.BindEnv("openstack.user", "openstack_user")
	viper.SetDefault("openstack.password", "")
	viper.BindEnv("openstack.password", "openstack_password")
	viper.SetDefault("openstack.projectDomainName", "")
	viper.BindEnv("openstack.projectDomainName", "openstack_projectDomainName")
	viper.SetDefault("openstack.projectName", "")
	viper.BindEnv("openstack.projectName", "openstack_projectName")
	viper.SetDefault("openstack.domainName", "")
	viper.BindEnv("openstack.domainName", "openstack_domainName")

	viper.SetDefault("arista.user", "")
	viper.BindEnv("arista.user", "arista_user")
	viper.SetDefault("arista.password", "")
	viper.BindEnv("arista.password", "arista_password")
	viper.SetDefault("arista.transport", "")
	viper.BindEnv("arista.transport", "arista_transport")
	viper.SetDefault("arista.port", 0)
	viper.BindEnv("arista.port", "arista_port")

	viper.SetDefault("aci.user", "")
	viper.BindEnv("aci.user", "aci_user")
	viper.SetDefault("aci.password", "")
	viper.BindEnv("aci.password", "aci_password")

	viper.SetDefault("deployment.image", "")
	viper.BindEnv("deployment.image", "deployment_image")
	viper.SetDefault("deployment.conductorZone", "")
	viper.BindEnv("deployment.conductorZone", "deployment_conductorZone")
	viper.SetDefault("deployment.flavor", "")
	viper.BindEnv("deployment.flavor", "deployment_flavor")
	viper.SetDefault("deployment.network", "")
	viper.BindEnv("deployment.network", "deployment_network")
	viper.SetDefault("deployment.openstack.user", "")
	viper.BindEnv("deployment.openstack.user", "deployment_openstack_user")
	viper.SetDefault("deployment.openstack.password", "")
	viper.BindEnv("deployment.openstack.password", "deployment_openstack_password")
	viper.SetDefault("deployment.openstack.projectDomainName", "")
	viper.BindEnv("deployment.openstack.projectDomainName", "deployment_openstack_projectDomainName")
	viper.SetDefault("deployment.openstack.projectName", "")
	viper.BindEnv("deployment.openstack.projectName", "deployment_openstack_projectName")
	viper.SetDefault("deployment.openstack.domainName", "")
	viper.BindEnv("deployment.openstack.domainName", "deployment_openstack_domainName")

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	}
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Info("Using config file:", viper.ConfigFileUsed())
	} else {
		log.Error(fmt.Sprintf("cannot read config file: %s", err.Error()))
		os.Exit(1)
	}
	viper.AutomaticEnv() // read in environment variables that match
	UnmarshalConfig(&cfg)
}

func UnmarshalConfig(cfg *config.Config) {
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal(err.Error())
	}
	if *cfg.Redfish.BootImage == "" {
		log.Warning("did not find boot image for cable check. run check without it")
	}
}
