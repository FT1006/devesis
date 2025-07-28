package main

import (
	"bufio"
	"fmt"
	"strings"
	"time"

	"github.com/spaceship/devesis/pkg/core"
)

// PreviewAndConfirm shows what effects will happen and asks for confirmation
func (g *GameManager) PreviewAndConfirm(action core.Action, reader *bufio.Reader) bool {
	// Create a copy of the state for preview
	previewState := core.DeepCopyGameState(*g.state)
	
	// Create a fresh log for preview
	previewLog := core.NewEffectLog()
	
	// Run the action on the copy to see what would happen
	core.Apply(previewState, action, previewLog)
	
	// Only show preview if there are effects to show
	if !previewLog.IsEmpty() {
		fmt.Println("\n— Effects Preview —")
		previewLog.PrintBulk()
		
		// Ask for confirmation
		fmt.Print("\nProceed? (y/n) > ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return false
		}
		
		response := strings.TrimSpace(strings.ToLower(input))
		return response == "y" || response == "yes"
	}
	
	// If no effects, proceed automatically
	return true
}

// ResolveWithLogging applies the action and streams the effects
func (g *GameManager) ResolveWithLogging(action core.Action) {
	// Create a fresh log for resolution
	resolveLog := core.NewEffectLog()
	
	// Apply the action with logging
	newState := core.Apply(*g.state, action, resolveLog)
	
	// Update the game state
	g.state = &newState
	
	// Stream the effects if any occurred
	if !resolveLog.IsEmpty() {
		fmt.Println("\n— Resolve —")
		resolveLog.StreamLines(1000 * time.Millisecond) // 1 second delay for clear line-by-line visibility
	}
}

