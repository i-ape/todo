package handlers

import "fmt"

func Help() {
	fmt.Print(`
📝 Usage:
  todo add [text] [due?]       → Add new task
  todo list                    → List all tasks
  todo done [id]               → Mark task done
  todo delete [id]             → Delete task
  todo edit [id]               → Edit task text
  todo search [keyword]        → Search tasks
  todo clear                   → Clear all tasks
  todo reset                   → Delete tasks.json
  todo help                    → Show this help

💡 Flags:
  --no-fzf    → Disable FZF interactive mode
  --tui       → Launch Bubble Tea interface

Aliases:
  a, ls, d, rm, clr, r, s, del, h, ?, -h, --help
`)
}
