package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"todo/config"
	todo "todo/td"
)

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

func handleShow() {
	if len(os.Args) < 3 {
		fail("Usage: todo show [id]")
		return
	}

	id, err := parseID(os.Args[2])
	if err != nil {
		fail("Invalid ID: %v", err)
		return
	}

	task, err := todo.GetTaskByID(id)
	if err != nil {
		if errors.Is(err, todo.ErrTaskNotFound) {
			warn("Task #%d not found.", id)
			return
		}
		fail("Failed to get task: %v", err)
		return
	}

	todo.PrintTaskDetails(*task)
}
