package utils

import (
	"os"
)

// CreateDirIfNotExists creates a directory if it does not exist
func CreateDirIfNotExists(dirPath string) {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
}
