package main

import (
	"bufio"
	"fmt"
	"strings"
	"time"

	"github.com/mattn/go-runewidth"
	"github.com/spaceship/devesis/pkg/core"
)

type GameManager struct {
	state *core.GameState
}

func NewGameManager() *GameManager {
	return &GameManager{}
}

func (g *GameManager) Initialize() error {
	fmt.Println("Welcome to Devesis: Tutorial Hell!")
	fmt.Println("Escape the corrupted spaceship before time runs out!\n")
	
	// Load card database
	if err := core.LoadCards("./data"); err != nil {
		return fmt.Errorf("failed to load cards: %w", err)
	}
	
	// Get player class selection
	playerClass, err := g.selectPlayerClass()
	if err != nil {
		return err
	}
	
	// Create initial game state using reducer
	emptyState := core.GameState{}
	initialAction := core.InitializeGameAction{
		Seed:        time.Now().UnixNano(),
		PlayerClass: playerClass,
	}
	
	newState := core.Apply(emptyState, initialAction)
	g.state = &newState
	
	fmt.Printf("You are a %s developer. Good luck!\n", g.getClassDisplayName(playerClass))
	fmt.Println("Type '?' for help\n")
	return nil
}

func (g *GameManager) selectPlayerClass() (core.DevClass, error) {
	classes := core.GetAvailableClasses()
	
	fmt.Println("Choose your developer class:")
	for _, class := range classes {
		fmt.Printf("%d. %-9s (HP: %d, Ammo: %d) - %s\n", 
			class.ID, class.DisplayName, class.HP, class.MaxAmmo, class.Description)
	}
	fmt.Print("Enter choice (1-4): ")
	
	var choice int
	_, err := fmt.Scanf("%d", &choice)
	if err != nil {
		return core.Frontend, fmt.Errorf("invalid input: %v", err)
	}
	
	if selectedClass, valid := core.ValidateClassChoice(choice); valid {
		return selectedClass, nil
	} else {
		fmt.Printf("Invalid choice %d, defaulting to Frontend\n", choice)
		return core.Frontend, nil
	}
}

func (g *GameManager) getClassDisplayName(class core.DevClass) string {
	classes := core.GetAvailableClasses()
	for _, c := range classes {
		if c.Class == class {
			return c.DisplayName
		}
	}
	return "Unknown"
}

func (g *GameManager) IsGameOver() bool {
	return core.IsGameOver(g.state)
}

func (g *GameManager) DisplayStatus() {
	player := core.GetActivePlayer(g.state)
	if player == nil {
		return
	}
	
	// Display hand first
	g.displayHandStatus(player)
	
	// Then display bordered status panel
	g.displayStatusPanel(player)
	
	// Finally display command help
	fmt.Printf("\n[ACTIONS] move(mv) play(c) search(s) shoot(f) melee(ml) room(ra) pass(p)\n")
	fmt.Printf("[INFO] hand(h) map(mp) status(st) help(?) quit/exit(q)\n")
}

func (g *GameManager) displayHandStatus(player *core.PlayerState) {
	fmt.Printf("\nHand (%d cards):\n", len(player.Hand))
	if len(player.Hand) == 0 {
		fmt.Println("  (empty)")
	} else {
		for i, cardID := range player.Hand {
			if card, exists := core.CardDB[cardID]; exists {
				fmt.Printf("  %d. %s - %s\n", i+1, card.Name, card.Description)
			} else {
				fmt.Printf("  %d. %s (unknown card)\n", i+1, cardID)
			}
		}
	}
}

// helper – visible width in cells, not bytes
func cells(s string) int { return runewidth.StringWidth(s) }

func (g *GameManager) displayStatusPanel(player *core.PlayerState) {
	const minWidth = 72  // never shrink below this
	
	// Get room info
	room := g.state.Rooms[player.Location]
	roomType := ""
	searchStatus := ""
	
	if room != nil {
		if room.Searched {
			searchStatus = "searched"
		} else {
			searchStatus = "not searched"  
		}
		roomType = g.getRoomTypeDisplayName(room)
	}
	
	// Count enemies in room
	loopCount := 0
	overflowCount := 0 
	pythogorasCount := 0
	for _, enemy := range g.state.Enemies {
		if enemy.Location == player.Location {
			switch enemy.Type {
			case core.InfiniteLoop:
				loopCount++
			case core.StackOverflow:
				overflowCount++
			case core.Pythogoras:
				pythogorasCount++
			}
		}
	}
	
	// Calculate rounds left
	roundsLeft := core.MaxRounds - g.state.Round
	if roundsLeft < 0 {
		roundsLeft = 0
	}
	
	// Corruption status
	corruptedStatus := "✘"
	if room != nil && room.Corrupted {
		corruptedStatus = "✓"
	}
	
	// ➊ create the raw content **without** borders
	lines := make([]string, 0, 4)
	
	lines = append(lines,
		fmt.Sprintf("HP   %2d / %2d     Ammo %2d / %2d        Cards  Hand:%d  Deck:%d  Discard:%d",
			player.HP, player.MaxHP, player.Ammo, player.MaxAmmo, 
			len(player.Hand), len(player.Deck), len(player.Discard)),
	)
	lines = append(lines,
		fmt.Sprintf("Turn   Actions %d / 2", g.state.ActionsLeft),
	)
	lines = append(lines,
		fmt.Sprintf("Room   Bugs:%d   Loop:%d   Overflow:%d   Pythogoras:%d   Corrupted: %s",
			room.BugMarkers, loopCount, overflowCount, pythogorasCount, corruptedStatus),
	)
	lines = append(lines,
		fmt.Sprintf("Game   Round: %d      Rounds left: %d", 
			g.state.Round, roundsLeft),
	)
	
	// ➊ work out how wide the panel really needs to be
	width := minWidth
	for _, l := range lines {
		if w := cells(l) + 2; w > width { // +2 for the side borders
			width = w
		}
	}
	
	// ➌ print top border
	className := g.getClassDisplayName(player.Class)
	header := fmt.Sprintf(" P1 %s ── Room %s (%s, %s) ",
		className, player.Location, roomType, searchStatus)
	pad := width - cells(header) - 2          // 2 for corner chars
	fmt.Printf("\n┌%s%s┐\n", header, strings.Repeat("─", pad))
	
	// ➍ print each content line
	for _, l := range lines {
		fmt.Printf("│%-*s│\n", width-2, l)     // %-*s pads to full width
	}
	
	// ➎ bottom border
	fmt.Printf("└%s┘\n", strings.Repeat("─", width-2))
}

func (g *GameManager) getRoomDisplayName(roomID core.RoomID) string {
	// Simple mapping for room names
	roomNames := map[string]string{
		"R01": "Key", "R02": "Store", "R03": "Comp",
		"R04": "Crew", "R05": "Lab", "R06": "Sys", 
		"R07": "Air", "R08": "Power", "R09": "Maint",
		"R10": "Cache", "R11": "Cache", "R12": "Start",
		"R13": "Data", "R14": "Log", "R15": "Engine",
		"R16": "Gen", "R17": "Engine", "R18": "Engine", 
		"R19": "Escape", "R20": "Escape",
	}
	
	if name, exists := roomNames[string(roomID)]; exists {
		return name
	}
	return string(roomID)
}

func (g *GameManager) getRoomTypeDisplayName(room *core.RoomState) string {
	switch room.Type {
	case core.AmmoCache:
		return "ammo cache"
	case core.MedBay:
		return "med bay"
	case core.CleanRoomType:
		return "clean room"
	case core.EnemySpawn:
		return "enemy spawn"
	case core.Empty:
		return "air circulation"
	default:
		return "unknown"
	}
}

func (g *GameManager) DisplayGameOver() {
	fmt.Println("\n========== GAME OVER ==========")
	if g.state.Round >= core.MaxRounds {
		fmt.Println("💀 TIME'S UP! The corruption consumed the ship...")
	} else {
		fmt.Println("💀 DEFEAT! All developers were lost to the corruption...")
	}
}

func (g *GameManager) ExecuteCommand(command string, args []string) error {
	switch command {
	// Turn-economy actions
	case "move", "mv":
		return g.executeMove(args)
	case "play", "c":
		return g.executePlayCard(args)
	case "search", "s":
		return g.executeSearch()
	case "shoot", "f":
		return g.executeShoot()
	case "melee", "ml":
		return g.executeMelee()
	case "room", "ra":
		return g.executeRoomAction()
	case "pass", "p":
		return g.executePass()
		
	// Information commands
	case "hand", "h":
		return g.showHand()
	case "map", "mp":
		return g.showMap()
	case "status", "st":
		g.DisplayStatus()
		return nil
	case "help", "?":
		return g.showHelp()
	case "quit", "q":
		fmt.Println("Thanks for playing!")
		return fmt.Errorf("quit")
		
	default:
		return fmt.Errorf("unknown command. Type '?' for help")
	}
}

// Phase execution methods for 4-phase round structure

func (g *GameManager) ExecuteDrawPhase() {
	fmt.Printf("\n=== ROUND %d: DRAW PHASE ===\n", g.state.Round)
	core.DrawPhase(g.state)
	
	player := core.GetActivePlayer(g.state)
	if player != nil {
		fmt.Printf("📋 Drew cards. Hand: %d cards, Deck: %d cards\n", len(player.Hand), len(player.Deck))
	}
}

func (g *GameManager) ExecutePlayerPhase(reader *bufio.Reader) error {
	fmt.Printf("\n=== PLAYER PHASE ===\n")
	fmt.Printf("Actions remaining: %d\n", g.state.ActionsLeft)
	
	// Player action loop - continue until actions exhausted or pass
	for g.state.ActionsLeft > 0 {
		// Display current status
		g.DisplayStatus()
		
		// Show command prompt
		fmt.Printf("(%d actions left) > ", g.state.ActionsLeft)
		
		// Read command
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("error reading input: %w", err)
		}
		
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}
		
		// Parse and execute command
		args := strings.Fields(input)
		if len(args) == 0 {
			continue
		}
		
		command := strings.ToLower(args[0])
		commandArgs := args[1:]
		
		// Execute command
		if err := g.ExecuteCommand(command, commandArgs); err != nil {
			if err.Error() == "quit" {
				return err
			}
			fmt.Printf("Error: %v\n", err)
		}
		
		// Check if player passed (ActionsLeft set to 0)
		if g.state.ActionsLeft == 0 {
			break
		}
	}
	
	fmt.Println("Player phase complete.")
	return nil
}

func (g *GameManager) ExecuteEventPhase() {
	fmt.Printf("\n=== EVENT PHASE ===\n")
	core.EventPhase(g.state)
	fmt.Printf("Time remaining: %d rounds\n", g.state.Time)
}

func (g *GameManager) ExecuteRoundMaintenance() {
	fmt.Printf("\n=== ROUND MAINTENANCE ===\n")
	core.EndRoundMaintenance(g.state)
	fmt.Printf("Round %d complete. Starting round %d...\n", g.state.Round-1, g.state.Round)
}

func (g *GameManager) CheckEndConditions() (ended bool, win bool) {
	return core.CheckEndSolo(g.state)
}

func (g *GameManager) DisplayGameResult(win bool) {
	fmt.Println("\n========== GAME COMPLETE ==========")
	if win {
		fmt.Println("🎉 VICTORY! You escaped Tutorial Hell!")
	} else {
		fmt.Println("💀 DEFEAT! You were consumed by the corruption...")
	}
}