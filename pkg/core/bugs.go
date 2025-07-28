package core

// PlaceBugs adds bug markers to rooms based on game events
func PlaceBugs(state *GameState, count uint8) {
	if count == 0 {
		return
	}

	// Get game RNG
	rng := GetGameRNG(state)
	
	// Track rooms that need spawn checks
	var spawnCheckRooms []RoomID

	// Get all valid rooms (not out of RAM)
	validRooms := getValidRoomsForBugs(state)
	if len(validRooms) == 0 {
		return
	}

	// Place bugs randomly
	for i := uint8(0); i < count; i++ {
		roomIdx := rng.Intn(len(validRooms))
		roomID := validRooms[roomIdx]
		room := state.Rooms[roomID]
		
		// Check if room is corrupted before adding bug
		wasCorruptedBefore := room.Corrupted
		
		// Add bug marker (max 9)
		if room.BugMarkers < MaxBugMarkers {
			room.BugMarkers++
			
			// Check corruption threshold
			if room.BugMarkers >= BugCorruptionThreshold {
				room.Corrupted = true
			}
			
			// If bug was added to already corrupted room, mark for spawn check
			if wasCorruptedBefore {
				spawnCheckRooms = append(spawnCheckRooms, roomID)
			}
		}
	}
	
	// Handle spawning for corrupted rooms that received bugs
	if len(spawnCheckRooms) > 0 {
		SpawnEnemiesForCorruptedRooms(state, spawnCheckRooms)
	}
}

// CheckOutOfRAMCondition returns true if 5+ rooms are OutOfRam (game over)
func CheckOutOfRAMCondition(state *GameState) bool {
	count := 0
	for _, room := range state.Rooms {
		if room.OutOfRam {
			count++
		}
	}
	return count >= 5
}

func getValidRoomsForBugs(state *GameState) []RoomID {
	var valid []RoomID
	for id, room := range state.Rooms {
		// Can place bugs in any room that's not out of RAM
		// Corrupted rooms can still receive bugs (which triggers spawns)
		if !room.OutOfRam {
			valid = append(valid, id)
		}
	}
	return valid
}

// SpawnEnemiesForCorruptedRooms handles enemy spawning when bugs are added to corrupted rooms
func SpawnEnemiesForCorruptedRooms(state *GameState, roomIDs []RoomID) {
	for _, roomID := range roomIDs {
		SpawnEnemyFromBag(state, roomID)
		// If bag is empty, stop spawning
		if IsSpawnBagEmpty(state) {
			break
		}
	}
}

// UpdateRoomCorruption checks and updates corruption status based on bug count
func UpdateRoomCorruption(state *GameState) {
	for _, room := range state.Rooms {
		// Auto-corrupt at 3+ bugs
		if room.BugMarkers >= BugCorruptionThreshold {
			room.Corrupted = true
		} else {
			// Auto-uncorrupt below 3 bugs
			room.Corrupted = false
		}
	}
}

// PlaceBugsInSpecificRooms adds bugs to specific rooms and handles spawning
func PlaceBugsInSpecificRooms(state *GameState, roomIDs []RoomID) {
	var spawnCheckRooms []RoomID
	
	for _, roomID := range roomIDs {
		room := state.Rooms[roomID]
		if room == nil || room.OutOfRam {
			continue
		}
		
		// Check if room is corrupted before adding bug
		wasCorruptedBefore := room.Corrupted
		
		// Add 2 bug markers (max 9) - increased penalty for wrong answers
		bugsToAdd := uint8(2)
		if room.BugMarkers + bugsToAdd > MaxBugMarkers {
			bugsToAdd = MaxBugMarkers - room.BugMarkers
		}
		if bugsToAdd > 0 {
			room.BugMarkers += bugsToAdd
			
			// Check corruption threshold
			if room.BugMarkers >= BugCorruptionThreshold {
				room.Corrupted = true
			}
			
			// If bug was added to already corrupted room, mark for spawn check
			if wasCorruptedBefore {
				spawnCheckRooms = append(spawnCheckRooms, roomID)
			}
		}
	}
	
	// Handle spawning for corrupted rooms that received bugs
	if len(spawnCheckRooms) > 0 {
		SpawnEnemiesForCorruptedRooms(state, spawnCheckRooms)
	}
}

// ApplyWrongAnswerPenalties handles all penalties for incorrect movement questions
func ApplyWrongAnswerPenalties(state *GameState, targetRoom RoomID) {
	// Add bugs to target room and all adjacent rooms
	roomsToInfect := []RoomID{targetRoom}
	adjacent := GetAdjacentRooms(targetRoom)
	roomsToInfect = append(roomsToInfect, adjacent...)
	
	// Use proper bug placement (respects limits, handles corruption, triggers spawns)
	PlaceBugsInSpecificRooms(state, roomsToInfect)
	
	// Drop all cards from active player's hand to discard pile
	if player := GetActivePlayer(state); player != nil {
		moveAllCards(&player.Hand, &player.Discard)
	}
}