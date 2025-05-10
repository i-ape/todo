package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	todo "todo/todo.int"
)

// --- Task Management Functions ---

func AddTask(text, due string) error {
	return todo.AddTaskWithDueDate(text, due)
}

func ListTasks() {
	todo.ListTasks()
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
		ListTasks()

	case "done":
		handleDone()

	case "due":
		handleDue()

	case "delete":
		handleDelete()

	case "clear":
		if err := ClearTasks(); err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("✅ All tasks cleared.")
		}

	case "reset":
		if err := ResetTasks(); err != nil {
			fmt.Println("⚠️ Reset failed:", err)
		} else {
			fmt.Println("🗑️ tasks.json deleted.")
		}

	case "search":
		handleSearch()

	case "help":
		printHelp()

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
	tasks, err := todo.LoadTasks()
	if err != nil {
		fmt.Println("Failed to load tasks:", err)
		return
	}
	selected, err := todo.SelectTaskFzf(tasks)
	if err != nil {
		fmt.Println("Select error:", err)
		return
	}

	fmt.Printf("✏️  Editing: %s\n> ", selected.Text)
	reader := bufio.NewReader(os.Stdin)
	newText, _ := reader.ReadString('\n')
	newText = strings.TrimSpace(newText)
	if newText == "" {
		fmt.Println("No changes made.")
		return
	}
	if err := todo.EditTaskText(strconv.Itoa(selected.ID), newText); err != nil {
		fmt.Println("Edit error:", err)
	}
}

func handleDone() {
	tasks, _ := todo.LoadTasks()
	selected, err := todo.SelectTaskFzf(tasks)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if err := MarkTaskDone(strconv.Itoa(selected.ID)); err != nil {
		fmt.Println("Error:", err)
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

func handleDelete() {
	tasks, err := todo.LoadTasks()
	if err != nil {
		fmt.Println("Error loading tasks:", err)
		return
	}
	selected, err := todo.SelectTaskFzf(tasks)
	if err != nil {
		fmt.Println("Error selecting task:", err)
		return
	}
	if err := todo.DeleteTask(strconv.Itoa(selected.ID)); err != nil {
		fmt.Println("Error deleting task:", err)
	}
}

func handleSearch() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: todo search [keyword]")
		return
	}
	SearchTasks(os.Args[2])
}

// --- Help Menu ---

func printHelp() {
	fmt.Println(`📝 Usage:
  todo add [text] [due?]       → Add new task
  todo list                    → List all tasks
  todo done [id|text]          → Mark task done
  todo due [id|text] [date]    → Set/change due date
  todo delete [id|text]        → Delete task
  todo search [keyword]        → Search task text
  todo clear                   → Clear all tasks
  todo reset                   → Delete tasks.json
  todo help                    → Show help

🔤 Aliases:
  a     → add
  ls    → list
  d     → done
  del   → delete
  rm    → delete
  clr   → clear
  r     → reset
  s     → search
  h, ?, -h, --help → help`)
}
