package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	todo "todo/todo.int"
)

// --- Task Management Functions ---

func AddTask(text, due string) error {
	return todo.AddTaskWithDueDate(text, due)
}

func handleList() {
	tags := os.Args[2:] // all arguments after `list` = filters
	todo.ListTasks(tags...)
}


func MarkTaskDone(input string) error {
	return todo.MarkTaskDone(input)
}

func SetDueDate(input, dueDate string) error {
	return todo.SetDueDate(input, dueDate)
}

func DeleteTask(input string) error {
	return todo.DeleteTask(input)
}

func ClearTasks() error {
	return todo.ClearTasks()
}

func ResetTasks() error {
	return os.Remove("todo/tasks.json")
}

func SearchTasks(keyword string) {
	todo.SearchTasks(keyword)
}

// --- CLI Command Dispatcher ---

func HandleCommands() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	aliases := map[string]string{
		"a":      "add",
		"ls":     "list",
		"d":      "done",
		"rm":     "delete",
		"del":    "delete",
		"clr":    "clear",
		"r":      "reset",
		"s":      "search",
		"h":      "help",
		"?":      "help",
		"-h":     "help",
		"--help": "help",
	}

	cmd := strings.ToLower(os.Args[1])
	if real, ok := aliases[cmd]; ok {
		cmd = real
	}

	switch cmd {
	case "add":
		handleAdd()
	case "edit":
		handleEdit()
	case "list":
		handleList()
	case "done":
		handleDone()
	case "due":
		handleDue()
	case "delete":
		handleDelete()
	case "clear":
		handleClear()
	case "reset":
		handleReset()
	case "search":
		handleSearch()
	case "tag":
		handleTag()
	case "help":
		printHelp()
	default:
		fmt.Println("❌ Unknown command:", cmd)
		printHelp()
	}
}

// --- FZF Selector ---

func selectTasksWithFzf(multi bool) ([]todo.Task, error) {
	tasks, err := todo.LoadTasks()
	if err != nil {
		return nil, fmt.Errorf("failed to load tasks: %w", err)
	}

	options := []string{}
	idMap := map[string]todo.Task{}
	for _, t := range tasks {
		label := fmt.Sprintf("%d: %s", t.ID, t.Text)
		options = append(options, label)
		idMap[label] = t
	}

	args := []string{}
	if multi {
		args = append(args, "--multi")
	}

	cmd := exec.Command("fzf", args...)
	cmd.Stdin = strings.NewReader(strings.Join(options, "\n"))
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("fzf error: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var selected []todo.Task
	for _, line := range lines {
		if task, ok := idMap[line]; ok {
			selected = append(selected, task)
		}
	}

	return selected, nil
}

// --- Handlers ---

func handleAdd() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: todo add [task text] [optional due date]")
		return
	}
	text := os.Args[2]
	due := ""
	if len(os.Args) > 3 {
		due = strings.Join(os.Args[3:], " ")
	}
	if err := AddTask(text, due); err != nil {
		fmt.Println("Error:", err)
	}
}

func handleEdit() {
	selected, err := selectTasksWithFzf(false)
	if err != nil || len(selected) == 0 {
		fmt.Println("Select error:", err)
		return
	}
	task := selected[0]

	fmt.Printf("✏️  Editing: %s\n> ", task.Text)
	reader := bufio.NewReader(os.Stdin)
	newText, _ := reader.ReadString('\n')
	newText = strings.TrimSpace(newText)
	if newText == "" {
		fmt.Println("No changes made.")
		return
	}
	if err := todo.EditTaskText(strconv.Itoa(task.ID), newText); err != nil {
		fmt.Println("Edit error:", err)
	}
}

func handleDone() {
	selected, err := selectTasksWithFzf(true)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for _, task := range selected {
		if err := MarkTaskDone(strconv.Itoa(task.ID)); err != nil {
			fmt.Println("❌", err)
		}
	}
}

func handleDelete() {
	selected, err := selectTasksWithFzf(true)
	if err != nil {
		fmt.Println("Error selecting task:", err)
		return
	}
	for _, task := range selected {
		if err := todo.DeleteTask(strconv.Itoa(task.ID)); err != nil {
			fmt.Println("❌", err)
		}
	}
}

func handleDue() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: todo due [task ID or task text] [date]")
		return
	}
	input := os.Args[2]
	dueDate := strings.Join(os.Args[3:], " ")
	if err := SetDueDate(input, dueDate); err != nil {
		fmt.Println("Error:", err)
	}
}

func handleSearch() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: todo search [keyword]")
		return
	}
	SearchTasks(os.Args[2])
}
func handleTag() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: todo tag [id|text] [tags...] [--remove #tag1 #tag2]")
		return
	}

	target := os.Args[2]
	addTags := []string{}
	removeTags := []string{}

	// Parse args after target
	for i := 3; i < len(os.Args); i++ {
		arg := os.Args[i]
		if arg == "--remove" {
			removeTags = os.Args[i+1:]
			break
		}
		addTags = append(addTags, arg)
	}

	if err := todo.UpdateTags(target, addTags, removeTags); err != nil {
		fmt.Println("Error:", err)
	}
}


func handleClear() {
	if err := ClearTasks(); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("✅ All tasks cleared.")
	}
}

func handleReset() {
	if err := ResetTasks(); err != nil {
		fmt.Println("⚠️ Reset failed:", err)
	} else {
		fmt.Println("🗑️ tasks.json deleted.")
	}
}

// --- Help ---

func printHelp() {
	fmt.Println(`📝 Usage:
  todo add [text] [due?]       → Add new task
  todo list                    → List all tasks
  todo done                    → Mark one or more tasks done
  todo due [id|text] [date]    → Set/change due date
  todo delete                  → Delete one or more tasks
  todo edit                    → Edit a task
  todo search [keyword]        → Search task text
  todo clear                   → Clear all tasks
  todo reset                   → Delete tasks.json
  todo help                    → Show help`)
}
