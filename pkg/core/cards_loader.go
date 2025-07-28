package core

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// CardDatabase represents the YAML structure
type CardDatabase struct {
	Cards struct {
		Action  []YAMLCard `yaml:"action"`
		Special []YAMLCard `yaml:"special"`
		Event   []YAMLCard `yaml:"event"`
	} `yaml:"cards"`
}

// YAMLCard represents a card as stored in YAML
type YAMLCard struct {
	ID       string    `yaml:"id"`
	Name     string    `yaml:"name"`
	Desc     string    `yaml:"desc"`
	Category string    `yaml:"category"`
	Source   string    `yaml:"source"`
	Rarity   string    `yaml:"rarity,omitempty"`
	Flavor   string    `yaml:"flavor,omitempty"`
	FX       []YAMLFx  `yaml:"fx"`
}

// YAMLFx represents an effect as stored in YAML
type YAMLFx struct {
	Op    string `yaml:"op"`
	Scope string `yaml:"scope"`
	N     int    `yaml:"n"`
}

var CardDB map[CardID]Card
var loadError error

// LoadCards loads the card database from YAML file
func LoadCards(dataPath string) error {
	cardFilePath := filepath.Join(dataPath, "cards.yaml")
	
	data, err := ioutil.ReadFile(cardFilePath)
	if err != nil {
		return fmt.Errorf("failed to read cards file: %w", err)
	}

	var db CardDatabase
	if err := yaml.Unmarshal(data, &db); err != nil {
		return fmt.Errorf("failed to parse cards YAML: %w", err)
	}

	CardDB = make(map[CardID]Card)

	// Process all card categories
	allCards := []YAMLCard{}
	allCards = append(allCards, db.Cards.Action...)
	allCards = append(allCards, db.Cards.Special...)
	allCards = append(allCards, db.Cards.Event...)

	for _, yamlCard := range allCards {
		card, err := convertYAMLToCard(yamlCard)
		if err != nil {
			return fmt.Errorf("failed to convert card %s: %w", yamlCard.ID, err)
		}
		CardDB[CardID(yamlCard.ID)] = card
	}

	return nil
}

// GetCard retrieves a card by ID
func GetCard(cardID CardID) (Card, error) {
	if CardDB == nil {
		return Card{}, fmt.Errorf("card database not loaded")
	}

	card, exists := CardDB[cardID]
	if !exists {
		return Card{}, fmt.Errorf("card not found: %s", cardID)
	}

	return card, nil
}

// convertYAMLToCard converts YAML card format to core.Card
func convertYAMLToCard(yamlCard YAMLCard) (Card, error) {
	// Convert source string to EffectSource
	var source EffectSource
	switch yamlCard.Source {
	case "action":
		source = SrcAction
	case "event":
		source = SrcEvent
	case "special":
		source = SrcSpecial
	default:
		return Card{}, fmt.Errorf("unknown source: %s", yamlCard.Source)
	}

	// Convert effects
	var effectsList []Effect
	for i, fx := range yamlCard.FX {
		// Convert op string to EffectOp
		op, err := stringToEffectOp(fx.Op)
		if err != nil {
			return Card{}, fmt.Errorf("effect %d: %w", i, err)
		}

		// Convert scope string to ScopeType
		scope, err := stringToScopeType(fx.Scope)
		if err != nil {
			return Card{}, fmt.Errorf("effect %d: %w", i, err)
		}

		effect := Effect{
			Op:    op,
			Scope: scope,
			N:     fx.N,
		}

		effectsList = append(effectsList, effect)
	}

	card := Card{
		ID:      yamlCard.ID,
		Name:    yamlCard.Name,
		Description: yamlCard.Desc,
		Source:  source,
		Effects: effectsList,
	}

	return card, nil
}

// stringToEffectOp converts string to EffectOp
func stringToEffectOp(s string) (EffectOp, error) {
	switch s {
	case "ModifyHP":
		return ModifyHP, nil
	case "ModifyAmmo":
		return ModifyAmmo, nil
	case "DrawCards":
		return DrawCards, nil
	case "DiscardCards":
		return DiscardCards, nil
	case "OutOfRam":
		return OutOfRam, nil
	case "ModifyBugs":
		return ModifyBugs, nil
	case "RevealRoom":
		return RevealRoom, nil
	case "CleanRoom":
		return CleanRoom, nil
	case "SetCorrupted":
		return SetCorrupted, nil
	case "SpawnEnemy":
		return SpawnEnemy, nil
	case "MoveEnemies":
		return MoveEnemies, nil
	default:
		return 0, fmt.Errorf("unknown effect op: %s", s)
	}
}

// stringToScopeType converts string to ScopeType
func stringToScopeType(s string) (ScopeType, error) {
	switch s {
	case "Self":
		return Self, nil
	case "CurrentRoom":
		return CurrentRoom, nil
	case "AdjacentRooms":
		return AdjacentRooms, nil
	case "AllRooms":
		return AllRooms, nil
	case "RoomWithMostBugs":
		return RoomWithMostBugs, nil
	case "RoomWithMostEnemies":
		return RoomWithMostEnemies, nil
	case "AllPlayers":
		return AllPlayers, nil
	default:
		return 0, fmt.Errorf("unknown scope type: %s", s)
	}
}