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
	if len(args) == 0 {
		return fmt.Errorf("usage: move <roomID>")
	}
	
	targetRoom := strings.ToUpper(args[0])
	player := core.GetActivePlayer(g.state)
	if player == nil {
		return fmt.Errorf("no active player")
	}
	
	// Check basic movement validity first - don't consume action if invalid
	if !core.CanMove(g.state, player.Location, core.RoomID(targetRoom)) {
		fmt.Printf("✗ Cannot move to %s (not adjacent).\n", targetRoom)
		return nil
	}
	
	if !g.consumeAction() {
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
			// Reward: Give a special card for correct answer using reducer
			rewardAction := core.GiveSpecialCardAction{
				PlayerID: core.GetActivePlayer(g.state).ID,
			}
			newState := core.ApplyWithoutLog(*g.state, rewardAction)
			
			// Check if a card was actually added
			oldPlayer := core.GetActivePlayer(g.state)
			newPlayer := core.GetActivePlayer(&newState)
			if len(newPlayer.Hand) > len(oldPlayer.Hand) {
				// Card was added, show reward message
				addedCard := newPlayer.Hand[len(newPlayer.Hand)-1]
				if card, exists := core.CardDB[addedCard]; exists {
					fmt.Printf("🎁 Reward: **%s** - %s\n", card.Name, card.Description)
				} else {
					fmt.Printf("🎁 Reward: Special card added to hand!\n")
				}
			}
			g.state = &newState
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
	
	// Apply the movement action with logging
	action := core.MoveAction{
		PlayerID: player.ID,
		To:       core.RoomID(targetRoom),
	}
	
	g.ResolveWithLogging(action)
	
	// Check if move succeeded (effects already shown via ResolveWithLogging)
	newPlayer := core.GetActivePlayer(g.state)
	if newPlayer.Location == core.RoomID(targetRoom) {
		// Show room type discovery if it's newly explored
		if room := g.state.Rooms[core.RoomID(targetRoom)]; room != nil && room.Explored {
			fmt.Printf("📍 You discover this is a %s.\n", g.getRoomTypeName(room))
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
	newState := core.ApplyWithoutLog(*g.state, action)
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

func (g *GameManager) getRoomTypeName(room *core.RoomState) string {
	// Handle predefined rooms first
	if room.Type == core.Predefined {
		switch room.ID {
		case "R01":
			return "key room"
		case "R12":
			return "start room"
		case "R15", "R17", "R18":
			return "engine room"
		case "R19", "R20":
			return "escape room"
		default:
			return "special room"
		}
	}
	
	// Handle regular room types
	switch room.Type {
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
	
	// Check if room was already searched - don't consume action if invalid
	currentRoom := g.state.Rooms[player.Location]
	if currentRoom != nil && currentRoom.Searched {
		fmt.Println("❌ This room has already been searched.")
		return nil
	}
	
	if !g.consumeAction() {
		return nil
	}
	
	// Execute search action with logging
	action := core.SearchAction{
		PlayerID: player.ID,
	}
	
	g.ResolveWithLogging(action)
	
	return nil
}

func (g *GameManager) executeShoot() error {
	player := core.GetActivePlayer(g.state)
	if player == nil {
		return fmt.Errorf("no active player")
	}
	
	// Check ammo before attempting - don't consume action if insufficient
	if player.Ammo < core.ShootAmmoCost {
		fmt.Printf("✗ Not enough ammo! Need %d ammo, have %d.\n", core.ShootAmmoCost, player.Ammo)
		return nil
	}
	
	// Get adjacent rooms and count targets - don't consume action if no targets
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
	
	if !g.consumeAction() {
		return nil
	}
	
	action := core.ShootAction{
		PlayerID: player.ID,
	}
	
	g.ResolveWithLogging(action)
	return nil
}

func (g *GameManager) executeMelee() error {
	player := core.GetActivePlayer(g.state)
	if player == nil {
		return fmt.Errorf("no active player")
	}
	
	// Count targets in current room - don't consume action if no targets
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
	
	if !g.consumeAction() {
		return nil
	}
	
	action := core.MeleeAction{
		PlayerID: player.ID,
	}
	
	g.ResolveWithLogging(action)
	return nil
}

func (g *GameManager) executePlayCard(args []string) error {
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
	
	if !g.consumeAction() {
		return nil
	}
	
	// Execute play card action with logging
	action := core.PlayCardAction{
		PlayerID: player.ID,
		CardID:   cardID,
	}
	
	fmt.Printf("✓ Playing %s\n", card.Name)
	g.ResolveWithLogging(action)
	return nil
}

func (g *GameManager) executeRoomAction() error {
	player := core.GetActivePlayer(g.state)
	if player == nil {
		return fmt.Errorf("no active player")
	}
	
	// Basic validation before consuming action
	if player.SpecialUsed {
		fmt.Println("✗ Room action already used this turn!")
		return nil
	}
	
	// Check if current room has a special ability
	room := g.state.Rooms[player.Location]
	if room == nil {
		fmt.Println("✗ No room found at current location!")
		return nil
	}
	
	// Check if room type supports actions
	if room.Type != core.MedBay && room.Type != core.AmmoCache && room.Type != core.CleanRoomType {
		fmt.Println("✗ No special room action available here")
		return nil
	}
	
	if !g.consumeAction() {
		return nil
	}
	
	// Execute room action through core reducer with logging
	action := core.RoomAction{
		PlayerID: player.ID,
	}
	
	g.ResolveWithLogging(action)
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