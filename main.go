package main

import (
	"fmt"
	"os"
	"task-cli/internal/cli"
	"task-cli/internal/task"
)

func main() {
	if err := task.InitTaskStore(); err != nil {
		fmt.Printf("Error initializing task store: %v\n", err)
		os.Exit(1)
	}

	cli.Run(os.Args[1:])
}
