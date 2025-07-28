package core

import (
	"testing"
)

func TestDrawPhaseRound1RefillsToFive(t *testing.T) {
	gs := GameState{
		ActivePlayer: "P1",
		RandSeed:    42,
		Round:       1, // First turn
		Players: map[PlayerID]*PlayerState{
			"P1": {
				ID:   "P1",
				Hand: []CardID{"C1", "C2"}, // 2 cards in hand
				Deck: []CardID{"C3", "C4", "C5", "C6", "C7"}, // 5 cards in deck
			},
		},
	}
	
	DrawPhase(&gs)
	
	// Should draw 3 cards to reach 5 total (round 1 behavior)
	player := gs.Players["P1"]
	if len(player.Hand) != 5 {
		t.Errorf("expected 5 cards in hand on round 1, got %d", len(player.Hand))
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

func TestDrawPhaseSubsequentRoundsDrawTwo(t *testing.T) {
	gs := GameState{
		ActivePlayer: "P1",
		RandSeed:    42,
		Round:       2, // Subsequent turn
		Players: map[PlayerID]*PlayerState{
			"P1": {
				ID:   "P1",
				Hand: []CardID{"C1", "C2", "C3"}, // 3 cards in hand
				Deck: []CardID{"C4", "C5", "C6", "C7"}, // 4 cards in deck
			},
		},
	}
	
	DrawPhase(&gs)
	
	// Should draw exactly 2 cards (subsequent round behavior)
	player := gs.Players["P1"]
	if len(player.Hand) != 5 {
		t.Errorf("expected 5 cards in hand (3 + 2 drawn), got %d", len(player.Hand))
	}
	
	// Deck should have 2 cards left
	if len(player.Deck) != 2 {
		t.Errorf("expected 2 cards in deck, got %d", len(player.Deck))
	}
}

func TestDrawPhaseShufflesDiscardWhenDeckEmpty(t *testing.T) {
	gs := GameState{
		ActivePlayer: "P1",
		RandSeed:    42,
		Round:       1, // Round 1: refill to 5
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
	
	// Should draw 4 cards to reach 5 total (round 1 behavior)
	player := gs.Players["P1"]
	if len(player.Hand) != 5 {
		t.Errorf("expected 5 cards in hand (1 + 4 drawn), got %d", len(player.Hand))
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

func TestDrawPhaseRound1NoDrawWhenHandFull(t *testing.T) {
	gs := GameState{
		ActivePlayer: "P1",
		RandSeed:    42,
		Round:       1, // Round 1: refill to 5
		Players: map[PlayerID]*PlayerState{
			"P1": {
				ID:   "P1",
				Hand: []CardID{"C1", "C2", "C3", "C4", "C5"}, // Already 5 cards
				Deck: []CardID{"C6", "C7"},
			},
		},
	}
	
	DrawPhase(&gs)
	
	// Should not draw any cards on round 1 when already at 5
	player := gs.Players["P1"]
	if len(player.Hand) != 5 {
		t.Errorf("expected 5 cards in hand (no draw needed), got %d", len(player.Hand))
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
	
	// Should draw all available cards (tries to draw 2, gets all 2 available)
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

func TestDrawPhaseEnforcesHandLimit(t *testing.T) {
	gs := GameState{
		ActivePlayer: "P1",
		RandSeed:    42,
		Round:       2, // Subsequent turn: draw 2 cards
		Players: map[PlayerID]*PlayerState{
			"P1": {
				ID:   "P1",
				Hand: []CardID{"C1", "C2", "C3", "C4", "C5", "C6"}, // Already at max (6 cards)
				Deck: []CardID{"C7", "C8", "C9", "C10"}, // More cards available
				Discard: []CardID{}, // Empty discard
			},
		},
	}
	
	DrawPhase(&gs)
	
	// Hand should remain at max size (6)
	player := gs.Players["P1"]
	if len(player.Hand) != 6 {
		t.Errorf("expected 6 cards in hand (max limit), got %d", len(player.Hand))
	}
	
	// Excess cards should go to discard (2 cards drawn, both should go to discard)
	if len(player.Discard) != 2 {
		t.Errorf("expected 2 cards in discard (excess from drawing), got %d", len(player.Discard))
	}
	
	// Deck should have 2 cards left (4 - 2 drawn)
	if len(player.Deck) != 2 {
		t.Errorf("expected 2 cards left in deck, got %d", len(player.Deck))
	}
}

func TestDrawPhaseRound1HandLimitEnforcement(t *testing.T) {
	gs := GameState{
		ActivePlayer: "P1",
		RandSeed:    42,
		Round:       1, // Round 1: target 5 cards
		Players: map[PlayerID]*PlayerState{
			"P1": {
				ID:   "P1",
				Hand: []CardID{"C1", "C2", "C3", "C4"}, // 4 cards in hand
				Deck: []CardID{"C5", "C6", "C7", "C8", "C9", "C10"}, // 6 cards in deck
				Discard: []CardID{}, // Empty discard
			},
		},
	}
	
	DrawPhase(&gs)
	
	// Should draw 1 card to reach target of 5 (round 1 behavior)
	player := gs.Players["P1"]
	if len(player.Hand) != 5 {
		t.Errorf("expected 5 cards in hand on round 1, got %d", len(player.Hand))
	}
	
	// Discard should remain empty (no overflow)
	if len(player.Discard) != 0 {
		t.Errorf("expected empty discard, got %d cards", len(player.Discard))
	}
	
	// Deck should have 5 cards left (6 - 1 drawn)
	if len(player.Deck) != 5 {
		t.Errorf("expected 5 cards left in deck, got %d", len(player.Deck))
	}
}

func TestDrawPhaseExcessiveRound1Draw(t *testing.T) {
	gs := GameState{
		ActivePlayer: "P1",
		RandSeed:    42,
		Round:       1, // Round 1: target 5 cards, but will try to draw many
		Players: map[PlayerID]*PlayerState{
			"P1": {
				ID:   "P1",
				Hand: []CardID{}, // Empty hand (will try to draw 5)
				Deck: []CardID{"C1", "C2", "C3", "C4", "C5", "C6", "C7", "C8"}, // 8 cards in deck
				Discard: []CardID{}, // Empty discard
			},
		},
	}
	
	DrawPhase(&gs)
	
	// Should draw 5 cards and stop (round 1 target)
	player := gs.Players["P1"]
	if len(player.Hand) != 5 {
		t.Errorf("expected 5 cards in hand on round 1, got %d", len(player.Hand))
	}
	
	// Discard should remain empty (no overflow since target was 5)
	if len(player.Discard) != 0 {
		t.Errorf("expected empty discard, got %d cards", len(player.Discard))
	}
	
	// Deck should have 3 cards left (8 - 5 drawn)
	if len(player.Deck) != 3 {
		t.Errorf("expected 3 cards left in deck, got %d", len(player.Deck))
	}
}