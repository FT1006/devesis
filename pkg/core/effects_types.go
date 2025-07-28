package core

// Magic number for "all"
const ALL = -999

// EffectOp enumeration
type EffectOp int

const (
	ModifyHP EffectOp = iota
	ModifyAmmo
	DrawCards
	DiscardCards
	SkipQuestion
	ModifyBugs
	RevealRoom
	CleanRoom
	SetCorrupted
	SpawnEnemy
)

// ScopeType enumeration
type ScopeType int

const (
	Self ScopeType = iota
	CurrentRoom
	AdjacentRooms
	AllRooms
	RoomWithMostBugs
	AllPlayers
)

// EffectSource enumeration
type EffectSource int

const (
	SrcAction EffectSource = iota
	SrcEvent
	SrcSpecial
)

// Effect represents a single effect to apply
type Effect struct {
	Op    EffectOp  `json:"op"`
	Scope ScopeType `json:"scope"`
	N     int       `json:"n"`
}

// Card represents a playable card with effects
type Card struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Source      EffectSource `json:"source"`
	Effects     []Effect     `json:"effects"`
}