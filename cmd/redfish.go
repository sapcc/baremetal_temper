package cmd

import (
	"sync"

	"github.com/sapcc/baremetal_temper/pkg/node"
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
		if err := loadNodes(); err != nil {
			log.Errorf("error loading nodes: %s", err.Error())
		}
		for _, n := range nodes {
			node, err := node.New(n, cfg)
			if err != nil {
				log.Errorf("error node %s: %s", n, err.Error())
			}
			wg.Add(1)
			node.AddTask(100, "boot_image").Exec = node.BootImage
			go node.Temper(netboxStatus, &wg)
		}
		wg.Wait()
		log.Info("boot image completed")
	},
}

func init() {
	redfishCmd.AddCommand(bootImage)
	rootCmd.AddCommand(redfishCmd)
}
