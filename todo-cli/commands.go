package main

import (
	"fmt"
	"os"
	"strconv"
)

// Add a new task
func AddTask(text string) error {
	tasks, _ := LoadTasks()
	newTask := Task{ID: len(tasks) + 1, Text: text, Completed: false}
	tasks = append(tasks, newTask)
	return SaveTasks(tasks)
}

// List all tasks
func ListTasks() {
	tasks, _ := LoadTasks()
	for _, task := range tasks {
		fmt.Println(task)
	}
}

// Mark a task as completed
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

// Delete a task
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

// Handle CLI commands
func HandleCommands() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: todo add|list|done|delete [task text or task ID]")
		return
	}

	cmd := os.Args[1]

	switch cmd {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: todo add [task text]")
			return
		}
		text := os.Args[2]
		if err := AddTask(text); err != nil {
			fmt.Println("Error:", err)
		}

	case "list":
		ListTasks()

	case "done":
		if len(os.Args) < 3 {
			fmt.Println("Usage: todo done [task ID]")
			return
		}
		id, _ := strconv.Atoi(os.Args[2])
		if err := MarkTaskDone(id); err != nil {
			fmt.Println("Error:", err)
		}

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: todo delete [task ID]")
			return
		}
		id, _ := strconv.Atoi(os.Args[2])
		if err := DeleteTask(id); err != nil {
			fmt.Println("Error:", err)
		}

	default:
		fmt.Println("Invalid command. Use add|list|done|delete")
	}
}
