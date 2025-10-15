
package main

import (
	"log"
	"os"
	"todo/config"
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

func main() {
	HandleCommands()
}
