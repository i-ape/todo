package main

import (
	"fmt"
	"os"
	"strings"
	"time"
	"todo/config"
	todo "todo/td"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
)

type model struct {
	tasks       []todo.Task
	cursor      int
	quitting    bool
	err         error
	errTimeout  time.Time       // For auto-clearing errors
	state       string          // Tracks TUI state (e.g., "normal", "new_task", "new_priority", "new_due", "new_tags")
	newTask     todo.Task       // Temporary task being created
	input       textinput.Model // Current input field
	filterTag   string          // Tag to filter tasks
	showPending bool            // Show only incomplete tasks
	showHelp    bool            // Show help screen
	saveMsg     string          // Temporary save confirmation message
	saveTimeout time.Time       // For auto-clearing save message
	viewMode    string          // "normal" or "sticky"

}

func NewModel() model {
	tasks, err := todo.LoadTasks()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load tasks: %v\n", err)
		os.Exit(1)
	}
	input := textinput.New()
	input.Focus()
	return model{
		tasks:       tasks,
		cursor:      0,
		state:       "normal",
		input:       input,
		err:         nil,
		newTask:     todo.Task{},
		quitting:    false,
		filterTag:   "",
		showPending: false,
		showHelp:    false,
		saveMsg:     "",
		viewMode:    "normal", // sticky
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Auto-clear error after 3 seconds
	if !m.errTimeout.IsZero() && time.Now().After(m.errTimeout) {
		m.err = nil
		m.errTimeout = time.Time{}
	}
	// Auto-clear save message after 2 seconds
	if m.saveMsg != "" && time.Now().After(m.saveTimeout) {
		m.saveMsg = ""
		m.saveTimeout = time.Time{}
	}

	// Handle help screen toggle
	if m.showHelp {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "q", "esc", "h":
				m.showHelp = false
				return m, nil
			}
		}
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "esc":
			if m.state != "normal" {
				m.state = "normal"
				m.input.Reset()
				m.err = nil
				return m, nil
			}
		case "v":
			// Toggle between normal and sticky view modes
			if m.viewMode == "normal" {
				m.viewMode = "sticky"
			} else {
				m.viewMode = "normal"
			}
			m.cursor = 0 // Reset cursor to avoid out-of-bounds errors
			return m, nil
		case "m":
			// Toggle the sticky status of the selected task
			if len(m.tasks) > 0 {
				m.tasks[m.cursor].Sticky = !m.tasks[m.cursor].Sticky
			}
			return m, nil
		case "h":
			m.showHelp = !m.showHelp
			return m, nil
		case "i":
			if len(m.tasks) > 0 {
				m.tasks[m.cursor].Important = !m.tasks[m.cursor].Important
				if err := todo.SaveTasks(m.tasks); err != nil {
					m.err = err
					m.errTimeout = time.Now().Add(3 * time.Second)
				} else {
					m.saveMsg = "Task importance toggled"
					m.saveTimeout = time.Now().Add(2 * time.Second)
				}
			}
		case "o": // Toggle showing only pending tasks
			m.showPending = !m.showPending
			m.cursor = 0 // Reset cursor when filter changes
			return m, nil
		case "f":
			if m.state == "normal" && !config.DisableFzf {
				m.state = "filter_tag"
				m.input = textinput.New()
				m.input.Placeholder = "🔍 Enter tag to filter (or enter to clear):"
				m.input.Focus()
				return m, textinput.Blink
			}
		}
	}

	if m.state != "normal" {
		var cmd tea.Cmd
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
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
					m.errTimeout = time.Now().Add(3 * time.Second)
					m.state = "normal"
					m.input.Reset()
					return m, nil
				case "new_priority":
					priority := config.GetDefaultPriority()
					if value != "" {
						if todo.ParsePriority(value) == "" {
							m.err = fmt.Errorf("invalid priority: %s", value)
							m.errTimeout = time.Now().Add(3 * time.Second)
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
							m.errTimeout = time.Now().Add(3 * time.Second)
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
						m.errTimeout = time.Now().Add(3 * time.Second)
					} else {
						m.saveMsg = "Task saved successfully"
						m.saveTimeout = time.Now().Add(2 * time.Second)
					}
					m.state = "normal"
					m.input.Reset()
					m.cursor = len(m.tasks) - 1 // Move cursor to new task
					return m, nil
				case "edit_text":
					if value != "" {
						m.tasks[m.cursor].Text = value
						if err := todo.SaveTasks(m.tasks); err != nil {
							m.err = err
							m.errTimeout = time.Now().Add(3 * time.Second)
						} else {
							m.saveMsg = "Task updated successfully"
							m.saveTimeout = time.Now().Add(2 * time.Second)
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
							m.errTimeout = time.Now().Add(3 * time.Second)
							m.state = "normal"
							m.input.Reset()
							return m, nil
						}
						m.tasks[m.cursor].DueDate = dueDate
						if err := todo.SaveTasks(m.tasks); err != nil {
							m.err = err
							m.errTimeout = time.Now().Add(3 * time.Second)
						} else {
							m.saveMsg = "Due date updated successfully"
							m.saveTimeout = time.Now().Add(2 * time.Second)
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
							m.errTimeout = time.Now().Add(3 * time.Second)
						} else {
							m.saveMsg = "Tags updated successfully"
							m.saveTimeout = time.Now().Add(2 * time.Second)
						}
					}
					m.state = "normal"
					m.input.Reset()
					return m, nil
				case "edit_priority":
					if value != "" {
						if todo.ParsePriority(value) == "" {
							m.err = fmt.Errorf("invalid priority: %s", value)
							m.errTimeout = time.Now().Add(3 * time.Second)
							m.state = "normal"
							m.input.Reset()
							return m, nil
						}
						m.tasks[m.cursor].Priority = todo.ParsePriority(value)
						if err := todo.SaveTasks(m.tasks); err != nil {
							m.err = err
							m.errTimeout = time.Now().Add(3 * time.Second)
						} else {
							m.saveMsg = "Priority updated successfully"
							m.saveTimeout = time.Now().Add(2 * time.Second)
						}
					}
					m.state = "normal"
					m.input.Reset()
					return m, nil
				case "filter_tag":
					if value != "" {
						m.filterTag = value
					} else {
						m.filterTag = "" // Clear filter
					}
					m.cursor = 0 // Reset cursor when filter changes
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
							m.errTimeout = time.Now().Add(3 * time.Second)
						} else {
							m.saveMsg = "Subtask added successfully"
							m.saveTimeout = time.Now().Add(2 * time.Second)
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
		case "j", "down":
			if m.cursor < len(m.visibleTasks())-1 {
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
					m.errTimeout = time.Now().Add(3 * time.Second)
				} else {
					m.saveMsg = "Task updated successfully"
					m.saveTimeout = time.Now().Add(2 * time.Second)
				}
			}
		case "x", "backspace":
			if len(m.tasks) > 0 {
				m.tasks = append(m.tasks[:m.cursor], m.tasks[m.cursor+1:]...)
				if m.cursor >= len(m.visibleTasks()) && m.cursor > 0 {
					m.cursor--
				}
				if err := todo.SaveTasks(m.tasks); err != nil {
					m.err = err
					m.errTimeout = time.Now().Add(3 * time.Second)
				} else {
					m.saveMsg = "Task deleted successfully"
					m.saveTimeout = time.Now().Add(2 * time.Second)
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
				m.errTimeout = time.Now().Add(3 * time.Second)
				return m, nil
			}
			selected, err := todo.SelectTasksWithFzf(false, config.DisableFzf)
			if err != nil {
				m.err = err
				m.errTimeout = time.Now().Add(3 * time.Second)
				return m, nil
			}
			if len(selected) == 0 {
				m.err = fmt.Errorf("no task selected")
				m.errTimeout = time.Now().Add(3 * time.Second)
				return m, nil
			}
			for i, task := range m.tasks {
				if task.ID == selected[0].ID {
					m.cursor = i
					m.err = nil
					break
				}
			}
			return m, nil
		case "r":
			todo.SortTasks(m.tasks, config.GetSortOrder())
			if err := todo.SaveTasks(m.tasks); err != nil {
				m.err = err
				m.errTimeout = time.Now().Add(3 * time.Second)
			} else {
				m.saveMsg = "Tasks sorted successfully"
				m.saveTimeout = time.Now().Add(2 * time.Second)
			}
			return m, nil
		case "c":
			m.err = nil
			m.errTimeout = time.Time{}
			return m, nil
		}
	}
	return m, nil
}

// visibleTasks returns tasks filtered by tag and completion status
func (m model) visibleTasks() []todo.Task {
	var result []todo.Task
	for _, task := range m.tasks {
		if m.viewMode == "sticky" && !task.Sticky {
			continue
		}
		if m.showPending && task.Completed {
			continue
		}
		if m.filterTag != "" {
			hasTag := false
			for _, tag := range task.Tags {
				if tag == m.filterTag {
					hasTag = true
					break
				}
			}
			if !hasTag {
				continue
			}
		}
		result = append(result, task)
	}
	return result
}

func (m model) View() string {
	if m.quitting {
		return "Goodbye 👋\n"
	}
	if m.showHelp {
		return m.helpView()
	}
	var b strings.Builder
	if m.err != nil {
		b.WriteString(color.RedString("Error: %v\n\n", m.err))
	}
	if m.saveMsg != "" {
		b.WriteString(color.GreenString("%s\n\n", m.saveMsg))
	}
	if m.viewMode == "sticky" && len(m.visibleTasks()) == 0 {
		b.WriteString("📌 No sticky tasks. Press 'm' to mark a task as sticky.\n\n")
	} else if len(m.visibleTasks()) == 0 {
		b.WriteString("📋 No tasks yet. Press 'n' to add one.\n\n")
	} else {
		b.WriteString("📋 Tasks:\n\n")
		taskMap := make(map[int][]todo.Task)
		var topLevel []todo.Task
		for _, task := range m.visibleTasks() {
			if task.ParentID == 0 {
				topLevel = append(topLevel, task)
			} else {
				taskMap[task.ParentID] = append(taskMap[task.ParentID], task)
			}
		}
		for i, task := range topLevel {
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
			if task.Sticky {
				label = "📌 " + label
			}
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
			for _, subtask := range taskMap[task.ID] {
				subStatus := color.CyanString("[ ]")
				if subtask.Completed {
					subStatus = color.GreenString("[✓]")
				} else if subtask.DueDate != "" && todo.IsOverdue(subtask.DueDate) {
					subStatus = color.RedString("[✗]")
				}
				subLabel := subtask.Text
				if subtask.Sticky {
					subLabel = "📌 " + subLabel
				}
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
	}
	if m.filterTag != "" {
		b.WriteString(fmt.Sprintf("\nFiltered by tag: %s\n", m.filterTag))
	}
	if m.showPending {
		b.WriteString("Showing only pending tasks\n")
	}
	b.WriteString("\n↑/↓ or j/k: navigate, [enter]: toggle, [x]: delete, [n]: new, [d]: due date, [e]: edit, [t]: tags, [p]: priority, [s]: subtask, [f]: filter/fzf, [o]: toggle pending, [r]: sort, [h]: help, [m]: toggle sticky, [v]: sticky view, [q]: quit\n")
	return b.String()
}

// helpView renders the help screen
func (m model) helpView() string {
	var b strings.Builder
	b.WriteString(color.CyanString("📋 To-Do List Help\n\n"))
	b.WriteString("Keybindings:\n")
	b.WriteString("  ↑/↓ or j/k: Move cursor up/down\n")
	b.WriteString("  enter or space: Toggle task completion\n")
	b.WriteString("  x or backspace: Delete task\n")
	b.WriteString("  n: Add new task\n")
	b.WriteString("  d: Edit due date\n")
	b.WriteString("  e: Edit task text\n")
	b.WriteString("  t: Edit tags\n")
	b.WriteString("  p: Edit priority\n")
	b.WriteString("  s: Add subtask\n")
	b.WriteString("  f: Filter by tag (or fzf select if enabled)\n")
	b.WriteString("  o: Toggle showing only pending tasks\n")
	b.WriteString("  r: Sort tasks\n")
	b.WriteString("  c: Clear error\n")
	b.WriteString("  h: Show/hide this help\n")
	b.WriteString("  q or ctrl+c: Quit\n")
	b.WriteString("  esc: Cancel input (in editing modes)\n\n")
	b.WriteString("Configuration (Environment Variables):\n")
	b.WriteString("  TODO_DEFAULT_PRIORITY: Set default priority (high, medium, low; default: medium)\n")
	b.WriteString("  TODO_DISABLE_FZF: Set to 'true' to disable fzf (default: false)\n")
	b.WriteString("  TODO_SORT_ORDER: Set sort order (due, priority, text; default: due)\n\n")
	b.WriteString("Press h, q, or esc to return to tasks")
	b.WriteString("  i: Toggle task importance\n")
	return b.String()
}

// it starts the TUI
func StartTUI() {
	p := tea.NewProgram(NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}
}
