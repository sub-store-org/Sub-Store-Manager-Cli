package cmd

import (
	"fmt"
	"os"

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
	newCmd.Flags().StringVarP(&inputName, "name", "n", vars.DockerNameBE, "The container name")
	newCmd.Flags().StringVarP(&inputPort, "port", "p", "3000", "The port to expose")
	newCmd.Flags().StringVarP(&inputType, "type", "t", vars.ContainerTypeBE, "The target type to create a sub-store container")
}

func newContainer() {
	// 检查指定版本
	var target string
	if inputVersion == "" {
		if t, err := lib.GetLatestVersionString(); err != nil {
			fmt.Printf("Failed to get latest version: %s\n", err)
			os.Exit(1)
		} else {
			target = t
			fmt.Printf("No version specified, using the latest version %s\n", target)
		}
	} else {
		isValid := false
		for _, v := range lib.GetVersionsString() {
			if v == inputVersion {
				isValid = true
				break
			}
		}

		if !isValid {
			fmt.Printf("The version %s is invalid，please select one of version in https://github.com/sub-store-org/Sub-Store/releases\n", inputVersion)
			os.Exit(1)
		} else {
			target = inputVersion
		}
	}

	// 检查是否已运行一个同名容器

	for _, c := range docker.GetAllContainers() {
		name := c.Names[0][1:]
		if name == inputName {
			lib.PrintError(fmt.Sprintf("The container %s is already exist, if you want run another backend at sametime, please specifed a container name.\n", name), nil)
		}
	}

	// 检查端口
	// if !lib.CheckPort(inputPort) {
	//     fmt.Printf("The port %s is already in use, please specify another port.\n", inputPort)
	//     os.Exit(1)
	// }

	imageName, imageType := getType()
	c := docker.Container{
		Name:          inputName,
		ImageName:     imageName,
		Version:       target,
		HostPort:      inputPort,
		ContainerType: imageType,
		DockerfileStr: docker.DockerfileStr.Node,
	}

	c.CreateImage()
	c.StartImage()
}
