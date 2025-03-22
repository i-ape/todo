package cmd

import (
	"fmt"
	"strconv"
	"todo-cli/todo"  // ✅ Import task management functions

	"github.com/spf13/cobra"
)

var doneCmd = &cobra.Command{
	Use:   "done [task ID]",
	Short: "Mark a task as completed",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Invalid task ID:", args[0])
			return
		}
		err = todo.MarkTaskDone(id)  // ✅ Use function from `todo/task.go`
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Printf("✅ Task %d marked as done!\n", id)
		}
	},
}

func init() {
	rootCmd.AddCommand(doneCmd)
}
