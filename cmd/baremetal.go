package cmd

import (
	"github.com/sapcc/baremetal_temper/pkg/model"
	"github.com/sapcc/baremetal_temper/pkg/temper"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var (
	user string
)

var baremetalCmd = &cobra.Command{
	Use:   "baremetal",
	Short: "interact with baremetal (ironic) service",
}

var create = &cobra.Command{
	Use:   "create",
	Short: "Triggers a baremetal node create",
	Run: func(cmd *cobra.Command, args []string) {
		t := temper.New(cfg)
		c, err := t.GetClients(node)
		n, err := t.LoadNodeInfos(node)
		if err != nil {
			log.Fatal(err)
		}
		if err = t.TemperNode(&n, c.Openstack.Create()); err != nil {
			log.Fatal(err)
		}
	},
}

var test = &cobra.Command{
	Use:   "test",
	Short: "Triggers prepare and test tasks",
	Run: func(cmd *cobra.Command, args []string) {
		t := temper.New(cfg)
		c, err := t.GetClients(node)
		n, err := t.LoadNodeInfos(node)
		if err != nil {
			log.Fatal(err)
		}
		if err = t.TemperNode(&n, c.Openstack.TestAndPrepare()); err != nil {
			log.Fatal(err)
		}
	},
}

var validate = &cobra.Command{
	Use:   "validate",
	Short: "Triggers a baremetal node validation",
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
		if err = t.TemperNode(&n, []func(n *model.Node) error{c.Openstack.Validate}); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	baremetalCmd.AddCommand(create)
	baremetalCmd.AddCommand(test)

	rootCmd.AddCommand(baremetalCmd)
}
