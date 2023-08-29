package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var inputVersion = ""
var inputName = ""
var inputPort = ""

var rootCmd = &cobra.Command{
	Use:   "ssm",
	Short: "A Sub-Store Manager CLI",
	Long:  `A Sub-Store Manager CLI for managing sub-store`,
}

func init() {
	rootCmd.AddCommand(versionCmd, lsCmd, newCmd, stopCmd, startCmd, updateCmd, deleteCmd)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}
