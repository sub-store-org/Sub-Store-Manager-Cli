package docker

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Masterminds/semver/v3"
	"github.com/docker/docker/api/types"
	"sub-store-manager-cli/lib"
	"sub-store-manager-cli/vars"
)

func formatDockerOutput(body io.ReadCloser) {
	// 解析构建输出并格式化输出
	decoder := json.NewDecoder(body)
	for {
		var message map[string]interface{}
		if err := decoder.Decode(&message); err == io.EOF {
			break
		} else if err != nil {
			lib.PrintError("Failed to decode message:", err)
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

func ImageIsExist(n string, v string) bool {
	// 检查镜像是否存在
	images, err := dc.ImageList(dcCtx, types.ImageListOptions{All: true})
	if err != nil {
		lib.PrintError("Failed to list images:", err)
	}

	for _, image := range images {
		for _, repoTag := range image.RepoTags {
			if repoTag == n+":"+v {
				return true
			}
		}
	}

	return false
}

func writeDockerfileToOS(d string, t string, v string) {
	// 检查 .ssm 目录是否存在，不存在则创建
	appDirIsExist := lib.CheckExist(vars.AppDir)
	if !appDirIsExist {
		lib.MakeDir(vars.AppDir)
	}

	// 检查 appFile 目录是否存在，不存在则创建
	appFileDirIsExist := lib.CheckExist(vars.AppFileDir)
	if !appFileDirIsExist {
		lib.MakeDir(vars.AppFileDir)
	}

	// 检查资源文件目录是否存在，不存在则创建
	var workDir string
	switch t {
	case vars.ContainerTypeFE:
		feFileDirIsExist := lib.CheckExist(vars.FEFileDir)
		if !feFileDirIsExist {
			lib.MakeDir(vars.FEFileDir)
		}
		workDir = vars.FEFileDir
	case vars.ContainerTypeBE:
		beFileDirIsExist := lib.CheckExist(vars.BEFileDir)
		if !beFileDirIsExist {
			lib.MakeDir(vars.BEFileDir)
		}
		workDir = vars.BEFileDir
	}

	// 移除旧版本目录 创建新版本目录
	versionDir := filepath.Join(workDir, v)
	lib.RemoveDir(versionDir)
	lib.MakeDir(versionDir)

	// 写入 Dockerfile
	dockerfilePath := filepath.Join(versionDir, "Dockerfile")
	err := os.WriteFile(dockerfilePath, []byte(d), 0666)
	if err != nil {
		lib.PrintError("Failed to write Dockerfile: ", err)
	}

	fmt.Println("Dockerfile created successfully.")
}

func (c *Container) SetDefaultName() {
	switch c.ContainerType {
	case vars.ContainerTypeFE:
		c.Name = vars.DockerNameFE
	case vars.ContainerTypeBE:
		c.Name = vars.DockerNameBE
	}
}

func (c *Container) SetDefaultPort() {
	switch c.ContainerType {
	case vars.ContainerTypeFE:
		c.HostPort = "80"
	case vars.ContainerTypeBE:
		c.HostPort = "3000"
	}
}

func (c *Container) SetDockerfile(flag string) {
	switch c.ContainerType {
	case vars.ContainerTypeFE:
		c.DockerfileStr = DockerfileStr.FE
	case vars.ContainerTypeBE:
		target, err := semver.NewVersion(c.Version)
		if err != nil {
			lib.PrintError("Failed to parse target version.", err)
		}
		bundleRule, err := semver.NewConstraint(">= 2.14.40")
		if err != nil {
			lib.PrintError("Failed to parse bundle version.", err)
		}
		envRule, err := semver.NewConstraint(">= 2.14.49")
		if err != nil {
			lib.PrintError("Failed to parse env version.", err)
		}

		if canUseBundle, _ := bundleRule.Validate(target); !canUseBundle {
			lib.PrintError("The version is not supported, please use a version after 2.14.40.", nil)
		} else if canUseEnv, _ := envRule.Validate(target); !canUseEnv {
			c.DockerfileStr = DockerfileStr.Node
		} else {
			c.DockerfileStr = DockerfileStr.NodeWithDataEnv
		}

		// switch flag {
		// case "node":
		//     c.DockerfileStr = DockerfileStr.Node
		// default:
		//     lib.PrintError("Not support backend container type.", nil)
		// }
	}
}
