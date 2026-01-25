package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// LoadConfig загружает конфигурационный файл
func LoadConfig(name string) (*os.File, error) {
	possiblePaths := []string{
		filepath.Join("config", name),
		filepath.Join("./config", name),
		filepath.Join("/app/config", name),
	}
	
	for _, path := range possiblePaths {
		if file, err := os.Open(path); err == nil {
			return file, nil
		}
	}
	
	return nil, fmt.Errorf("config file %s not found", name)
}

// FindConfigFile находит конфигурационный файл
func FindConfigFile(filename string) (*os.File, error) {
	return LoadConfig(filename)
}
