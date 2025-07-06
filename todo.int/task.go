package todo

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

// Task struct represents a single task
type Task struct {
	ID        int      `json:"id"`
	Text      string   `json:"text"`
	Completed bool     `json:"completed"`
	DueDate   string   `json:"due_date,omitempty"`
	Tags      []string `json:"tags,omitempty"`
	Priority  string   `json:"priority,omitempty"`
	Recurring string   `json:"recurring,omitempty"`
	Notes     string   `json:"notes,omitempty"`
	ParentID  int      `json:"parent_id,omitempty"` // ID of parent task, 0 = top-level
}

// AddTaskWithDueDate adds a task with an optional due date
func AddTaskWithDueDate(text, due string) error {
	tasks, _ := LoadTasks()
	parsed := ""
	if due != "" {
		dt, err := parseNaturalDate(due)
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
		status := color.CyanString("[ ] %d: %s", task.ID, task.Text)
		if task.Completed {
			status = color.GreenString("[✓] %d: %s", task.ID, task.Text)
		}

		if task.DueDate != "" {
			status += color.MagentaString(" (Due: %s)", task.DueDate)
		}

		fmt.Println(status)
	}
}

// ListFilterOptions defines filters that can be applied to a task list
type ListFilterOptions struct {
	ShowDone    bool
	ShowPending bool
	TodayOnly   bool
	OverdueOnly bool
	JSONOutput  bool
	Tags        string
	Priority    string
}

// FilterAndSortTasks filters and returns tasks based on ListFilterOptions
func FilterAndSortTasks(tasks []Task, opts ListFilterOptions) []Task {
	var result []Task
	today := time.Now().Format("2006-01-02")

	for _, task := range tasks {
		if opts.ShowDone && !task.Completed {
			continue
		}
		if opts.ShowPending && task.Completed {
			continue
		}
		if opts.TodayOnly && task.DueDate != today {
			continue
		}
		if opts.OverdueOnly && (task.DueDate == "" || !IsOverdue(task.DueDate)) {
			continue
		}
		if opts.Tags != "" {
			found := false
			for _, tag := range task.Tags {
				if strings.EqualFold(tag, opts.Tags) {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		if opts.Priority != "" && !strings.EqualFold(task.Priority, opts.Priority) {
			continue
		}

		result = append(result, task)
	}

	return result
}

// task.go
func FilterTasks(tasks []Task, options ListFilterOptions) []Task {
	var filtered []Task
	today := time.Now().Format("2006-01-02")

	for _, task := range tasks {
		if options.ShowDone && !task.Completed {
			continue
		}
		if options.ShowPending && task.Completed {
			continue
		}
		if options.TodayOnly && task.DueDate != today {
			continue
		}
		if options.OverdueOnly && !IsOverdue(task.DueDate) {
			continue
		}
		filtered = append(filtered, task)
	}
	return filtered
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

// EditTaskText updates a task's text
func EditTaskText(idOrText, newText string) error {
	tasks, err := LoadTasks()
	if err != nil {
		return err
	}

	id, idErr := strconv.Atoi(idOrText)
	found := false

	for i := range tasks {
		if (idErr == nil && tasks[i].ID == id) || tasks[i].Text == idOrText {
			tasks[i].Text = strings.TrimSpace(newText)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("task not found")
	}

	return SaveTasks(tasks)
}

// SearchTasks prints tasks that match the keyword
func SearchTasks(keyword string) {
	tasks, err := LoadTasks()
	if err != nil {
		fmt.Println("Error loading tasks:", err)
		return
	}
	found := false
	for _, task := range tasks {
		if strings.Contains(strings.ToLower(task.Text), strings.ToLower(keyword)) {
			fmt.Printf("🔍 %d: %s\n", task.ID, task.Text)
			found = true
		}
	}
	if !found {
		fmt.Println("No matching tasks found.")
	}
}

// ClearTasks deletes all tasks
func ClearTasks() error {
	return SaveTasks([]Task{})
}

// GetSubtasks returns all subtasks of a given task
func GetSubtasks(tasks []Task, parentID int) []Task {
	var subs []Task
	for _, t := range tasks {
		if t.ParentID == parentID {
			subs = append(subs, t)
		}
	}
	return subs
}

// HasSubtasks returns true if the task has child tasks
func HasSubtasks(tasks []Task, taskID int) bool {
	for _, t := range tasks {
		if t.ParentID == taskID {
			return true
		}
	}
	return false
}
