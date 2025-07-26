package core

func Apply(state GameState, action Action) GameState {
	// Deep copy the state to avoid mutations
	newState := deepCopyGameState(state)
	
	switch a := action.(type) {
	case MoveAction:
		player, exists := newState.Players[a.PlayerID]
		if !exists {
			return newState
		}
		
		// Validate move using CanMove
		if CanMove(&newState, player.Location, a.To) {
			player.Location = a.To
		}
		
	case SearchAction:
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
	}
	
	return newState
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
		Bag:           make([]Token, len(state.Bag)),
		Enemies:       make(map[EnemyID]*Enemy),
		UsedQuestions: make([]int, len(state.UsedQuestions)),
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
			Hand:         make([]Card, len(player.Hand)),
			Deck:         make([]Card, len(player.Deck)),
			Discard:      make([]Card, len(player.Discard)),
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
	
	// Copy events, bag, and used questions
	copy(newState.Events, state.Events)
	copy(newState.Bag, state.Bag)
	copy(newState.UsedQuestions, state.UsedQuestions)
	
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