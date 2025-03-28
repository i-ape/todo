package todo

import (
	"fmt"
	"github.com/fatih/color"
)

// Task struct
type Task struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
	DueDate   string `json:"due_date,omitempty"`
}

// AddTask adds a task
func AddTask(text string) error {
	tasks, _ := LoadTasks()
	newTask := Task{ID: len(tasks) + 1, Text: text, Completed: false}
	tasks = append(tasks, newTask)
	return SaveTasks(tasks)
}

// ListTasks prints all tasks
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
		fmt.Println(status)
	}
}

// MarkTaskDone marks a task as completed
func MarkTaskDone(id int) error {
	tasks, _ := LoadTasks()
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Completed = true
			return SaveTasks(tasks)
		}
	}
	return fmt.Errorf("task not found")
}

// DeleteTask deletes a task (✅ Ensure function name is capitalized for export)
func DeleteTask(id int) error {
	tasks, _ := LoadTasks()
	newTasks := []Task{}
	for _, task := range tasks {
		if task.ID != id {
			newTasks = append(newTasks, task)
		}
	}
	if len(tasks) == len(newTasks) {
		return fmt.Errorf("task not found")
	}
	return SaveTasks(newTasks)
}
