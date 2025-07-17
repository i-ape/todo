// cmd/config.go
package config

import (
	"os"
	"path/filepath"
)

var (
	TaskFile   = getEnv("TODO_PATH", "tasks.json")
	DisableFzf = getEnvBool("TODO_NO_FZF", false)
	EnableTui  = getEnvBool("TODO_TUI", false)
)

// Default filename for storing tasks
var TaskFile = "tasks.json"

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		return value == "1" || value == "true"
	}
	return fallback
}

// Optional: override via $TODO_PATH env var
func GetTaskFilePath() string {
	if custom := os.Getenv("TODO_PATH"); custom != "" {
		return custom
	}
	return filepath.Join(".", TaskFile)
}
