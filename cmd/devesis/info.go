package main

import (
	"fmt"
	"strings"

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
		for i, card := range player.Hand {
			fmt.Printf("  %d. %v\n", i+1, card) // TODO: Proper card display
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
	fmt.Println("  quit           (q)   - Exit game")
	fmt.Println()
	
	return nil
}