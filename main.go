package main

import (
	"fmt"
	"log"
	"os"
	"todo/cmd"
	"todo/config"
	///"todo/handlers"
)

func init() {
	if err := config.Load(); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		if arg == "--no-fzf" {
			config.DisableFzf = true
			os.Args = append(os.Args[:i], os.Args[i+1:]...)
			i--
		}
		if arg == "--tui" {
			config.EnableTui = true
			os.Args = append(os.Args[:i], os.Args[i+1:]...)
			i--
		}
	}
	///config.loadDotEnv(".env")

}

// Entry point
func main() {
	// Optional: print debug info about CLI args
	cmd.Debug("CLI args: %v", os.Args)

	// Simple argument handling
	switch arg(1) {
	case "help":
		fmt.Println("Usage: todo [command]")
		fmt.Println("Commands:")
		fmt.Println("  help   Show this help message")
		fmt.Println("  tui    Launch the interactive TUI")
		fmt.Println("  list   List tasks (CLI mode)")
		fmt.Println("  add    Add a new task (CLI mode)")
		return

	case "tui", "":
		// Default: launch the TUI
		cmd.StartTUI()
		return

	default:
		// Otherwise, handle CLI commands (like add, list, delete, etc.)
		cmd.HandleCommands()
	}
}
