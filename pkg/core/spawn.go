package core

import (
	"fmt"
	"math/rand"
)

// GetGameRNG creates a consistent RNG for game operations
func GetGameRNG(state *GameState) *rand.Rand {
	// Create deterministic RNG from game state
	// Use Round * 1000 + Time to ensure unique seeds for different game moments
	seed := state.RandSeed + int64(state.Round*1000) + int64(state.Time)
	return rand.New(rand.NewSource(seed))
}

// SpawnEnemyFromBag draws an enemy from the spawn bag and places it in the specified room
// This is the ONLY way to spawn enemies in the game
func SpawnEnemyFromBag(state *GameState, roomID RoomID) bool {
	if state.SpawnBag == nil || len(state.SpawnBag.Tokens) == 0 {
		return false // No enemies left to spawn
	}
	
	// Use the game's RNG
	rng := GetGameRNG(state)
	
	// Draw random enemy from spawn bag
	tokenIndex := rng.Intn(len(state.SpawnBag.Tokens))
	enemyType := state.SpawnBag.Tokens[tokenIndex]
	
	// Remove token from bag (draw without replacement)
	state.SpawnBag.Tokens = append(
		state.SpawnBag.Tokens[:tokenIndex],
		state.SpawnBag.Tokens[tokenIndex+1:]...,
	)
	
	// Create the enemy
	// Validate enemy type
	if enemyType < InfiniteLoop || enemyType > Pythogoras {
		return false
	}
	
	// Generate unique enemy ID
	enemyID := EnemyID(fmt.Sprintf("E%d", len(state.Enemies)+1))
	
	// Create enemy with stats from constants
	stats := ENEMY_STATS[enemyType]
	enemy := &Enemy{
		ID:       enemyID,
		Type:     enemyType,
		HP:       stats.HP,
		MaxHP:    stats.HP,
		Damage:   stats.Damage,
		Location: roomID,
	}
	
	state.Enemies[enemyID] = enemy
	return true
}

// GetSpawnBagStatus returns info about the current spawn bag
func GetSpawnBagStatus(state *GameState) (int, map[EnemyType]int) {
	if state.SpawnBag == nil {
		return 0, nil
	}
	
	total := len(state.SpawnBag.Tokens)
	counts := make(map[EnemyType]int)
	
	for _, enemyType := range state.SpawnBag.Tokens {
		counts[enemyType]++
	}
	
	return total, counts
}

// IsSpawnBagEmpty returns true if no more enemies can be spawned
func IsSpawnBagEmpty(state *GameState) bool {
	return state.SpawnBag == nil || len(state.SpawnBag.Tokens) == 0
}

