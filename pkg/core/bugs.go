package core

import (
	"math/rand"
)

// PlaceBugs adds bug markers to rooms based on game events
func PlaceBugs(state *GameState, count uint8, rng *rand.Rand) {
	if count == 0 {
		return
	}

	// Get all valid rooms (not corrupted, not out of RAM)
	validRooms := getValidRoomsForBugs(state)
	if len(validRooms) == 0 {
		return
	}

	// Place bugs randomly
	for i := uint8(0); i < count; i++ {
		roomIdx := rng.Intn(len(validRooms))
		roomID := validRooms[roomIdx]
		room := state.Rooms[roomID]
		
		// Add bug marker (max 255)
		if room.BugMarkers < MaxBugMarkers {
			room.BugMarkers++
		}
	}
}

// GetRoomWithMostBugs finds the room with highest bug count for enemy targeting
func GetRoomWithMostBugs(state *GameState) RoomID {
	var mostBugs uint8 = 0
	var targetRoom RoomID = "R12" // Default to start room
	
	for id, room := range state.Rooms {
		if room.BugMarkers > mostBugs {
			mostBugs = room.BugMarkers
			targetRoom = id
		}
	}
	
	return targetRoom
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
		if !room.Corrupted && !room.OutOfRam {
			valid = append(valid, id)
		}
	}
	return valid
}