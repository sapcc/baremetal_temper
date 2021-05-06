package cmd

import (
	"github.com/sapcc/baremetal_temper/pkg/model"
	"github.com/sapcc/baremetal_temper/pkg/temper"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var dnsCmd = &cobra.Command{
	Use:   "dns",
	Short: "interact with dns service",
}

var createDns = &cobra.Command{
	Use:   "create",
	Short: "creates a nodes dns records based on netbox info",
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
		if err = t.TemperNode(&n, []func(n *model.Node) error{c.Openstack.CreateDNSRecords}); err != nil {
			log.Fatal(err)
		}
		log.Info("records created")
	},
}

func init() {
	dnsCmd.AddCommand(createDns)
	rootCmd.AddCommand(dnsCmd)
}
