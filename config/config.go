// config/config.go
package config

import (
	"fmt"
	"os"
	"path/filepath"

	//"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

var Cfg Config // global config struct

type Config struct {
	Path            string `yaml:"path"`
	BackupPath      string `yaml:"backup_path"`
	PathMode        string `yaml:"path_mode"`
	Output          string `yaml:"output"`
	TuiTheme        string `yaml:"tui_theme"`
	DefaultPriority string `yaml:"default_priority"`
	MaxTasks        int    `yaml:"max_tasks"`
}

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

	// Start with defaults
	Cfg = Config{
		Path:            defaultTaskFile(),
		BackupPath:      defaultBackupTaskFile(),
		PathMode:        "home",
		Output:          "text",
		TuiTheme:        "dark",
		DefaultPriority: "medium",
		MaxTasks:        1000,
	}

	// Load YAML overrides
	loadYAMLConfig()

	// Finally, apply env overrides
	if v := os.Getenv("TODO_PATH"); v != "" {
		Cfg.Path = v
	}
	if v := os.Getenv("TODO_BACKUP_PATH"); v != "" {
		Cfg.BackupPath = v
	}
	if v := os.Getenv("TODO_PATH_MODE"); v != "" {
		Cfg.PathMode = v
	}
}

// ----------------------------
// yaml loader
// ----------------------------
func loadYAMLConfig() {
	cfgPath := filepath.Join(os.Getenv("HOME"), ".todo", "config.yaml")
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return // silently skip if no file
	}

	var parsed Config
	if err := yaml.Unmarshal(data, &parsed); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse %s: %v\n", cfgPath, err)
		return
	}

	// Merge parsed values into global Cfg
	if parsed.Path != "" {
		Cfg.Path = parsed.Path
	}
	if parsed.BackupPath != "" {
		Cfg.BackupPath = parsed.BackupPath
	}
	if parsed.PathMode != "" {
		Cfg.PathMode = parsed.PathMode
	}
	if parsed.Output != "" {
		Cfg.Output = parsed.Output
	}
	if parsed.TuiTheme != "" {
		Cfg.TuiTheme = parsed.TuiTheme
	}
	if parsed.DefaultPriority != "" {
		Cfg.DefaultPriority = parsed.DefaultPriority
	}
	if parsed.MaxTasks != 0 {
		Cfg.MaxTasks = parsed.MaxTasks
	}
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
// File path resolution
// ----------------------------
func GetTaskFilePath() string {
	return resolvePath(Cfg.Path, Cfg.PathMode)
}

func GetBackupTaskFilePath() string {
	return resolvePath(Cfg.BackupPath, Cfg.PathMode)
}

func resolvePath(file, mode string) string {
	var path string

	switch mode {
	case "cwd":
		if !filepath.IsAbs(file) {
			if cwd, err := os.Getwd(); err == nil {
				path = filepath.Join(cwd, file)
			}
		} else {
			path = file
		}
	default:
		path = file
	}

	dir := filepath.Dir(path)
	_ = os.MkdirAll(dir, 0755)
	return path
}

// ----------------------------
// Other config accessors
// ----------------------------
func GetSortOrder() string {
	order := Cfg.Output // or add `SortOrder string` to Config struct
	allowed := []string{"id", "due", "priority", "text"}
	for _, a := range allowed {
		if order == a {
			return order
		}
	}
	fmt.Fprintf(os.Stderr, "Invalid sort order: %s, falling back to 'id'\n", order)
	return "id"
}

func GetOutputFormat() string    { return DefaultOutputFormat }
func GetTuiTheme() string        { return TuiTheme }
func GetDefaultPriority() string { return DefaultPriority }
func GetFzfArgs() []string       { return strings.Fields(FzfArgs) }
func GetMaxTasks() int           { return MaxTasks }
