package cmd

import (
	"github.com/sapcc/baremetal_temper/pkg/diagnostics"
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
		t := temper.New(cfg)
		c, err := t.GetClients(node)
		if err != nil {
			log.Fatal(err)
		}
		n, err := t.LoadNodeInfos(node)
		if err != nil {
			log.Fatal(err)
		}
		tasks, err := diagnostics.GetHardwareCheckTasks(*c.Redfish.ClientConfig, cfg, ctxLogger)
		if err != nil {
			log.Fatal(err)
		}
		if err = t.TemperNode(&n, tasks); err != nil {
			log.Fatal(err)
		}
		log.Info("node synced")
	},
}

var hardwareCheck = &cobra.Command{
	Use:   "cablecheck",
	Short: "runs a cable check (lldp)",
	Run: func(cmd *cobra.Command, args []string) {
		ctxLogger := log.WithFields(log.Fields{
			"cli": "cableCheck",
		})
		t := temper.New(cfg)
		n, err := t.LoadNodeInfos(node)
		if err != nil {
			log.Fatal(err)
		}
		if err = t.TemperNode(&n, diagnostics.GetCableCheckTasks(cfg, ctxLogger)); err != nil {
			log.Fatal(err)
		}

		log.Info("cable check successful")
	},
}

func init() {
	diagCmd.AddCommand(cableCheck)
	diagCmd.AddCommand(hardwareCheck)
	rootCmd.AddCommand(diagCmd)
}
