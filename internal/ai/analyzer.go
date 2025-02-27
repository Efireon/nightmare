package ai

import (
	"fmt"
	"math"
	"sort"
	"time"

	"nightmare/internal/entity"
)

// PlayerPattern представляет шаблон поведения игрока
type PlayerPattern struct {
	Name        string
	Description string
	Weight      float64
}

// MovementAnalysis содержит результаты анализа передвижения игрока
type MovementAnalysis struct {
	AverageSpeed     float64
	DirectionChanges int
	ExplorationArea  float64
	PathRepetition   float64
	PreferredAreas   []entity.Vector2D
}

// InteractionAnalysis содержит результаты анализа взаимодействий
type InteractionAnalysis struct {
	InteractionRate       float64
	PreferredInteractions map[string]int
	ResponseToScareEvents map[ScareEventType]float64
	HealthLossRate        float64
	SanityLossRate        float64
}

// Analyzer анализирует поведение игрока
type Analyzer struct {
	player              *entity.Player
	detectedPatterns    []PlayerPattern
	movementAnalysis    MovementAnalysis
	interactionAnalysis InteractionAnalysis
	scareHistory        []ScareEvent
	lastAnalysisTime    time.Time

	positionHistory []entity.Vector2D
	areaVisits      map[string]int  // ключ: "x,y" для сектора, значение: количество посещений
	sectorsExplored map[string]bool // ключ: "x,y" для сектора, значение: был ли исследован

	sectorSize float64     // размер одного сектора для анализа
	heatmap    [][]float64 // тепловая карта посещений

	scareResponses map[ScareEventType][]float64 // реакции на испуг
}

// NewAnalyzer создает новый анализатор
func NewAnalyzer(player *entity.Player) *Analyzer {
	return &Analyzer{
		player:           player,
		detectedPatterns: []PlayerPattern{},
		movementAnalysis: MovementAnalysis{
			AverageSpeed:     0,
			DirectionChanges: 0,
			ExplorationArea:  0,
			PathRepetition:   0,
			PreferredAreas:   []entity.Vector2D{},
		},
		interactionAnalysis: InteractionAnalysis{
			InteractionRate:       0,
			PreferredInteractions: make(map[string]int),
			ResponseToScareEvents: make(map[ScareEventType]float64),
			HealthLossRate:        0,
			SanityLossRate:        0,
		},
		scareHistory:     []ScareEvent{},
		lastAnalysisTime: time.Now(),
		positionHistory:  []entity.Vector2D{},
		areaVisits:       make(map[string]int),
		sectorsExplored:  make(map[string]bool),
		sectorSize:       5.0,                   // Размер сектора в единицах мира
		heatmap:          make([][]float64, 50), // 50x50 тепловая карта
		scareResponses:   make(map[ScareEventType][]float64),
	}
}

// AnalyzePlayer проводит комплексный анализ поведения игрока
func (a *Analyzer) AnalyzePlayer() {
	// Записываем текущую позицию игрока
	a.recordPlayerPosition()

	// Анализируем перемещения
	a.analyzeMovement()

	// Анализируем взаимодействия
	a.analyzeInteractions()

	// Обновляем тепловую карту
	a.updateHeatmap()

	// Определяем шаблоны поведения
	a.detectPatterns()

	// Обновляем время последнего анализа
	a.lastAnalysisTime = time.Now()
}

// recordPlayerPosition записывает текущую позицию игрока
func (a *Analyzer) recordPlayerPosition() {
	a.positionHistory = append(a.positionHistory, a.player.Position)

	// Ограничиваем размер истории
	if len(a.positionHistory) > 1000 {
		a.positionHistory = a.positionHistory[len(a.positionHistory)-1000:]
	}

	// Обновляем посещения секторов
	sectorX := int(a.player.Position.X / a.sectorSize)
	sectorY := int(a.player.Position.Y / a.sectorSize)
	sectorKey := makeKey(sectorX, sectorY)

	a.areaVisits[sectorKey]++
	a.sectorsExplored[sectorKey] = true
}

// makeKey создает строковый ключ из координат
func makeKey(x, y int) string {
	return fmt.Sprintf("%d,%d", x, y)
}

// analyzeMovement анализирует движение игрока
func (a *Analyzer) analyzeMovement() {
	if len(a.positionHistory) < 2 {
		return
	}

	// Вычисляем среднюю скорость
	totalDistance := 0.0
	directionChanges := 0
	for i := 1; i < len(a.positionHistory); i++ {
		// Расстояние между последовательными точками
		dist := distance(a.positionHistory[i-1], a.positionHistory[i])
		totalDistance += dist

		// Изменения направления
		if i > 1 {
			prev := a.positionHistory[i-2]
			curr := a.positionHistory[i-1]
			next := a.positionHistory[i]

			dir1 := math.Atan2(curr.Y-prev.Y, curr.X-prev.X)
			dir2 := math.Atan2(next.Y-curr.Y, next.X-curr.X)

			// Нормализуем разницу углов
			diff := math.Abs(dir2 - dir1)
			if diff > math.Pi {
				diff = 2*math.Pi - diff
			}

			// Считаем изменение направления, если угол больше порогового значения
			if diff > math.Pi/4 {
				directionChanges++
			}

		}
	}

	// Обновляем результаты анализа
	a.movementAnalysis.AverageSpeed = totalDistance / float64(len(a.positionHistory)-1)
	a.movementAnalysis.DirectionChanges = directionChanges

	// Вычисляем площадь исследованной области
	a.movementAnalysis.ExplorationArea = float64(len(a.sectorsExplored)) * (a.sectorSize * a.sectorSize)

	// Находим предпочитаемые области (наиболее посещаемые секторы)
	a.findPreferredAreas()

	// Вычисляем повторяемость пути
	a.calculatePathRepetition()
}

// analyzeInteractions анализирует взаимодействия игрока
func (a *Analyzer) analyzeInteractions() {
	// Анализируем записи о действиях
	totalActions := len(a.player.ActionLog)
	if totalActions == 0 {
		return
	}

	interactions := 0
	for _, action := range a.player.ActionLog {
		if action.Action == entity.ActionInteract {
			interactions++

			// Добавляем тип взаимодействия, если он есть
			if action.InteractionType != "" {
				a.interactionAnalysis.PreferredInteractions[action.InteractionType]++
			}
		}
	}

	// Вычисляем частоту взаимодействий
	a.interactionAnalysis.InteractionRate = float64(interactions) / float64(totalActions)
}

// RecordScareResponse записывает реакцию игрока на пугающее событие
func (a *Analyzer) RecordScareResponse(event ScareEvent, response float64) {
	// Записываем событие в историю
	a.scareHistory = append(a.scareHistory, event)

	// Ограничиваем размер истории
	if len(a.scareHistory) > 50 {
		a.scareHistory = a.scareHistory[len(a.scareHistory)-50:]
	}

	// Записываем реакцию
	if _, ok := a.scareResponses[event.Type]; !ok {
		a.scareResponses[event.Type] = []float64{}
	}
	a.scareResponses[event.Type] = append(a.scareResponses[event.Type], response)

	// Вычисляем среднюю реакцию
	total := 0.0
	for _, resp := range a.scareResponses[event.Type] {
		total += resp
	}
	avgResponse := total / float64(len(a.scareResponses[event.Type]))

	// Обновляем анализ
	a.interactionAnalysis.ResponseToScareEvents[event.Type] = avgResponse
}

// updateHeatmap обновляет тепловую карту посещений
func (a *Analyzer) updateHeatmap() {
	// Инициализируем тепловую карту, если нужно
	if len(a.heatmap) == 0 || len(a.heatmap[0]) == 0 {
		for i := range a.heatmap {
			a.heatmap[i] = make([]float64, 50)
		}
	}

	// Обновляем значения тепловой карты
	for key, visits := range a.areaVisits {
		var x, y int
		fmt.Sscanf(key, "%d,%d", &x, &y)

		// Преобразуем координаты мира в координаты тепловой карты
		heatmapX := x * 50 / 256 // Предполагаем, что мир 256x256
		heatmapY := y * 50 / 256

		// Проверяем, что координаты в пределах карты
		if heatmapX >= 0 && heatmapX < 50 && heatmapY >= 0 && heatmapY < 50 {
			// Увеличиваем значение в зависимости от количества посещений
			a.heatmap[heatmapY][heatmapX] = math.Min(1.0, float64(visits)/10.0)
		}
	}
}

// findPreferredAreas находит предпочитаемые области
func (a *Analyzer) findPreferredAreas() {
	// Очищаем предыдущие результаты
	a.movementAnalysis.PreferredAreas = []entity.Vector2D{}

	// Находим наиболее посещаемые секторы
	type sectorVisit struct {
		key    string
		visits int
	}

	// Преобразуем карту в срез для сортировки
	visits := []sectorVisit{}
	for key, count := range a.areaVisits {
		visits = append(visits, sectorVisit{key: key, visits: count})
	}

	// Сортируем по убыванию количества посещений
	sort.Slice(visits, func(i, j int) bool {
		return visits[i].visits > visits[j].visits
	})

	// Берем топ-5 или меньше
	count := min(5, len(visits))
	for i := 0; i < count; i++ {
		var x, y int
		fmt.Sscanf(visits[i].key, "%d,%d", &x, &y)

		// Преобразуем координаты сектора в координаты мира
		worldX := float64(x)*a.sectorSize + a.sectorSize/2
		worldY := float64(y)*a.sectorSize + a.sectorSize/2

		a.movementAnalysis.PreferredAreas = append(a.movementAnalysis.PreferredAreas,
			entity.Vector2D{X: worldX, Y: worldY})
	}
}

// calculatePathRepetition вычисляет повторяемость пути
func (a *Analyzer) calculatePathRepetition() {
	totalSectors := len(a.sectorsExplored)
	if totalSectors == 0 {
		a.movementAnalysis.PathRepetition = 0
		return
	}

	// Количество посещений каждого сектора
	totalVisits := 0
	for _, count := range a.areaVisits {
		totalVisits += count
	}

	// Вычисляем среднее количество посещений на сектор
	avgVisits := float64(totalVisits) / float64(totalSectors)

	// Повторяемость - отношение среднего количества посещений к ожидаемому (1)
	a.movementAnalysis.PathRepetition = math.Max(0, avgVisits-1)
}

// detectPatterns определяет шаблоны поведения игрока
func (a *Analyzer) detectPatterns() {
	// Очищаем предыдущие результаты
	a.detectedPatterns = []PlayerPattern{}

	// Анализируем на основе передвижения
	a.detectMovementPatterns()

	// Анализируем на основе взаимодействий
	a.detectInteractionPatterns()

	// Анализируем на основе реакций на испуг
	a.detectFearResponsePatterns()

	// Сортируем шаблоны по весу
	sort.Slice(a.detectedPatterns, func(i, j int) bool {
		return a.detectedPatterns[i].Weight > a.detectedPatterns[j].Weight
	})
}

// detectMovementPatterns определяет шаблоны передвижения
func (a *Analyzer) detectMovementPatterns() {
	// Исследователь - много исследует, мало повторяется
	if a.movementAnalysis.ExplorationArea > 500 && a.movementAnalysis.PathRepetition < 2 {
		a.addPattern(PlayerPattern{
			Name:        "explorer",
			Description: "Игрок активно исследует мир, не задерживаясь на одном месте",
			Weight:      0.8 - (a.movementAnalysis.PathRepetition / 10),
		})
	}

	// Осторожный - медленно двигается, много меняет направление
	if a.movementAnalysis.AverageSpeed < 1.5 && a.movementAnalysis.DirectionChanges > 30 {
		a.addPattern(PlayerPattern{
			Name:        "cautious",
			Description: "Игрок осторожен, медленно двигается и часто меняет направление",
			Weight:      0.9 - (a.movementAnalysis.AverageSpeed / 3),
		})
	}

	// Целеустремленный - быстро и целенаправленно движется
	if a.movementAnalysis.AverageSpeed > 2.0 && a.movementAnalysis.DirectionChanges < 15 {
		a.addPattern(PlayerPattern{
			Name:        "determined",
			Description: "Игрок быстро движется в выбранном направлении",
			Weight:      0.7 + (a.movementAnalysis.AverageSpeed / 5),
		})
	}

	// Нерешительный - много топчется на месте
	if a.movementAnalysis.ExplorationArea < 200 && a.movementAnalysis.PathRepetition > 3 {
		a.addPattern(PlayerPattern{
			Name:        "indecisive",
			Description: "Игрок нерешителен, часто возвращается в одни и те же места",
			Weight:      0.6 + (a.movementAnalysis.PathRepetition / 5),
		})
	}
}

// detectInteractionPatterns определяет шаблоны взаимодействия
func (a *Analyzer) detectInteractionPatterns() {
	// Интерактивный - часто взаимодействует с миром
	if a.interactionAnalysis.InteractionRate > 0.3 {
		a.addPattern(PlayerPattern{
			Name:        "interactive",
			Description: "Игрок активно взаимодействует с окружением",
			Weight:      0.7 + a.interactionAnalysis.InteractionRate,
		})
	}

	// Пассивный - редко взаимодействует с миром
	if a.interactionAnalysis.InteractionRate < 0.1 {
		a.addPattern(PlayerPattern{
			Name:        "passive",
			Description: "Игрок редко взаимодействует с окружением",
			Weight:      0.6 + (0.1 - a.interactionAnalysis.InteractionRate),
		})
	}
}

// detectFearResponsePatterns определяет шаблоны реакции на страх
func (a *Analyzer) detectFearResponsePatterns() {
	// Общая реакция на испуг
	avgResponse := 0.0
	count := 0

	for _, response := range a.interactionAnalysis.ResponseToScareEvents {
		avgResponse += response
		count++
	}

	if count > 0 {
		avgResponse /= float64(count)

		// Устойчивый к испугу
		if avgResponse < 0.3 {
			a.addPattern(PlayerPattern{
				Name:        "fearless",
				Description: "Игрок слабо реагирует на пугающие события",
				Weight:      0.8 - avgResponse,
			})
		}

		// Легко пугается
		if avgResponse > 0.7 {
			a.addPattern(PlayerPattern{
				Name:        "easily_scared",
				Description: "Игрок сильно реагирует на пугающие события",
				Weight:      0.7 + avgResponse,
			})
		}
	}

	// Реакция на конкретные типы испуга
	if response, ok := a.interactionAnalysis.ResponseToScareEvents[EventSuddenNoise]; ok && response > 0.8 {
		a.addPattern(PlayerPattern{
			Name:        "startles_easily",
			Description: "Игрок особенно чувствителен к внезапным звукам",
			Weight:      0.7 + response,
		})
	}

	if response, ok := a.interactionAnalysis.ResponseToScareEvents[EventCreatureAppearance]; ok && response > 0.8 {
		a.addPattern(PlayerPattern{
			Name:        "monster_phobia",
			Description: "Игрок особенно боится встреч с существами",
			Weight:      0.7 + response,
		})
	}
}

// addPattern добавляет шаблон поведения
func (a *Analyzer) addPattern(pattern PlayerPattern) {
	// Проверяем, есть ли уже такой шаблон
	for i, p := range a.detectedPatterns {
		if p.Name == pattern.Name {
			// Обновляем вес существующего шаблона
			a.detectedPatterns[i].Weight = (a.detectedPatterns[i].Weight + pattern.Weight) / 2
			return
		}
	}

	// Добавляем новый шаблон
	a.detectedPatterns = append(a.detectedPatterns, pattern)
}

// GetTopPatterns возвращает наиболее выраженные шаблоны поведения
func (a *Analyzer) GetTopPatterns(count int) []PlayerPattern {
	if count > len(a.detectedPatterns) {
		count = len(a.detectedPatterns)
	}

	return a.detectedPatterns[:count]
}

// GetHeatmap возвращает тепловую карту посещений
func (a *Analyzer) GetHeatmap() [][]float64 {
	return a.heatmap
}

// GetMovementAnalysis возвращает результаты анализа передвижения
func (a *Analyzer) GetMovementAnalysis() MovementAnalysis {
	return a.movementAnalysis
}

// GetInteractionAnalysis возвращает результаты анализа взаимодействий
func (a *Analyzer) GetInteractionAnalysis() InteractionAnalysis {
	return a.interactionAnalysis
}

// distance вычисляет расстояние между двумя точками
func distance(a, b entity.Vector2D) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// min возвращает минимальное из двух чисел
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
