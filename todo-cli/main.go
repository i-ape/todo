package main

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "todo",
	Short: "A simple CLI to-do list",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to the Go To-Do List! Use 'todo --help' to see available commands.")
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
