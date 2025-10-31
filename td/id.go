package todo

import (
	"fmt"
)

// NextTaskID calculates the next available task ID.
// It ensures new tasks always have a unique, incremented ID.
func NextTaskID(tasks []Task) int {
	maxID := 0
	for _, task := range tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}
	return maxID + 1
}

// GetTaskByID retrieves a task by its ID from disk.
// Returns a copy of the task, so modifying it won't affect stored data directly.
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

// RemoveTaskByID removes the task with the given ID from the slice.
// Returns the updated slice and the removed task (or nil if not found).
func RemoveTaskByID(tasks []Task, id int) ([]Task, *Task) {
	var deleted *Task
	newTasks := make([]Task, 0, len(tasks))
	for _, t := range tasks {
		if t.ID == id {
			tmp := t // ensure pointer doesn’t reference loop variable
			deleted = &tmp
			continue
		}
		newTasks = append(newTasks, t)
	}
	return newTasks, deleted
}

// UpdateTaskByID updates a task by ID using the provided update function.
// It loads tasks, applies the update, and saves them automatically.
func UpdateTaskByID(id int, updateFn func(*Task)) error {
	tasks, err := LoadTasks()
	if err != nil {
		return fmt.Errorf("failed to load tasks: %w", err)
	}

	found := false
	for i := range tasks {
		if tasks[i].ID == id {
			updateFn(&tasks[i])
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("task with ID %d not found", id)
	}

	if err := SaveTasks(tasks); err != nil {
		return fmt.Errorf("failed to save updated tasks: %w", err)
	}

	return nil
}

// MustGetTaskIndex returns the index of a task in the slice by ID.
// It panics if the task does not exist. Useful for internal logic.
func MustGetTaskIndex(tasks []Task, id int) int {
	for i, t := range tasks {
		if t.ID == id {
			return i
		}
	}
	panic(fmt.Sprintf("task with ID %d not found", id))
}

// ReplaceTaskByID replaces a task in the slice and returns a new slice.
func ReplaceTaskByID(tasks []Task, updated Task) []Task {
	for i, t := range tasks {
		if t.ID == updated.ID {
			tasks[i] = updated
			break
		}
	}
	return tasks
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
