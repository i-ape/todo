// task.go
package todo  // ✅ Must use package "todo"

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

// Task struct represents a to-do task
type Task struct {
	ID        int       `json:"id"`
	Text      string    `json:"text"`
	Completed bool      `json:"completed"`
	DueDate   string    `json:"due_date"`
}

// AddTask adds a new task with an optional due date
func AddTask(text, dueDate string) error {
	tasks, _ := LoadTasks()
	newTask := Task{ID: len(tasks) + 1, Text: text, Completed: false, DueDate: dueDate}
	tasks = append(tasks, newTask)
	return SaveTasks(tasks)
}

// ListTasks displays all tasks with due dates
func ListTasks() {
	tasks, _ := LoadTasks()
	if len(tasks) == 0 {
		color.Yellow("📭 No tasks available.")
		return
	}
	for _, task := range tasks {
		var status string
		if task.Completed {
			status = color.GreenString("[✅] %d: %s (Due: %s)", task.ID, task.Text, task.DueDate)
		} else {
			// Check if task is overdue
			due, _ := time.Parse("2006-01-02", task.DueDate)
			if time.Now().After(due) {
				status = color.RedString("[❗ OVERDUE ❗] %d: %s (Due: %s)", task.ID, task.Text, task.DueDate)
			} else {
				status = color.CyanString("[❌] %d: %s (Due: %s)", task.ID, task.Text, task.DueDate)
			}
		}
		fmt.Println(status)
	}
}

// ClearTasks removes all tasks by deleting the tasks.json file
func ClearTasks() error {
	return os.Remove("tasks.json")
}