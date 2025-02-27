package entity

import (
	"math"
	"time"

	"nightmare/internal/common"
	"nightmare/internal/world"
)

// Константы игрока
const (
	MaxHealth     = 100
	MaxSanity     = 100
	MoveSpeed     = 3.0
	RotationSpeed = 0.05
)

// PlayerAction представляет действие игрока
type PlayerAction int

const (
	ActionMove PlayerAction = iota
	ActionInteract
	ActionRun
	ActionHide
)

// PlayerActionRecord записывает действия игрока с временной меткой
type PlayerActionRecord struct {
	Action          common.PlayerActionType
	Timestamp       time.Time
	Position        common.Vector2D
	InteractionType string // Добавим это поле для решения ошибки
}

// Player представляет игрока
type Player struct {
	Position  Vector2D
	Direction float64 // угол в радианах
	Health    float64
	Sanity    float64
	Inventory []Item
	ActionLog []PlayerActionRecord // история действий для анализа ИИ
}

// Vector2D представляет 2D вектор
type Vector2D struct {
	X, Y float64
}

// Item представляет предмет в инвентаре
type Item struct {
	ID          int
	Name        string
	Description string
	// Другие свойства предмета
}

// NewPlayer создает нового игрока
func NewPlayer() *Player {
	return &Player{
		Position:  Vector2D{X: 128, Y: 128}, // Начальная позиция в центре мира
		Direction: 0,                        // Начальное направление (вперед)
		Health:    MaxHealth,
		Sanity:    MaxSanity,
		Inventory: []Item{},
		ActionLog: []PlayerActionRecord{},
	}
}

// Update обновляет состояние игрока
func (p *Player) Update() {
	// Здесь может быть логика автоматического изменения состояния игрока
	// Например, медленное восстановление здоровья или снижение рассудка в темноте
}

// MoveForward перемещает игрока вперед
func (p *Player) MoveForward() {
	dx := MoveSpeed * math.Cos(p.Direction)
	dy := MoveSpeed * math.Sin(p.Direction)
	p.Position.X += dx
	p.Position.Y += dy

	p.recordAction(ActionMove)
}

// MoveBackward перемещает игрока назад
func (p *Player) MoveBackward() {
	dx := MoveSpeed * math.Cos(p.Direction)
	dy := MoveSpeed * math.Sin(p.Direction)
	p.Position.X -= dx
	p.Position.Y -= dy

	p.recordAction(ActionMove)
}

// TurnLeft поворачивает игрока влево
func (p *Player) TurnLeft() {
	p.Direction -= RotationSpeed
	if p.Direction < 0 {
		p.Direction += 2 * math.Pi
	}
}

// TurnRight поворачивает игрока вправо
func (p *Player) TurnRight() {
	p.Direction += RotationSpeed
	if p.Direction >= 2*math.Pi {
		p.Direction -= 2 * math.Pi
	}
}

// Interact взаимодействует с миром
func (p *Player) Interact(w *world.World) {
	// Здесь будет логика взаимодействия с объектами мира
	p.recordAction(ActionInteract)
}

// TakeDamage наносит урон игроку
func (p *Player) TakeDamage(amount float64) {
	p.Health -= amount
	if p.Health < 0 {
		p.Health = 0
	}
}

// ReduceSanity снижает рассудок игрока
func (p *Player) ReduceSanity(amount float64) {
	p.Sanity -= amount
	if p.Sanity < 0 {
		p.Sanity = 0
	}
}

// AddItem добавляет предмет в инвентарь
func (p *Player) AddItem(item Item) {
	p.Inventory = append(p.Inventory, item)
}

// recordAction записывает действие игрока в лог
func (p *Player) recordAction(action PlayerAction) {
	record := PlayerActionRecord{
		Action:    action,
		Timestamp: time.Now(),
		Position:  p.Position,
	}
	p.ActionLog = append(p.ActionLog, record)

	// Ограничиваем размер лога, чтобы не расходовать слишком много памяти
	if len(p.ActionLog) > 1000 {
		p.ActionLog = p.ActionLog[len(p.ActionLog)-1000:]
	}
}
