package main

import (
	"os"
)

func init() {
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		if arg == "--no-fzf" {
			disableFzf = true
			os.Args = append(os.Args[:i], os.Args[i+1:]...)
			i--
		}
		if arg == "--tui" {
			enableTui = true
			os.Args = append(os.Args[:i], os.Args[i+1:]...)
			i--
		}
	}
}

func main() {
	HandleCommands()
}
