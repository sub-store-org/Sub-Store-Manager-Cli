package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"sub-store-manager-cli/docker"
	"sub-store-manager-cli/lib"
	"sub-store-manager-cli/vars"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "create a new sub-store docker container",
	Run: func(cmd *cobra.Command, args []string) {
		newContainer()
	},
}

func init() {
	newCmd.Flags().StringVarP(&inputVersion, "version", "v", "", "The target version to launch of the sub-store")
	newCmd.Flags().StringVarP(&inputName, "name", "n", "", "The container name")
	newCmd.Flags().StringVarP(&inputPort, "port", "p", "", "The port to expose")
	newCmd.Flags().BoolVarP(&inputType, "interface", "i", false, "The target type to create a sub-store container")
	newCmd.Flags().StringVarP(&inputNetwork, "network", "", "", "The docker network to connect")
	newCmd.Flags().BoolVarP(&inputPrivate, "private", "", false, "Host IP is private")
}

func newContainer() {
	imageName, imageType := getType()
	c := docker.Container{
		ImageName:     imageName,
		ContainerType: imageType,
		Network:       inputNetwork,
		Private:       inputPrivate,
	}

	// 检查是否已有同名容器
	if inputName == "" {
		c.SetDefaultName()
	} else {
		c.Name = inputName
	}
	_, isExist := docker.GetContainerByName(c.Name)
	if isExist {
		lib.PrintError("A container with the same name already exists.", nil)
	}

	// 检查指定版本
	if inputVersion == "" {
		c.SetLatestVersion()
		fmt.Println("No version specified, using the latest version")
	} else {
		if c.ContainerType == vars.ContainerTypeFE {
			c.SetLatestVersion()
			lib.PrintInfo("The version flag is ignored when creating a front-end container.")
		} else {
			c.Version = inputVersion
			isValid := c.CheckVersionValid()
			if !isValid {
				lib.PrintError("The version is not valid.", nil)
			}
		}
	}

	// 设置端口
	if inputPort == "" {
		c.SetDefaultPort()
	} else {
		c.HostPort = inputPort
	}

	// 检查端口
	if portOk := lib.CheckPort(c.HostPort); !portOk {
		lib.PrintError("The port is unavailable.", nil)
	}

	c.SetDockerfile("")
	c.CreateImage()
	c.StartImage()
}
