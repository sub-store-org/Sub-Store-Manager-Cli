package main

import (
	"sub-store-manager-cli/cmd"
	"sub-store-manager-cli/lib"
)

func main() {
	lib.Init()
	cmd.Execute()
}
