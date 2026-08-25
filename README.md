# Task-CLI: TODO Task Manager

<a href="https://roadmap.sh/projects/task-tracker" target="_blank">
<img src="https://img.shields.io/badge/roadmap.sh-task%20tracker-blue" alt="Task Tracker">
</a>

## 📖 Description

**Task-CLI** is a lightweight command-line interface (CLI) application for efficient task management.  
Built in Go, it allows you to **add, update, delete, list, and track tasks** directly from your terminal with persistent JSON storage.

---

## ✨ Features

- **Add a Task** → Create tasks with descriptions. Each task gets a unique ID and a default `todo` status.
- **Update a Task** → Modify the description of an existing task.
- **Mark as In-Progress** → Quickly change a task's status to `in-progress`.
- **Mark as Done** → Quickly change a task's status to `done`.
- **Delete a Task** → Remove tasks by their ID.
- **List Tasks** → Display all tasks or filter them by status: `todo`, `in-progress`, or `done`.

---

## 🗂 Project Structure

```
task-cli/
├── main.go                 # Application entry point
├── go.mod                  # Go module definition
├── internal/
│   ├── task/
│   │   ├── types.go        # Task and Status definitions
│   │   ├── store.go        # JSON persistence and storage
│   │   ├── task.go         # Task business logic
│   │   ├── store_test.go   # Unit tests for storage
│   │   └── task_test.go    # Unit tests for business logic
│   └── cli/
│       ├── commands.go     # CLI command handlers
│       └── commands_test.go # CLI integration tests
└── data/
    └── tasks.json          # Task data storage (auto-created)
```

---

## ⚡ Installation

### Build from Source

```bash
git clone https://github.com/yourusername/task-cli.git
cd task-cli
go build -o task-cli
```

### Install via Go

```bash
go install github.com/yourusername/task-cli@latest
```

---

## 🚀 Usage

```bash
task-cli <command> [arguments]
```

### Commands

| Command | Description |
|---------|-------------|
| `add <description>` | Add a new task |
| `update <id> <description>` | Update a task description |
| `delete <id>` | Delete a task |
| `mark-in-progress <id>` | Mark a task as in-progress |
| `mark-done <id>` | Mark a task as done |
| `list` | List all tasks |
| `list <status>` | List tasks by status (todo, in-progress, done) |
| `help` | Show help |

---

## 📝 Examples

```bash
# Add a task
task-cli add "Buy groceries"
# Output: Task added successfully (ID: 1)

# Update a task
task-cli update 1 "Buy groceries and cook dinner"
# Output: Task 1 updated successfully

# Mark as in-progress
task-cli mark-in-progress 1
# Output: Task 1 marked as in-progress

# Mark as done
task-cli mark-done 1
# Output: Task 1 marked as done

# List all tasks
task-cli list
# Output:
# ID: 1
# Description: Buy groceries and cook dinner
# Status: done
# Created At: 2026-08-24 10:30:00
# Updated At: 2026-08-24 14:30:00

# List tasks by status
task-cli list todo
task-cli list in-progress
task-cli list done

# Delete a task
task-cli delete 1
# Output: Task 1 deleted successfully

# Show help
task-cli help
```

---

## 📦 Data Storage

Tasks are stored in `data/tasks.json` in the current working directory. The file is automatically created on first use.

### Example JSON Structure

```json
{
  "tasks": [
    {
      "id": 1,
      "description": "Buy groceries and cook dinner",
      "status": "done",
      "created_at": "2026-08-24T10:30:00Z",
      "updated_at": "2026-08-24T14:30:00Z"
    },
    {
      "id": 2,
      "description": "Write quarterly report",
      "status": "todo",
      "created_at": "2026-08-24T10:31:00Z",
      "updated_at": null
    }
  ],
  "last_id": 2
}
```

---

## 🧪 Testing

The project includes a comprehensive test suite covering storage operations, business logic, and CLI integration.

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage report
go test -cover ./...

# Generate detailed coverage profile
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Test Coverage

Tests are organized by package:

| Package | Tests Cover |
|---------|-------------|
| `internal/task` | Store initialization, CRUD operations, ID generation, error handling, status filtering |
| `internal/cli` | Command parsing, success/failure outputs, help, argument validation, integration with task package |

All tests use temporary directories to avoid affecting your real data.
