package cmd

import (
	"context"
	"sync"

	"github.com/sapcc/baremetal_temper/pkg/node"
	"github.com/sapcc/baremetal_temper/pkg/temper"
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
		t := temper.New(cfg, context.Background(), netboxStatus)
		if len(nodes) > 0 {
			for _, n := range nodes {
				wg.Add(1)
				go createDNSExec(n, t, &wg)
			}
		}
		wg.Wait()
		log.Info("dns create completed")
	},
}

func createDNSExec(n string, t *temper.Temper, wg *sync.WaitGroup) {
	defer wg.Done()
	node, err := node.New(n, cfg)
	if err != nil {
		log.Errorf("error node %s: %s", n, err.Error())
	}
	node.CreateDNSRecords()
}

func init() {
	dnsCmd.AddCommand(createDNS)
	rootCmd.AddCommand(dnsCmd)
}
