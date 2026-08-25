package main

import (
	"fmt"
	"os"
	"task-cli/internal/cli"
	"task-cli/internal/task"
)

func main() {
	if err := task.InitTaskStore(); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing task store: %v\n", err)
		os.Exit(1)
	}

	if err := cli.Run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
