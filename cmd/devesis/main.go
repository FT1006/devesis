package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	game := NewGameManager()

	// Initialize new game or load saved state
	if err := game.Initialize(); err != nil {
		fmt.Printf("Failed to initialize game: %v\n", err)
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)

	// Main game loop with 4-phase structure
	quitByPlayer := false
	for !game.IsGameOver() {
		// Phase 1: Draw Phase
		game.ExecuteDrawPhase()
		
		// Phase 2: Player Phase (action-driven commands)
		if err := game.ExecutePlayerPhase(reader); err != nil {
			if err.Error() == "quit" {
				quitByPlayer = true
				break
			}
			fmt.Printf("Player phase error: %v\n", err)
		}
		
		// Phase 3: Event Phase
		game.ExecuteEventPhase()
		
		// Phase 4: Round Maintenance
		game.ExecuteRoundMaintenance()
		
		// Check end conditions after each round
		if ended, win := game.CheckEndConditions(); ended {
			game.DisplayGameResult(win)
			return
		}
	}

	if !quitByPlayer {
		// Game over due to normal end conditions
		game.DisplayGameOver()
	}
}