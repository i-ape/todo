package handlers

import (
	"fmt"
	"os"
)

func Reset() {
	if err := os.Remove("todo/tasks.json"); err != nil {
		fmt.Println("❌ Failed to reset:", err)
		return
	}
	fmt.Println("🧨 Reset todo list (deleted tasks.json).")
}
