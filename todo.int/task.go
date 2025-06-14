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

// parseNaturalDate handles natural language date inputs
func parseNaturalDate(input string) (string, error) {
	now := time.Now()
	input = strings.ToLower(strings.TrimSpace(input))

	switch input {
	case "today":
		return now.Format("2006-01-02"), nil
	case "tomorrow":
		return now.AddDate(0, 0, 1).Format("2006-01-02"), nil
	case "next week":
		return now.AddDate(0, 0, 7).Format("2006-01-02"), nil
	case "next month":
		return now.AddDate(0, 1, 0).Format("2006-01-02"), nil
	case "next year":
		return now.AddDate(1, 0, 0).Format("2006-01-02"), nil
	default:
		// in N days or weeks
		if strings.HasPrefix(input, "in ") {
			parts := strings.Split(input, " ")
			if len(parts) == 3 {
				num, err := strconv.Atoi(parts[1])
				if err != nil {
					return "", fmt.Errorf("invalid number in relative date")
				}
				switch parts[2] {
				case "day", "days":
					return now.AddDate(0, 0, num).Format("2006-01-02"), nil
				case "week", "weeks":
					return now.AddDate(0, 0, num*7).Format("2006-01-02"), nil
				}
			}
		}

		// try DD-MM-YYYY
		t, err := time.Parse("02-01-2006", input)
		if err == nil {
			return t.Format("2006-01-02"), nil
		}
		// try YYYY-MM-DD
		t, err = time.Parse("2006-01-02", input)
		if err == nil {
			return t.Format("2006-01-02"), nil
		}
		return "", fmt.Errorf("invalid date format or unsupported natural keyword")
	}
}

// SetDueDate assigns a due date to a task
func SetDueDate(input string, dueDate string) error {
	tasks, _ := LoadTasks()
	found := false

	parsedDate, err := parseNaturalDate(dueDate)
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
