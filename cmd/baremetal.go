package cmd

import (
	"sync"

	"github.com/sapcc/baremetal_temper/pkg/node"
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
		err := loadNodes()
		if err != nil {
			log.Errorf("error loading nodes: %s", err.Error())
			return
		}
		if len(nodes) > 0 {
			for _, n := range nodes {
				nd, err := node.New(n, cfg)
				if err != nil {
					log.Errorf("error node %s: %s", n, err.Error())
					continue
				}
				nd.AddTask("temper_import-ironic")
				wg.Add(1)
				go nd.Temper(netboxStatus, &wg)
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
		err := loadNodes()
		if err != nil {
			log.Errorf("error loading nodes: %s", err.Error())
			return
		}
		for _, n := range nodes {
			node, err := node.New(n, cfg)
			if err != nil {
				log.Errorf("error node %s: %s", n, err.Error())
				continue
			}
			wg.Add(1)
			node.AddTask("temper_ironic-test-deployment")
			node.Temper(netboxStatus, &wg)
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
		err := loadNodes()
		if err != nil {
			log.Errorf("error loading nodes: %s", err.Error())
			return
		}

		for _, n := range nodes {
			node, err := node.New(n, cfg)
			if err != nil {
				log.Errorf("error node %s: %s", n, err.Error())
				continue
			}
			node.AddTask("validate_node")
			wg.Add(1)
			go node.Temper(netboxStatus, &wg)
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
		err := loadNodes()
		if err != nil {
			log.Errorf("error loading nodes: %s", err.Error())
			return
		}
		for _, n := range nodes {
			node, err := node.New(n, cfg)
			if err != nil {
				log.Errorf("error node %s: %s", n, err.Error())
				continue
			}
			node.AddTask("prepare_node")
			wg.Add(1)
			go node.Temper(netboxStatus, &wg)
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
