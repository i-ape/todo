package handlers

import (
	"fmt"
	"os"
	todo "todo/td"
)

func Delete() {
	if len(os.Args) < 3 {
		fmt.Println("❌ Usage: todo delete [task id or text]")
		return
	}
	arg := os.Args[2]
	if err := todo.DeleteTask(arg); err != nil {
		fmt.Println("❌", err)
		return
	}
	fmt.Printf("🗑️ Deleted task: %s\n", arg)
}
