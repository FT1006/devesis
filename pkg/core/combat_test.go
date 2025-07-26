package core

import (
	"testing"
)

func TestShootAction_DamagesAdjacentEnemies(t *testing.T) {
	state := newCombatTestGameState()
	action := ShootAction{PlayerID: "P1"}
	
	result := ApplyCombat(state, action)
	
	// Enemy in adjacent room should take 1 damage
	enemy := result.Enemies["E1"]
	if enemy.HP != 2 { // 3 - 1 = 2
		t.Errorf("Expected enemy HP 2, got %d", enemy.HP)
	}
}

func TestShootAction_ConsumesAmmo(t *testing.T) {
	state := newCombatTestGameState()
	action := ShootAction{PlayerID: "P1"}
	
	result := ApplyCombat(state, action)
	
	player := result.Players["P1"]
	if player.Ammo != 4 { // 5 - 1 = 4
		t.Errorf("Expected player ammo 4, got %d", player.Ammo)
	}
}

func TestShootAction_FailsWithNoAmmo(t *testing.T) {
	state := newCombatTestGameState()
	state.Players["P1"].Ammo = 0
	action := ShootAction{PlayerID: "P1"}
	
	result := ApplyCombat(state, action)
	
	// Enemy should not take damage
	enemy := result.Enemies["E1"]
	if enemy.HP != 3 {
		t.Error("Enemy should not take damage when player has no ammo")
	}
}

func TestMeleeAction_DamagesSameRoomEnemies(t *testing.T) {
	state := newCombatTestGameState()
	// Move enemy to same room as player
	state.Enemies["E1"].Location = "R12"
	action := MeleeAction{PlayerID: "P1"}
	
	result := ApplyCombat(state, action)
	
	enemy := result.Enemies["E1"]
	if enemy.HP != 2 { // 3 - 1 = 2
		t.Errorf("Expected enemy HP 2, got %d", enemy.HP)
	}
}

func TestMeleeAction_NoAmmoCost(t *testing.T) {
	state := newCombatTestGameState()
	state.Enemies["E1"].Location = "R12" // Same room
	action := MeleeAction{PlayerID: "P1"}
	
	result := ApplyCombat(state, action)
	
	player := result.Players["P1"]
	if player.Ammo != 5 { // Should remain 5
		t.Errorf("Expected player ammo 5, got %d", player.Ammo)
	}
}

func TestCombatAction_EnemyDeath(t *testing.T) {
	state := newCombatTestGameState()
	state.Enemies["E1"].HP = 1 // One hit away from death
	action := ShootAction{PlayerID: "P1"}
	
	result := ApplyCombat(state, action)
	
	// Enemy should be removed from game
	if _, exists := result.Enemies["E1"]; exists {
		t.Error("Dead enemy should be removed from game")
	}
}

// Helper for combat tests
func newCombatTestGameState() GameState {
	return GameState{
		Rooms: map[RoomID]*RoomState{
			"R12": {ID: "R12", Type: Predefined}, // Player location
			"R07": {ID: "R07", Type: AmmoCache},   // Adjacent to R12
		},
		Players: map[PlayerID]*PlayerState{
			"P1": {
				ID:       "P1",
				Location: "R12",
				Ammo:     5,
				HP:       10,
			},
		},
		Enemies: map[EnemyID]*Enemy{
			"E1": {
				ID:       "E1",
				Type:     InfiniteLoop,
				HP:       3,
				MaxHP:    3,
				Damage:   1,
				Location: "R07", // Adjacent to player
			},
		},
	}
}