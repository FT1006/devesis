package core

// ApplyCombat handles combat actions and returns a modified GameState
func ApplyCombat(state GameState, action Action) GameState {
	// Deep copy the state to avoid mutations
	result := deepCopyGameState(state)
	
	switch a := action.(type) {
	case ShootAction:
		applyShootAction(&result, a)
	case MeleeAction:
		applyMeleeAction(&result, a)
	default:
		// Return unchanged state for invalid actions
		return result
	}
	
	// Remove dead enemies
	removeDeadEnemies(&result)
	
	return result
}

func applyShootAction(state *GameState, action ShootAction) {
	player, exists := state.Players[action.PlayerID]
	if !exists {
		return
	}
	
	// Check if player has ammo
	if player.Ammo < ShootAmmoCost {
		return
	}
	
	// Consume ammo
	player.Ammo -= ShootAmmoCost
	
	// Get adjacent rooms
	adjacentRooms := GetAdjacentRooms(player.Location)
	
	// Damage all enemies in adjacent rooms
	for _, enemy := range state.Enemies {
		for _, roomID := range adjacentRooms {
			if enemy.Location == roomID {
				if enemy.HP > 0 {
					enemy.HP--
				}
				break
			}
		}
	}
}

func applyMeleeAction(state *GameState, action MeleeAction) {
	player, exists := state.Players[action.PlayerID]
	if !exists {
		return
	}
	
	// Damage all enemies in same room (no ammo cost)
	for _, enemy := range state.Enemies {
		if enemy.Location == player.Location {
			if enemy.HP > 0 {
				enemy.HP--
			}
		}
	}
}

func removeDeadEnemies(state *GameState) {
	for enemyID, enemy := range state.Enemies {
		if enemy.HP == 0 {
			delete(state.Enemies, enemyID)
		}
	}
}

