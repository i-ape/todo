package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"todo/config"
)

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

func SaveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(config.GetTaskFilePath(), data, 0644)
}

func GetTaskByID(id int) (*Task, error) {
	tasks, err := LoadTasks()
	if err != nil {
		return nil, err
	}
	for _, task := range tasks {
		if task.ID == id {
			return &task, nil
		}
	}
	return nil, fmt.Errorf("task with ID %d not found", id)
}
