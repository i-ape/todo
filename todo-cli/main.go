package main

import (
	"fmt"
	"os"
	"todo-cli/todo"  // ✅ Import `todo` package
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
		text := os.Args[2]
		if err := todo.AddTask(text); err != nil {  // ✅ Call function from `todo`
			fmt.Println("Error:", err)
		}

	case "list":
		todo.ListTasks()  // ✅ Call function from `todo`

	case "done":
		if len(os.Args) < 3 {
			fmt.Println("Usage: todo done [task ID]")
			return
		}
		id := os.Args[2]
		if err := todo.MarkTaskDone(id); err != nil {  // ✅ Call function from `todo`
			fmt.Println("Error:", err)
		}

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: todo delete [task ID]")
			return
		}
		id := os.Args[2]
		if err := todo.DeleteTask(id); err != nil {  // ✅ Call function from `todo`
			fmt.Println("Error:", err)
		}

	default:
		fmt.Println("Invalid command. Use add|list|done|delete")
	}
}
