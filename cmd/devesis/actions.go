package main

import (
	"fmt"
	"strings"

	"github.com/spaceship/devesis/pkg/core"
)

func (g *GameManager) executeMove(args []string) error {
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
	player := core.GetActivePlayer(g.state)
	if player == nil {
		return fmt.Errorf("no active player")
	}
	
	// Check ammo before attempting
	if player.Ammo < core.ShootAmmoCost {
		fmt.Printf("✗ Not enough ammo! Need %d ammo, have %d.\\n", core.ShootAmmoCost, player.Ammo)
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
	fmt.Printf("🔫 Shooting adjacent rooms! (-%d ammo)\\n", core.ShootAmmoCost)
	for roomID, enemyIDs := range targets {
		for _, enemyID := range enemyIDs {
			oldEnemy := g.state.Enemies[enemyID]
			newEnemy, exists := newState.Enemies[enemyID]
			
			if !exists {
				fmt.Printf("   💀 %s in %s destroyed!\\n", g.getEnemyName(oldEnemy.Type), roomID)
			} else if newEnemy.HP < oldEnemy.HP {
				fmt.Printf("   🎯 %s in %s: %d → %d HP\\n", 
					g.getEnemyName(oldEnemy.Type), roomID, oldEnemy.HP, newEnemy.HP)
			}
		}
	}
	
	g.state = &newState
	return nil
}

func (g *GameManager) executeMelee() error {
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
			fmt.Printf("   💀 %s destroyed!\\n", g.getEnemyName(oldEnemy.Type))
		} else if newEnemy.HP < oldEnemy.HP {
			fmt.Printf("   🎯 %s: %d → %d HP\\n", 
				g.getEnemyName(oldEnemy.Type), oldEnemy.HP, newEnemy.HP)
		}
	}
	
	g.state = &newState
	return nil
}

func (g *GameManager) executePlayCard(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: play <cardID>")
	}
	
	player := core.GetActivePlayer(g.state)
	if player == nil {
		return fmt.Errorf("no active player")
	}
	
	if len(player.Hand) == 0 {
		fmt.Println("✗ Your hand is empty!")
		return nil
	}
	
	cardID := core.CardID(args[0])
	
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
	// TODO: Implement room actions
	return fmt.Errorf("room actions not yet implemented")
}

func (g *GameManager) executePass() error {
	player := core.GetActivePlayer(g.state)
	if player == nil {
		return fmt.Errorf("no active player")
	}
	
	action := core.PassAction{
		PlayerID: player.ID,
	}
	
	newState := core.Apply(*g.state, action)
	g.state = &newState
	fmt.Println("You pass your turn.")
	
	return nil
}