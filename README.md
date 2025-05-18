# todo

i wanna get good at go since typescript uses it :)

## 📝 Go CLI To-Do List

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

---

## Installation

```bash
go build -o todo-cli
chmod +x todo-cli
./todo-cli list

## 📂 Project Structure

todo-cli/

- │── main.go          # Entry point, calls command handlers
- │── task.go          # Task struct and related functions
- │── storage.go       # Reads/Writes tasks to a JSON file
- │── commands.go      # CLI command handlers
- │── tasks.json       # JSON file (created at runtime)

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

## Commands

- todo add "Write blog post"
- todo add "Submit tax return" tomorrow
- todo list
- todo done 1
- todo due 2 fri
- todo search "blog"
- todo delete 1
- todo clear
- todo tags [@tag|#tag]         → Filter tasks by tag

## Example
```sh
todo add "Finish writing blog post"
todo due 1 2024-04-10
todo list
todo add "Write report @work #priority"
todo tags @work


current puzzle, more human but short
todo add "Gym @health" every mon,wed,fri
todo add "Standup meeting @work" every weekday @ 09:00
todo add "Call mom" every sunday @ 18:00

## 🧠 Learning Goals

✅ Structs & methods

✅ File I/O with JSON

✅ Command-line tools

✅ Natural language date parsing

⏳ Potential: concurrency, custom date DSL, Bubble Tea UI



---

🔮 Future Ideas

Cobra or urfave CLI parser

Task priorities, x

Tags & filters, x

Weekly/agenda views

Notification integrations
