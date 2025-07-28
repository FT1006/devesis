package core

import (
	"testing"
)

// Layer 1 unit test for ApplyMoveEnemies following established patterns
func TestApplyMoveEnemies(t *testing.T) {
	tests := []struct {
		name     string
		effect   Effect
		playerID PlayerID
		setup    func(GameState) GameState
		validate func(*testing.T, GameState)
	}{
		{
			name:     "moves single enemy one step",
			effect:   Effect{Op: MoveEnemies, Scope: AllRooms, N: 1},
			playerID: "P1",
			setup: func(s GameState) GameState {
				s.Enemies = map[EnemyID]*Enemy{
					"E1": {ID: "E1", Type: InfiniteLoop, Location: "R01"},
				}
				return s
			},
			validate: func(t *testing.T, s GameState) {
				// Enemy should have moved to an adjacent room (R02 due to deterministic RNG)
				if s.Enemies["E1"].Location != "R02" {
					t.Errorf("Expected enemy at R02, got %s", s.Enemies["E1"].Location)
				}
			},
		},
		{
			name:     "avoids corrupted rooms",
			effect:   Effect{Op: MoveEnemies, Scope: AllRooms, N: 1},
			playerID: "P1",
			setup: func(s GameState) GameState {
				s.Rooms["R02"].Corrupted = true // Block adjacent room
				s.Enemies = map[EnemyID]*Enemy{
					"E1": {ID: "E1", Type: InfiniteLoop, Location: "R01"},
				}
				return s
			},
			validate: func(t *testing.T, s GameState) {
				// Enemy should not move to corrupted R02
				if s.Enemies["E1"].Location == "R02" {
					t.Error("Enemy moved to corrupted room")
				}
			},
		},
		{
			name:     "moves multiple steps",
			effect:   Effect{Op: MoveEnemies, Scope: AllRooms, N: 3},
			playerID: "P1",
			setup: func(s GameState) GameState {
				s.Enemies = map[EnemyID]*Enemy{
					"E1": {ID: "E1", Type: InfiniteLoop, Location: "R01"},
				}
				return s
			},
			validate: func(t *testing.T, s GameState) {
				// After 3 steps, enemy should have moved from R01
				if s.Enemies["E1"].Location == "R01" {
					t.Error("Enemy should have moved after 3 steps")
				}
			},
		},
		{
			name:     "handles empty enemy list",
			effect:   Effect{Op: MoveEnemies, Scope: AllRooms, N: 1},
			playerID: "P1",
			setup: func(s GameState) GameState {
				s.Enemies = map[EnemyID]*Enemy{} // No enemies
				return s
			},
			validate: func(t *testing.T, s GameState) {
				// Should not crash with empty enemy list
				if len(s.Enemies) != 0 {
					t.Error("Expected no enemies")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Start with minimal state
			state := minimalEnemyState()
			state = tt.setup(state)

			// Apply effect
			log := NewEffectLog()
			err := ApplyMoveEnemies(&state, tt.effect, tt.playerID, log)
			if err != nil {
				t.Fatalf("ApplyMoveEnemies failed: %v", err)
			}

			// Validate result
			tt.validate(t, state)
		})
	}
}

// Helper: create minimal state for enemy tests
func minimalEnemyState() GameState {
	return GameState{
		Round:    1,
		Time:     10,
		RandSeed: 42,
		Rooms: map[RoomID]*RoomState{
			"R01": {ID: "R01", Type: Empty, BugMarkers: 0, Corrupted: false},
			"R02": {ID: "R02", Type: Empty, BugMarkers: 1, Corrupted: false},
			"R03": {ID: "R03", Type: Empty, BugMarkers: 2, Corrupted: false},
		},
		Players: map[PlayerID]*PlayerState{
			"P1": {
				ID:       "P1",
				Class:    Frontend,
				Location: "R01",
			},
		},
		Enemies: map[EnemyID]*Enemy{},
	}
}