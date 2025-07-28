package core

import (
	"math"
	"sort"
)

// Direction sets for neighbor computation
var (
	orthoDirs = []Coord{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	diagDirs  = []Coord{
		{-1, -1}, {-1, 1},
		{1, -1}, {1, 1},
	}
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
			// Skip if room doesn't exist
			if gs.Rooms[neighbor] == nil {
				continue
			}
			// Note: Corrupted rooms are passable (blocking not implemented for solo play)
			
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

// computeNeighbors is the shared neighbor computation algorithm
func computeNeighbors(pos Coord, dirs []Coord) []RoomID {
	var ids []string
	for _, d := range dirs {
		target := Coord{pos.Row + d.Row, pos.Col + d.Col}
		for id, p := range ROOM_POSITIONS {
			if p == target {
				ids = append(ids, id)
				break
			}
		}
	}
	sort.Strings(ids) // deterministic order
	out := make([]RoomID, len(ids))
	for i, id := range ids {
		out[i] = RoomID(id)
	}
	return out
}

// buildAdjacency pre-computes all orthogonal adjacencies
func buildAdjacency() map[RoomID][]RoomID {
	adj := make(map[RoomID][]RoomID, len(ROOM_POSITIONS))
	for id, pos := range ROOM_POSITIONS {
		adj[RoomID(id)] = computeNeighbors(pos, orthoDirs)
	}
	return adj
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
	return computeNeighbors(pos, diagDirs)
}

// GetAdjacentRooms maintains backward compatibility (orthogonal only, fast)
func GetAdjacentRooms(roomID RoomID) []RoomID {
	return adjacencyMap[roomID]
}

// GetDiagonalNeighbors provides diagonal lookup (rare, compute on demand)
func GetDiagonalNeighbors(roomID RoomID) []RoomID {
	pos, ok := ROOM_POSITIONS[string(roomID)]
	if !ok {
		return nil
	}
	return computeNeighbors(pos, diagDirs)
}