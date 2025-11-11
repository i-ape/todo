package handlers

import (
	"fmt"
	"os"
	todo "todo/td"
)

func Done() {
	if len(os.Args) < 3 {
		fmt.Println("❌ Usage: todo done [task id or text]")
		return
	}
	arg := os.Args[2]
	if err := todo.MarkTaskDone(arg); err != nil {
		fmt.Println("❌", err)
		return
	}
	fmt.Printf("✅ Marked done: %s\n", arg)
}
