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

## Example
```sh
todo add "Finish writing blog post"
todo due 1 2024-04-10
todo list

🎯 Learning Goals

This project helps explore:

🏗 Structs & Methods

📂 File I/O with JSON

⚡ Concurrency (future enhancement)

🖥 Command-line tools in Go

🔮 Future Enhancements

🏗 Better CLI handling with Cobra

🎨 Colored output for better readability

🖥 Terminal UI with Bubble Tea

📆 Task due dates and priorities

