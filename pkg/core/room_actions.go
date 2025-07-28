package core

// ApplyRoomAction handles room-specific special abilities
func ApplyRoomAction(state GameState, action RoomAction, log *EffectLog) GameState {
	// Deep copy the state to avoid mutations
	newState := deepCopyGameState(state)
	
	player, exists := newState.Players[action.PlayerID]
	if !exists {
		return newState
	}
	
	// Get current room
	room, roomExists := newState.Rooms[player.Location]
	if !roomExists {
		return newState
	}
	
	// Check if room is corrupted
	if room.Corrupted {
		log.Add("✗ Cannot use room actions in corrupted rooms!")
		return state // Return original state unchanged
	}
	
	// Check if player already used room action this turn
	if player.SpecialUsed {
		log.Add("✗ Room action already used this turn!")
		return state // Return original state unchanged
	}
	
	// Apply room-specific effect
	switch room.Type {
	case MedBay:
		applyMedBayAction(&newState, player, log)
	case AmmoCache:
		applyAmmoCacheAction(&newState, player, log)
	case CleanRoomType:
		applyCleanRoomAction(&newState, player, log)
	default:
		log.Add("✗ No special room action available here")
		return state // Return original state unchanged
	}
	
	return newState
}

func applyMedBayAction(state *GameState, player *PlayerState, log *EffectLog) {
	// Check if healing would be beneficial
	if player.HP >= player.MaxHP {
		log.Add("💡 Already at full health - MedBay does nothing")
		return
	}
	
	// Apply healing
	oldHP := player.HP
	newHP := player.HP + MedBayHealAmount
	if newHP > player.MaxHP {
		newHP = player.MaxHP
	}
	player.HP = newHP
	
	// Mark as used
	player.SpecialUsed = true
	
	log.Add("🏥 MedBay healing! HP: %d → %d (+%d)", oldHP, newHP, newHP-oldHP)
	log.Add("✓ Room action complete (1 action used)")
}

func applyAmmoCacheAction(state *GameState, player *PlayerState, log *EffectLog) {
	// Check if ammo refill would be beneficial
	if player.Ammo >= player.MaxAmmo {
		log.Add("💡 Already at full ammo - AmmoCache does nothing")
		return
	}
	
	// Apply ammo refill
	oldAmmo := player.Ammo
	newAmmo := player.Ammo + AmmoCacheAmount
	if newAmmo > player.MaxAmmo {
		newAmmo = player.MaxAmmo
	}
	player.Ammo = newAmmo
	
	// Mark as used
	player.SpecialUsed = true
	
	log.Add("🔫 AmmoCache refill! Ammo: %d → %d (+%d)", oldAmmo, newAmmo, newAmmo-oldAmmo)
	log.Add("✓ Room action complete (1 action used)")
}

func applyCleanRoomAction(state *GameState, player *PlayerState, log *EffectLog) {
	// Get adjacent rooms
	adjacentRoomIDs := GetAdjacentRooms(player.Location)
	
	// Count rooms that have bugs to clean
	bugsCleaned := 0
	roomsCleaned := 0
	
	log.Add("🧹 CleanRoom decontamination activated!")
	
	for _, roomID := range adjacentRoomIDs {
		adjacentRoom := state.Rooms[roomID]
		if adjacentRoom != nil && adjacentRoom.BugMarkers > 0 {
			oldBugs := adjacentRoom.BugMarkers
			adjacentRoom.BugMarkers--
			bugsCleaned++
			roomsCleaned++
			
			// Update corruption status
			adjacentRoom.Corrupted = adjacentRoom.BugMarkers >= BugCorruptionThreshold
			
			log.Add("   %s: %d → %d bugs (-1)", roomID, oldBugs, adjacentRoom.BugMarkers)
		} else if adjacentRoom != nil {
			log.Add("   %s: 0 bugs (no effect)", roomID)
		}
	}
	
	if bugsCleaned == 0 {
		log.Add("💡 No bugs in adjacent rooms - CleanRoom does nothing")
		return
	}
	
	// Mark as used
	player.SpecialUsed = true
	
	log.Add("✓ Cleaned %d bugs from %d rooms (1 action used)", bugsCleaned, roomsCleaned)
}