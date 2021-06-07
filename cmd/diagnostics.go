package cmd

import (
	"sync"

	"github.com/sapcc/baremetal_temper/pkg/node"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var diagCmd = &cobra.Command{
	Use:   "diagnostics",
	Short: "interact with diagnostics services",
}

var hardwareCheck = &cobra.Command{
	Use:   "hardwarecheck",
	Short: "runs a vendor specific hardware check",
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		if err := loadNodes(); err != nil {
			log.Errorf("error loading nodes: %s", err.Error())
		}
		for _, n := range nodes {
			node, err := node.New(n, cfg)
			if err != nil {
				log.Errorf("error node %s: %s", n, err.Error())
				return
			}
			node.AddTask("temper_hardware-check")
			wg.Add(1)
			go node.Temper(netboxStatus, &wg)
		}
		wg.Wait()
		log.Info("check completed")
	},
}

var cableCheck = &cobra.Command{
	Use:   "cablecheck",
	Short: "runs a cable check (lldp)",
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		if err := loadNodes(); err != nil {
			log.Errorf("error loading nodes: %s", err.Error())
		}
		for _, n := range nodes {
			nd, err := node.New(n, cfg)
			if err != nil {
				log.Errorf("error node %s: %s", n, err.Error())
				continue
			}

			nd.AddTask("temper_cable-check")
			wg.Add(1)
			go nd.Temper(netboxStatus, &wg)
		}
		wg.Wait()
		log.Info("cable check completed")
	},
}

func init() {
	diagCmd.PersistentFlags().BoolVar(&bootImg, "bootImage", false, "boots an image before running cablecheck")

	diagCmd.AddCommand(cableCheck)
	diagCmd.AddCommand(hardwareCheck)
	rootCmd.AddCommand(diagCmd)
}
