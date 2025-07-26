package core

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestApplyReturnsNewGameState(t *testing.T) {
	state := newTestGameState()
	action := MoveAction{PlayerID: "P1", To: "R02"}
	
	result := Apply(state, action)
	
	if reflect.DeepEqual(state, result) {
		t.Fatal("expected state to change")
	}
}

func TestApplyNeverMutatesOriginalState(t *testing.T) {
	state := newTestGameState()
	originalLocation := state.Players["P1"].Location
	action := MoveAction{PlayerID: "P1", To: "R02"}
	
	Apply(state, action)
	
	if state.Players["P1"].Location != originalLocation {
		t.Error("Apply mutated original state")
	}
}

func TestApplyWithInvalidActionReturnsStateUnchanged(t *testing.T) {
	state := newTestGameState()
	action := &invalidAction{}
	
	result := Apply(state, action)
	
	if result.Players["P1"].Location != state.Players["P1"].Location {
		t.Error("Invalid action should not change state")
	}
}

func TestMoveActionUpdatesPlayerLocation(t *testing.T) {
	state := newTestGameState()
	action := MoveAction{PlayerID: "P1", To: "R07"} // Adjacent to R12
	
	result := Apply(state, action)
	
	if result.Players["P1"].Location != "R07" {
		t.Errorf("Expected player location R07, got %s", result.Players["P1"].Location)
	}
}

func TestMoveActionToNonAdjacentRoomFails(t *testing.T) {
	state := newTestGameState()
	action := MoveAction{PlayerID: "P1", To: "R01"} // Non-adjacent to R12
	
	result := Apply(state, action)
	
	if result.Players["P1"].Location != state.Players["P1"].Location {
		t.Error("Non-adjacent move should fail")
	}
}

func TestSearchActionMarksRoomSearched(t *testing.T) {
	state := newTestGameState()
	state.Players["P1"].Hand = []Card{{}, {}} // 2 cards
	action := SearchAction{PlayerID: "P1"}
	
	result := Apply(state, action)
	
	if !result.Rooms["R12"].Searched {
		t.Error("Search should mark room as searched")
	}
	if len(result.Players["P1"].Hand) != len(state.Players["P1"].Hand)-1 {
		t.Error("Search should discard 1 card")
	}
}

func TestJSONRoundTrip(t *testing.T) {
	s := newTestGameState()
	b, _ := json.Marshal(s)
	var back GameState
	_ = json.Unmarshal(b, &back)
	if !reflect.DeepEqual(s, back) {
		t.Fatal("JSON round-trip failed")
	}
}

func TestRandSeedUnchanged(t *testing.T) {
	s := newTestGameState()
	seed := s.RandSeed
	Apply(s, SearchAction{PlayerID: "P1"})
	if s.RandSeed != seed {
		t.Fatal("RandSeed must not change inside reducer")
	}
}

// Test helpers
func newTestGameState() GameState {
	return GameState{
		Round:    1,
		Time:     15,
		RandSeed: 42,
		Rooms: map[RoomID]*RoomState{
			"R12": {ID: "R12", Type: Predefined}, // Start
			"R07": {ID: "R07", Type: AmmoCache},   // Adjacent to R12
			"R01": {ID: "R01", Type: Predefined},  // Non-adjacent to R12
		},
		Players: map[PlayerID]*PlayerState{
			"P1": {
				ID:       "P1",
				Location: "R12",
				Hand:     []Card{{}, {}},
			},
		},
		Enemies: map[EnemyID]*Enemy{},
	}
}

type invalidAction struct{}
func (invalidAction) isAction() {}