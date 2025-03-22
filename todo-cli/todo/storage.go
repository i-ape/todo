package todo

import (
	"encoding/json"
	"errors"
	"os"
)

const filename = "tasks.json"

// ✅ Exported function (uppercase L)
func LoadTasks() ([]Task, error) {
	var tasks []Task
	file, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []Task{}, nil
		}
		return nil, err
	}
	err = json.Unmarshal(file, &tasks)
	return tasks, err
}

// ✅ Exported function (uppercase S)
func SaveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}
