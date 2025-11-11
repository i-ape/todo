package handlers

import (
	"fmt"
	"os"
	todo "todo/td"
)

func Search() {
	if len(os.Args) < 3 {
		fmt.Println("❌ Usage: todo search [keyword]")
		return
	}
	todo.SearchTasks(os.Args[2])
}
