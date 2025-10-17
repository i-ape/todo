// cmd/utils.go
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"todo/config"
	todo "todo/td"
)

// --- Output helpers ---
func info(msg string, args ...any) {
	fmt.Printf("ℹ️  "+msg+"\n", args...)
}

func success(msg string, args ...any) {
	fmt.Printf("✅ "+msg+"\n", args...)
}

func warn(msg string, args ...any) {
	fmt.Printf("⚠️  "+msg+"\n", args...)
}

func fail(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, "❌ "+msg+"\n", args...)
}

// --- Debug logging (optional) ---
var debugEnabled = os.Getenv("TODO_DEBUG") == "true"

func debug(msg string, args ...any) {
	if debugEnabled {
		fmt.Printf("[debug] "+msg+"\n", args...)
	}
}

// --- CLI arg helper ---
func arg(i int) string {
	if len(os.Args) > i {
		return os.Args[i]
	}
	return ""
}

// --- Task selection helper ---
func selectSingleTask() (*todo.Task, error) {
	selected, err := todo.SelectTasksWithFzf(false, config.DisableFzf)
	if err != nil {
		return nil, err
	}
	if len(selected) == 0 {
		return nil, fmt.Errorf("no task selected")
	}
	return &selected[0], nil
}

// --- Task lookup by ID or partial text ---
func getTaskByInput(input string) (*todo.Task, error) {
	tasks, err := todo.LoadTasks()
	if err != nil {
		return nil, err
	}

	// Try numeric ID first
	if id, err := strconv.Atoi(input); err == nil {
		for _, t := range tasks {
			if t.ID == id {
				return &t, nil
			}
		}
		return nil, fmt.Errorf("no task found with ID %d", id)
	}

	// Try fuzzy text match
	lower := strings.ToLower(input)
	for _, t := range tasks {
		if strings.Contains(strings.ToLower(t.Text), lower) {
			return &t, nil
		}
	}
	return nil, fmt.Errorf("no task found matching text: %q", input)
}

// containsArg checks if a given CLI argument is present (e.g. --json, --edit).
func containsArg(flag string) bool {
	for _, a := range os.Args {
		if a == flag {
			return true
		}
	}
	return false
}
