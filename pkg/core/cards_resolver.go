package core

// ApplyCardEffects executes all effects from a card on the game state
func ApplyCardEffects(state GameState, card Card, playerID PlayerID, log *EffectLog) GameState {
	// Validate all effects before applying any
	if err := ValidateCard(card); err != nil {
		// Return original state if card invalid
		return state
	}

	// Apply effects in order  
	newState := deepCopyGameState(state)
	for _, effect := range card.Effects {
		if err := applyEffect(&newState, effect, playerID, log); err != nil {
			// Error logging is handled centrally in ApplyEffect
			// Return original state if any effect fails
			return state
		}
	}

	return newState
}

// applyEffect executes a single effect on the game state using the centralized handler
func applyEffect(state *GameState, effect Effect, playerID PlayerID, log *EffectLog) error {
	return ApplyEffect(state, effect, playerID, log)
}