package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"task-cli/internal/task"
)

func Run(args []string) {
	if len(args) == 0 {
		printHelp()
		return
	}

	command := args[0]
	cmdArgs := args[1:]

	switch command {
	case "add", "a":
		addTask(cmdArgs)
	case "update", "u":
		updateTask(cmdArgs)
	case "delete", "del", "rm", "d":
		deleteTask(cmdArgs)
	case "list", "ls", "l":
		listTasks(cmdArgs)
	case "mark-in-progress", "start":
		markInProgress(cmdArgs)
	case "mark-done", "done", "complete":
		markDone(cmdArgs)
	case "help", "-h", "--help", "h", "?":
		printHelp()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Run 'taskcli help' for available commands")
		os.Exit(1)
	}
}

func addTask(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: taskcli add <description>")
		fmt.Println("Example: taskcli add \"Buy groceries\"")
		os.Exit(1)
	}
	description := strings.Join(args, " ")
	task.AddTask(description)
}

func updateTask(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: taskcli update <id> <description>")
		fmt.Println("Example: taskcli update 1 \"Buy milk and eggs\"")
		os.Exit(1)
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("Error: '%s' is not a valid number\n", args[0])
		os.Exit(1)
	}

	description := strings.Join(args[1:], " ")
	task.UpdateTaskTitle(id, description)
}

func deleteTask(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: taskcli delete <id>")
		fmt.Println("Example: taskcli delete 1")
		os.Exit(1)
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("Error: '%s' is not a valid number\n", args[0])
		os.Exit(1)
	}

	task.DeleteTask(id)
}

func listTasks(args []string) {
	if len(args) == 0 {
		task.ListAll()
		return
	}

	filter := args[0]
	switch filter {
	case "todo":
		task.List(task.Todo)
	case "in-progress":
		task.List(task.InProgress)
	case "done":
		task.List(task.Done)
	default:
		fmt.Printf("Invalid filter: %s\n", filter)
		fmt.Println("Available filters: todo, in-progress, done")
		os.Exit(1)
	}
}

func markInProgress(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: taskcli mark-in-progress <id>")
		fmt.Println("Example: taskcli mark-in-progress 1")
		os.Exit(1)
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("Error: '%s' is not a valid number\n", args[0])
		os.Exit(1)
	}

	task.ChangeState(id, task.InProgress)
}

func markDone(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: taskcli mark-done <id>")
		fmt.Println("Example: taskcli mark-done 1")
		os.Exit(1)
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("Error: '%s' is not a valid number\n", args[0])
		os.Exit(1)
	}

	task.ChangeState(id, task.Done)
}

func printHelp() {
	fmt.Print(`
Task Manager CLI - Manage your tasks from the command line

USAGE:
  taskcli [command] [options]

COMMANDS:
  add, a      <description>    Add a new task
  update, u   <id> <desc>      Update task description
  delete, del <id>             Delete a task
  list, ls    [filter]         List tasks (filter: todo, in-progress, done, all)
  mark-in-progress, start <id> Mark task as in-progress
  mark-done, done, complete <id> Mark task as done
  help, -h                     Show this help

EXAMPLES:
  taskcli add "Buy groceries"
  taskcli add "Write quarterly report"
  taskcli list
  taskcli list todo
  taskcli list in-progress
  taskcli update 1 "Buy milk and eggs"
  taskcli mark-in-progress 2
  taskcli mark-done 1
  taskcli delete 3

NOTES:
  - All data is stored in data/tasks.json
  - The app auto-creates the data directory if it doesn't exist
`)
}
