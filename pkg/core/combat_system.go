package core

// ApplyCombat handles combat actions and returns a modified GameState
func ApplyCombat(state GameState, action Action, log *EffectLog) GameState {
	// Deep copy the state to avoid mutations
	result := deepCopyGameState(state)
	
	switch a := action.(type) {
	case ShootAction:
		applyShootAction(&result, a, log)
	case MeleeAction:
		applyMeleeAction(&result, a, log)
	default:
		// Return unchanged state for invalid actions
		return result
	}
	
	// Remove dead enemies
	removeDeadEnemies(&result)
	
	return result
}

// ApplyCombatWithoutLog is a compatibility wrapper
func ApplyCombatWithoutLog(state GameState, action Action) GameState {
	log := NewEffectLog()
	return ApplyCombat(state, action, log)
}


func applyShootAction(state *GameState, action ShootAction, log *EffectLog) {
	player, exists := state.Players[action.PlayerID]
	if !exists {
		return
	}
	
	// Check if player has ammo
	if player.Ammo < ShootAmmoCost {
		return
	}
	
	// Consume ammo
	oldAmmo := player.Ammo
	player.Ammo -= ShootAmmoCost
	log.Add("🔫 %s shoots! Ammo: %d → %d", action.PlayerID, oldAmmo, player.Ammo)
	
	// Get adjacent rooms
	adjacentRooms := GetAdjacentRooms(player.Location)
	
	// Damage all enemies in adjacent rooms
	for _, enemy := range state.Enemies {
		for _, roomID := range adjacentRooms {
			if enemy.Location == roomID {
				oldHP := enemy.HP
				if enemy.HP > ShootDamage {
					enemy.HP -= ShootDamage
				} else if enemy.HP > 0 {
					enemy.HP = 0
				}
				if oldHP != enemy.HP {
					log.Add("💥 Hit %s in %s! HP: %d → %d", getEnemyDisplayName(enemy.Type), roomID, oldHP, enemy.HP)
				}
				break
			}
		}
	}
}

func applyMeleeAction(state *GameState, action MeleeAction, log *EffectLog) {
	player, exists := state.Players[action.PlayerID]
	if !exists {
		return
	}
	
	log.Add("⚔️ %s attacks with melee!", action.PlayerID)
	
	// Damage all enemies in same room (no ammo cost)
	for _, enemy := range state.Enemies {
		if enemy.Location == player.Location {
			oldHP := enemy.HP
			if enemy.HP > MeleeDamage {
				enemy.HP -= MeleeDamage
			} else if enemy.HP > 0 {
				enemy.HP = 0
			}
			if oldHP != enemy.HP {
				log.Add("💥 Hit %s! HP: %d → %d", getEnemyDisplayName(enemy.Type), oldHP, enemy.HP)
			}
		}
	}
}

func removeDeadEnemies(state *GameState) {
	for enemyID, enemy := range state.Enemies {
		if enemy.HP <= 0 {
			delete(state.Enemies, enemyID)
		}
	}
}

