package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"sub-store-manager-cli/lib"
	"sub-store-manager-cli/vars"
)

var deleteConfig bool

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete a sub-store docker container by name",
	Run: func(cmd *cobra.Command, args []string) {
		var n string
		if len(args) == 0 {
			n = vars.DockerName
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
	fmt.Println("delete container", n)

	// 检查是否存在正在运行的名字为n的容器
	isExist := false
	for _, c := range lib.GetSSMContainers() {
		if c.Name == n {
			c.Delete()
			// 删除配置文件
			if deleteConfig {
				err := os.RemoveAll(filepath.Join(vars.ConfigDir, n))
				if err != nil {
					fmt.Println("container is deleted, but failed to clear config file:", err)
				}
			}
			isExist = true
		}
	}
	if !isExist {
		fmt.Printf("container %s not found\n", n)
	}
}
