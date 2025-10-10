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
// It holds all user settings from defaults, YAML, and environment.
var Cfg Config

// Config defines all available configuration options for the todo CLI.
type Config struct {
	Path            string `yaml:"path"`             // path to tasks.json
	BackupPath      string `yaml:"backup_path"`      // path to backup file
	PathMode        string `yaml:"path_mode"`        // "home" or "cwd"
	Output          string `yaml:"output"`           // output format: text/json
	SortOrder       string `yaml:"sort_order"`       // how tasks are sorted
	TuiTheme        string `yaml:"tui_theme"`        // TUI theme: dark/light
	DefaultPriority string `yaml:"default_priority"` // fallback priority
	MaxTasks        int    `yaml:"max_tasks"`        // max tasks to keep

	// Optional / advanced fields
	EnableTui      bool `yaml:"enable_tui"`       // auto-launch TUI on startup
	DisableFzf     bool `yaml:"disable_fzf"`      // disable fuzzy finder
	AutoBackup     bool `yaml:"auto_backup"`      // automatically back up
	ShowCompleted  bool `yaml:"show_completed"`   // show completed tasks
	DefaultDueDays int  `yaml:"default_due_days"` // auto-due date offset (days)
}

// ----------------------------
// Init: load defaults → YAML → .env
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
	// 1. Load .env first (in case it defines config file paths)
	loadDotEnv(".env")

	// 2. Apply defaults
	Cfg = Config{
		Path:            defaultTaskFile(),
		BackupPath:      defaultBackupTaskFile(),
		PathMode:        "home",
		Output:          "text",
		SortOrder:       "id",
		TuiTheme:        "dark",
		DefaultPriority: "medium",
		MaxTasks:        1000,
		EnableTui:       false,
		DisableFzf:      false,
		AutoBackup:      true,
		ShowCompleted:   false,
		DefaultDueDays:  7,
	}

	// 3. Apply overrides from YAML (if present)
	loadYAMLConfig()

	// 4. Apply overrides from environment variables (highest precedence)
	applyEnvOverrides()

	// 5. Validate and print final config
	validateConfig()
	fmt.Printf("[config] Loaded configuration: %+v\n", Cfg)
}

// ----------------------------
// YAML loader
// ----------------------------
func loadYAMLConfig() {
	cfgPath := filepath.Join(os.Getenv("HOME"), ".todo", "config.yaml")
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return // silently skip if no YAML file
	}

	var parsed Config
	if err := yaml.Unmarshal(data, &parsed); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse %s: %v\n", cfgPath, err)
		return
	}

	// Merge YAML values into global Cfg (only non-zero / non-empty)
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
	if parsed.EnableTui {
		Cfg.EnableTui = parsed.EnableTui
	}
	if parsed.DisableFzf {
		Cfg.DisableFzf = parsed.DisableFzf
	}
	if parsed.AutoBackup {
		Cfg.AutoBackup = parsed.AutoBackup
	}
	if parsed.ShowCompleted {
		Cfg.ShowCompleted = parsed.ShowCompleted
	}
	if parsed.DefaultDueDays != 0 {
		Cfg.DefaultDueDays = parsed.DefaultDueDays
	}
}

// ----------------------------
// .env loader
// ----------------------------
func loadDotEnv(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return // skip silently if no .env
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
		val := strings.Trim(strings.TrimSpace(parts[1]), `"`)
		os.Setenv(key, val)
	}
}

// ----------------------------
// Environment overrides
// ----------------------------
func applyEnvOverrides() {
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
		if n, err := strconv.Atoi(v); err == nil {
			Cfg.MaxTasks = n
		}
	}
	if v := os.Getenv("TODO_ENABLE_TUI"); v != "" {
		Cfg.EnableTui = v == "true"
	}
	if v := os.Getenv("TODO_DISABLE_FZF"); v != "" {
		Cfg.DisableFzf = v == "true"
	}
	if v := os.Getenv("TODO_AUTO_BACKUP"); v != "" {
		Cfg.AutoBackup = v == "true"
	}
	if v := os.Getenv("TODO_SHOW_COMPLETED"); v != "" {
		Cfg.ShowCompleted = v == "true"
	}
	if v := os.Getenv("TODO_DEFAULT_DUE_DAYS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			Cfg.DefaultDueDays = n
		}
	}
}

// ----------------------------
// Validation
// ----------------------------
func validateConfig() {
	validSorts := map[string]bool{"id": true, "due": true, "priority": true, "text": true}
	if !validSorts[Cfg.SortOrder] {
		fmt.Fprintf(os.Stderr, "Warning: invalid sort_order %q, using 'id'\n", Cfg.SortOrder)
		Cfg.SortOrder = "id"
	}
	if Cfg.MaxTasks <= 0 {
		Cfg.MaxTasks = 1000
	}
}

// ----------------------------
// Default paths
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
func GetSortOrder() string        { return Cfg.SortOrder }
func GetOutputFormat() string     { return Cfg.Output }
func GetTuiTheme() string         { return Cfg.TuiTheme }
func GetDefaultPriority() string  { return Cfg.DefaultPriority }
func GetMaxTasks() int            { return Cfg.MaxTasks }
func IsTuiEnabled() bool          { return Cfg.EnableTui }
func IsFzfDisabled() bool         { return Cfg.DisableFzf }
func ShouldAutoBackup() bool      { return Cfg.AutoBackup }
func ShowCompletedTasks() bool    { return Cfg.ShowCompleted }
func GetDefaultDueDays() int      { return Cfg.DefaultDueDays }
