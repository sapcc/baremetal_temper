package cmd

import (
	"context"
	"sync"

	"github.com/sapcc/baremetal_temper/pkg/node"
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
		t := temper.New(cfg, context.Background(), netboxStatus)
		if len(nodes) > 0 {
			for _, n := range nodes {
				wg.Add(1)
				go syncExec(n, t, &wg)
			}
		}
		wg.Wait()
		log.Info("sync completed")
	},
}

func syncExec(n string, t *temper.Temper, wg *sync.WaitGroup) {
	defer wg.Done()
	nd, err := node.New(n, cfg)
	if err != nil {
		log.Errorf("error node %s: %s", n, err.Error())
	}
	nd.Update()
}

func init() {
	netboxCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(netboxCmd)
}
