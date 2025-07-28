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
	log.Add("🚶 MoveEnemies: Processing %d enemies", len(state.Enemies))
	moveCount := 0
	for _, enemy := range state.Enemies {
		log.Add("🚶 Processing %s at %s", getEnemyDisplayName(enemy.Type), enemy.Location)
		// Find closest living player
		var closestPlayer RoomID
		minDistance := effect.N + 1
		
		for _, player := range state.Players {
			if player.HP > 0 {
				path := CanTraverse(state, PathQuery{From: enemy.Location, To: player.Location, MaxSteps: effect.N})
				if path.Valid && len(path.Path)-1 < minDistance {
					closestPlayer = player.Location
					minDistance = len(path.Path) - 1
				}
			}
		}
		
		if closestPlayer == "" {
			continue
		}
		
		// Move toward closest player
		path := CanTraverse(state, PathQuery{From: enemy.Location, To: closestPlayer, MaxSteps: effect.N})
		if path.Valid && len(path.Path) > 1 {
			oldLocation := enemy.Location
			maxIndex := effect.N
			if len(path.Path)-1 < maxIndex {
				maxIndex = len(path.Path) - 1
			}
			enemy.Location = path.Path[maxIndex]
			
			if oldLocation != enemy.Location {
				log.Add("🚶 %s moves %s → %s", getEnemyDisplayName(enemy.Type), oldLocation, enemy.Location)
				moveCount++
			}
		}
	}
	
	log.Add("🚶 %d enemies moved", moveCount)
	return nil
}

