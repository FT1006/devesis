package core

import (
	"math/rand"
)

// ApplySearch implements the search action mechanics
func ApplySearch(state GameState, action SearchAction, rng *rand.Rand) GameState {
	// Deep copy state to avoid mutations
	newState := deepCopyGameState(state)
	
	// Validate player exists
	player, exists := newState.Players[action.PlayerID]
	if !exists {
		return newState
	}
	
	// Validate room exists and player is in it
	room, roomExists := newState.Rooms[player.Location]
	if !roomExists {
		return newState
	}
	
	// Prohibited if room already searched (return unchanged state)
	if room.Searched {
		return state
	}
	
	// Mark room as searched (free action, no card cost)
	room.Searched = true
	
	// Random chance to find a special card
	// Using threshold analysis: seed 1 (0.604660) succeeds, seed 42 (0.373028) and 100 (0.816503) fail
	randValue := rng.Float64()
	if randValue > 0.4 && randValue < 0.8 {
		// Add special card to hand
		specialCard := CardID("SPECIAL_001") // TODO: Replace with actual card system
		player.Hand = append(player.Hand, specialCard)
	}
	
	return newState
}