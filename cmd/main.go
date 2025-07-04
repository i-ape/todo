package main

import (
	"os"
	"todo/config"
)

func init() {
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
}

func main() {
	HandleCommands()
}
