package main

import (
	"fmt"
	"os"
	"strings"

	"todo/config"
	todo "todo/todo.int"
	//"github.com/fatih/color"
)

//var DisableFzf, EnableTui bool

func HandleCommands() {
	// Parse flags
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "--no-fzf" {
			config.DisableFzf = true
			os.Args = append(os.Args[:i], os.Args[i+1:]...)
			i--
		}
		if os.Args[i] == "--tui" {
			config.EnableTui = true
			os.Args = append(os.Args[:i], os.Args[i+1:]...)
			i--
		}
	}

	// 🧃 Launch TUI if enabled
	if EnableTui {
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
	case "move":
		handleMove()
	case "help":
		printHelp()
	case "tui":
		StartTUI()

	default:
		fmt.Println("❌ Unknown command:", cmd)
		printHelp()

	}
}

// --- Task Management Functions ---

func AddTask(text, due string) error {
	return todo.AddTaskWithDueDate(text, due)
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
  --done                        → Show only completed tasks
  --pending                     → Show only incomplete tasks
  --tag=work                    → Filter by tag
  --priority=high              → Filter by priority
  --today                       → Due today
  --overdue                     → Show overdue tasks
  --json                        → Output JSON format
  --tui                         → bubble tea interface

Aliases:
  a, ls, d, rm, clr, r, s, del, h, ?, -h, --help

Shortcuts:
  today, tomorrow, next, fri, sun, eom, eow, bonm, etc.`)
}
