package todo

import (
	"encoding/json"
	"errors"
	"os"
	"todo/todo.int/config"
)

// LoadTasks reads tasks from disk using config-defined path
func LoadTasks() ([]Task, error) {
	var tasks []Task
	file, err := os.ReadFile(config.GetTaskFilePath())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []Task{}, nil
		}
		return nil, err
	}
	err = json.Unmarshal(file, &tasks)
	return tasks, err
}

// SaveTasks writes tasks to disk
func SaveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(config.GetTaskFilePath(), data, 0644)
}
