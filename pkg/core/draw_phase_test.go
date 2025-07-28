package core

import (
	"testing"
)

func TestDrawPhaseRefillsHandToFive(t *testing.T) {
	gs := GameState{
		ActivePlayer: "P1",
		RandSeed:    42,
		Round:       1,
		Players: map[PlayerID]*PlayerState{
			"P1": {
				ID:   "P1",
				Hand: []CardID{"C1", "C2"}, // 2 cards in hand
				Deck: []CardID{"C3", "C4", "C5", "C6", "C7"}, // 5 cards in deck
			},
		},
	}
	
	DrawPhase(&gs)
	
	// Should draw 3 cards to reach 5
	player := gs.Players["P1"]
	if len(player.Hand) != 5 {
		t.Errorf("expected 5 cards in hand, got %d", len(player.Hand))
	}
	
	// Deck should have 2 cards left
	if len(player.Deck) != 2 {
		t.Errorf("expected 2 cards in deck, got %d", len(player.Deck))
	}
	
	// Should set phase to player
	if gs.Phase != "player" {
		t.Errorf("expected phase 'player', got %s", gs.Phase)
	}
	
	// Should set actions to 2
	if gs.ActionsLeft != 2 {
		t.Errorf("expected 2 actions, got %d", gs.ActionsLeft)
	}
}

func TestDrawPhaseShufflesDiscardWhenDeckEmpty(t *testing.T) {
	gs := GameState{
		ActivePlayer: "P1",
		RandSeed:    42,
		Round:       1,
		Players: map[PlayerID]*PlayerState{
			"P1": {
				ID:      "P1",
				Hand:    []CardID{"C1"}, // 1 card in hand
				Deck:    []CardID{"C2"}, // Only 1 card in deck
				Discard: []CardID{"C3", "C4", "C5", "C6", "C7"}, // 5 cards in discard
			},
		},
	}
	
	DrawPhase(&gs)
	
	// Should draw to 5 cards
	player := gs.Players["P1"]
	if len(player.Hand) != 5 {
		t.Errorf("expected 5 cards in hand, got %d", len(player.Hand))
	}
	
	// Discard should be empty (shuffled into deck)
	if len(player.Discard) != 0 {
		t.Errorf("expected empty discard, got %d cards", len(player.Discard))
	}
	
	// Deck should have remaining cards
	if len(player.Deck) != 2 {
		t.Errorf("expected 2 cards left in deck, got %d", len(player.Deck))
	}
}

func TestDrawPhaseNoDrawWhenHandFull(t *testing.T) {
	gs := GameState{
		ActivePlayer: "P1",
		RandSeed:    42,
		Round:       1,
		Players: map[PlayerID]*PlayerState{
			"P1": {
				ID:   "P1",
				Hand: []CardID{"C1", "C2", "C3", "C4", "C5"}, // Already 5 cards
				Deck: []CardID{"C6", "C7"},
			},
		},
	}
	
	DrawPhase(&gs)
	
	// Should not draw any cards
	player := gs.Players["P1"]
	if len(player.Hand) != 5 {
		t.Errorf("expected 5 cards in hand, got %d", len(player.Hand))
	}
	
	// Deck should be unchanged
	if len(player.Deck) != 2 {
		t.Errorf("expected 2 cards in deck, got %d", len(player.Deck))
	}
}

func TestDrawPhaseHandlesInsufficientCards(t *testing.T) {
	gs := GameState{
		ActivePlayer: "P1",
		RandSeed:    42,
		Round:       1,
		Players: map[PlayerID]*PlayerState{
			"P1": {
				ID:      "P1",
				Hand:    []CardID{"C1"}, // 1 card in hand
				Deck:    []CardID{"C2"}, // Only 1 card in deck
				Discard: []CardID{"C3"}, // Only 1 card in discard
			},
		},
	}
	
	DrawPhase(&gs)
	
	// Should draw all available cards (3 total)
	player := gs.Players["P1"]
	if len(player.Hand) != 3 {
		t.Errorf("expected 3 cards in hand, got %d", len(player.Hand))
	}
	
	// Both deck and discard should be empty
	if len(player.Deck) != 0 {
		t.Errorf("expected empty deck, got %d cards", len(player.Deck))
	}
	if len(player.Discard) != 0 {
		t.Errorf("expected empty discard, got %d cards", len(player.Discard))
	}
}