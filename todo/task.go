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

// 📌 Centralized abbreviation map
var abbreviationMap = map[string]func(time.Time) string{
	"td":      func(t time.Time) string { return t.Format("2006-01-02") },
	"tdy":     func(t time.Time) string { return t.Format("2006-01-02") },
	"today":   func(t time.Time) string { return t.Format("2006-01-02") },
	"tm":      func(t time.Time) string { return t.AddDate(0, 0, 1).Format("2006-01-02") },
	"tmmrw":   func(t time.Time) string { return t.AddDate(0, 0, 1).Format("2006-01-02") },
	"af":      func(t time.Time) string { return t.AddDate(0, 0, 2).Format("2006-01-02") },
	"aft":     func(t time.Time) string { return t.AddDate(0, 0, 2).Format("2006-01-02") },
	"yd":      func(t time.Time) string { return t.AddDate(0, 0, -1).Format("2006-01-02") },
	"yst":     func(t time.Time) string { return t.AddDate(0, 0, -1).Format("2006-01-02") },
	"nw":      func(t time.Time) string { return t.AddDate(0, 0, 7).Format("2006-01-02") },
	"nxtwk":   func(t time.Time) string { return t.AddDate(0, 0, 7).Format("2006-01-02") },
	"n2w":     func(t time.Time) string { return t.AddDate(0, 0, 14).Format("2006-01-02") },
	"n3w":     func(t time.Time) string { return t.AddDate(0, 0, 21).Format("2006-01-02") },
	"eod":     func(t time.Time) string { return t.Format("2006-01-02") },
	"someday": func(t time.Time) string { return "" },
	"soon":    func(t time.Time) string { return t.AddDate(0, 0, 3).Format("2006-01-02") },
	"nm":      func(t time.Time) string { return t.AddDate(0, 1, 0).Format("2006-01-02") },
	"em": func(t time.Time) string {
		nm := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location())
		return nm.AddDate(0, 0, -1).Format("2006-01-02")
	},
	"ew": func(t time.Time) string {
		return t.AddDate(0, 0, 7-int(t.Weekday())).Format("2006-01-02")
	},
	"nxtmon": func(t time.Time) string { return t.AddDate(0, 0, weekdayOffset(t, time.Monday)).Format("2006-01-02") },
	"nxfri":  func(t time.Time) string { return t.AddDate(0, 0, weekdayOffset(t, time.Friday)).Format("2006-01-02") },
	"mon":    func(t time.Time) string { return t.AddDate(0, 0, weekdayOffset(t, time.Monday)).Format("2006-01-02") },
	"tue":    func(t time.Time) string { return t.AddDate(0, 0, weekdayOffset(t, time.Tuesday)).Format("2006-01-02") },
	"wed": func(t time.Time) string {
		return t.AddDate(0, 0, weekdayOffset(t, time.Wednesday)).Format("2006-01-02")
	},
	"thu": func(t time.Time) string { return t.AddDate(0, 0, weekdayOffset(t, time.Thursday)).Format("2006-01-02") },
	"fri": func(t time.Time) string { return t.AddDate(0, 0, weekdayOffset(t, time.Friday)).Format("2006-01-02") },
	"sat": func(t time.Time) string { return t.AddDate(0, 0, weekdayOffset(t, time.Saturday)).Format("2006-01-02") },
	"sun": func(t time.Time) string { return t.AddDate(0, 0, weekdayOffset(t, time.Sunday)).Format("2006-01-02") },
}

// 📌 Calculate offset to next weekday
func weekdayOffset(t time.Time, target time.Weekday) int {
	offset := (int(target) - int(t.Weekday()) + 7) % 7
	if offset == 0 {
		offset = 7
	}
	return offset
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

// Parse natural date or fallback to standard formats
func parseNaturalDate(input string) (string, error) {
	input = strings.ToLower(strings.TrimSpace(input))
	today := time.Now()

	if f, ok := abbreviationMap[input]; ok {
		return f(today), nil
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
		switch parts[1] {
		case "d", "day", "days":
			return today.AddDate(0, 0, num).Format("2006-01-02"), nil
		case "w", "week", "weeks":
			return today.AddDate(0, 0, 7*num).Format("2006-01-02"), nil
		case "m", "month", "months":
			return today.AddDate(0, num, 0).Format("2006-01-02"), nil
		default:
			return "", fmt.Errorf("unsupported unit: %s", parts[1])
		}
	}

	for _, format := range []string{"2006-01-02", "02-01-2006"} {
		if t, err := time.Parse(format, input); err == nil {
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
