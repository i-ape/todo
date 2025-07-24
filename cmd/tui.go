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
	tasks     []todo.Task
	cursor    int
	quitting  bool
	err       error
	state     string          // Tracks TUI state (e.g., "normal", "new_task", "new_priority", "new_due", "new_tags")
	newTask   todo.Task       // Temporary task being created
	input     textinput.Model // Current input field
	filterTag string
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	fmt.Fprintf(os.Stderr, "Key pressed: %s\n", msg) // Debug log
	if m.state != "normal" {
		var cmd tea.Cmd
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "q", "esc":
				m.state = "normal"
				m.input.Reset()
				m.err = nil
				return m, nil
			case "enter":
				value := strings.TrimSpace(m.input.Value())
				switch m.state {
				case "new_task":
					if value != "" {
						m.newTask = todo.Task{
							ID:   todo.NextTaskID(m.tasks),
							Text: value,
						}
						m.state = "new_priority"
						m.input.Reset()
						m.input.Placeholder = "🔥 Enter priority (high, medium, low, or enter for default):"
						return m, textinput.Blink
					}
					m.err = fmt.Errorf("task name cannot be empty")
					m.state = "normal"
					m.input.Reset()
					return m, nil
				case "new_priority":
					priority := config.GetDefaultPriority()
					if value != "" {
						if todo.ParsePriority(value) == "" {
							m.err = fmt.Errorf("invalid priority: %s", value)
							m.state = "normal"
							m.input.Reset()
							return m, nil
						}
						priority = value
					}
					m.newTask.Priority = todo.ParsePriority(priority)
					m.state = "new_due"
					m.input.Reset()
					m.input.Placeholder = "📅 Enter due date (e.g., today, 2025-12-31, or enter to skip):"
					return m, textinput.Blink
				case "new_due":
					if value != "" {
						dueDate, err := todo.ParseNaturalDate(value)
						if err != nil {
							m.err = err
							m.state = "normal"
							m.input.Reset()
							return m, nil
						}
						m.newTask.DueDate = dueDate
					}
					m.state = "new_tags"
					m.input.Reset()
					m.input.Placeholder = "🏷️ Enter tags (comma-separated, or enter to skip):"
					return m, textinput.Blink
				case "new_tags":
					if value != "" {
						m.newTask.Tags = todo.ParseTags(value)
					}
					m.tasks = append(m.tasks, m.newTask)
					if err := todo.SaveTasks(m.tasks); err != nil {
						m.err = err
					} else {
						m.err = nil
					}
					m.state = "normal"
					m.input.Reset()
					return m, nil
				case "edit_text":
					if value != "" {
						m.tasks[m.cursor].Text = value
						if err := todo.SaveTasks(m.tasks); err != nil {
							m.err = err
						} else {
							m.err = nil
						}
					}
					m.state = "normal"
					m.input.Reset()
					return m, nil
				case "edit_due":
					if value != "" {
						dueDate, err := todo.ParseNaturalDate(value)
						if err != nil {
							m.err = err
							m.state = "normal"
							m.input.Reset()
							return m, nil
						}
						m.tasks[m.cursor].DueDate = dueDate
						if err := todo.SaveTasks(m.tasks); err != nil {
							m.err = err
						} else {
							m.err = nil
						}
					}
					m.state = "normal"
					m.input.Reset()
					return m, nil
				case "edit_tags":
					if value != "" {
						m.tasks[m.cursor].Tags = todo.ParseTags(value)
						if err := todo.SaveTasks(m.tasks); err != nil {
							m.err = err
						} else {
							m.err = nil
						}
					}
					m.state = "normal"
					m.input.Reset()
					return m, nil
				case "edit_priority":
					if value != "" {
						if todo.ParsePriority(value) == "" {
							m.err = fmt.Errorf("invalid priority: %s", value)
							m.state = "normal"
							m.input.Reset()
							return m, nil
						}
						m.tasks[m.cursor].Priority = todo.ParsePriority(value)
						if err := todo.SaveTasks(m.tasks); err != nil {
							m.err = err
						} else {
							m.err = nil
						}
					}
					m.state = "normal"
					m.input.Reset()
					return m, nil
				case "filter_tag":
					if value != "" {
						m.filterTag = value
					} else {
						m.filterTag = "" // Clear filter if empty
					}
					m.state = "normal"
					m.input.Reset()
					return m, nil
				case "new_subtask":
					if value != "" {
						subtask := todo.Task{
							ID:        todo.NextTaskID(m.tasks),
							Text:      value,
							ParentID:  m.tasks[m.cursor].ID,
							Completed: false,
							Priority:  todo.ParsePriority(config.GetDefaultPriority()),
						}
						m.tasks = append(m.tasks, subtask)
						if err := todo.SaveTasks(m.tasks); err != nil {
							m.err = err
						} else {
							m.err = nil
						}
					}
					m.state = "normal"
					m.input.Reset()
					return m, nil
				}
			default:
				m.input, cmd = m.input.Update(msg)
			}
		}
		return m, cmd
	}

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
					fmt.Fprintf(os.Stderr, "Error saving tasks: %v\n", err)
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
					fmt.Fprintf(os.Stderr, "Error saving tasks: %v\n", err)
				} else {
					m.err = nil
				}
			}
		case "d":
			m.state = "edit_due"
			m.input = textinput.New()
			m.input.Placeholder = "📅 Enter new due date (e.g., today, 2025-12-31):"
			m.input.Focus()
			return m, textinput.Blink
		case "e":
			m.state = "edit_text"
			m.input = textinput.New()
			m.input.Placeholder = "✏️ Edit task text:"
			m.input.Focus()
			return m, textinput.Blink
		case "n":
			m.state = "new_task"
			m.input = textinput.New()
			m.input.Placeholder = "➕ New task:"
			m.input.Focus()
			return m, textinput.Blink
		case "t":
			m.state = "edit_tags"
			m.input = textinput.New()
			m.input.Placeholder = "🏷️ Enter tags (comma-separated):"
			m.input.Focus()
			return m, textinput.Blink
		case "p":
			m.state = "edit_priority"
			m.input = textinput.New()
			m.input.Placeholder = "🔥 Enter priority (high, medium, low):"
			m.input.Focus()
			return m, textinput.Blink
		case "s":
			m.state = "new_subtask"
			m.input = textinput.New()
			m.input.Placeholder = "📌 Subtask text:"
			m.input.Focus()
			return m, textinput.Blink
		case "f":
			if config.DisableFzf {
				m.err = fmt.Errorf("fzf is disabled")
				fmt.Fprintf(os.Stderr, "Error: fzf is disabled\n")
				return m, nil
			}
			selected, err := todo.SelectTasksWithFzf(false, config.DisableFzf)
			if err != nil {
				m.err = err
				fmt.Fprintf(os.Stderr, "Error selecting task with fzf: %v\n", err)
				return m, nil
			}
			if len(selected) == 0 {
				m.err = fmt.Errorf("no task selected")
				fmt.Fprintf(os.Stderr, "No task selected with fzf\n")
				return m, nil
			}
			for i, task := range m.tasks {
				if task.ID == selected[0].ID {
					m.cursor = i
					m.err = nil
					fmt.Fprintf(os.Stderr, "Selected task ID %d at cursor %d\n", task.ID, i)
					break
				}
			}
			return m, nil
		case "r":
			todo.SortTasks(m.tasks, config.GetSortOrder())
			if err := todo.SaveTasks(m.tasks); err != nil {
				m.err = err
				fmt.Fprintf(os.Stderr, "Error saving tasks: %v\n", err)
			} else {
				m.err = nil
			}
		case "c": // Clear error
			m.err = nil
			fmt.Fprintf(os.Stderr, "Error cleared\n")
			return m, nil
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
	taskMap := make(map[int][]todo.Task)
	topLevel := make([]todo.Task, 0, len(m.tasks)) // Pre-allocate
	for _, task := range m.tasks {
		if task.ParentID == 0 {
			topLevel = append(topLevel, task)
		} else {
			taskMap[task.ParentID] = append(taskMap[task.ParentID], task)
		}
	}
	for i, task := range topLevel { // Iterate topLevel instead of m.tasks
		cursor := "  "
		if i == m.cursor {
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
		// Display subtasks
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

func StartTUI() {
	tasks, err := todo.LoadTasks()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load tasks: %v\n", err)
		os.Exit(1)
	}
	p := tea.NewProgram(model{tasks: tasks}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}
}
