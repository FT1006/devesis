package core

import (
	"fmt"
)

// ApplySpawnEnemy creates new enemies in target locations  
func ApplySpawnEnemy(state *GameState, effect Effect, playerID PlayerID, log *EffectLog) error {
	targets := getRoomTargets(state, effect.Scope, playerID)
	for _, room := range targets {
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
	}
	return nil
}