package cmd

import (
	"context"
	"sync"

	"github.com/sapcc/baremetal_temper/pkg/node"
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
				wg.Add(1)
				go bootImageExec(n, t, &wg)
			}
		}
		wg.Wait()
		log.Info("boot image completed")
	},
}

func bootImageExec(n string, t *temper.Temper, wg *sync.WaitGroup) {
	defer wg.Done()
	node, err := node.New(n, cfg)
	if err != nil {
		log.Errorf("error node %s: %s", n, err.Error())
	}
	node.BootImage()
}

func init() {
	redfishCmd.AddCommand(bootImage)
	rootCmd.AddCommand(redfishCmd)
}
