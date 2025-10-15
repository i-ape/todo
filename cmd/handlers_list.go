// cmd/handlers_list.go
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"todo/config"
	todo "todo/td"
)

// --- List tasks ---
func handleList() {
	path := config.GetTaskFilePath()
	data, err := os.ReadFile(path)
	if err != nil {
		warn("No tasks found.")
		return
	}

	var tasks []todo.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		fail("Failed to read tasks: %v", err)
		return
	}

	if len(tasks) == 0 {
		warn("No tasks in your list.")
		return
	}

	fmt.Println("🗒️  Your tasks:")
	for _, t := range tasks {
		status := " "
		if t.Completed {
			status = "x"
		}
		fmt.Printf("[%s] #%d %s\n", status, t.ID, t.Text)
	}
}

// --- Search tasks ---
func handleSearch() {
	if len(os.Args) < 3 {
		fail("Usage: todo search [keyword]")
		return
	}
	query := strings.ToLower(strings.Join(os.Args[2:], " "))

	tasks, err := todo.LoadTasks()
	if err != nil {
		fail("Failed to load tasks: %v", err)
		return
	}

	fmt.Printf("🔍 Search results for \"%s\":\n", query)
	found := false
	for _, t := range tasks {
		if strings.Contains(strings.ToLower(t.Text), query) {
			fmt.Printf(" - #%d %s\n", t.ID, t.Text)
			found = true
		}
	}

	if !found {
		warn("No matching tasks found.")
	}
}

// --- Show task details ---
func handleShow() {
	if len(os.Args) < 3 {
		fail("Usage: todo show [id]")
		return
	}
	id := os.Args[2]
	task, err := todo.GetTaskByInput(id)
	if err != nil {
		fail("%v", err)
		return
	}
	todo.PrintTaskDetails(*task)
}
