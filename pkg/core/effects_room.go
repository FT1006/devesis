package core

import "fmt"

// getCorruptionRoomTargets resolves which rooms are affected by corruption effects
// For RoomWithMostBugs, it uses spawn-like logic that doesn't require bugs > 0
func getCorruptionRoomTargets(state *GameState, scope ScopeType, playerID PlayerID) []*RoomState {
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

// ApplyModifyBugs adds or removes bug markers from rooms
func ApplyModifyBugs(state *GameState, effect Effect, playerID PlayerID, log *EffectLog) error {
	targets := getRoomTargets(state, effect.Scope, playerID)
	if effect.Scope == RoomWithMostBugs && len(targets) > 0 {
		log.Add("🐛 ModifyBugs: Target room %s selected (bugs=%d)", targets[0].ID, targets[0].BugMarkers)
	}
	for _, room := range targets {
		if room.OutOfRam {
			continue // Skip OutOfRam rooms
		}
		
		oldBugs := room.BugMarkers
		var newBugs uint8
		if effect.N == ALL {
			newBugs = 0
		} else {
			newBugCount := int(room.BugMarkers) + effect.N
			if newBugCount < 0 {
				newBugCount = 0
			}
			if newBugCount > MaxBugMarkers {
				newBugCount = MaxBugMarkers
			}
			newBugs = uint8(newBugCount)
		}
		
		room.BugMarkers = newBugs
		
		// Auto-corruption at 3+ bugs
		wasCorrupted := room.Corrupted
		room.Corrupted = newBugs >= BugCorruptionThreshold
		
		if oldBugs != newBugs {
			if effect.N == ALL {
				log.Add("🪫 %s bugs: %d → 0 (cleared)", room.ID, oldBugs)
			} else if effect.N > 0 {
				log.Add("🪫 %s bugs: %d → %d (+%d)", room.ID, oldBugs, newBugs, effect.N)
			} else {
				log.Add("🪫 %s bugs: %d → %d (%d)", room.ID, oldBugs, newBugs, effect.N)
			}
		}
		
		// Log corruption state changes
		if !wasCorrupted && room.Corrupted {
			log.Add("⚠️ %s corrupted due to %d bugs!", room.ID, room.BugMarkers)
		} else if wasCorrupted && !room.Corrupted {
			log.Add("✨ %s restored (bugs below threshold)", room.ID)
		}
	}
	return nil
}

// ApplyRevealRoom marks rooms as explored
func ApplyRevealRoom(state *GameState, effect Effect, playerID PlayerID, log *EffectLog) error {
	targets := getRoomTargets(state, effect.Scope, playerID)
	log.Add("👁️ RevealRoom: Found %d target rooms for scope %s", len(targets), getScopeName(effect.Scope))
	for _, room := range targets {
		if !room.Explored {
			room.Explored = true
			log.Add("🗺️ %s revealed", room.ID)
		}
	}
	return nil
}

// ApplyCleanRoom removes all bugs from rooms
func ApplyCleanRoom(state *GameState, effect Effect, playerID PlayerID, log *EffectLog) error {
	targets := getRoomTargets(state, effect.Scope, playerID)
	for _, room := range targets {
		oldBugs := room.BugMarkers
		room.BugMarkers = 0
		room.Corrupted = false
		if oldBugs > 0 {
			log.Add("🧹 %s cleaned: %d bugs removed", room.ID, oldBugs)
		}
	}
	return nil
}

// ApplySetCorrupted finds a random room with <3 bugs and sets it to 3 bugs (triggering corruption)
func ApplySetCorrupted(state *GameState, effect Effect, playerID PlayerID, log *EffectLog) error {
	if effect.N != 1 {
		return fmt.Errorf("SetCorrupted only supports n=1 (corrupt room)")
	}
	
	// Find all rooms with < 3 bugs that can be corrupted
	candidateRooms := make([]*RoomState, 0)
	for _, room := range state.Rooms {
		if room.BugMarkers < BugCorruptionThreshold && !room.OutOfRam {
			candidateRooms = append(candidateRooms, room)
		}
	}
	
	if len(candidateRooms) == 0 {
		log.Add("⚠️ SetCorrupted: No rooms available for corruption (all have ≥3 bugs or are OutOfRam)")
		return nil
	}
	
	// Use game RNG to select random room
	rng := GetGameRNG(state)
	selectedRoom := candidateRooms[rng.Intn(len(candidateRooms))]
	
	// Set bugs to exactly 3 (corruption threshold)
	oldBugs := selectedRoom.BugMarkers
	selectedRoom.BugMarkers = BugCorruptionThreshold
	selectedRoom.Corrupted = true
	
	log.Add("⚠️ %s corrupted! Bugs: %d → %d", selectedRoom.ID, oldBugs, selectedRoom.BugMarkers)
	return nil
}

// ApplyOutOfRam forces the room with most enemies out of RAM
func ApplyOutOfRam(state *GameState, effect Effect, playerID PlayerID, log *EffectLog) error {
	targets := getRoomTargets(state, effect.Scope, playerID)
	if len(targets) == 0 {
		log.Add("⚠️ OutOfRam: No target rooms found for scope %s", getScopeName(effect.Scope))
		return nil
	}
	
	for _, room := range targets {
		if !room.OutOfRam {
			room.OutOfRam = true
			log.Add("💻 %s forced out of RAM! System crashes will deal extra damage", room.ID)
		} else {
			log.Add("💻 %s already out of RAM", room.ID)
		}
	}
	
	return nil
}

// getRoomTargets resolves which rooms are affected by the effect
func getRoomTargets(state *GameState, scope ScopeType, playerID PlayerID) []*RoomState {
	player := state.Players[playerID]
	// For event cards (empty playerID), use active player for CurrentRoom/AdjacentRooms
	if player == nil && (scope == CurrentRoom || scope == AdjacentRooms) {
		player = state.Players[state.ActivePlayer]
	}
	if player == nil {
		return nil
	}

	switch scope {
	case CurrentRoom:
		if room := state.Rooms[player.Location]; room != nil {
			return []*RoomState{room}
		}
		return nil
	case AdjacentRooms:
		adjacentIDs := GetAdjacentRooms(player.Location)
		targets := make([]*RoomState, 0, len(adjacentIDs))
		for _, roomID := range adjacentIDs {
			if room := state.Rooms[roomID]; room != nil {
				targets = append(targets, room)
			}
		}
		return targets
	case AllRooms:
		targets := make([]*RoomState, 0, len(state.Rooms))
		for _, room := range state.Rooms {
			targets = append(targets, room)
		}
		return targets
	case RoomWithMostBugs:
		targetRoomID := GetRoomWithMostBugs(state)
		if targetRoomID == "" {
			// No rooms have bugs
			return nil
		}
		if room := state.Rooms[targetRoomID]; room != nil {
			return []*RoomState{room}
		}
		return nil
	case RoomWithMostEnemies:
		targetRoomID := GetRoomWithMostEnemies(state)
		if targetRoomID == "" {
			// No enemies exist
			return nil
		}
		if room := state.Rooms[targetRoomID]; room != nil {
			return []*RoomState{room}
		}
		return nil
	default:
		return nil
	}
}