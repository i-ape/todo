// cmd/handlers_edit.go
package main

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"todo/config"
	todo "todo/td"
)

// --- Edit task text ---
func handleEdit() {
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

	err = todo.UpdateTaskByID(id, func(t *todo.Task) {
		t.Text = newText
	})
	if err != nil {
		if errors.Is(err, todo.ErrTaskNotFound) {
			warn("Task with ID %d not found.", id)
			return
		}
		fail("Failed to edit task: %v", err)
		return
	}

	info("✅ Task #%d updated successfully!", id)
}

// --- Mark as done ---
func handleDone() {
	tasks, err := todo.SelectTasksWithFzf(false, config.DisableFzf)
	if err != nil {
		fail("Selection error: %v", err)
		return
	}

	for _, t := range tasks {
		if err := todo.MarkTaskDone(strconv.Itoa(t.ID)); err != nil {
			fail("Could not mark done: %v", err)
		} else {
			success("Marked done: %s", t.Text)
		}
	}
}

// --- Change priority ---
func handlePriority() {
	task, err := selectSingleTask()
	if err != nil {
		fail("Error selecting task: %v", err)
		return
	}

	newPriority := todo.PromptInput("🔥 Set priority (high/medium/low)", task.Priority)
	if newPriority == "" {
		warn("No changes made.")
		return
	}
	parsed := todo.ParsePriority(newPriority)

	err = todo.UpdateTaskByID(task.ID, func(t *todo.Task) {
		t.Priority = parsed
	})
	if err != nil {
		fail("Failed to update priority: %v", err)
		return
	}
	success("Priority updated to %s.", newPriority)
}

// --- Edit tags ---
func handleTags() {
	task, err := selectSingleTask()
	if err != nil {
		fail("Error selecting task: %v", err)
		return
	}

	raw := todo.PromptInput("🏷️ Edit tags (comma-separated)", strings.Join(task.Tags, ", "))
	if raw == "" {
		warn("No changes made.")
		return
	}
	tags := todo.ParseTags(raw)

	err = todo.UpdateTaskByID(task.ID, func(t *todo.Task) {
		t.Tags = tags
	})
	if err != nil {
		fail("Failed to update tags: %v", err)
		return
	}
	success("Tags updated.")
}

// --- Edit recurrence ---
func handleRecurring() {
	task, err := selectSingleTask()
	if err != nil {
		fail("Error selecting task: %v", err)
		return
	}

	rec := todo.PromptInput("🔁 Set recurrence (daily, weekly, etc)", task.Recurring)
	if rec == "" {
		warn("No changes made.")
		return
	}

	err = todo.UpdateTaskByID(task.ID, func(t *todo.Task) {
		t.Recurring = rec
	})
	if err != nil {
		fail("Failed to update recurrence: %v", err)
		return
	}
	success("Recurrence set to %s.", rec)
}
