package lib

import (
	"archive/tar"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"sub-store-manager-cli/vars"
)

var (
	dcIsInit bool
	DC       *client.Client
	DCCtx    context.Context
)

func initDockerClient() {
	if dcIsInit {
		return
	}

	// 创建 Docker 客户端
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalln("Failed to create Docker client:", err)
	}

	DC = dockerClient
	DCCtx = context.Background()
	dcIsInit = true
}

type SSMContainer struct {
	Id          string
	Name        string
	HostPort    string
	Port        string
	NetworkType string
	Version     string
	Status      string
}

func GetSSMContainers() []SSMContainer {
	// 获取容器列表
	containers, err := DC.ContainerList(DCCtx, types.ContainerListOptions{All: true})
	if err != nil {
		log.Fatalln("Failed to list containers:", err)
	}

	var ssmList []SSMContainer

	// 遍历容器列表并解析镜像名称
	for _, c := range containers {
		imageNameParts := strings.Split(c.Image, ":")
		if len(imageNameParts) > 1 {
			n, v := imageNameParts[0], imageNameParts[1]
			if n == vars.DockerName {
				ssmC := SSMContainer{
					Id:      c.ID[0:24],
					Name:    c.Names[0][1:],
					Version: v,
				}
				ssmC.Status = strings.Split(c.Status, " ")[0]

				if ssmC.Status == "Up" {
					ssmC.HostPort = strconv.Itoa(int(c.Ports[0].PublicPort))
					ssmC.Port = strconv.Itoa(int(c.Ports[0].PrivatePort))
					ssmC.NetworkType = c.Ports[0].Type
				} else {
					ssmC.HostPort = "none"
				}

				ssmList = append(ssmList, ssmC)
			}
		}
	}

	return ssmList
}

func CreateDockerfile(v string) {
	// 检查 .ssm 目录是否存在，不存在则创建
	appDirIsExist := checkExist(vars.AppDir)
	if !appDirIsExist {
		makeDir(vars.AppDir)
	}

	// 检查版本目录是否存在，不存在则创建
	appFileDirIsExist := checkExist(vars.AppFileDir)
	if !appFileDirIsExist {
		makeDir(vars.AppFileDir)
	}

	// 创建版本目录
	versionDir := filepath.Join(vars.AppFileDir, v)
	versionDirIsExist := checkExist(versionDir)
	if !versionDirIsExist {
		makeDir(versionDir)
	}

	// 移除旧 Dockerfile 并创建新的
	dockerfilePath := filepath.Join(versionDir, "Dockerfile")
	dockerfileIsExist := checkExist(dockerfilePath)
	if dockerfileIsExist {
		rmFile(dockerfilePath)
	}
	makeFile(dockerfilePath)

	// 写入 Dockerfile
	content := `FROM node:16.18-slim

# 设置工作目录
WORKDIR /app

# 复制 package.json 和 pnpm-lock.yaml 并安装依赖项
COPY package.json pnpm-lock.yaml ./
RUN mkdir config && npm install -g pnpm && pnpm install && npm cache clean --force

# 复制项目文件
COPY . .

# 暴露端口
EXPOSE 3000

# 启动应用
CMD cd config && node ../sub-store.min.js
`

	err := os.WriteFile(dockerfilePath, []byte(content), 0666)
	if err != nil {
		fmt.Println("Failed to write Dockerfile: ", err)
		os.Exit(1)
	}

	fmt.Println("Dockerfile created successfully.")

	// 下载对应版本后端程序文件
	packageJson := filepath.Join(versionDir, "package.json")
	downloadFile("https://raw.githubusercontent.com/sub-store-org/Sub-Store/master/backend/package.json", packageJson)
	lockfile := filepath.Join(versionDir, "pnpm-lock.yaml")
	downloadFile("https://raw.githubusercontent.com/sub-store-org/Sub-Store/master/backend/pnpm-lock.yaml", lockfile)
	minFile := filepath.Join(versionDir, "sub-store.min.js")
	downloadFile(fmt.Sprintf("https://github.com/sub-store-org/Sub-Store/releases/download/%s/sub-store.min.js", v), minFile)

	fmt.Println("Files downloaded successfully.")
}

func (c *SSMContainer) Start() {
	fmt.Println("Start container...")

	err := DC.ContainerStart(DCCtx, c.Id, types.ContainerStartOptions{})
	if err != nil {
		fmt.Println("Failed to start container: ", err)
		os.Exit(1)
	}

	fmt.Println("Container started successfully.")
}

func (c *SSMContainer) Stop() {
	fmt.Println("Stop container...")

	err := DC.ContainerStop(DCCtx, c.Id, container.StopOptions{})
	if err != nil {
		fmt.Println("Failed to stop container: ", err)
		os.Exit(1)
	}

	fmt.Println("Container stopped successfully.")
}

func (c *SSMContainer) Delete() {
	fmt.Println("Delete container...")

	err := DC.ContainerRemove(DCCtx, c.Id, types.ContainerRemoveOptions{Force: true})
	if err != nil {
		fmt.Println("Failed to delete container: ", err)
		os.Exit(1)
	}

	fmt.Println("Container deleted successfully.")
}

func StartImage(v, name, port string) {
	fmt.Println("Start container...")

	configDir := filepath.Join(vars.ConfigDir, name)

	// 创建一个容器并运行它
	containerConfig := &container.Config{
		Image: vars.DockerName + ":" + v,
	}
	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"3000/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: port,
				},
			},
		},
		Binds: []string{configDir + ":/app/config"},
	}
	resp, err := DC.ContainerCreate(DCCtx, containerConfig, hostConfig, nil, nil, name)
	if err != nil {
		fmt.Println("Failed to create container: ", err)
		os.Exit(1)
	}

	err = DC.ContainerStart(DCCtx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		fmt.Println("Failed to start container: ", err)
		err := DC.ContainerRemove(DCCtx, resp.ID, types.ContainerRemoveOptions{Force: true})
		if err != nil {
			fmt.Println("Failed to kill container, please kill it manually.\n", err)
		}
		os.Exit(1)
	}

	now := time.Now()
	seconds := now.Unix()
	nanoseconds := now.UnixNano()
	untilValue := strconv.FormatInt(seconds, 10) + "." + strconv.FormatInt(nanoseconds, 10)
	pruneFilters := filters.NewArgs()
	pruneFilters.Add("until", untilValue)
	_, err = DC.ContainersPrune(DCCtx, pruneFilters)
	if err != nil {
		fmt.Printf("Failed to prune containers: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Container started successfully. You can use `ssm ls` to view the container status.")
}

func BuildContainer(v string) {
	if ImageIsExist(v) {
		fmt.Printf("The image %s is already exist, skip build. if you want rebuild it, please remove image first.\n", vars.DockerName+":"+v)
		return
	}

	fmt.Println("Start building docker image, please waiting...")

	// 打开 Dockerfile 文件
	buildDir := filepath.Join(vars.AppFileDir, v)
	tarPath := filepath.Join(vars.AppDir, "temp.tar")

	// 如果旧 tar 文件存在则删除
	if checkExist(tarPath) {
		rmFile(tarPath)
	}
	err := createTarArchive(buildDir, tarPath)
	if err != nil {
		fmt.Printf("Failed to create tar archive: %s\n", err)
		os.Exit(1)
	}
	defer os.Remove(tarPath)

	tarFile, err := os.Open(tarPath)
	if err != nil {
		fmt.Printf("Failed to open tar archive: %s\n", err)
		os.Exit(1)
	}
	defer tarFile.Close()

	// 构建容器
	buildOptions := types.ImageBuildOptions{
		Context:    tarFile,
		Dockerfile: "Dockerfile",
		Tags:       []string{vars.DockerName + ":" + v},
	}

	buildResponse, err := DC.ImageBuild(DCCtx, buildOptions.Context, buildOptions)
	if err != nil {
		fmt.Printf("Failed to build image: %s\n", err)
		os.Exit(1)
	}
	defer buildResponse.Body.Close()

	// 构建输出
	formatDockerOutput(buildResponse.Body)

	fmt.Println("\nPrune build cache...")
	_, err = DC.BuildCachePrune(DCCtx, types.BuildCachePruneOptions{})
	if err != nil {
		fmt.Printf("Failed to prune build cache: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Docker image build successfully.")
}

func ImageIsExist(v string) bool {
	// 检查镜像是否存在
	images, err := DC.ImageList(DCCtx, types.ImageListOptions{All: true})
	if err != nil {
		fmt.Println("Failed to list images:", err)
		os.Exit(1)
	}

	for _, image := range images {
		for _, repoTag := range image.RepoTags {
			if repoTag == vars.DockerName+":"+v {
				return true
			}
		}
	}

	return false
}

func createTarArchive(srcPath, destPath string) error {
	tarFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer tarFile.Close()

	tarWriter := tar.NewWriter(tarFile)
	defer tarWriter.Close()

	return filepath.Walk(srcPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(srcPath, filePath)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		header, err := tar.FileInfoHeader(info, relPath)
		if err != nil {
			return err
		}

		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(tarWriter, file)
		if err != nil {
			return err
		}

		return nil
	})
}

func formatDockerOutput(body io.ReadCloser) {
	// 解析构建输出并格式化输出
	decoder := json.NewDecoder(body)
	for {
		var message map[string]interface{}
		if err := decoder.Decode(&message); err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Failed to decode JSON:", err)
			os.Exit(1)
		}

		if stream, ok := message["stream"].(string); ok {
			fmt.Print(stream)
		}
		// else if status, ok := message["status"].(string); ok {
		//     fmt.Print(status)
		//     if progressDetail, ok := message["progressDetail"].(map[string]interface{}); ok {
		//         if _, exists := progressDetail["current"]; exists {
		//             fmt.Printf(" %v/%v", progressDetail["current"], progressDetail["total"])
		//         }
		//     }
		//     fmt.Print("\n")
		// } else if errorMessage, ok := message["error"].(string); ok {
		//     fmt.Printf("Error: %s\n", errorMessage)
		// }
	}
}
