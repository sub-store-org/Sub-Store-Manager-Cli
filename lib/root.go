package lib

import (
	"fmt"
	"net"
)

func Init() {
	initDockerClient()
	initHttpClient()
}

// CheckPort 检查端口是否可用，可用-true 不可用-false
func CheckPort(port string) bool {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return false
	}
	defer listener.Close()
	return true
}
