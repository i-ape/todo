// cmd/tui.go
package main

import (
	"fmt"
	"os"
	"strings"
	todo "todo/td"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
)

type model struct {
	tasks    []todo.Task
	cursor   int
	quitting bool
	err      error // Added for error display
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
		case " ", "enter":
			if len(m.tasks) > 0 {
				m.tasks[m.cursor].Completed = !m.tasks[m.cursor].Completed
				_ = todo.SaveTasks(m.tasks)
			}
		case "x", "backspace":
			if len(m.tasks) > 0 {
				m.tasks = append(m.tasks[:m.cursor], m.tasks[m.cursor+1:]...)
				if m.cursor >= len(m.tasks) && m.cursor > 0 {
					m.cursor--
				}
				_ = todo.SaveTasks(m.tasks)
			}
		case "d":
			newDue, ok := prompt("📅 Enter new due date (e.g., today, 2025-12-31):")
			if ok && strings.TrimSpace(newDue) != "" {
				parsed, err := todo.ParseNaturalDate(newDue)
				if err != nil {
					m.err = err
				} else {
					m.tasks[m.cursor].DueDate = parsed
					m.err = nil
					_ = todo.SaveTasks(m.tasks)
				}
			}
		case "e":
			newText, ok := prompt("✏️ Edit task text:")
			if ok && strings.TrimSpace(newText) != "" {
				m.tasks[m.cursor].Text = strings.TrimSpace(newText)
				m.err = nil
				_ = todo.SaveTasks(m.tasks)
			}
		case "n":
			newTask, ok := prompt("➕ New task:")
			if ok && strings.TrimSpace(newTask) != "" {
				tasks, _ := todo.LoadTasks()
				task := todo.Task{
					ID:       todo.NextTaskID(tasks),
					Text:     strings.TrimSpace(newTask),
					Priority: todo.ParsePriority("medium"),
				}
				m.tasks = append(m.tasks, task)
				m.err = nil
				_ = todo.SaveTasks(m.tasks)
			}
		case "t":
			tags, ok := prompt("🏷️ Enter tags (comma-separated):")
			if ok && strings.TrimSpace(tags) != "" {
				m.tasks[m.cursor].Tags = todo.ParseTags(tags)
				m.err = nil
				_ = todo.SaveTasks(m.tasks)
			}
		case "p":
			priority, ok := prompt("🔥 Enter priority (high, medium, low):")
			if ok && strings.TrimSpace(priority) != "" {
				m.tasks[m.cursor].Priority = todo.ParsePriority(priority)
				m.err = nil
				_ = todo.SaveTasks(m.tasks)
			}
		case "s":
			text, ok := prompt("📌 Subtask text:")
			if ok && strings.TrimSpace(text) != "" {
				tasks, _ := todo.LoadTasks()
				subtask := todo.Task{
					ID:        todo.NextTaskID(tasks),
					Text:      strings.TrimSpace(text),
					ParentID:  m.tasks[m.cursor].ID,
					Completed: false,
					Priority:  todo.ParsePriority("medium"),
				}
				m.tasks = append(m.tasks, subtask)
				m.err = nil
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
	if m.err != nil {
		b.WriteString(color.RedString("Error: %v\n\n", m.err))
	}
	b.WriteString("📋 Tasks:\n\n")
	for i, task := range m.tasks {
		cursor := "  "
		if m.cursor == i {
			cursor = "▶"
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
		switch task.Priority {
		case "high":
			label += color.RedString(" 🔥")
		case "low":
			label += color.BlueString(" ⬇")
		}
		indent := ""
		if task.ParentID > 0 {
			indent = "  ↳ "
		}
		b.WriteString(fmt.Sprintf("%s %s %s%s\n", cursor, status, indent, label))
	}
	b.WriteString("\n↑/↓ or j/k: navigate, [enter]: toggle, [x]: delete, [n]: new, [d]: due date, [e]: edit, [t]: tags, [p]: priority, [s]: subtask, [q]: quit\n")
	return b.String()
}

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
		fmt.Fprintf(os.Stderr, "Failed to load tasks: %v\n", err)
		os.Exit(1)
	}
	p := tea.NewProgram(model{tasks: tasks})
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}
}
