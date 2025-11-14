package handlers

import (
	"errors"
	"os"
	td"todo/td"
)

func Show() {
	if len(os.Args) < 3 {
		fail("Usage: todo show [id]")
		return
	}

	id, err := parseID(os.Args[2])
	if err != nil {
		fail("Invalid ID: %v", err)
		return
	}

	task, err := td.GetTaskByID(id)
	if err != nil {
		if errors.Is(err, td.ErrTaskNotFound) {
			warn("Task #%d not found.", id)
			return
		}
		fail("Failed to load task: %v", err)
		return
	}

	td.PrintTaskDetails(*task)
}
