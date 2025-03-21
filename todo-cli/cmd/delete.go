package cmd

import (
	"fmt"
	"strconv"
	"todo-cli/todo"  // ✅ Import correct package

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [task ID]",
	Short: "Delete a task",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Invalid task ID:", args[0])
			return
		}
		err = todo.DeleteTask(id)  // ✅ Use correct function call
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Printf("❌ Task %d deleted!\n", id)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
