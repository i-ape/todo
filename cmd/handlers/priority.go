package handlers

import (
	td"todo/td"
)

func Priority() {
	task, err := selectSingleTask()
	if err != nil {
		fail("Error selecting task: %v", err)
		return
	}

	newPriority := td.PromptInput("🔥 Set priority (high/medium/low)", task.Priority)
	if newPriority == "" {
		warn("No changes made.")
		return
	}

	parsed := td.ParsePriority(newPriority)

	err = td.UpdateTaskByID(task.ID, func(t *td.Task) {
		t.Priority = parsed
	})
	if err != nil {
		fail("Failed to update priority: %v", err)
		return
	}

	success("🔥 Priority updated to %s for #%d", newPriority, task.ID)
}
