package lib

import (
	"archive/tar"
	"errors"
	"io"
	"os"
	"path/filepath"
)

func CheckExist(path string) bool {
	_, err := os.Stat(path)

	switch {
	case err == nil:
		return true
	case errors.Is(err, os.ErrNotExist):
		return false
	default:
		PrintError("Error while checking filesystem: ", err)
		return false
	}
}

func MakeDir(path string) {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		PrintError("Failed to create directory: ", err)
	}
}

func MakeFile(path string) {
	_, err := os.Create(path)
	if err != nil {
		PrintError("Failed to create file: ", err)
	}
}

func RemoveFile(path string) {
	err := os.Remove(path)
	if err != nil {
		PrintError("Failed to remove file: ", err)
	}
}

func RemoveDir(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		PrintError("Failed to remove directory: ", err)
	}
}

func CreateTarArchive(srcPath, destPath string) error {
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
