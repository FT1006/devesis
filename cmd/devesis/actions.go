package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spaceship/devesis/pkg/core"
)

// consumeAction decrements the action counter and returns true if action is allowed
func (g *GameManager) consumeAction() bool {
	if g.state.ActionsLeft <= 0 {
		fmt.Println("✗ No actions remaining this turn!")
		return false
	}
	g.state.ActionsLeft--
	return true
}

func (g *GameManager) executeMove(args []string) error {
	if !g.consumeAction() {
		return nil
	}
	
	if len(args) == 0 {
		return fmt.Errorf("usage: move <roomID>")
	}
	
	targetRoom := strings.ToUpper(args[0])
	player := core.GetActivePlayer(g.state)
	if player == nil {
		return fmt.Errorf("no active player")
	}
	
	// Check basic movement validity first
	if !core.CanMove(g.state, player.Location, core.RoomID(targetRoom)) {
		fmt.Printf("✗ Cannot move to %s (not adjacent or blocked).\n", targetRoom)
		return nil
	}
	
	// Check if target room is already explored - no question needed
	targetRoomState := g.state.Rooms[core.RoomID(targetRoom)]
	if targetRoomState != nil && targetRoomState.Explored {
		// Skip question for explored rooms
		fmt.Printf("Moving to explored room %s (no question needed).\n", targetRoom)
	} else {
		// Warn about unexplored room and get confirmation
		fmt.Printf("⚠️  Warning: %s is unexplored! You'll need to answer a coding question.\n", targetRoom)
		fmt.Printf("Wrong answers cause bugs to spread and you lose all cards!\n")
		fmt.Print("Continue? (y/n): ")
		
		var confirm string
		fmt.Scanf("%s", &confirm)
		if confirm != "y" && confirm != "Y" && confirm != "yes" && confirm != "Yes" {
			fmt.Println("Movement cancelled.")
			return nil
		}
		
		// Get a coding question before allowing movement to unexplored rooms
		question, questionState := core.GetRandomQuestion(*g.state)
		if question.ID == -1 {
			// No questions available, allow movement without question
			fmt.Println("⚠ No questions available, movement allowed.")
		} else {
		// Show the question
		fmt.Printf("\n[CODING CHALLENGE] Answer correctly to move to %s:\n", targetRoom)
		fmt.Printf("%s\n\n", question.Text)
		for i, option := range question.Options {
			fmt.Printf("%d) %s\n", i+1, option)
		}
		
		fmt.Print("Answer (1-4): ")
		var choice int
		_, err := fmt.Scanf("%d", &choice)
		if err != nil || choice < 1 || choice > 4 {
			fmt.Println("✗ Invalid answer format. Movement failed.")
			return nil
		}
		
		// Check if answer is correct
		if core.CheckAnswer(question, choice-1) {
			fmt.Println("✓ Correct! You may proceed.")
		} else {
			fmt.Println("✗ Incorrect answer! Bugs spread everywhere...")
			// Update state to mark question as used
			g.state = &questionState
			// Apply severe penalties for wrong answer
			return g.applyWrongAnswerPenalties(core.RoomID(targetRoom))
		}
		// Update state to mark question as used for correct answers
		g.state = &questionState
		}
	}
	
	// Apply the movement action
	action := core.MoveAction{
		PlayerID: player.ID,
		To:       core.RoomID(targetRoom),
	}
	
	newState := core.Apply(*g.state, action)
	
	// Check if move succeeded
	newPlayer := core.GetActivePlayer(&newState)
	if newPlayer.Location == core.RoomID(targetRoom) {
		g.state = &newState
		fmt.Printf("✓ You move to %s.\n", targetRoom)
		
		// Mark room as explored when entering
		if room := g.state.Rooms[core.RoomID(targetRoom)]; room != nil && !room.Explored {
			room.Explored = true
			fmt.Printf("📍 You discover this is a %s.\n", g.getRoomTypeName(room.Type))
		}
	} else {
		fmt.Printf("✗ Movement to %s failed.\n", targetRoom)
	}
	
	return nil
}

func (g *GameManager) applyWrongAnswerPenalties(targetRoom core.RoomID) error {
	// 1. Allow movement first
	action := core.MoveAction{
		PlayerID: core.GetActivePlayer(g.state).ID,
		To:       targetRoom,
	}
	newState := core.Apply(*g.state, action)
	g.state = &newState
	fmt.Printf("✓ You move to %s.\n", targetRoom)
	
	// 2. Show what rooms will be affected
	roomsToInfect := []core.RoomID{targetRoom}
	adjacent := core.GetAdjacentRooms(targetRoom)
	roomsToInfect = append(roomsToInfect, adjacent...)
	
	fmt.Printf("💀 Bugs spread to %d rooms: ", len(roomsToInfect))
	for i, roomID := range roomsToInfect {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Print(roomID)
	}
	fmt.Println()
	
	// Count cards before penalty
	player := core.GetActivePlayer(g.state)
	cardCount := len(player.Hand)
	
	// 3. Apply all penalties using core function (bugs, spawns, card drop)
	core.ApplyWrongAnswerPenalties(g.state, targetRoom)
	
	// 4. Show feedback
	if cardCount > 0 {
		fmt.Printf("💸 You drop all %d cards from your hand!\n", cardCount)
	}
	
	fmt.Println("⏰ Your turn ends immediately due to the wrong answer.")
	
	return nil
}

func (g *GameManager) getRoomTypeName(roomType core.RoomType) string {
	switch roomType {
	case core.AmmoCache:
		return "Ammo Cache"
	case core.MedBay:
		return "Medical Bay"
	case core.CleanRoomType:
		return "Clean Room"
	case core.EnemySpawn:
		return "Enemy Spawn"
	case core.Empty:
		return "Air Circulation room"
	default:
		return "Unknown room type"
	}
}

func (g *GameManager) getEnemyName(enemyType core.EnemyType) string {
	switch enemyType {
	case core.InfiniteLoop:
		return "Infinite Loop"
	case core.StackOverflow:
		return "Stack Overflow"
	case core.Pythogoras:
		return "Pythogoras"
	default:
		return "Unknown Enemy"
	}
}

func (g *GameManager) executeSearch() error {
	if !g.consumeAction() {
		return nil
	}
	
	player := core.GetActivePlayer(g.state)
	if player == nil {
		return fmt.Errorf("no active player")
	}
	
	// Check if room was already searched
	currentRoom := g.state.Rooms[player.Location]
	if currentRoom != nil && currentRoom.Searched {
		fmt.Println("❌ This room has already been searched.")
		return nil
	}
	
	// Count cards before search
	cardsBefore := len(player.Hand)
	
	action := core.SearchAction{
		PlayerID: player.ID,
	}
	
	newState := core.Apply(*g.state, action)
	
	// Check if search succeeded by comparing room searched status
	room := newState.Rooms[player.Location]
	if room != nil && room.Searched {
		// Search succeeded
		newPlayer := newState.Players[player.ID]
		cardsAfter := len(newPlayer.Hand)
		
		g.state = &newState
		fmt.Println("🔍 You search the room thoroughly...")
		
		if cardsAfter > cardsBefore {
			// Found a card!
			newCardID := newPlayer.Hand[cardsAfter-1] // Last card added
			if card, exists := core.CardDB[newCardID]; exists {
				fmt.Printf("💎 Found special card: **%s** - %s\n", card.Name, card.Description)
				fmt.Printf("💾 **%s** added to your hand.\n", card.Name)
			} else {
				fmt.Println("💎 Found a special card!")
			}
		} else {
			// No card found
			fmt.Println("❌ Nothing useful found in this room.")
		}
	} else {
		fmt.Println("✗ Cannot search this room.")
	}
	
	return nil
}

func (g *GameManager) executeShoot() error {
	if !g.consumeAction() {
		return nil
	}
	
	player := core.GetActivePlayer(g.state)
	if player == nil {
		return fmt.Errorf("no active player")
	}
	
	// Check ammo before attempting
	if player.Ammo < core.ShootAmmoCost {
		fmt.Printf("✗ Not enough ammo! Need %d ammo, have %d.\n", core.ShootAmmoCost, player.Ammo)
		return nil
	}
	
	// Get adjacent rooms and count targets
	adjacentRooms := core.GetAdjacentRooms(player.Location)
	targets := make(map[core.RoomID][]core.EnemyID)
	totalTargets := 0
	
	for enemyID, enemy := range g.state.Enemies {
		for _, roomID := range adjacentRooms {
			if enemy.Location == roomID {
				targets[roomID] = append(targets[roomID], enemyID)
				totalTargets++
				break
			}
		}
	}
	
	if totalTargets == 0 {
		fmt.Println("✗ No enemies in adjacent rooms to shoot!")
		return nil
	}
	
	action := core.ShootAction{
		PlayerID: player.ID,
	}
	
	newState := core.ApplyCombat(*g.state, action)
	
	// Show combat results
	fmt.Printf("🔫 Shooting adjacent rooms! (-%d ammo)\n", core.ShootAmmoCost)
	for roomID, enemyIDs := range targets {
		for _, enemyID := range enemyIDs {
			oldEnemy := g.state.Enemies[enemyID]
			newEnemy, exists := newState.Enemies[enemyID]
			
			if !exists {
				fmt.Printf("   💀 %s in %s destroyed!\n", g.getEnemyName(oldEnemy.Type), roomID)
			} else if newEnemy.HP < oldEnemy.HP {
				fmt.Printf("   🎯 %s in %s: %d → %d HP\n", 
					g.getEnemyName(oldEnemy.Type), roomID, oldEnemy.HP, newEnemy.HP)
			}
		}
	}
	
	g.state = &newState
	return nil
}

func (g *GameManager) executeMelee() error {
	if !g.consumeAction() {
		return nil
	}
	
	player := core.GetActivePlayer(g.state)
	if player == nil {
		return fmt.Errorf("no active player")
	}
	
	// Count targets in current room
	targets := make([]core.EnemyID, 0)
	for enemyID, enemy := range g.state.Enemies {
		if enemy.Location == player.Location {
			targets = append(targets, enemyID)
		}
	}
	
	if len(targets) == 0 {
		fmt.Println("✗ No enemies in current room to attack!")
		return nil
	}
	
	action := core.MeleeAction{
		PlayerID: player.ID,
	}
	
	newState := core.ApplyCombat(*g.state, action)
	
	// Show combat results
	fmt.Println("⚔️ Melee attack!")
	for _, enemyID := range targets {
		oldEnemy := g.state.Enemies[enemyID]
		newEnemy, exists := newState.Enemies[enemyID]
		
		if !exists {
			fmt.Printf("   💀 %s destroyed!\n", g.getEnemyName(oldEnemy.Type))
		} else if newEnemy.HP < oldEnemy.HP {
			fmt.Printf("   🎯 %s: %d → %d HP\n", 
				g.getEnemyName(oldEnemy.Type), oldEnemy.HP, newEnemy.HP)
		}
	}
	
	g.state = &newState
	return nil
}

func (g *GameManager) executePlayCard(args []string) error {
	if !g.consumeAction() {
		return nil
	}
	
	if len(args) == 0 {
		return fmt.Errorf("usage: play <cardNumber> or play <cardID>")
	}
	
	player := core.GetActivePlayer(g.state)
	if player == nil {
		return fmt.Errorf("no active player")
	}
	
	if len(player.Hand) == 0 {
		fmt.Println("✗ Your hand is empty!")
		return nil
	}
	
	var cardID core.CardID
	
	// Try to parse as card number first (1-based index)
	if cardNum, err := strconv.Atoi(args[0]); err == nil {
		if cardNum < 1 || cardNum > len(player.Hand) {
			fmt.Printf("✗ Card number must be between 1 and %d!\n", len(player.Hand))
			return nil
		}
		// Convert to 0-based index and get the card ID
		cardID = player.Hand[cardNum-1]
	} else {
		// Treat as direct card ID
		cardID = core.CardID(args[0])
	}
	
	// Check if card is in hand
	found := false
	for _, handCardID := range player.Hand {
		if handCardID == cardID {
			found = true
			break
		}
	}
	
	if !found {
		fmt.Printf("✗ Card %s not in your hand!\n", cardID)
		return nil
	}
	
	// Get card details for feedback
	card, err := core.GetCard(cardID)
	if err != nil {
		fmt.Printf("✗ Unknown card: %s\n", cardID)
		return nil
	}
	
	// Execute play card action
	action := core.PlayCardAction{
		PlayerID: player.ID,
		CardID:   cardID,
	}
	
	newState := core.Apply(*g.state, action)
	g.state = &newState
	
	fmt.Printf("✓ Played %s\n", card.Name)
	return nil
}

func (g *GameManager) executeRoomAction() error {
	if !g.consumeAction() {
		return nil
	}
	
	player := core.GetActivePlayer(g.state)
	if player == nil {
		return fmt.Errorf("no active player")
	}
	
	// Get current room
	currentRoom := g.state.Rooms[player.Location]
	if currentRoom == nil {
		return fmt.Errorf("invalid player location")
	}
	
	// Check if room is corrupted
	if currentRoom.Corrupted {
		fmt.Println("✗ Cannot use room actions in corrupted rooms!")
		return nil
	}
	
	// Check if player already used room action this turn
	if player.SpecialUsed {  // Using existing field temporarily
		fmt.Println("✗ You've already used your room action this turn!")
		return nil
	}
	
	// Apply room-specific effect
	switch currentRoom.Type {
	case core.MedBay:
		return g.executeRoomActionMedBay(player, currentRoom)
	case core.AmmoCache:
		return g.executeRoomActionAmmoCache(player, currentRoom)
	case core.CleanRoomType:
		return g.executeRoomActionCleanRoom(player, currentRoom)
	default:
		fmt.Println("✗ No special room action available here.")
		return nil
	}
}

func (g *GameManager) executeRoomActionMedBay(player *core.PlayerState, room *core.RoomState) error {
	// Check if healing would be beneficial
	if player.HP >= player.MaxHP {
		fmt.Println("💡 You're already at full health - MedBay does nothing.")
		return nil
	}
	
	// Apply healing
	oldHP := player.HP
	newHP := player.HP + core.MedBayHealAmount
	if newHP > player.MaxHP {
		newHP = player.MaxHP
	}
	player.HP = newHP
	
	// Mark as used
	player.SpecialUsed = true
	
	fmt.Printf("🏥 MedBay healing! HP: %d → %d (+%d)\n", 
		oldHP, newHP, newHP-oldHP)
	fmt.Println("✓ Room action complete (1 action used)")
	
	return nil
}

func (g *GameManager) executeRoomActionAmmoCache(player *core.PlayerState, room *core.RoomState) error {
	// Check if ammo refill would be beneficial
	if player.Ammo >= player.MaxAmmo {
		fmt.Println("💡 You're already at full ammo - AmmoCache does nothing.")
		return nil
	}
	
	// Apply ammo refill
	oldAmmo := player.Ammo
	newAmmo := player.Ammo + core.AmmoCacheAmount
	if newAmmo > player.MaxAmmo {
		newAmmo = player.MaxAmmo
	}
	player.Ammo = newAmmo
	
	// Mark as used
	player.SpecialUsed = true
	
	fmt.Printf("🔫 AmmoCache refill! Ammo: %d → %d (+%d)\n", 
		oldAmmo, newAmmo, newAmmo-oldAmmo)
	fmt.Println("✓ Room action complete (1 action used)")
	
	return nil
}

func (g *GameManager) executeRoomActionCleanRoom(player *core.PlayerState, room *core.RoomState) error {
	// Get adjacent rooms
	adjacentRoomIDs := core.GetAdjacentRooms(player.Location)
	
	// Count rooms that have bugs to clean
	bugsCleaned := 0
	roomsCleaned := 0
	
	fmt.Println("🧹 CleanRoom decontamination activated!")
	
	for _, roomID := range adjacentRoomIDs {
		adjacentRoom := g.state.Rooms[roomID]
		if adjacentRoom != nil && adjacentRoom.BugMarkers > 0 {
			oldBugs := adjacentRoom.BugMarkers
			adjacentRoom.BugMarkers--
			bugsCleaned++
			roomsCleaned++
			
			// Update corruption status
			adjacentRoom.Corrupted = adjacentRoom.BugMarkers >= core.BugCorruptionThreshold
			
			fmt.Printf("   %s: %d → %d bugs (-1)\n", roomID, oldBugs, adjacentRoom.BugMarkers)
		} else if adjacentRoom != nil {
			fmt.Printf("   %s: 0 bugs (no effect)\n", roomID)
		}
	}
	
	if bugsCleaned == 0 {
		fmt.Println("💡 No bugs in adjacent rooms - CleanRoom does nothing.")
		return nil
	}
	
	// Mark as used
	player.SpecialUsed = true
	
	fmt.Printf("✓ Cleaned %d bugs from %d rooms (1 action used)\n", bugsCleaned, roomsCleaned)
	
	return nil
}

func (g *GameManager) executePass() error {
	player := core.GetActivePlayer(g.state)
	if player == nil {
		return fmt.Errorf("no active player")
	}
	
	// Pass ends the player phase immediately
	actionsSkipped := g.state.ActionsLeft
	g.state.ActionsLeft = 0
	fmt.Printf("You pass your turn. (%d actions skipped)\n", actionsSkipped)
	
	return nil
}