package handlers

import (
	"os"
	"strings"
	td "todo/td"
)

func Add() {
	if len(os.Args) < 3 {
		warn("Usage: todo add [task text] [optional due date] [priority: high/medium/low] [recurrence: every ...]")
		return
	}

	text := os.Args[2]
	args := os.Args[3:]

	var dueParts []string
	var priority, recurring string

	for i := 0; i < len(args); i++ {
		arg := strings.ToLower(args[i])
		switch {
		case arg == "high" || arg == "medium" || arg == "low":
			priority = arg
		case strings.HasPrefix(arg, "every"):
			recurring = strings.Join(args[i:], " ")
			i = len(args) // stop processing after recurring
		default:
			dueParts = append(dueParts, args[i])
		}
	}

	due := strings.Join(dueParts, " ")

	if err := td.AddTaskWithDueDateAndRecurring(text, due, recurring, priority); err != nil {
		fail("Error adding task: %v", err)
		return
	}

	success("✅ Task added: %s", text)
}
