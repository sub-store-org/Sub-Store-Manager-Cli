package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ssm",
	Short: "A Sub-Store Manager CLI",
	Long:  `A Sub-Store Manager CLI for managing sub-store in Linux`,
}

func init() {
	rootCmd.AddCommand(versionCmd, lsCmd, newCmd)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}
