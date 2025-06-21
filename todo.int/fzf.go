// fzf.go
package todo

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func SelectTasksWithFzf(multi bool) ([]Task, error) {
	tasks, err := LoadTasks()
	if err != nil {
		return nil, fmt.Errorf("failed to load tasks: %w", err)
	}

	if _, err := exec.LookPath("fzf"); err == nil {
		if multi {
			return SelectMultipleTasksFzf(tasks)
		}
		task, err := SelectTaskFzf(tasks)
		if err != nil {
			return nil, err
		}
		return []Task{task}, nil
	}

	fmt.Println("FZF not found, falling back to manual selection:")
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
			return []Task{t}, nil
		}
	}
	return nil, fmt.Errorf("task not found")
}
