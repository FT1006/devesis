package core

import (
	"testing"
)

func TestMedBayHealsPlayer(t *testing.T) {
	gs := GameState{
		ActivePlayer: "P1",
		Rooms: map[RoomID]*RoomState{
			"R01": {ID: "R01", Type: MedBay, Searched: true},
		},
		Players: map[PlayerID]*PlayerState{
			"P1": {
				ID:       "P1",
				Location: "R01",
				HP:       3,
				MaxHP:    5,
			},
		},
		ActionsLeft: 1,
	}
	
	action := RoomAction{PlayerID: "P1"}
	log := NewEffectLog()
	
	result := Apply(gs, action, log)
	
	// Should heal 2 HP
	if result.Players["P1"].HP != 5 {
		t.Errorf("expected HP 5, got %d", result.Players["P1"].HP)
	}
	
	// Room actions don't consume normal actions (they use SpecialUsed flag instead)
	if result.ActionsLeft != 1 {
		t.Errorf("expected 1 action left, got %d", result.ActionsLeft)
	}
	
	// Should mark special used
	if !result.Players["P1"].SpecialUsed {
		t.Error("expected SpecialUsed to be true")
	}
}

func TestAmmoCacheRefillsAmmo(t *testing.T) {
	gs := GameState{
		ActivePlayer: "P1",
		Rooms: map[RoomID]*RoomState{
			"R01": {ID: "R01", Type: AmmoCache, Searched: true},
		},
		Players: map[PlayerID]*PlayerState{
			"P1": {
				ID:       "P1",
				Location: "R01",
				Ammo:     1,
				MaxAmmo:  5,
			},
		},
		ActionsLeft: 1,
	}
	
	action := RoomAction{PlayerID: "P1"}
	log := NewEffectLog()
	
	result := Apply(gs, action, log)
	
	// Should add 3 ammo
	if result.Players["P1"].Ammo != 4 {
		t.Errorf("expected ammo 4, got %d", result.Players["P1"].Ammo)
	}
}

func TestCleanRoomRemovesBugs(t *testing.T) {
	gs := GameState{
		ActivePlayer: "P1",
		Rooms: map[RoomID]*RoomState{
			"R01": {ID: "R01", Type: CleanRoomType, Searched: true, BugMarkers: 3},
		},
		Players: map[PlayerID]*PlayerState{
			"P1": {
				ID:       "P1",
				Location: "R01",
			},
		},
		ActionsLeft: 1,
	}
	
	action := RoomAction{PlayerID: "P1"}
	log := NewEffectLog()
	
	result := Apply(gs, action, log)
	
	// Should remove all bugs
	if result.Rooms["R01"].BugMarkers != 0 {
		t.Errorf("expected 0 bugs, got %d", result.Rooms["R01"].BugMarkers)
	}
}

func TestRoomActionBlockedWhenCorrupted(t *testing.T) {
	gs := GameState{
		ActivePlayer: "P1",
		Rooms: map[RoomID]*RoomState{
			"R01": {
				ID:         "R01",
				Type:       MedBay,
				Searched:   true,
				Corrupted:  true, // Corrupted room
				BugMarkers: 3,
			},
		},
		Players: map[PlayerID]*PlayerState{
			"P1": {
				ID:       "P1",
				Location: "R01",
				HP:       3,
				MaxHP:    5,
			},
		},
		ActionsLeft: 1,
	}
	
	action := RoomAction{PlayerID: "P1"}
	log := NewEffectLog()
	
	result := Apply(gs, action, log)
	
	// Should not heal (room corrupted)
	if result.Players["P1"].HP != 3 {
		t.Errorf("expected HP unchanged at 3, got %d", result.Players["P1"].HP)
	}
	
	// Should not consume action (invalid)
	if result.ActionsLeft != 1 {
		t.Errorf("expected 1 action left, got %d", result.ActionsLeft)
	}
}

func TestRoomActionOncePerTurn(t *testing.T) {
	gs := GameState{
		ActivePlayer: "P1",
		Rooms: map[RoomID]*RoomState{
			"R01": {
				ID:       "R01",
				Type:     MedBay,
				Searched: true,
			},
		},
		Players: map[PlayerID]*PlayerState{
			"P1": {
				ID:          "P1",
				Location:    "R01",
				HP:          3,
				MaxHP:       5,
				SpecialUsed: true, // Already used room action this turn
			},
		},
		ActionsLeft: 1,
	}
	
	action := RoomAction{PlayerID: "P1"}
	log := NewEffectLog()
	
	result := Apply(gs, action, log)
	
	// Should not heal (already used)
	if result.Players["P1"].HP != 3 {
		t.Errorf("expected HP unchanged at 3, got %d", result.Players["P1"].HP)
	}
	
	// Should not consume action
	if result.ActionsLeft != 1 {
		t.Errorf("expected 1 action left, got %d", result.ActionsLeft)
	}
}

func TestRoomActionRequiresSearched(t *testing.T) {
	gs := GameState{
		ActivePlayer: "P1",
		Rooms: map[RoomID]*RoomState{
			"R01": {
				ID:       "R01",
				Type:     MedBay,
				Searched: false, // Not searched yet
			},
		},
		Players: map[PlayerID]*PlayerState{
			"P1": {
				ID:       "P1",
				Location: "R01",
				HP:       3,
				MaxHP:    5,
			},
		},
		ActionsLeft: 1,
	}
	
	action := RoomAction{PlayerID: "P1"}
	log := NewEffectLog()
	
	result := Apply(gs, action, log)
	
	// Should not heal (room not searched)
	if result.Players["P1"].HP != 3 {
		t.Errorf("expected HP unchanged at 3, got %d", result.Players["P1"].HP)
	}
}