package handlers

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	cmd "todo/cmd"
	"todo/config"
	td "todo/td"
)

// Move allows reordering tasks
func Move() {
	tasks, err := td.LoadTasks()
	if err != nil {
		fail("Failed to load tasks: %v", err)
		return
	}

	// Select task to move
	selected, err := cmd.SelectTasksWithFzf(false, config.DisableFzf)
	if err != nil || len(selected) == 0 {
		fail("No task selected: %v", err)
		return
	}

	task := selected[0]

	fmt.Printf("📦 Selected: #%d %s\n", task.ID, task.Text)
	fmt.Println("➡️  Enter new position (1 -", len(tasks), "):")

	newIndex := promptNumber()

	if newIndex < 1 || newIndex > len(tasks) {
		fail("Invalid position")
		return
	}

	// Remove from old position
	var updated []td.Task
	for _, t := range tasks {
		if t.ID != task.ID {
			updated = append(updated, t)
		}
	}

	// Insert at new position
	newIndex-- // zero-based index
	updated = append(updated[:newIndex], append([]td.Task{task}, updated[newIndex:]...)...)

	// Reassign IDs cleanly
	for i := range updated {
		updated[i].ID = i + 1
	}

	if err := td.SaveTasks(updated); err != nil {
		fail("Failed to save reordered tasks: %v", err)
		return
	}

	success("Moved #%d → position %d", task.ID, newIndex+1)
}

func promptNumber() int {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("> ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	n, _ := strconv.Atoi(input)
	return n
}
