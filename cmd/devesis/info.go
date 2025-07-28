package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"io/ioutil"

	"github.com/spaceship/devesis/pkg/core"
)

func (g *GameManager) showHand() error {
	player := core.GetActivePlayer(g.state)
	if player == nil {
		return fmt.Errorf("no active player")
	}
	
	fmt.Printf("Hand (%d cards):\n", len(player.Hand))
	if len(player.Hand) == 0 {
		fmt.Println("  (empty)")
	} else {
		for i, cardID := range player.Hand {
			if card, exists := core.CardDB[cardID]; exists {
				fmt.Printf("  %d. %s (%s) - %s\n", i+1, card.Name, cardID, card.Description)
			} else {
				fmt.Printf("  %d. %s (unknown card)\n", i+1, cardID)
			}
		}
	}
	
	return nil
}

func (g *GameManager) showMap() error {
	fmt.Println(g.renderMap())
	return nil
}

func (g *GameManager) renderMap() string {
	var result strings.Builder
	overflowRooms := make(map[core.RoomID][]string)
	
	// Header row  
	result.WriteString("   A            B            C            D            E            F            G\n")
	
	// Render each row
	for row := 0; row < core.GridRows; row++ {
		line1, line2, line3 := g.renderGridRow(row, overflowRooms)
		result.WriteString(line1 + "\n")
		result.WriteString(line2 + "\n") 
		result.WriteString(line3 + "\n")
		
		if row < core.GridRows-1 {
			result.WriteString("\n") // Empty line between rows
		}
	}
	
	// Add overflow details if any
	if len(overflowRooms) > 0 {
		result.WriteString("\nOverflow Details:\n")
		for roomID, objects := range overflowRooms {
			result.WriteString(fmt.Sprintf("%s: %s\n", roomID, strings.Join(objects, ",")))
		}
	}
	
	// Add mini-guide with examples
	result.WriteString("\n")
	result.WriteString("Examples: [R12,+,0] = Room R12, not searched, 0 bugs | [R07,-,2] = Room R07, searched, 2 bugs\n")
	result.WriteString("Content:  P1 = you, P2-P4 = other players | IL = Infinite Loop, SO = Stack Overflow, PY = Pythogoras\n")
	result.WriteString("\n")
	result.WriteString("[PREDEFINED] KEY:R01 STR:R12 EN1:R15 EN2:R18 EN3:R17 ESC:R19,R20\n")
	result.WriteString("[ROOM TYPES] AMO×3 MED×3 CLN×2 AIR×3 SPN×4\n")
	
	return result.String()
}

// formatGridCell ensures exactly 12 characters for consistent grid alignment  
func (g *GameManager) formatGridCell(content string) string {
	const cellWidth = 12
	const maxContentWidth = cellWidth - 4 // [ content ]
	
	if len(content) > maxContentWidth {
		content = content[:maxContentWidth] // Truncate if too long
	}
	
	return fmt.Sprintf("[ %-*s ] ", maxContentWidth, content)
}

func (g *GameManager) renderGridRow(row int, overflowRooms map[core.RoomID][]string) (string, string, string) {
	var line1, line2, line3 strings.Builder
	
	line1.WriteString(fmt.Sprintf("%d ", row+1))
	line2.WriteString("  ")
	line3.WriteString("  ")
	
	for col := 0; col < core.GridCols; col++ {
		roomID := g.findRoomAt(row, col)
		
		if roomID == "" {
			// Empty cell - exactly 9 characters each line
			line1.WriteString(g.formatGridCell(""))
			line2.WriteString(g.formatGridCell(""))
			line3.WriteString(g.formatGridCell(""))
		} else {
			room := g.state.Rooms[core.RoomID(roomID)]
			
			// Line 1: [RoomID,SearchStatus,BugCount]
			searchStatus := "+"
			if room.Searched {
				searchStatus = "-"
			}
			roomInfo := fmt.Sprintf("%s,%s,%d", roomID, searchStatus, room.BugMarkers)
			line1.WriteString(g.formatGridCell(roomInfo))
			
			// Line 2: Room type
			roomType := g.getRoomTypeDisplay(room)
			line2.WriteString(g.formatGridCell(roomType))
			
			// Line 3: Contents
			contents := g.formatGridContents(room, overflowRooms)
			line3.WriteString(g.formatGridCell(contents))
		}
	}
	
	return line1.String(), line2.String(), line3.String()
}

func (g *GameManager) findRoomAt(row, col int) string {
	target := core.Coord{Row: row, Col: col}
	for roomID, pos := range core.ROOM_POSITIONS {
		if pos == target {
			return roomID
		}
	}
	return ""
}

func (g *GameManager) getRoomTypeDisplay(room *core.RoomState) string {
	// Predefined rooms are always known
	switch room.ID {
	case "R01":
		return "KEY"
	case "R12":
		return "STR"
	case "R15":
		return "EN1"
	case "R17":
		return "EN3"
	case "R18":
		return "EN2"
	case "R19", "R20":
		return "ESC"
	}
	
	// If not explored, show XXX
	if !room.Explored {
		return "XXX"
	}
	
	// Otherwise show actual room type
	switch room.Type {
	case core.AmmoCache:
		return "AMO"
	case core.MedBay:
		return "MED"
	case core.CleanRoomType:
		return "CLN"
	case core.EnemySpawn:
		return "SPN"
	case core.Empty:
		return "AIR"
	default:
		return "???"
	}
}

func (g *GameManager) formatGridContents(room *core.RoomState, overflowRooms map[core.RoomID][]string) string {
	var objects []string
	
	// Add players
	for pid, player := range g.state.Players {
		if player.Location == room.ID {
			objects = append(objects, string(pid))
		}
	}
	
	// Add enemies
	for _, enemy := range g.state.Enemies {
		if enemy.Location == room.ID {
			objects = append(objects, g.getEnemyAbbrev(enemy.Type))
		}
	}
	
	// Bug markers are NOT objects - they're shown in the room ID line [R07,+,1]
	// Don't add bug markers to objects list
	
	// Handle display based on object count
	if len(objects) == 0 {
		return ""
	} else if len(objects) <= 2 {
		joined := strings.Join(objects, ",")
		if len(joined) <= 8 { // Fits in 8 character display area
			return joined
		} else {
			// Even 2 objects are too long, store in overflow
			overflowRooms[room.ID] = objects
			return fmt.Sprintf("%d+ OBJ", len(objects))
		}
	} else {
		// 3+ objects, always store in overflow
		overflowRooms[room.ID] = objects
		return fmt.Sprintf("%d+ OBJ", len(objects))
	}
}

func (g *GameManager) getEnemyAbbrev(enemyType core.EnemyType) string {
	switch enemyType {
	case core.InfiniteLoop:
		return "IL"
	case core.StackOverflow:
		return "SO"
	case core.Pythogoras:
		return "PY"
	default:
		return "??"
	}
}

func (g *GameManager) showHelp() error {
	fmt.Println("\n=== Devesis: Tutorial Hell - Commands ===")
	fmt.Println()
	fmt.Println("Turn-economy actions (cost a turn):")
	fmt.Println("  move <roomID>  (mv)  - Move to adjacent room")
	fmt.Println("  play <cardID>  (c)   - Play a card from hand")
	fmt.Println("  search         (s)   - Search current room")
	fmt.Println("  shoot          (f)   - Attack enemies in adjacent rooms")
	fmt.Println("  melee          (ml)  - Attack enemies in current room")
	fmt.Println("  room           (ra)  - Use room's special ability")
	fmt.Println("  pass           (p)   - End turn without action")
	fmt.Println()
	fmt.Println("Information commands (free):")
	fmt.Println("  hand           (h)   - Show your cards")
	fmt.Println("  map            (mp)  - Display game map")
	fmt.Println("  status         (st)  - Show current status")
	fmt.Println("  help           (?)   - Show this help")
	fmt.Println("  rule           (ru)  - Show game rules (pager view)")
	fmt.Println("  quit           (q)   - Exit game")
	fmt.Println()
	
	return nil
}

func (g *GameManager) showRules() error {
	rulesContent := `
DEVESIS: TUTORIAL HELL - GAME RULES
==================================

OBJECTIVE
---------
Escape the corrupted spaceship USS Boot.dev before time runs out! Collect engine components
and activate the escape pods while surviving the programming malware that infests the ship.

WIN CONDITION
-------------
1. Search engine rooms (R15, R17, R18) to collect 3 Engine Core cards
2. Navigate to escape rooms (R19 or R20)  
3. Play an Engine Core card at the escape room (if no Pythogoras present)
4. Victory! You've escaped Tutorial Hell!

TURN STRUCTURE (4 Phases)
-------------------------
1. DRAW PHASE: Auto-refill hand to 5 cards from deck
2. PLAYER PHASE: Take up to 2 actions per turn
3. EVENT PHASE: Time decreases, enemies attack/move, corruption spreads
4. ROUND MAINTENANCE: Advance to next round

PLAYER ACTIONS (Cost 1 Action Each)
----------------------------------
• move <roomID>    - Move to adjacent room (triggers coding question)
• play <cardID>    - Play a card from your hand  
• search           - Search current room for special items
• shoot            - Attack enemies in adjacent rooms (costs 1 ammo)
• melee            - Attack enemies in current room (no ammo cost)
• room             - Use current room's special ability
• pass             - End turn early

INFORMATION COMMANDS (Free)
--------------------------
• hand             - Show cards in hand
• map              - Display ship layout
• status           - Show player stats and room info
• help             - Show command help
• rule             - Show these rules (you're here!)

SPECIAL ROOMS
-------------
• R01 (KEY): Search to gain BOOT.dev KEY (increases damage from 1 to 3)
• R15/R17/R18 (Engines): Search to gain 3 Engine Core cards
• R19/R20 (Escape): Play Engine Core here to win (if no Pythogoras)
• R12 (Start): Your starting location

ENEMIES
-------
• Infinite Loop (1 HP, 1 DMG): Weak but numerous
• Stack Overflow (3 HP, 1 DMG): Medium threat  
• Pythogoras (6 HP, 1 DMG): Powerful boss - blocks escape rooms

RESOURCES
---------
• HP: Health points - game over if reduced to 0
• Ammo: Required for shooting attacks
• Cards: Hand limit of 6 cards, excess go to discard
• Time: 15 rounds total - game over if time expires

MOVEMENT & QUESTIONS
------------------
Moving between rooms triggers coding questions. Wrong answers add bug markers
to rooms. At 3+ bugs, rooms become corrupted and spawn more enemies.

CARD SYSTEM
-----------
• Deck starts with 10 random action cards
• Hand refills to 5 cards each turn
• When deck empty, discard pile shuffles back into deck
• Special cards found by searching rooms

TIPS FOR SURVIVAL
-----------------
1. Search the KEY room (R01) early for damage boost
2. Collect all 3 engines before heading to escape
3. Clear Pythogoras from escape rooms before playing Engine Core
4. Manage ammo and HP carefully - use room abilities when possible
5. Answer coding questions correctly to avoid corruption

Press 'q' to exit this view.
`

	return g.showInPager(rulesContent)
}

func (g *GameManager) showInPager(content string) error {
	// Try to use system pager (less, more, etc.)
	pager := os.Getenv("PAGER")
	if pager == "" {
		// Default pagers in order of preference
		pagers := []string{"less", "more", "cat"}
		for _, p := range pagers {
			if _, err := exec.LookPath(p); err == nil {
				pager = p
				break
			}
		}
	}

	if pager == "" {
		// Fallback: just print to stdout
		fmt.Print(content)
		return nil
	}

	// Create temporary file with content
	tmpfile, err := ioutil.TempFile("", "devesis-rules-*.txt")
	if err != nil {
		// Fallback: just print to stdout
		fmt.Print(content)
		return nil
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		tmpfile.Close()
		fmt.Print(content)
		return nil
	}
	tmpfile.Close()

	// Launch pager with the temporary file
	cmd := exec.Command(pager, tmpfile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	return cmd.Run()
}