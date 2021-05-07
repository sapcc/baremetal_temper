package cmd

import (
	"sync"

	"github.com/sapcc/baremetal_temper/pkg/model"
	"github.com/sapcc/baremetal_temper/pkg/temper"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var (
	user           string
	testDeployment bool
)

var baremetalCmd = &cobra.Command{
	Use:   "baremetal",
	Short: "interact with baremetal (ironic) service",
}

var create = &cobra.Command{
	Use:   "create",
	Short: "Triggers a baremetal node create",
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
				tasks := make([]func(n *model.Node) error, 0)
				tasks = append(tasks, c.Openstack.Create()...)
				if testDeployment {
					tasks = append(tasks, c.Openstack.DeploymentTest()...)
				}
				tasks = append(tasks, c.Openstack.Prepare)
				go temperNode(t, n, tasks, &wg)
			}
		}
		wg.Wait()
		log.Info("create completed")
	},
}

var test = &cobra.Command{
	Use:   "test",
	Short: "Triggers prepare and test tasks",
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
				go temperNode(t, n, c.Openstack.DeploymentTest(), &wg)
			}
		}
		wg.Wait()
		log.Info("test completed")
	},
}

var validate = &cobra.Command{
	Use:   "validate",
	Short: "Triggers a baremetal node validation",
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
				go temperNode(t, n, []func(n *model.Node) error{c.Openstack.Validate}, &wg)
			}
		}
		wg.Wait()
		log.Info("validate completed")
	},
}

var prepare = &cobra.Command{
	Use:   "prepare",
	Short: "Triggers a baremetal node preparation",
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
				go temperNode(t, n, []func(n *model.Node) error{c.Openstack.Prepare}, &wg)
			}
		}
		wg.Wait()
		log.Info("validate completed")
	},
}

func init() {
	create.PersistentFlags().BoolVar(&testDeployment, "testDeployment", false, "run baremetal create with a test deployment")

	baremetalCmd.AddCommand(create)
	baremetalCmd.AddCommand(test)
	baremetalCmd.AddCommand(prepare)
	baremetalCmd.AddCommand(validate)

	rootCmd.AddCommand(baremetalCmd)
}
