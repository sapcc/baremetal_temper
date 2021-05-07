package cmd

import (
	"sync"

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
		var wg sync.WaitGroup
		t := temper.New(cfg)
		if len(nodes) > 0 {
			for _, n := range nodes {
				c, err := t.GetClients(n)
				if err != nil {
					log.Errorf("error node %s: %s", n, err.Error())
					continue
				}
				wg.Add(1)
				go temperNode(t, n, []func(n *model.Node) error{c.Netbox.Update}, &wg)
			}
		}
		wg.Wait()
		log.Info("sync completed")
	},
}

func init() {
	netboxCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(netboxCmd)
}
