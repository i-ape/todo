package handlers

import (
	"os"
	"strings"
	td "todo/td"
)

func Note() {
	if len(os.Args) < 4 {
		fail("Usage: todo note [task ID or text] [note text]")
		return
	}

	input := os.Args[2]
	note := strings.Join(os.Args[3:], " ")

	task, err := getTaskByInput(input)
	if err != nil {
		fail("Task not found: %v", err)
		return
	}

	err = td.UpdateTaskByID(task.ID, func(t *td.Task) {
		t.Notes = note
	})
	if err != nil {
		fail("Failed to update note: %v", err)
		return
	}

	success("📝 Note added to #%d", task.ID)
}
