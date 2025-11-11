package handlers

import (
	"fmt"
	"os"
	"strings"
	todo "todo/td"
)

func Add() {
	if len(os.Args) < 3 {
		fmt.Println("❌ Usage: todo add [task text] [optional due date]")
		return
	}

	text := os.Args[2]
	due := ""
	if len(os.Args) > 3 {
		due = strings.Join(os.Args[3:], " ")
	}

	if err := todo.AddTaskWithDueDate(text, due); err != nil {
		fmt.Println("❌ Failed to add task:", err)
		return
	}
	fmt.Printf("✅ Added: %s\n", text)
}
