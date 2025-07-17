package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"todo/config"
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
	selected, err := td.SelectTasksWithFzf(false, config.DisableFzf)

	if err != nil || len(selected) == 0 {
		fmt.Println("Select error:", err)
		return
	}
	task := selected[0]

	newText := td.PromptInput(fmt.Sprintf("✏️  Edit task \"%s\"", task.Text), task.Text)

	if newText == "" {
		fmt.Println("No changes made.")
		return
	}
	if err := td.EditTaskText(strconv.Itoa(task.ID), newText); err != nil {
		fmt.Println("Edit error:", err)
	}
}

func handleDone() {
	selected, err := td.SelectTasksWithFzf(false, config.DisableFzf)

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
	selected, err := td.SelectTasksWithFzf(false, config.DisableFzf)

	if err != nil {
		fmt.Println("Error selecting task:", err)
		return
	}
	for _, task := range selected {
		if err := td.DeleteTask(strconv.Itoa(task.ID)); err != nil {
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
	selected, err := td.SelectTasksWithFzf(false, config.DisableFzf)
	if err != nil || len(selected) == 0 {
		fmt.Println("Error selecting task:", err)
		return
	}
	task := selected[0]

	newPriority := td.PromptInput("🔥 Set priority (high/medium/low)", task.Priority)
	if newPriority == "" {
		fmt.Println("No changes made.")
		return
	}
	parsed := td.ParsePriority(newPriority)

	err = td.UpdateTaskByID(task.ID, func(t *td.Task) {
		t.Priority = parsed
	})
	if err != nil {
		fmt.Println("Failed to update priority:", err)
		return
	}
	fmt.Println("✅ Priority updated.")
}

func handleTags() {
	selected, err := td.SelectTasksWithFzf(false, config.DisableFzf)

	if err != nil {
		fmt.Println("Error selecting task:", err)
		return
	}
	task := selected[0]

	raw := td.PromptInput("🏷️ Edit tags (comma-separated)", strings.Join(task.Tags, ", "))
	if raw == "" {
		fmt.Println("No changes made.")
		return
	}
	tags := td.ParseTags(raw)

	err = td.UpdateTaskByID(task.ID, func(t *td.Task) {
		t.Tags = tags
	})
	if err != nil {
		fmt.Println("❌", err)
		return
	}
	fmt.Println("✅ Tags updated.")
}

func handleRecurring() {
	selected, err := td.SelectTasksWithFzf(false, config.DisableFzf)

	if err != nil {
		fmt.Println("Error selecting task:", err)
		return
	}
	task := selected[0]

	rec := td.PromptInput("🔁 Set recurrence (daily, weekly, etc)", task.Recurring)
	if rec == "" {
		fmt.Println("No changes made.")
		return
	}

	err = td.UpdateTaskByID(task.ID, func(t *td.Task) {
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

	tasks, err := td.LoadTasks()
	if err != nil {
		fmt.Println("❌ Failed to load tasks:", err)
		return
	}

	fmt.Printf("🪄 Task: \"%s\"\n", task.Text)
	posStr := td.PromptInput("↕️ New position (1-based index)", "")
	pos, err := strconv.Atoi(posStr)
	if err != nil || pos < 1 || pos > len(tasks) {
		fmt.Println("❌ Invalid position")
		return
	}

	// Remove task from current position
	var moved td.Task
	newList := []td.Task{}
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
		newList = append(newList[:pos], append([]td.Task{moved}, newList[pos:]...)...)
	}

	if err := td.SaveTasks(newList); err != nil {
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
	opts := td.ListFilterOptions{}
	showNotes := false

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
		case arg == "--notes", arg == "--verbose":
			showNotes = true
		case strings.HasPrefix(arg, "--tag="):
			opts.Tags = strings.TrimPrefix(arg, "--tag=")
		case strings.HasPrefix(arg, "--priority="):
			opts.Priority = strings.TrimPrefix(arg, "--priority=")
		}
	}

	tasks, err := td.LoadTasks()
	if err != nil {
		fmt.Println("❌ Failed to load tasks:", err)
		return
	}
	filtered := td.FilterTasks(tasks, opts)

	// map parent ID → children
	children := map[int][]td.Task{}
	parents := []td.Task{}
	for _, task := range filtered {
		if task.ParentID > 0 {
			children[task.ParentID] = append(children[task.ParentID], task)
		} else {
			parents = append(parents, task)
		}
	}

	if opts.JSONOutput {
		data, _ := json.MarshalIndent(filtered, "", "  ")
		fmt.Println(string(data))
		return
	}

	for _, parent := range parents {
		td.PrintTask(parent, showNotes)
		for _, child := range children[parent.ID] {
			td.PrintTaskIndented(child, showNotes)
		}
	}
}

func mustSelectSingleTask() (*td.Task, error) {
	tasks, err := td.SelectTasksWithFzf(false, config.DisableFzf)
	if err != nil {
		return nil, err
	}
	return &tasks[0], nil
}

func HandleNoteInteractive() {
	selected, err := td.SelectTasksWithFzf(false, config.DisableFzf)
	if err != nil || len(selected) == 0 {
		fmt.Println("❌ Error selecting task:", err)
		return
	}
	task := selected[0]

	note := td.PromptMultiline("📝 Edit note", task.Notes)
	if note == "" {
		fmt.Println("No changes made.")
		return
	}

	err = td.UpdateTaskByID(task.ID, func(t *td.Task) {
		t.Notes = note
	})
	if err != nil {
		fmt.Println("❌ Failed to update notes:", err)
		return
	}
	fmt.Println("✅ Notes saved.")
}

func handleNote() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: todo note [task ID or text] [note text]")
		return
	}
	input := os.Args[2]
	note := strings.Join(os.Args[3:], " ")
	id, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("❌ Invalid task ID")
		return
	}

	err = td.UpdateTaskByID(id, func(t *td.Task) {
		t.Notes = note
	})
	if err != nil {
		fmt.Println("❌ Failed to update notes:", err)
		return
	}
	fmt.Println("✅ Notes updated.")
}

func handleSubtask() {
	selected, err := td.SelectTasksWithFzf(false, config.DisableFzf)
	if err != nil || len(selected) == 0 {
		fmt.Println("❌ Error selecting parent task:", err)
		return
	}
	parent := selected[0]

	text := td.PromptInput("📌 Subtask text", "")
	if text == "" {
		fmt.Println("Subtask not created.")
		return
	}

	tasks, _ := td.LoadTasks()
	subtask := td.Task{
		ID:        len(tasks) + 1,
		Text:      text,
		ParentID:  parent.ID,
		Completed: false,
	}
	tasks = append(tasks, subtask)

	if err := td.SaveTasks(tasks); err != nil {
		fmt.Println("❌ Failed to save subtask:", err)
		return
	}
	fmt.Println("✅ Subtask added.")
}
func handlePick() {
	tasks, err := td.SelectTasksWithFzf(true, config.DisableFzf)
	if err != nil || len(tasks) == 0 {
		fmt.Println("❌ No task selected:", err)
		return
	}

	switch {
	case containsArg("--json"):
		data, _ := json.MarshalIndent(tasks, "", "  ")
		fmt.Println(string(data))

	case containsArg("--edit"):
		for _, task := range tasks {
			newText := td.PromptInput(fmt.Sprintf("✏️ Edit \"%s\"", task.Text), task.Text)
			if newText != "" {
				_ = td.UpdateTaskByID(task.ID, func(t *td.Task) {
					t.Text = newText
				})
				fmt.Println("✅ Edited:", task.ID)
			}
		}

	case containsArg("--note"):
		for _, task := range tasks {
			note := td.PromptMultiline("📝 Note for "+task.Text, task.Notes)
			if note != "" {
				_ = td.UpdateTaskByID(task.ID, func(t *td.Task) {
					t.Notes = note
				})
				fmt.Println("✅ Note saved:", task.ID)
			}
		}

	default:
		for _, t := range tasks {
			fmt.Println(t.ID)
		}
	}
}

func containsArg(flag string) bool {
	for _, arg := range os.Args {
		if arg == flag {
			return true
		}
	}
	return false
}
func handleShow() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: todo show [id]")
		return
	}
	id, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Invalid ID")
		return
	}
	task, err := td.GetTaskByID(id)
	if err != nil {
		fmt.Println("❌", err)
		return
	}
	td.PrintTaskDetails(*task) // new function to show full info
}

func handleHelp() {
	if len(os.Args) == 2 {
		printHelp()
		return
	}
	switch os.Args[2] {
	case "add":
		fmt.Println("Usage: todo add \"task text\" [due]")
	case "note":
		fmt.Println("Usage: todo note 1 \"Note content\"")
	// ... add for others
	default:
		fmt.Println("Unknown command. Try `todo help`")
	}
}
