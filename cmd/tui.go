// cmd/tui.go
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
	"github.com/charmbracelet/lipgloss"
)

// ==============================
// 🎨 Styles
// ==============================
var (
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00BFFF"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555"))
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#50FA7B"))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#BBBBBB"))
	cursorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF"))
	taskStyle    = lipgloss.NewStyle().PaddingLeft(2)
	stickyStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700"))
	dueStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA500"))
	tagStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#87CEEB"))
	importantStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF69B4"))
)

// ==============================
// 📦 Model
// ==============================
type model struct {
	tasks       []todo.Task     // List of loaded tasks
	cursor      int             // Current cursor position in the task list
	quitting    bool            // Flag to indicate if the TUI is quitting
	err         error           // Current error, if any
	errTimeout  time.Time       // Time when the error should auto-clear
	state       string          // Current TUI state (e.g., "normal", "new_task")
	newTask     todo.Task       // Temporary task struct for creation
	input       textinput.Model // Input field for user entries
	filterTag   string          // Tag filter for tasks
	showPending bool            // Flag to show only pending tasks
	showHelp    bool            // Flag to show the help screen
	saveMsg     string          // Temporary save success message
	saveTimeout time.Time       // Time when the save message should auto-clear
	viewMode    string          // View mode ("normal" or "sticky")
}

// ==============================
// 🧩 Init
// ==============================
func NewModel() model {
	tasks, err := todo.LoadTasks()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load tasks: %v\n", err)
		tasks = []todo.Task{}
	}
	input := textinput.New()
	input.Focus()
	return model{
		tasks:       tasks,
		state:       "normal",
		input:       input,
		viewMode:    "normal",
	}
}

func (m model) Init() tea.Cmd { return textinput.Blink }

// ==============================
// 🔁 Update
// ==============================
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// auto-clear messages
	if !m.errTimeout.IsZero() && time.Now().After(m.errTimeout) {
		m.err = nil
		m.errTimeout = time.Time{}
	}
	if m.saveMsg != "" && time.Now().After(m.saveTimeout) {
		m.saveMsg = ""
		m.saveTimeout = time.Time{}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "h":
			m.showHelp = !m.showHelp
			return m, nil
		case "n":
			m.state = "new_task"
			m.input = textinput.New()
			m.input.Placeholder = "➕ New task..."
			m.input.Focus()
			return m, textinput.Blink
		case "esc":
			m.state = "normal"
			m.input.Reset()
			return m, nil
		case "j", "down":
			if m.cursor < len(m.visibleTasks())-1 {
				m.cursor++
			}
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "enter", " ":
			tasks := m.visibleTasks()
			if len(tasks) == 0 {
				return m, nil
			}
			t := tasks[m.cursor]
			for i := range m.tasks {
				if m.tasks[i].ID == t.ID {
					m.tasks[i].Completed = !m.tasks[i].Completed
					break
				}
			}
			todo.SaveTasks(m.tasks)
			m.saveMsg = "Task toggled"
			m.saveTimeout = time.Now().Add(2 * time.Second)
		}
	}

	if m.state == "new_task" {
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		if key, ok := msg.(tea.KeyMsg); ok && key.String() == "enter" {
			text := strings.TrimSpace(m.input.Value())
			if text == "" {
				m.err = fmt.Errorf("task cannot be empty")
				m.errTimeout = time.Now().Add(3 * time.Second)
				m.state = "normal"
				m.input.Reset()
				return m, nil
			}
			task := todo.Task{ID: todo.NextTaskID(m.tasks), Text: text}
			m.tasks = append(m.tasks, task)
			_ = todo.SaveTasks(m.tasks)
			m.saveMsg = "Task added"
			m.saveTimeout = time.Now().Add(2 * time.Second)
			m.state = "normal"
			m.input.Reset()
		}
		return m, cmd
	}

	return m, nil
}

// ==============================
// 👁️ View
// ==============================
func (m model) View() string {
	if m.quitting {
		return "👋 Goodbye!\n"
	}
	if m.showHelp {
		return m.helpView()
	}

	var b strings.Builder
	b.WriteString(m.headerView())
	b.WriteString("\n")

	if m.state == "new_task" {
		b.WriteString(fmt.Sprintf("📝 %s\n\n", m.input.View()))
	}

	tasks := m.visibleTasks()
	if len(tasks) == 0 {
		b.WriteString(helpStyle.Render("No tasks yet — press 'n' to add one.\n"))
	} else {
		for i, t := range tasks {
			cursor := "  "
			if i == m.cursor {
				cursor = cursorStyle.Render("▶ ")
			}
			status := "[ ]"
			if t.Completed {
				status = successStyle.Render("[✓]")
			}
			line := fmt.Sprintf("%s %s %s", cursor, status, t.Text)
			if t.Sticky {
				line += stickyStyle.Render(" 📌")
			}
			if t.DueDate != "" {
				line += dueStyle.Render(fmt.Sprintf(" 📅 %s", t.DueDate))
			}
			if len(t.Tags) > 0 {
				line += tagStyle.Render(" 🏷️ " + strings.Join(t.Tags, ", "))
			}
			if t.Important {
				line += importantStyle.Render(" ❗")
			}
			b.WriteString(taskStyle.Render(line) + "\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(m.statusView())
	return b.String()
}

// ==============================
// 🧭 Subviews
// ==============================
func (m model) headerView() string {
	header := titleStyle.Render("🗒️ Todo TUI")
	if m.filterTag != "" {
		header += tagStyle.Render(fmt.Sprintf("  [filter: %s]", m.filterTag))
	}
	if m.showPending {
		header += dueStyle.Render("  [pending only]")
	}
	if m.viewMode == "sticky" {
		header += stickyStyle.Render("  [sticky view]")
	}
	return header
}

func (m model) statusView() string {
	switch {
	case m.err != nil:
		return errorStyle.Render(fmt.Sprintf("Error: %v\n", m.err))
	case m.saveMsg != "":
		return successStyle.Render(fmt.Sprintf("%s\n", m.saveMsg))
	default:
		return helpStyle.Render("↑/↓ navigate  ⏎ toggle  n new  h help  q quit\n")
	}
}

// ==============================
// ℹ️ Help
// ==============================
func (m model) helpView() string {
	return lipgloss.NewStyle().Padding(1, 2).Render(
		titleStyle.Render("📘 Todo Help\n\n") +
			"↑/↓ or j/k  Move cursor\n" +
			"⏎ or space  Toggle completion\n" +
			"n           Add new task\n" +
			"h           Show/hide help\n" +
			"q, ctrl+c   Quit\n" +
			"esc         Cancel input\n",
	)
}

// ==============================
// 🧠 Helpers
// ==============================
func (m model) visibleTasks() []todo.Task {
	var result []todo.Task
	for _, t := range m.tasks {
		if m.showPending && t.Completed {
			continue
		}
		if m.filterTag != "" {
			found := false
			for _, tag := range t.Tags {
				if tag == m.filterTag {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		result = append(result, t)
	}
	return result
}

// ==============================
// 🚀 Entry point
// ==============================
func StartTUI() {
	p := tea.NewProgram(NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}
}
