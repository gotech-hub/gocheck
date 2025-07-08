package scanner

import (
	"fmt"
	"os"
	"path/filepath"
)

func ScanDir(root string) []string {
	var files []string
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && filepath.Ext(path) == ".go" {
			fmt.Println("Scanning:", path)
			files = append(files, path)
		}
		return nil
	})
	return files
}
