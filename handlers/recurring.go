package handlers

import (
	"todo/cmd"
	td "todo/td"
)

func Recurring() {
	task, err := cmd.SelectSingleTask()
	if err != nil {
		fail("Error selecting task: %v", err)
		return
	}

	rec := td.PromptInput("🔁 Set recurrence (daily, weekly, etc.)", task.Recurring)
	if rec == "" {
		warn("No changes made.")
		return
	}

	err = td.UpdateTaskByID(task.ID, func(t *td.Task) {
		t.Recurring = rec
	})
	if err != nil {
		fail("Failed to set recurrence: %v", err)
		return
	}

	success("🔁 Recurrence set to %s for #%d", rec, task.ID)
}
