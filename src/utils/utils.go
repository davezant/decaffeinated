package utils

import (
	"os"
	"strings"
)

func IsPath(path string) bool {
	if len(path) == 0 {
		return false
	}

	lastChar := path[len(path)-1]
	if lastChar == '/' || lastChar == '\\' {
		return true
	}
	return false
}

func DoesExecutableExist(path, executable string) bool {
	fullPath := path + string(os.PathSeparator) + executable

	info, err := os.Stat(fullPath)
	if err != nil {
		return false
	}

	if info.Mode().IsRegular() {
		return true
	}
	return false
}

func EqualIgnoreCase(a, b string) bool {
	return strings.EqualFold(a, b)
}
