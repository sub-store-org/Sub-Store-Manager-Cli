package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"sub-store-manager-cli/vars"
)

var inputVersion = ""
var inputName = ""
var inputPort = ""
var inputType = ""

var rootCmd = &cobra.Command{
	Use:   "ssm",
	Short: "A Sub-Store Manager CLI",
	Long:  `A Sub-Store Manager CLI for managing sub-store`,
}

func init() {
	rootCmd.AddCommand(versionCmd, lsCmd, newCmd, stopCmd, startCmd, updateCmd, deleteCmd)
}

func getType() (string, string) {
	if inputType == "" || inputType == "be" {
		return vars.DockerNameBE, vars.ContainerTypeBE
	} else if inputType == "fe" {
		return vars.DockerNameFE, vars.ContainerTypeFE
	} else {
		return "", ""
	}
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}
