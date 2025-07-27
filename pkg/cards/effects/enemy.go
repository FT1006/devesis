package effects

import (
	"fmt"
	"github.com/spaceship/devesis/pkg/core"
)

// ApplySpawnEnemy creates new enemies in target locations  
func ApplySpawnEnemy(state *core.GameState, effect Effect, playerID core.PlayerID) error {
	targets := getRoomTargets(state, effect.Scope, playerID)
	for _, room := range targets {
		// Convert N to enemy type
		var enemyType core.EnemyType
		switch effect.N {
		case 1:
			enemyType = core.InfiniteLoop
		case 2:
			enemyType = core.StackOverflow
		case 3:
			enemyType = core.Pythogoras
		default:
			return fmt.Errorf("invalid enemy type: %d", effect.N)
		}
		
		// Create enemy directly (bypass spawn bag for special effects)
		enemyID := core.EnemyID(fmt.Sprintf("E%d", len(state.Enemies)+1))
		stats := core.ENEMY_STATS[enemyType]
		enemy := &core.Enemy{
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