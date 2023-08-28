package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"sub-store-manager-cli/vars"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of ssm",
	Run: func(cmd *cobra.Command, args []string) {
		version()
	},
}

func version() {
	fmt.Printf("Sub-Store Manager CLI v%s\n", vars.Version)
}
