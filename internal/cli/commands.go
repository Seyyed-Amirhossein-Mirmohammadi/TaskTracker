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
	case "add", "a":
		return addTask(cmdArgs)
	case "update", "u":
		return updateTask(cmdArgs)
	case "delete", "del", "rm", "d":
		return deleteTask(cmdArgs)
	case "list", "ls", "l":
		return listTasks(cmdArgs)
	case "mark-in-progress", "start":
		return markInProgress(cmdArgs)
	case "mark-done", "done", "complete":
		return markDone(cmdArgs)
	case "help", "-h", "--help", "h", "?":
		printHelp()
		return nil
	default:
		return fmt.Errorf("unknown command: %s\nRun 'taskcli help' for available commands", command)
	}
}

func addTask(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: taskcli add <description>\nExample: taskcli add \"Buy groceries\"")
	}

	description := strings.Join(args, " ")
	id, err := task.AddTask(description)
	if err != nil {
		return err
	}

	fmt.Printf("Task added successfully. ID: %d\n", id)
	return nil
}

func updateTask(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: taskcli update <id> <description>\nExample: taskcli update 1 \"Buy milk and eggs\"")
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
		return fmt.Errorf("usage: taskcli delete <id>\nExample: taskcli delete 1")
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
		return fmt.Errorf("usage: taskcli mark-in-progress <id>\nExample: taskcli mark-in-progress 1")
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
		return fmt.Errorf("usage: taskcli mark-done <id>\nExample: taskcli mark-done 1")
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

	fmt.Printf("Found %d task(s):\n\n", len(tasks))
	for _, t := range tasks {
		displayTask(t)
	}
}

func displayTask(t task.Task) {
	statusStr := ""
	switch t.Status {
	case task.Todo:
		statusStr = "Todo"
	case task.InProgress:
		statusStr = "In Progress"
	case task.Done:
		statusStr = "Done"
	}

	fmt.Printf("[%d] %s\n", t.Id, t.Description)
	fmt.Printf("  Status: %s | Created: %s\n",
		statusStr,
		t.CreatedAt.Format("2006-01-02 15:04"))

	if t.UpdatedAt != nil {
		fmt.Printf("  Updated: %s\n", t.UpdatedAt.Format("2006-01-02 15:04"))
	} else {
		fmt.Printf("  Updated: Never\n")
	}
	fmt.Println()
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
  list, ls    [filter]         List tasks (filter: todo, in-progress, done)
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
