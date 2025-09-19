// config/config.go
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	// File paths
	TaskFile       string
	BackupTaskFile string
	PathMode       string

	// Features
	DisableFzf          bool
	EnableTui           bool
	DefaultOutputFormat string
	TuiTheme            string
	DefaultSortOrder    string
	DefaultPriority     string
	FzfArgs             string
	MaxTasks            int
)

func init() {
	// Load .env before assigning config vars
	loadDotEnv(".env")

	// Debug: show what was read from env
	fmt.Println(">>> ENV TODO_PATH:", os.Getenv("TODO_PATH"))
	fmt.Println(">>> ENV TODO_BACKUP_PATH:", os.Getenv("TODO_BACKUP_PATH"))
	fmt.Println(">>> ENV TODO_PATH_MODE:", os.Getenv("TODO_PATH_MODE"))

	// Assign values
	TaskFile = getEnv("TODO_PATH", defaultTaskFile())
	BackupTaskFile = getEnv("TODO_BACKUP_PATH", defaultBackupTaskFile())
	PathMode = getEnv("TODO_PATH_MODE", "home")

	DisableFzf = getEnvBool("TODO_NO_FZF", false)
	EnableTui = getEnvBool("TODO_TUI", false)
	DefaultOutputFormat = getEnv("TODO_OUTPUT", "text")
	TuiTheme = getEnv("TODO_TUI_THEME", "dark")
	DefaultSortOrder = getEnv("TODO_SORT", "id")
	DefaultPriority = getEnv("TODO_DEFAULT_PRIORITY", "medium")
	FzfArgs = getEnv("TODO_FZF_ARGS", "--prompt='Task> '")
	MaxTasks = getEnvInt("TODO_MAX_TASKS", 1000)
}

// ----------------------------
// .env loader
// ----------------------------
func loadDotEnv(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return // silently skip if missing
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		// trim spaces, quotes, and CRLF (\r)
		val := strings.Trim(strings.TrimSpace(parts[1]), `"`)
		val = strings.ReplaceAll(val, "\r", "")
		os.Setenv(key, val)
	}
}

// ----------------------------
// Defaults
// ----------------------------
func defaultTaskFile() string {
	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Join(home, "todo", "tasks.json")
	}
	return "tasks.json"
}

func defaultBackupTaskFile() string {
	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Join(home, "todo", "tasks.json.bak")
	}
	return "tasks.json.bak"
}

// ----------------------------
// Env helpers
// ----------------------------
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return strings.TrimSpace(value)
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		v := strings.ToLower(strings.TrimSpace(value))
		return v == "1" || v == "true" || v == "yes"
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if i, err := strconv.Atoi(strings.TrimSpace(value)); err == nil {
			return i
		}
	}
	return fallback
}

// ----------------------------
// File path resolution
// ----------------------------
func GetTaskFilePath() string {
	return resolvePath(TaskFile)
}

func GetBackupTaskFilePath() string {
	return resolvePath(BackupTaskFile)
}

func resolvePath(file string) string {
	var path string

	switch PathMode {
	case "cwd":
		// Always relative to current working directory
		if !filepath.IsAbs(file) {
			if cwd, err := os.Getwd(); err == nil {
				path = filepath.Join(cwd, file)
			}
		} else {
			path = file
		}
	default: // "home" or anything else
		path = file
	}

	fmt.Println(">>> Resolving TaskFile:", path, " PathMode:", PathMode)

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create directory %s: %v\n", dir, err)
	}
	return path
}

// ----------------------------
// Other config accessors
// ----------------------------
func GetSortOrder() string {
	order := getEnv("TODO_SORT", "id")
	allowed := []string{"id", "due", "priority", "text"}
	for _, a := range allowed {
		if order == a {
			return order
		}
	}
	fmt.Fprintf(os.Stderr, "Invalid TODO_SORT: %s, falling back to 'id'\n", order)
	return "id"
}

func GetOutputFormat() string    { return DefaultOutputFormat }
func GetTuiTheme() string        { return TuiTheme }
func GetDefaultPriority() string { return DefaultPriority }
func GetFzfArgs() []string       { return strings.Fields(FzfArgs) }
func GetMaxTasks() int           { return MaxTasks }
