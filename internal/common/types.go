// internal/common/types.go
package common

import (
	"math"
	"time"
)

// Vector2D представляет 2D вектор - общий тип для разных пакетов
type Vector2D struct {
	X, Y float64
}

// PlayerActionType представляет тип действия игрока
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

// FearType представляет тип страха
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

// ScareEventType представляет тип пугающего события
type ScareEventType int

const (
	EventAmbientSound ScareEventType = iota
	EventSuddenNoise
	EventCreatureAppearance
	EventEnvironmentChange
	EventHallucination
	EventWhisper
)

// TileType представляет тип тайла в мире
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

// WorldObject представляет базовую структуру для объекта в мире
type WorldObject struct {
	ID          int
	Type        string
	Position    Vector2D
	Solid       bool
	Interactive bool
}

// PlayerAction представляет запись действия игрока
type PlayerAction struct {
	Type            PlayerActionType
	Position        Vector2D
	Direction       Vector2D
	Timestamp       time.Time
	InteractionType string
}

// Distance вычисляет расстояние между двумя точками
func Distance(a, b Vector2D) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// ConvertPosition создает новый Vector2D на основе другого типа координат
func ConvertPosition(x, y float64) Vector2D {
	return Vector2D{X: x, Y: y}
}
