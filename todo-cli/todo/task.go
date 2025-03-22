package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

const filename = "tasks.json"

type Task struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
}

// ✅ LoadTasks reads tasks from the JSON file
func LoadTasks() ([]Task, error) {
	var tasks []Task
	file, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []Task{}, nil // Return empty list if file doesn't exist
		}
		return nil, err
	}
	err = json.Unmarshal(file, &tasks)
	return tasks, err
}

// ✅ SaveTasks writes tasks to the JSON file
func SaveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

// ✅ AddTask adds a new task
func AddTask(text string) error {
	tasks, _ := LoadTasks()
	newTask := Task{ID: len(tasks) + 1, Text: text, Completed: false}
	tasks = append(tasks, newTask)
	return SaveTasks(tasks)
}

// ✅ ListTasks displays all tasks
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
		fmt.Printf("[%d] %s %s\n", task.ID, status, task.Text)
	}
}

// ✅ MarkTaskDone marks a task as completed
func MarkTaskDone(id int) error {
	tasks, _ := LoadTasks()
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Completed = true
			return SaveTasks(tasks)
		}
	}
	return fmt.Errorf("task with ID %d not found", id)
}

// ✅ DeleteTask removes a task
func DeleteTask(id int) error {
	tasks, _ := LoadTasks()
	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return SaveTasks(tasks)
		}
	}
	return fmt.Errorf("task with ID %d not found", id)
}
