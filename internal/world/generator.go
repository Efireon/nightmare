package world

import (
	"math"
	"math/rand"

	"nightmare/internal/ai"
	"nightmare/internal/entity"
	"nightmare/internal/util"
)

// ThemeType представляет тип темы мира
type ThemeType int

const (
	ThemeDecay      ThemeType = iota // Разложение, гниение
	ThemeBlood                       // Кровь, внутренности
	ThemeShadow                      // Тени, тьма
	ThemeChild                       // Детская, школа, игрушки
	ThemeHospital                    // Больница, сумасшедший дом
	ThemeAncient                     // Древние руины, ритуалы
	ThemeIndustrial                  // Завод, механизмы
	ThemeForest                      // Лес, природа
	ThemeVoid                        // Пустота, космос
)

// ZoneType представляет тип зоны
type ZoneType int

const (
	ZoneSafe        ZoneType = iota // Безопасная зона
	ZoneExploration                 // Зона исследования
	ZoneDanger                      // Зона опасности
	ZoneNightmare                   // Зона кошмара
	ZoneTransition                  // Переходная зона
)

// Generator отвечает за процедурную генерацию мира
type Generator struct {
	world  *World
	random *util.RandomGenerator
	noise  *util.NoiseGenerator

	width  int
	height int

	mainTheme    ThemeType
	activeThemes map[ThemeType]float64

	zones     []Zone
	roomCount int

	creatureFactory *entity.CreatureGenerator

	observer      *ai.ObserverSystem
	fearInfluence map[ai.FearType]ThemeType
}

// Zone представляет зону в мире
type Zone struct {
	Type        ZoneType
	Position    entity.Vector2D
	Radius      float64
	Density     float64
	Theme       ThemeType
	Corruption  float64
	Connections []int // Индексы зон, с которыми соединена
}

// NewGenerator создает новый генератор мира
func NewGenerator(world *World, width, height int) *Generator {
	return &Generator{
		world:  world,
		random: util.NewRandomGenerator(rand.Int63()),
		noise:  util.NewNoiseGenerator(rand.Int63()),

		width:  width,
		height: height,

		mainTheme:    ThemeForest,
		activeThemes: make(map[ThemeType]float64),

		zones:     []Zone{},
		roomCount: 0,

		creatureFactory: entity.NewCreatureGenerator(),

		observer: nil,
		fearInfluence: map[ai.FearType]ThemeType{
			ai.FearDarkness:       ThemeShadow,
			ai.FearCreatures:      ThemeBlood,
			ai.FearSuddenNoises:   ThemeIndustrial,
			ai.FearIsolation:      ThemeVoid,
			ai.FearChasing:        ThemeForest,
			ai.FearGore:           ThemeBlood,
			ai.FearClaustrophobia: ThemeIndustrial,
			ai.FearOpenSpaces:     ThemeVoid,
			ai.FearUnknown:        ThemeAncient,
		},
	}
}

// SetObserver устанавливает систему наблюдения
func (g *Generator) SetObserver(observer *ai.ObserverSystem) {
	g.observer = observer
}

// GenerateWorld генерирует новый мир
func (g *Generator) GenerateWorld() error {
	// Инициализируем мир
	g.world.Width = g.width
	g.world.Height = g.height
	g.world.Tiles = make([][]Tile, g.height)
	for y := range g.world.Tiles {
		g.world.Tiles[y] = make([]Tile, g.width)
	}

	// Выбираем основную тему
	g.selectMainTheme()

	// Генерируем базовый ландшафт
	g.generateBaseTerrain()

	// Генерируем зоны
	g.generateZones()

	// Генерируем объекты
	g.generateObjects()

	// Генерируем существ
	g.generateCreatures()

	return nil
}

// selectMainTheme выбирает основную тему мира
func (g *Generator) selectMainTheme() {
	// По умолчанию используем тему леса
	g.mainTheme = ThemeForest

	// Если есть система наблюдения, учитываем страхи игрока
	if g.observer != nil {
		// Получаем доминирующий страх
		dominantFear := g.observer.GetDominantFear()

		// Сопоставляем страх с темой
		if theme, ok := g.fearInfluence[dominantFear]; ok {
			g.mainTheme = theme
		}
	}

	// Инициализируем активные темы
	g.activeThemes[g.mainTheme] = 1.0

	// Добавляем случайные второстепенные темы
	secondaryThemeCount := 2
	for i := 0; i < secondaryThemeCount; i++ {
		theme := ThemeType(g.random.RangeInt(0, int(ThemeVoid)+1))
		if theme != g.mainTheme {
			g.activeThemes[theme] = 0.3 + g.random.Float64()*0.3 // 30-60% влияния
		}
	}
}

// generateBaseTerrain генерирует базовый ландшафт
func (g *Generator) generateBaseTerrain() {
	// Генерируем базовые слои шума для разных характеристик местности
	elevationMap := g.generateNoiseMap(0.03, 0.7)
	moistureMap := g.generateNoiseMap(0.04, 0.5)
	densityMap := g.generateNoiseMap(0.05, 0.6)

	// Генерируем карту коррупции (влияние темы)
	corruptionMap := g.generateCorruptionMap()

	// Заполняем мир тайлами на основе сгенерированных карт
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			// Вычисляем характеристики местности
			elevation := elevationMap[y][x]
			moisture := moistureMap[y][x]
			density := densityMap[y][x]
			corruption := corruptionMap[y][x]

			// Определяем тип местности
			tileType := g.determineTileType(elevation, moisture, density, corruption)

			// Создаем тайл
			g.world.Tiles[y][x] = Tile{
				Type:       tileType,
				Elevation:  elevation,
				Moisture:   moisture,
				Corruption: corruption,
				Objects:    []WorldObject{},
			}
		}
	}
}

// generateNoiseMap генерирует карту шума
func (g *Generator) generateNoiseMap(scale, persistence float64) [][]float64 {
	noiseMap := make([][]float64, g.height)
	for y := range noiseMap {
		noiseMap[y] = make([]float64, g.width)
		for x := range noiseMap[y] {
			// Генерируем шум Перлина
			nx := float64(x) * scale
			ny := float64(y) * scale
			noiseMap[y][x] = g.noise.Perlin2D(nx, ny, persistence)
		}
	}
	return noiseMap
}

// generateCorruptionMap генерирует карту коррупции
func (g *Generator) generateCorruptionMap() [][]float64 {
	corruptionMap := make([][]float64, g.height)
	for y := range corruptionMap {
		corruptionMap[y] = make([]float64, g.width)
	}

	// Создаем базовую карту коррупции
	baseMap := g.generateNoiseMap(0.02, 0.8)

	// Добавляем очаги коррупции
	corruptionSpots := 5 + g.random.RangeInt(0, 10)
	for i := 0; i < corruptionSpots; i++ {
		// Случайная позиция
		cx := g.random.RangeInt(0, g.width)
		cy := g.random.RangeInt(0, g.height)

		// Случайный радиус
		radius := 10.0 + g.random.Float64()*30.0

		// Случайная интенсивность
		intensity := 0.6 + g.random.Float64()*0.4

		// Добавляем пятно коррупции
		g.applyCorruptionSpot(corruptionMap, cx, cy, radius, intensity)
	}

	// Смешиваем с базовой картой
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			corruptionMap[y][x] = math.Min(1.0, corruptionMap[y][x]+baseMap[y][x]*0.3)
		}
	}

	return corruptionMap
}

// applyCorruptionSpot добавляет пятно коррупции
func (g *Generator) applyCorruptionSpot(corruptionMap [][]float64, cx, cy int, radius, intensity float64) {
	radiusSq := radius * radius

	for y := int(math.Max(0, float64(cy)-radius)); y <= int(math.Min(float64(g.height-1), float64(cy)+radius)); y++ {
		for x := int(math.Max(0, float64(cx)-radius)); x <= int(math.Min(float64(g.width-1), float64(cx)+radius)); x++ {
			// Вычисляем расстояние до центра
			dx := float64(x - cx)
			dy := float64(y - cy)
			distSq := dx*dx + dy*dy

			// Если точка в пределах радиуса
			if distSq <= radiusSq {
				// Вычисляем силу воздействия (ближе к центру - сильнее)
				factor := (1.0 - math.Sqrt(distSq)/radius) * intensity

				// Применяем коррупцию
				corruptionMap[y][x] = math.Min(1.0, corruptionMap[y][x]+factor)
			}
		}
	}
}

// determineTileType определяет тип тайла на основе характеристик местности
func (g *Generator) determineTileType(elevation, moisture, density, corruption float64) TileType {
	// Если уровень коррупции высокий, используем искаженную местность
	if corruption > 0.7 {
		return TileCorrupted
	}

	// Вода в низинах с высокой влажностью
	if elevation < 0.3 && moisture > 0.6 {
		if moisture > 0.8 {
			return TileWater
		}
		return TileSwamp
	}

	// Равнины и леса
	if elevation < 0.6 {
		if density > 0.7 {
			return TileDenseForest
		}
		if density > 0.4 {
			return TileForest
		}
		return TileGrass
	}

	// Скалы и горы
	if elevation > 0.8 {
		return TileRocks
	}

	// Тропинки
	if 0.45 < density && density < 0.55 {
		return TilePath
	}

	// По умолчанию - трава
	return TileGrass
}

// generateZones генерирует зоны мира
func (g *Generator) generateZones() {
	// Очищаем существующие зоны
	g.zones = []Zone{}

	// Определяем количество зон
	zoneCount := 5 + g.random.RangeInt(0, 10)

	// Генерируем безопасную зону в центре
	g.zones = append(g.zones, Zone{
		Type:        ZoneSafe,
		Position:    entity.Vector2D{X: float64(g.width) / 2, Y: float64(g.height) / 2},
		Radius:      15.0 + g.random.Float64()*10.0,
		Density:     0.2 + g.random.Float64()*0.3,
		Theme:       g.mainTheme,
		Corruption:  0.1,
		Connections: []int{},
	})

	// Генерируем остальные зоны
	for i := 1; i < zoneCount; i++ {
		// Выбираем случайную позицию
		x := g.random.Float64() * float64(g.width)
		y := g.random.Float64() * float64(g.height)

		// Выбираем тип зоны (чем дальше от центра, тем опаснее)
		distFromCenter := math.Sqrt(math.Pow(x-float64(g.width)/2, 2) + math.Pow(y-float64(g.height)/2, 2))
		normalizedDist := distFromCenter / (math.Sqrt(math.Pow(float64(g.width)/2, 2) + math.Pow(float64(g.height)/2, 2)))

		zoneType := ZoneExploration
		if normalizedDist > 0.8 {
			zoneType = ZoneNightmare
		} else if normalizedDist > 0.6 {
			zoneType = ZoneDanger
		} else if normalizedDist > 0.4 {
			zoneType = ZoneExploration
		} else {
			zoneType = ZoneTransition
		}

		// Выбираем тему зоны
		theme := g.selectZoneTheme(zoneType)

		// Определяем уровень коррупции
		corruption := normalizedDist*0.7 + g.random.Float64()*0.3

		// Создаем зону
		zone := Zone{
			Type:        zoneType,
			Position:    entity.Vector2D{X: x, Y: y},
			Radius:      10.0 + g.random.Float64()*20.0,
			Density:     0.3 + g.random.Float64()*0.5,
			Theme:       theme,
			Corruption:  corruption,
			Connections: []int{},
		}

		// Проверяем на перекрытие с существующими зонами
		overlaps := false
		for j, existingZone := range g.zones {
			dist := distance(zone.Position, existingZone.Position)
			if dist < zone.Radius+existingZone.Radius {
				overlaps = true

				// Добавляем связь между зонами
				zone.Connections = append(zone.Connections, j)
				g.zones[j].Connections = append(g.zones[j].Connections, i)

				break
			}
		}

		// Если зона не перекрывается, добавляем ее
		if !overlaps {
			// Находим ближайшую зону и добавляем связь
			closestIdx := -1
			closestDist := math.MaxFloat64

			for j, existingZone := range g.zones {
				dist := distance(zone.Position, existingZone.Position)
				if dist < closestDist {
					closestDist = dist
					closestIdx = j
				}
			}

			if closestIdx >= 0 {
				zone.Connections = append(zone.Connections, closestIdx)
				g.zones[closestIdx].Connections = append(g.zones[closestIdx].Connections, i)
			}
		}

		g.zones = append(g.zones, zone)
	}

	// Применяем зоны к миру
	g.applyZonesToWorld()
}

// selectZoneTheme выбирает тему для зоны
func (g *Generator) selectZoneTheme(zoneType ZoneType) ThemeType {
	// Для безопасной зоны используем основную тему
	if zoneType == ZoneSafe {
		return g.mainTheme
	}

	// Для остальных зон выбираем случайную активную тему
	activeThemes := []ThemeType{}
	weights := []float64{}

	for theme, weight := range g.activeThemes {
		activeThemes = append(activeThemes, theme)
		weights = append(weights, weight)
	}

	// Если нет активных тем, используем основную
	if len(activeThemes) == 0 {
		return g.mainTheme
	}

	// Выбираем тему с учетом весов
	idx := g.random.WeightedChoiceIndex(weights)
	return activeThemes[idx]
}

// applyZonesToWorld применяет зоны к миру
func (g *Generator) applyZonesToWorld() {
	// Применяем характеристики зон к тайлам
	for _, zone := range g.zones {
		radiusSq := zone.Radius * zone.Radius

		for y := int(math.Max(0, zone.Position.Y-zone.Radius)); y <= int(math.Min(float64(g.height-1), zone.Position.Y+zone.Radius)); y++ {
			for x := int(math.Max(0, zone.Position.X-zone.Radius)); x <= int(math.Min(float64(g.width-1), zone.Position.X+zone.Radius)); x++ {
				// Вычисляем расстояние до центра зоны
				dx := float64(x) - zone.Position.X
				dy := float64(y) - zone.Position.Y
				distSq := dx*dx + dy*dy

				// Если точка в пределах зоны
				if distSq <= radiusSq {
					// Вычисляем фактор влияния (ближе к центру - сильнее)
					factor := (1.0 - math.Sqrt(distSq)/zone.Radius)

					// Модифицируем тайл в зависимости от зоны
					tile := &g.world.Tiles[y][x]

					// Изменяем уровень коррупции
					tile.Corruption = math.Min(1.0, tile.Corruption+zone.Corruption*factor*0.5)

					// Изменяем тип тайла в зависимости от темы зоны
					g.applyThemeToTile(tile, zone.Theme, factor)
				}
			}
		}
	}

	// Создаем пути между зонами
	g.generatePathsBetweenZones()
}

// applyThemeToTile применяет тему к тайлу
func (g *Generator) applyThemeToTile(tile *Tile, theme ThemeType, factor float64) {
	// Если фактор влияния слишком мал, выходим
	if factor < 0.2 {
		return
	}

	// Применяем тему в зависимости от ее типа
	switch theme {
	case ThemeDecay:
		// Увеличиваем влажность, создаем болота
		tile.Moisture = math.Min(1.0, tile.Moisture+factor*0.3)
		if factor > 0.5 && tile.Type == TileGrass {
			tile.Type = TileSwamp
		}

	case ThemeBlood:
		// Увеличиваем коррупцию
		tile.Corruption = math.Min(1.0, tile.Corruption+factor*0.4)

	case ThemeShadow:
		// Увеличиваем плотность леса
		if tile.Type == TileGrass && factor > 0.7 {
			tile.Type = TileForest
		} else if tile.Type == TileForest && factor > 0.8 {
			tile.Type = TileDenseForest
		}

	case ThemeChild:
		// Открытые пространства с тропинками
		if tile.Type != TileWater && tile.Type != TileSwamp && factor > 0.6 {
			if g.random.Float64() < 0.2 {
				tile.Type = TilePath
			} else {
				tile.Type = TileGrass
			}
		}

	case ThemeHospital:
		// Регулярная структура
		if tile.Type != TileWater && factor > 0.7 {
			x := int(tile.Position.X)
			y := int(tile.Position.Y)
			if (x+y)%5 == 0 {
				tile.Type = TilePath
			}
		}

	case ThemeAncient:
		// Каменистая местность
		if factor > 0.7 && g.random.Float64() < 0.3 {
			tile.Type = TileRocks
		}

	case ThemeIndustrial:
		// Открытые пространства с тропинками
		if tile.Type != TileWater && tile.Type != TileSwamp && factor > 0.6 {
			if g.random.Float64() < 0.3 {
				tile.Type = TilePath
			}
		}

	case ThemeForest:
		// Увеличиваем плотность леса
		if tile.Type == TileGrass && factor > 0.4 {
			tile.Type = TileForest
		} else if tile.Type == TileForest && factor > 0.7 {
			tile.Type = TileDenseForest
		}

	case ThemeVoid:
		// Пустые пространства
		if tile.Type != TileWater && tile.Type != TileSwamp && factor > 0.6 {
			tile.Type = TileGrass
		}
	}

	// Если уровень коррупции высокий, используем искаженную местность
	if tile.Corruption > 0.7 {
		tile.Type = TileCorrupted
	}
}

// generatePathsBetweenZones генерирует пути между зонами
func (g *Generator) generatePathsBetweenZones() {
	// Проходим по всем зонам
	for _, zone := range g.zones {
		// Проходим по всем связям зоны
		for _, connectedIdx := range zone.Connections {
			connectedZone := g.zones[connectedIdx]

			// Генерируем путь между зонами
			g.generatePath(zone.Position, connectedZone.Position)
		}
	}
}

// generatePath генерирует путь между двумя точками
func (g *Generator) generatePath(start, end entity.Vector2D) {
	// Шаг пути
	stepSize := 1.0

	// Вычисляем направление
	dx := end.X - start.X
	dy := end.Y - start.Y
	distance := math.Sqrt(dx*dx + dy*dy)

	// Нормализуем направление
	dx /= distance
	dy /= distance

	// Добавляем случайность к пути
	noise := g.noise

	// Проходим по пути
	for t := 0.0; t < distance; t += stepSize {
		// Базовая позиция на пути
		x := start.X + dx*t
		y := start.Y + dy*t

		// Добавляем случайность
		noiseScale := 0.05
		noiseVal := noise.Perlin2D(t*noiseScale, 0, 0.5) * 5.0

		// Вычисляем перпендикулярное направление
		perpX := -dy
		perpY := dx

		// Смещаем позицию
		x += perpX * noiseVal
		y += perpY * noiseVal

		// Проверяем, что позиция в пределах мира
		if x >= 0 && x < float64(g.width) && y >= 0 && y < float64(g.height) {
			// Создаем путь
			tileX := int(x)
			tileY := int(y)

			// Ширина пути
			pathWidth := 1 + g.random.RangeInt(0, 2)

			// Создаем путь указанной ширины
			for offsetY := -pathWidth; offsetY <= pathWidth; offsetY++ {
				for offsetX := -pathWidth; offsetX <= pathWidth; offsetX++ {
					nx := tileX + offsetX
					ny := tileY + offsetY

					// Проверяем, что координаты в пределах мира
					if nx >= 0 && nx < g.width && ny >= 0 && ny < g.height {
						// Устанавливаем тип тайла в путь
						g.world.Tiles[ny][nx].Type = TilePath
					}
				}
			}
		}
	}
}

// generateObjects генерирует объекты в мире
func (g *Generator) generateObjects() {
	// Очищаем существующие объекты
	g.world.Objects = []WorldObject{}

	// Генерируем объекты для каждой зоны
	for _, zone := range g.zones {
		// Определяем количество объектов в зависимости от плотности зоны
		objectCount := int(zone.Density * zone.Radius * zone.Radius * 0.05)

		// Генерируем объекты
		for i := 0; i < objectCount; i++ {
			// Выбираем случайную позицию в пределах зоны
			angle := g.random.Float64() * 2 * math.Pi
			distance := g.random.Float64() * zone.Radius

			x := zone.Position.X + math.Cos(angle)*distance
			y := zone.Position.Y + math.Sin(angle)*distance

			// Проверяем, что позиция в пределах мира
			if x < 0 || x >= float64(g.width) || y < 0 || y >= float64(g.height) {
				continue
			}

			// Получаем тип тайла
			tileX := int(x)
			tileY := int(y)
			tile := g.world.GetTileAt(tileX, tileY)

			// Пропускаем, если тайл непроходимый
			if tile == nil || tile.Type == TileWater || tile.Type == TileSwamp {
				continue
			}

			// Создаем объект в зависимости от темы зоны
			obj := g.createObjectForTheme(zone.Theme, entity.Vector2D{X: x, Y: y})

			// Добавляем объект в мир
			g.world.Objects = append(g.world.Objects, obj)

			// Добавляем объект в тайл
			tile.Objects = append(tile.Objects, obj)
		}
	}
}

// createObjectForTheme создает объект в соответствии с темой
func (g *Generator) createObjectForTheme(theme ThemeType, position entity.Vector2D) WorldObject {
	// Базовый ID объекта
	id := g.world.nextID
	g.world.nextID++

	// Выбираем тип объекта в зависимости от темы
	objectType := ""
	solid := true
	interactive := false

	switch theme {
	case ThemeDecay:
		objects := []string{"rotten_log", "dead_tree", "decomposed_body", "fungus", "slime"}
		objectType = objects[g.random.RangeInt(0, len(objects))]
		solid = objectType != "slime"
		interactive = objectType == "fungus" || objectType == "decomposed_body"

	case ThemeBlood:
		objects := []string{"blood_pool", "giblets", "meat_pile", "bone_pile", "hanging_corpse"}
		objectType = objects[g.random.RangeInt(0, len(objects))]
		solid = objectType == "hanging_corpse"
		interactive = objectType == "blood_pool" || objectType == "meat_pile"

	case ThemeShadow:
		objects := []string{"shadow_pillar", "dark_statue", "void_crack", "black_obelisk", "shadow_pool"}
		objectType = objects[g.random.RangeInt(0, len(objects))]
		solid = objectType != "shadow_pool" && objectType != "void_crack"
		interactive = objectType == "void_crack" || objectType == "shadow_pool"

	case ThemeChild:
		objects := []string{"broken_toy", "empty_swing", "school_desk", "cradle", "doll"}
		objectType = objects[g.random.RangeInt(0, len(objects))]
		solid = objectType != "broken_toy" && objectType != "doll"
		interactive = objectType == "doll" || objectType == "cradle"

	case ThemeHospital:
		objects := []string{"hospital_bed", "wheelchair", "medical_cabinet", "surgery_table", "iv_stand"}
		objectType = objects[g.random.RangeInt(0, len(objects))]
		solid = objectType != "iv_stand"
		interactive = objectType == "medical_cabinet" || objectType == "surgery_table"

	case ThemeAncient:
		objects := []string{"stone_altar", "ruined_pillar", "ancient_statue", "ritual_circle", "totem"}
		objectType = objects[g.random.RangeInt(0, len(objects))]
		solid = objectType != "ritual_circle"
		interactive = objectType == "stone_altar" || objectType == "ritual_circle"

	case ThemeIndustrial:
		objects := []string{"machinery", "pipe", "barrel", "control_panel", "generator"}
		objectType = objects[g.random.RangeInt(0, len(objects))]
		solid = objectType != "control_panel"
		interactive = objectType == "control_panel" || objectType == "generator"

	case ThemeForest:
		objects := []string{"tree", "stump", "bush", "rock", "fallen_log"}
		objectType = objects[g.random.RangeInt(0, len(objects))]
		solid = objectType != "bush"
		interactive = objectType == "stump" || objectType == "fallen_log"

	case ThemeVoid:
		objects := []string{"void_hole", "floating_rocks", "energy_pillar", "reality_tear", "cosmic_dust"}
		objectType = objects[g.random.RangeInt(0, len(objects))]
		solid = objectType != "cosmic_dust" && objectType != "void_hole"
		interactive = objectType == "void_hole" || objectType == "reality_tear"

	default:
		objects := []string{"tree", "rock", "bush", "stump", "fallen_log"}
		objectType = objects[g.random.RangeInt(0, len(objects))]
		solid = objectType != "bush"
		interactive = objectType == "stump"
	}

	// Создаем объект
	return WorldObject{
		ID:          id,
		Type:        objectType,
		Position:    position,
		Solid:       solid,
		Interactive: interactive,
	}
}

// generateCreatures генерирует существ в мире
func (g *Generator) generateCreatures() {
	// Очищаем существующих существ
	g.world.Entities = []*entity.Entity{}

	// Генерируем существ для каждой зоны
	for _, zone := range g.zones {
		// Определяем количество существ в зависимости от зоны
		var creatureCount int

		switch zone.Type {
		case ZoneSafe:
			creatureCount = 0 // Нет существ в безопасной зоне
		case ZoneTransition:
			creatureCount = 1 + g.random.RangeInt(0, 2)
		case ZoneExploration:
			creatureCount = 2 + g.random.RangeInt(0, 3)
		case ZoneDanger:
			creatureCount = 3 + g.random.RangeInt(0, 4)
		case ZoneNightmare:
			creatureCount = 4 + g.random.RangeInt(0, 5)
		default:
			creatureCount = 1
		}

		// Генерируем существ
		for i := 0; i < creatureCount; i++ {
			// Выбираем случайную позицию в пределах зоны
			angle := g.random.Float64() * 2 * math.Pi
			distance := g.random.Float64() * zone.Radius * 0.8 // Чуть ближе к центру

			x := zone.Position.X + math.Cos(angle)*distance
			y := zone.Position.Y + math.Sin(angle)*distance

			// Проверяем, что позиция в пределах мира
			if x < 0 || x >= float64(g.width) || y < 0 || y >= float64(g.height) {
				continue
			}

			// Получаем тип тайла
			tileX := int(x)
			tileY := int(y)
			tile := g.world.GetTileAt(tileX, tileY)

			// Пропускаем, если тайл непроходимый
			if tile == nil || tile.Type == TileWater || tile.Type == TileSwamp || tile.Type == TileRocks {
				continue
			}

			// Выбираем тип существа в зависимости от зоны и темы
			creatureType := g.selectCreatureType(zone.Type, zone.Theme, tile.Corruption)

			// Создаем существо
			creature := g.world.SpawnCreature(creatureType, entity.Vector2D{X: x, Y: y})

			// Настраиваем существо в зависимости от типа зоны
			g.configureCreatureForZone(creature, zone.Type)
		}
	}
}

// selectCreatureType выбирает тип существа
func (g *Generator) selectCreatureType(zoneType ZoneType, theme ThemeType, corruption float64) string {
	// Базовые типы существ
	creatureTypes := []string{"shadow", "spider", "phantom", "wendigo", "faceless", "doppelganger"}

	// Если уровень коррупции высокий, добавляем более искаженных существ
	if corruption > 0.7 {
		creatureTypes = append(creatureTypes, "abomination", "thing", "watcher", "nightmare")
	}

	// Добавляем существ в зависимости от темы
	switch theme {
	case ThemeDecay:
		creatureTypes = append(creatureTypes, "rotter", "slime", "mold_creature", "fungal_horror")
	case ThemeBlood:
		creatureTypes = append(creatureTypes, "flesh_beast", "blood_feeder", "gore_walker", "butcher")
	case ThemeShadow:
		creatureTypes = append(creatureTypes, "shadow_walker", "void_entity", "darkness", "shade")
	case ThemeChild:
		creatureTypes = append(creatureTypes, "toy_creature", "imaginary_friend", "lost_child", "puppet")
	case ThemeHospital:
		creatureTypes = append(creatureTypes, "patient", "doctor", "nurse", "experiment")
	case ThemeAncient:
		creatureTypes = append(creatureTypes, "cultist", "elder_thing", "idol", "forgotten_god")
	case ThemeIndustrial:
		creatureTypes = append(creatureTypes, "automaton", "worker", "living_machine", "rust_creature")
	case ThemeForest:
		creatureTypes = append(creatureTypes, "beast", "spriggan", "wolf", "stag")
	case ThemeVoid:
		creatureTypes = append(creatureTypes, "starving_void", "cosmic_horror", "beyond_one", "traveler")
	}

	// Выбираем случайный тип
	return creatureTypes[g.random.RangeInt(0, len(creatureTypes))]
}

// configureCreatureForZone настраивает существо в зависимости от типа зоны
func (g *Generator) configureCreatureForZone(creature *entity.Entity, zoneType ZoneType) {
	// Настраиваем существо в зависимости от типа зоны
	switch zoneType {
	case ZoneTransition:
		creature.Behavior = entity.NewBasicCreatureBehavior()

	case ZoneExploration:
		creature.Behavior = entity.NewPatrolBehavior()

	case ZoneDanger:
		creature.Behavior = entity.NewAggressiveBehavior()

	case ZoneNightmare:
		creature.Behavior = entity.NewStalkerBehavior()
	}
}

// ModifyWorldBasedOnPlayerFears модифицирует мир на основе страхов игрока
func (g *Generator) ModifyWorldBasedOnPlayerFears() {
	// Если нет системы наблюдения, выходим
	if g.observer == nil {
		return
	}

	// Получаем профиль страхов игрока
	fearProfile := g.observer.GetPlayerFearProfile()

	// Проходим по всем типам страха
	for fearType, fearValue := range fearProfile {
		// Если уровень страха достаточно высокий
		if fearValue > 0.5 {
			// Модифицируем мир в соответствии с типом страха
			g.modifyWorldForFearType(fearType, fearValue)
		}
	}
}

// modifyWorldForFearType модифицирует мир в соответствии с типом страха
func (g *Generator) modifyWorldForFearType(fearType ai.FearType, intensity float64) {
	// Выбираем тему на основе типа страха
	theme, ok := g.fearInfluence[fearType]
	if !ok {
		return
	}

	// Определяем количество модификаций
	modCount := int(intensity*5) + 1

	// Применяем модификации
	for i := 0; i < modCount; i++ {
		// Выбираем случайную позицию
		x := g.random.Float64() * float64(g.width)
		y := g.random.Float64() * float64(g.height)

		// Радиус воздействия
		radius := 10.0 + g.random.Float64()*20.0

		// Применяем тему к области
		g.applyThemeToArea(entity.Vector2D{X: x, Y: y}, radius, theme, intensity)
	}
}

// applyThemeToArea применяет тему к области
func (g *Generator) applyThemeToArea(position entity.Vector2D, radius float64, theme ThemeType, intensity float64) {
	radiusSq := radius * radius

	// Модифицируем тайлы в радиусе
	for y := int(math.Max(0, position.Y-radius)); y <= int(math.Min(float64(g.height-1), position.Y+radius)); y++ {
		for x := int(math.Max(0, position.X-radius)); x <= int(math.Min(float64(g.width-1), position.X+radius)); x++ {
			// Вычисляем расстояние до центра
			dx := float64(x) - position.X
			dy := float64(y) - position.Y
			distSq := dx*dx + dy*dy

			// Если точка в пределах радиуса
			if distSq <= radiusSq {
				// Вычисляем фактор влияния (ближе к центру - сильнее)
				factor := (1.0 - math.Sqrt(distSq)/radius) * intensity

				// Модифицируем тайл
				tile := &g.world.Tiles[y][x]
				g.applyThemeToTile(tile, theme, factor)
			}
		}
	}
}
