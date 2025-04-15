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

// parseNaturalDate handles natural language date inputs
func parseNaturalDate(input string) (string, error) {
	input = strings.ToLower(strings.TrimSpace(input))
	today := time.Now()

	switch input {
	case "td", "tdy", "today":
		return today.Format("2006-01-02"), nil
	case "tm", "tmmrw", "tomorrow":
		return today.AddDate(0, 0, 1).Format("2006-01-02"), nil
	case "af", "aft", "after tomorrow":
		return today.AddDate(0, 0, 2).Format("2006-01-02"), nil
	case "yd", "yst", "yesterday":
		return today.AddDate(0, 0, -1).Format("2006-01-02"), nil
	case "nw", "nxtwk", "next week":
		return today.AddDate(0, 0, 7).Format("2006-01-02"), nil
	case "n2w":
		return today.AddDate(0, 0, 14).Format("2006-01-02"), nil
	case "n3w":
		return today.AddDate(0, 0, 21).Format("2006-01-02"), nil
	case "ew", "end of week":
		weekday := int(today.Weekday())
		daysUntilSunday := 7 - weekday
		return today.AddDate(0, 0, daysUntilSunday).Format("2006-01-02"), nil
	case "nm", "next month":
		return today.AddDate(0, 1, 0).Format("2006-01-02"), nil
	case "em", "end of month":
		firstOfNextMonth := time.Date(today.Year(), today.Month()+1, 1, 0, 0, 0, 0, today.Location())
		endOfMonth := firstOfNextMonth.AddDate(0, 0, -1)
		return endOfMonth.Format("2006-01-02"), nil
	case "nxtmon":
		offset := (int(time.Monday) - int(today.Weekday()) + 7) % 7
		if offset == 0 {
			offset = 7
		}
		return today.AddDate(0, 0, offset).Format("2006-01-02"), nil
	case "nxfri":
		offset := (int(time.Friday) - int(today.Weekday()) + 7) % 7
		if offset == 0 {
			offset = 7
		}
		return today.AddDate(0, 0, offset).Format("2006-01-02"), nil
	}
	

	if strings.HasPrefix(input, "in ") {
		parts := strings.Split(input[3:], " ")
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid relative date format: %s", input)
		}
		num, err := strconv.Atoi(parts[0])
		if err != nil {
			return "", fmt.Errorf("invalid number in relative date: %v", err)
		}
		unit := parts[1]
		switch unit {
		case "d", "day", "days":
			return today.AddDate(0, 0, num).Format("2006-01-02"), nil
		case "week", "weeks", "w":
			return today.AddDate(0, 0, 7*num).Format("2006-01-02"), nil
		case "month", "months", "m":
			return today.AddDate(0, num, 0).Format("2006-01-02"), nil
		default:
			return "", fmt.Errorf("unsupported time unit: %s", unit)
		}
	}

	// Try standard date formats
	formats := []string{"2006-01-02", "02-01-2006"}
	for _, layout := range formats {
		if t, err := time.Parse(layout, input); err == nil {
			return t.Format("2006-01-02"), nil
		}
	}

	return "", fmt.Errorf("could not parse date: %s", input)
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
