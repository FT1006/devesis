package core

import (
	"fmt"
)

// getSpawnRoomTargets resolves which rooms are affected by spawn effects
// For RoomWithMostBugs, it uses the spawn-specific logic that doesn't require bugs > 0
func getSpawnRoomTargets(state *GameState, scope ScopeType, playerID PlayerID) []*RoomState {
	if scope == RoomWithMostBugs {
		targetRoomID := GetRoomWithMostBugsForSpawn(state)
		if room := state.Rooms[targetRoomID]; room != nil {
			return []*RoomState{room}
		}
		return nil
	}
	
	// For all other scopes, use the regular targeting logic
	return getRoomTargets(state, scope, playerID)
}

// ApplySpawnEnemy creates new enemies in target locations  
func ApplySpawnEnemy(state *GameState, effect Effect, playerID PlayerID, log *EffectLog) error {
	targets := getSpawnRoomTargets(state, effect.Scope, playerID)
	if len(targets) == 0 {
		log.Add("⚠️ SpawnEnemy: No target rooms found for scope %s", getScopeName(effect.Scope))
		return nil
	}
	
	log.Add("🎯 SpawnEnemy targeting %d room(s)", len(targets))
	for _, room := range targets {
		log.Add("🎯 Target room: %s", room.ID)
		
		// Convert N to enemy type
		var enemyType EnemyType
		switch effect.N {
		case 1:
			enemyType = InfiniteLoop
		case 2:
			enemyType = StackOverflow
		case 3:
			enemyType = Pythogoras
		default:
			return fmt.Errorf("invalid enemy type: %d", effect.N)
		}
		
		// Create enemy directly (bypass spawn bag for special effects)
		enemyID := EnemyID(fmt.Sprintf("E%d", len(state.Enemies)+1))
		stats := ENEMY_STATS[enemyType]
		enemy := &Enemy{
			ID:       enemyID,
			Type:     enemyType,
			HP:       stats.HP,
			MaxHP:    stats.HP,
			Damage:   stats.Damage,
			Location: room.ID,
		}
		state.Enemies[enemyID] = enemy
		
		log.Add("👹 %s spawned in %s", getEnemyDisplayName(enemyType), room.ID)
	}
	return nil
}

// ApplyMoveEnemies moves all enemies N steps toward players using the movement system
func ApplyMoveEnemies(state *GameState, effect Effect, playerID PlayerID, log *EffectLog) error {
	maxStep := effect.N
	moved := 0

	for _, enemy := range state.Enemies {
		// 1️⃣ Find shortest path to ANY living player (no step limit)
		bestPath := PathResult{} // PathResult{Valid bool, Path []RoomID}
		bestLen := 999           // Use large number instead of math.MaxInt for simplicity

		for _, player := range state.Players {
			if player.HP == 0 {
				continue // Skip dead players
			}
			path := CanTraverse(state, PathQuery{
				From:     enemy.Location,
				To:       player.Location,
				MaxSteps: 99, // Effectively no cap - find any reachable player
			})
			if !path.Valid {
				continue // Player not reachable
			}
			if pathLen := len(path.Path) - 1; pathLen < bestLen {
				bestLen, bestPath = pathLen, path
			}
		}

		if !bestPath.Valid {
			continue // No reachable player found
		}

		// 2️⃣ Move up to N steps along that path
		step := maxStep
		if bestLen < step {
			step = bestLen // Can't move more steps than path length
		}
		if step == 0 {
			continue // Already at player location
		}

		oldLocation := enemy.Location
		enemy.Location = bestPath.Path[step]
		log.Add("🚶 %s moves %s → %s", getEnemyDisplayName(enemy.Type), oldLocation, enemy.Location)
		moved++
	}
	
	log.Add("🚶 %d enemies moved", moved)
	return nil
}

