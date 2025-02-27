package entity

import (
	"math"
	"time"

	"nightmare/internal/common"
)

// Constants for player
const (
	MaxHealth     = 100
	MaxSanity     = 100
	MoveSpeed     = 3.0
	RotationSpeed = 0.05
)

// PlayerAction represents a player action
type PlayerAction int

const (
	ActionMove PlayerAction = iota
	ActionInteract
	ActionRun
	ActionHide
)

// PlayerActionRecord records player actions with a timestamp
type PlayerActionRecord struct {
	Action          PlayerAction // Using entity's PlayerAction, not common.PlayerActionType
	Timestamp       time.Time
	Position        Vector2D // Using entity's Vector2D
	InteractionType string
}

// Player represents the player
type Player struct {
	Position  Vector2D
	Direction float64 // angle in radians
	Health    float64
	Sanity    float64
	Inventory []Item
	ActionLog []PlayerActionRecord // action history for AI analysis
}

// Vector2D represents a 2D vector
type Vector2D struct {
	X, Y float64
}

// Item represents an item in the inventory
type Item struct {
	ID          int
	Name        string
	Description string
	// Other item properties
}

// NewPlayer creates a new player
func NewPlayer() *Player {
	return &Player{
		Position:  Vector2D{X: 128, Y: 128}, // Initial position in the center of the world
		Direction: 0,                        // Initial direction (forward)
		Health:    MaxHealth,
		Sanity:    MaxSanity,
		Inventory: []Item{},
		ActionLog: []PlayerActionRecord{},
	}
}

// Update updates the player's state
func (p *Player) Update() {
	// Logic for automatic changes to player state can go here
	// For example, slow health regeneration or sanity decrease in darkness
}

// MoveForward moves the player forward
func (p *Player) MoveForward() {
	dx := MoveSpeed * math.Cos(p.Direction)
	dy := MoveSpeed * math.Sin(p.Direction)
	p.Position.X += dx
	p.Position.Y += dy

	p.recordAction(ActionMove)
}

// MoveBackward moves the player backward
func (p *Player) MoveBackward() {
	dx := MoveSpeed * math.Cos(p.Direction)
	dy := MoveSpeed * math.Sin(p.Direction)
	p.Position.X -= dx
	p.Position.Y -= dy

	p.recordAction(ActionMove)
}

// TurnLeft turns the player left
func (p *Player) TurnLeft() {
	p.Direction -= RotationSpeed
	if p.Direction < 0 {
		p.Direction += 2 * math.Pi
	}
}

// TurnRight turns the player right
func (p *Player) TurnRight() {
	p.Direction += RotationSpeed
	if p.Direction >= 2*math.Pi {
		p.Direction -= 2 * math.Pi
	}
}

// Interact interacts with the world
func (p *Player) Interact(w interface{}) {
	// Logic for interacting with world objects will go here
	p.recordAction(ActionInteract)
}

// TakeDamage damages the player
func (p *Player) TakeDamage(amount float64) {
	p.Health -= amount
	if p.Health < 0 {
		p.Health = 0
	}
}

// ReduceSanity reduces the player's sanity
func (p *Player) ReduceSanity(amount float64) {
	p.Sanity -= amount
	if p.Sanity < 0 {
		p.Sanity = 0
	}
}

// AddItem adds an item to the inventory
func (p *Player) AddItem(item Item) {
	p.Inventory = append(p.Inventory, item)
}

// recordAction records a player action in the log
func (p *Player) recordAction(action PlayerAction) {
	record := PlayerActionRecord{
		Action:    action,
		Timestamp: time.Now(),
		Position:  p.Position,
	}
	p.ActionLog = append(p.ActionLog, record)

	// Limit log size to avoid using too much memory
	if len(p.ActionLog) > 1000 {
		p.ActionLog = p.ActionLog[len(p.ActionLog)-1000:]
	}
}

// ConvertToCommonAction converts the entity PlayerAction to common.PlayerActionType
func ConvertToCommonAction(action PlayerAction) common.PlayerActionType {
	switch action {
	case ActionMove:
		return common.ActionMove
	case ActionInteract:
		return common.ActionInteract
	case ActionRun:
		return common.ActionRun
	case ActionHide:
		return common.ActionHide
	default:
		return common.ActionMove
	}
}

// ToCommonVector converts entity Vector2D to common.Vector2D
func (v Vector2D) ToCommonVector() common.Vector2D {
	return common.Vector2D{X: v.X, Y: v.Y}
}

// FromCommonVector converts common.Vector2D to entity Vector2D
func FromCommonVector(v common.Vector2D) Vector2D {
	return Vector2D{X: v.X, Y: v.Y}
}

// ToCommonRecord converts entity PlayerActionRecord to common.PlayerAction
func (r PlayerActionRecord) ToCommonRecord() common.PlayerAction {
	return common.PlayerAction{
		Type:            ConvertToCommonAction(r.Action),
		Position:        r.Position.ToCommonVector(),
		Direction:       common.Vector2D{X: 0, Y: 0}, // Default direction
		Timestamp:       r.Timestamp,
		InteractionType: r.InteractionType,
	}
}
