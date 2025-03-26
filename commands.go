// commands.go
package main

import (
	"fmt"
	"os"
	"strconv"

	"todo-cli/todo"  // ✅ Ensure this matches your `go.mod`
)


// AddTask adds a new task (calls `todo.AddTask`)
func AddTask(text string) error {
	return todo.AddTask(text)  // ✅ Call function from `todo`
}

// ListTasks displays all tasks (calls `todo.ListTasks`)
func ListTasks() {
	todo.ListTasks()  // ✅ Call function from `todo`
}

// MarkTaskDone marks a task as completed (calls `todo.MarkTaskDone`)
func MarkTaskDone(id int) error {
	return todo.MarkTaskDone(id)  // ✅ Call function from `todo`
}

// DeleteTask removes a task (calls `todo.DeleteTask`)
func DeleteTask(id int) error {
	return todo.DeleteTask(id)  // ✅ Call function from `todo`
}

// Handle CLI commands
func HandleCommands() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: todo add|list|done|delete [task text or task ID]")
		return
	}

	cmd := os.Args[1]

	switch cmd {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: todo add [task text]")
			return
		}
		text := os.Args[2]
		if err := AddTask(text); err != nil {
			fmt.Println("Error:", err)
		}

	case "list":
		ListTasks()

	case "done":
		if len(os.Args) < 3 {
			fmt.Println("Usage: todo done [task ID]")
			return
		}
		id, _ := strconv.Atoi(os.Args[2])
		if err := MarkTaskDone(id); err != nil {
			fmt.Println("Error:", err)
		}

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: todo delete [task ID]")
			return
		}
		id, _ := strconv.Atoi(os.Args[2])
		if err := DeleteTask(id); err != nil {
			fmt.Println("Error:", err)
		}

	default:
		fmt.Println("Invalid command. Use add|list|done|delete")
	}
}