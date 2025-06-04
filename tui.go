package main

import (
	"fmt"
	"os"
	"strings"

	todo "todo/todo.int"

	tea "github.com/charmbracelet/bubbletea"
	color "github.com/fatih/color"
)

type model struct {
	tasks    []todo.Task
	cursor   int
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

		case "enter", " ":
			m.tasks[m.cursor].Completed = !m.tasks[m.cursor].Completed
			if err := todo.SaveTasks(m.tasks); err != nil {
				fmt.Println("❌ Failed to save task:", err)
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
			cursor = "▶" // or "▸", "▶", "›", "→", "➤"
		}
		status := color.CyanString("[ ]")
		if task.Completed {
			status = color.GreenString("[✓]")
		} else if task.DueDate != "" && todo.IsOverdue(task.DueDate) {
			status = color.RedString("[✗]")
		}
		b.WriteString(fmt.Sprintf("%s %s %s\n", cursor, status, task.Text))
	}
	b.WriteString("\n↑/↓ or j/k to navigate, q to quit\n")
	return b.String()
}

func StartTUI() {
	tasks, err := todo.LoadTasks()
	if err != nil {
		fmt.Println("Failed to load tasks:", err)
		os.Exit(1)
	}

	p := tea.NewProgram(model{tasks: tasks})
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running TUI:", err)
		os.Exit(1)
	}
}
