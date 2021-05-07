package cmd

import (
	"sync"

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
				go temperNode(t, n, []func(n *model.Node) error{c.Openstack.CreateDNSRecords}, &wg)
			}
		}
		wg.Wait()
		log.Info("dns create completed")
	},
}

func init() {
	dnsCmd.AddCommand(createDns)
	rootCmd.AddCommand(dnsCmd)
}
