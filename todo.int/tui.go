package todo

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	tasks []todo.Task
	cursor int
	quitting bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "j", "down":
			if m.cursor < len(m.tasks)-1 {
				m.cursor++
			}
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return "Goodbye 👋\n"
	}

	var b strings.Builder
	b.WriteString("📋 Tasks:\n\n")
	for i, task := range m.tasks {
		cursor := "  "
		if m.cursor == i {
			cursor = "👉"
		}
		status := "[ ]"
		if task.Completed {
			status = "[✓]"
		}
		b.WriteString(fmt.Sprintf("%s %s %s\n", cursor, status, task.Text))
	}
	b.WriteString("\n↑/↓ or j/k to navigate, q to quit\n")
	return b.String()
}

func StartTui() {
	tasks, err := todo.LoadTasks()
	if err != nil {
		fmt.Println("failed to load tasks:", err)
		os.Exit(1)
	}

	m := model{tasks: tasks}
	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		fmt.Println("TUI error:", err)
		os.Exit(1)
	}
}
