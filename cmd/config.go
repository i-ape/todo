// cmd/config.go
package main

import (
	"os"
	"path/filepath"
)

var (
	DisableFzf bool
	EnableTui  bool
)

// Default filename for storing tasks
var TaskFile = "tasks.json"

// Optional: override via $TODO_PATH env var
func GetTaskFilePath() string {
	if custom := os.Getenv("TODO_PATH"); custom != "" {
		return custom
	}
	return filepath.Join(".", TaskFile)
}
