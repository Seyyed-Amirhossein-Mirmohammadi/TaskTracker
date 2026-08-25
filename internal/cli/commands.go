package cli

import (
	"fmt"
	"strconv"
	"strings"
	"task-cli/internal/task"
)

func Run(args []string) error {
	if len(args) == 0 {
		printHelp()
		return nil
	}

	command := args[0]
	cmdArgs := args[1:]

	switch command {
	case "add":
		return addTask(cmdArgs)
	case "update":
		return updateTask(cmdArgs)
	case "delete":
		return deleteTask(cmdArgs)
	case "list":
		return listTasks(cmdArgs)
	case "mark-in-progress":
		return markInProgress(cmdArgs)
	case "mark-done":
		return markDone(cmdArgs)
	case "help", "-h", "--help":
		printHelp()
		return nil
	default:
		return fmt.Errorf("unknown command: %s\nRun 'task-cli help' for available commands", command)
	}
}

func addTask(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: task-cli add <description>")
	}

	description := strings.Join(args, " ")
	id, err := task.AddTask(description)
	if err != nil {
		return err
	}

	fmt.Printf("Task added successfully (ID: %d)\n", id)
	return nil
}

func updateTask(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: task-cli update <id> <description>")
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("'%s' is not a valid number", args[0])
	}

	description := strings.Join(args[1:], " ")
	if err := task.UpdateTaskTitle(id, description); err != nil {
		return err
	}

	fmt.Printf("Task %d updated successfully\n", id)
	return nil
}

func deleteTask(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: task-cli delete <id>")
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("'%s' is not a valid number", args[0])
	}

	if err := task.DeleteTask(id); err != nil {
		return err
	}

	fmt.Printf("Task %d deleted successfully\n", id)
	return nil
}

func listTasks(args []string) error {
	var tasks []task.Task
	var err error

	if len(args) == 0 {
		tasks, err = task.ListAll()
	} else {
		filter := args[0]
		switch filter {
		case "todo":
			tasks, err = task.List(task.Todo)
		case "in-progress":
			tasks, err = task.List(task.InProgress)
		case "done":
			tasks, err = task.List(task.Done)
		default:
			return fmt.Errorf("invalid filter: %s\nAvailable filters: todo, in-progress, done", filter)
		}
	}

	if err != nil {
		return err
	}

	displayTasks(tasks)
	return nil
}

func markInProgress(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: task-cli mark-in-progress <id>")
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("'%s' is not a valid number", args[0])
	}

	if err := task.ChangeState(id, task.InProgress); err != nil {
		return err
	}

	fmt.Printf("Task %d marked as in-progress\n", id)
	return nil
}

func markDone(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: task-cli mark-done <id>")
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("'%s' is not a valid number", args[0])
	}

	if err := task.ChangeState(id, task.Done); err != nil {
		return err
	}

	fmt.Printf("Task %d marked as done\n", id)
	return nil
}

func displayTasks(tasks []task.Task) {
	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}

	for _, t := range tasks {
		displayTask(t)
	}
}

func displayTask(t task.Task) {
	statusStr := t.Status.String()

	fmt.Printf("ID: %d\n", t.ID)
	fmt.Printf("Description: %s\n", t.Description)
	fmt.Printf("Status: %s\n", statusStr)
	fmt.Printf("Created At: %s\n", t.CreatedAt.Format("2006-01-02 15:04:05"))
	if t.UpdatedAt != nil {
		fmt.Printf("Updated At: %s\n", t.UpdatedAt.Format("2006-01-02 15:04:05"))
	}
	fmt.Println()
}

func printHelp() {
	fmt.Print(`
Task Manager CLI - Manage your tasks from the command line

Usage:
  task-cli <command> [arguments]

Commands:
  add <description>              Add a new task
  update <id> <description>      Update a task description
  delete <id>                    Delete a task
  mark-in-progress <id>          Mark a task as in progress
  mark-done <id>                 Mark a task as done
  list                           List all tasks
  list <status>                  List tasks by status (todo, in-progress, done)
  help                           Show this help

Examples:
  task-cli add "Buy groceries"
  task-cli update 1 "Buy groceries and cook dinner"
  task-cli delete 1
  task-cli mark-in-progress 1
  task-cli mark-done 1
  task-cli list
  task-cli list done
  task-cli list todo
  task-cli list in-progress

Notes:
  - All data is stored in data/tasks.json
  - The app auto-creates the data directory if it doesn't exist
`)
}
