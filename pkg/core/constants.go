package core

const (
	MaxRounds     = 15
	MaxHandSize   = 6
	MaxBugMarkers = 9  // Max bugs per room
	BugCorruptionThreshold = 3  // Rooms corrupt at 3+ bugs

	GridRows = 7
	GridCols = 7

	// Action costs
	SearchDiscardCost = 1
	MeleeAmmoCost     = 0
	ShootAmmoCost     = 1
	
	// Combat damage
	BasicDamage = 1    // Default player damage
	BootDevDamage = 3  // KEY item damage bonus
	
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

// Class Stats: HP and Ammo capacity by developer class
var CLASS_STATS = map[DevClass]struct {
	HP       uint8
	MaxAmmo  uint8
}{
	Frontend:  {HP: 6, MaxAmmo: 3},
	Backend:   {HP: 3, MaxAmmo: 6},
	DevOps:    {HP: 4, MaxAmmo: 5},
	Fullstack: {HP: 5, MaxAmmo: 4},
}

// Enemy Stats: HP and attack damage by enemy type
var ENEMY_STATS = map[EnemyType]struct {
	HP     uint8
	Damage uint8
}{
	InfiniteLoop:   {HP: 1, Damage: 1},
	StackOverflow:  {HP: 3, Damage: 1},
	Pythogoras:     {HP: 6, Damage: 1},
}