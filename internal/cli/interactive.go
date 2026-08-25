package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func RunInteractive() {
	scanner := bufio.NewScanner(os.Stdin)

	printWelcome()

	for {
		fmt.Print("\nmycli> ")

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())

		if input == "" {
			continue
		}

		if input == "exit" || input == "quit" || input == "q" {
			fmt.Println("Goodbye! 👋")
			break
		}

		executeCommand(input)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
	}
}

func printWelcome() {
	fmt.Println(`
╔═══════════════════════════════════════════╗
║          MY CLI - Task Manager            ║
╠═══════════════════════════════════════════╣
║  Type 'help' for available commands       ║
║  Type 'exit' or 'quit' to quit            ║
╚═══════════════════════════════════════════╝`)
}
