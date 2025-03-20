package cmd

import (
	"fmt"
	"todo-cli"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	Run: func(cmd *cobra.Command, args []string) {
		todo_cli.ListTasks()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
