package cmd

import (
	"todo-cli/todo"  // ✅ Ensure this matches your `go.mod`

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	Run: func(cmd *cobra.Command, args []string) {
		todo.ListTasks()  // ✅ Call functions using `todo.`
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
