package todo

import (
	"fmt"
	"os/exec"
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

type ListFilterOptions struct {
	ShowDone    bool
	ShowPending bool
	TodayOnly   bool
	OverdueOnly bool
	JSONOutput  bool
	Tag         string
	Priority    string
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

// EditTaskText updates a task's text
func EditTaskText(idOrText, newText string) error {
	tasks, err := LoadTasks()
	if err != nil {
		return err
	}
	id, idErr := strconv.Atoi(idOrText)
	updated := false
	for i := range tasks {
		if (idErr == nil && tasks[i].ID == id) || tasks[i].Text == idOrText {
			tasks[i].Text = newText
			updated = true
			break
		}
	}
	if !updated {
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

// SelectTaskFzf allows user to choose a single task
func SelectTaskFzf(tasks []Task) (Task, error) {
	if _, err := exec.LookPath("fzf"); err != nil {
		return Task{}, fmt.Errorf("fzf not found")
	}
	opts := []string{}
	ref := map[string]Task{}
	for _, t := range tasks {
		label := fmt.Sprintf("%d: %s", t.ID, t.Text)
		opts = append(opts, label)
		ref[label] = t
	}
	cmd := exec.Command("fzf")
	cmd.Stdin = strings.NewReader(strings.Join(opts, "\n"))
	out, err := cmd.Output()
	if err != nil {
		return Task{}, fmt.Errorf("fzf error: %w", err)
	}
	choice := strings.TrimSpace(string(out))
	task, ok := ref[choice]
	if !ok {
		return Task{}, fmt.Errorf("invalid selection")
	}
	return task, nil
}

// SelectMultipleTasksFzf allows multiple task selection
func SelectMultipleTasksFzf(tasks []Task) ([]Task, error) {
	if _, err := exec.LookPath("fzf"); err != nil {
		return nil, fmt.Errorf("fzf not found")
	}
	opts := []string{}
	ref := map[string]Task{}
	for _, t := range tasks {
		label := fmt.Sprintf("%d: %s", t.ID, t.Text)
		opts = append(opts, label)
		ref[label] = t
	}
	cmd := exec.Command("fzf", "--multi")
	cmd.Stdin = strings.NewReader(strings.Join(opts, "\n"))
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("fzf error: %w", err)
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var result []Task
	for _, l := range lines {
		if task, ok := ref[l]; ok {
			result = append(result, task)
		}
	}
	return result, nil
}
