package todo

import (
	"fmt"
)

// Task struct represents a to-do task
type Task struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
}

// AddTask adds a new task
func AddTask(text string) error {
	tasks, _ := LoadTasks() // Calls `LoadTasks()` from storage.go
	newTask := Task{ID: len(tasks) + 1, Text: text, Completed: false}
	tasks = append(tasks, newTask)
	return SaveTasks(tasks) // Calls `SaveTasks()` from storage.go
}

// ListTasks displays all tasks
func ListTasks() {
	tasks, _ := LoadTasks() // Calls `LoadTasks()` from storage.go
	if len(tasks) == 0 {
		fmt.Println("📭 No tasks available.")
		return
	}
	for _, task := range tasks {
		status := "❌"
		if task.Completed {
			status = "✅"
		}
		fmt.Printf("[%d] %s %s\n", task.ID, status, task.Text)
	}
}

// MarkTaskDone marks a task as completed
func MarkTaskDone(id int) error {
	tasks, _ := LoadTasks() // Calls `LoadTasks()` from storage.go
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Completed = true
			return SaveTasks(tasks) // Calls `SaveTasks()` from storage.go
		}
	}
	return fmt.Errorf("task with ID %d not found", id)
}

// DeleteTask removes a task
func DeleteTask(id int) error {
	tasks, _ := LoadTasks() // Calls `LoadTasks()` from storage.go
	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return SaveTasks(tasks) // Calls `SaveTasks()` from storage.go
		}
	}
	return fmt.Errorf("task with ID %d not found", id)
}
