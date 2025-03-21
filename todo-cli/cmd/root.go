package cmd

import (
	"os"
	"todo-cli/todo"


	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "todo",
	Short: "A simple CLI to-do list",
	Long:  "A command-line to-do list application written in Go using Cobra.",
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}