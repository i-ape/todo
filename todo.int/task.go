package todo

import (
	"fmt"
	"strconv"
	"strings"
	"time"

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




// AddTaskWithDueDate adds a task with an optional due date
func AddTaskWithDueDate(text, due string) error {
	tasks, _ := LoadTasks()
	parsed := ""
	if due != "" {
		dt, err := ParseNaturalDate(due)
		if err != nil {
			return err
		}
		parsed = dt
	}
	newTask := Task{ID: len(tasks) + 1, Text: text, Completed: false, DueDate: parsed}
	tasks = append(tasks, newTask)
	return SaveTasks(tasks)
}

// ListTasks displays all tasks
func ListTasks() {
	tasks, _ := LoadTasks()
	if len(tasks) == 0 {
		color.Yellow("📭 No tasks available.")
		return
	}

	for _, task := range tasks {
		label := fmt.Sprintf("%d: %s", task.ID, task.Text)
		if task.DueDate != "" {
			label += fmt.Sprintf(" (Due: %s)", task.DueDate)
		}

		if task.Completed {
			fmt.Println(color.GreenString("[✓] %s", label))
			continue
		}

		if task.DueDate != "" {
			due, err := time.Parse("2006-01-02", task.DueDate)
			if err == nil && time.Now().After(due) {
				fmt.Println(color.RedString("[✗] %s", label))
				continue
			}
		}

		fmt.Println(color.CyanString("[ ] %s", label))
	}
}

// MarkTaskDone marks a task as completed
func MarkTaskDone(input string) error {
	tasks, _ := LoadTasks()
	found := false

	id, err := strconv.Atoi(input)
	for i, task := range tasks {
		if (err == nil && task.ID == id) || task.Text == input {
			tasks[i].Completed = true
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("task not found")
	}

	return SaveTasks(tasks)
}


// SetDueDate assigns a due date to a task
func SetDueDate(input string, dueDate string) error {
	tasks, _ := LoadTasks()
	found := false

	parsedDate, err := ParseNaturalDate(dueDate)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(input)
	for i, task := range tasks {
		if (err == nil && task.ID == id) || task.Text == input {
			tasks[i].DueDate = parsedDate
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("task not found")
	}

	return SaveTasks(tasks)
}
func EditTaskText(input, newText string) error {
	tasks, _ := LoadTasks()
	found := false
	id, err := strconv.Atoi(input)
	for i := range tasks {
		if (err == nil && tasks[i].ID == id) || tasks[i].Text == input {
			tasks[i].Text = newText
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("task not found")
	}
	return SaveTasks(tasks)
}


// DeleteTask removes a task by ID or text
func DeleteTask(input string) error {
	tasks, _ := LoadTasks()
	newTasks := []Task{}
	found := false

	id, err := strconv.Atoi(input)
	for _, task := range tasks {
		if (err == nil && task.ID == id) || task.Text == input {
			found = true
			continue
		}
		newTasks = append(newTasks, task)
	}

	if !found {
		return fmt.Errorf("task not found")
	}

	return SaveTasks(newTasks)
}
// ClearTasks deletes all tasks (empties the task list)
func ClearTasks() {
	_ = SaveTasks([]Task{})
}

// SearchTasks displays tasks that contain a keyword
func SearchTasks(keyword string) {
	tasks, _ := LoadTasks()
	matched := false
	for _, task := range tasks {
		if strings.Contains(strings.ToLower(task.Text), strings.ToLower(keyword)) {
			fmt.Printf("🔎 [%d] %s\n", task.ID, task.Text)
			matched = true
		}
	}
	if !matched {
		fmt.Println("🔍 No matching tasks found.")
	}
}
