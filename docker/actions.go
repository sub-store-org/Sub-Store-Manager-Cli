package docker

import (
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

// Start 启动容器
func (c *Container) Start() {
	fmt.Printf("Start container %s...", c.Name)

	err := dc.ContainerStart(dcCtx, c.DockerContainer.ID, types.ContainerStartOptions{})
	if err != nil {
		fmt.Println("Failed to start container: ", err)
		os.Exit(1)
	}

	fmt.Println("Container started successfully.")
}

// Stop 停止容器
func (c *Container) Stop() {
	fmt.Printf("Stop container %s...", c.Name)

	err := dc.ContainerStop(dcCtx, c.DockerContainer.ID, container.StopOptions{})
	if err != nil {
		fmt.Println("Failed to stop container: ", err)
		os.Exit(1)
	}

	fmt.Println("Container stopped successfully.")
}

// Delete 删除容器
func (c *Container) Delete() {
	fmt.Printf("Delete container %s...", c.Name)

	err := dc.ContainerRemove(dcCtx, c.DockerContainer.ID, types.ContainerRemoveOptions{Force: true})
	if err != nil {
		fmt.Println("Failed to delete container: ", err)
		os.Exit(1)
	}

	fmt.Println("Container deleted successfully.")
}
