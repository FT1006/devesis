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

// EventPhase executes the 6-step event sequence from ruleset with logging
func EventPhase(state *GameState, log *EffectLog) {
	state.Phase = "event"
	
	log.Add("=== EVENT PHASE ===")
	
	// Step 1: Time Marker -1
	oldTime := state.Time
	state.Time--
	log.Add("⏰ Step 1: Time passes - Round %d → %d", oldTime, state.Time)
	
	// Step 2: Malware attacks co-located developers
	log.Add("👹 Step 2: Malware attacks...")
	malwareAttackPhase(state, log)
	
	// Step 3: System crashes damage malware in OutOfRam rooms
	log.Add("💥 Step 3: System crashes...")
	systemCrashPhase(state, log)
	
	// Step 4: Draw Event Card (apply effect from fixed deck)
	log.Add("🃏 Step 4: Event card...")
	drawEventCardPhase(state, log)
	
	// Step 5: Enemy Development (draw tokens based on round)
	log.Add("🧬 Step 5: Enemy development...")
	enemyDevelopmentPhase(state, log)
	
	// Step 5.5: Corrupted Room Spawns (enemies spawn in all corrupted rooms)
	log.Add("👹 Corruption spawns...")
	corruptedRoomSpawnPhase(state, log)
	
	// Step 6: Check End Triggers (handled by caller)
	log.Add("🎯 Step 6: End condition checks...")
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

func malwareAttackPhase(state *GameState, log *EffectLog) {
	// For each enemy, attack any co-located players
	attacksOccurred := false
	for _, enemy := range state.Enemies {
		for _, player := range state.Players {
			if player.Location == enemy.Location && player.HP > 0 {
				// Apply enemy damage
				oldHP := player.HP
				damage := enemy.Damage
				if player.HP <= damage {
					player.HP = 0
				} else {
					player.HP -= damage
				}
				log.Add("💔 %s attacks %s in %s! HP: %d → %d", getEnemyDisplayName(enemy.Type), player.ID, player.Location, oldHP, player.HP)
				attacksOccurred = true
			}
		}
	}
	if !attacksOccurred {
		log.Add("✅ No co-located enemies - players are safe")
	}
}

func systemCrashPhase(state *GameState, log *EffectLog) {
	// Damage all malware in OutOfRam rooms
	crashesOccurred := false
	for enemyID, enemy := range state.Enemies {
		room := state.Rooms[enemy.Location]
		if room != nil && room.OutOfRam {
			oldHP := enemy.HP
			damage := uint8(2)
			
			// Check if damage would kill the enemy (prevent uint8 underflow)
			if enemy.HP <= damage {
				log.Add("💥 %s destroyed by system crash in %s!", getEnemyDisplayName(enemy.Type), enemy.Location)
				delete(state.Enemies, enemyID)
			} else {
				enemy.HP -= damage
				log.Add("💥 %s damaged by system crash in %s! HP: %d → %d", getEnemyDisplayName(enemy.Type), enemy.Location, oldHP, enemy.HP)
			}
			crashesOccurred = true
		}
	}
	if !crashesOccurred {
		log.Add("✅ No enemies in OutOfRam rooms")
	}
}

func drawEventCardPhase(state *GameState, log *EffectLog) {
	if len(state.Events) == 0 {
		log.Add("🃏 No event cards available")
		return
	}
	
	// Draw the current event card
	eventCard := state.Events[state.EventIndex]
	log.Add("🃏 Event: %s (%s) - %s", eventCard.Name, eventCard.ID, eventCard.Description)
	
	// Apply the event card effects
	for _, effect := range eventCard.Effects {
		log.Add("🔧 Applying effect: %s (scope: %s, n: %d)", 
			getEffectOpName(effect.Op), getScopeName(effect.Scope), effect.N)
		applyEventEffect(state, effect, log)
	}
	
	// Advance to next event card
	state.EventIndex = (state.EventIndex + 1) % uint8(len(state.Events))
}

func enemyDevelopmentPhase(state *GameState, log *EffectLog) {
	if state.SpawnBag == nil {
		log.Add("🧬 No spawn bag available")
		return
	}
	
	// Draw count = (round + 1) / 2
	drawCount := (state.Round + 1) / 2
	log.Add("🧬 Drawing %d enemy tokens (round %d)", drawCount, state.Round)
	
	rng := rand.New(rand.NewSource(state.RandSeed + int64(state.Round)*1000))
	
	spawned := 0
	for i := 0; i < drawCount; i++ {
		if len(state.SpawnBag.Tokens) == 0 {
			log.Add("🧬 Spawn bag empty - no more enemies")
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
		spawnLocation := spawnEnemy(state, token, rng, log)
		if spawnLocation != "" {
			spawned++
		}
		
		// Add stronger token back to bag
		addStrongerToken(state.SpawnBag, token, log)
	}
	
	if spawned == 0 {
		log.Add("🧬 No enemies spawned this round")
	}
}

func spawnEnemy(state *GameState, enemyType EnemyType, rng *rand.Rand, log *EffectLog) RoomID {
	// Find a random room to spawn in
	roomIDs := make([]RoomID, 0, len(state.Rooms))
	for roomID := range state.Rooms {
		roomIDs = append(roomIDs, roomID)
	}
	
	if len(roomIDs) == 0 {
		return ""
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
	
	log.Add("👹 %s spawned in %s", getEnemyDisplayName(enemyType), spawnRoom)
	return spawnRoom
}

// applyEventEffect applies a single effect from an event card using the centralized handler
func applyEventEffect(state *GameState, effect Effect, log *EffectLog) {
	// Event cards don't have a specific player, so use empty PlayerID
	// The effect system will handle this appropriately
	playerID := PlayerID("")
	
	// Error logging is handled centrally in ApplyEffect
	ApplyEffect(state, effect, playerID, log)
}

func addStrongerToken(bag *SpawnBag, token EnemyType, log *EffectLog) {
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
	log.Add("🧬 %s token added to spawn bag", getEnemyDisplayName(stronger))
}

func corruptedRoomSpawnPhase(state *GameState, log *EffectLog) {
	spawnCount := 0
	
	for _, room := range state.Rooms {
		if room.Corrupted && !room.OutOfRam {
			// Spawn Infinite Loop (weakest enemy) in each corrupted room
			enemyID := EnemyID(fmt.Sprintf("CORRUPT_%s_%d", room.ID, state.Round))
			stats := ENEMY_STATS[InfiniteLoop]
			
			enemy := &Enemy{
				ID:       enemyID,
				Type:     InfiniteLoop,
				HP:       stats.HP,
				MaxHP:    stats.HP,
				Damage:   stats.Damage,
				Location: room.ID,
			}
			state.Enemies[enemyID] = enemy
			
			log.Add("👹 Infinite Loop spawned in corrupted %s", room.ID)
			spawnCount++
		}
	}
	
	if spawnCount == 0 {
		log.Add("👹 No corrupted rooms - no corruption spawns")
	} else {
		log.Add("👹 %d enemies spawned from corruption", spawnCount)
	}
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


