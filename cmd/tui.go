// cmd/tui.go
package main

import (
	"fmt"
	"os"
	"strings"
	"todo/config"
	todo "todo/td"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
)

type model struct {
	tasks    []todo.Task
	cursor   int
	quitting bool
	err      error
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
				if err := todo.SaveTasks(m.tasks); err != nil {
					m.err = err
				} else {
					m.err = nil
				}
			}
		case "x", "backspace":
			if len(m.tasks) > 0 {
				m.tasks = append(m.tasks[:m.cursor], m.tasks[m.cursor+1:]...)
				if m.cursor >= len(m.tasks) && m.cursor > 0 {
					m.cursor--
				}
				if err := todo.SaveTasks(m.tasks); err != nil {
					m.err = err
				} else {
					m.err = nil
				}
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
					if err := todo.SaveTasks(m.tasks); err != nil {
						m.err = err
					}
				}
			}
		case "e":
			newText, ok := prompt("✏️ Edit task text:")
			if ok && strings.TrimSpace(newText) != "" {
				m.tasks[m.cursor].Text = strings.TrimSpace(newText)
				m.err = nil
				if err := todo.SaveTasks(m.tasks); err != nil {
					m.err = err
				}
			}
		case "n":
			newTask, ok := prompt("➕ New task:")
			if ok && strings.TrimSpace(newTask) != "" {
				tasks, err := todo.LoadTasks()
				if err != nil {
					m.err = err
				} else {
					task := todo.Task{
						ID:       todo.NextTaskID(tasks),
						Text:     strings.TrimSpace(newTask),
						Priority: todo.ParsePriority(config.GetDefaultPriority()),
					}
					m.tasks = append(m.tasks, task)
					m.err = nil
					if err := todo.SaveTasks(m.tasks); err != nil {
						m.err = err
					}
				}
			}
		case "t":
			tags, ok := prompt("🏷️ Enter tags (comma-separated):")
			if ok && strings.TrimSpace(tags) != "" {
				m.tasks[m.cursor].Tags = todo.ParseTags(tags)
				m.err = nil
				if err := todo.SaveTasks(m.tasks); err != nil {
					m.err = err
				}
			}
		case "p":
			priority, ok := prompt("🔥 Enter priority (high, medium, low):")
			if ok && strings.TrimSpace(priority) != "" {
				m.tasks[m.cursor].Priority = todo.ParsePriority(priority)
				m.err = nil
				if err := todo.SaveTasks(m.tasks); err != nil {
					m.err = err
				}
			}
		case "s":
			text, ok := prompt("📌 Subtask text:")
			if ok && strings.TrimSpace(text) != "" {
				tasks, err := todo.LoadTasks()
				if err != nil {
					m.err = err
				} else {
					subtask := todo.Task{
						ID:        todo.NextTaskID(tasks),
						Text:      strings.TrimSpace(text),
						ParentID:  m.tasks[m.cursor].ID,
						Completed: false,
						Priority:  todo.ParsePriority(config.GetDefaultPriority()),
					}
					m.tasks = append(m.tasks, subtask)
					m.err = nil
					if err := todo.SaveTasks(m.tasks); err != nil {
						m.err = err
					}
				}
			}
		case "f":
			if !config.DisableFzf {
				selected, err := todo.SelectTasksWithFzf(false, config.DisableFzf)
				if err != nil {
					m.err = err
				} else if len(selected) > 0 {
					for i, task := range m.tasks {
						if task.ID == selected[0].ID {
							m.cursor = i
							m.err = nil
							break
						}
					}
				}
			} else {
				m.err = fmt.Errorf("fzf is disabled")
			}
		case "r":
			tasks, err := todo.LoadTasks()
			if err != nil {
				m.err = err
			} else {
				m.tasks = tasks
				todo.SortTasks(m.tasks, config.GetSortOrder())
				m.err = nil
				if err := todo.SaveTasks(m.tasks); err != nil {
					m.err = err
				}
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
	// Group tasks by ParentID
	taskMap := make(map[int][]todo.Task)
	var topLevel []todo.Task
	for _, task := range m.tasks {
		if task.ParentID == 0 {
			topLevel = append(topLevel, task)
		} else {
			taskMap[task.ParentID] = append(taskMap[task.ParentID], task)
		}
	}
	cursorIndex := -1
	for i, task := range topLevel {
		if m.cursor == len(topLevel) {
			cursorIndex = i
		}
		cursor := "  "
		if i == cursorIndex {
			cursor = "▶ "
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
		b.WriteString(fmt.Sprintf("%s %s %s\n", cursor, status, label))
		// Add subtasks
		for _, subtask := range taskMap[task.ID] {
			subStatus := color.CyanString("[ ]")
			if subtask.Completed {
				subStatus = color.GreenString("[✓]")
			} else if subtask.DueDate != "" && todo.IsOverdue(subtask.DueDate) {
				subStatus = color.RedString("[✗]")
			}
			subLabel := subtask.Text
			if subtask.DueDate != "" {
				subLabel += color.YellowString(" 📅 %s", subtask.DueDate)
			}
			if len(subtask.Tags) > 0 {
				subLabel += " 🏷️ " + strings.Join(subtask.Tags, ", ")
			}
			switch subtask.Priority {
			case "high":
				subLabel += color.RedString(" 🔥")
			case "low":
				subLabel += color.BlueString(" ⬇")
			}
			b.WriteString(fmt.Sprintf("    ↳ %s %s\n", subStatus, subLabel))
		}
	}
	b.WriteString("\n↑/↓ or j/k: navigate, [enter]: toggle, [x]: delete, [n]: new, [d]: due date, [e]: edit, [t]: tags, [p]: priority, [s]: subtask, [f]: fzf select, [r]: sort, [q]: quit\n")
	return b.String()
}

type promptModel struct {
	input textinput.Model
}

func (p promptModel) Init() tea.Cmd {
	return textinput.Blink
}

func (p promptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	p.input, cmd = p.input.Update(msg)
	return p, cmd
}

func (p promptModel) View() string {
	return fmt.Sprintf("\n%s\n> %s", p.input.Placeholder, p.input.View())
}

func prompt(promptText string) (string, bool) {
	input := textinput.New()
	input.Placeholder = promptText
	input.Focus()
	p := tea.NewProgram(promptModel{input: input})
	m, err := p.Run()
	if err != nil {
		return "", false
	}
	final := m.(promptModel)
	value := strings.TrimSpace(final.input.Value())
	return value, value != ""
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
