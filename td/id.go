package td

import (
	"errors"
	"fmt"
)

var ErrTaskNotFound = errors.New("task not found")

// NextTaskID calculates the next available task ID.
func NextTaskID(tasks []Task) int {
	maxID := 0
	for _, task := range tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}
	return maxID + 1
}

// HasTaskID checks if a task with the given ID exists.
func HasTaskID(tasks []Task, id int) bool {
	for _, t := range tasks {
		if t.ID == id {
			return true
		}
	}
	return false
}

// GetTaskByID retrieves a task by its ID from disk.
func GetTaskByID(id int) (*Task, error) {
	tasks, err := LoadTasks()
	if err != nil {
		return nil, fmt.Errorf("GetTaskByID: %w", err)
	}
	for _, task := range tasks {
		if task.ID == id {
			t := task
			return &t, nil
		}
	}
	return nil, fmt.Errorf("%w: ID %d", ErrTaskNotFound, id)
}

// GetTaskByIDFromSlice finds a task by ID from an in-memory slice.
func GetTaskByIDFromSlice(tasks []Task, id int) (*Task, error) {
	for i := range tasks {
		if tasks[i].ID == id {
			return &tasks[i], nil
		}
	}
	return nil, fmt.Errorf("%w: ID %d", ErrTaskNotFound, id)
}

// UpdateTaskByID updates a task by ID using the provided function.
func UpdateTaskByID(id int, updateFn func(*Task)) error {
	tasks, err := LoadTasks()
	if err != nil {
		return fmt.Errorf("UpdateTaskByID: %w", err)
	}

	for i := range tasks {
		if tasks[i].ID == id {
			updateFn(&tasks[i])
			if err := SaveTasks(tasks); err != nil {
				return fmt.Errorf("UpdateTaskByID: save failed: %w", err)
			}
			return nil
		}
	}

	return fmt.Errorf("%w: ID %d", ErrTaskNotFound, id)
}

// DeleteTaskByID removes a task from storage and saves the updated list.
func DeleteTaskByID(id int) (*Task, error) {
	tasks, err := LoadTasks()
	if err != nil {
		return nil, fmt.Errorf("DeleteTaskByID: %w", err)
	}
	newTasks, deleted := RemoveTaskByID(tasks, id)
	if deleted == nil {
		return nil, fmt.Errorf("%w: ID %d", ErrTaskNotFound, id)
	}
	if err := SaveTasks(newTasks); err != nil {
		return nil, fmt.Errorf("DeleteTaskByID: save failed: %w", err)
	}
	return deleted, nil
}

// RemoveTaskByID removes the task with the given ID from the slice.
func RemoveTaskByID(tasks []Task, id int) ([]Task, *Task) {
	var deleted *Task
	newTasks := make([]Task, 0, len(tasks))
	for _, t := range tasks {
		if t.ID == id {
			tmp := t
			deleted = &tmp
			continue
		}
		newTasks = append(newTasks, t)
	}
	return newTasks, deleted
}

// ReplaceTaskByID replaces a task in the slice.
func ReplaceTaskByID(tasks []Task, updated Task) ([]Task, bool) {
	for i := range tasks {
		if tasks[i].ID == updated.ID {
			tasks[i] = updated
			return tasks, true
		}
	}
	return tasks, false
}

// MustGetTaskIndex returns the index of a task in the slice by ID or panics.
func MustGetTaskIndex(tasks []Task, id int) int {
	for i, t := range tasks {
		if t.ID == id {
			return i
		}
	}
	panic(fmt.Sprintf("task with ID %d not found", id))
}
