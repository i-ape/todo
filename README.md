# todo

i wanna get good at go since typescript uses it :)

## 📝 Go ToDo List

## 📌 Project Overview

This is a simple **command-line To-Do List application** written in **Go**.  
The goal of this project is to **explore Go as a language** while building a useful tool.  

It covers:

- Structs and JSON serialization
- File handling for data persistence
- Command-line arguments processing
- Basic concurrency concepts (can be extended)

## Features

- Add tasks with optional due dates
- Support for **natural language dates** (e.g. `tomorrow`, `in 3 days`, `fri`)
- Mark tasks as complete
- Edit due dates after creation
- Delete or clear tasks
- Search by keyword
- Reset/delete the entire task database
- Color-coded task listing (due, overdue, complete)
- Add and edit tags with `@` or `#` syntax
- Filter tasks by tag (e.g. `todo tags @work`)
- bubble tea ui

---

## Installation

```bash
go build -o todo
chmod +x todo
./todo list

## 📂 Project Structure

- todo/
- ├── main.go          # Entry point, dispatch commands
- ├── commands.go      # All CLI logic (add/edit/list/etc)
- ├── task.go          # Task struct and core logic
- ├── storage.go       # JSON file read/write
- ├── tasks.json       # Auto-generated task data
+ .
+ ├── cmd/
+ │   ├── main.go
+ │   ├── commands.go
+ │   ├── handlers.go
+ │   └── tui.go
+ ├── todo.int/
+ │   ├── task.go
+ │   ├── storage.go
+ │   ├── due.go
+ │   ├── fzf.go
+ │   ├── utils.go
+ ├── config/
+ │   └── config.go
+ ├── tasks.json       # Generated user data


## 🔧 Installation

1. Install Go: [Download Go](https://go.dev/dl/)
2. Clone this repository:

   ```sh
   git clone https://github.com/i-ape/todo.git
   cd todo

⚡️ Optional Dependency:
- Install [`fzf`](https://github.com/junegunn/fzf) for interactive task selection:
  ```bash
  brew install fzf # or apt install fzf

3. Initialize Go module:

go mod init todo

🏃 Usage

1️⃣ Build the Program

go build -o todo
chmod +x todo
./todo list

## Commands

| Command            | Example usage                      | What to check                                 |
| ------------------ | ---------------------------------- | --------------------------------------------- |
| `add`              | `todo add "Write blog post" today` | Adds task with optional due date              |
| `list`             | `todo list`                        | Shows task list with correct labels, subtasks |
| `done`             | `todo done`                        | Marks selected task(s) done                   |
| `edit`             | `todo edit`                        | Updates task text                             |
| `delete`           | `todo delete`                      | Deletes selected task(s)                      |
| `due`              | `todo due 1 next week`             | Updates due date                              |
| `priority`         | `todo priority`                    | Prompts and updates priority                  |
| `tag`              | `todo tag`                         | Adds or edits tags                            |
| `search`           | `todo search blog`                 | Matches text and returns task(s)              |
| `note`             | `todo note 1 "Update with links"`  | Adds or edits notes                           |
| `move`             | `todo move`                        | Reorders task in list                         |
| `sub` or `subtask` | `todo sub 1 "Write draft"`         | Creates subtask under parent                  |
| `clear`            | `todo clear`                       | Clears all tasks (confirm it empties list)    |
| `reset`            | `todo reset`                       | Deletes `tasks.json`                          |


## Example

```sh
todo add "Write blog post"
todo add "Call mom every sunday"
todo list --tag=home
todo done 1
todo due 2 tomorrow
todo delete 1
todo edit
todo search "report"
todo tag
todo clear
todo reset
todo add "Gym @health" every mon,wed,fri
todo add "Standup meeting @work" every weekday @ 09:00
todo add "Call mom" every sunday @ 18:00
todo add "Quarterly report" bonm
todo add "End of month close" eom
todo due 4 eow
todo list --tag=work --priority=high --pending --json
todo add "Meeting @work" friday @ 14:00 for 45m
todo add "Call mom @family" sunday @ 18:00 for 1h for 3weeks
todo tui        # launch interactive interface
todo --tui list # use tui selection for list
todo list --fzf         # Launch FZF to filter
todo list --json        # Export filtered task list
todo move                  → Reorder a task to new position


## flags

| Flag                   | Example                     | Expected                        |
| ---------------------- | --------------------------- | ------------------------------- |
| `--tui`                | `todo --tui`                | Launches TUI (if implemented)   |
| `--no-fzf`             | `todo list --no-fzf`        | Bypasses interactive FZF        |
| `--json`               | `todo list --json`          | Returns JSON output             |
| `--today`              | `todo list --today`         | Filters today’s tasks           |
| `--done` / `--pending` | `todo list --done`          | Shows only completed/incomplete |
| `--tag=work`           | `todo list --tag=work`      | Filters by tag                  |
| `--priority=high`      | `todo list --priority=high` | Filters by priority             |

---

### 🧪 `todo pick` Flags

| Flag         | Example               | Description                                  |
| ------------ | --------------------- | -------------------------------------------- |
| `--json`     | `todo pick --json`    | Output selected task(s) in JSON              |
| `--edit`     | `todo pick --edit`    | Edit selected task(s) in place               |
| `--note`     | `todo pick --note`    | Add/edit notes for selected task(s)          |
|              | `todo pick`           | Print ID(s) of selected task(s)       


## 🧠 Learning Goals

✅ Structs & methods
✅ File I/O with JSON
✅ Command-line tools
✅ Natural language date parsing
✅ Optional interactivity with FZF
⏳ Potential: concurrency, custom DSL, Bubble Tea UI, cron



---

## 🔮 Future Ideas

Advanced tagging system (#tag, @context, filter/search by t

Recurring tasks (e.g. daily, weekly, every Mon,Wed)

Time-specific reminders (e.g. at 10am, 8:30p)

Agenda/week views (todo agenda, todo week)

fzf-powered interactive mode (todo pick, todo edit --fzf)

Task priorities (high, medium, low — visual indicators)

Archiving/completed history (todo archive, todo stats)

Multi-task selection (bulk delete/done/edit)

Task templates/snippets (e.g. meeting prep, habit task)

Custom aliasing system (e.g. shortcut commands or templates)

Sync/export support (e.g. sync to GitHub issues, export CSV)

Optional Bubble Tea TUI (interactive full-screen mode)

Push notifications/integrations (via cron, ntfy, or APIs)

 🏷️ Advanced tagging system (#tag, @context, filter/search/assign)

 🔁 Recurring tasks (e.g. every Monday at 10:00)

 🗓️ Calendar / agenda views (todo week)

 🚨 Notifications (ntfy, cron, push)

 📊 Stats & history (archive completed tasks)

 🔌 Integration with GitHub issues, CSV export

 🎭 Templates/snippets (todo new meeting)

 💅 Bubble Tea full-screen TUI

 🧠 Custom aliases / shortcuts


 test
