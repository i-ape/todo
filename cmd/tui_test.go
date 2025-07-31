package main

import (
	"strings"
	"testing"
	"time"

	"todo/td"

	tea "github.com/charmbracelet/bubbletea"
)

// Mock Task for testing
type Task = td.Task // Assuming Task is defined in td/task.go

// Mock LoadTasks/SaveTasks
var mockTasks = []Task{
	{ID: 1, Text: "Buy coffee", DueDate: time.Now().Local().AddDate(0, 0, -1).Format("2006-01-02")},
	{ID: 2, Text: "Code app", DueDate: ""},
}

func mockLoadTasks() ([]Task, error) {
	return mockTasks, nil
}

func mockSaveTasks(tasks []Task) error {
	mockTasks = tasks
	return nil
}

func TestModelInit(t *testing.T) {
	td.LoadTasks = mockLoadTasks
	td.SaveTasks = mockSaveTasks
	defer func() {
		td.LoadTasks = td.LoadTasks // Restore original
		td.SaveTasks = td.SaveTasks
	}()

	m := NewModel()
	if m.state != "normal" {
		t.Errorf("NewModel() state = %q, want %q", m.state, "normal")
	}
	if m.input == nil {
		t.Error("NewModel() input is nil")
	}
	if len(m.tasks) != len(mockTasks) {
		t.Errorf("NewModel() tasks len = %d, want %d", len(m.tasks), len(mockTasks))
	}
}

func TestUpdateNewTask(t *testing.T) {
	td.LoadTasks = mockLoadTasks
	td.SaveTasks = mockSaveTasks
	defer func() {
		td.LoadTasks = td.LoadTasks
		td.SaveTasks = td.SaveTasks
	}()

	m := NewModel()
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	if m.state != "new_task" {
		t.Errorf("Update('n') state = %q, want %q", m.state, "new_task")
	}
	if m.input.Placeholder != "📝 Enter task description:" {
		t.Errorf("Update('n') placeholder = %q, want %q", m.input.Placeholder, "📝 Enter task description:")
	}
}

func TestUpdateSetDueDate(t *testing.T) {
	td.LoadTasks = mockLoadTasks
	td.SaveTasks = mockSaveTasks
	defer func() {
		td.LoadTasks = td.LoadTasks
		td.SaveTasks = td.SaveTasks
	}()

	m := NewModel()
	m.cursor = 0
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	if m.state != "edit_due" {
		t.Errorf("Update('d') state = %q, want %q", m.state, "edit_due")
	}

	m.input.SetValue("tomorrow")
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if m.state != "normal" {
		t.Errorf("Update(enter 'tomorrow') state = %q, want %q", m.state, "normal")
	}
	expectedDue := time.Now().Local().AddDate(0, 0, 1).Format("2006-01-02")
	if len(m.tasks) < 1 || m.tasks[0].DueDate != expectedDue {
		t.Errorf("Update(enter 'tomorrow') task[0].DueDate = %q, want %q", m.tasks[0].DueDate, expectedDue)
	}
}

func TestUpdateAddTask(t *testing.T) {
	td.LoadTasks = mockLoadTasks
	td.SaveTasks = mockSaveTasks
	defer func() {
		td.LoadTasks = td.LoadTasks
		td.SaveTasks = td.SaveTasks
	}()

	m := NewModel()
	mockTasks = []Task{} // Start with empty task list
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	m.input.SetValue("New task")
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m.input.SetValue("today")
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if m.state != "normal" {
		t.Errorf("Update(enter 'today') state = %q, want %q", m.state, "normal")
	}
	if len(m.tasks) != 1 {
		t.Errorf("Update(add task) tasks len = %d, want 1", len(m.tasks))
	}
	expectedDue := time.Now().Local().Format("2006-01-02")
	if m.tasks[0].Text != "New task" || m.tasks[0].DueDate != expectedDue {
		t.Errorf("Update(add task) task = %+v, want Text='New task', DueDate=%q", m.tasks[0], expectedDue)
	}
}

func TestView(t *testing.T) {
	td.LoadTasks = mockLoadTasks
	td.SaveTasks = mockSaveTasks
	defer func() {
		td.LoadTasks = td.LoadTasks
		td.SaveTasks = td.SaveTasks
	}()

	m := NewModel()
	view := m.View()
	if !strings.Contains(view, "Buy coffee") || !strings.Contains(view, "Code app") {
		t.Errorf("View() missing tasks, got: %s", view)
	}
	if strings.Contains(view, "ERROR") && m.err == nil {
		t.Errorf("View() shows ERROR when m.err is nil")
	}
}
