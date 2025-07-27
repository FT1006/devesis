package cards

import "github.com/spaceship/devesis/pkg/cards/effects"

// Re-export types for convenience
type EffectOp = effects.EffectOp
type ScopeType = effects.ScopeType
type EffectSource = effects.EffectSource
type Effect = effects.Effect
type Card = effects.Card

// Re-export constants
const ALL = effects.ALL

const (
	ModifyHP     = effects.ModifyHP
	ModifyAmmo   = effects.ModifyAmmo
	DrawCards    = effects.DrawCards
	DiscardCards = effects.DiscardCards
	SkipQuestion = effects.SkipQuestion
	ModifyBugs   = effects.ModifyBugs
	RevealRoom   = effects.RevealRoom
	CleanRoom    = effects.CleanRoom
	SetCorrupted = effects.SetCorrupted
	SpawnEnemy   = effects.SpawnEnemy
)

const (
	Self             = effects.Self
	CurrentRoom      = effects.CurrentRoom
	AdjacentRooms    = effects.AdjacentRooms
	AllRooms         = effects.AllRooms
	RoomWithMostBugs = effects.RoomWithMostBugs
	AllPlayers       = effects.AllPlayers
)

const (
	SrcAction  = effects.SrcAction
	SrcEvent   = effects.SrcEvent
	SrcEnemyAI = effects.SrcEnemyAI
)