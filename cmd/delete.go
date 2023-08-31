package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"sub-store-manager-cli/docker"
	"sub-store-manager-cli/vars"
)

var deleteConfig bool

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete a sub-store docker container by name",
	Run: func(cmd *cobra.Command, args []string) {
		var n string
		if len(args) == 0 {
			n = vars.DockerNameBE
		} else {
			n = args[0]
		}
		deleteContainer(n)
	},
}

func init() {
	deleteCmd.Flags().BoolVarP(&deleteConfig, "clear", "c", false, "delete config file simultaneously")
}

func deleteContainer(n string) {
	c := docker.GetContainerByName(n)
	c.Delete()

	if c.ContainerType == vars.ContainerTypeBE && deleteConfig {
		err := os.RemoveAll(filepath.Join(vars.ConfigDir, n))
		if err != nil {
			fmt.Printf("container %s is deleted, but failed to clear config file: %s", n, err)
			os.Exit(0)
		}
	}
}
