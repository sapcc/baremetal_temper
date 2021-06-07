package cmd

import (
	"sync"

	"github.com/sapcc/baremetal_temper/pkg/node"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var (
	netboxToken     string
	netboxHost      string
	redfishUser     string
	redfishPassword string
)

var netboxCmd = &cobra.Command{
	Use:   "netbox",
	Short: "interact with netbox service",
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "syncs a node's netbox information based on it's redfish data",
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		if err := loadNodes(); err != nil {
			log.Errorf("error loading nodes: %s", err.Error())
		}
		for _, n := range nodes {
			wg.Add(1)
			nd, err := node.New(n, cfg)
			if err != nil {
				log.Errorf("error node %s: %s", n, err.Error())
			}
			nd.AddTask("temper_sync-netbox")
			go nd.Temper(netboxStatus, &wg)
		}
		wg.Wait()
		log.Info("sync completed")
	},
}

func init() {
	netboxCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(netboxCmd)
}
