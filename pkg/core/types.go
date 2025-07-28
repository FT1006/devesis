package core

type GameState struct {
	Round      int
	Time       int
	RandSeed   int64
	EventIndex uint8
	
	// Turn controller fields
	ActionsLeft   int      // 0-2 actions remaining for current player
	Phase         string   // "player" or "event"
	ActivePlayer  PlayerID
	
	Rooms         map[RoomID]*RoomState
	Players       map[PlayerID]*PlayerState
	Events        []EventCard
	SpawnBag      *SpawnBag
	Enemies       map[EnemyID]*Enemy
	// Question system using pre-shuffle approach
	QuestionOrder []int // Pre-shuffled order of question IDs 0-49
	NextQuestion  int   // Index of next question to use
}

type RoomState struct {
	ID         RoomID
	Type       RoomType
	Explored   bool
	Searched   bool
	Corrupted  bool
	OutOfRam   bool
	BugMarkers uint8
}

type PlayerState struct {
	ID           PlayerID
	Class        DevClass
	HP           uint8
	MaxHP        uint8
	Ammo         uint8
	MaxAmmo      uint8
	Hand         []CardID
	Deck         []CardID
	Discard      []CardID
	Location     RoomID
	HasActed     bool
	SpecialUsed  bool
	PersonalObj  ObjectiveID
	CorporateObj ObjectiveID
}

type Enemy struct {
	ID       EnemyID
	Type     EnemyType
	HP       uint8
	MaxHP    uint8
	Damage   uint8
	Location RoomID
}

type RoomID string
type PlayerID string
type EnemyID string
type CardID string
type ObjectiveID string

type DevClass int
const (
	Frontend DevClass = iota
	Backend
	DevOps
	Fullstack
)

type EnemyType int
const (
	InfiniteLoop EnemyType = iota
	StackOverflow
	Pythogoras
)

type RoomType int
const (
	Predefined RoomType = iota
	AmmoCache
	MedBay
	CleanRoomType
	EnemySpawn
	Empty
)

type Token int
const (
	NoSpawn Token = iota
	LoopToken
	OverflowToken
	PythagorasToken
)

// SpawnBag represents the enemy spawn pool
type SpawnBag struct {
	Tokens []EnemyType // Each entry is one spawn chance
}

type Coord struct {
	Row, Col int
}

type EventCard struct{}
// Card reference - actual cards live in pkg/cards
// Players hold CardIDs, cards are resolved when played

type Question struct {
	ID            int
	Text          string
	Options       [4]string
	CorrectAnswer int
}