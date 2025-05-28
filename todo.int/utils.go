package todo

import (
	"fmt"
	"strconv"
	"strings"
	"time"
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
		if len(tokens) == 2 && tokens[0] == "for" {
			delta := tokens[1]
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
			}
			recurring = "custom"
		}
	}

	// Parse natural date
	d, err := ParseNaturalDate(strings.TrimSpace(main))
	return d, t, dur, recurring, until, err
}
