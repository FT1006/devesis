# Devesis: Tutorial Hell

```
████████▄     ▄████████  ▄█    █▄     ▄████████    ▄████████  ▄█     ▄████████ 
███   ▀███   ███    ███ ███    ███   ███    ███   ███    ███ ███    ███    ███ 
███    ███   ███    █▀  ███    ███   ███    █▀    ███    █▀  ███▌   ███    █▀  
███    ███  ▄███▄▄▄     ███    ███  ▄███▄▄▄       ███        ███▌   ███      
███    ███ ▀▀███▀▀▀     ███    ███ ▀▀███▀▀▀     ▀███████████ ███▌ ▀███████████ 
███    ███   ███    █▄  ███    ███   ███    █▄           ███ ███           ███ 
███   ▄███   ███    ███ ███    ███   ███    ███    ▄█    ███ ███     ▄█    ███ 
████████▀    ██████████  ▀██████▀    ██████████  ▄████████▀  █▀    ▄████████▀  
```

**A CLI survival horror game inspired by Nemesis and Boot.dev where hackathon developers escape Tutorial Hell**

*You are a self-learning dev get trapped in an endless cycle of coding tutorials during a hackathon. You have only 15 turns to survive and escape with your PROGRAMMING KNOWLEDGE.*

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org)
[![Build Status](https://img.shields.io/badge/build-passing-green.svg)]()
[![License](https://img.shields.io/badge/license-MIT-blue.svg)]()

## 🚀 Quick Start (< 5 Minutes)

### Prerequisites

- Go 1.21 or higher
- Terminal that supports ANSI colors

### Installation & Run

```bash
# Clone and run immediately
git clone https://github.com/yourusername/devesis.git
cd devesis
go run ./cmd/devesis

# Or build and run
go build -o devesis ./cmd/devesis
./devesis
```

That's it! The game will start immediately with character selection.

## 🎮 Game Overview

### The Premise

You're part of a team of self-learning developers participating in a hackathon. But something goes wrong—you get trapped in **Tutorial Hell**, an infinite loop of corrupted programming tutorials where the only way out is to debug your way to freedom before your `sanity.exe` stops responding.

### Choose Your Class

- **Frontend Developer** (HP: 6, Ammo: 3) - High survivability, limited resources
- **Backend Developer** (HP: 3, Ammo: 6) - Glass cannon with lots of firepower
- **DevOps Engineer** (HP: 4, Ammo: 5) - Balanced stats for infrastructure warfare
- **Fullstack Developer** (HP: 5, Ammo: 4) - Jack of all trades, master of survival

### Face Programming Enemies

- **Infinite Loop** (1 HP) - Weak but numerous, spawn from corrupted rooms
- **Stack Overflow** (3 HP) - Mid-level threats that pack a punch
- **Pythogoras** (6 HP) - The serpent god of tutorials, blocks your escape

### Game Mechanics

**Turn Structure**: Each round has 4 phases - Draw cards, Player actions (2 per turn), Event phase (enemies attack/move), Round maintenance.

**Movement & Learning**: Moving between rooms triggers coding questions. Correct answers = safe passage. Wrong answers spawn bugs that corrupt rooms and attract enemies.

**Card System**: Draw cards each turn and play them for actions. Hand limit of 6 cards - excess goes to discard pile.

**Bug Corruption**: Rooms with 3+ bugs become corrupted and spawn Infinite Loop enemies every event phase.

**Combat**: Battle with buggy enemies using melee attacks (free but dangerous) or shooting (costs ammo but can target adjacent rooms).

**Search & Discovery**: Search rooms to find special cards and items. Engine rooms contain Engine Core cards needed for victory.

**Room Actions**: Each room type has special abilities - Medical rooms heal HP, Ammo Caches refill ammunition, Clean Rooms remove bugs.

**Resource Management**: Balance HP, ammo, and cards across 15 intense rounds. Some rooms provide healing, ammo, or bug cleanup.

### How a Round Plays Out

**Win Condition**: Find and play an Engine card to activate escape pods in escape room when no Pythogoras is present.

Each round has **two main stages**:

**1. Player Stage** - You get 2 actions to:

- Move between rooms (triggers coding questions)
- Search rooms for cards and special items
- Play cards from your hand
- Fight enemies (melee or shoot)
- Use room abilities (heal, get ammo, clean bugs)

**2. Event Stage** - The system fights back:

- Enemies attack any players in their room
- System crashes damage enemies in OutOfRam rooms
- Event cards trigger various effects
- New enemies spawn and existing ones get stronger
- Corrupted rooms spawn additional Infinite Loops

After 15 rounds, if you haven't escaped or died, the system crashes and you lose.

## 🕹️ How to Play

### Basic Commands

```bash
# Movement and exploration
move R07          # Move to room R07 (triggers coding question)
search            # Search current room for special items
map               # Display the ship layout

# Combat and survival  
shoot             # Attack enemies in adjacent rooms (costs ammo)
melee             # Fight enemies in current room (no ammo cost)
play ACTION_001   # Play a card from your hand

# Information
hand              # Show cards in your hand
status            # Display player stats and room info
rule              # View complete game rules
list              # Browse all available cards
help              # Show all commands
```

### The Game Map

```
[          ] [          ] [          ] [ R19,+,0  ] [          ] [          ] 
[          ] [          ] [          ] [ ESC      ] [          ] [          ] 
[          ] [          ] [          ] [          ] [          ] [          ] 

[          ] [          ] [ R05,+,0  ] [ R10,+,0  ] [ R15,+,0  ] [          ] 
[          ] [          ] [ XXX      ] [ XXX      ] [ EN1      ] [          ] 
[          ] [          ] [          ] [          ] [          ] [          ] 

[          ] [ R02,+,0  ] [ R06,+,0  ] [ R11,+,0  ] [          ] [          ] 
[          ] [ XXX      ] [ XXX      ] [ XXX      ] [          ] [          ] 
[          ] [          ] [          ] [          ] [          ] [          ] 

[ R01,+,0  ] [ R03,+,0  ] [ R07,+,0  ] [ R12,+,0  ] [ R16,+,0  ] [ R18,+,0  ] 
[ KEY      ] [ XXX      ] [ XXX      ] [ STR      ] [ XXX      ] [ EN2      ] 
[          ] [          ] [          ] [ P1       ] [          ] [          ] 

[          ] [ R04,+,0  ] [ R08,+,0  ] [ R13,+,0  ] [          ] [          ] 
[          ] [ XXX      ] [ XXX      ] [ XXX      ] [          ] [          ] 
[          ] [          ] [          ] [          ] [          ] [          ] 

[          ] [          ] [ R09,+,0  ] [ R14,+,0  ] [ R17,+,0  ] [          ] 
[          ] [          ] [ XXX      ] [ XXX      ] [ EN3      ] [          ] 
[          ] [          ] [          ] [          ] [          ] [          ] 

[          ] [          ] [          ] [ R20,+,0  ] [          ] [          ] 
[          ] [          ] [          ] [ ESC      ] [          ] [          ] 
[          ] [          ] [          ] [          ] [          ] [          ] 
```

🗺️  **MAP LEGEND**:
• **Rooms**: [ID,±,B*] = [Room ID, Searched(+/-), Bug count, OutOfRam(*)]
• **Types**: KEY=Key STR=Start EN#=Engine ESC=Escape
         AMO=Ammo MED=Medical CLN=Clean AIR=Air SPN=Spawn
• **Units**: P#=Player IL=Infinite Loop SO=Stack Overflow PY=Pythogoras
• **Status**: XXX=Unexplored room, * = OutOfRam

### Sample Game Session

```
📍 CURRENT MAP STATE:
╔════════════════════════════════════════════════════════════════════════╗
║ 🐛 Total Bugs: 3    💀 Corrupted: 1   👹 Enemies: IL:2 SO:1 PY:0        ║
╚════════════════════════════════════════════════════════════════════════╝

┌ P1 Frontend ── Room R12 (start room, searched) ──────────────────────┐
│HP    6 /  6     Ammo  3 /  3   Damage  1                             │
│Turn   Actions 2 / 2      Cards  Hand:5  Deck:5  Discard:0            │
│Room   Bugs:0   Loop:0   Overflow:0   Pythogoras:0   Corrupted: ✘     │
│Game   Round: 1      Rounds left: 14                                  │
└──────────────────────────────────────────────────────────────────────┘

(2 actions left) > move R07

🤔 CODING QUESTION: What does 'git rebase' do?
A) Merges branches
B) Rewrites commit history  
C) Creates a new branch
D) Deletes commits

Your answer: B
✅ Correct! Safe movement to R07.
```

## 🏗️ Technical Architecture

### Built for Quality

- **Pure Reducer Pattern**: Deterministic game state transitions
- **Test-Driven Development**: Comprehensive test coverage for all core mechanics
- **Clean Architecture**: Separation of concerns with clear boundaries
- **Type Safety**: Strong typing throughout with Go's type system
- **REPL-like Experience**: Interactive command-line interface with immediate feedback

### Key Features

- **ASCII Art Rendering**: Beautiful terminal-based map display
- **Pager System**: Built-in paging for rules and card browsing
- **Effect System**: Flexible card effects with scoped targeting
- **Question Database**: Real programming questions with educational value
- **Save/Load**: JSON-based game state persistence

### Code Quality Highlights

```bash
# Run the comprehensive test suite
go test ./pkg/core -v

# Example: 30+ test cases covering edge cases
=== RUN   TestDrawPhaseEnforcesHandLimit
=== RUN   TestWinConditionsAllScenarios  
=== RUN   TestMovementWithCodingQuestions
=== RUN   TestBugCorruptionMechanics
```

---

*"Escape Tutorial Hell before your sanity.exe stops responding!"*

**Built with ❤️ during the Boot.dev Hackathon 2024**
