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
	debug("Reading tasks from %s", path)

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

	info("Found %d tasks:", len(tasks))
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
	query := strings.ToLower(strings.Join(os.Args[2:], " "))
	if query == "" {
		fail("Usage: todo search [keyword]")
		return
	}
	debug("Searching tasks for keyword: %q", query)

	tasks, err := todo.LoadTasks()
	if err != nil {
		fail("Failed to load tasks: %v", err)
		return
	}

	info("Searching %d tasks...", len(tasks))
	fmt.Printf("🔍 Results for \"%s\":\n", query)
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
	id := arg(2)
	if id == "" {
		fail("Usage: todo show [id]")
		return
	}

	debug("Showing details for task input: %s", id)

	task, err := getTaskByInput(id)
	if err != nil {
		fail("%v", err)
		return
	}

	info("Displaying details for #%d", task.ID)
	todo.PrintTaskDetails(*task)
}
