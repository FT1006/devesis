package core

import (
	"math/rand"
)

func Apply(state GameState, action Action) GameState {
	switch a := action.(type) {
	case InitializeGameAction:
		// Create initial game state - no deep copy needed
		return initializeGameState(a.Seed, a.PlayerClass)

	case MoveAction:
		// Deep copy the state to avoid mutations
		newState := deepCopyGameState(state)
		player, exists := newState.Players[a.PlayerID]
		if !exists {
			return newState
		}
		
		// Validate move using CanMove
		if CanMove(&newState, player.Location, a.To) {
			player.Location = a.To
		}
		return newState
		
	case SearchAction:
		// Deep copy the state to avoid mutations
		newState := deepCopyGameState(state)
		player, exists := newState.Players[a.PlayerID]
		if !exists {
			return newState
		}
		
		// Mark current room as searched
		if room, roomExists := newState.Rooms[player.Location]; roomExists {
			room.Searched = true
		}
		
		// Discard 1 card if player has cards
		if len(player.Hand) > 0 {
			player.Hand = player.Hand[:len(player.Hand)-1]
		}
		return newState

	case PlayCardAction:
		// Deep copy the state to avoid mutations
		newState := deepCopyGameState(state)
		player, exists := newState.Players[a.PlayerID]
		if !exists {
			return newState
		}

		// Find and remove card from player's hand
		cardIndex := -1
		for i, cardID := range player.Hand {
			if cardID == a.CardID {
				cardIndex = i
				break
			}
		}

		if cardIndex == -1 {
			// Card not found in hand
			return newState
		}

		// Remove card from hand
		player.Hand = append(player.Hand[:cardIndex], player.Hand[cardIndex+1:]...)

		// Apply card effects using the effects engine
		if card, exists := CardDB[a.CardID]; exists {
			newState = ApplyCardEffects(newState, card, a.PlayerID)
		}

		return newState
}
	
	// Default case - return original state unchanged
	return state
}

// Deep copy helper function
func deepCopyGameState(state GameState) GameState {
	newState := GameState{
		Round:         state.Round,
		Time:          state.Time,
		RandSeed:      state.RandSeed,
		EventIndex:    state.EventIndex,
		Rooms:         make(map[RoomID]*RoomState),
		Players:       make(map[PlayerID]*PlayerState),
		Events:        make([]EventCard, len(state.Events)),
		SpawnBag:      nil,
		Enemies:       make(map[EnemyID]*Enemy),
		QuestionOrder: make([]int, len(state.QuestionOrder)),
		NextQuestion:  state.NextQuestion,
	}
	
	// Copy rooms
	for id, room := range state.Rooms {
		newState.Rooms[id] = &RoomState{
			ID:         room.ID,
			Type:       room.Type,
			Explored:   room.Explored,
			Searched:   room.Searched,
			Corrupted:  room.Corrupted,
			OutOfRam:   room.OutOfRam,
			BugMarkers: room.BugMarkers,
		}
	}
	
	// Copy players
	for id, player := range state.Players {
		newState.Players[id] = &PlayerState{
			ID:           player.ID,
			Class:        player.Class,
			HP:           player.HP,
			MaxHP:        player.MaxHP,
			Ammo:         player.Ammo,
			MaxAmmo:      player.MaxAmmo,
			Hand:         make([]CardID, len(player.Hand)),
			Deck:         make([]CardID, len(player.Deck)),
			Discard:      make([]CardID, len(player.Discard)),
			Location:     player.Location,
			HasActed:     player.HasActed,
			SpecialUsed:  player.SpecialUsed,
			PersonalObj:  player.PersonalObj,
			CorporateObj: player.CorporateObj,
		}
		copy(newState.Players[id].Hand, player.Hand)
		copy(newState.Players[id].Deck, player.Deck)
		copy(newState.Players[id].Discard, player.Discard)
	}
	
	// Copy events and question order
	copy(newState.Events, state.Events)
	copy(newState.QuestionOrder, state.QuestionOrder)
	
	// Deep copy spawn bag
	if state.SpawnBag != nil {
		newState.SpawnBag = &SpawnBag{
			Tokens: make([]EnemyType, len(state.SpawnBag.Tokens)),
		}
		copy(newState.SpawnBag.Tokens, state.SpawnBag.Tokens)
	}

	// Copy enemies
	for id, enemy := range state.Enemies {
		newState.Enemies[id] = &Enemy{
			ID:       enemy.ID,
			Type:     enemy.Type,
			HP:       enemy.HP,
			MaxHP:    enemy.MaxHP,
			Damage:   enemy.Damage,
			Location: enemy.Location,
		}
	}
	
	return newState
}

// initializeGameState creates a fresh game state
func initializeGameState(seed int64, playerClass DevClass) GameState {
	state := GameState{
		Round:         1,
		Time:          0,
		RandSeed:      seed,
		EventIndex:    0,
		Rooms:         make(map[RoomID]*RoomState),
		Players:       make(map[PlayerID]*PlayerState),
		Events:        []EventCard{},
		SpawnBag:      initializeSpawnBag(),
		Enemies:       make(map[EnemyID]*Enemy),
		// Initialize pre-shuffled question order
		QuestionOrder: initializeQuestionOrder(seed),
		NextQuestion:  0,
	}

	// Initialize all 20 rooms from the spaceship layout
	for roomIDStr := range ROOM_POSITIONS {
		roomID := RoomID(roomIDStr)
		roomType := Empty
		explored := false

		// Check if this room has a predefined type
		if predefinedType, exists := PREDEFINED_ROOMS[roomIDStr]; exists {
			roomType = predefinedType
			explored = true // Predefined rooms are already explored
		}

		state.Rooms[roomID] = &RoomState{
			ID:         roomID,
			Type:       roomType,
			Explored:   explored,
			Searched:   false, // No rooms are searched initially
			Corrupted:  false,
			OutOfRam:   false,
			BugMarkers: 0,
		}
	}

	// Initialize single player with class-specific stats
	playerID := PlayerID("P1")
	classStats := CLASS_STATS[playerClass]

	state.Players[playerID] = &PlayerState{
		ID:           playerID,
		Class:        playerClass,
		HP:           classStats.HP,
		MaxHP:        classStats.HP,
		Ammo:         classStats.MaxAmmo,
		MaxAmmo:      classStats.MaxAmmo,
		Hand:         []CardID{},
		Deck:         []CardID{},
		Discard:      []CardID{},
		Location:     "R12", // Start room
		HasActed:     false,
		SpecialUsed:  false,
		PersonalObj:  ObjectiveID(""),
		CorporateObj: ObjectiveID(""),
	}

	return state
}

// GetActivePlayer returns the currently active player (P1 for single player)
func GetActivePlayer(state *GameState) *PlayerState {
	if player, exists := state.Players["P1"]; exists {
		return player
	}
	return nil
}

// IsGameOver checks if the game has ended
func IsGameOver(state *GameState) bool {
	// Check time limit
	if state.Round >= MaxRounds {
		return true
	}

	// Check if all players are dead
	for _, player := range state.Players {
		if player.HP > 0 {
			return false // At least one player alive
		}
	}

	return true // All players dead
}

// ClassOption represents a selectable developer class
type ClassOption struct {
	ID          int
	Class       DevClass
	DisplayName string
	HP          uint8
	MaxAmmo     uint8
	Description string
}

// GetAvailableClasses returns all selectable developer classes with their stats
func GetAvailableClasses() []ClassOption {
	return []ClassOption{
		{
			ID:          1,
			Class:       Frontend,
			DisplayName: "Frontend",
			HP:          CLASS_STATS[Frontend].HP,
			MaxAmmo:     CLASS_STATS[Frontend].MaxAmmo,
			Description: "Balanced survivability",
		},
		{
			ID:          2,
			Class:       Backend,
			DisplayName: "Backend",
			HP:          CLASS_STATS[Backend].HP,
			MaxAmmo:     CLASS_STATS[Backend].MaxAmmo,
			Description: "High firepower, low health",
		},
		{
			ID:          3,
			Class:       DevOps,
			DisplayName: "DevOps",
			HP:          CLASS_STATS[DevOps].HP,
			MaxAmmo:     CLASS_STATS[DevOps].MaxAmmo,
			Description: "Well-rounded specialist",
		},
		{
			ID:          4,
			Class:       Fullstack,
			DisplayName: "Fullstack",
			HP:          CLASS_STATS[Fullstack].HP,
			MaxAmmo:     CLASS_STATS[Fullstack].MaxAmmo,
			Description: "Jack of all trades",
		},
	}
}

// ValidateClassChoice validates a class selection ID and returns the DevClass
func ValidateClassChoice(choiceID int) (DevClass, bool) {
	classes := GetAvailableClasses()
	for _, class := range classes {
		if class.ID == choiceID {
			return class.Class, true
		}
	}
	return Frontend, false // Default fallback
}

// initializeSpawnBag creates the initial enemy spawn pool
func initializeSpawnBag() *SpawnBag {
	bag := &SpawnBag{
		Tokens: []EnemyType{},
	}

	// Add enemy tokens based on difficulty distribution
	// More weak enemies, fewer strong ones

	// 10 Infinite Loops (weakest)
	for i := 0; i < 10; i++ {
		bag.Tokens = append(bag.Tokens, InfiniteLoop)
	}

	// 6 Stack Overflows (medium)
	for i := 0; i < 6; i++ {
		bag.Tokens = append(bag.Tokens, StackOverflow)
	}

	// 2 Pythogoras (strongest)
	for i := 0; i < 2; i++ {
		bag.Tokens = append(bag.Tokens, Pythogoras)
	}

	return bag
}

// initializeQuestionOrder creates a pre-shuffled order of question IDs 0-49
func initializeQuestionOrder(seed int64) []int {
	rng := rand.New(rand.NewSource(seed))
	return rng.Perm(50) // Creates [0,1,2,...,49] in random order
}