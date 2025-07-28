package core

import (
	"fmt"
)

// ApplyModifyHP changes player health points
func ApplyModifyHP(state *GameState, effect Effect, playerID PlayerID, log *EffectLog) error {
	targets := getPlayerTargets(state, effect.Scope, playerID)
	for _, player := range targets {
		oldHP := player.HP
		newHP := int(player.HP) + effect.N
		if newHP < 0 {
			newHP = 0
		}
		if newHP > int(player.MaxHP) {
			newHP = int(player.MaxHP)
		}
		player.HP = uint8(newHP)
		
		if oldHP != player.HP {
			if effect.N > 0 {
				log.Add("🩹 %s HP: %d → %d (+%d)", player.ID, oldHP, player.HP, effect.N)
			} else {
				log.Add("💔 %s HP: %d → %d (%d)", player.ID, oldHP, player.HP, effect.N)
			}
		}
	}
	return nil
}

// ApplyModifyAmmo changes player ammunition
func ApplyModifyAmmo(state *GameState, effect Effect, playerID PlayerID, log *EffectLog) error {
	targets := getPlayerTargets(state, effect.Scope, playerID)
	for _, player := range targets {
		oldAmmo := player.Ammo
		newAmmo := int(player.Ammo) + effect.N
		if newAmmo < 0 {
			newAmmo = 0
		}
		if newAmmo > int(player.MaxAmmo) {
			newAmmo = int(player.MaxAmmo)
		}
		player.Ammo = uint8(newAmmo)
		
		if oldAmmo != player.Ammo {
			if effect.N > 0 {
				log.Add("🔫 %s ammo: %d → %d (+%d)", player.ID, oldAmmo, player.Ammo, effect.N)
			} else {
				log.Add("🔫 %s ammo: %d → %d (%d)", player.ID, oldAmmo, player.Ammo, effect.N)
			}
		}
	}
	return nil
}

// ApplyDrawCards draws cards from deck to hand
func ApplyDrawCards(state *GameState, effect Effect, playerID PlayerID, log *EffectLog) error {
	targets := getPlayerTargets(state, effect.Scope, playerID)
	
	// Create RNG for deterministic card drawing
	rng := GetGameRNG(state)
	
	for _, player := range targets {
		oldHandSize := len(player.Hand)
		drawCards(&player.Hand, &player.Deck, &player.Discard, effect.N, rng)
		
		// Enforce hand limit if drawing would exceed it
		enforceHandLimitWithDiscard(&player.Hand, &player.Discard)
		
		cardsDrawn := len(player.Hand) - oldHandSize
		if cardsDrawn > 0 {
			log.Add("🃏 %s draws %d cards", player.ID, cardsDrawn)
		}
	}
	
	return nil
}

// ApplyDiscardCards removes cards from hand and moves them to discard pile
func ApplyDiscardCards(state *GameState, effect Effect, playerID PlayerID, log *EffectLog) error {
	targets := getPlayerTargets(state, effect.Scope, playerID)
	for _, player := range targets {
		oldHandSize := len(player.Hand)
		if effect.N == ALL {
			// Move all cards to discard pile
			moveAllCards(&player.Hand, &player.Discard)
			log.Add("🗑️ %s discards all %d cards", player.ID, oldHandSize)
		} else {
			// Move cards from beginning of hand to discard pile (deterministic for testing)
			cardsToDiscard := effect.N
			if cardsToDiscard > len(player.Hand) {
				cardsToDiscard = len(player.Hand)
			}
			if cardsToDiscard > 0 {
				moveCards(&player.Hand, &player.Discard, 0, cardsToDiscard)
				log.Add("🗑️ %s discards %d cards", player.ID, cardsToDiscard)
			}
		}
	}
	return nil
}

// ApplySkipQuestion allows bypassing movement questions
func ApplySkipQuestion(state *GameState, effect Effect, playerID PlayerID, log *EffectLog) error {
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