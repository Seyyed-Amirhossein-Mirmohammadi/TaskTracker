package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Seyyed-Amirhossein-Mirmohammadi/TaskTracker.git/internal/task"
)

func executeCommand(input string) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return
	}

	command := parts[0]
	args := parts[1:]

	switch command {
	case "add":
		AddTaskInteractive(args)
	case "update":
		UpdateTaskInteractive(args)
	case "delete":
		DeleteTaskInteractive(args)
	case "list":
		ListTasksInteractive(args)
	case "mark-in-progress":
		MarkInProgressInteractive(args)
	case "mark-done":
		MarkDoneInterctive(args)
	case "help", "h", "?":
		printHelp()
	default:
		fmt.Printf("Unknown command: %s\nType 'help' for available commands\n", command)
	}
}

func AddTaskInteractive(args []string) {
	if len(args) != 1 {
		fmt.Println("Invalid argument count. Just give a description.")
		return
	}
	task.AddTask(args[0])
}

func UpdateTaskInteractive(args []string) {
	if len(args) != 2 {
		fmt.Println("Invalid argument count. Just give and Id and a description.")
		return
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println(args[0] + " is not a number")
		return
	}

	task.UpdateTaskTitle(id, args[1])
}

func DeleteTaskInteractive(args []string) {
	if len(args) != 1 {
		fmt.Println("Invalid argument count. Just give an id.")
		return
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println(args[0] + " is not a number")
		return
	}
	task.DeleteTask(id)
}

func ListTasksInteractive(args []string) {
	if len(args) > 1 {
		fmt.Println("Invalid argument count.")
		return
	}

	if len(args) == 0 {
		task.ListAll()
	} else {
		status := args[0]
		if status == "done" {
			task.List(task.Done)
		} else if status == "todo" {
			task.List(task.Todo)
		} else if status == "in-progress" {
			task.List(task.InProgress)
		} else {
			fmt.Println("Invalid status.")
		}
	}
}

func MarkInProgressInteractive(args []string) {
	if len(args) != 1 {
		fmt.Println("Invalid argument count. Just give and Id.")
		return
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println(args[0] + " is not a number")
		return
	}

	task.ChangeState(id, task.InProgress)
}

func MarkDoneInterctive(args []string) {
	if len(args) != 1 {
		fmt.Println("Invalid argument count. Just give and Id.")
		return
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println(args[0] + " is not a number")
		return
	}

	task.ChangeState(id, task.Done)
}

func printHelp() {
	fmt.Println(`
Available Commands:
  add, a     - Add a new task
  list, ls   - List tasks
  complete, done, c - Mark a task as completed
  delete, del, rm, d - Delete a task
  help, h, ? - Show this help
  exit, quit, q - Exit the application

Examples:
  add Buy groceries --priority high --due 2026-08-30
  list --filter pending
  complete 1
  delete 1 --force`)
}
