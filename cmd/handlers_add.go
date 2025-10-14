package main

import (
	"os"
	"strings"
	todo "todo/td"
)

func handleAdd() {
	if len(os.Args) < 3 {
		warn("Usage: todo add [task text] [optional due date] [optional priority: high/medium/low] [optional recurrence: every ...]")
		return
	}

	text := os.Args[2]
	args := os.Args[3:]

	var (
		dueParts  []string
		priority  string
		recurring string
	)

	for i := 0; i < len(args); i++ {
		arg := strings.ToLower(args[i])
		switch {
		case arg == "high" || arg == "medium" || arg == "low":
			priority = arg
		case strings.HasPrefix(arg, "every"):
			recurring = strings.Join(args[i:], " ")
			break
		default:
			dueParts = append(dueParts, args[i])
		}
	}

	due := strings.Join(dueParts, " ")

	if err := todo.AddTaskWithDueDateAndRecurring(text, due, recurring, priority); err != nil {
		fail("Error adding task: %v", err)
		return
	}
	success("Task added.")
}
