package core

import (
	"testing"
)

func TestCanTraverse_Length1(t *testing.T) {
	state := newMovementTestGameState()
	
	// R12 to R07 is adjacent (1 step)
	result := CanTraverse(&state, PathQuery{
		From: "R12", To: "R07", MaxSteps: 1,
	})
	
	if !result.Valid {
		t.Error("Expected valid 1-step path R12→R07")
	}
	if len(result.Path) != 2 || result.Path[0] != "R12" || result.Path[1] != "R07" {
		t.Errorf("Expected path [R12, R07], got %v", result.Path)
	}
}

func TestCanTraverse_Length2(t *testing.T) {
	state := newMovementTestGameState()
	
	// R12 → R07 → R06 (2 steps)
	result := CanTraverse(&state, PathQuery{
		From: "R12", To: "R06", MaxSteps: 2,
	})
	
	if !result.Valid {
		t.Error("Expected valid 2-step path R12→R06")
	}
	if len(result.Path) != 3 {
		t.Errorf("Expected 3-room path, got %d: %v", len(result.Path), result.Path)
	}
}

func TestCanTraverse_ExceedsMaxSteps(t *testing.T) {
	state := newMovementTestGameState()
	
	// R12 to R06 requires 2 steps, but only allow 1
	result := CanTraverse(&state, PathQuery{
		From: "R12", To: "R06", MaxSteps: 1,
	})
	
	if result.Valid {
		t.Error("Expected invalid path when exceeding MaxSteps")
	}
}

func TestCanTraverse_SameRoom(t *testing.T) {
	state := newMovementTestGameState()
	
	result := CanTraverse(&state, PathQuery{
		From: "R12", To: "R12", MaxSteps: 1,
	})
	
	if !result.Valid {
		t.Error("Expected valid path for same room")
	}
	if len(result.Path) != 1 || result.Path[0] != "R12" {
		t.Errorf("Expected path [R12], got %v", result.Path)
	}
}

func TestCanMove_BackwardCompatibility(t *testing.T) {
	state := newMovementTestGameState()
	
	// Test that old CanMove function still works
	if !CanMove(&state, "R12", "R07") {
		t.Error("CanMove should work for adjacent rooms")
	}
	if CanMove(&state, "R12", "R01") {
		t.Error("CanMove should fail for non-adjacent rooms")
	}
}

// Helper for movement tests - includes rooms needed for path testing
func newMovementTestGameState() GameState {
	return GameState{
		Rooms: map[RoomID]*RoomState{
			"R12": {ID: "R12", Type: Predefined}, // Start {3,3}
			"R07": {ID: "R07", Type: AmmoCache},   // Adjacent to R12 {3,2}
			"R06": {ID: "R06", Type: AmmoCache},   // Adjacent to R07 {2,2}
			"R01": {ID: "R01", Type: Predefined},  // Non-adjacent to R12 {3,0}
		},
	}
}