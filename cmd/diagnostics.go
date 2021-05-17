package cmd

import (
	"context"
	"sync"
	"time"

	"github.com/sapcc/baremetal_temper/pkg/diagnostics"
	"github.com/sapcc/baremetal_temper/pkg/model"
	"github.com/sapcc/baremetal_temper/pkg/temper"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var diagCmd = &cobra.Command{
	Use:   "diagnostics",
	Short: "interact with diagnostics services",
}

var cableCheck = &cobra.Command{
	Use:   "hardwarecheck",
	Short: "runs a vendor specific hardware check",
	Run: func(cmd *cobra.Command, args []string) {
		ctxLogger := log.WithFields(log.Fields{
			"cli": "hardwareCheck",
		})
		var wg sync.WaitGroup
		t := temper.New(cfg, ctxLogger.Context, netboxStatus)
		if len(nodes) > 0 {
			for _, n := range nodes {
				c, err := t.GetClients(n)
				if err != nil {
					log.Errorf("error node %s: %s", n, err.Error())
					continue
				}
				tasks, err := diagnostics.GetHardwareCheckTasks(*c.Redfish.ClientConfig, cfg, ctxLogger)
				if err != nil {
					log.Errorf("error node %s: %s", n, err.Error())
					continue
				}
				wg.Add(1)
				temperNode(t, n, tasks, &wg)
			}
		}
		wg.Wait()
		log.Info("check completed")
	},
}

var hardwareCheck = &cobra.Command{
	Use:   "cablecheck",
	Short: "runs a cable check (lldp)",
	Run: func(cmd *cobra.Command, args []string) {
		ctxLogger := log.WithFields(log.Fields{
			"cli": "cableCheck",
		})
		var wg sync.WaitGroup
		t := temper.New(cfg, context.Background(), netboxStatus)
		if len(nodes) > 0 {
			for _, n := range nodes {
				wg.Add(1)
				tasks := make([]func(n *model.Node) error, 0)

				if bootImg && cfg.Redfish.BootImage != nil {
					c, err := t.GetClients(n)
					if err != nil {
						log.Errorf("error node %s: %s", n, err.Error())
						continue
					}
					tasks = append(tasks, c.Redfish.BootImage, t.GetTimeoutTask(5*time.Minute))
				}
				tasks = append(tasks, diagnostics.GetCableCheckTasks(cfg, ctxLogger)...)
				go temperNode(t, n, tasks, &wg)
			}
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
