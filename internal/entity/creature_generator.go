package entity

import (
	"math"
	"math/rand"

	"nightmare/internal/util"
)

// CreatureGenerator отвечает за процедурную генерацию существ
type CreatureGenerator struct {
	noise        *util.NoiseGenerator
	nextID       int
	textureAtlas []int // IDs для доступных текстур
}

// NewCreatureGenerator создает новый генератор существ
func NewCreatureGenerator() *CreatureGenerator {
	return &CreatureGenerator{
		noise:        util.NewNoiseGenerator(rand.Int63()),
		nextID:       1,
		textureAtlas: make([]int, 0),
	}
}

// RegisterTexture регистрирует текстуру для использования
func (g *CreatureGenerator) RegisterTexture(textureID int) {
	g.textureAtlas = append(g.textureAtlas, textureID)
}

// GenerateCreature создает новое существо указанного типа
func (g *CreatureGenerator) GenerateCreature(creatureType string, position Vector2D) *Creature {
	creature := NewCreature(g.nextID, creatureType, position)
	g.nextID++

	// Генерируем части тела существа
	g.generateParts(creature)

	return creature
}

// generateParts генерирует части тела существа
func (g *CreatureGenerator) generateParts(creature *Creature) {
	switch creature.Type {
	case "shadow":
		g.generateShadowParts(creature)
	case "spider":
		g.generateSpiderParts(creature)
	case "phantom":
		g.generatePhantomParts(creature)
	case "wendigo":
		g.generateWendigoParts(creature)
	case "faceless":
		g.generateFacelessParts(creature)
	default:
		g.generateGenericParts(creature)
	}
}

// generateShadowParts генерирует части тела для тени
func (g *CreatureGenerator) generateShadowParts(creature *Creature) {
	// Тело - основная форма
	bodyTexture := g.getRandomTextureID()
	creature.Parts = append(creature.Parts, CreaturePart{
		Type:        "body",
		TextureID:   bodyTexture,
		Position:    Vector2D{X: 0, Y: 0},
		Rotation:    0,
		Scale:       1.0 + rand.Float64()*0.3,
		AnimFrames:  []int{0, 1, 2, 3},
		CurrentAnim: 0,
	})

	// Конечности (случайное количество от 2 до 5)
	numLimbs := 2 + rand.Intn(4)
	for i := 0; i < numLimbs; i++ {
		angle := float64(i) * (2 * math.Pi / float64(numLimbs))
		dist := 0.5 + rand.Float64()*0.3
		limbTexture := g.getRandomTextureID()

		creature.Parts = append(creature.Parts, CreaturePart{
			Type:        "limb",
			TextureID:   limbTexture,
			Position:    Vector2D{X: math.Cos(angle) * dist, Y: math.Sin(angle) * dist},
			Rotation:    angle,
			Scale:       0.5 + rand.Float64()*0.5,
			AnimFrames:  []int{0, 1, 2, 3},
			CurrentAnim: 0,
		})
	}

	// Лицо (опционально)
	if rand.Float64() < 0.7 {
		faceTexture := g.getRandomTextureID()
		creature.Parts = append(creature.Parts, CreaturePart{
			Type:        "face",
			TextureID:   faceTexture,
			Position:    Vector2D{X: 0, Y: 0},
			Rotation:    0,
			Scale:       0.7 + rand.Float64()*0.3,
			AnimFrames:  []int{0, 1, 2, 3},
			CurrentAnim: 0,
		})
	}
}

// generateSpiderParts генерирует части тела для паука
func (g *CreatureGenerator) generateSpiderParts(creature *Creature) {
	// Тело
	bodyTexture := g.getRandomTextureID()
	creature.Parts = append(creature.Parts, CreaturePart{
		Type:        "body",
		TextureID:   bodyTexture,
		Position:    Vector2D{X: 0, Y: 0},
		Rotation:    0,
		Scale:       0.8 + rand.Float64()*0.4,
		AnimFrames:  []int{0, 1, 2, 3},
		CurrentAnim: 0,
	})

	// Голова
	headTexture := g.getRandomTextureID()
	headOffset := Vector2D{X: 0.4 + rand.Float64()*0.2, Y: 0}
	creature.Parts = append(creature.Parts, CreaturePart{
		Type:        "head",
		TextureID:   headTexture,
		Position:    headOffset,
		Rotation:    0,
		Scale:       0.5 + rand.Float64()*0.3,
		AnimFrames:  []int{0, 1, 2, 3},
		CurrentAnim: 0,
	})

	// Ноги (8 штук)
	for i := 0; i < 8; i++ {
		angle := float64(i) * (2 * math.Pi / 8)
		dist := 0.5 + rand.Float64()*0.2
		legTexture := g.getRandomTextureID()

		scale := 0.7 + rand.Float64()*0.4
		// Более длинные передние ноги
		if i < 2 {
			scale *= 1.3
		}

		creature.Parts = append(creature.Parts, CreaturePart{
			Type:        "leg",
			TextureID:   legTexture,
			Position:    Vector2D{X: math.Cos(angle) * dist, Y: math.Sin(angle) * dist},
			Rotation:    angle,
			Scale:       scale,
			AnimFrames:  []int{0, 1, 2, 3},
			CurrentAnim: 0,
		})
	}
}

// generatePhantomParts генерирует части тела для призрака
func (g *CreatureGenerator) generatePhantomParts(creature *Creature) {
	// Основное тело (полупрозрачное)
	bodyTexture := g.getRandomTextureID()
	creature.Parts = append(creature.Parts, CreaturePart{
		Type:        "body",
		TextureID:   bodyTexture,
		Position:    Vector2D{X: 0, Y: 0},
		Rotation:    0,
		Scale:       1.2 + rand.Float64()*0.6,
		AnimFrames:  []int{0, 1, 2, 3},
		CurrentAnim: 0,
	})

	// "Хвост" или нижняя часть
	tailTexture := g.getRandomTextureID()
	creature.Parts = append(creature.Parts, CreaturePart{
		Type:        "tail",
		TextureID:   tailTexture,
		Position:    Vector2D{X: 0, Y: 0.5 + rand.Float64()*0.3},
		Rotation:    0,
		Scale:       0.8 + rand.Float64()*0.4,
		AnimFrames:  []int{0, 1, 2, 3},
		CurrentAnim: 0,
	})

	// Руки (2-4 штуки)
	numArms := 2 + rand.Intn(3)
	for i := 0; i < numArms; i++ {
		angle := float64(i) * (2 * math.Pi / float64(numArms))
		dist := 0.4 + rand.Float64()*0.3
		armTexture := g.getRandomTextureID()

		creature.Parts = append(creature.Parts, CreaturePart{
			Type:        "arm",
			TextureID:   armTexture,
			Position:    Vector2D{X: math.Cos(angle) * dist, Y: math.Sin(angle)*dist - 0.2},
			Rotation:    angle,
			Scale:       0.6 + rand.Float64()*0.5,
			AnimFrames:  []int{0, 1, 2, 3},
			CurrentAnim: 0,
		})
	}

	// Лицо или маска
	faceTexture := g.getRandomTextureID()
	creature.Parts = append(creature.Parts, CreaturePart{
		Type:        "face",
		TextureID:   faceTexture,
		Position:    Vector2D{X: 0, Y: -0.3 - rand.Float64()*0.1},
		Rotation:    0,
		Scale:       0.6 + rand.Float64()*0.3,
		AnimFrames:  []int{0, 1, 2, 3},
		CurrentAnim: 0,
	})
}

// generateWendigoParts генерирует части тела для вендиго
func (g *CreatureGenerator) generateWendigoParts(creature *Creature) {
	// Тело
	bodyTexture := g.getRandomTextureID()
	creature.Parts = append(creature.Parts, CreaturePart{
		Type:        "body",
		TextureID:   bodyTexture,
		Position:    Vector2D{X: 0, Y: 0},
		Rotation:    0,
		Scale:       1.0 + rand.Float64()*0.4,
		AnimFrames:  []int{0, 1, 2, 3},
		CurrentAnim: 0,
	})

	// Голова (с рогами)
	headTexture := g.getRandomTextureID()
	creature.Parts = append(creature.Parts, CreaturePart{
		Type:        "head",
		TextureID:   headTexture,
		Position:    Vector2D{X: 0, Y: -0.7 - rand.Float64()*0.2},
		Rotation:    0,
		Scale:       0.9 + rand.Float64()*0.3,
		AnimFrames:  []int{0, 1, 2, 3},
		CurrentAnim: 0,
	})

	// Рога
	for i := 0; i < 2; i++ {
		angle := (float64(i)*2 - 1) * (math.Pi / 8)
		hornTexture := g.getRandomTextureID()

		creature.Parts = append(creature.Parts, CreaturePart{
			Type:        "horn",
			TextureID:   hornTexture,
			Position:    Vector2D{X: math.Cos(angle) * 0.4, Y: -0.9 - rand.Float64()*0.3},
			Rotation:    angle - math.Pi/2,
			Scale:       0.7 + rand.Float64()*0.5,
			AnimFrames:  []int{0, 1, 2, 3},
			CurrentAnim: 0,
		})
	}

	// Руки (длинные)
	for i := 0; i < 2; i++ {
		side := float64(i*2 - 1)
		armTexture := g.getRandomTextureID()

		creature.Parts = append(creature.Parts, CreaturePart{
			Type:        "arm",
			TextureID:   armTexture,
			Position:    Vector2D{X: side * (0.5 + rand.Float64()*0.1), Y: -0.2},
			Rotation:    side * math.Pi / 8,
			Scale:       1.2 + rand.Float64()*0.6,
			AnimFrames:  []int{0, 1, 2, 3},
			CurrentAnim: 0,
		})
	}

	// Ноги
	for i := 0; i < 2; i++ {
		side := float64(i*2 - 1)
		legTexture := g.getRandomTextureID()

		creature.Parts = append(creature.Parts, CreaturePart{
			Type:        "leg",
			TextureID:   legTexture,
			Position:    Vector2D{X: side * (0.3 + rand.Float64()*0.1), Y: 0.6},
			Rotation:    side * math.Pi / 10,
			Scale:       1.0 + rand.Float64()*0.4,
			AnimFrames:  []int{0, 1, 2, 3},
			CurrentAnim: 0,
		})
	}
}

// generateFacelessParts генерирует части тела для безликого
func (g *CreatureGenerator) generateFacelessParts(creature *Creature) {
	// Тело (высокое и тонкое)
	bodyTexture := g.getRandomTextureID()
	creature.Parts = append(creature.Parts, CreaturePart{
		Type:        "body",
		TextureID:   bodyTexture,
		Position:    Vector2D{X: 0, Y: 0},
		Rotation:    0,
		Scale:       1.5 + rand.Float64()*0.5,
		AnimFrames:  []int{0, 1, 2, 3},
		CurrentAnim: 0,
	})

	// Голова (без лица)
	headTexture := g.getRandomTextureID()
	creature.Parts = append(creature.Parts, CreaturePart{
		Type:        "head",
		TextureID:   headTexture,
		Position:    Vector2D{X: 0, Y: -0.8 - rand.Float64()*0.2},
		Rotation:    0,
		Scale:       0.7 + rand.Float64()*0.2,
		AnimFrames:  []int{0, 1, 2, 3},
		CurrentAnim: 0,
	})

	// Руки (несколько пар, длинные)
	numArmPairs := 2 + rand.Intn(2)
	for i := 0; i < numArmPairs; i++ {
		for j := 0; j < 2; j++ {
			side := float64(j*2 - 1)
			armTexture := g.getRandomTextureID()

			yOffset := -0.4 + float64(i)*0.4

			creature.Parts = append(creature.Parts, CreaturePart{
				Type:        "arm",
				TextureID:   armTexture,
				Position:    Vector2D{X: side * (0.4 + rand.Float64()*0.2), Y: yOffset},
				Rotation:    side * (math.Pi/4 + rand.Float64()*math.Pi/8),
				Scale:       1.3 + rand.Float64()*0.7,
				AnimFrames:  []int{0, 1, 2, 3},
				CurrentAnim: 0,
			})
		}
	}
}

// generateGenericParts генерирует части тела для неизвестного типа существа
func (g *CreatureGenerator) generateGenericParts(creature *Creature) {
	// Используем шум для определения формы
	seed := rand.Float64() * 100
	complexity := 0.5 + rand.Float64()*0.5

	// Основное тело
	bodyTexture := g.getRandomTextureID()
	creature.Parts = append(creature.Parts, CreaturePart{
		Type:        "body",
		TextureID:   bodyTexture,
		Position:    Vector2D{X: 0, Y: 0},
		Rotation:    0,
		Scale:       1.0 + rand.Float64()*0.5,
		AnimFrames:  []int{0, 1, 2, 3},
		CurrentAnim: 0,
	})

	// Добавляем случайные выступы
	numProtrusions := 3 + rand.Intn(7)
	for i := 0; i < numProtrusions; i++ {
		// Используем шум для определения положения
		angle := float64(i) * (2 * math.Pi / float64(numProtrusions))
		noise := g.noise.Perlin2D(seed+float64(i), seed+10, complexity)
		dist := 0.3 + rand.Float64()*0.4 + noise*0.3

		protrusionTexture := g.getRandomTextureID()

		// Определяем тип выступа
		partTypes := []string{"limb", "tentacle", "spike", "bulb"}
		partType := partTypes[rand.Intn(len(partTypes))]

		creature.Parts = append(creature.Parts, CreaturePart{
			Type:        partType,
			TextureID:   protrusionTexture,
			Position:    Vector2D{X: math.Cos(angle) * dist, Y: math.Sin(angle) * dist},
			Rotation:    angle,
			Scale:       0.4 + rand.Float64()*0.6 + noise*0.3,
			AnimFrames:  []int{0, 1, 2, 3},
			CurrentAnim: 0,
		})
	}

	// Случайно добавляем "голову" или "глаза"
	if rand.Float64() < 0.7 {
		eyeTexture := g.getRandomTextureID()
		eyeCount := 1 + rand.Intn(4)

		for i := 0; i < eyeCount; i++ {
			angle := float64(i) * (2 * math.Pi / float64(eyeCount))
			eyeDist := 0.2 + rand.Float64()*0.3

			creature.Parts = append(creature.Parts, CreaturePart{
				Type:        "eye",
				TextureID:   eyeTexture,
				Position:    Vector2D{X: math.Cos(angle) * eyeDist, Y: math.Sin(angle) * eyeDist},
				Rotation:    0,
				Scale:       0.2 + rand.Float64()*0.3,
				AnimFrames:  []int{0, 1, 2, 3},
				CurrentAnim: 0,
			})
		}
	}
}

// GenerateRandomCreature создает случайное существо
func (g *CreatureGenerator) GenerateRandomCreature(position Vector2D) *Creature {
	// Типы существ
	creatureTypes := []string{
		"shadow", "spider", "phantom", "wendigo", "faceless",
	}

	// Выбираем случайный тип
	creatureType := creatureTypes[rand.Intn(len(creatureTypes))]

	// Иногда создаем полностью случайное существо
	if rand.Float64() < 0.2 {
		creatureType = "random"
	}

	return g.GenerateCreature(creatureType, position)
}

// getRandomTextureID возвращает случайный ID текстуры
func (g *CreatureGenerator) getRandomTextureID() int {
	if len(g.textureAtlas) == 0 {
		return 0
	}
	return g.textureAtlas[rand.Intn(len(g.textureAtlas))]
}
