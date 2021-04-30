package cmd

import (
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
	Args:  cobra.ExactArgs(0),
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
		if err = t.TemperNode(&n, []func(n *model.Node) error{c.Redfish.BootImage}); err != nil {
			log.Fatal(err)
		}
		log.Info("booting image")
	},
}

func init() {
	redfishCmd.AddCommand(bootImage)
	rootCmd.AddCommand(redfishCmd)
}
