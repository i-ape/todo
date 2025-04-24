package main

import (
	"fmt"
	"os"
	"strings"

	"todo/todo"
)

// AddTask adds a new task
func AddTask(text string, due string) error {
	return todo.AddTaskWithDueDate(text, due)
}

// ListTasks displays all tasks
func ListTasks() {
	todo.ListTasks()
}

// MarkTaskDone marks a task as completed
func MarkTaskDone(input string) error {
	return todo.MarkTaskDone(input)
}

// SetDueDate assigns a due date to a task
// Supports formats like "YYYY-MM-DD", "DD-MM-YYYY", "tomorrow", "next week", "in 3 days"
func SetDueDate(input string, dueDate string) error {
	return todo.SetDueDate(input, dueDate)
}

// DeleteTask removes a task
func DeleteTask(input string) error {
	return todo.DeleteTask(input)
}

func ClearTasks() {
	todo.ClearTasks()
	fmt.Println("✅ All tasks cleared.")
}

func ResetTasks() {
	err := os.Remove("todo/tasks.json")
	if err != nil {
		fmt.Println("⚠️ Reset failed:", err)
	} else {
		fmt.Println("🗑️  tasks.json deleted.")
	}
}

func SearchTasks(keyword string) {
	todo.SearchTasks(keyword)
}

// HandleCommands processes CLI input
func HandleCommands() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	cmd := strings.ToLower(os.Args[1])

	switch cmd {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: todo add [task text] [optional due date]")
			return
		}
		text := os.Args[2]
		due := ""
		if len(os.Args) > 3 {
			due = os.Args[3]
		}
		if err := AddTask(text, due); err != nil {
			fmt.Println("Error:", err)
		}
	

	case "list":
		ListTasks()

	case "done":
		if len(os.Args) < 3 {
			fmt.Println("Usage: todo done [task ID or task text]")
			return
		}
		if err := MarkTaskDone(os.Args[2]); err != nil {
			fmt.Println("Error:", err)
		}

	case "due":
		if len(os.Args) < 4 {
			fmt.Println("Usage: todo due [task ID or task text] [date string]")
			return
		}
		input := os.Args[2]
		dueDate := strings.Join(os.Args[3:], " ") // ✅ Fix multi-word due date
		if err := SetDueDate(input, dueDate); err != nil {
			fmt.Println("Error:", err)
		}
	

	case "delete", "rm":
		if len(os.Args) < 3 {
			fmt.Println("Usage: todo delete [task ID or task text]")
			return
		}
		if err := DeleteTask(os.Args[2]); err != nil {
			fmt.Println("Error:", err)
		}

	case "clear":
		ClearTasks()

	case "reset":
		ResetTasks()

	case "search":
		if len(os.Args) < 3 {
			fmt.Println("Usage: todo search [keyword]")
			return
		}
		SearchTasks(os.Args[2])

	case "help":
		printHelp()

	default:
		fmt.Println("❌ Invalid command.")
		printHelp()
	}
}

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
  todo help                    → Show help`)
}
