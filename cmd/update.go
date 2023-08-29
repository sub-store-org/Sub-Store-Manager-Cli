package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"sub-store-manager-cli/lib"
	"sub-store-manager-cli/vars"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update a sub-store container version",
	Run: func(cmd *cobra.Command, args []string) {
		var n string
		if len(args) == 0 {
			n = vars.DockerName
		} else {
			n = args[0]
		}
		updateContainer(n)
	},
}

func init() {
	updateCmd.Flags().StringVarP(&inputVersion, "version", "v", "", "The target version to update")
}

func updateContainer(n string) {
	v := inputVersion
	if v == "" {
		if latest, err := lib.GetLatestVersionString(); err != nil {
			fmt.Println("get latest version error:", err)
			os.Exit(1)
		} else {
			v = latest
		}
	}

	// 找到目标容器
	var c lib.SSMContainer
	for _, container := range lib.GetSSMContainers() {
		if container.Name == n {
			c = container
			break
		}
	}
	if c.Name == "" {
		fmt.Println("container not found, please check the name")
		os.Exit(1)
	}
	if c.Version == v {
		fmt.Println("the container is already the same version, no need to update")
		os.Exit(0)
	}

	p := c.HostPort

	c.Delete()

	// 查询是否已有对应版本的镜像
	lib.CreateDockerfile(v)
	lib.BuildContainer(v)
	lib.StartImage(v, n, p)

	fmt.Println("Container updated successfully.")
}
