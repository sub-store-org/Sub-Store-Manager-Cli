package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"sub-store-manager-cli/vars"
)

var inputVersion = ""
var inputName = ""
var inputPort = ""
var inputNetwork = ""
var inputType bool

var rootCmd = &cobra.Command{
	Use:   "ssm",
	Short: "A Sub-Store Manager CLI",
	Long:  `A Sub-Store Manager CLI for managing sub-store`,
}

func init() {
	rootCmd.AddCommand(versionCmd, lsCmd, newCmd, stopCmd, startCmd, deleteCmd, updateCmd)
}

func getType() (string, string) {
	if !inputType {
		return vars.DockerNameBE, vars.ContainerTypeBE
	} else {
		return vars.DockerNameFE, vars.ContainerTypeFE
	}
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}
