package lib

import (
	"errors"
	"fmt"
	"os"
)

func checkExist(path string) bool {
	_, err := os.Stat(path)

	switch {
	case err == nil:
		return true
	case errors.Is(err, os.ErrNotExist):
		return false
	default:
		fmt.Println("Error while checking filesystem: ", err)
		os.Exit(1)
		return false
	}
}

func makeDir(path string) {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		fmt.Println("Failed to create directory: ", err)
		os.Exit(1)
	}
}

func makeFile(path string) {
	_, err := os.Create(path)
	if err != nil {
		fmt.Println("Failed to create file: ", err)
		os.Exit(1)
	}
}

func rmFile(path string) {
	err := os.Remove(path)
	if err != nil {
		fmt.Println("Failed to remove file: ", err)
		os.Exit(1)
	}
}
