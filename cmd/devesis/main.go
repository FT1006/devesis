package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	game := NewGameManager()

	// Initialize new game or load saved state
	if err := game.Initialize(); err != nil {
		fmt.Printf("Failed to initialize game: %v\n", err)
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)

	// Main game loop
	for !game.IsGameOver() {
		// Display game status
		game.DisplayStatus()

		// Show command prompt
		fmt.Print("> ")

		// Read command
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		// Parse and execute command
		args := strings.Fields(input)
		if len(args) == 0 {
			continue
		}

		command := strings.ToLower(args[0])
		commandArgs := args[1:]

		// Execute command
		if err := game.ExecuteCommand(command, commandArgs); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}

	// Game over
	game.DisplayGameOver()
}