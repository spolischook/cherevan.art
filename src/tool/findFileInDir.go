package tool

import (
	"fmt"
	"os"
	"path/filepath"
)

func FindFileInDir(dirPath string, fileName string) (string, error) {
	var filePath string
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if info.Name() == fileName {
			filePath = path
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	if filePath == "" {
		return "", fmt.Errorf("file not found: %s\n", fileName)
	}
	return filePath, nil
}
