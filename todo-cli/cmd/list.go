package cmd

import (
	"todo-cli/todo" 

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	Run: func(cmd *cobra.Command, args []string) {
		todo.ListTasks() 
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
