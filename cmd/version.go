package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version of powder-monkey
var Version string

// ReportVersion reports the binary version
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print version",
	Long:  "Shows current version of baremetal-temper",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
