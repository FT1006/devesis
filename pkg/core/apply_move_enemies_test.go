package core

import (
	"testing"
)

func TestApplyMoveEnemiesMovesShortestPath(t *testing.T) {
	// Use actual game rooms that have adjacency: R12 - R07 - R06
	gs := GameState{
		ActivePlayer: "P1",
		Rooms: map[RoomID]*RoomState{
			"R12": {ID: "R12"},
			"R07": {ID: "R07"},
			"R06": {ID: "R06"},
		},
		Players: map[PlayerID]*PlayerState{
			"P1": {
				ID:       "P1",
				Location: "R06",
				HP:       5,
			},
		},
		Enemies: map[EnemyID]*Enemy{
			"E1": {
				ID:       "E1",
				Type:     InfiniteLoop,
				Location: "R12",
				HP:       1,
			},
		},
	}
	
	// Enemy should move 1 step toward player
	eff := Effect{Op: MoveEnemies, Scope: AllRooms, N: 1}
	log := NewEffectLog()
	
	if err := ApplyMoveEnemies(&gs, eff, "P1", log); err != nil {
		t.Fatalf("apply error: %v", err)
	}
	
	if loc := gs.Enemies["E1"].Location; loc != "R07" {
		t.Fatalf("expected E1 at R07, got %s", loc)
	}
}

func TestApplyMoveEnemiesNoMovementWhenPlayerDead(t *testing.T) {
	gs := GameState{
		ActivePlayer: "P1",
		Rooms: map[RoomID]*RoomState{
			"R01": {ID: "R01"},
			"R02": {ID: "R02"},
		},
		Players: map[PlayerID]*PlayerState{
			"P1": {
				ID:       "P1",
				Location: "R02",
				HP:       0, // Dead player
			},
		},
		Enemies: map[EnemyID]*Enemy{
			"E1": {
				ID:       "E1",
				Type:     InfiniteLoop,
				Location: "R12",
				HP:       1,
			},
		},
	}
	
	eff := Effect{Op: MoveEnemies, Scope: AllRooms, N: 1}
	log := NewEffectLog()
	
	if err := ApplyMoveEnemies(&gs, eff, "P1", log); err != nil {
		t.Fatalf("apply error: %v", err)
	}
	
	// Enemy should not move when all players are dead
	if loc := gs.Enemies["E1"].Location; loc != "R01" {
		t.Fatalf("expected E1 to stay at R01, got %s", loc)
	}
}

func TestApplyMoveEnemiesPythogorasMovesTowardEscape(t *testing.T) {
	// Pythogoras should move toward escape rooms (R19/R20) instead of players
	gs := GameState{
		ActivePlayer: "P1",
		Rooms: map[RoomID]*RoomState{
			"R12": {ID: "R12"},
			"R13": {ID: "R13"}, 
			"R19": {ID: "R19"}, // Escape room
			"R01": {ID: "R01"},
		},
		Players: map[PlayerID]*PlayerState{
			"P1": {
				ID:       "P1",
				Location: "R01", // Player far from escape
				HP:       5,
			},
		},
		Enemies: map[EnemyID]*Enemy{
			"PY1": {
				ID:       "PY1",
				Type:     Pythogoras,
				Location: "R12",
				HP:       6,
			},
		},
	}
	
	eff := Effect{Op: MoveEnemies, Scope: AllRooms, N: 1}
	log := NewEffectLog()
	
	if err := ApplyMoveEnemies(&gs, eff, "P1", log); err != nil {
		t.Fatalf("apply error: %v", err)
	}
	
	// Pythogoras should move toward escape room R19, not toward player at R01
	if loc := gs.Enemies["PY1"].Location; loc != "R13" {
		t.Fatalf("expected Pythogoras to move toward escape (R13), got %s", loc)
	}
}

func TestApplyMoveEnemiesLimitedByN(t *testing.T) {
	// Enemy can only move N steps even if player is farther
	gs := GameState{
		ActivePlayer: "P1",
		Rooms: map[RoomID]*RoomState{
			"R01": {ID: "R01"},
			"R02": {ID: "R02"},
			"R03": {ID: "R03"},
			"R04": {ID: "R04"},
			"R05": {ID: "R05"},
		},
		Players: map[PlayerID]*PlayerState{
			"P1": {
				ID:       "P1",
				Location: "R05",
				HP:       5,
			},
		},
		Enemies: map[EnemyID]*Enemy{
			"E1": {
				ID:       "E1",
				Type:     StackOverflow,
				Location: "R01",
				HP:       3,
			},
		},
	}
	
	// Enemy can move 2 steps maximum
	eff := Effect{Op: MoveEnemies, Scope: AllRooms, N: 2}
	log := NewEffectLog()
	
	if err := ApplyMoveEnemies(&gs, eff, "P1", log); err != nil {
		t.Fatalf("apply error: %v", err)
	}
	
	// Should move 2 steps from R01 to R03 (not all the way to R05)
	if loc := gs.Enemies["E1"].Location; loc != "R03" {
		t.Fatalf("expected E1 at R03 (2 steps), got %s", loc)
	}
}