package cmd

import (
	"github.com/sapcc/baremetal_temper/pkg/model"
	"github.com/sapcc/baremetal_temper/pkg/temper"
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
		t := temper.New(cfg)
		c, err := t.GetClients(node)
		if err != nil {
			log.Fatal(err)
		}
		n, err := t.LoadNodeInfos(node)
		if err != nil {
			log.Fatal(err)
		}
		if err = t.TemperNode(&n, []func(n *model.Node) error{c.Netbox.Update}); err != nil {
			log.Fatal(err)
		}
		log.Info("node synced")
	},
}

func init() {
	netboxCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(netboxCmd)
}
