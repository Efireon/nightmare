package world

import (
	"math"

	"nightmare/internal/entity"
)

// CollisionSystem управляет коллизиями в игровом мире
type CollisionSystem struct {
	world        *World
	cellSize     float64
	collisionMap [][]bool
}

// CollisionResult представляет результат проверки коллизий
type CollisionResult struct {
	HasCollision bool
	Normal       entity.Vector2D
	Object       interface{}
	Distance     float64
}

// NewCollisionSystem создает новую систему коллизий
func NewCollisionSystem(world *World, cellSize float64) *CollisionSystem {
	// Создаем карту коллизий
	width := int(math.Ceil(float64(world.Width) / cellSize))
	height := int(math.Ceil(float64(world.Height) / cellSize))

	collisionMap := make([][]bool, height)
	for y := range collisionMap {
		collisionMap[y] = make([]bool, width)
	}

	return &CollisionSystem{
		world:        world,
		cellSize:     cellSize,
		collisionMap: collisionMap,
	}
}

// UpdateCollisionMap обновляет карту коллизий
func (cs *CollisionSystem) UpdateCollisionMap() {
	// Сбрасываем карту коллизий
	for y := range cs.collisionMap {
		for x := range cs.collisionMap[y] {
			cs.collisionMap[y][x] = false
		}
	}

	// Обновляем коллизии на основе тайлов
	for y := 0; y < cs.world.Height; y++ {
		for x := 0; x < cs.world.Width; x++ {
			tile := cs.world.GetTileAt(x, y)
			if tile == nil {
				continue
			}

			// Проверяем, является ли тайл непроходимым
			if cs.isTileSolid(tile) {
				// Вычисляем индексы в карте коллизий
				cellX := int(float64(x) / cs.cellSize)
				cellY := int(float64(y) / cs.cellSize)

				// Проверяем, что индексы в пределах карты
				if cellX >= 0 && cellX < len(cs.collisionMap[0]) && cellY >= 0 && cellY < len(cs.collisionMap) {
					cs.collisionMap[cellY][cellX] = true
				}
			}
		}
	}

	// Добавляем коллизии от объектов
	for _, obj := range cs.world.Objects {
		if obj.Solid {
			// Вычисляем индексы в карте коллизий
			cellX := int(obj.Position.X / cs.cellSize)
			cellY := int(obj.Position.Y / cs.cellSize)

			// Проверяем, что индексы в пределах карты
			if cellX >= 0 && cellX < len(cs.collisionMap[0]) && cellY >= 0 && cellY < len(cs.collisionMap) {
				cs.collisionMap[cellY][cellX] = true
			}
		}
	}
}

// CheckCollision проверяет коллизию в указанной позиции
func (cs *CollisionSystem) CheckCollision(position entity.Vector2D) bool {
	// Вычисляем индексы в карте коллизий
	cellX := int(position.X / cs.cellSize)
	cellY := int(position.Y / cs.cellSize)

	// Проверяем, что индексы в пределах карты
	if cellX < 0 || cellX >= len(cs.collisionMap[0]) || cellY < 0 || cellY >= len(cs.collisionMap) {
		return true // За пределами мира считаем коллизией
	}

	return cs.collisionMap[cellY][cellX]
}

// CheckMovement проверяет возможность движения из текущей позиции в новую
func (cs *CollisionSystem) CheckMovement(from, to entity.Vector2D) (entity.Vector2D, bool) {
	// Проверяем коллизию в новой позиции
	if cs.CheckCollision(to) {
		// Есть коллизия, движение невозможно
		return from, false
	}

	// Проверяем путь от from до to
	direction := entity.Vector2D{
		X: to.X - from.X,
		Y: to.Y - from.Y,
	}

	distance := math.Sqrt(direction.X*direction.X + direction.Y*direction.Y)

	// Если расстояние меньше порогового значения, считаем движение возможным
	if distance < cs.cellSize*0.1 {
		return to, true
	}

	// Нормализуем направление
	direction.X /= distance
	direction.Y /= distance

	// Проверяем несколько точек на пути
	steps := int(distance / (cs.cellSize * 0.5))
	if steps < 2 {
		steps = 2
	}

	stepSize := distance / float64(steps)

	for i := 1; i < steps; i++ {
		checkPoint := entity.Vector2D{
			X: from.X + direction.X*stepSize*float64(i),
			Y: from.Y + direction.Y*stepSize*float64(i),
		}

		if cs.CheckCollision(checkPoint) {
			// Нашли коллизию, движение до этой точки невозможно
			// Возвращаем последнюю безопасную позицию
			safePoint := entity.Vector2D{
				X: from.X + direction.X*stepSize*float64(i-1),
				Y: from.Y + direction.Y*stepSize*float64(i-1),
			}

			return safePoint, false
		}
	}

	// Движение возможно
	return to, true
}

// CheckMovementWithSliding проверяет возможность движения с учетом скольжения вдоль стен
func (cs *CollisionSystem) CheckMovementWithSliding(from, to entity.Vector2D) entity.Vector2D {
	// Проверяем коллизию напрямую
	newPos, canMove := cs.CheckMovement(from, to)
	if canMove {
		return to
	}

	// Если движение невозможно, пробуем движение только по X
	xMove := entity.Vector2D{X: to.X, Y: from.Y}
	newPos, canMoveX := cs.CheckMovement(from, xMove)

	// Если движение по X возможно, используем его
	if canMoveX {
		return newPos
	}

	// Пробуем движение только по Y
	yMove := entity.Vector2D{X: from.X, Y: to.Y}
	newPos, canMoveY := cs.CheckMovement(from, yMove)

	// Если движение по Y возможно, используем его
	if canMoveY {
		return newPos
	}

	// Если ни один из вариантов не работает, остаемся на месте
	return from
}

// CastRay выполняет рейкаст от точки в указанном направлении
func (cs *CollisionSystem) CastRay(start entity.Vector2D, direction entity.Vector2D, maxDistance float64) (CollisionResult, bool) {
	// Нормализуем направление
	length := math.Sqrt(direction.X*direction.X + direction.Y*direction.Y)
	if length > 0 {
		direction.X /= length
		direction.Y /= length
	} else {
		return CollisionResult{}, false
	}

	// Шаг рейкаста
	stepSize := cs.cellSize * 0.5

	// Проверяем точки на пути луча
	for dist := 0.0; dist <= maxDistance; dist += stepSize {
		checkPoint := entity.Vector2D{
			X: start.X + direction.X*dist,
			Y: start.Y + direction.Y*dist,
		}

		// Проверяем, находится ли точка в пределах мира
		if checkPoint.X < 0 || checkPoint.X >= float64(cs.world.Width) ||
			checkPoint.Y < 0 || checkPoint.Y >= float64(cs.world.Height) {
			// Достигли границы мира
			return CollisionResult{
				HasCollision: true,
				Normal:       entity.Vector2D{X: -direction.X, Y: -direction.Y},
				Object:       nil,
				Distance:     dist,
			}, true
		}

		// Проверяем коллизию в точке
		if cs.CheckCollision(checkPoint) {
			// Нашли коллизию

			// Определяем объект, с которым произошла коллизия
			var hitObject interface{}

			// Проверяем тайл
			tileX, tileY := int(checkPoint.X), int(checkPoint.Y)
			tile := cs.world.GetTileAt(tileX, tileY)
			if tile != nil && cs.isTileSolid(tile) {
				hitObject = tile
			}

			// Проверяем объекты
			for _, obj := range cs.world.Objects {
				if obj.Solid {
					objDist := distance(obj.Position, checkPoint)
					if objDist < 1.0 {
						hitObject = obj
						break
					}
				}
			}

			// Вычисляем нормаль поверхности (упрощенно)
			normal := entity.Vector2D{X: -direction.X, Y: -direction.Y}

			return CollisionResult{
				HasCollision: true,
				Normal:       normal,
				Object:       hitObject,
				Distance:     dist,
			}, true
		}
	}

	// Луч не пересек ни один объект
	return CollisionResult{HasCollision: false}, false
}

// CheckLineOfSight проверяет прямую видимость между двумя точками
func (cs *CollisionSystem) CheckLineOfSight(from, to entity.Vector2D) bool {
	// Вычисляем направление и расстояние
	direction := entity.Vector2D{
		X: to.X - from.X,
		Y: to.Y - from.Y,
	}

	distance := math.Sqrt(direction.X*direction.X + direction.Y*direction.Y)

	// Выполняем рейкаст
	result, hit := cs.CastRay(from, direction, distance)

	// Если не было коллизий, или расстояние до коллизии больше или равно расстоянию до цели,
	// то прямая видимость есть
	return !hit || result.Distance >= distance
}

// CheckCollisionWithEntities проверяет коллизию с существами
func (cs *CollisionSystem) CheckCollisionWithEntities(position entity.Vector2D, radius float64) (bool, *entity.Creature) {
	// Проверяем коллизию со всеми существами
	for _, e := range cs.world.Entities {
		creature, ok := e.(*entity.Creature)
		if !ok {
			continue
		}

		// Вычисляем расстояние между позицией и существом
		dist := distance(position, creature.Position)

		// Если расстояние меньше суммы радиусов, есть коллизия
		if dist < radius {
			return true, creature
		}
	}

	return false, nil
}

// isTileSolid проверяет, является ли тайл непроходимым
func (cs *CollisionSystem) isTileSolid(tile *Tile) bool {
	// Определяем, какие типы тайлов считаются непроходимыми
	switch tile.Type {
	case TileWater, TileSwamp, TileRocks, TileDenseForest:
		return true
	case TileCorrupted:
		// Искаженная местность может быть частично проходимой в зависимости от уровня искажения
		return tile.Corruption > 0.7
	default:
		return false
	}
}

// distance вычисляет расстояние между двумя точками
func distance(a, b entity.Vector2D) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return math.Sqrt(dx*dx + dy*dy)
}
