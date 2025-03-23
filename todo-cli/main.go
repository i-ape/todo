package main

import (
	"fmt"
	"os"
	"strconv"

	"todo-cli/todo"  // ✅ Ensure this matches `go.mod`
)
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: todo add|list|done|delete [task text or task ID]")
		return
	}

	command := os.Args[1]

	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: todo add [task text]")
			return
		}
		taskText := os.Args[2]
		if err := todo.AddTask(taskText); err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("✅ Task added:", taskText)
		}

	case "list":
		todo.ListTasks()

	case "done":
		if len(os.Args) < 3 {
			fmt.Println("Usage: todo done [task ID]")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid task ID:", os.Args[2])
			return
		}
		if err := todo.MarkTaskDone(id); err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Printf("✅ Task %d marked as done!\n", id)
		}

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: todo delete [task ID]")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid task ID:", os.Args[2])
			return
		}
		if err := todo.DeleteTask(id); err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Printf("❌ Task %d deleted!\n", id)
		}

	default:
		fmt.Println("Invalid command. Use add|list|done|delete")
	}
}
