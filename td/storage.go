package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"todo/config"
)

// LoadTasks loads tasks from the JSON file, returning empty if not found
func LoadTasks() ([]Task, error) {
	filePath := config.GetTaskFilePath()
	data, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []Task{}, nil
		}
		return nil, fmt.Errorf("failed to read tasks file %s: %w", filePath, err)
	}

	var tasks []Task
	if err = json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tasks: %w", err)
	}

	// Optional validation: Check for duplicate IDs
	idMap := make(map[int]struct{})
	for _, task := range tasks {
		if _, exists := idMap[task.ID]; exists {
			return nil, fmt.Errorf("duplicate task ID %d found in file", task.ID)
		}
		idMap[task.ID] = struct{}{}
	}

	return tasks, nil
}

// SaveTasks saves tasks to the JSON file with backup
func SaveTasks(tasks []Task) error {
	filePath := config.GetTaskFilePath()
	backupPath := config.GetBackupTaskFilePath()

	// Create backup if file exists
	if _, err := os.Stat(filePath); err == nil {
		data, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to backup tasks: %w", err)
		}
		if err = os.WriteFile(backupPath, data, 0644); err != nil {
			return fmt.Errorf("failed to write backup %s: %w", backupPath, err)
		}
	}

	// Marshal and save
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tasks: %w", err)
	}

	if err = os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write tasks file %s: %w", filePath, err)
	}

	return nil
}

// ResetTasks deletes the tasks file and backup
func ResetTasks() error {
	filePath := config.GetTaskFilePath()
	backupPath := config.GetBackupTaskFilePath()

	if err := os.Remove(filePath); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to delete tasks file %s: %w", filePath, err)
	}
	if err := os.Remove(backupPath); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to delete backup file %s: %w", backupPath, err)
	}
	return nil
}

// NextTaskID calculates the next available task ID
func NextTaskID(tasks []Task) int {
	maxID := 0
	for _, task := range tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}
	return maxID + 1
}

// GetTaskByID retrieves a task by ID
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
