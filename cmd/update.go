package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"sub-store-manager-cli/docker"
	"sub-store-manager-cli/lib"
	"sub-store-manager-cli/vars"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update a sub-store docker container",
	Run: func(cmd *cobra.Command, args []string) {
		updateContainer()
	},
}

func init() {
	updateCmd.Flags().StringVarP(&inputVersion, "version", "v", "", "The target version to update")
	updateCmd.Flags().StringVarP(&inputName, "name", "n", "", "The target sub-store container name to update")
}

func updateContainer() {
	name := inputName
	if name == "" {
		name = vars.DockerNameBE
	}

	oldContainer, isExist := docker.GetContainerByName(name)
	if !isExist {
		lib.PrintError("The container does not exist.", nil)
	}

	c := docker.Container{
		Name:          inputName,
		ImageName:     oldContainer.ImageName,
		ContainerType: oldContainer.ContainerType,
		HostPort:      oldContainer.HostPort,
		Version:       inputVersion,
	}

	// 检查指定版本
	if inputVersion == "" {
		c.SetLatestVersion()
		fmt.Println("No version specified, using the latest version")
	} else {
		if c.ContainerType == vars.ContainerTypeFE {
			c.SetLatestVersion()
			lib.PrintInfo("The version flag is ignored when target is a front-end container.")
		} else {
			isValid := c.CheckVersionValid()
			if !isValid {
				lib.PrintError("The version is not valid.", nil)
			}
		}
	}

	// 如果当前运行的版本就是目标版本，则不更新
	if oldContainer.Version == c.Version {
		lib.PrintInfo("The current version is the target version, no need to update.")
		return
	}

	// 获取旧容器的端口信息
	if p, err := oldContainer.GetPortInfo(); err != nil {
		lib.PrintError("Failed to get port info:", err)
	} else {
		c.HostPort = p.Public
	}

	c.SetDockerfile("")
	c.CreateImage()

	// 删除旧容器, 启动新容器
	oldContainer.Stop()
	oldContainer.Delete()
	c.StartImage()
}
