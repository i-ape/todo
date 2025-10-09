// config/config.go
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// Cfg is the global configuration object.
// All config values should come from here.
var Cfg Config

// Config holds all user-adjustable settings.
// Defaults are applied first, then overridden by config.yaml,
// and finally overridden by environment variables.
type Config struct {
	Path            string `yaml:"path"`             // path to tasks.json
	BackupPath      string `yaml:"backup_path"`      // path to backup file
	PathMode        string `yaml:"path_mode"`        // "home" or "cwd"
	Output          string `yaml:"output"`           // output format: text/json
	SortOrder       string `yaml:"sort_order"`       // how tasks are sorted
	TuiTheme        string `yaml:"tui_theme"`        // TUI theme: dark/light
	DefaultPriority string `yaml:"default_priority"` // fallback priority
	MaxTasks        int    `yaml:"max_tasks"`        // max tasks to keep

	// Optional / future fields
	EnableTui      bool `yaml:"enable_tui"`       // auto-launch TUI on startup
	DisableFzf     bool `yaml:"disable_fzf"`      // disable fuzzy finder
	AutoBackup     bool `yaml:"auto_backup"`      // automatically back up
	ShowCompleted  bool `yaml:"show_completed"`   // list completed tasks
	DefaultDueDays int  `yaml:"default_due_days"` // auto-due date offset
}

// ----------------------------
// Init: load defaults, then YAML, then env
// ----------------------------

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
	// Load .env first, in case it sets config paths
	loadDotEnv(".env")

	// Start with defaults
	Cfg = Config{
		Path:            defaultTaskFile(),
		BackupPath:      defaultBackupTaskFile(),
		PathMode:        "home",
		Output:          "text",
		SortOrder:       "id",
		TuiTheme:        "dark",
		DefaultPriority: "medium",
		MaxTasks:        1000,
	}

	// Apply config.yaml overrides
	loadYAMLConfig()

	// Apply env overrides (highest precedence)
	if v := os.Getenv("TODO_PATH"); v != "" {
		Cfg.Path = v
	}
	if v := os.Getenv("TODO_BACKUP_PATH"); v != "" {
		Cfg.BackupPath = v
	}
	if v := os.Getenv("TODO_PATH_MODE"); v != "" {
		Cfg.PathMode = v
	}
	if v := os.Getenv("TODO_OUTPUT"); v != "" {
		Cfg.Output = v
	}
	if v := os.Getenv("TODO_SORT"); v != "" {
		Cfg.SortOrder = v
	}
	if v := os.Getenv("TODO_TUI_THEME"); v != "" {
		Cfg.TuiTheme = v
	}
	if v := os.Getenv("TODO_DEFAULT_PRIORITY"); v != "" {
		Cfg.DefaultPriority = v
	}
	if v := os.Getenv("TODO_MAX_TASKS"); v != "" {
		// simple atoi parse
		if n, err := strconv.Atoi(v); err == nil {
			Cfg.MaxTasks = n
		}
	}
}

// ----------------------------
// YAML loader
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
	if parsed.SortOrder != "" {
		Cfg.SortOrder = parsed.SortOrder
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
// Accessors
// ----------------------------
func GetSortOrder() string {
	allowed := []string{"id", "due", "priority", "text"}
	for _, a := range allowed {
		if Cfg.SortOrder == a {
			return a
		}
	}
	fmt.Fprintf(os.Stderr, "Invalid sort order: %s, falling back to 'id'\n", Cfg.SortOrder)
	return "id"
}

func GetOutputFormat() string    { return Cfg.Output }
func GetTuiTheme() string        { return Cfg.TuiTheme }
func GetDefaultPriority() string { return Cfg.DefaultPriority }
func GetMaxTasks() int           { return Cfg.MaxTasks }
func IsTuiEnabled() bool         { return Cfg.EnableTui }
func IsFzfDisabled() bool        { return Cfg.DisableFzf }
func ShouldAutoBackup() bool     { return Cfg.AutoBackup }
func ShowCompletedTasks() bool   { return Cfg.ShowCompleted }
func DefaultDueDays() int        { return Cfg.DefaultDueDays }
