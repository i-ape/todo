// task.go
ppackage main

import (
	"fmt"
	"time"
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
		fmt.Println("📭 No tasks available.")
		return
	}
	for _, task := range tasks {
		status := "❌"
		if task.Completed {
			status = "✅"
		}
		fmt.Printf("[%d] %s %s (Due: %s)\n", task.ID, status, task.Text, task.DueDate)
	}
}