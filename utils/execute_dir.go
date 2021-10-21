package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func GetExecuteFilePath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	srcPath := filepath.Dir(ex)
	return formatPath(srcPath), nil
}

func GetWorkingPath() (string, error) {
	srcPath, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return formatPath(srcPath), nil
}

func formatPath(srcPath string) string {
	if runtime.GOOS == "windows" {
		srcPath = strings.ReplaceAll(srcPath, "\\", "/")
	}
	return srcPath
}