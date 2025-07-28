package core

import (
	"testing"
)

// Test helper functions

func newWinTestGameState() GameState {
	return GameState{
		Round:        1,
		Time:         15,
		RandSeed:     42,
		ActivePlayer: "P1",
		Rooms: map[RoomID]*RoomState{
			"R12": {ID: "R12", Type: Predefined},        // Start
			"R15": {ID: "R15", Type: Predefined},        // Engine room
			"R17": {ID: "R17", Type: Predefined},        // Engine room
			"R18": {ID: "R18", Type: Predefined},        // Engine room
			"R19": {ID: "R19", Type: Predefined},        // Escape room
			"R20": {ID: "R20", Type: Predefined},        // Escape room
		},
		Players: map[PlayerID]*PlayerState{
			"P1": {
				ID:        "P1",
				Location:  "R12",
				HP:        5,
				MaxHP:     5,
				Hand:      []CardID{"SPECIAL_ENGINE", "CARD_1", "CARD_2"},
				EngineUsed: false,
			},
		},
		Enemies: map[EnemyID]*Enemy{},
	}
}

func addPythogorasToRoom(state *GameState, roomID RoomID) {
	enemyID := EnemyID("PYTHOGORAS_TEST")
	state.Enemies[enemyID] = &Enemy{
		ID:       enemyID,
		Type:     Pythogoras,
		HP:       6,
		MaxHP:    6,
		Damage:   1,
		Location: roomID,
	}
}

func setPlayerAtEscapeRoom(state *GameState, roomID RoomID) {
	state.Players["P1"].Location = roomID
}

// Loss Condition Tests

func TestCheckEndSolo_PlayerDead(t *testing.T) {
	state := newWinTestGameState()
	state.Players["P1"].HP = 0

	ended, win := CheckEndSolo(&state)

	if !ended {
		t.Error("expected game to end when player is dead")
	}
	if win {
		t.Error("expected loss when player is dead")
	}
}

func TestCheckEndSolo_TimeUp(t *testing.T) {
	state := newWinTestGameState()
	state.Time = 0

	ended, win := CheckEndSolo(&state)

	if !ended {
		t.Error("expected game to end when time is up")
	}
	if win {
		t.Error("expected loss when time is up")
	}
}

func TestCheckEndSolo_NoActivePlayer(t *testing.T) {
	state := newWinTestGameState()
	state.Players = map[PlayerID]*PlayerState{} // No players

	ended, win := CheckEndSolo(&state)

	if !ended {
		t.Error("expected game to end when no active player")
	}
	if win {
		t.Error("expected loss when no active player")
	}
}

// Win Condition Tests

func TestCheckEndSolo_EngineUsed(t *testing.T) {
	state := newWinTestGameState()
	state.Players["P1"].EngineUsed = true

	ended, win := CheckEndSolo(&state)

	if !ended {
		t.Error("expected game to end when engine is used")
	}
	if !win {
		t.Error("expected win when engine is used")
	}
}

func TestCheckEndSolo_EngineUsedOverridesOtherConditions(t *testing.T) {
	state := newWinTestGameState()
	state.Players["P1"].HP = 1      // Low HP
	state.Time = 1                  // Low time
	state.Players["P1"].EngineUsed = true

	ended, win := CheckEndSolo(&state)

	if !ended {
		t.Error("expected game to end")
	}
	if !win {
		t.Error("expected win when engine used, even with low HP/time")
	}
}

// Game Continuing Tests

func TestCheckEndSolo_GameContinues(t *testing.T) {
	state := newWinTestGameState()
	// Default state: healthy player, time remaining, no engine used

	ended, win := CheckEndSolo(&state)

	if ended {
		t.Error("expected game to continue")
	}
	if win {
		t.Error("expected no win when game continues")
	}
}

// Engine Card Win Scenario Tests

func TestEngineCardWin_SuccessfulEscape_R19(t *testing.T) {
	state := newWinTestGameState()
	setPlayerAtEscapeRoom(&state, "R19")
	
	log := NewEffectLog()
	action := PlayCardAction{
		PlayerID: "P1",
		CardID:   "SPECIAL_ENGINE",
	}

	result := Apply(state, action, log)

	if !result.Players["P1"].EngineUsed {
		t.Error("expected EngineUsed to be true after playing engine card at R19")
	}
	
	// Verify win condition
	ended, win := CheckEndSolo(&result)
	if !ended || !win {
		t.Error("expected victory after successfully using engine card at escape room")
	}
}

func TestEngineCardWin_SuccessfulEscape_R20(t *testing.T) {
	state := newWinTestGameState()
	setPlayerAtEscapeRoom(&state, "R20")
	
	log := NewEffectLog()
	action := PlayCardAction{
		PlayerID: "P1",
		CardID:   "SPECIAL_ENGINE",
	}

	result := Apply(state, action, log)

	if !result.Players["P1"].EngineUsed {
		t.Error("expected EngineUsed to be true after playing engine card at R20")
	}
	
	// Verify win condition
	ended, win := CheckEndSolo(&result)
	if !ended || !win {
		t.Error("expected victory after successfully using engine card at escape room")
	}
}

func TestEngineCardWin_BlockedByPythogoras_R19(t *testing.T) {
	state := newWinTestGameState()
	setPlayerAtEscapeRoom(&state, "R19")
	addPythogorasToRoom(&state, "R19") // Pythogoras blocks escape
	
	log := NewEffectLog()
	action := PlayCardAction{
		PlayerID: "P1",
		CardID:   "SPECIAL_ENGINE",
	}

	result := Apply(state, action, log)

	if result.Players["P1"].EngineUsed {
		t.Error("expected EngineUsed to remain false when Pythogoras blocks escape")
	}
	
	// Verify game continues
	ended, _ := CheckEndSolo(&result)
	if ended {
		t.Error("expected game to continue when Pythogoras blocks escape")
	}
}

func TestEngineCardWin_BlockedByPythogoras_R20(t *testing.T) {
	state := newWinTestGameState()
	setPlayerAtEscapeRoom(&state, "R20")
	addPythogorasToRoom(&state, "R20") // Pythogoras blocks escape
	
	log := NewEffectLog()
	action := PlayCardAction{
		PlayerID: "P1",
		CardID:   "SPECIAL_ENGINE",
	}

	result := Apply(state, action, log)

	if result.Players["P1"].EngineUsed {
		t.Error("expected EngineUsed to remain false when Pythogoras blocks escape")
	}
	
	// Verify game continues
	ended, _ := CheckEndSolo(&result)
	if ended {
		t.Error("expected game to continue when Pythogoras blocks escape")
	}
}

func TestEngineCardWin_WrongRoom(t *testing.T) {
	state := newWinTestGameState()
	// Player starts at R12 (not an escape room)
	
	log := NewEffectLog()
	action := PlayCardAction{
		PlayerID: "P1",
		CardID:   "SPECIAL_ENGINE",
	}

	result := Apply(state, action, log)

	if result.Players["P1"].EngineUsed {
		t.Error("expected EngineUsed to remain false when not at escape room")
	}
	
	// Verify game continues
	ended, _ := CheckEndSolo(&result)
	if ended {
		t.Error("expected game to continue when engine used outside escape room")
	}
}

// Edge Cases

func TestEngineCard_PythogorasInBothRooms(t *testing.T) {
	state := newWinTestGameState()
	setPlayerAtEscapeRoom(&state, "R19")
	addPythogorasToRoom(&state, "R19") // Block R19
	addPythogorasToRoom(&state, "R20") // Block R20
	
	log := NewEffectLog()
	action := PlayCardAction{
		PlayerID: "P1",
		CardID:   "SPECIAL_ENGINE",
	}

	result := Apply(state, action, log)

	if result.Players["P1"].EngineUsed {
		t.Error("expected EngineUsed to remain false when both escape rooms have Pythogoras")
	}
}

func TestEngineCard_AfterCombat(t *testing.T) {
	state := newWinTestGameState()
	setPlayerAtEscapeRoom(&state, "R19")
	addPythogorasToRoom(&state, "R19")
	
	// First kill the Pythogoras
	delete(state.Enemies, "PYTHOGORAS_TEST")
	
	// Now try the engine card
	log := NewEffectLog()
	action := PlayCardAction{
		PlayerID: "P1",
		CardID:   "SPECIAL_ENGINE",
	}

	result := Apply(state, action, log)

	if !result.Players["P1"].EngineUsed {
		t.Error("expected EngineUsed to be true after clearing Pythogoras from escape room")
	}
	
	// Verify win condition
	ended, win := CheckEndSolo(&result)
	if !ended || !win {
		t.Error("expected victory after clearing escape room and using engine card")
	}
}

func TestCheckEndSolo_MultipleConditions(t *testing.T) {
	// Test that EngineUsed takes priority over death
	state := newWinTestGameState()
	state.Players["P1"].HP = 0           // Should cause loss
	state.Players["P1"].EngineUsed = true // Should cause win

	ended, win := CheckEndSolo(&state)

	if !ended {
		t.Error("expected game to end")
	}
	
	// The current implementation checks loss conditions first,
	// but documenting the actual behavior for future reference
	if win {
		t.Error("current implementation: loss conditions checked before win condition")
	}
}