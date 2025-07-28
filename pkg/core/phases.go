package core

import (
	"fmt"
	"math/rand"
)

// DrawPhase refills active player to 5 cards from their deck
func DrawPhase(state *GameState) {
	player := GetActivePlayer(state)
	if player == nil {
		return
	}
	
	targetHandSize := 5
	cardsToDraw := targetHandSize - len(player.Hand)
	
	if cardsToDraw > 0 {
		rng := rand.New(rand.NewSource(state.RandSeed + int64(state.Round)*100))
		drawCards(&player.Hand, &player.Deck, &player.Discard, cardsToDraw, rng)
	}
	
	// Set actions for player phase
	state.ActionsLeft = 2
	state.Phase = "player"
}

// EventPhase executes the 6-step event sequence from ruleset
func EventPhase(state *GameState) {
	state.Phase = "event"
	
	// Step 1: Time Marker -1
	state.Time--
	
	// Step 2: Malware attacks co-located developers
	malwareAttackPhase(state)
	
	// Step 3: System crashes damage malware in OutOfRam rooms
	systemCrashPhase(state)
	
	// Step 4: Draw Event Card (apply effect from fixed deck)
	drawEventCardPhase(state)
	
	// Step 5: Enemy Development (draw tokens based on round)
	enemyDevelopmentPhase(state)
	
	// Step 6: Check End Triggers (handled by caller)
}

// EndRoundMaintenance resets per-round flags and advances round
func EndRoundMaintenance(state *GameState) {
	// Reset per-round player flags
	for _, player := range state.Players {
		player.HasActed = false
		player.SpecialUsed = false
	}
	
	// Advance round
	state.Round++
	
	// Reset phase to player for next round
	state.Phase = "player"
	state.ActionsLeft = 0 // Will be set by DrawPhase
}

// CheckEndSolo checks provisional win/loss conditions for solo play
func CheckEndSolo(state *GameState) (ended bool, win bool) {
	player := GetActivePlayer(state)
	if player == nil {
		return true, false // No player = loss
	}
	
	// Loss conditions
	if player.HP == 0 {
		return true, false // Death = loss
	}
	if state.Time <= 0 {
		return true, false // Time up = loss
	}
	
	// Win condition: player successfully used engine card at escape room
	if player.EngineUsed {
		return true, true
	}
	
	return false, false
}

// Helper functions for event phase steps

func malwareAttackPhase(state *GameState) {
	// For each enemy, attack any co-located players
	for _, enemy := range state.Enemies {
		for _, player := range state.Players {
			if player.Location == enemy.Location && player.HP > 0 {
				// Apply enemy damage
				damage := enemy.Damage
				if player.HP <= damage {
					player.HP = 0
				} else {
					player.HP -= damage
				}
			}
		}
	}
}

func systemCrashPhase(state *GameState) {
	// Damage all malware in OutOfRam rooms
	for enemyID, enemy := range state.Enemies {
		room := state.Rooms[enemy.Location]
		if room != nil && room.OutOfRam {
			enemy.HP--
			if enemy.HP <= 0 {
				delete(state.Enemies, enemyID)
			}
		}
	}
}

func drawEventCardPhase(state *GameState) {
	// TODO: Implement event cards from YAML
	// For now, just advance the event index
	if len(state.Events) > 0 {
		state.EventIndex = (state.EventIndex + 1) % uint8(len(state.Events))
	}
}

func enemyDevelopmentPhase(state *GameState) {
	if state.SpawnBag == nil {
		return
	}
	
	// Draw count = (round + 1) / 2
	drawCount := (state.Round + 1) / 2
	
	rng := rand.New(rand.NewSource(state.RandSeed + int64(state.Round)*1000))
	
	for i := 0; i < drawCount; i++ {
		if len(state.SpawnBag.Tokens) == 0 {
			break
		}
		
		// Draw random token from bag
		tokenIndex := rng.Intn(len(state.SpawnBag.Tokens))
		token := state.SpawnBag.Tokens[tokenIndex]
		
		// Remove token from bag
		state.SpawnBag.Tokens = append(
			state.SpawnBag.Tokens[:tokenIndex], 
			state.SpawnBag.Tokens[tokenIndex+1:]...)
		
		// Spawn enemy if not NoSpawn
		if token != InfiniteLoop && token != StackOverflow && token != Pythogoras {
			continue // Skip invalid tokens
		}
		
		// Spawn the enemy
		spawnEnemy(state, token, rng)
		
		// Add stronger token back to bag
		addStrongerToken(state.SpawnBag, token)
	}
}

func spawnEnemy(state *GameState, enemyType EnemyType, rng *rand.Rand) {
	// Find a random room to spawn in
	roomIDs := make([]RoomID, 0, len(state.Rooms))
	for roomID := range state.Rooms {
		roomIDs = append(roomIDs, roomID)
	}
	
	if len(roomIDs) == 0 {
		return
	}
	
	spawnRoom := roomIDs[rng.Intn(len(roomIDs))]
	stats := ENEMY_STATS[enemyType]
	
	// Generate unique enemy ID
	enemyID := EnemyID(fmt.Sprintf("%s_%d_%d", getEnemyTypeName(enemyType), state.Round, rng.Int31()))
	
	state.Enemies[enemyID] = &Enemy{
		ID:       enemyID,
		Type:     enemyType,
		HP:       stats.HP,
		MaxHP:    stats.HP,
		Damage:   stats.Damage,
		Location: spawnRoom,
	}
}

func addStrongerToken(bag *SpawnBag, token EnemyType) {
	var stronger EnemyType
	switch token {
	case InfiniteLoop:
		stronger = StackOverflow
	case StackOverflow:
		stronger = Pythogoras
	case Pythogoras:
		stronger = Pythogoras // Pythogoras is max level
	default:
		return
	}
	
	bag.Tokens = append(bag.Tokens, stronger)
}

func getEnemyTypeName(enemyType EnemyType) string {
	switch enemyType {
	case InfiniteLoop:
		return "LOOP"
	case StackOverflow:
		return "OVERFLOW"
	case Pythogoras:
		return "PYTHO"
	default:
		return "UNKNOWN"
	}
}

