package core

type Action interface {
	isAction()
}

type MoveAction struct {
	PlayerID PlayerID
	To       RoomID
}

func (MoveAction) isAction() {}

type SearchAction struct {
	PlayerID PlayerID
}

func (SearchAction) isAction() {}

type ShootAction struct {
	PlayerID PlayerID
	// Hits all enemies in surrounding rooms
}

func (ShootAction) isAction() {}

type MeleeAction struct {
	PlayerID PlayerID
	// Hits all enemies in same room
}

func (MeleeAction) isAction() {}

type RoomAction struct {
	PlayerID PlayerID
}

func (RoomAction) isAction() {}

type SpecialAction struct {
	PlayerID PlayerID
}

func (SpecialAction) isAction() {}

type PlayCardAction struct {
	PlayerID PlayerID
	CardID   CardID
}

func (PlayCardAction) isAction() {}

type PassAction struct {
	PlayerID PlayerID
}

func (PassAction) isAction() {}

type InitializeGameAction struct {
	Seed        int64
	PlayerClass DevClass
}

func (InitializeGameAction) isAction() {}