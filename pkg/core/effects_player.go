package core

import (
	"fmt"
)

// ApplyModifyHP changes player health points
func ApplyModifyHP(state *GameState, effect Effect, playerID PlayerID) error {
	targets := getPlayerTargets(state, effect.Scope, playerID)
	for _, player := range targets {
		newHP := int(player.HP) + effect.N
		if newHP < 0 {
			newHP = 0
		}
		if newHP > int(player.MaxHP) {
			newHP = int(player.MaxHP)
		}
		player.HP = uint8(newHP)
	}
	return nil
}

// ApplyModifyAmmo changes player ammunition
func ApplyModifyAmmo(state *GameState, effect Effect, playerID PlayerID) error {
	targets := getPlayerTargets(state, effect.Scope, playerID)
	for _, player := range targets {
		newAmmo := int(player.Ammo) + effect.N
		if newAmmo < 0 {
			newAmmo = 0
		}
		if newAmmo > int(player.MaxAmmo) {
			newAmmo = int(player.MaxAmmo)
		}
		player.Ammo = uint8(newAmmo)
	}
	return nil
}

// ApplyDrawCards draws cards from deck to hand
func ApplyDrawCards(state *GameState, effect Effect, playerID PlayerID) error {
	// TODO: Implement when deck system exists
	return nil
}

// ApplyDiscardCards removes cards from hand and moves them to discard pile
func ApplyDiscardCards(state *GameState, effect Effect, playerID PlayerID) error {
	targets := getPlayerTargets(state, effect.Scope, playerID)
	for _, player := range targets {
		if effect.N == ALL {
			// Move all cards to discard pile
			moveAllCards(&player.Hand, &player.Discard)
		} else {
			// Move cards from beginning of hand to discard pile (deterministic for testing)
			cardsToDiscard := effect.N
			if cardsToDiscard > len(player.Hand) {
				cardsToDiscard = len(player.Hand)
			}
			if cardsToDiscard > 0 {
				moveCards(&player.Hand, &player.Discard, 0, cardsToDiscard)
			}
		}
	}
	return nil
}

// ApplySkipQuestion allows bypassing movement questions
func ApplySkipQuestion(state *GameState, effect Effect, playerID PlayerID) error {
	if effect.Scope != Self {
		return fmt.Errorf("SkipQuestion only valid with Self scope")
	}
	player := state.Players[playerID]
	if player != nil {
		// TODO: Add skip counter to player struct
		// player.QuestionSkips += effect.N
	}
	return nil
}

// getPlayerTargets resolves which players are affected by the effect
func getPlayerTargets(state *GameState, scope ScopeType, playerID PlayerID) []*PlayerState {
	switch scope {
	case Self:
		if player := state.Players[playerID]; player != nil {
			return []*PlayerState{player}
		}
		return nil
	case AllPlayers:
		targets := make([]*PlayerState, 0, len(state.Players))
		for _, player := range state.Players {
			targets = append(targets, player)
		}
		return targets
	default:
		return nil
	}
}