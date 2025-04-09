// commands.go
package main

import (
	"fmt"
	"os"

	"todo-cli/todo"
)

// AddTask adds a new task
func AddTask(text string) error {
	return todo.AddTask(text)
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
func SetDueDate(input string, dueDate string) error {
	return todo.SetDueDate(input, dueDate)
}

// DeleteTask removes a task
func DeleteTask(input string) error {
	return todo.DeleteTask(input) // Pass string instead of int
}

// HandleCommands processes CLI input
func HandleCommands() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: todo add|list|done|delete|due [task text or task ID] [due date]")
		return
	}

	cmd := os.Args[1]

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
		if err := todo.AddTaskWithDueDate(text, due); err != nil {
			fmt.Println("Error:", err)
		}
	
	case "list":
		ListTasks()

	case "done":
		if len(os.Args) < 3 {
			fmt.Println("Usage: todo done [task ID or task text]")
			return
		}
		input := os.Args[2]
		if err := MarkTaskDone(input); err != nil {
			fmt.Println("Error:", err)
		}

	case "due":
		if len(os.Args) < 4 {
			fmt.Println("Usage: todo due [task ID or task text] [YYYY-MM-DD]")
			return
		}
		input := os.Args[2]
		dueDate := os.Args[3]
		if err := SetDueDate(input, dueDate); err != nil {
			fmt.Println("Error:", err)
		}

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: todo delete [task ID or task text]")
			return
		}
		input := os.Args[2]
		if err := DeleteTask(input); err != nil {
			fmt.Println("Error:", err)
		}

	default:
		fmt.Println("Invalid command. Use add|list|done|delete|due")
	}
}
