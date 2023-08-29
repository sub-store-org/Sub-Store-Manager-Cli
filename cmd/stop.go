package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
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
			n = vars.DockerName
		} else {
			n = args[0]
		}
		stopContainer(n)
	},
}

func stopContainer(n string) {
	fmt.Println("stop container", n)

	// 检查是否存在正在运行的名字为n的容器
	isExist := false
	for _, c := range lib.GetSSMContainers() {
		if c.Name == n {
			c.Stop()
			isExist = true
			break
		}
	}
	if !isExist {
		fmt.Printf("container %s not found\n", n)
	}
}
