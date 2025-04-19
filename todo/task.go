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

// 📌 Centralized abbreviation map — natural language → date string
var abbreviationMap = map[string]func(time.Time) string{
	// 📅 Absolute terms
	"td":       formatToday,
	"tdy":      formatToday,
	"today":    formatToday,
	"tm":       inDays(1),
	"tmmrw":    inDays(1),
	"next":     inDays(1),
	"af":       inDays(2),
	"aft":      inDays(2),
	"yd":       inDays(-1),
	"yst":      inDays(-1),
	"soon":     inDays(3),
	"later":    inDays(7),
	"someday":  func(t time.Time) string { return "" },
	"now":      formatToday,

	// 📆 Weeks
	"nw":       inDays(7),
	"nxtwk":    inDays(7),
	"n2w":      inDays(14),
	"n3w":      inDays(21),
	"eowk":     nextWeekday(time.Friday),

	// 📅 Months
	"nm":       inMonths(1),
	"em":       endOfMonth,

	// 🕐 Time-based
	"eod":      formatToday, // End of day = today, could change later

	// 🗓️ Weekdays (next occurrence)
	"mon":      nextWeekday(time.Monday),
	"tue":      nextWeekday(time.Tuesday),
	"wed":      nextWeekday(time.Wednesday),
	"thu":      nextWeekday(time.Thursday),
	"fri":      nextWeekday(time.Friday),
	"sat":      nextWeekday(time.Saturday),
	"sun":      nextWeekday(time.Sunday),

	"nxtmon":   nextWeekday(time.Monday),
	"nxfri":    nextWeekday(time.Friday),

	// 📅 End of current week
	"ew": func(t time.Time) string {
		return t.AddDate(0, 0, 7-int(t.Weekday())).Format("2006-01-02")
	},
}

func formatToday(t time.Time) string {
	return t.Format("2006-01-02")
}

func inDays(n int) func(time.Time) string {
	return func(t time.Time) string {
		return t.AddDate(0, 0, n).Format("2006-01-02")
	}
}

func inMonths(n int) func(time.Time) string {
	return func(t time.Time) string {
		return t.AddDate(0, n, 0).Format("2006-01-02")
	}
}

func endOfMonth(t time.Time) string {
	firstNext := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location())
	return firstNext.AddDate(0, 0, -1).Format("2006-01-02")
}

func nextWeekday(wd time.Weekday) func(time.Time) string {
	return func(t time.Time) string {
		offset := (int(wd) - int(t.Weekday()) + 7) % 7
		if offset == 0 {
			offset = 7
		}
		return t.AddDate(0, 0, offset).Format("2006-01-02")
	}
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
