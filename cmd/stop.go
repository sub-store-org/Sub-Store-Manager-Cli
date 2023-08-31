package cmd

import (
	"github.com/spf13/cobra"
	"sub-store-manager-cli/docker"
	"sub-store-manager-cli/lib"
	"sub-store-manager-cli/vars"
)

// var inputName = ""

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "stop a sub-store docker container by name",
	Run: func(cmd *cobra.Command, args []string) {
		var n string
		if len(args) == 0 {
			n = vars.DockerNameBE
		} else {
			n = args[0]
		}
		stopContainer(n)
	},
}

func stopContainer(n string) {
	c, isExist := docker.GetContainerByName(n)
	if !isExist {
		lib.PrintError("container not found", nil)
	}
	c.Stop()
}
