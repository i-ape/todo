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
