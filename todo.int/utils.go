package todo

import (
	"fmt"
	"strings"
)

func promptInput(prompt string, current string) string {
	fmt.Printf("%s [%s]: ", prompt, current)
	var input string
	fmt.Scanln(&input)
	return strings.TrimSpace(input)
}

func parseTags(input string) []string {
	parts := strings.Split(input, ",")
	var tags []string
	for _, tag := range parts {
		trimmed := strings.TrimSpace(tag)
		if trimmed != "" {
			tags = append(tags, trimmed)
		}
	}
	return tags
}
func selectSingleTask() (Task, error) {
	tasks, err := SelectTasksWithFzf(false)
	if err != nil || len(tasks) == 0 {
		return Task{}, err
	}
	return tasks[0], nil
}