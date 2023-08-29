package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"sub-store-manager-cli/lib"
	"sub-store-manager-cli/vars"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start a sub-store docker container by name",
	Run: func(cmd *cobra.Command, args []string) {
		var n string
		if len(args) == 0 {
			n = vars.DockerName
		} else {
			n = args[0]
		}
		startContainer(n)
	},
}

func startContainer(n string) {
	fmt.Println("start container", n)

	// 检查是否存在名字为n的容器
	isExist := false
	for _, c := range lib.GetSSMContainers() {
		if c.Name == n {
			c.Start()
			isExist = true
			break
		}
	}
	if !isExist {
		fmt.Printf("container %s not found, if you want to create a new one, please use `ssm new` command\n", n)
	}
}
