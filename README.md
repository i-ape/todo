# todo-cli

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

## 🚀 Features

- ✅ Add new tasks  
- 📋 List all tasks  
- ✔️ Mark tasks as completed  
- ❌ Delete tasks  
- 💾 Persistent storage using `tasks.json`  

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
   git clone https://github.com/i-ape/todo-cli.git
   cd todo-cli

3. Initialize Go module:

go mod init todo-cli

🏃 Usage

1️⃣ Build the Program

go build -o todo

## Commands

- `todo add "task text"` — Add a new task
- `todo list` — List all tasks
- `todo done [task ID or text]` — Mark a task as completed
- `todo delete [task ID or text]` — Delete a task
- `todo due [task ID or text] [YYYY-MM-DD]` — Set a due date

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

