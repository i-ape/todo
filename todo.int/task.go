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
	tasks, err := LoadTasks()
	if err != nil {
		return err
	}
	newTask := Task{ID: len(tasks) + 1, Text: text, Completed: false}
	tasks = append(tasks, newTask)
	return SaveTasks(tasks)
}

// AddTaskWithDueDate adds a task with an optional due date
func AddTaskWithDueDate(text, due string) error {
	tasks, err := LoadTasks()
	if err != nil {
		return err
	}
	parsedDue := ""
	if due != "" {
		parsedDue, err = ParseNaturalDate(due)
		if err != nil {
			return err
		}
	}
	newTask := Task{ID: len(tasks) + 1, Text: text, Completed: false, DueDate: parsedDue}
	tasks = append(tasks, newTask)
	return SaveTasks(tasks)
}

// ListTasks displays all tasks
func ListTasks() {
	tasks, err := LoadTasks()
	if err != nil {
		color.Red("Failed to load tasks: %v", err)
		return
	}
	if len(tasks) == 0 {
		color.Yellow("📭 No tasks available.")
		return
	}

	for _, task := range tasks {
		displayTask(task)
	}
}

func displayTask(task Task) {
	label := fmt.Sprintf("%d: %s", task.ID, task.Text)
	if task.DueDate != "" {
		label += fmt.Sprintf(" (Due: %s)", task.DueDate)
	}
	if task.Completed {
		fmt.Println(color.GreenString("[✓] %s", label))
	} else if task.DueDate != "" && isOverdue(task.DueDate) {
		fmt.Println(color.RedString("[✗] %s", label))
	} else {
		fmt.Println(color.CyanString("[ ] %s", label))
	}
}

// isOverdue marks a task as past due date
func isOverdue(dueDate string) bool {
	due, err := time.Parse("2006-01-02", dueDate)
	return err == nil && time.Now().After(due)
}

// MarkTaskDone marks a task as completed
func MarkTaskDone(input string) error {
	tasks, err := LoadTasks()
	if err != nil {
		return err
	}
	found := false
	id, idErr := strconv.Atoi(input)

	for i, task := range tasks {
		if (idErr == nil && task.ID == id) || task.Text == input {
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
	tasks, err := LoadTasks()
	if err != nil {
		return err
	}
	parsedDate, err := ParseNaturalDate(dueDate)
	if err != nil {
		return err
	}

	found := false
	id, idErr := strconv.Atoi(input)

	for i, task := range tasks {
		if (idErr == nil && task.ID == id) || task.Text == input {
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
	tasks, err := LoadTasks()
	if err != nil {
		return err
	}
	found := false
	id, idErr := strconv.Atoi(input)

	for i, task := range tasks {
		if (idErr == nil && task.ID == id) || task.Text == input {
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
	tasks, err := LoadTasks()
	if err != nil {
		return err
	}

	newTasks := []Task{}
	found := false
	id, idErr := strconv.Atoi(input)

	for _, task := range tasks {
		if (idErr == nil && task.ID == id) || task.Text == input {
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

func ClearTasks() error {
	return SaveTasks([]Task{})
}

// SearchTasks displays tasks that contain a keyword
func SearchTasks(keyword string) {
	tasks, err := LoadTasks()
	if err != nil {
		color.Red("Failed to load tasks: %v", err)
		return
	}

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
