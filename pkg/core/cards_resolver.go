package core

import "fmt"

// ApplyCardEffects executes all effects from a card on the game state
func ApplyCardEffects(state GameState, card Card, playerID PlayerID) GameState {
	// Validate all effects before applying any
	if err := ValidateCard(card); err != nil {
		// Return original state if card invalid
		return state
	}

	// Apply effects in order  
	newState := deepCopyGameState(state)
	for _, effect := range card.Effects {
		if err := applyEffect(&newState, effect, playerID); err != nil {
			// Return original state if any effect fails
			return state
		}
	}

	return newState
}

// applyEffect executes a single effect on the game state
func applyEffect(state *GameState, effect Effect, playerID PlayerID) error {
	switch effect.Op {
	case ModifyHP:
		return ApplyModifyHP(state, effect, playerID)
	case ModifyAmmo:
		return ApplyModifyAmmo(state, effect, playerID)
	case DrawCards:
		return ApplyDrawCards(state, effect, playerID)
	case DiscardCards:
		return ApplyDiscardCards(state, effect, playerID)
	case SkipQuestion:
		return ApplySkipQuestion(state, effect, playerID)
	case ModifyBugs:
		return ApplyModifyBugs(state, effect, playerID)
	case RevealRoom:
		return ApplyRevealRoom(state, effect, playerID)
	case CleanRoom:
		return ApplyCleanRoom(state, effect, playerID)
	case SetCorrupted:
		return ApplySetCorrupted(state, effect, playerID)
	case SpawnEnemy:
		return ApplySpawnEnemy(state, effect, playerID)
	default:
		return fmt.Errorf("unknown effect op: %v", effect.Op)
	}
}