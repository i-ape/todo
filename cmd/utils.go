// cmd/utils.go
package main

import (
	"fmt"
	"os"
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
