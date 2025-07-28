package core

import (
	"math/rand"
)

// moveCards removes `count` cards starting at `start` from src,
// appends them to dst, and returns the removed slice.
// This handles duplicates correctly by working with indices, not values.
func moveCards(src *[]CardID, dst *[]CardID, start, count int) []CardID {
	if count <= 0 || start < 0 || start >= len(*src) {
		return nil
	}
	if start+count > len(*src) {
		count = len(*src) - start
	}

	// Copy cards to be moved
	moved := make([]CardID, count)
	copy(moved, (*src)[start:start+count])

	// Remove from source
	*src = append((*src)[:start], (*src)[start+count:]...)

	// Add to destination
	*dst = append(*dst, moved...)

	return moved
}

// moveAllCards moves all cards from src to dst
func moveAllCards(src *[]CardID, dst *[]CardID) []CardID {
	if len(*src) == 0 {
		return nil
	}
	return moveCards(src, dst, 0, len(*src))
}

// moveCardByIndex moves a single card at the specified index
func moveCardByIndex(src *[]CardID, dst *[]CardID, index int) CardID {
	moved := moveCards(src, dst, index, 1)
	if len(moved) > 0 {
		return moved[0]
	}
	return ""
}

// enforceHandLimitWithDiscard ensures hand doesn't exceed MaxHandSize
// by moving excess cards from the beginning to discard
func enforceHandLimitWithDiscard(hand *[]CardID, discard *[]CardID) {
	if len(*hand) > MaxHandSize {
		overflow := len(*hand) - MaxHandSize
		moveCards(hand, discard, 0, overflow)
	}
}

// shuffleCards shuffles a slice of CardIDs using Fisher-Yates algorithm
func shuffleCards(cards []CardID, rng *rand.Rand) {
	for i := len(cards) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		cards[i], cards[j] = cards[j], cards[i]
	}
}

// shuffleRooms shuffles a slice of RoomIDs using Fisher-Yates algorithm
func shuffleRooms(rooms []RoomID, rng *rand.Rand) {
	for i := len(rooms) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		rooms[i], rooms[j] = rooms[j], rooms[i]
	}
}

// shuffleRoomTypes shuffles a slice of RoomTypes using Fisher-Yates algorithm
func shuffleRoomTypes(roomTypes []RoomType, rng *rand.Rand) {
	for i := len(roomTypes) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		roomTypes[i], roomTypes[j] = roomTypes[j], roomTypes[i]
	}
}

// getEnemyDisplayName returns the display name for an enemy type
func getEnemyDisplayName(enemyType EnemyType) string {
	switch enemyType {
	case InfiniteLoop:
		return "Infinite Loop"
	case StackOverflow:
		return "Stack Overflow"
	case Pythogoras:
		return "Pythogoras"
	default:
		return "Enemy"
	}
}

// getEffectOpName returns the readable name for an effect operation
func getEffectOpName(op EffectOp) string {
	switch op {
	case ModifyHP:
		return "ModifyHP"
	case ModifyAmmo:
		return "ModifyAmmo"
	case DrawCards:
		return "DrawCards"
	case DiscardCards:
		return "DiscardCards"
	case OutOfRam:
		return "OutOfRam"
	case ModifyBugs:
		return "ModifyBugs"
	case RevealRoom:
		return "RevealRoom"
	case CleanRoom:
		return "CleanRoom"
	case SetCorrupted:
		return "SetCorrupted"
	case SpawnEnemy:
		return "SpawnEnemy"
	case MoveEnemies:
		return "MoveEnemies"
	default:
		return "Unknown"
	}
}

// getScopeName returns the readable name for an effect scope
func getScopeName(scope ScopeType) string {
	switch scope {
	case Self:
		return "Self"
	case CurrentRoom:
		return "CurrentRoom"
	case AdjacentRooms:
		return "AdjacentRooms"
	case AllRooms:
		return "AllRooms"
	case RoomWithMostBugs:
		return "RoomWithMostBugs"
	case RoomWithMostEnemies:
		return "RoomWithMostEnemies"
	case AllPlayers:
		return "AllPlayers"
	default:
		return "Unknown"
	}
}

// GetEffectOpName returns the readable name for an effect operation (exported)
func GetEffectOpName(op EffectOp) string {
	return getEffectOpName(op)
}

// GetScopeName returns the readable name for an effect scope (exported)
func GetScopeName(scope ScopeType) string {
	return getScopeName(scope)
}

// selectRoomWithTieBreaking selects the best room from candidates using pathfinding-based tie-breaking
// Returns the room closest to the active player's position in case of ties
func selectRoomWithTieBreaking(state *GameState, candidates []RoomID) RoomID {
	if len(candidates) == 0 {
		return ""
	}
	if len(candidates) == 1 {
		return candidates[0]
	}
	
	// Use active player's current position as anchor for tie-breaking
	anchor := state.Players[state.ActivePlayer].Location
	
	best := candidates[0]
	bestDistance := 999
	
	// Find the room closest to active player
	for _, roomID := range candidates {
		path := CanTraverse(state, PathQuery{From: roomID, To: anchor, MaxSteps: 0})
		distance := 999
		if path.Valid {
			distance = len(path.Path) - 1
		}
		
		if distance < bestDistance {
			best = roomID
			bestDistance = distance
		}
	}
	
	return best
}

// GetRoomWithMostBugs finds the room with highest bug count for bug-related effects
// Returns empty string if no rooms have bugs (> 0)
func GetRoomWithMostBugs(state *GameState) RoomID {
	return getRoomWithMostBugs(state, true)
}

// GetRoomWithMostBugsForSpawn finds the room with highest bug count for enemy spawning
// Returns a room even if it has 0 bugs (fallback to closest to active player)
func GetRoomWithMostBugsForSpawn(state *GameState) RoomID {
	return getRoomWithMostBugs(state, false)
}

// getRoomWithMostBugs is a helper that finds the room with highest bug count
// Returns empty string if requireBugs is true and no rooms have bugs (> 0)
// In case of ties, returns the room closest to active player's position using actual pathfinding
func getRoomWithMostBugs(state *GameState, requireBugs bool) RoomID {
	minBugs := 0
	if requireBugs {
		minBugs = 1 // Must have at least 1 bug to be considered
	}
	
	maxBugs := minBugs - 1 // Start below minimum
	var candidates []RoomID
	
	for id, room := range state.Rooms {
		bugs := int(room.BugMarkers)
		
		if bugs < minBugs {
			continue // Below minimum threshold
		}
		
		if bugs > maxBugs {
			// Found room(s) with more bugs - reset candidates
			maxBugs = bugs
			candidates = []RoomID{id}
		} else if bugs == maxBugs {
			// Tie - add to candidates
			candidates = append(candidates, id)
		}
	}
	
	// Return empty if requireBugs and no rooms have bugs
	if requireBugs && maxBugs < minBugs {
		return ""
	}
	
	// Use tie-breaking helper to select from candidates
	return selectRoomWithTieBreaking(state, candidates)
}

// GetRoomWithMostEnemies finds the room with highest enemy count for OutOfRam effects
// Returns empty string if no enemies exist
func GetRoomWithMostEnemies(state *GameState) RoomID {
	// Count enemies per room
	enemyCount := make(map[RoomID]int)
	for _, enemy := range state.Enemies {
		enemyCount[enemy.Location]++
	}
	
	if len(enemyCount) == 0 {
		return "" // No enemies exist
	}
	
	maxEnemies := 0
	var candidates []RoomID
	
	for roomID, count := range enemyCount {
		if count > maxEnemies {
			// Found room(s) with more enemies - reset candidates
			maxEnemies = count
			candidates = []RoomID{roomID}
		} else if count == maxEnemies {
			// Tie - add to candidates
			candidates = append(candidates, roomID)
		}
	}
	
	// Use tie-breaking helper to select from candidates
	return selectRoomWithTieBreaking(state, candidates)
}

// drawCards draws up to `count` cards from deck to hand, handling deck reshuffling
// Returns the number of cards actually drawn
func drawCards(hand *[]CardID, deck *[]CardID, discard *[]CardID, count int, rng *rand.Rand) int {
	cardsDrawn := 0
	
	for i := 0; i < count; i++ {
		// If deck is empty, shuffle discard pile into deck
		if len(*deck) == 0 && len(*discard) > 0 {
			// Move all discard cards to deck
			moveAllCards(discard, deck)
			
			// Shuffle the new deck
			shuffleCards(*deck, rng)
		}
		
		// Draw from deck if available
		if len(*deck) > 0 {
			moveCardByIndex(deck, hand, 0) // Move first card from deck to hand
			cardsDrawn++
		} else {
			// No more cards available in deck or discard
			break
		}
	}
	
	return cardsDrawn
}