package todo

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	tasks []todo.Task
	index int
}

func (m model) Init() tea.Cmd { return nil }

func (m model) View() string {
	s := "🧃 TUI Task Picker:\n\n"
	for i, t := range m.tasks {
		cursor := " "
		if m.index == i {
			cursor = "👉"
		}
		s += fmt.Sprintf("%s %d: %s\n", cursor, t.ID, t.Text)
	}
	s += "\n(q to quit, ↑↓ to move, enter to select)"
	return s
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.index > 0 {
				m.index--
			}
		case "down", "j":
			if m.index < len(m.tasks)-1 {
				m.index++
			}
		case "enter":
			fmt.Println(m.tasks[m.index])
			return m, tea.Quit
		}
	}
	return m, nil
}

func StartTUI() {
	tasks, _ := todo.LoadTasks()
	p := tea.NewProgram(model{tasks: tasks})
	p.Start()
}

func StartTUISelect(multi bool) ([]todo.Task, error) {
	StartTUI()
	return nil, fmt.Errorf("TUI selection not yet implemented")
}
