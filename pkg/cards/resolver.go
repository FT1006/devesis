package cards

import (
	"fmt"
	"github.com/spaceship/devesis/pkg/cards/effects"
	"github.com/spaceship/devesis/pkg/core"
)

// ApplyEffects executes all effects from a card on the game state
func ApplyEffects(state *core.GameState, card effects.Card, playerID core.PlayerID) (*core.GameState, error) {
	// Validate all effects before applying any
	if err := effects.ValidateCard(card); err != nil {
		return state, fmt.Errorf("invalid card: %w", err)
	}

	// Apply effects in order
	newState := *state
	for i, effect := range card.Effects {
		if err := applyEffect(&newState, effect, playerID); err != nil {
			return state, fmt.Errorf("effect %d failed: %w", i, err)
		}
	}

	return &newState, nil
}

// applyEffect executes a single effect on the game state
func applyEffect(state *core.GameState, effect effects.Effect, playerID core.PlayerID) error {
	switch effect.Op {
	case ModifyHP:
		return effects.ApplyModifyHP(state, effect, playerID)
	case ModifyAmmo:
		return effects.ApplyModifyAmmo(state, effect, playerID)
	case DrawCards:
		return effects.ApplyDrawCards(state, effect, playerID)
	case DiscardCards:
		return effects.ApplyDiscardCards(state, effect, playerID)
	case SkipQuestion:
		return effects.ApplySkipQuestion(state, effect, playerID)
	case ModifyBugs:
		return effects.ApplyModifyBugs(state, effect, playerID)
	case RevealRoom:
		return effects.ApplyRevealRoom(state, effect, playerID)
	case CleanRoom:
		return effects.ApplyCleanRoom(state, effect, playerID)
	case SetCorrupted:
		return effects.ApplySetCorrupted(state, effect, playerID)
	case SpawnEnemy:
		return effects.ApplySpawnEnemy(state, effect, playerID)
	default:
		return fmt.Errorf("unknown effect op: %v", effect.Op)
	}
}