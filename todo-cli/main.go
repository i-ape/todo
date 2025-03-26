package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: todo add|list|done|delete|clear [task text or task ID]")
		return
	}

	command := os.Args[1]

	switch command {
	case "add":
		if len(os.Args) < 4 {
			fmt.Println("Usage: todo add [task text] [due date YYYY-MM-DD]")
			return
		}
		taskText := os.Args[2]
		dueDate := os.Args[3]
		_, err := time.Parse("2006-01-02", dueDate)
		if err != nil {
			fmt.Println("Invalid date format. Use YYYY-MM-DD")
			return
		}
		if err := AddTask(taskText, dueDate); err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("✅ Task added:", taskText, "(Due:", dueDate, ")")
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

	case "clear":
		if err := todo.ClearTasks(); err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("🗑️ All tasks deleted!")
		}

	default:
		fmt.Println("Invalid command. Use add|list|done|delete|clear")
	}
}