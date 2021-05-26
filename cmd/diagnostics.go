package cmd

import (
	"context"
	"sync"
	"time"

	"github.com/sapcc/baremetal_temper/pkg/node"
	"github.com/sapcc/baremetal_temper/pkg/temper"
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
		t := temper.New(cfg, context.Background(), netboxStatus)
		if len(nodes) > 0 {
			for _, n := range nodes {
				wg.Add(1)
				go hardwareCheckExec(n, t, &wg)
			}
		}
		wg.Wait()
		log.Info("check completed")
	},
}

func hardwareCheckExec(n string, t *temper.Temper, wg *sync.WaitGroup) {
	defer wg.Done()
	node, err := node.New(n, cfg)
	if err != nil {
		log.Errorf("error node %s: %s", n, err.Error())
		return
	}
	err = node.RunHardwareChecks()
	if err != nil {
		log.Errorf("error node %s: %s", n, err.Error())
		return
	}
}

var cableCheck = &cobra.Command{
	Use:   "cablecheck",
	Short: "runs a cable check (lldp)",
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		t := temper.New(cfg, context.Background(), netboxStatus)
		if len(nodes) > 0 {
			for _, n := range nodes {
				wg.Add(1)
				go cableCheckExec(n, t, &wg)
			}
		}
		wg.Wait()
		log.Info("cable check completed")
	},
}

func cableCheckExec(n string, t *temper.Temper, wg *sync.WaitGroup) {
	defer wg.Done()
	node, err := node.New(n, cfg)
	if err != nil {
		log.Errorf("error node %s: %s", n, err.Error())
		return
	}
	if bootImg && cfg.Redfish.BootImage != nil {
		node.BootImage()
		time.Sleep(5 * time.Minute)
	}
	if err = node.RunAristaCheck(); err != nil {
		log.Errorf("error node %s: %s", n, err.Error())
	}
	if err = node.RunACICheck(); err != nil {
		log.Errorf("error node %s: %s", n, err.Error())
	}
}

func init() {
	diagCmd.PersistentFlags().BoolVar(&bootImg, "bootImage", false, "boots an image before running cablecheck")

	diagCmd.AddCommand(cableCheck)
	diagCmd.AddCommand(hardwareCheck)
	rootCmd.AddCommand(diagCmd)
}
