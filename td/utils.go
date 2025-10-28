package todo

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

// PromptInput prompts the user for input, pre-filling with current value
func PromptInput(prompt string, current string) string {
	fmt.Printf("%s [%s]: ", prompt, current)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return current // Return current if empty
	}
	return input
}

// ParseTags parses comma-separated tags, handling @ or # prefixes
func ParseTags(input string) []string {
	parts := strings.Split(input, ",")
	var tags []string
	for _, tag := range parts {
		trimmed := strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(tag, "@"), "#"))
		if trimmed != "" {
			tags = append(tags, trimmed)
		}
	}
	return tags
}

// ConfirmPrompt prompts for yes/no confirmation, defaulting to no
func ConfirmPrompt(question string) bool {
	fmt.Printf("%s [y/N]: ", question)
	var input string
	fmt.Scanln(&input)
	input = strings.ToLower(strings.TrimSpace(input))
	return input == "y" || input == "yes"
}

// Slugify converts string to slug format
func Slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, " ", "-")
	return s
}

// JoinNonEmpty joins non-empty strings with separator
func JoinNonEmpty(elems []string, sep string) string {
	var result []string
	for _, e := range elems {
		if trimmed := strings.TrimSpace(e); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return strings.Join(result, sep)
}

// ParsePriority parses and normalizes priority string
func ParsePriority(s string) string {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "low":
		return "low"
	case "high":
		return "high"
	default:
		return "medium"
	}
}

// ParseFlags parses command-line flags and args
func ParseFlags(args []string) (command string, commandArgs []string, flags map[string]string) {
	flags = make(map[string]string)
	command = ""
	commandArgs = []string{}

	for _, arg := range args {
		if strings.HasPrefix(arg, "--") {
			if strings.Contains(arg, "=") {
				parts := strings.SplitN(arg[2:], "=", 2)
				flags[parts[0]] = parts[1]
			} else {
				flags[arg[2:]] = "true"
			}
		} else if strings.HasPrefix(arg, "-") { // Handle short flags, e.g., -j for --json
			flags[arg[1:]] = "true"
		} else if command == "" {
			command = arg
		} else {
			commandArgs = append(commandArgs, arg)
		}
	}
	return
}

// ParseDateTimeDurationRepeat parses strings like "friday @ 14:00 for 45m every weekday"
func ParseDateTimeDurationRepeat(input string) (date, t, dur, recurring, until string, err error) {
	input = strings.ToLower(input)
	main := input
	if idx := strings.Index(input, "@"); idx != -1 {
		main = input[:idx]
		tPart := strings.TrimSpace(input[idx+1:])
		tokens := strings.Fields(tPart)

		// Parse time
		if len(tokens) > 0 && strings.Contains(tokens[0], ":") {
			t = tokens[0]
			tokens = tokens[1:]
		}

		// Parse duration
		if len(tokens) > 1 && tokens[0] == "for" {
			dur = tokens[1]
			tokens = tokens[2:]
		}

		// Parse recurrence
		if len(tokens) > 0 && tokens[0] == "every" {
			recurring = strings.Join(tokens, " ")
			// Parse optional "for X duration" after recurrence
			if strings.Contains(recurring, "for") {
				parts := strings.SplitN(recurring, "for", 2)
				recurring = strings.TrimSpace(parts[0])
				delta := strings.TrimSpace(parts[1])
				now := time.Now()
				switch {
				case strings.HasSuffix(delta, "d"):
					d, _ := strconv.Atoi(strings.TrimSuffix(delta, "d"))
					until = now.AddDate(0, 0, d).Format("2006-01-02")
				case strings.HasSuffix(delta, "w"):
					w, _ := strconv.Atoi(strings.TrimSuffix(delta, "w"))
					until = now.AddDate(0, 0, w*7).Format("2006-01-02")
				case strings.HasSuffix(delta, "m"):
					m, _ := strconv.Atoi(strings.TrimSuffix(delta, "m"))
					until = now.AddDate(0, m, 0).Format("2006-01-02")
				case strings.HasSuffix(delta, "weeks"):
					w, _ := strconv.Atoi(strings.TrimSuffix(delta, "weeks"))
					until = now.AddDate(0, 0, w*7).Format("2006-01-02")
					// Add more for "3weeks", etc.
				}
			}
		}
}

	// Parse natural date
	d, err := ParseNaturalDate(strings.TrimSpace(main))
	return d, t, dur, recurring, until, err
}

// UpdateTaskByID updates a task by ID using the provided function
func UpdateTaskByID(id int, updateFn func(*Task)) error {
	tasks, err := LoadTasks()
	if err != nil {
		return err
	}
	for i := range tasks {
		if tasks[i].ID == id {
			updateFn(&tasks[i])
			return SaveTasks(tasks)
		}
	}
	return fmt.Errorf("task not found")
}

// PromptMultiline prompts for multiline input using editor
func PromptMultiline(prompt, initial string) string {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi" // Fallback
	}
	tmpfile, err := os.CreateTemp("", "todo_note_*.txt")
	if err != nil {
		fmt.Println("Error creating temp file:", err)
		return initial
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(initial)); err != nil {
		fmt.Println("Error writing initial note:", err)
		return initial
	}
	tmpfile.Close()

	cmd := exec.Command(editor, tmpfile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error running editor:", err)
		return initial
	}

	data, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		fmt.Println("Error reading note:", err)
		return initial
	}
	return strings.TrimSpace(string(data))
}

// PrintTaskDetails prints detailed task info
func PrintTaskDetails(task Task) {
	fmt.Printf("🆔 ID: %d\n", task.ID)
	fmt.Printf("📝 Text: %s\n", task.Text)
	if task.DueDate != "" {
		fmt.Printf("📅 Due: %s\n", task.DueDate)
	}
	if task.Priority != "" {
		fmt.Printf("🔥 Priority: %s\n", task.Priority)
	}
	if len(task.Tags) > 0 {
		fmt.Printf("🏷️ Tags: %s\n", strings.Join(task.Tags, ", "))
	}
	if task.Recurring != "" {
		fmt.Printf("🔁 Recurring: %s\n", task.Recurring)
	}
	if task.Notes != "" {
		fmt.Printf("📝 Notes: %s\n", task.Notes)
	}
	fmt.Printf("✅ Completed: %v\n", task.Completed)
	if task.ParentID > 0 {
		fmt.Printf("📌 Parent ID: %d\n", task.ParentID)
	}
}

// RenderTask renders a task string for TUI/CLI
func RenderTask(task Task, isSelected bool) string {
	cursor := "  "
	if isSelected {
		cursor = "▶ "
	}
	status := color.CyanString("[ ]")
	if task.Completed {
		status = color.GreenString("[✓]")
	} else if task.DueDate != "" && IsOverdue(task.DueDate) {
		status = color.RedString("[✗]")
	}
	label := task.Text
	if task.Sticky {
		label = "📌 " + label
	}
	if task.DueDate != "" {
		label += color.YellowString(" 📅 %s", task.DueDate)
	}
	if len(task.Tags) > 0 {
		label += " 🏷️ " + strings.Join(task.Tags, ", ")
	}
	switch task.Priority {
	case "high":
		label += color.RedString(" 🔥")
	case "low":
		label += color.BlueString(" ⬇")
	}
	return fmt.Sprintf("%s %s %s", cursor, status, label)
}

// SortTasks sorts tasks by order
func SortTasks(tasks []Task, order string) {
	sort.Slice(tasks, func(i, j int) bool {
		switch order {
		case "due":
			return tasks[i].DueDate < tasks[j].DueDate
		case "priority":
			priorityOrder := map[string]int{"high": 3, "medium": 2, "low": 1}
			return priorityOrder[tasks[i].Priority] > priorityOrder[tasks[j].Priority]
		case "text":
			return tasks[i].Text < tasks[j].Text
		default:
			return tasks[i].ID < tasks[j].ID
		}
	})
}

// TruncateString truncates a string to maxLen with ellipsis
func TruncateString(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen-3] + "..."
	}
	return s
}

// NextTaskID calculates the next available task ID
func NextTaskID(tasks []Task) int {
	maxID := 0
	for _, task := range tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}
	return maxID + 1
}

// GetTaskByID retrieves a task by ID
func GetTaskByID(id int) (*Task, error) {
	tasks, err := LoadTasks()
	if err != nil {
		return nil, err
	}
	for _, task := range tasks {
		if task.ID == id {
			return &task, nil
		}
	}
	return nil, fmt.Errorf("task with ID %d not found", id)
}

// RemoveTaskByID removes a task with a matching ID from the slice.
func RemoveTaskByID(tasks []Task, id int) ([]Task, *Task) {
	var deleted *Task
	newTasks := make([]Task, 0, len(tasks))
	for _, t := range tasks {
		if t.ID == id {
			deleted = &t
			continue
		}
		newTasks = append(newTasks, t)
	}
	return newTasks, deleted
}
