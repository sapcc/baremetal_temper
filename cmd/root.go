package cmd

import (
	"github.com/mitchellh/go-homedir"
	"github.com/sapcc/baremetal_temper/pkg/config"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	logLevel int
	node     string
	cfgFile  string
	cfg      config.Config
)

var rootCmd = &cobra.Command{
	Use:   "temper",
	Short: "temper manages your node instances",
	Long: `temper manages your node instances
including the backend database.
Supports warmups and backups.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().IntVarP(&logLevel, "log-level", "l", 5, "temper log level")
	rootCmd.PersistentFlags().StringVarP(&node, "node", "n", "", "node name")
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/config.yaml)")

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

// initConfig reads in config file and ENV variables if set.
func initConfig() {
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
	viper.SetDefault("openstack.userDomainName", "")
	viper.BindEnv("openstack.userDomainName", "openstack_userDomainName")
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
	viper.SetDefault("arista.port", "")
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
	viper.SetDefault("deployment.openstack.userDomainName", "")
	viper.BindEnv("deployment.openstack.userDomainName", "deployment_openstack_userDomainName")
	viper.SetDefault("deployment.openstack.projectName", "")
	viper.BindEnv("deployment.openstack.projectName", "deployment_openstack_projectName")
	viper.SetDefault("deployment.openstack.domainName", "")
	viper.BindEnv("deployment.openstack.domainName", "deployment_openstack_domainName")

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Fatal(err.Error())
		}
		viper.AddConfigPath(home)
		viper.SetConfigName("config.yaml")
	}
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Info("Using config file:", viper.ConfigFileUsed())
	}
	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal(err.Error())
	}
}
