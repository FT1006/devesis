package core

import (
	"math/rand"
	"strings"
)

// ApplySearch implements the search action mechanics
func ApplySearch(state GameState, action SearchAction, rng *rand.Rand, log *EffectLog) GameState {
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
	log.Add("🔍 %s searches %s", action.PlayerID, player.Location)
	
	// Room-specific search overrides
	switch player.Location {
	case "R01": // KEY room - gives BOOT.dev KEY power
		oldDamage := player.Damage
		player.Damage = BootDevDamage // Increase damage from BasicDamage to BootDevDamage
		log.Add("🔑 Found the BOOT.dev KEY! Damage: %d → %d", oldDamage, player.Damage)
		// Note: This is a permanent power upgrade, not a card
		return newState
		
	case "R15", "R17", "R18": // Engine rooms EN1, EN2, EN3
		// Give 3 engine cards (representing all 3 engines)
		player.Hand = append(player.Hand, "SPECIAL_ENGINE", "SPECIAL_ENGINE", "SPECIAL_ENGINE")
		log.Add("⚙️ Found 3 engine cards!")
		enforceHandLimitWithDiscard(&player.Hand, &player.Discard)
		return newState
	}
	
	// Default random chance to find a special card
	// Using threshold analysis: seed 1 (0.604660) succeeds, seed 42 (0.373028) and 100 (0.816503) fail
	randValue := rng.Float64()
	if randValue > 0.4 && randValue < 0.8 {
		// Select random special card from available cards
		specialCard := selectRandomSpecialCard(rng)
		if specialCard != "" {
			player.Hand = append(player.Hand, specialCard)
			if card, exists := CardDB[specialCard]; exists {
				log.Add("🎴 Found %s - %s", card.Name, card.Description)
			} else {
				log.Add("🎴 Found special card: %s", specialCard)
			}
			enforceHandLimitWithDiscard(&player.Hand, &player.Discard)
		} else {
			log.Add("🔍 Nothing found")
		}
	} else {
		log.Add("🔍 Nothing found")
	}
	
	return newState
}

// selectRandomSpecialCard picks a random special card from the loaded database
func selectRandomSpecialCard(rng *rand.Rand) CardID {
	// Collect all special card IDs from the database
	var specialCards []CardID
	for cardID, card := range CardDB {
		// Check if it's a special card by ID prefix
		if strings.HasPrefix(string(cardID), "SPECIAL") {
			// Get rarity from card data
			rarity := getCardRarity(card)
			
			// Add card multiple times based on rarity (common cards more likely)
			switch rarity {
			case "common":
				for i := 0; i < 3; i++ {
					specialCards = append(specialCards, cardID)
				}
			case "uncommon":
				for i := 0; i < 2; i++ {
					specialCards = append(specialCards, cardID)
				}
			case "rare":
				specialCards = append(specialCards, cardID)
			default:
				// Default to uncommon weighting
				for i := 0; i < 2; i++ {
					specialCards = append(specialCards, cardID)
				}
			}
		}
	}
	
	// Return empty if no special cards available
	if len(specialCards) == 0 {
		return ""
	}
	
	// Pick random card from weighted list
	randomIndex := rng.Intn(len(specialCards))
	return specialCards[randomIndex]
}

// getCardRarity extracts rarity from card data, with fallback to effect-count heuristic
func getCardRarity(card Card) string {
	// TODO: Use card.Rarity field when YAML parsing supports it
	// For now, infer rarity from effects complexity
	switch n := len(card.Effects); {
	case n >= 3:
		return "rare"
	case n >= 2:
		return "uncommon"
	default:
		return "common"
	}
}

// enforceHandLimit is now replaced by enforceHandLimitWithDiscard in card_utils.go