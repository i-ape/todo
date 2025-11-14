package handlers

import (
	"errors"
	"os"
	td"todo/td"
)

func Edit() {
	if len(os.Args) < 4 {
		fail("Usage: todo edit [id] [new text]")
		return
	}

	id, err := parseID(os.Args[2])
	if err != nil {
		fail("Invalid ID: %v", err)
		return
	}

	newText := os.Args[3]
	err = td.UpdateTaskByID(id, func(t *td.Task) { t.Text = newText })
	if err != nil {
		if errors.Is(err, td.ErrTaskNotFound) {
			warn("Task #%d not found.", id)
			return
		}
		fail("Failed to update: %v", err)
		return
	}
	success("✏️ Updated task #%d", id)
}
