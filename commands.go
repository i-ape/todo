// commands.go
package main

import (
	"fmt"
	"os"
	"strconv"

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
func MarkTaskDone(id int) error {
	return todo.MarkTaskDone(id)
}

// DeleteTask removes a task
func DeleteTask(input string) error {
	return todo.DeleteTask(input) // Pass string instead of int
}

// HandleCommands processes CLI input
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
		id, _ := strconv.Atoi(os.Args[2]) // Fix: Ensure it's an integer
		if err := MarkTaskDone(id); err != nil {
			fmt.Println("Error:", err)
		}

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: todo delete [task ID or task text]")
			return
		}
	
		input := os.Args[2] // Keep input as string
		if err := DeleteTask(input); err != nil {
			fmt.Println("Error:", err)
		}
	
	
	default:
		fmt.Println("Invalid command. Use add|list|done|delete")
	}
}
