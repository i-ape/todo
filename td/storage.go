package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"todo/config"
)

// LoadTasks loads tasks from the JSON file.
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

	// Basic sanity check
	idMap := make(map[int]struct{})
	for _, task := range tasks {
		if _, exists := idMap[task.ID]; exists {
			return nil, fmt.Errorf("duplicate task ID %d found in file", task.ID)
		}
		idMap[task.ID] = struct{}{}
	}

	return tasks, nil
}

// SaveTasks saves tasks to the JSON file with backup.
func SaveTasks(tasks []Task) error {
	filePath := config.GetTaskFilePath()
	backupPath := config.GetBackupTaskFilePath()

	// Backup if file exists
	if _, err := os.Stat(filePath); err == nil {
		data, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to backup tasks: %w", err)
		}
		if err = os.WriteFile(backupPath, data, 0644); err != nil {
			return fmt.Errorf("failed to write backup %s: %w", backupPath, err)
		}
	}

	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tasks: %w", err)
	}

	if err = os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write tasks file %s: %w", filePath, err)
	}

	fmt.Println(">>> Saving tasks to:", filePath)
	return nil
}

// ResetTasks deletes both the task file and its backup.
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
