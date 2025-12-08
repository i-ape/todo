// fzf.go
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	td "todo/td"
)

// --- FZF Selector ---
func SelectTasksWithFzf(multi bool, disable bool) ([]td.Task, error) {
	tasks, err := td.LoadTasks()
	if err != nil {
		return nil, fmt.Errorf("failed to load tasks: %w", err)
	}

	// If not disabled and fzf exists, try FZF
	if !disable {
		if _, err := exec.LookPath("fzf"); err == nil {
			if multi {
				return SelectMultipleTasksFzf(tasks)
			}
			task, err := SelectTaskFzf(tasks)
			if err != nil {
				return nil, err
			}
			return []td.Task{task}, nil
		}
	}

	// Fallback manual selection
	fmt.Println("FZF disabled or not found. Manual selection:")
	for _, t := range tasks {
		fmt.Printf("%d: %s\n", t.ID, t.Text)
	}
	fmt.Print("> Enter task ID: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	id, err := strconv.Atoi(input)
	if err != nil {
		return nil, fmt.Errorf("invalid ID")
	}

	for _, t := range tasks {
		if t.ID == id {
			return []td.Task{t}, nil
		}
	}
	return nil, fmt.Errorf("task not found")
}

// SelectTaskFzf allows user to choose a single task
func SelectTaskFzf(tasks []td.Task) (td.Task, error) {
	if _, err := exec.LookPath("fzf"); err != nil {
		return td.Task{}, fmt.Errorf("fzf not found")
	}
	opts := []string{}
	ref := map[string]td.Task{}
	for _, t := range tasks {
		label := fmt.Sprintf("%d: %s", t.ID, t.Text)
		opts = append(opts, label)
		ref[label] = t
	}
	cmd := exec.Command("fzf")
	cmd.Stdin = strings.NewReader(strings.Join(opts, "\n"))
	out, err := cmd.Output()
	if err != nil {
		return td.Task{}, fmt.Errorf("fzf error: %w", err)
	}
	choice := strings.TrimSpace(string(out))
	task, ok := ref[choice]
	if !ok {
		return td.Task{}, fmt.Errorf("invalid selection")
	}
	return task, nil
}

// SelectMultipleTasksFzf allows multiple task selection
func SelectMultipleTasksFzf(tasks []td.Task) ([]td.Task, error) {
	if _, err := exec.LookPath("fzf"); err != nil {
		return nil, fmt.Errorf("fzf not found")
	}
	opts := []string{}
	ref := map[string]td.Task{}
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
	var result []td.Task
	for _, l := range lines {
		if task, ok := ref[l]; ok {
			result = append(result, task)
		}
	}
	return result, nil
}
