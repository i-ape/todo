package handlers

import (
	"fmt"
	todo "todo/td"
)

func Clear() {
	if err := todo.ClearTasks(); err != nil {
		fmt.Println("❌ Failed to clear tasks:", err)
		return
	}
	fmt.Println("🧹 Cleared all tasks.")
}
