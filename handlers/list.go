package handlers

import (
	"encoding/json"
	"fmt"
	"os"
	"todo/config"
	td "todo/td"
)

func List() {
	path := config.GetTaskFilePath()
	data, err := os.ReadFile(path)
	if err != nil {
		warn("No tasks found.")
		return
	}

	var tasks []td.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		fail("Failed to parse tasks: %v", err)
		return
	}

	if len(tasks) == 0 {
		warn("No tasks in your list.")
		return
	}

	fmt.Println("🗒️  Tasks:")
	for _, t := range tasks {
		status := " "
		if t.Completed {
			status = "x"
		}
		fmt.Printf("[%s] #%d %s\n", status, t.ID, t.Text)
	}
}
