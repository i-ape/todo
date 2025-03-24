// task.go
package main

import "fmt"

type Task struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
}

// Display task
func (t Task) String() string {
	status := "❌"
	if t.Completed {
		status = "✅"
	}
	return fmt.Sprintf("[%d] %s %s", t.ID, status, t.Text)
}
