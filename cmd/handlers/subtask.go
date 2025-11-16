package handlers

import (
	td "todo/td"
	cmd"todo/cmd"
)

func Subtask() {
	parent, err := cmd.SelectSingleTask()
	if err != nil {
		fail("No parent task selected: %v", err)
		return
	}

	text := td.PromptInput("📌 Subtask text", "")
	if text == "" {
		warn("Subtask not created.")
		return
	}

	tasks, err := td.LoadTasks()
	if err != nil {
		fail("Failed to load tasks: %v", err)
		return
	}

	sub := td.Task{
		ID:        len(tasks) + 1,
		Text:      text,
		ParentID:  parent.ID,
		Completed: false,
	}

	tasks = append(tasks, sub)

	if err := td.SaveTasks(tasks); err != nil {
		fail("Failed to save subtask: %v", err)
		return
	}

	success("📌 Subtask added under #%d", parent.ID)
}
