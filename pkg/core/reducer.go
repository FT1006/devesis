package core

import (
	"math/rand"
)

func Apply(state GameState, action Action, log *EffectLog) GameState {
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
			oldLocation := player.Location
			player.Location = a.To
			log.Add("🚶 %s moves from %s → %s", a.PlayerID, oldLocation, a.To)
			
			// Mark target room as explored when entering
			if room := newState.Rooms[a.To]; room != nil && !room.Explored {
				room.Explored = true
				log.Add("🏷️  %s marked as explored", a.To)
			}
			
			// Movement consequence: RNG bug placement (3 equal outcomes)
			rng := rand.New(rand.NewSource(newState.RandSeed + int64(newState.Round)*1000 + int64(len(newState.Players))))
			bugOutcome := rng.Intn(3) // 0, 1, or 2
			
			switch bugOutcome {
			case 0:
				// 1 bug in current room (where player moved FROM)
				log.Add("⚠️ Movement consequence: Bug left behind in departure room")
				if oldRoom := newState.Rooms[oldLocation]; oldRoom != nil {
					oldBugs := oldRoom.BugMarkers
					oldRoom.BugMarkers += 1
					if oldRoom.BugMarkers > MaxBugMarkers {
						oldRoom.BugMarkers = MaxBugMarkers
					}
					log.Add("🪲 %s bugs: %d → %d (left behind)", oldLocation, oldBugs, oldRoom.BugMarkers)
				}
			case 1:
				// 1 bug in max 2 surrounding rooms of old location
				log.Add("⚠️ Movement consequence: Bugs spread to adjacent rooms")
				adjacentRooms := GetAdjacentRooms(oldLocation)
				if len(adjacentRooms) > 0 {
					// Shuffle adjacent rooms and pick max 2
					shuffledRooms := make([]RoomID, len(adjacentRooms))
					copy(shuffledRooms, adjacentRooms)
					shuffleRooms(shuffledRooms, rng)
					
					maxRooms := 2
					if len(shuffledRooms) < maxRooms {
						maxRooms = len(shuffledRooms)
					}
					
					bugsSpread := 0
					for i := 0; i < maxRooms; i++ {
						if room := newState.Rooms[shuffledRooms[i]]; room != nil {
							oldBugs := room.BugMarkers
							room.BugMarkers += 1
							if room.BugMarkers > MaxBugMarkers {
								room.BugMarkers = MaxBugMarkers
							}
							log.Add("🪲 %s bugs: %d → %d (spread from %s)", shuffledRooms[i], oldBugs, room.BugMarkers, oldLocation)
							bugsSpread++
						}
					}
					if bugsSpread > 0 {
						log.Add("📡 %d adjacent room(s) affected by movement", bugsSpread)
					}
				}
			case 2:
				// Safe - no bugs added
				log.Add("✅ Movement consequence: Safe passage (no bugs added)")
			}
		}
		return newState

	case GiveSpecialCardAction:
		// Deep copy the state to avoid mutations
		newState := deepCopyGameState(state)
		player, exists := newState.Players[a.PlayerID]
		if !exists {
			return newState
		}

		// Use game RNG to select random special card
		rng := rand.New(rand.NewSource(newState.RandSeed + int64(newState.Round)*500 + int64(len(newState.Players))))
		
		// Collect all special card IDs
		specialCards := make([]CardID, 0)
		for cardID, card := range CardDB {
			if card.Source == SrcSpecial {
				specialCards = append(specialCards, cardID)
			}
		}
		
		if len(specialCards) > 0 {
			// Pick random special card
			cardIndex := rng.Intn(len(specialCards))
			selectedCard := specialCards[cardIndex]
			
			// Add to hand
			player.Hand = append(player.Hand, selectedCard)
			
			// Enforce hand limit
			enforceHandLimitWithDiscard(&player.Hand, &player.Discard)
		}
		
		return newState
		
	case SearchAction:
		// Use the proper search logic with RNG
		rng := rand.New(rand.NewSource(state.RandSeed + int64(state.Round)*1000))
		return ApplySearch(state, a, rng, log)

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

		// Move card from hand to discard pile using helper
		moveCardByIndex(&player.Hand, &player.Discard, cardIndex)
		
		// Log card play with name and description
		if card, exists := CardDB[a.CardID]; exists {
			log.Add("🃏 %s plays %s - %s", a.PlayerID, card.Name, card.Description)
		} else {
			log.Add("🃏 %s plays %s", a.PlayerID, a.CardID)
		}

		// Check for engine card usage at escape room
		if a.CardID == "SPECIAL_ENGINE" {
			if player.Location == "R19" || player.Location == "R20" {
				// Check if no Pythogoras in escape room
				pythogorasInEscapeRoom := false
				for _, enemy := range newState.Enemies {
					if enemy.Type == Pythogoras && (enemy.Location == "R19" || enemy.Location == "R20") {
						pythogorasInEscapeRoom = true
						break
					}
				}
				if !pythogorasInEscapeRoom {
					player.EngineUsed = true // Mark as victory condition met
					log.Add("🚀 Engine activated! Victory condition met")
				}
			}
		}

		// Apply card effects using the effects engine
		if card, exists := CardDB[a.CardID]; exists {
			newState = ApplyCardEffects(newState, card, a.PlayerID, log)
		}

		return newState

	case ShootAction, MeleeAction:
		// Handle combat actions
		return ApplyCombat(state, a, log)
		
	case RoomAction:
		// Handle room-specific actions
		return ApplyRoomAction(state, a, log)
}
	
	// Default case - return original state unchanged
	return state
}

// ApplyWithoutLog is a convenience wrapper for backwards compatibility
func ApplyWithoutLog(state GameState, action Action) GameState {
	log := NewEffectLog()
	return Apply(state, action, log)
}

// DeepCopyGameState creates a deep copy of the game state (exported)
func DeepCopyGameState(state GameState) GameState {
	return deepCopyGameState(state)
}

// Deep copy helper function
func deepCopyGameState(state GameState) GameState {
	newState := GameState{
		Round:         state.Round,
		Time:          state.Time,
		RandSeed:      state.RandSeed,
		EventIndex:    state.EventIndex,
		ActionsLeft:   state.ActionsLeft,
		Phase:         state.Phase,
		ActivePlayer:  state.ActivePlayer,
		Rooms:         make(map[RoomID]*RoomState),
		Players:       make(map[PlayerID]*PlayerState),
		Events:        make([]EventCard, len(state.Events)),
		SpawnBag:      nil,
		Enemies:       make(map[EnemyID]*Enemy),
		QuestionOrder: make([]int, len(state.QuestionOrder)),
		NextQuestion:  state.NextQuestion,
		ScratchLog:    NewEffectLog(), // Initialize effect log
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
			Damage:       player.Damage,
			Hand:         make([]CardID, len(player.Hand)),
			Deck:         make([]CardID, len(player.Deck)),
			Discard:      make([]CardID, len(player.Discard)),
			Location:     player.Location,
			HasActed:     player.HasActed,
			SpecialUsed:  player.SpecialUsed,
			EngineUsed:   player.EngineUsed,
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
		Time:          15, // Start with 15 time units
		RandSeed:      seed,
		EventIndex:    0,
		Rooms:         make(map[RoomID]*RoomState),
		Players:       make(map[PlayerID]*PlayerState),
		Events:        initializeEventCards(),
		SpawnBag:      initializeSpawnBag(),
		Enemies:       make(map[EnemyID]*Enemy),
		// Initialize pre-shuffled question order
		QuestionOrder: initializeQuestionOrder(seed),
		NextQuestion:  0,
		ScratchLog:    NewEffectLog(), // Initialize effect log
	}

	// Initialize all 20 rooms from the spaceship layout
	// Create room type distribution for non-predefined rooms
	roomTypePool := []RoomType{
		AmmoCache, AmmoCache, AmmoCache,        // 3 AmmoCache
		MedBay, MedBay, MedBay,                 // 3 MedBay  
		CleanRoomType,                          // 1 CleanRoomType
		EnemySpawn, EnemySpawn, EnemySpawn, EnemySpawn, // 4 EnemySpawn
		Empty, Empty,                           // 2 Empty
	}
	
	// Use game RNG for consistent room assignment
	tempState := GameState{RandSeed: seed, Round: 1, Time: 15}
	rng := GetGameRNG(&tempState)
	
	// Shuffle the room type pool for random assignment
	shuffleRoomTypes(roomTypePool, rng)
	
	poolIndex := 0
	for roomIDStr := range ROOM_POSITIONS {
		roomID := RoomID(roomIDStr)
		roomType := Empty
		explored := false

		// Check if this room has a predefined type
		if predefinedType, exists := PREDEFINED_ROOMS[roomIDStr]; exists {
			roomType = predefinedType
			explored = true // Predefined rooms are already explored
		} else {
			// Assign from shuffled pool
			roomType = roomTypePool[poolIndex]
			poolIndex++
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
		Damage:       BasicDamage, // Base damage
		Hand:         []CardID{},
		Deck:         createRandomStartingDeck(seed),
		Discard:      []CardID{},
		Location:     "R12", // Start room
		HasActed:     false,
		SpecialUsed:  false,
		EngineUsed:   false,
		PersonalObj:  ObjectiveID(""),
		CorporateObj: ObjectiveID(""),
	}

	// Set turn controller state
	state.ActivePlayer = playerID  // "P1" for solo mode
	state.Phase = "player"
	state.ActionsLeft = 0 // Will be set by DrawPhase

	return state
}

// GetActivePlayer returns the currently active player
func GetActivePlayer(state *GameState) *PlayerState {
	if player, exists := state.Players[state.ActivePlayer]; exists {
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

// initializeEventCards loads all event cards from the CardDB
func initializeEventCards() []EventCard {
	eventCards := make([]EventCard, 0)
	for cardID, card := range CardDB {
		if card.Source == SrcEvent {
			eventCards = append(eventCards, EventCard{
				ID:          cardID,
				Name:        card.Name,
				Description: card.Description,
				Effects:     card.Effects,
			})
		}
	}
	return eventCards
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

// createRandomStartingDeck creates a random 10-card starting deck from action cards
func createRandomStartingDeck(seed int64) []CardID {
	rng := rand.New(rand.NewSource(seed + 1000)) // Offset seed for deck generation
	
	// Get all action card IDs from the loaded CardDB
	actionCards := make([]CardID, 0)
	for cardID, card := range CardDB {
		if card.Source == SrcAction {
			actionCards = append(actionCards, cardID)
		}
	}
	
	// If no cards loaded, return empty deck (fallback)
	if len(actionCards) == 0 {
		return []CardID{}
	}
	
	// Shuffle action cards using Fisher-Yates
	for i := len(actionCards) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		actionCards[i], actionCards[j] = actionCards[j], actionCards[i]
	}
	
	// Take first 10 (or all if fewer than 10)
	deckSize := 10
	if len(actionCards) < 10 {
		deckSize = len(actionCards)
	}
	
	return actionCards[:deckSize]
}