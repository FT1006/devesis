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