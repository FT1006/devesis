package core

import (
	"fmt"
)

// ValidateCard checks if all effects in a card are valid
func ValidateCard(card Card) error {
	for i, effect := range card.Effects {
		if err := ValidateEffect(effect, card.Source); err != nil {
			return fmt.Errorf("effect %d: %w", i, err)
		}
	}
	return nil
}

// ValidateEffect checks if a single effect is valid
func ValidateEffect(effect Effect, source EffectSource) error {
	// Check Op-Scope compatibility
	if !isValidOpScope(effect.Op, effect.Scope) {
		return fmt.Errorf("invalid op-scope combination: %v with %v", effect.Op, effect.Scope)
	}

	// Check N value ranges
	if !isValidNValue(effect.Op, effect.N) {
		return fmt.Errorf("invalid N value %d for op %v", effect.N, effect.Op)
	}

	// Check phase compatibility
	if !isValidPhaseOp(source, effect.Op) {
		return fmt.Errorf("op %v not allowed in phase %v", effect.Op, source)
	}

	return nil
}

// isValidOpScope checks if an operation is compatible with a scope
func isValidOpScope(op EffectOp, scope ScopeType) bool {
	validScopes := map[EffectOp][]ScopeType{
		ModifyHP:     {Self, AllPlayers},
		ModifyAmmo:   {Self, AllPlayers},
		DrawCards:    {Self, AllPlayers},
		DiscardCards: {Self, AllPlayers},
		OutOfRam:     {RoomWithMostEnemies},
		ModifyBugs:   {CurrentRoom, AdjacentRooms, AllRooms, RoomWithMostBugs},
		RevealRoom:   {CurrentRoom, AdjacentRooms, AllRooms},
		CleanRoom:    {CurrentRoom, AdjacentRooms, AllRooms},
		SetCorrupted: {CurrentRoom, AdjacentRooms, AllRooms, RoomWithMostBugs},
		SpawnEnemy:   {CurrentRoom, RoomWithMostBugs},
		MoveEnemies:  {AllRooms}, // Enemy movement affects all enemies
	}

	scopes, exists := validScopes[op]
	if !exists {
		return false
	}

	for _, validScope := range scopes {
		if scope == validScope {
			return true
		}
	}
	return false
}

// isValidNValue checks if N value is in valid range for the operation
func isValidNValue(op EffectOp, n int) bool {
	switch op {
	case ModifyHP, ModifyAmmo:
		return n >= -10 && n <= 10
	case ModifyBugs:
		return n == ALL || (n >= -9 && n <= 9)
	case DrawCards, DiscardCards:
		return n >= 1 && n <= 5
	case SpawnEnemy:
		return n >= 1 && n <= 3
	case OutOfRam:
		return n == 1
	case RevealRoom, CleanRoom:
		return n == 1
	case SetCorrupted:
		return n == 0 || n == 1
	case MoveEnemies:
		return n >= 1 && n <= 3 // 1-3 steps movement
	default:
		return false
	}
}

// isValidPhaseOp checks if operation is allowed in the given phase
func isValidPhaseOp(source EffectSource, op EffectOp) bool {
	switch source {
	case SrcAction:
		return true // All ops allowed in action phase
	case SrcEvent:
		allowedOps := []EffectOp{SpawnEnemy, ModifyBugs, SetCorrupted, CleanRoom, RevealRoom, MoveEnemies}
		for _, allowedOp := range allowedOps {
			if op == allowedOp {
				return true
			}
		}
		return false
	case SrcSpecial:
		return true // Special cards allow all ops (like action cards)
	default:
		return false
	}
}