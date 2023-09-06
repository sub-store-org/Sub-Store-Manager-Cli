package lib

import (
	"fmt"
	"os"
)

func PrintError(text string, err error) {
	// 添加红色字体
	if err != nil {
		fmt.Println("\033[31mError: ", text, err, "\033[0m")
	} else {
		fmt.Println("\033[31mError: ", text, "\033[0m")
	}
	os.Exit(1)
}

func PrintInfo(info string) {
	fmt.Println("Info:", info)
}
