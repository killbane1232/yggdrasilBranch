package utils

import (
	"os"
)

// GetServiceName возвращает имя сервиса из переменной окружения или hostname
func GetServiceName() string {
	branchName := os.Getenv("BRANCHNAME")
	if branchName != "" {
		return branchName
	}

	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}
