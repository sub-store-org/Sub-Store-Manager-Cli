package vars

import (
	"log"
	"os"
	"path/filepath"
)

const (
	Version         = "0.0.8"
	DockerNameBE    = "ssm-backend"
	DockerNameFE    = "ssm-frontend"
	ContainerTypeFE = "frontend"
	ContainerTypeBE = "backend"
)

var (
	HomeDir    string
	AppDir     string
	AppFileDir string
	FEFileDir  string
	BEFileDir  string
	ConfigDir  string
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln("Failed to get user home directory: ", err)
		return
	}
	HomeDir = homeDir
	AppDir = filepath.Join(HomeDir, ".ssm")
	AppFileDir = filepath.Join(AppDir, "appFiles")
	FEFileDir = filepath.Join(AppFileDir, "frontend")
	BEFileDir = filepath.Join(AppFileDir, "backend")
	ConfigDir = filepath.Join(AppDir, "configs")
}
