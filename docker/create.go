package docker

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"sub-store-manager-cli/lib"
	"sub-store-manager-cli/vars"
)

// createImageFromPath 从指定路径创建镜像
func createImageFromPath(c *Container, bPath string, tPath string) {
	// 如果旧 tar 文件存在则删除
	if lib.CheckExist(tPath) {
		lib.RemoveFile(tPath)
	}
	err := lib.CreateTarArchive(bPath, tPath)
	if err != nil {
		lib.PrintError("Failed to create tar archive:", err)
	}
	defer os.Remove(tPath)

	tarFile, err := os.Open(tPath)
	if err != nil {
		lib.PrintError("Failed to open tar file:", err)
	}
	defer tarFile.Close()

	// 构建容器
	c.buildImage(tarFile)
}

// CreateFEImage 创建前端镜像
func createFEImage(c *Container) {
	// 打开 Dockerfile 文件
	buildDir := filepath.Join(vars.FEFileDir, c.Version)
	tarPath := filepath.Join(vars.FEFileDir, "temp.tar")

	createImageFromPath(c, buildDir, tarPath)
}

// CreateBEImage 创建后端镜像
func createBEImage(c *Container) {
	// 打开 Dockerfile 文件
	buildDir := filepath.Join(vars.BEFileDir, c.Version)
	tarPath := filepath.Join(vars.BEFileDir, "temp.tar")

	// 下载对应版本后端程序文件
	// packageJson := filepath.Join(versionDir, "package.json")
	// downloadFile("https://raw.githubusercontent.com/sub-store-org/Sub-Store/master/backend/package.json", packageJson)
	// lockfile := filepath.Join(versionDir, "pnpm-lock.yaml")
	// downloadFile("https://raw.githubusercontent.com/sub-store-org/Sub-Store/master/backend/pnpm-lock.yaml", lockfile)
	minFile := filepath.Join(buildDir, "sub-store.bundle.js")
	lib.DownloadFile(fmt.Sprintf("https://github.com/sub-store-org/Sub-Store/releases/download/%s/sub-store.bundle.js", c.Version), minFile)

	fmt.Println("Files downloaded successfully.")
	createImageFromPath(c, buildDir, tarPath)
}

// CreateImage 创建镜像
func (c *Container) CreateImage() {
	// if ImageIsExist(v) {
	//     fmt.Printf("The image %s is already exist, skip build. if you want rebuild it, please remove image first.\n", vars.DockerNameBE+":"+v)
	//     return
	// }

	fmt.Println("Start building docker image, please waiting...")

	if c.ImageName == "" || c.Version == "" {
		lib.PrintError("ImageName or Version is empty.", nil)
	}

	if c.DockerfileStr == "" {
		lib.PrintError("Not provide Dockerfile, please check.", nil)
	}

	writeDockerfileToOS(c.DockerfileStr, c.ContainerType, c.Version)
	switch c.ContainerType {
	case vars.ContainerTypeFE:
		createFEImage(c)
	case vars.ContainerTypeBE:
		createBEImage(c)
	default:
		lib.PrintError("Not support container type.", nil)
	}
}

func (c *Container) buildImage(tar *os.File) {
	// 构建容器
	buildOptions := types.ImageBuildOptions{
		Context:    tar,
		Dockerfile: "Dockerfile",
		Tags:       []string{c.ImageName + ":" + c.Version},
	}

	buildResponse, err := dc.ImageBuild(dcCtx, buildOptions.Context, buildOptions)
	if err != nil {
		lib.PrintError("Failed to build image:", err)
	}
	defer buildResponse.Body.Close()

	// 构建输出
	formatDockerOutput(buildResponse.Body)

	fmt.Println("\nPrune build cache...")
	_, err = dc.BuildCachePrune(dcCtx, types.BuildCachePruneOptions{})
	if err != nil {
		lib.PrintError("Failed to prune build cache:", err)
	}

	fmt.Println("Docker image build successfully.")
}

func (c *Container) StartImage() {
	fmt.Println("Start container...")

	// 创建一个容器并运行它
	containerConfig := &container.Config{
		Image: c.ImageName + ":" + c.Version,
	}

	hostConfig := &container.HostConfig{
		RestartPolicy: container.RestartPolicy{
			Name: "unless-stopped",
		},
	}

	networkConfig := &network.NetworkingConfig{}
	if c.Network != "" {
		_, err := c.GetNetworkID()
		if err != nil {
			os.Exit(1)
		}

		networkConfig = &network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				c.Network: {},
			},
		}
	}

	switch c.ContainerType {
	case vars.ContainerTypeFE:
		hostConfig.PortBindings = nat.PortMap{
			"80/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: c.HostPort,
				},
			},
		}
	case vars.ContainerTypeBE:
		configDir := filepath.Join(vars.ConfigDir, c.Name)
		hostConfig.Binds = append(hostConfig.Binds, configDir+":/app/config")
		hostConfig.PortBindings = nat.PortMap{
			"3000/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: c.HostPort,
				},
			},
		}
	}

	resp, err := dc.ContainerCreate(dcCtx, containerConfig, hostConfig, networkConfig, nil, c.Name)
	if err != nil {
		lib.PrintError("Failed to create container: ", err)
	}

	// 启动容器
	err = dc.ContainerStart(dcCtx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		fmt.Println("Failed to start container: ", err)
		err = dc.ContainerRemove(dcCtx, resp.ID, types.ContainerRemoveOptions{Force: true})
		if err != nil {
			lib.PrintError("Failed to kill container, please kill it manually.\n", err)
		}
		os.Exit(1)
	}

	now := time.Now()
	seconds := now.Unix()
	nanoseconds := now.UnixNano()
	untilValue := strconv.FormatInt(seconds, 10) + "." + strconv.FormatInt(nanoseconds, 10)
	pruneFilters := filters.NewArgs()
	pruneFilters.Add("until", untilValue)
	_, err = dc.ContainersPrune(dcCtx, pruneFilters)
	if err != nil {
		fmt.Printf("Failed to prune containers: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Container started successfully. You can use `ssm ls` to view the container status.")
}
