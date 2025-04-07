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
	DueDate   time.Time `json:"due_date,omitempty"`
}

// AddTask adds a new task
func AddTask(text string) error {
	tasks, err := LoadTasks()
	if err != nil {
		return err
	}
	newTask := Task{ID: len(tasks) + 1, Text: text, Completed: false}
	tasks = append(tasks, newTask)
	return SaveTasks(tasks)
}

// ListTasks displays all tasks
func ListTasks() {
	tasks, err := LoadTasks()
	if err != nil {
		color.Red("Error loading tasks: %v", err)
		return
	}

	if len(tasks) == 0 {
		color.Yellow("📭 No tasks available.")
		return
	}

	now := time.Now()
	for _, task := range tasks {
		status := color.CyanString("[ ] %d: %s", task.ID, task.Text)
		if task.Completed {
			status = color.GreenString("[✓] %d: %s", task.ID, task.Text)
		}

		if !task.DueDate.IsZero() {
			dateStr := task.DueDate.Format("2006-01-02")
			if !task.Completed && task.DueDate.Before(now) {
				status += color.RedString(" (OVERDUE: %s)", dateStr)
			} else {
				status += color.MagentaString(" (Due: %s)", dateStr)
			}
		}

		fmt.Println(status)
	}
}

// MarkTaskDone marks a task as completed
func MarkTaskDone(input string) error {
	tasks, err := LoadTasks()
	if err != nil {
		return err
	}
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
	tasks, err := LoadTasks()
	if err != nil {
		return err
	}
	found := false

	id, parseErr := strconv.Atoi(input)
	for i, task := range tasks {
		if (parseErr == nil && task.ID == id) || task.Text == input {
			var parsed time.Time
			parsed, err = time.Parse("02-01-2006", dueDate)
			if err != nil {
				parsed, err = time.Parse("2006-01-02", dueDate)
				if err != nil {
					return fmt.Errorf("invalid date format, use DD-MM-YYYY or YYYY-MM-DD")
				}
			}
			tasks[i].DueDate = parsed
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
	tasks, err := LoadTasks()
	if err != nil {
		return err
	}
	newTasks := []Task{}
	found := false

	id, parseErr := strconv.Atoi(input)
	for _, task := range tasks {
		if (parseErr == nil && task.ID == id) || task.Text == input {
			found = true
			continue
		}
		newTasks = append(newTasks, task)
	}

	if !found {
		return fmt.Errorf("task not found")
	}

	return SaveTasks(newTasks)
}