package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	todo "todo/todo.int"

	"github.com/fatih/color"
)

// --- Handlers ---

func handleAdd() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: todo add [task text] [optional due date]")
		return
	}
	text := os.Args[2]
	due := ""
	if len(os.Args) > 3 {
		due = strings.Join(os.Args[3:], " ")
	}
	if err := AddTask(text, due); err != nil {
		fmt.Println("Error:", err)
	}
}

func handleEdit() {
	selected, err := todo.SelectTasksWithFzf(false, DisableFzf)

	if err != nil || len(selected) == 0 {
		fmt.Println("Select error:", err)
		return
	}
	task := selected[0]

	newText := todo.PromptInput(fmt.Sprintf("✏️  Edit task \"%s\"", task.Text), task.Text)

	if newText == "" {
		fmt.Println("No changes made.")
		return
	}
	if err := todo.EditTaskText(strconv.Itoa(task.ID), newText); err != nil {
		fmt.Println("Edit error:", err)
	}
}

func handleDone() {
	selected, err := todo.SelectTasksWithFzf(false, DisableFzf)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for _, task := range selected {
		if err := MarkTaskDone(strconv.Itoa(task.ID)); err != nil {
			fmt.Println("❌", err)
		}
	}
}

func handleDelete() {
	selected, err := todo.SelectTasksWithFzf(false, DisableFzf)

	if err != nil {
		fmt.Println("Error selecting task:", err)
		return
	}
	for _, task := range selected {
		if err := todo.DeleteTask(strconv.Itoa(task.ID)); err != nil {
			fmt.Println("❌", err)
		}
	}
}

func handleDue() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: todo due [task ID or task text] [date]")
		return
	}
	input := os.Args[2]
	dueDate := strings.Join(os.Args[3:], " ")
	if err := SetDueDate(input, dueDate); err != nil {
		fmt.Println("Error:", err)
	}
}

func handleSearch() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: todo search [keyword]")
		return
	}
	SearchTasks(os.Args[2])
}

func handlePriority() {
	selected, err := todo.SelectTasksWithFzf(false, DisableFzf)
	if err != nil || len(selected) == 0 {
		fmt.Println("Error selecting task:", err)
		return
	}
	task := selected[0]

	newPriority := todo.PromptInput("🔥 Set priority (high/medium/low)", task.Priority)
	if newPriority == "" {
		fmt.Println("No changes made.")
		return
	}
	parsed := todo.ParsePriority(newPriority)

	err = todo.UpdateTaskByID(task.ID, func(t *todo.Task) {
		t.Priority = parsed
	})
	if err != nil {
		fmt.Println("Failed to update priority:", err)
		return
	}
	fmt.Println("✅ Priority updated.")
}

func handleTags() {
	selected, err := todo.SelectTasksWithFzf(false, DisableFzf)

	if err != nil {
		fmt.Println("Error selecting task:", err)
		return
	}
	task := selected[0]

	raw := todo.PromptInput("🏷️ Edit tags (comma-separated)", strings.Join(task.Tags, ", "))
	if raw == "" {
		fmt.Println("No changes made.")
		return
	}
	tags := todo.ParseTags(raw)

	err = todo.UpdateTaskByID(task.ID, func(t *todo.Task) {
		t.Tags = tags
	})
	if err != nil {
		fmt.Println("❌", err)
		return
	}
	fmt.Println("✅ Tags updated.")
}

func handleRecurring() {
	selected, err := todo.SelectTasksWithFzf(false, DisableFzf)

	if err != nil {
		fmt.Println("Error selecting task:", err)
		return
	}
	task := selected[0]

	rec := todo.PromptInput("🔁 Set recurrence (daily, weekly, etc)", task.Recurring)
	if rec == "" {
		fmt.Println("No changes made.")
		return
	}

	err = todo.UpdateTaskByID(task.ID, func(t *todo.Task) {
		t.Recurring = rec
	})
	if err != nil {
		fmt.Println("❌", err)
	} else {
		fmt.Println("✅ Recurrence set.")
	}
}

func handleMove() {
	task, err := mustSelectSingleTask()
	if err != nil {
		fmt.Println("❌ Selection error:", err)
		return
	}

	tasks, err := todo.LoadTasks()
	if err != nil {
		fmt.Println("❌ Failed to load tasks:", err)
		return
	}

	fmt.Printf("🪄 Task: \"%s\"\n", task.Text)
	posStr := todo.PromptInput("↕️ New position (1-based index)", "")
	pos, err := strconv.Atoi(posStr)
	if err != nil || pos < 1 || pos > len(tasks) {
		fmt.Println("❌ Invalid position")
		return
	}

	// Remove task from current position
	var moved todo.Task
	newList := []todo.Task{}
	for _, t := range tasks {
		if t.ID == task.ID {
			moved = t
		} else {
			newList = append(newList, t)
		}
	}

	// Insert at new position (slice-safe)
	pos-- // convert to 0-based index
	if pos >= len(newList) {
		newList = append(newList, moved)
	} else {
		newList = append(newList[:pos], append([]todo.Task{moved}, newList[pos:]...)...)
	}

	if err := todo.SaveTasks(newList); err != nil {
		fmt.Println("❌ Failed to save:", err)
		return
	}

	fmt.Println("✅ Task moved successfully.")
}

func handleClear() {
	if err := ClearTasks(); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("✅ All tasks cleared.")
	}
}

func handleReset() {
	if err := ResetTasks(); err != nil {
		fmt.Println("⚠️ Reset failed:", err)
	} else {
		fmt.Println("🗑️ tasks.json deleted.")
	}
}
func handleList() {
	args := os.Args[2:]
	opts := todo.ListFilterOptions{}

	for _, arg := range args {
		switch {
		case arg == "--json":
			opts.JSONOutput = true
		case arg == "--done":
			opts.ShowDone = true
		case arg == "--pending":
			opts.ShowPending = true
		case arg == "--today":
			opts.TodayOnly = true
		case arg == "--overdue":
			opts.OverdueOnly = true
		case strings.HasPrefix(arg, "--tag="):
			opts.Tags = strings.TrimPrefix(arg, "--tag=")
		case strings.HasPrefix(arg, "--priority="):
			opts.Priority = strings.TrimPrefix(arg, "--priority=")
		}
	}

	tasks, err := todo.LoadTasks()
	if err != nil {
		fmt.Println("❌ Failed to load tasks:", err)
		return
	}

	filtered := todo.FilterAndSortTasks(tasks, opts)

	if opts.JSONOutput {
		data, _ := json.MarshalIndent(filtered, "", "  ")
		fmt.Println(string(data))
		return
	}

	for _, task := range filtered {
		label := fmt.Sprintf("%d: %s", task.ID, task.Text)
		if task.DueDate != "" {
			label += fmt.Sprintf(" (Due: %s)", task.DueDate)
		}
		if len(task.Tags) > 0 {
			label += " [" + strings.Join(task.Tags, ", ") + "]"
		}
		switch {
		case task.Completed:
			fmt.Println(color.GreenString("[✓] " + label))
		case todo.IsOverdue(task.DueDate):
			fmt.Println(color.RedString("[✗] " + label))
		default:
			fmt.Println(color.CyanString("[ ] " + label))
		}
	}
}

func mustSelectSingleTask() (*todo.Task, error) {
	tasks, err := todo.SelectTasksWithFzf(false, DisableFzf)
	if err != nil {
		return nil, err
	}
	return &tasks[0], nil
}
