package cmd

import (
	"fmt"

	"todo-cli/todo"  // ✅ Ensure this matches your module name from `go.mod`

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [task]",
	Short: "Add a new task",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		taskText := args[0]
		err := todo.AddTask(taskText)  // ✅ Call function from `todo/task.go`
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("✅ Task added:", taskText)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
