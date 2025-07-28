package core

import "fmt"

// ApplyEffect is the centralized effect handler that consolidates all effect.Op switch statements
func ApplyEffect(state *GameState, effect Effect, playerID PlayerID, log *EffectLog) error {
	var err error
	switch effect.Op {
	case ModifyHP:
		err = ApplyModifyHP(state, effect, playerID, log)
	case ModifyAmmo:
		err = ApplyModifyAmmo(state, effect, playerID, log)
	case DrawCards:
		err = ApplyDrawCards(state, effect, playerID, log)
	case DiscardCards:
		err = ApplyDiscardCards(state, effect, playerID, log)
	case OutOfRam:
		err = ApplyOutOfRam(state, effect, playerID, log)
	case ModifyBugs:
		err = ApplyModifyBugs(state, effect, playerID, log)
	case RevealRoom:
		err = ApplyRevealRoom(state, effect, playerID, log)
	case CleanRoom:
		err = ApplyCleanRoom(state, effect, playerID, log)
	case SetCorrupted:
		err = ApplySetCorrupted(state, effect, playerID, log)
	case SpawnEnemy:
		err = ApplySpawnEnemy(state, effect, playerID, log)
	case MoveEnemies:
		err = ApplyMoveEnemies(state, effect, playerID, log)
	default:
		err = fmt.Errorf("unknown effect op: %v", effect.Op)
	}
	
	// Centralized error logging
	if err != nil {
		log.Add("⚠️ Effect failed - Op: %s, Scope: %s, N: %d - Error: %v", 
			getEffectOpName(effect.Op), getScopeName(effect.Scope), effect.N, err)
	}
	
	return err
}