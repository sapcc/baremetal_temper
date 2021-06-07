package cmd

import (
	"sync"

	"github.com/sapcc/baremetal_temper/pkg/node"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var dnsCmd = &cobra.Command{
	Use:   "dns",
	Short: "interact with dns service",
}

var createDNS = &cobra.Command{
	Use:   "create",
	Short: "creates a nodes dns records based on netbox info",
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		if err := loadNodes(); err != nil {
			log.Errorf("error loading nodes: %s", err.Error())
		}
		for _, n := range nodes {
			wg.Add(1)
			node, err := node.New(n, cfg)
			if err != nil {
				log.Errorf("error node %s: %s", n, err.Error())
			}
			node.AddTask("temper_dns")
			go node.Temper(netboxStatus, &wg)
		}
		wg.Wait()
		log.Info("dns create completed")
	},
}

func init() {
	dnsCmd.AddCommand(createDNS)
	rootCmd.AddCommand(dnsCmd)
}
