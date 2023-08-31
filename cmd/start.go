package cmd

import (
	"github.com/spf13/cobra"
	"sub-store-manager-cli/docker"
	"sub-store-manager-cli/vars"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start a sub-store docker container by name",
	Run: func(cmd *cobra.Command, args []string) {
		var n string
		if len(args) == 0 {
			n = vars.DockerNameBE
		} else {
			n = args[0]
		}
		startContainer(n)
	},
}

func startContainer(n string) {
	c := docker.GetContainerByName(n)
	c.Start()
}
