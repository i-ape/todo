package main

import (
	// "bufio"
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

		case " " , "enter":
			m.tasks[m.cursor].Completed = !m.tasks[m.cursor].Completed
			_ = todo.SaveTasks(m.tasks)

		case "x", "backspace":
			m.tasks = append(m.tasks[:m.cursor], m.tasks[m.cursor+1:]...)
			if m.cursor >= len(m.tasks) && m.cursor > 0 {
				m.cursor--
			}
			_ = todo.SaveTasks(m.tasks)

		case "d":
			newDue, ok := prompt("📅 Enter new due date:")
			if ok {
				parsed, err := todo.ParseNaturalDate(newDue)
				if err == nil {
					m.tasks[m.cursor].DueDate = parsed
					_ = todo.SaveTasks(m.tasks)
				}
			}

		case "e":
			newText, ok := prompt("✏️ Edit task text:")
			if ok && strings.TrimSpace(newText) != "" {
				m.tasks[m.cursor].Text = strings.TrimSpace(newText)
				_ = todo.SaveTasks(m.tasks)
			}

		case "n":
			newTask, ok := prompt("➕ New task:")
			if ok && strings.TrimSpace(newTask) != "" {
				task := todo.Task{
					ID:   nextID(m.tasks),
					Text: strings.TrimSpace(newTask),
				}
				m.tasks = append(m.tasks, task)
				_ = todo.SaveTasks(m.tasks)
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
		label := task.Text

		if task.DueDate != "" {
			label += color.YellowString(" 📅 %s", task.DueDate)
		}
		if len(task.Tags) > 0 {
			label += " 🏷️ " + strings.Join(task.Tags, ", ")
		}
		if task.Priority == "high" {
			label += color.RedString(" 🔥")
		} else if task.Priority == "low" {
			label += color.BlueString(" ⬇")
		}

		b.WriteString(fmt.Sprintf("%s %s %s\n", cursor, status, label))
	}
	b.WriteString("\n↑/↓ or j/k to navigate, [n] new task, [enter] toggle complete, [q] quit\n")
	return b.String()
}

// Prompt Model

type promptModel struct {
	prompt  string
	value   string
	confirm bool
	cancel  bool
}

func (p promptModel) Init() tea.Cmd { return nil }

func (p promptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			p.confirm = true
			return p, tea.Quit
		case "esc":
			p.cancel = true
			return p, tea.Quit
		case "backspace":
			if len(p.value) > 0 {
				p.value = p.value[:len(p.value)-1]
			}
		default:
			p.value += msg.String()
		}
	}
	return p, nil
}

func (p promptModel) View() string {
	return fmt.Sprintf("\n%s\n> %s", p.prompt, p.value)
}

func prompt(promptText string) (string, bool) {
	pm := promptModel{prompt: promptText}
	p := tea.NewProgram(pm)
	m, err := p.Run()
	if err != nil {
		return "", false
	}
	final := m.(promptModel)
	return strings.TrimSpace(final.value), final.confirm && !final.cancel
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
func nextID(tasks []todo.Task) int {
	max := 0
	for _, t := range tasks {
		if t.ID > max {
			max = t.ID
		}
	}
	return max + 1
}