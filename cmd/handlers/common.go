package handlers

import (
	"fmt"
	"os"
	"strconv"
	todo "todo/td"
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

func getTaskByInput(input string) (*todo.Task, error) {
	id, err := strconv.Atoi(input)
	if err == nil {
		return todo.GetTaskByID(id)
	}
	tasks, err := todo.LoadTasks()
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
