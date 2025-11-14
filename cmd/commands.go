package main

import (
	"fmt"
	"os"
	"strings"

	"todo/handlers" // 👈 import your handlers package
	todo "todo/td"
)

func HandleCommands() {
	cmd := strings.ToLower(arg(1))
	if cmd == "" {
		handlers.HandleList()
		return
	}

	aliases := map[string]string{
		"a": "add", "ls": "list", "d": "done",
		"rm": "delete", "del": "delete",
		"clr": "clear", "r": "reset", "s": "search",
		"h": "help", "?": "help", "-h": "help", "--help": "help",
	}

	if real, ok := aliases[cmd]; ok {
		cmd = real
	}

	switch cmd {
	case "add":
		handlers.Add()
	case "edit":
		handlers.Edit()
	case "list":
		handlers.List()
	case "done":
		handlers.Done()
	case "due":
		handlers.Due()
	case "delete":
		handlers.Delete()
	case "clear":
		handlers.Clear()
	case "reset":
		handlers.Reset()
	case "search":
		handlers.Search()
	case "priority":
		handlers.Priority()
	case "tag":
		handlers.Tags()
	case "recurring":
		handlers.Recurring()
	case "move":
		handlers.Move()
	case "tui":
		StartTUI()
	case "note":
		handlers.Note()
	case "sub":
		handlers.Subtask()
	case "pick":
		handlers.Pick()
	case "show":
		handlers.Show()
	case "help":
		handlers.Help()
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
  todo due [id] [date]         → Change due date (supports natural lang)
  todo delete                  → Delete selected task(s)
  todo edit                    → Edit selected task
  todo note [id] [text?]       → Edit notes (or use prompt)
  todo sub [parent id] [text]  → Add subtask to task
  todo move                    → Reorder task
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
