package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	todo "todo/todo.int"

	"github.com/fatih/color"
)

// --- FZF Selector ---

func selectTasksWithFzf(multi bool) ([]todo.Task, error) {
	tasks, err := todo.LoadTasks()
	if err != nil {
		return nil, fmt.Errorf("failed to load tasks: %w", err)
	}

	if !disableFzf {
		if _, err := exec.LookPath("fzf"); err == nil {
			if multi {
				return todo.SelectMultipleTasksFzf(tasks)
			}
			task, err := todo.SelectTaskFzf(tasks)
			if err != nil {
				return nil, err
			}
			return []todo.Task{task}, nil
		}
	}

	fmt.Println("FZF disabled or not found. Manual selection:")
	for _, t := range tasks {
		fmt.Printf("%d: %s\n", t.ID, t.Text)
	}
	fmt.Print("> Enter task ID: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	id, err := strconv.Atoi(input)
	if err != nil {
		return nil, fmt.Errorf("invalid ID")
	}

	for _, t := range tasks {
		if t.ID == id {
			return []todo.Task{t}, nil
		}
	}
	return nil, fmt.Errorf("task not found")
}

// --- Task Management Functions ---

var disableFzf, enableTui bool

func init() {
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		if arg == "--no-fzf" {
			disableFzf = true
			os.Args = append(os.Args[:i], os.Args[i+1:]...)
			i--
		}
		if arg == "--tui" {
			enableTui = true
			os.Args = append(os.Args[:i], os.Args[i+1:]...)
			i--
		}
	}
} // TODO: launch tea.Program(model)

func AddTask(text, due string) error {
	return todo.AddTaskWithDueDate(text, due)
}
func handleList() {
	args := os.Args[2:]
	opts := todo.ListFilterOptions{}

	for _, arg := range args {
		switch {
		case arg == "--json":
			opts.JSONOutput = true
		case arg == "--done":
			opts.ShowDone = true
		case arg == "--pending":
			opts.ShowPending = true
		case arg == "--today":
			opts.TodayOnly = true
		case arg == "--overdue":
			opts.OverdueOnly = true
		case strings.HasPrefix(arg, "--tag="):
			opts.Tags = strings.TrimPrefix(arg, "--tag=")
		case strings.HasPrefix(arg, "--priority="):
			opts.Priority = strings.TrimPrefix(arg, "--priority=")
		}
	}

	tasks, err := todo.LoadTasks()
	if err != nil {
		fmt.Println("❌ Failed to load tasks:", err)
		return
	}

	filtered := todo.FilterAndSortTasks(tasks, opts)

	if opts.JSONOutput {
		data, _ := json.MarshalIndent(filtered, "", "  ")
		fmt.Println(string(data))
		return
	}

	for _, task := range filtered {
		label := fmt.Sprintf("%d: %s", task.ID, task.Text)
		if task.DueDate != "" {
			label += fmt.Sprintf(" (Due: %s)", task.DueDate)
		}
		if len(task.Tags) > 0 {
			label += " [" + strings.Join(task.Tags, ", ") + "]"
		}
		switch {
		case task.Completed:
			fmt.Println(color.GreenString("[✓] " + label))
		case todo.IsOverdue(task.DueDate):
			fmt.Println(color.RedString("[✗] " + label))
		default:
			fmt.Println(color.CyanString("[ ] " + label))
		}
	}
}

func MarkTaskDone(input string) error {
	return todo.MarkTaskDone(input)
}

func SetDueDate(input, dueDate string) error {
	return todo.SetDueDate(input, dueDate)
}

func DeleteTask(input string) error {
	return todo.DeleteTask(input)
}

func ClearTasks() error {
	return todo.ClearTasks()
}

func ResetTasks() error {
	return os.Remove("todo/tasks.json")
}

func SearchTasks(keyword string) {
	todo.SearchTasks(keyword)
}

// --- CLI Command Dispatcher ---

func HandleCommands() {
	// Parse flags
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "--no-fzf" {
			disableFzf = true
			os.Args = append(os.Args[:i], os.Args[i+1:]...)
			i--
		}
		if os.Args[i] == "--tui" {
			enableTui = true
			os.Args = append(os.Args[:i], os.Args[i+1:]...)
			i--
		}
	}

	// 🧃 Launch TUI if enabled
	if enableTui {
		StartTUI()
		return
	}

	if len(os.Args) < 2 {
		printHelp()
		return
	}

	aliases := map[string]string{
		"a":      "add",
		"ls":     "list",
		"d":      "done",
		"rm":     "delete",
		"del":    "delete",
		"clr":    "clear",
		"r":      "reset",
		"s":      "search",
		"h":      "help",
		"?":      "help",
		"-h":     "help",
		"--help": "help",
	}

	cmd := strings.ToLower(os.Args[1])
	if real, ok := aliases[cmd]; ok {
		cmd = real
	}

	switch cmd {
	case "add":
		handleAdd()
	case "edit":
		handleEdit()
	case "list":
		handleList()
	case "done":
		handleDone()
	case "due":
		handleDue()
	case "delete":
		handleDelete()
	case "clear":
		handleClear()
	case "reset":
		handleReset()
	case "search":
		handleSearch()
	case "priority":
		handlePriority()
	case "tag":
		handleTags()
	case "recurring":
		handleRecurring()

	case "help":
		printHelp()
	case "tui":
		StartTUI()

	default:
		fmt.Println("❌ Unknown command:", cmd)
		printHelp()
	}
}

// --- Handlers ---

func handleAdd() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: todo add [task text] [optional due date]")
		return
	}
	text := os.Args[2]
	due := ""
	if len(os.Args) > 3 {
		due = strings.Join(os.Args[3:], " ")
	}
	if err := AddTask(text, due); err != nil {
		fmt.Println("Error:", err)
	}
}

func handleEdit() {
	selected, err := selectTasksWithFzf(false)
	if err != nil || len(selected) == 0 {
		fmt.Println("Select error:", err)
		return
	}
	task := selected[0]

	newText := todo.PromptInput(fmt.Sprintf("✏️  Edit task \"%s\"", task.Text), task.Text)

	if newText == "" {
		fmt.Println("No changes made.")
		return
	}
	if err := todo.EditTaskText(strconv.Itoa(task.ID), newText); err != nil {
		fmt.Println("Edit error:", err)
	}
}

func handleDone() {
	selected, err := selectTasksWithFzf(true)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for _, task := range selected {
		if err := MarkTaskDone(strconv.Itoa(task.ID)); err != nil {
			fmt.Println("❌", err)
		}
	}
}

func handleDelete() {
	selected, err := selectTasksWithFzf(true)
	if err != nil {
		fmt.Println("Error selecting task:", err)
		return
	}
	for _, task := range selected {
		if err := todo.DeleteTask(strconv.Itoa(task.ID)); err != nil {
			fmt.Println("❌", err)
		}
	}
}

func handleDue() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: todo due [task ID or task text] [date]")
		return
	}
	input := os.Args[2]
	dueDate := strings.Join(os.Args[3:], " ")
	if err := SetDueDate(input, dueDate); err != nil {
		fmt.Println("Error:", err)
	}
}

func handleSearch() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: todo search [keyword]")
		return
	}
	SearchTasks(os.Args[2])
}
func handlePriority() {
	selected, err := selectTasksWithFzf(false)
	if err != nil || len(selected) == 0 {
		fmt.Println("Error selecting task:", err)
		return
	}
	task := selected[0]

	newPriority := todo.PromptInput("🔥 Set priority (high/medium/low)", task.Priority)
	if newPriority == "" {
		fmt.Println("No changes made.")
		return
	}
	task.Priority = todo.ParsePriority(newPriority)

	tasks, _ := todo.LoadTasks()
	for i := range tasks {
		if tasks[i].ID == task.ID {
			tasks[i] = task
			break
		}
	}
	if err := todo.SaveTasks(tasks); err != nil {
		fmt.Println("Failed to save:", err)
		return
	}
	fmt.Println("✅ Priority updated.")
}

func handleTags() {
	selected, err := selectTasksWithFzf(false)
	if err != nil || len(selected) == 0 {
		fmt.Println("Error selecting task:", err)
		return
	}
	task := selected[0]

	existing := strings.Join(task.Tags, ", ")
	raw := todo.PromptInput("🏷️  Edit tags (comma-separated)", existing)
	if raw == "" {
		fmt.Println("No changes made.")
		return
	}
	task.Tags = todo.ParseTags(raw)

	tasks, _ := todo.LoadTasks()
	for i := range tasks {
		if tasks[i].ID == task.ID {
			tasks[i] = task
			break
		}
	}
	if err := todo.SaveTasks(tasks); err != nil {
		fmt.Println("Failed to save:", err)
		return
	}
	fmt.Println("✅ Tags updated.")
}

func handleClear() {
	if err := ClearTasks(); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("✅ All tasks cleared.")
	}
}

func handleReset() {
	if err := ResetTasks(); err != nil {
		fmt.Println("⚠️ Reset failed:", err)
	} else {
		fmt.Println("🗑️ tasks.json deleted.")
	}
}

// --- Help ---

func printHelp() {
	fmt.Println(`
📝 Usage:
  todo add [text] [due?]       → Add new task
  todo list                    → List all tasks
  todo done                    → Mark one or more tasks done
  todo due [id|text] [date]    → Set/change due date
  todo delete                  → Delete one or more tasks
  todo edit                    → Edit a task
  todo search [keyword]        → Search task text
  todo tag                     → Edit task tags
  todo clear                   → Clear all tasks
  todo reset                   → Delete tasks.json
  todo help                    → Show help

💡 Flags:
  --no-fzf                      → Disable FZF interactive mode
  --done						→ Show only completed tasks
  --pending						→ Show only incomplete tasks
  --tag=work					→ Filter by tag
  --priority=high				→ Filter by priority
  --today						→ Due today
  --overdue						→ Show overdue tasks
  --json 						→ Output JSON format
  --tui 						→ bubble tea interface


🔤 Aliases:
  a, ls, d, rm, clr, r, s, del, h, ?, -h, --help`)
}
