package vars

import (
	"log"
	"os"
	"path/filepath"
)

const (
	Version    = "0.0.3"
	DockerName = "sub-store-manager-backend"
)

var (
	HomeDir    string
	AppDir     string
	AppFileDir string
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
	ConfigDir = filepath.Join(AppDir, "configs")
}
