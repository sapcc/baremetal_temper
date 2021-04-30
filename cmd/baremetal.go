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
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		t := temper.New(cfg)
		c, err := t.GetClients(node)
		n, err := t.LoadNodeInfos(node)
		if err != nil {
			log.Fatal(err)
		}
		t.TemperNode(&n, []func(n *model.Node) error{c.Openstack.Create})
		if err != nil {
			log.Fatal(err)
		}
	},
}

var deploy = &cobra.Command{
	Use:   "deploy",
	Short: "Triggers a baremetal node test deployment",
	Args:  cobra.ExactArgs(1),
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
		t.TemperNode(&n, []func(n *model.Node) error{c.Openstack.DeployTestInstance})
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	baremetalCmd.AddCommand(create)
	baremetalCmd.AddCommand(deploy)

	rootCmd.AddCommand(baremetalCmd)
}
