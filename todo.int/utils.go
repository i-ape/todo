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

func parseFlags(args []string) (command string, commandArgs []string, flags map[string]string) {
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
		} else if command == "" {
			command = arg
		} else {
			commandArgs = append(commandArgs, arg)
		}
	}
	return
}