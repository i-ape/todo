// config/config.go
package config

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	TaskFile            = getEnv("TODO_PATH", defaultTaskFile())
	BackupTaskFile      = getEnv("TODO_BACKUP_PATH", defaultBackupTaskFile())
	DisableFzf          = getEnvBool("TODO_NO_FZF", false)
	EnableTui           = getEnvBool("TODO_TUI", false)
	DefaultOutputFormat = getEnv("TODO_OUTPUT", "text")
	TuiTheme            = getEnv("TODO_TUI_THEME", "dark")
	DefaultSortOrder    = getEnv("TODO_SORT", "id")
	DefaultPriority     = getEnv("TODO_DEFAULT_PRIORITY", "medium")
	FzfArgs             = getEnv("TODO_FZF_ARGS", "--prompt='Task> '")
	MaxTasks            = getEnvInt("TODO_MAX_TASKS", 1000)
)

func defaultTaskFile() string {
	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Join(home, ".todo", "tasks.json")
	}
	return "tasks.json"
}

func defaultBackupTaskFile() string {
	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Join(home, ".todo", "tasks.json.bak")
	}
	return "tasks.json.bak"
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

func getEnvInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return fallback
}

func GetTaskFilePath() string {
	path := TaskFile
	if !filepath.IsAbs(path) {
		if cwd, err := os.Getwd(); err == nil {
			path = filepath.Join(cwd, path)
		}
	}
	dir := filepath.Dir(path)
	_ = os.MkdirAll(dir, 0755)
	return path
}

func GetBackupTaskFilePath() string {
	path := BackupTaskFile
	if !filepath.IsAbs(path) {
		if cwd, err := os.Getwd(); err == nil {
			path = filepath.Join(cwd, path)
		}
	}
	dir := filepath.Dir(path)
	_ = os.MkdirAll(dir, 0755)
	return path
}

func GetOutputFormat() string {
	return DefaultOutputFormat
}

func GetTuiTheme() string {
	return TuiTheme
}

func GetSortOrder() string {
	return DefaultSortOrder
}

func GetDefaultPriority() string {
	return DefaultPriority
}

func GetFzfArgs() []string {
	return strings.Fields(FzfArgs)
}

func GetMaxTasks() int {
	return MaxTasks
}
