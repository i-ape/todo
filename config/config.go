// cmd/config.go
package config

import (
	"os"
	"path/filepath"
)

var (
	TaskFile   = getEnv("TODO_PATH", defaultTaskFile())
	DisableFzf = getEnvBool("TODO_NO_FZF", false)
	EnableTui  = getEnvBool("TODO_TUI", false)
)

// defaultTaskFile returns the default task file path, preferring a user-specific directory
func defaultTaskFile() string {
	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Join(home, ".todo", "tasks.json")
	}
	return "tasks.json"
}

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

// GetTaskFilePath returns the path to the tasks file, ensuring the directory exists
func GetTaskFilePath() string {
	path := TaskFile
	// If path is relative, make it absolute relative to the current working directory
	if !filepath.IsAbs(path) {
		if cwd, err := os.Getwd(); err == nil {
			path = filepath.Join(cwd, path)
		}
	}
	// Ensure the parent directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		// Log error but return path anyway to allow the app to attempt file creation
		// Alternatively, you could panic or return an error in stricter scenarios
		// fmt.Fprintf(os.Stderr, "Failed to create directory %s: %v\n", dir, err)
	}
	return path
}
