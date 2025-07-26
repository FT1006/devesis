package core

type GameState struct {
	Round      int
	Time       int
	RandSeed   int64
	EventIndex uint8
	
	Rooms   map[RoomID]*RoomState
	Players map[PlayerID]*PlayerState
	Events  []EventCard
	Bag     []Token
	Enemies map[EnemyID]*Enemy
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
	Hand         []Card
	Deck         []Card
	Discard      []Card
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
	CleanRoom
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

type Coord struct {
	Row, Col int
}

type EventCard struct{}
type Card struct{}