package cmd

import (
	"context"
	"sync"

	"github.com/sapcc/baremetal_temper/pkg/model"
	"github.com/sapcc/baremetal_temper/pkg/temper"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var redfishCmd = &cobra.Command{
	Use:   "redfish",
	Short: "interact with a node's redfish api",
}

var bootImage = &cobra.Command{
	Use:   "bootimage",
	Short: "mounts and boots an image via redfish",
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		t := temper.New(cfg, context.Background(), netboxStatus)
		if len(nodes) > 0 {
			for _, n := range nodes {
				c, err := t.GetClients(n)
				if err != nil {
					log.Errorf("error node %s: %s", n, err.Error())
					continue
				}
				wg.Add(1)
				go temperNode(t, n, []func(n *model.Node) error{c.Redfish.BootImage}, &wg)
			}
		}
		wg.Wait()
		log.Info("boot image completed")
	},
}

func init() {
	redfishCmd.AddCommand(bootImage)
	rootCmd.AddCommand(redfishCmd)
}
