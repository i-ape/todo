package handlers

import (
	"os"
	"strings"
	td "todo/td"
)

func Due() {
	if len(os.Args) < 4 {
		fail("Usage: todo due [task ID or text] [date]")
		return
	}

	input := os.Args[2]
	due := strings.Join(os.Args[3:], " ")

	task, err := getTaskByInput(input)
	if err != nil {
		fail("Task not found: %v", err)
		return
	}

	err = td.UpdateTaskByID(task.ID, func(t *td.Task) {
		t.Due = due
	})
	if err != nil {
		fail("Failed to update due date: %v", err)
		return
	}

	success("📅 Updated due date for #%d → %s", task.ID, due)
}
