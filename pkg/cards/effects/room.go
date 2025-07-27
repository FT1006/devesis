package effects

import (
	"github.com/spaceship/devesis/pkg/core"
)

// ApplyModifyBugs adds or removes bug markers from rooms
func ApplyModifyBugs(state *core.GameState, effect Effect, playerID core.PlayerID) error {
	targets := getRoomTargets(state, effect.Scope, playerID)
	for _, room := range targets {
		if room.OutOfRam {
			continue // Skip OutOfRam rooms
		}
		
		var newBugs uint8
		if effect.N == ALL {
			newBugs = 0
		} else {
			newBugCount := int(room.BugMarkers) + effect.N
			if newBugCount < 0 {
				newBugCount = 0
			}
			if newBugCount > core.MaxBugMarkers {
				newBugCount = core.MaxBugMarkers
			}
			newBugs = uint8(newBugCount)
		}
		
		room.BugMarkers = newBugs
		
		// Auto-corruption at 3+ bugs
		room.Corrupted = newBugs >= core.BugCorruptionThreshold
	}
	return nil
}

// ApplyRevealRoom marks rooms as explored
func ApplyRevealRoom(state *core.GameState, effect Effect, playerID core.PlayerID) error {
	targets := getRoomTargets(state, effect.Scope, playerID)
	for _, room := range targets {
		room.Explored = true
	}
	return nil
}

// ApplyCleanRoom removes all bugs from rooms
func ApplyCleanRoom(state *core.GameState, effect Effect, playerID core.PlayerID) error {
	targets := getRoomTargets(state, effect.Scope, playerID)
	for _, room := range targets {
		room.BugMarkers = 0
		room.Corrupted = false
	}
	return nil
}

// ApplySetCorrupted forces rooms into/out of corrupted state
func ApplySetCorrupted(state *core.GameState, effect Effect, playerID core.PlayerID) error {
	targets := getRoomTargets(state, effect.Scope, playerID)
	for _, room := range targets {
		if effect.N == 1 {
			room.Corrupted = true
		} else if effect.N == 0 && room.BugMarkers < core.BugCorruptionThreshold {
			room.Corrupted = false
		}
	}
	return nil
}

// getRoomTargets resolves which rooms are affected by the effect
func getRoomTargets(state *core.GameState, scope ScopeType, playerID core.PlayerID) []*core.RoomState {
	player := state.Players[playerID]
	if player == nil {
		return nil
	}

	switch scope {
	case CurrentRoom:
		if room := state.Rooms[player.Location]; room != nil {
			return []*core.RoomState{room}
		}
		return nil
	case AdjacentRooms:
		adjacentIDs := core.GetAdjacentRooms(player.Location)
		targets := make([]*core.RoomState, 0, len(adjacentIDs))
		for _, roomID := range adjacentIDs {
			if room := state.Rooms[roomID]; room != nil {
				targets = append(targets, room)
			}
		}
		return targets
	case AllRooms:
		targets := make([]*core.RoomState, 0, len(state.Rooms))
		for _, room := range state.Rooms {
			targets = append(targets, room)
		}
		return targets
	case RoomWithMostBugs:
		targetRoomID := core.GetRoomWithMostBugs(state)
		if room := state.Rooms[targetRoomID]; room != nil {
			return []*core.RoomState{room}
		}
		return nil
	default:
		return nil
	}
}