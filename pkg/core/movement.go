package core

import (
	"math"
	"sort"
)

// Pre-computed adjacency map for performance
var adjacencyMap = buildAdjacency()

type PathQuery struct {
	From, To RoomID
	MaxSteps int  // 1 = current behavior, 0 = unlimited
	Diagonal bool // false for Devesis (orthogonal only)
}

type PathResult struct {
	Valid bool
	Path  []RoomID // [From, ..., To] if Valid
}

// CanTraverse performs BFS path-finding up to MaxSteps
func CanTraverse(gs *GameState, q PathQuery) PathResult {
	if q.From == q.To {
		return PathResult{Valid: true, Path: []RoomID{q.From}}
	}

	maxSteps := q.MaxSteps
	if maxSteps <= 0 {
		maxSteps = math.MaxInt32 // unlimited
	}

	// BFS with path tracking
	type pathNode struct {
		room RoomID
		path []RoomID
	}

	queue := []pathNode{{q.From, []RoomID{q.From}}}
	visited := map[RoomID]bool{q.From: true}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if len(current.path)-1 >= maxSteps {
			continue // exceeded max steps
		}

		neighbors := getNeighbors(current.room, q.Diagonal)
		
		for _, neighbor := range neighbors {
			// Skip if room doesn't exist or is corrupted
			if gs.Rooms[neighbor] == nil || gs.Rooms[neighbor].Corrupted {
				continue
			}
			
			if neighbor == q.To {
				return PathResult{
					Valid: true,
					Path:  append(current.path, neighbor),
				}
			}

			if !visited[neighbor] {
				visited[neighbor] = true
				newPath := make([]RoomID, len(current.path)+1)
				copy(newPath, current.path)
				newPath[len(current.path)] = neighbor
				queue = append(queue, pathNode{neighbor, newPath})
			}
		}
	}

	return PathResult{Valid: false, Path: nil}
}

// CanMove maintains backward compatibility (1-step movement)
func CanMove(gs *GameState, from, to RoomID) bool {
	return CanTraverse(gs, PathQuery{
		From: from, To: to, MaxSteps: 1,
	}).Valid
}

// getNeighbors returns adjacent rooms with deterministic ordering
func getNeighbors(roomID RoomID, diagonal bool) []RoomID {
	if !diagonal {
		// Use pre-computed adjacency for orthogonal movement
		return adjacencyMap[roomID]
	}
	
	// For diagonal movement, compute on demand (rare in Devesis)
	pos, exists := ROOM_POSITIONS[string(roomID)]
	if !exists {
		return nil
	}

	directions := []Coord{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	var neighbors []string
	for _, dir := range directions {
		newPos := Coord{pos.Row + dir.Row, pos.Col + dir.Col}
		for id, roomPos := range ROOM_POSITIONS {
			if roomPos == newPos {
				neighbors = append(neighbors, id)
				break
			}
		}
	}

	// Sort for deterministic ordering
	sort.Strings(neighbors)
	
	result := make([]RoomID, len(neighbors))
	for i, id := range neighbors {
		result[i] = RoomID(id)
	}
	return result
}

// GetAdjacentRooms maintains backward compatibility
func GetAdjacentRooms(roomID RoomID) []RoomID {
	return getNeighbors(roomID, false)
}

// buildAdjacency pre-computes all orthogonal adjacencies
func buildAdjacency() map[RoomID][]RoomID {
	adj := make(map[RoomID][]RoomID)
	directions := []Coord{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	
	for id, pos := range ROOM_POSITIONS {
		var neighbors []string
		for _, dir := range directions {
			newPos := Coord{pos.Row + dir.Row, pos.Col + dir.Col}
			for neighborID, neighborPos := range ROOM_POSITIONS {
				if neighborPos == newPos {
					neighbors = append(neighbors, neighborID)
					break
				}
			}
		}
		
		// Sort for deterministic ordering
		sort.Strings(neighbors)
		
		roomNeighbors := make([]RoomID, len(neighbors))
		for i, nid := range neighbors {
			roomNeighbors[i] = RoomID(nid)
		}
		adj[RoomID(id)] = roomNeighbors
	}
	
	return adj
}