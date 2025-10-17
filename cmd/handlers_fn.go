// cmd/handlers_fn.go
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"todo/config"
	todo "todo/td"
)

// ==============================
// 🗓️ DUE DATE
// ==============================

func handleDue() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: todo due [task ID or text] [date]")
		return
	}
	input := os.Args[2]
	dueDate := strings.Join(os.Args[3:], " ")

	task, err := getTaskByInput(input)
	if err != nil {
		fail("Task not found: %v", err)
		return
	}

	err = todo.UpdateTaskByID(task.ID, func(t *todo.Task) {
		t.Due = dueDate
	})
	if err != nil {
		fail("Failed to update due date: %v", err)
		return
	}
	success("Due date updated for #%d → %s", task.ID, dueDate)
}

// ==============================
// 🗑️ DELETE / CLEAR / RESET
// ==============================

func handleDelete() {
	selected, err := todo.SelectTasksWithFzf(false, config.DisableFzf)
	if err != nil || len(selected) == 0 {
		fail("No task selected: %v", err)
		return
	}

	for _, task := range selected {
		if err := todo.DeleteTask(strconv.Itoa(task.ID)); err != nil {
			fail("Could not delete #%d: %v", task.ID, err)
		} else {
			success("Deleted #%d", task.ID)
		}
	}
}

func handleClear() {
	if err := todo.SaveTasks([]todo.Task{}); err != nil {
		fail("Error clearing tasks: %v", err)
		return
	}
	success("All tasks cleared.")
}

func handleReset() {
	path := config.GetTaskFilePath()
	if err := os.Remove(path); err != nil {
		fail("Failed to delete %s: %v", path, err)
		return
	}
	success("tasks.json deleted.")
}

// ==============================
// ↕️ MOVE TASK
// ==============================

func handleMove() {
	task, err := selectSingleTask()
	if err != nil {
		fail("Selection error: %v", err)
		return
	}

	tasks, err := todo.LoadTasks()
	if err != nil {
		fail("Failed to load tasks: %v", err)
		return
	}

	posStr := todo.PromptInput(fmt.Sprintf("↕️ New position for \"%s\" (1-%d)", task.Text, len(tasks)), "")
	pos, err := strconv.Atoi(posStr)
	if err != nil || pos < 1 || pos > len(tasks) {
		fail("Invalid position.")
		return
	}

	// Remove and reinsert task
	var moved todo.Task
	newList := []todo.Task{}
	for _, t := range tasks {
		if t.ID == task.ID {
			moved = t
		} else {
			newList = append(newList, t)
		}
	}
	pos-- // 0-based
	if pos >= len(newList) {
		newList = append(newList, moved)
	} else {
		newList = append(newList[:pos], append([]todo.Task{moved}, newList[pos:]...)...)
	}

	if err := todo.SaveTasks(newList); err != nil {
		fail("Failed to save: %v", err)
		return
	}
	success("Task moved successfully.")
}

// ==============================
// 📝 NOTES & SUBTASKS
// ==============================

func handleNote() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: todo note [task ID or text] [note text]")
		return
	}
	input := os.Args[2]
	note := strings.Join(os.Args[3:], " ")

	task, err := getTaskByInput(input)
	if err != nil {
		fail("Task not found: %v", err)
		return
	}

	err = todo.UpdateTaskByID(task.ID, func(t *todo.Task) {
		t.Notes = note
	})
	if err != nil {
		fail("Failed to update note: %v", err)
		return
	}
	success("Note added to #%d", task.ID)
}

func handleSubtask() {
	selected, err := todo.SelectTasksWithFzf(false, config.DisableFzf)
	if err != nil || len(selected) == 0 {
		fail("No parent task selected: %v", err)
		return
	}
	parent := selected[0]

	text := todo.PromptInput("📌 Subtask text", "")
	if text == "" {
		warn("Subtask not created.")
		return
	}

	tasks, _ := todo.LoadTasks()
	subtask := todo.Task{
		ID:        len(tasks) + 1,
		Text:      text,
		ParentID:  parent.ID,
		Completed: false,
	}
	tasks = append(tasks, subtask)

	if err := todo.SaveTasks(tasks); err != nil {
		fail("Failed to save subtask: %v", err)
		return
	}
	success("Subtask added under #%d", parent.ID)
}

// ==============================
// 🔍 PICK / SEARCH
// ==============================

func handlePick() {
	tasks, err := todo.SelectTasksWithFzf(true, config.DisableFzf)
	if err != nil || len(tasks) == 0 {
		fail("No task selected: %v", err)
		return
	}

	switch {
	case containsArg("--json"):
		data, _ := json.MarshalIndent(tasks, "", "  ")
		fmt.Println(string(data))

	case containsArg("--edit"):
		for _, task := range tasks {
			newText := todo.PromptInput(fmt.Sprintf("✏️ Edit \"%s\"", task.Text), task.Text)
			if newText != "" {
				_ = todo.UpdateTaskByID(task.ID, func(t *todo.Task) { t.Text = newText })
				success("Edited #%d", task.ID)
			}
		}

	case containsArg("--note"):
		for _, task := range tasks {
			note := todo.PromptMultiline("📝 Note for "+task.Text, task.Notes)
			if note != "" {
				_ = todo.UpdateTaskByID(task.ID, func(t *todo.Task) { t.Notes = note })
				success("Note saved for #%d", task.ID)
			}
		}

	default:
		for _, t := range tasks {
			fmt.Println(t.ID)
		}
	}
}

// ==============================
// 🆘 HELP
// ==============================

func handleHelp() {
	fmt.Println(`
Todo CLI — Commands:
  add        Add a new task
  edit       Edit an existing task
  list       Show all tasks
  due        Set a due date for a task
  note       Add or edit a note
  move       Move task position
  delete     Delete selected task(s)
  clear      Remove all tasks
  reset      Delete task file
  subtask    Create a subtask
  pick       Select task(s) interactively
  help       Show this help message

Examples:
  todo add "Buy milk" tomorrow
  todo due 1 Monday
  todo note 2 "Remember lactose-free"
`)
}