package core

import (
	"testing"
)

func TestSpawnEnemyFromBag_Success(t *testing.T) {
	state := GameState{
		RandSeed: 42,
		Round:    1,
		Time:     0,
		SpawnBag: &SpawnBag{
			Tokens: []EnemyType{InfiniteLoop, StackOverflow},
		},
		Enemies: make(map[EnemyID]*Enemy),
	}
	
	// Spawn enemy
	success := SpawnEnemyFromBag(&state, "R01")
	
	if !success {
		t.Error("Expected spawn to succeed")
	}
	
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

func TestSpawnEnemyFromBag_EmptyBag(t *testing.T) {
	state := GameState{
		SpawnBag: &SpawnBag{
			Tokens: []EnemyType{}, // Empty bag
		},
		Enemies: make(map[EnemyID]*Enemy),
	}
	
	// Try to spawn (should fail)
	success := SpawnEnemyFromBag(&state, "R01")
	
	if success {
		t.Error("Expected spawn to fail with empty bag")
	}
	
	// Should have no enemies
	if len(state.Enemies) != 0 {
		t.Errorf("Expected 0 enemies, got %d", len(state.Enemies))
	}
}

func TestSpawnEnemyFromBag_MultipleRooms(t *testing.T) {
	state := GameState{
		RandSeed: 42,
		Round:    1,
		Time:     0,
		SpawnBag: &SpawnBag{
			Tokens: []EnemyType{InfiniteLoop, StackOverflow, Pythogoras},
		},
		Enemies: make(map[EnemyID]*Enemy),
	}
	
	roomIDs := []RoomID{"R01", "R02", "R03", "R04"}
	
	// Spawn enemies using only SpawnEnemyFromBag
	spawned := 0
	for _, roomID := range roomIDs {
		if SpawnEnemyFromBag(&state, roomID) {
			spawned++
		}
		if IsSpawnBagEmpty(&state) {
			break
		}
	}
	
	// Should spawn 3 enemies (bag size)
	if spawned != 3 {
		t.Errorf("Expected 3 spawned, got %d", spawned)
	}
	
	// Should have 3 enemies total
	if len(state.Enemies) != 3 {
		t.Errorf("Expected 3 enemies, got %d", len(state.Enemies))
	}
	
	// Bag should be empty
	if len(state.SpawnBag.Tokens) != 0 {
		t.Errorf("Expected empty bag, got %d tokens", len(state.SpawnBag.Tokens))
	}
}

func TestGetSpawnBagStatus(t *testing.T) {
	state := GameState{
		SpawnBag: &SpawnBag{
			Tokens: []EnemyType{
				InfiniteLoop, InfiniteLoop,
				StackOverflow,
				Pythogoras,
			},
		},
	}
	
	total, counts := GetSpawnBagStatus(&state)
	
	if total != 4 {
		t.Errorf("Expected total 4, got %d", total)
	}
	
	if counts[InfiniteLoop] != 2 {
		t.Errorf("Expected 2 InfiniteLoops, got %d", counts[InfiniteLoop])
	}
	
	if counts[StackOverflow] != 1 {
		t.Errorf("Expected 1 StackOverflow, got %d", counts[StackOverflow])
	}
	
	if counts[Pythogoras] != 1 {
		t.Errorf("Expected 1 Pythogoras, got %d", counts[Pythogoras])
	}
}

func TestIsSpawnBagEmpty(t *testing.T) {
	// Test with tokens
	state1 := GameState{
		SpawnBag: &SpawnBag{
			Tokens: []EnemyType{InfiniteLoop},
		},
	}
	
	if IsSpawnBagEmpty(&state1) {
		t.Error("Bag should not be empty")
	}
	
	// Test empty bag
	state2 := GameState{
		SpawnBag: &SpawnBag{
			Tokens: []EnemyType{},
		},
	}
	
	if !IsSpawnBagEmpty(&state2) {
		t.Error("Bag should be empty")
	}
	
	// Test nil bag
	state3 := GameState{
		SpawnBag: nil,
	}
	
	if !IsSpawnBagEmpty(&state3) {
		t.Error("Nil bag should be empty")
	}
}

func TestGetGameRNG_Deterministic(t *testing.T) {
	state := GameState{
		RandSeed: 42,
		Round:    1,
		Time:     5,
	}
	
	// Should get same sequence with same state
	rng1 := GetGameRNG(&state)
	val1 := rng1.Intn(100)
	
	rng2 := GetGameRNG(&state)
	val2 := rng2.Intn(100)
	
	if val1 != val2 {
		t.Errorf("RNG should be deterministic: %d != %d", val1, val2)
	}
	
	// Should get different sequence with different state
	state.Round = 2
	rng3 := GetGameRNG(&state)
	val3 := rng3.Intn(100)
	
	if val1 == val3 {
		t.Error("Different game states should produce different RNG sequences")
	}
}