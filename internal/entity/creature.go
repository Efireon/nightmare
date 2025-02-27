package entity

import (
	"math"
	"math/rand"
	"time"
)

// Типы поведения существ
const (
	BehaviorPassive = iota
	BehaviorAggressive
	BehaviorStalker
	BehaviorFleeing
	BehaviorPatrol
	BehaviorHunter
)

// Creature представляет собой существо в игре
type Creature struct {
	ID             int
	Position       Vector2D
	Direction      float64
	Speed          float64
	Health         float64
	Type           string
	BehaviorType   int
	TargetPos      Vector2D
	PlayerTarget   *Player
	DetectionRange float64
	AttackRange    float64
	AttackDamage   float64
	SanityDamage   float64
	Parts          []CreaturePart
	CurrentState   string
	StateTime      int
	LastSeen       time.Time
	IsVisible      bool
	StalkingTime   int
}

// CreaturePart представляет собой часть существа
type CreaturePart struct {
	Type        string
	TextureID   int
	Position    Vector2D
	Rotation    float64
	Scale       float64
	AnimFrames  []int
	CurrentAnim int
}

// NewCreature создает новое существо
func NewCreature(id int, creatureType string, position Vector2D) *Creature {
	c := &Creature{
		ID:             id,
		Position:       position,
		Direction:      rand.Float64() * 2 * math.Pi,
		Speed:          1.0 + rand.Float64()*1.5,
		Health:         50 + rand.Float64()*50,
		Type:           creatureType,
		DetectionRange: 10 + rand.Float64()*15,
		AttackRange:    1.5 + rand.Float64(),
		AttackDamage:   5 + rand.Float64()*10,
		SanityDamage:   2 + rand.Float64()*8,
		Parts:          []CreaturePart{},
		CurrentState:   "idle",
		StateTime:      0,
		LastSeen:       time.Now().Add(-10 * time.Minute), // Давно не видели
		IsVisible:      false,
		StalkingTime:   0,
	}

	// Устанавливаем поведение на основе типа
	switch creatureType {
	case "shadow":
		c.BehaviorType = BehaviorStalker
		c.Speed = 0.8 + rand.Float64()*0.4
		c.SanityDamage = 10 + rand.Float64()*15
	case "spider":
		c.BehaviorType = BehaviorAggressive
		c.Speed = 1.5 + rand.Float64()*1.0
		c.AttackDamage = 15 + rand.Float64()*10
	case "phantom":
		c.BehaviorType = BehaviorPatrol
		c.Speed = 1.0 + rand.Float64()*0.5
		c.SanityDamage = 5 + rand.Float64()*10
	case "wendigo":
		c.BehaviorType = BehaviorHunter
		c.Speed = 2.0 + rand.Float64()*1.0
		c.AttackDamage = 20 + rand.Float64()*15
	case "faceless":
		c.BehaviorType = BehaviorStalker
		c.Speed = 0.5 + rand.Float64()*0.3
		c.SanityDamage = 15 + rand.Float64()*10
	default:
		c.BehaviorType = BehaviorPassive
	}

	return c
}

// Update обновляет состояние существа
func (c *Creature) Update(worldWidth, worldHeight int) {
	c.StateTime++

	// Обновляем состояние на основе текущего поведения
	switch c.CurrentState {
	case "idle":
		// В состоянии покоя существо периодически меняет направление
		if c.StateTime > 60+rand.Intn(120) {
			c.Direction = rand.Float64() * 2 * math.Pi
			c.StateTime = 0

			// Иногда переходим в состояние блуждания
			if rand.Float64() < 0.7 {
				c.CurrentState = "wander"
				// Выбираем случайную точку назначения
				c.TargetPos = Vector2D{
					X: rand.Float64() * float64(worldWidth),
					Y: rand.Float64() * float64(worldHeight),
				}
			}
		}

	case "wander":
		// Двигаемся к целевой точке
		dir := c.getDirectionTo(c.TargetPos)
		c.Direction = c.smoothDirection(c.Direction, dir, 0.1)
		c.moveForward()

		// Проверяем, достигли ли цели
		dist := c.distanceTo(c.TargetPos)
		if dist < 1.0 || c.StateTime > 300 {
			c.CurrentState = "idle"
			c.StateTime = 0
		}

	case "chase":
		if c.PlayerTarget != nil {
			// Двигаемся к игроку
			dir := c.getDirectionTo(c.PlayerTarget.Position)
			c.Direction = c.smoothDirection(c.Direction, dir, 0.2)
			c.moveForward()

			// Проверяем, достаточно ли близко для атаки
			dist := c.distanceTo(c.PlayerTarget.Position)
			if dist < c.AttackRange {
				c.CurrentState = "attack"
				c.StateTime = 0
			}

			// Если игрок убежал слишком далеко, прекращаем погоню
			if dist > c.DetectionRange*1.5 {
				c.CurrentState = "search"
				c.StateTime = 0
				c.TargetPos = c.PlayerTarget.Position // Последнее известное положение игрока
			}
		} else {
			c.CurrentState = "idle"
			c.StateTime = 0
		}

	case "attack":
		// Анимация атаки длится 30 кадров
		if c.StateTime == 15 { // Середина анимации - момент атаки
			if c.PlayerTarget != nil {
				dist := c.distanceTo(c.PlayerTarget.Position)
				if dist < c.AttackRange {
					// Наносим урон игроку
					c.PlayerTarget.TakeDamage(c.AttackDamage)
					c.PlayerTarget.ReduceSanity(c.SanityDamage)
				}
			}
		}

		if c.StateTime >= 30 {
			c.CurrentState = "chase"
			c.StateTime = 0
		}

	case "search":
		// Ищем игрока в последнем известном местоположении
		dir := c.getDirectionTo(c.TargetPos)
		c.Direction = c.smoothDirection(c.Direction, dir, 0.1)
		c.moveForward()

		// Если достигли места и не нашли игрока, начинаем блуждать
		dist := c.distanceTo(c.TargetPos)
		if dist < 1.0 || c.StateTime > 180 {
			c.CurrentState = "wander"
			c.StateTime = 0
			// Выбираем новую случайную точку
			c.TargetPos = Vector2D{
				X: c.TargetPos.X + (rand.Float64()*20 - 10),
				Y: c.TargetPos.Y + (rand.Float64()*20 - 10),
			}
			// Ограничиваем координаты в пределах мира
			if c.TargetPos.X < 0 {
				c.TargetPos.X = 0
			}
			if c.TargetPos.Y < 0 {
				c.TargetPos.Y = 0
			}
			if c.TargetPos.X >= float64(worldWidth) {
				c.TargetPos.X = float64(worldWidth) - 1
			}
			if c.TargetPos.Y >= float64(worldHeight) {
				c.TargetPos.Y = float64(worldHeight) - 1
			}
		}

	case "stalk":
		// Преследуем игрока, но держимся на расстоянии
		if c.PlayerTarget != nil {
			targetDist := 8.0 + rand.Float64()*4.0 // Дистанция преследования
			dist := c.distanceTo(c.PlayerTarget.Position)

			if dist < targetDist-2.0 {
				// Слишком близко, отходим
				dir := c.getDirectionTo(c.PlayerTarget.Position) + math.Pi // Противоположное направление
				c.Direction = c.smoothDirection(c.Direction, dir, 0.1)
				c.moveForward()
			} else if dist > targetDist+2.0 {
				// Слишком далеко, приближаемся
				dir := c.getDirectionTo(c.PlayerTarget.Position)
				c.Direction = c.smoothDirection(c.Direction, dir, 0.1)
				c.moveForward()
			} else {
				// На оптимальной дистанции, просто перемещаемся вокруг игрока
				tangent := c.getDirectionTo(c.PlayerTarget.Position) + math.Pi/2
				c.Direction = c.smoothDirection(c.Direction, tangent, 0.05)
				c.moveForward()
			}

			// Иногда переходим в режим атаки
			if c.StalkingTime > 300 && rand.Float64() < 0.01 {
				c.CurrentState = "chase"
				c.StateTime = 0
				c.StalkingTime = 0
			} else {
				c.StalkingTime++
			}
		} else {
			c.CurrentState = "idle"
			c.StateTime = 0
		}

	case "flee":
		// Убегаем от игрока
		if c.PlayerTarget != nil {
			dir := c.getDirectionTo(c.PlayerTarget.Position) + math.Pi // Противоположное направление
			c.Direction = c.smoothDirection(c.Direction, dir, 0.2)
			c.moveForward()

			// Если убежали достаточно далеко, переходим в режим блуждания
			dist := c.distanceTo(c.PlayerTarget.Position)
			if dist > c.DetectionRange*2 || c.StateTime > 300 {
				c.CurrentState = "wander"
				c.StateTime = 0
			}
		} else {
			c.CurrentState = "idle"
			c.StateTime = 0
		}
	}

	// Обновляем анимацию для всех частей
	for i := range c.Parts {
		c.updatePartAnimation(&c.Parts[i])
	}
}

// SetTarget устанавливает игрока в качестве цели
func (c *Creature) SetTarget(player *Player) {
	c.PlayerTarget = player
	c.LastSeen = time.Now()

	// Реагируем на обнаружение игрока в зависимости от типа поведения
	switch c.BehaviorType {
	case BehaviorPassive:
		// Пассивные существа не реагируют или убегают
		if rand.Float64() < 0.7 {
			c.CurrentState = "flee"
		} else {
			c.CurrentState = "idle"
		}

	case BehaviorAggressive:
		// Агрессивные сразу атакуют
		c.CurrentState = "chase"

	case BehaviorStalker:
		// Сталкеры наблюдают издалека
		c.CurrentState = "stalk"

	case BehaviorFleeing:
		// Убегающие всегда убегают
		c.CurrentState = "flee"

	case BehaviorPatrol:
		// Патрульные могут атаковать или продолжать патрулирование
		if rand.Float64() < 0.5 {
			c.CurrentState = "chase"
		}

	case BehaviorHunter:
		// Охотники всегда преследуют
		c.CurrentState = "chase"
	}

	c.StateTime = 0
}

// LoseTarget теряет игрока из виду
func (c *Creature) LoseTarget() {
	if c.PlayerTarget != nil {
		// Запоминаем последнюю позицию игрока
		c.TargetPos = c.PlayerTarget.Position
		c.CurrentState = "search"
		c.StateTime = 0
	}
}

// TakeDamage наносит урон существу
func (c *Creature) TakeDamage(amount float64) {
	c.Health -= amount

	// Реакция на получение урона
	if c.Health > 0 {
		switch c.BehaviorType {
		case BehaviorPassive, BehaviorFleeing:
			// Пассивные существа убегают
			c.CurrentState = "flee"

		case BehaviorAggressive, BehaviorHunter:
			// Агрессивные существа атакуют
			c.CurrentState = "chase"

		case BehaviorStalker:
			// Сталкеры могут атаковать или скрыться
			if rand.Float64() < 0.5 {
				c.CurrentState = "chase"
			} else {
				c.CurrentState = "flee"
			}

		case BehaviorPatrol:
			// Патрульные могут атаковать или отступить
			if rand.Float64() < 0.7 {
				c.CurrentState = "chase"
			} else {
				c.CurrentState = "flee"
			}
		}

		c.StateTime = 0
	}
}

// IsDead проверяет, мертво ли существо
func (c *Creature) IsDead() bool {
	return c.Health <= 0
}

// moveForward перемещает существо вперед
func (c *Creature) moveForward() {
	dx := c.Speed * math.Cos(c.Direction)
	dy := c.Speed * math.Sin(c.Direction)
	c.Position.X += dx
	c.Position.Y += dy
}

// getDirectionTo возвращает направление к точке
func (c *Creature) getDirectionTo(target Vector2D) float64 {
	dx := target.X - c.Position.X
	dy := target.Y - c.Position.Y
	return math.Atan2(dy, dx)
}

// distanceTo возвращает расстояние до точки
func (c *Creature) distanceTo(target Vector2D) float64 {
	dx := target.X - c.Position.X
	dy := target.Y - c.Position.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// smoothDirection плавно изменяет направление
func (c *Creature) smoothDirection(current, target, factor float64) float64 {
	// Вычисляем разницу между углами
	diff := target - current

	// Нормализуем разницу в пределах [-π, π]
	for diff > math.Pi {
		diff -= 2 * math.Pi
	}
	for diff < -math.Pi {
		diff += 2 * math.Pi
	}

	// Плавно изменяем направление
	newDir := current + diff*factor

	// Нормализуем результат в пределах [0, 2π]
	for newDir < 0 {
		newDir += 2 * math.Pi
	}
	for newDir >= 2*math.Pi {
		newDir -= 2 * math.Pi
	}

	return newDir
}

// updatePartAnimation обновляет анимацию части существа
func (c *Creature) updatePartAnimation(part *CreaturePart) {
	// В зависимости от текущего состояния выбираем анимацию
	var frames []int

	switch c.CurrentState {
	case "idle":
		frames = []int{0, 1, 2, 1} // Простая анимация дыхания

	case "wander", "search":
		frames = []int{3, 4, 5, 6, 5, 4} // Анимация ходьбы

	case "chase":
		frames = []int{7, 8, 9, 10, 9, 8} // Быстрая анимация бега

	case "attack":
		frames = []int{11, 12, 13, 14, 15} // Анимация атаки

	case "stalk":
		frames = []int{16, 17, 18, 17} // Анимация скрытного передвижения

	case "flee":
		frames = []int{19, 20, 21, 22, 21, 20} // Анимация бегства

	default:
		frames = []int{0}
	}

	// Обновляем текущий кадр анимации
	if len(frames) > 0 {
		frameIdx := (c.StateTime / 10) % len(frames)
		part.CurrentAnim = frames[frameIdx]
	}
}
