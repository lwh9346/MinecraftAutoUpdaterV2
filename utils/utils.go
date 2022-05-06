package utils

import (
	"os"
	"path/filepath"
)

func RemoveEmptyDirectories(path string) bool {
	de, err := os.ReadDir(path)
	if err != nil {
		return false
	}
	empty := true
	for _, d := range de {
		if d.IsDir() {
			empty = empty && RemoveEmptyDirectories(filepath.Join(path, d.Name()))
		}
	}
	if empty {
		os.Remove(path)
	}
	return empty
}
