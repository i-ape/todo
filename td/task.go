package todo

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
	"todo/config"

	"github.com/fatih/color"
)

// Task represents a single task with all attributes.
type Task struct {
	ID        int      `json:"id"`
	Text      string   `json:"text"`
	Completed bool     `json:"completed"`
	DueDate   string   `json:"due_date,omitempty"`
	Tags      []string `json:"tags,omitempty"`
	Priority  string   `json:"priority,omitempty"`
	Recurring string   `json:"recurring,omitempty"`
	Notes     string   `json:"notes,omitempty"`
	Due       string   `json:"due,omitempty"`
	ParentID  int      `json:"parent_id,omitempty"`
	Category  string   `json:"category,omitempty"`
	Sticky    bool     `json:"sticky"`
	Important bool     `json:"important"`
}

// ListFilterOptions defines filters for task lists.
type ListFilterOptions struct {
	ShowDone    bool
	ShowPending bool
	TodayOnly   bool
	OverdueOnly bool
	JSONOutput  bool
	Tags        string
	Priority    string
	Recurring   bool
	Sticky      bool
}

// FilterTasks filters tasks based on options.
func FilterTasks(tasks []Task, opts ListFilterOptions) []Task {
	var filtered []Task
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
		filtered = append(filtered, task)
	}
	// Sort based on config
	SortTasks(filtered, config.GetSortOrder())
	return filtered
}

// AddTaskWithDueDateAndRecurring adds a task with optional due date, recurrence, and priority.
func AddTaskWithDueDateAndRecurring(text, due, recurring, priority string) error {
	tasks, err := LoadTasks()
	if err != nil {
		return err
	}
	parsedDue := ""
	if due != "" {
		dt, err := ParseNaturalDate(due)
		if err != nil {
			return err
		}
		parsedDue = dt
	}
	// Parse tags from text
	tags := ParseTagsFromText(text)

	newTask := Task{
		ID:        NextTaskID(tasks),
		Text:      text,
		Completed: false,
		DueDate:   parsedDue,
		Recurring: recurring,
		Tags:      tags,
		Priority:  ParsePriority(priority),
	}
	tasks = append(tasks, newTask)
	return SaveTasks(tasks)
}

// ParseTagsFromText extracts @/# tags from task text.
func ParseTagsFromText(text string) []string {
	var tags []string
	words := strings.Fields(text)
	for _, word := range words {
		if strings.HasPrefix(word, "@") || strings.HasPrefix(word, "#") {
			tags = append(tags, strings.TrimPrefix(strings.TrimPrefix(word, "@"), "#"))
		}
	}
	return tags
}

// MarkTaskDone marks a task as completed and handles recurrence.
func MarkTaskDone(input string) error {
	tasks, err := LoadTasks()
	if err != nil {
		return err
	}
	found := false
	id, err := strconv.Atoi(input)
	if err != nil {
		return fmt.Errorf("invalid task ID: %s", input)
	}
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Completed = true
			found = true
			if task.Recurring != "" {
				newTask := tasks[i]
				newTask.ID = NextTaskID(tasks)
				newTask.Completed = false
				newTask.DueDate = CalculateNextDueDate(task.DueDate, task.Recurring)
				tasks = append(tasks, newTask)
			}
			break
		}
	}
	if !found {
		return fmt.Errorf("task not found")
	}
	return SaveTasks(tasks)
}

// DeleteTask removes a task by ID or text.
func DeleteTask(input string) error {
	tasks, err := LoadTasks()
	if err != nil {
		return err
	}
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

// EditTaskText updates a task's text.
func EditTaskText(idOrText, newText string) error {
	tasks, err := LoadTasks()
	if err != nil {
		return err
	}

	id, idErr := strconv.Atoi(idOrText)
	found := false

	for i := range tasks {
		if (idErr == nil && tasks[i].ID == id) || tasks[i].Text == idOrText {
			if strings.TrimSpace(newText) == "" {
				return fmt.Errorf("new text cannot be empty")
			}
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

// SearchTasks prints tasks that match the keyword.
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

// ClearTasks deletes all tasks.
func ClearTasks() error {
	return SaveTasks([]Task{})
}

// GetSubtasks returns all subtasks of a given task.
func GetSubtasks(tasks []Task, parentID int) []Task {
	var subs []Task
	for _, t := range tasks {
		if t.ParentID == parentID {
			subs = append(subs, t)
		}
	}
	return subs
}

// HasSubtasks returns true if the task has child tasks.
func HasSubtasks(tasks []Task, taskID int) bool {
	for _, t := range tasks {
		if t.ParentID == taskID {
			return true
		}
	}
	return false
}

// ListTasks displays all tasks with hierarchy and details
func ListTasks() {
	tasks, err := LoadTasks()
	if err != nil {
		color.Red("Error loading tasks: %v", err)
		return
	}
	if len(tasks) == 0 {
		color.Yellow("📭 No tasks available.")
		return
	}

	// Optional: Filter with defaults (e.g., show all)
	opts := ListFilterOptions{}
	filtered := FilterTasks(tasks, opts)

	// Build hierarchy: map parent ID to subtasks
	taskMap := make(map[int][]Task)
	var topLevel []Task
	for _, task := range filtered {
		if task.ParentID == 0 {
			topLevel = append(topLevel, task)
		} else {
			taskMap[task.ParentID] = append(taskMap[task.ParentID], task)
		}
	}

	// Sort top-level by ID (or customize)
	sort.Slice(topLevel, func(i, j int) bool {
		return topLevel[i].ID < topLevel[j].ID
	})

	for _, task := range topLevel {
		printTaskLine(task, "") // No indent for top-level
		// Print subtasks indented
		subtasks := taskMap[task.ID]
		sort.Slice(subtasks, func(i, j int) bool {
			return subtasks[i].ID < subtasks[j].ID
		})
		for _, sub := range subtasks {
			printTaskLine(sub, "  ↳ ") // Indent subtasks
		}
	}
}

// printTaskLine prints a single task line with details (helper for ListTasks)
func printTaskLine(task Task, indent string) {
	status := color.CyanString("[ ] %d: %s", task.ID, task.Text)
	if task.Completed {
		status = color.GreenString("[✓] %d: %s", task.ID, task.Text)
	}

	if task.DueDate != "" {
		status += color.MagentaString(" (Due: %s)", task.DueDate)
	}
	if len(task.Tags) > 0 {
		status += color.BlueString(" (Tags: %s)", strings.Join(task.Tags, ", "))
	}
	if task.Priority != "" {
		status += color.YellowString(" (Priority: %s)", task.Priority)
	}
	if task.Recurring != "" {
		status += color.CyanString(" (Recurring: %s)", task.Recurring)
	}
	if task.Notes != "" {
		notesTrunc := task.Notes
		if len(notesTrunc) > 50 {
			notesTrunc = notesTrunc[:47] + "..."
		}
		status += color.RedString(" (Notes: %s)", notesTrunc)
	}

	fmt.Println(indent + status)
}

// AddTaskWithDueDate is deprecated; use AddTaskWithDueDateAndRecurring instead.
func AddTaskWithDueDate(text, due string) error {
	return AddTaskWithDueDateAndRecurring(text, due, "", "")
}
