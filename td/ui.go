// todo.int/ui.go
package todo

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

// PrintTask renders a single top-level task
func PrintTask(task Task, showNotes bool) {
	status := "[ ]"
	if task.Completed {
		status = "[✓]"
	} else if IsOverdue(task.DueDate) {
		status = "[✗]"
	}
	line := fmt.Sprintf("%s %d: %s", status, task.ID, task.Text)
	if task.DueDate != "" {
		line += color.MagentaString(" (Due: %s)", task.DueDate)
	}
	if len(task.Tags) > 0 {
		line += " [" + strings.Join(task.Tags, ", ") + "]"
	}
	fmt.Println(color.CyanString(line))
	if showNotes && task.Notes != "" {
		fmt.Println("   📝", task.Notes)
	}
}

// PrintTaskIndented renders a subtask
func PrintTaskIndented(task Task, showNotes bool) {
	status := "   [ ]"
	if task.Completed {
		status = "   [✓]"
	} else if IsOverdue(task.DueDate) {
		status = "   [✗]"
	}
	line := fmt.Sprintf("%s %d: %s", status, task.ID, task.Text)
	if task.DueDate != "" {
		line += color.MagentaString(" (Due: %s)", task.DueDate)
	}
	fmt.Println(color.CyanString(line))
	if showNotes && task.Notes != "" {
		fmt.Println("      📝", task.Notes)
	}
}
