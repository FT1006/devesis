package core

const (
	MaxRounds     = 15
	MaxHandSize   = 6
	MaxBugMarkers = 255

	GridRows = 6
	GridCols = 7

	// Action costs
	SearchDiscardCost = 1
	MeleeAmmoCost     = 0
	ShootAmmoCost     = 1
	
	// Room abilities
	MedBayHealAmount  = 2
	AmmoCacheAmount   = 3
)

var ROOM_POSITIONS = map[string]Coord{
	"R01": {3, 0}, "R02": {2, 1}, "R03": {3, 1},
	"R04": {4, 1}, "R05": {1, 2}, "R06": {2, 2},
	"R07": {3, 2}, "R08": {4, 2}, "R09": {5, 2},
	"R10": {1, 3}, "R11": {2, 3}, "R12": {3, 3},
	"R13": {4, 3}, "R14": {5, 3}, "R15": {1, 4},
	"R16": {3, 4}, "R17": {5, 4}, "R18": {3, 5},
	"R19": {0, 3}, "R20": {6, 3},
}

// Fixed room types - others assigned dynamically from AmmoCache, MedBay, CleanRoom, EnemySpawn, Empty
var PREDEFINED_ROOMS = map[string]RoomType{
	"R01": Predefined, // Key
	"R12": Predefined, // Start
	"R15": Predefined, // Engine
	"R17": Predefined, // Engine
	"R18": Predefined, // Engine
	"R19": Predefined, // Escape
	"R20": Predefined, // Escape
}