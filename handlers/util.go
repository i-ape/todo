package handlers

import (
	"fmt"
	"os"
	"strconv"
	td"todo/td"
)

func parseID(s string) (int, error) {
	return strconv.Atoi(s)
}

func containsArg(target string) bool {
	for _, a := range os.Args {
		if a == target {
			return true
		}
	}
	return false
}

func getTaskByInput(input string) (*td.Task, error) {
	id, err := strconv.Atoi(input)
	if err == nil {
		return td.GetTaskByID(id)
	}
	tasks, err := td.LoadTasks()
	if err != nil {
		return nil, err
	}
	for _, t := range tasks {
		if t.Text == input {
			return &t, nil
		}
	}
	return nil, fmt.Errorf("no match for input: %s", input)
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "❌ "+format+"\n", args...)
}

func warn(format string, args ...any) {
	fmt.Fprintf(os.Stdout, "⚠️  "+format+"\n", args...)
}

func success(format string, args ...any) {
	fmt.Fprintf(os.Stdout, format+"\n", args...)
}