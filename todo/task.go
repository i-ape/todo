package todo

import (
	"fmt"
	"strconv"
	"time"

	"github.com/fatih/color"
)

// Task struct represents a single task
type Task struct {
	ID        int       `json:"id"`
	Text      string    `json:"text"`
	Completed bool      `json:"completed"`
	DueDate   string    `json:"due_date,omitempty"`
}

// AddTask adds a new task
func AddTask(text string) error {
	tasks, _ := LoadTasks()
	newTask := Task{ID: len(tasks) + 1, Text: text, Completed: false}
	tasks = append(tasks, newTask)
	return SaveTasks(tasks)
}

// ListTasks displays all tasks
func ListTasks() {
	tasks, _ := LoadTasks()
	if len(tasks) == 0 {
		color.Yellow("📭 No tasks available.")
		return
	}

	for _, task := range tasks {
		status := color.CyanString("[ ] %d: %s", task.ID, task.Text)
		if task.Completed {
			status = color.GreenString("[✓] %d: %s", task.ID, task.Text)
		}

		// Show due date if available
		if task.DueDate != "" {
			status += color.MagentaString(" (Due: %s)", task.DueDate)
		}

		fmt.Println(status)
	}
}

// MarkTaskDone marks a task as completed
func MarkTaskDone(input string) error {
	tasks, _ := LoadTasks()
	found := false

	id, err := strconv.Atoi(input)
	for i, task := range tasks {
		if (err == nil && task.ID == id) || task.Text == input {
			tasks[i].Completed = true
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("task not found")
	}

	return SaveTasks(tasks)
}

// SetDueDate assigns a due date to a task
func SetDueDate(input string, dueDate string) error {
	tasks, _ := LoadTasks()
	found := false

	id, err := strconv.Atoi(input)
	for i, task := range tasks {
		if (err == nil && task.ID == id) || task.Text == input {
			_, err := time.Parse("2006-01-02", dueDate) // Validate format
			if err != nil {
				return fmt.Errorf("invalid date format, use YYYY-MM-DD")
			}
			tasks[i].DueDate = dueDate
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("task not found")
	}

	return SaveTasks(tasks)
}

// DeleteTask removes a task by ID or text
func DeleteTask(input string) error {
	tasks, _ := LoadTasks()
	newTasks := []Task{}
	found := false

	id, err := strconv.Atoi(input)
	for _, task := range tasks {
		if (err == nil && task.ID == id) || task.Text == input {
			found = true
			continue // Skip this task (delete it)
		}
		newTasks = append(newTasks, task)
	}

	if !found {
		return fmt.Errorf("task not found")
	}

	return SaveTasks(newTasks)
}
