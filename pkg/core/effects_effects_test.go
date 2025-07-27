package core

import (
	"testing"
)

func TestApplyModifyHP_Self(t *testing.T) {
	state := &GameState{
		Players: map[PlayerID]*PlayerState{
			"P1": {
				ID:    "P1",
				HP:    5,
				MaxHP: 8,
			},
		},
	}

	effect := Effect{
		Op:    ModifyHP,
		Scope: Self,
		N:     2,
	}

	err := ApplyModifyHP(state, effect, "P1")
	if err != nil {
		t.Errorf("ApplyModifyHP failed: %v", err)
	}

	if state.Players["P1"].HP != 7 {
		t.Errorf("Expected HP 7, got %d", state.Players["P1"].HP)
	}
}

func TestApplyModifyBugs_CurrentRoom(t *testing.T) {
	state := &GameState{
		Players: map[PlayerID]*PlayerState{
			"P1": {
				ID:       "P1",
				Location: "R01",
			},
		},
		Rooms: map[RoomID]*RoomState{
			"R01": {
				ID:         "R01",
				BugMarkers: 2,
				Corrupted:  false,
			},
		},
	}

	effect := Effect{
		Op:    ModifyBugs,
		Scope: CurrentRoom,
		N:     2, // Will hit corruption threshold
	}

	err := ApplyModifyBugs(state, effect, "P1")
	if err != nil {
		t.Errorf("ApplyModifyBugs failed: %v", err)
	}

	room := state.Rooms["R01"]
	if room.BugMarkers != 4 {
		t.Errorf("Expected 4 bugs, got %d", room.BugMarkers)
	}

	if !room.Corrupted {
		t.Error("Room should be corrupted at 3+ bugs")
	}
}

func TestValidateCard_Success(t *testing.T) {
	card := Card{
		ID:     "QUICK_FIX",
		Name:   "Quick Fix",
		Source: SrcAction,
		Effects: []Effect{
			{Op: ModifyBugs, Scope: CurrentRoom, N: -1},
		},
	}

	err := ValidateCard(card)
	if err != nil {
		t.Errorf("Valid card failed validation: %v", err)
	}
}

func TestValidateCard_InvalidOpScope(t *testing.T) {
	card := Card{
		ID:     "INVALID",
		Name:   "Invalid Card",
		Source: SrcAction,
		Effects: []Effect{
			{Op: ModifyHP, Scope: CurrentRoom, N: 1}, // Invalid: ModifyHP can't target rooms
		},
	}

	err := ValidateCard(card)
	if err == nil {
		t.Error("Invalid card should have failed validation")
	}
}