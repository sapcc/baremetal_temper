package cmd

import (
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
		if len(nodes) > 0 {
			for _, n := range nodes {
				nd, err := node.New(n, cfg)
				if err != nil {
					log.Errorf("error node %s: %s", n, err.Error())
					continue
				}
				if err = nd.Create(); err != nil {
					log.Errorf("error node %s: %s", n, err.Error())
					continue
				}
				if testDeployment {
					if err = nd.DeployTestInstance(); err != nil {
						log.Errorf("error node %s: %s", n, err.Error())
						continue
					}
				}
				if err = nd.Prepare(); err != nil {
					log.Errorf("error node %s: %s", n, err.Error())
					continue
				}
			}
		}
		log.Info("create completed")
	},
}

var test = &cobra.Command{
	Use:   "test",
	Short: "Triggers prepare and test tasks",
	Run: func(cmd *cobra.Command, args []string) {
		//var wg sync.WaitGroup
		if len(nodes) > 0 {
			for _, n := range nodes {
				node, err := node.New(n, cfg)
				if err != nil {
					log.Errorf("error node %s: %s", n, err.Error())
					continue
				}
				if err = node.DeployTestInstance(); err != nil {
					log.Errorf("error node %s: %s", n, err.Error())
					continue
				}
			}
		}
		log.Info("test completed")
	},
}

var validate = &cobra.Command{
	Use:   "validate",
	Short: "Triggers a baremetal node validation",
	Run: func(cmd *cobra.Command, args []string) {
		//var wg sync.WaitGroup
		if len(nodes) > 0 {
			for _, n := range nodes {
				node, err := node.New(n, cfg)
				if err != nil {
					log.Errorf("error node %s: %s", n, err.Error())
					continue
				}
				if err = node.Validate(); err != nil {
					log.Errorf("error node %s: %s", n, err.Error())
					continue
				}
			}
		}
		//wg.Wait()
		log.Info("validate completed")
	},
}

var prepare = &cobra.Command{
	Use:   "prepare",
	Short: "Triggers a baremetal node preparation",
	Run: func(cmd *cobra.Command, args []string) {
		//var wg sync.WaitGroup
		if len(nodes) > 0 {
			for _, n := range nodes {
				node, err := node.New(n, cfg)
				if err != nil {
					log.Errorf("error node %s: %s", n, err.Error())
					continue
				}
				if err = node.Prepare(); err != nil {
					log.Errorf("error node %s: %s", n, err.Error())
					continue
				}
			}
		}
		//wg.Wait()
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
