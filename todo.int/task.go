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
	ID           int    `json:"id"`
	Text         string `json:"text"`
	Completed    bool   `json:"completed"`
	DueDate      string `json:"due_date,omitempty"`
	Recurring    string `json:"recurring,omitempty"`
	Until        string `json:"until,omitempty"`        
	Count        int    `json:"count,omitempty"`        
	Priority     string `json:"priority,omitempty"`    
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
	tasks, _ := LoadTasks()
	parsedDue := ""
	recurring := ""
	until := ""
	count := 0
	priority := "medium" // default

	// Detect priority keywords
	textLower := strings.ToLower(text)
	if strings.Contains(textLower, "[high]") {
		priority = "high"
		text = strings.ReplaceAll(text, "[high]", "")
	} else if strings.Contains(textLower, "[low]") {
		priority = "low"
		text = strings.ReplaceAll(text, "[low]", "")
	}

	// Parse due date or recurring type
	if due != "" {
		dueLower := strings.ToLower(due)

		switch dueLower {
		case "daily", "weekly", "monthly", "yearly", "every monday", "every friday":
			recurring = dueLower
			// Optional: set a future limit
			until = time.Now().AddDate(0, 3, 0).Format("2006-01-02") // example: 3 months from now
			count = 0 // 0 = unlimited
		default:
			dt, err := ParseNaturalDate(due)
			if err != nil {
				return err
			}
			parsedDue = dt
		}
	}

	newTask := Task{
		ID:        len(tasks) + 1,
		Text:      strings.TrimSpace(text),
		Completed: false,
		DueDate:   parsedDue,
		Recurring: recurring,
		Until:     until,
		Count:     count,
		Priority:  priority,
	}

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
		label := fmt.Sprintf("%d: %s", task.ID, task.Text)
		if task.DueDate != "" {
			label += fmt.Sprintf(" (Due: %s)", task.DueDate)
		}
		if task.Recurring != "" {
			label += color.MagentaString(" (Repeats: %s)", task.Recurring)
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
	tasks, _ := LoadTasks()
	found := false

	id, err := strconv.Atoi(input)
	for i, task := range tasks {
		if (err == nil && task.ID == id) || task.Text == input {
			// Recognize recurring keywords
			lower := strings.ToLower(strings.TrimSpace(dueDate))
			switch lower {
			case "daily", "weekly", "monthly", "yearly",
				"every monday", "every friday":
				tasks[i].Recurring = lower
				tasks[i].DueDate = "" // clear if previously set
			default:
				parsedDate, err := ParseNaturalDate(dueDate)
				if err != nil {
					return err
				}
				tasks[i].DueDate = parsedDate
				tasks[i].Recurring = "" // reset if previously set
			}
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
func SelectTaskFzf(tasks []Task) (Task, error) {
	if len(tasks) == 0 {
		return Task{}, fmt.Errorf("no tasks available")
	}

	options := []string{}
	idMap := map[string]Task{}
	for _, t := range tasks {
		label := fmt.Sprintf("%d: %s", t.ID, t.Text)
		options = append(options, label)
		idMap[label] = t
	}

	cmd := exec.Command("fzf")
	cmd.Stdin = strings.NewReader(strings.Join(options, "\n"))
	out, err := cmd.Output()
	if err != nil {
		return Task{}, err
	}

	selected := strings.TrimSpace(string(out))
	return idMap[selected], nil
}
