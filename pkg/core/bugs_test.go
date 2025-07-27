package core

import (
	"testing"
)

func TestPlaceBugs_RespectsMaxLimit(t *testing.T) {
	// Create state with one room
	state := GameState{
		RandSeed: 42,
		Round:    1,
		Time:     0,
		Rooms: map[RoomID]*RoomState{
			"R01": {
				ID:         "R01",
				BugMarkers: 8, // Near max
				Corrupted:  true,
			},
		},
		SpawnBag: &SpawnBag{
			Tokens: []EnemyType{InfiniteLoop},
		},
		Enemies: make(map[EnemyID]*Enemy),
	}
	
	// Try to add 5 bugs (should cap at 9)
	PlaceBugs(&state, 5)
	
	// Should have exactly 9 bugs (not 13)
	if state.Rooms["R01"].BugMarkers != MaxBugMarkers {
		t.Errorf("Expected %d bugs, got %d", MaxBugMarkers, state.Rooms["R01"].BugMarkers)
	}
}

func TestPlaceBugs_AutoCorruption(t *testing.T) {
	// Create state with room at 2 bugs
	state := GameState{
		RandSeed: 42,
		Round:    1,
		Time:     0,
		Rooms: map[RoomID]*RoomState{
			"R01": {
				ID:         "R01",
				BugMarkers: 2,
				Corrupted:  false,
			},
		},
		SpawnBag: &SpawnBag{
			Tokens: []EnemyType{InfiniteLoop},
		},
		Enemies: make(map[EnemyID]*Enemy),
	}
	
	// Add 2 bugs (should trigger corruption at 3+)
	PlaceBugs(&state, 2)
	
	// Should be corrupted now
	if !state.Rooms["R01"].Corrupted {
		t.Error("Room should be corrupted at 3+ bugs")
	}
	
	// Should have 4 bugs
	if state.Rooms["R01"].BugMarkers != 4 {
		t.Errorf("Expected 4 bugs, got %d", state.Rooms["R01"].BugMarkers)
	}
}

func TestSpawnEnemiesForCorruptedRooms_DrawsFromBag(t *testing.T) {
	state := GameState{
		RandSeed: 42,
		Round:    1,
		Time:     0,
		SpawnBag: &SpawnBag{
			Tokens: []EnemyType{InfiniteLoop, StackOverflow},
		},
		Enemies: make(map[EnemyID]*Enemy),
	}
	
	// Spawn enemy for room R01
	SpawnEnemiesForCorruptedRooms(&state, []RoomID{"R01"})
	
	// Should have created one enemy
	if len(state.Enemies) != 1 {
		t.Errorf("Expected 1 enemy, got %d", len(state.Enemies))
	}
	
	// Bag should have one fewer token
	if len(state.SpawnBag.Tokens) != 1 {
		t.Errorf("Expected 1 token remaining, got %d", len(state.SpawnBag.Tokens))
	}
	
	// Enemy should be in correct location
	for _, enemy := range state.Enemies {
		if enemy.Location != "R01" {
			t.Errorf("Enemy in wrong location: %s", enemy.Location)
		}
	}
}

func TestSpawnEnemiesForCorruptedRooms_EmptyBag(t *testing.T) {
	state := GameState{
		SpawnBag: &SpawnBag{
			Tokens: []EnemyType{}, // Empty bag
		},
		Enemies: make(map[EnemyID]*Enemy),
	}
	
	// Try to spawn (should do nothing)
	SpawnEnemiesForCorruptedRooms(&state, []RoomID{"R01"})
	
	// Should have no enemies
	if len(state.Enemies) != 0 {
		t.Errorf("Expected 0 enemies, got %d", len(state.Enemies))
	}
}

func TestUpdateRoomCorruption_AutoCorruptAndUncorrupt(t *testing.T) {
	state := GameState{
		Rooms: map[RoomID]*RoomState{
			"R01": {BugMarkers: 2, Corrupted: true},  // Should uncorrupt
			"R02": {BugMarkers: 3, Corrupted: false}, // Should corrupt
			"R03": {BugMarkers: 5, Corrupted: true},  // Stay corrupted
		},
	}
	
	UpdateRoomCorruption(&state)
	
	if state.Rooms["R01"].Corrupted {
		t.Error("R01 should be uncorrupted (2 bugs)")
	}
	
	if !state.Rooms["R02"].Corrupted {
		t.Error("R02 should be corrupted (3 bugs)")
	}
	
	if !state.Rooms["R03"].Corrupted {
		t.Error("R03 should stay corrupted (5 bugs)")
	}
}

func TestInitializeSpawnBag_CorrectDistribution(t *testing.T) {
	bag := initializeSpawnBag()
	
	// Count enemy types
	loops := 0
	overflows := 0
	pythogoras := 0
	
	for _, enemy := range bag.Tokens {
		switch enemy {
		case InfiniteLoop:
			loops++
		case StackOverflow:
			overflows++
		case Pythogoras:
			pythogoras++
		}
	}
	
	// Check expected distribution
	if loops != 10 {
		t.Errorf("Expected 10 InfiniteLoops, got %d", loops)
	}
	if overflows != 6 {
		t.Errorf("Expected 6 StackOverflows, got %d", overflows)
	}
	if pythogoras != 2 {
		t.Errorf("Expected 2 Pythogoras, got %d", pythogoras)
	}
	
	// Total should be 18
	if len(bag.Tokens) != 18 {
		t.Errorf("Expected 18 total tokens, got %d", len(bag.Tokens))
	}
}