// internal/common/types.go
package common

import (
	"math"
	"time"
)

// Vector2D represents a 2D vector - common type for different packages
type Vector2D struct {
	X, Y float64
}

// PlayerActionType represents a player action type
type PlayerActionType int

const (
	ActionMove PlayerActionType = iota
	ActionInteract
	ActionRun
	ActionHide
	ActionAttack
	ActionInvestigate
	ActionRetreat
	ActionFreeze
)

// FearType represents a type of fear
type FearType int

const (
	FearDarkness FearType = iota
	FearCreatures
	FearSuddenNoises
	FearIsolation
	FearChasing
	FearGore
	FearClaustrophobia
	FearOpenSpaces
	FearUnknown
)

// ScareEventType represents a type of scare event
type ScareEventType int

const (
	EventAmbientSound ScareEventType = iota
	EventSuddenNoise
	EventCreatureAppearance
	EventEnvironmentChange
	EventHallucination
	EventWhisper
)

// TileType represents a tile type in the world
type TileType int

const (
	TileGrass TileType = iota
	TileForest
	TileDenseForest
	TilePath
	TileRocks
	TileWater
	TileSwamp
	TileCorrupted
)

// EntityType represents the type of entity
type EntityType string

// CreatureBehaviorType represents the type of creature behavior
type CreatureBehaviorType int

const (
	BehaviorPassive CreatureBehaviorType = iota
	BehaviorAggressive
	BehaviorStalker
	BehaviorFleeing
	BehaviorPatrol
	BehaviorHunter
)

// WorldObject represents a base structure for an object in the world
type WorldObject struct {
	ID          int
	Type        string
	Position    Vector2D
	Solid       bool
	Interactive bool
}

// PlayerAction represents a player action record
type PlayerAction struct {
	Type            PlayerActionType
	Position        Vector2D
	Direction       Vector2D
	Timestamp       time.Time
	InteractionType string
}

// ScareEvent represents a scare event
type ScareEvent struct {
	Type         ScareEventType
	Intensity    float64 // From 0 to 1
	Position     Vector2D
	Duration     time.Duration
	CreatureType string // Type of creature if the event is related to a creature
	Timestamp    time.Time
}

// Distance calculates the distance between two points
func Distance(a, b Vector2D) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// ConvertPosition creates a new Vector2D based on other coordinate type
func ConvertPosition(x, y float64) Vector2D {
	return Vector2D{X: x, Y: y}
}

// VectorFromEntity converts an entity Vector2D to common Vector2D
func VectorFromEntity(vec interface{}) Vector2D {
	// Type assertion for different vector types
	// This is a generic conversion function that will handle
	// different vector types that might be in the codebase

	// For example, if vec is already a common.Vector2D
	if v, ok := vec.(Vector2D); ok {
		return v
	}

	// Add other type conversions as needed
	// For now, return a zero vector as fallback
	return Vector2D{0, 0}
}
