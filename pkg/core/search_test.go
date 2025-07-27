package core

import (
	"math/rand"
	"testing"
)

func TestSearchAction_MarksRoomSearched(t *testing.T) {
	state := newSearchTestGameState()
	action := SearchAction{PlayerID: "P1"}
	rng := rand.New(rand.NewSource(42))
	
	result := ApplySearch(state, action, rng)
	
	room := result.Rooms["R12"]
	if !room.Searched {
		t.Error("Search should mark room as searched")
	}
}

func TestSearchAction_NoCardCost(t *testing.T) {
	state := newSearchTestGameState()
	action := SearchAction{PlayerID: "P1"}
	rng := rand.New(rand.NewSource(42))
	originalHandSize := len(state.Players["P1"].Hand)
	
	result := ApplySearch(state, action, rng)
	
	player := result.Players["P1"]
	if len(player.Hand) != originalHandSize {
		t.Error("Search should not consume cards")
	}
}

func TestSearchAction_ProhibitedAlreadySearched(t *testing.T) {
	state := newSearchTestGameState()
	state.Rooms["R12"].Searched = true
	action := SearchAction{PlayerID: "P1"}
	rng := rand.New(rand.NewSource(42))
	
	result := ApplySearch(state, action, rng)
	
	// Should return unchanged state (prohibited action)
	if result.Round != state.Round {
		t.Error("Searching already searched room should be prohibited - no state change")
	}
}

func TestSearchAction_FindsSpecialCard(t *testing.T) {
	state := newSearchTestGameState()
	action := SearchAction{PlayerID: "P1"}
	rng := rand.New(rand.NewSource(1)) // Seed for success
	originalHandSize := len(state.Players["P1"].Hand)
	
	result := ApplySearch(state, action, rng)
	
	player := result.Players["P1"]
	if len(player.Hand) != originalHandSize+1 {
		t.Error("Successful search should add special card to hand")
	}
}

func TestSearchAction_FailsToFindCard(t *testing.T) {
	state := newSearchTestGameState()
	action := SearchAction{PlayerID: "P1"}
	rng := rand.New(rand.NewSource(100)) // Seed for failure
	originalHandSize := len(state.Players["P1"].Hand)
	
	result := ApplySearch(state, action, rng)
	
	player := result.Players["P1"]
	if len(player.Hand) != originalHandSize {
		t.Error("Failed search should not change hand size")
	}
	
	// Room should still be marked searched even on failure
	room := result.Rooms["R12"]
	if !room.Searched {
		t.Error("Failed search should still mark room as searched")
	}
}

func TestSearchAction_WorksInAnyRoomType(t *testing.T) {
	// Test that search works regardless of room type
	roomTypes := []RoomType{Empty, AmmoCache, MedBay, CleanRoomType, EnemySpawn}
	
	for _, roomType := range roomTypes {
		state := newSearchTestGameState()
		state.Rooms["R12"].Type = roomType
		action := SearchAction{PlayerID: "P1"}
		rng := rand.New(rand.NewSource(1)) // Success seed
		
		result := ApplySearch(state, action, rng)
		
		if !result.Rooms["R12"].Searched {
			t.Errorf("Search should work in %v room type", roomType)
		}
	}
}

// Helper for search tests
func newSearchTestGameState() GameState {
	return GameState{
		Round: 1, // For testing prohibited action
		Rooms: map[RoomID]*RoomState{
			"R12": {ID: "R12", Type: Empty, Searched: false},
		},
		Players: map[PlayerID]*PlayerState{
			"P1": {
				ID:       "P1",
				Location: "R12",
				Hand:     []CardID{"CARD_1", "CARD_2"}, // 2 cards
				HP:       8,
				MaxHP:    10,
				Ammo:     3,
				MaxAmmo:  6,
			},
		},
		Enemies: map[EnemyID]*Enemy{},
	}
}