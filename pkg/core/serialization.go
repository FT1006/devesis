package core

import (
	"encoding/json"
	"fmt"
)

// SaveGameState serializes game state to JSON
func SaveGameState(state *GameState) ([]byte, error) {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal game state: %w", err)
	}
	return data, nil
}

// LoadGameState deserializes game state from JSON
func LoadGameState(data []byte) (*GameState, error) {
	var state GameState
	err := json.Unmarshal(data, &state)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal game state: %w", err)
	}
	return &state, nil
}