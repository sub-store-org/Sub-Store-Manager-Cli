package main

import (
	"sub-store-manager-cli/cmd"
	"sub-store-manager-cli/docker"
	"sub-store-manager-cli/lib"
)

func main() {
	docker.InitDockerClient()
	lib.InitHttpClient()
	cmd.Execute()
}
