package world

import (
	"math"
	"math/rand"

	"github.com/ojrac/opensimplex-go"
)

// TileType представляет тип тайла
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

// Tile представляет тайл в мире
type Tile struct {
	Type       TileType
	Elevation  float64
	Moisture   float64
	Corruption float64 // Уровень "кошмарности" от 0 до 1
	Objects    []WorldObject
}

// WorldObject представляет объект в мире
type WorldObject struct {
	ID          int
	Type        string
	Position    Vector2D
	Solid       bool
	Interactive bool
}

// Entity представляет сущность в мире
type Entity struct {
	ID        int
	Type      string
	Position  Vector2D
	Direction float64
	Model     *EntityModel
	Behavior  EntityBehavior
}

// EntityModel представляет модель сущности
type EntityModel struct {
	Parts      []EntityPart
	Animations map[string][]int
}

// EntityPart представляет часть модели сущности
type EntityPart struct {
	Type     string
	Texture  int
	Offset   Vector2D
	Scale    float64
	Rotation float64
}

// EntityBehavior определяет поведение сущности
type EntityBehavior interface {
	Update(world *World, entity *Entity)
}

// Vector2D представляет 2D вектор
type Vector2D struct {
	X, Y float64
}

// World представляет игровой мир
type World struct {
	Width    int
	Height   int
	Tiles    [][]Tile
	Entities []*Entity
	Objects  []WorldObject
	nextID   int
	noise    opensimplex.Noise // Шумовой генератор для процедурной генерации
}

// NewWorld создает новый мир
func NewWorld(width, height int) (*World, error) {
	// Создаем мир
	world := &World{
		Width:    width,
		Height:   height,
		Tiles:    make([][]Tile, height),
		Entities: []*Entity{},
		Objects:  []WorldObject{},
		nextID:   1,
		noise:    opensimplex.New(rand.Int63()),
	}

	// Инициализируем тайлы
	for y := 0; y < height; y++ {
		world.Tiles[y] = make([]Tile, width)
		for x := 0; x < width; x++ {
			world.Tiles[y][x] = Tile{
				Type:       TileGrass,
				Elevation:  0,
				Moisture:   0,
				Corruption: 0,
				Objects:    []WorldObject{},
			}
		}
	}

	// Генерируем ландшафт
	world.generateTerrain()

	// Размещаем объекты
	world.placeObjects()

	return world, nil
}

// generateTerrain генерирует ландшафт с использованием шумовых алгоритмов
func (w *World) generateTerrain() {
	const elevationScale = 0.05
	const moistureScale = 0.07

	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			// Генерируем высоту с помощью шума Симплекса
			elevation := (w.noise.Eval2(float64(x)*elevationScale, float64(y)*elevationScale) + 1) / 2

			// Генерируем влажность с другой октавой шума
			moisture := (w.noise.Eval2(float64(x)*moistureScale+100, float64(y)*moistureScale+100) + 1) / 2

			// Сохраняем базовые значения
			w.Tiles[y][x].Elevation = elevation
			w.Tiles[y][x].Moisture = moisture

			// Определяем тип тайла на основе высоты и влажности
			w.Tiles[y][x].Type = w.determineTileType(elevation, moisture)
		}
	}
}

// determineTileType определяет тип тайла на основе высоты и влажности
func (w *World) determineTileType(elevation, moisture float64) TileType {
	if elevation < 0.3 {
		if moisture > 0.7 {
			return TileSwamp
		} else if moisture > 0.4 {
			return TileWater
		} else {
			return TileGrass
		}
	} else if elevation < 0.6 {
		if moisture > 0.6 {
			return TileForest
		} else if moisture > 0.3 {
			return TileGrass
		} else {
			return TilePath
		}
	} else {
		if moisture > 0.7 {
			return TileDenseForest
		} else if moisture > 0.4 {
			return TileForest
		} else {
			return TileRocks
		}
	}
}

// placeObjects размещает объекты в мире
func (w *World) placeObjects() {
	// Размещаем деревья в лесных областях
	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			tile := &w.Tiles[y][x]

			if tile.Type == TileForest {
				if rand.Float64() < 0.2 {
					tree := WorldObject{
						ID:          w.nextID,
						Type:        "tree",
						Position:    Vector2D{X: float64(x), Y: float64(y)},
						Solid:       true,
						Interactive: false,
					}
					tile.Objects = append(tile.Objects, tree)
					w.Objects = append(w.Objects, tree)
					w.nextID++
				}
			} else if tile.Type == TileDenseForest {
				if rand.Float64() < 0.5 {
					tree := WorldObject{
						ID:          w.nextID,
						Type:        "dense_tree",
						Position:    Vector2D{X: float64(x), Y: float64(y)},
						Solid:       true,
						Interactive: false,
					}
					tile.Objects = append(tile.Objects, tree)
					w.Objects = append(w.Objects, tree)
					w.nextID++
				}
			} else if tile.Type == TileRocks {
				if rand.Float64() < 0.1 {
					rock := WorldObject{
						ID:          w.nextID,
						Type:        "rock",
						Position:    Vector2D{X: float64(x), Y: float64(y)},
						Solid:       true,
						Interactive: false,
					}
					tile.Objects = append(tile.Objects, rock)
					w.Objects = append(w.Objects, rock)
					w.nextID++
				}
			}
		}
	}
}

// Update обновляет состояние мира
func (w *World) Update() {
	// Обновляем все сущности
	for _, entity := range w.Entities {
		if entity.Behavior != nil {
			entity.Behavior.Update(w, entity)
		}
	}
}

// GetTileAt возвращает тайл в указанной позиции
func (w *World) GetTileAt(x, y int) *Tile {
	if x < 0 || y < 0 || x >= w.Width || y >= w.Height {
		return nil
	}
	return &w.Tiles[y][x]
}

// SpawnCreature создает существо указанного типа в указанной позиции
func (w *World) SpawnCreature(creatureType string, position Vector2D) *Entity {
	// Создаем модель существа
	model := generateCreatureModel(creatureType)

	// Создаем сущность
	entity := &Entity{
		ID:        w.nextID,
		Type:      creatureType,
		Position:  position,
		Direction: rand.Float64() * 2 * math.Pi,
		Model:     model,
		Behavior:  NewBasicCreatureBehavior(),
	}

	w.Entities = append(w.Entities, entity)
	w.nextID++

	return entity
}

// ModifyEnvironment изменяет окружение вокруг указанной позиции
func (w *World) ModifyEnvironment(position Vector2D, intensity float64) {
	// Радиус влияния
	radius := 10.0 + intensity*20.0
	radiusSq := radius * radius

	// Изменяем тайлы вокруг позиции
	for y := int(position.Y - radius); y <= int(position.Y+radius); y++ {
		for x := int(position.X - radius); x <= int(position.X+radius); x++ {
			// Проверяем, что координаты в пределах мира
			if x < 0 || y < 0 || x >= w.Width || y >= w.Height {
				continue
			}

			// Вычисляем расстояние до центра
			dx := float64(x) - position.X
			dy := float64(y) - position.Y
			distSq := dx*dx + dy*dy

			// Если точка внутри радиуса
			if distSq <= radiusSq {
				// Вычисляем силу воздействия (ближе к центру - сильнее)
				strength := (1.0 - distSq/radiusSq) * intensity

				// Увеличиваем уровень "кошмарности"
				tile := &w.Tiles[y][x]
				tile.Corruption += strength * 0.3

				// Ограничиваем значение
				if tile.Corruption > 1.0 {
					tile.Corruption = 1.0
				}

				// При высоком уровне коррупции меняем тип тайла
				if tile.Corruption > 0.7 {
					tile.Type = TileCorrupted
				}
			}
		}
	}
}

// generateCreatureModel генерирует модель существа указанного типа
func generateCreatureModel(creatureType string) *EntityModel {
	model := &EntityModel{
		Parts:      []EntityPart{},
		Animations: make(map[string][]int),
	}

	switch creatureType {
	case "shadow":
		// Тело
		model.Parts = append(model.Parts, EntityPart{
			Type:     "body",
			Texture:  1, // ID текстуры
			Offset:   Vector2D{X: 0, Y: 0},
			Scale:    1.0,
			Rotation: 0,
		})

		// Руки
		model.Parts = append(model.Parts, EntityPart{
			Type:     "arm",
			Texture:  2,
			Offset:   Vector2D{X: -0.5, Y: 0},
			Scale:    0.7,
			Rotation: 0,
		})

		model.Parts = append(model.Parts, EntityPart{
			Type:     "arm",
			Texture:  2,
			Offset:   Vector2D{X: 0.5, Y: 0},
			Scale:    0.7,
			Rotation: 0,
		})

	case "spider":
		// Тело
		model.Parts = append(model.Parts, EntityPart{
			Type:     "body",
			Texture:  10,
			Offset:   Vector2D{X: 0, Y: 0},
			Scale:    1.0,
			Rotation: 0,
		})

		// Ноги (8 штук)
		for i := 0; i < 8; i++ {
			angle := float64(i) * math.Pi / 4
			model.Parts = append(model.Parts, EntityPart{
				Type:     "leg",
				Texture:  11,
				Offset:   Vector2D{X: math.Cos(angle) * 0.5, Y: math.Sin(angle) * 0.5},
				Scale:    0.6,
				Rotation: angle,
			})
		}

	// Другие типы существ...
	default:
		// Базовая модель для неизвестного типа
		model.Parts = append(model.Parts, EntityPart{
			Type:     "body",
			Texture:  0,
			Offset:   Vector2D{X: 0, Y: 0},
			Scale:    1.0,
			Rotation: 0,
		})
	}

	// Добавляем базовые анимации
	model.Animations["idle"] = []int{0, 1, 2, 1}
	model.Animations["walk"] = []int{3, 4, 5, 6}
	model.Animations["attack"] = []int{7, 8, 9}

	return model
}

// BasicCreatureBehavior представляет базовое поведение существа
type BasicCreatureBehavior struct {
	state          string
	targetPoint    Vector2D
	idleTime       int
	attackRange    float64
	detectionRange float64
}

// NewBasicCreatureBehavior создает новое базовое поведение существа
func NewBasicCreatureBehavior() *BasicCreatureBehavior {
	return &BasicCreatureBehavior{
		state:          "wander",
		targetPoint:    Vector2D{X: 0, Y: 0},
		idleTime:       0,
		attackRange:    5.0,
		detectionRange: 20.0,
	}
}

// Update обновляет поведение существа
func (b *BasicCreatureBehavior) Update(world *World, entity *Entity) {
	// Здесь будет логика поведения существа
	// ...
}
