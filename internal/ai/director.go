package ai

import (
	"math"
	"math/rand"
	"time"

	"nightmare/internal/entity"
	"nightmare/internal/world"
)

// ScareEvent представляет пугающее событие
type ScareEvent struct {
	Type         ScareEventType
	Intensity    float64 // От 0 до 1
	Position     entity.Vector2D
	Duration     time.Duration
	CreatureType string // Тип существа, если событие связано с существом
}

// ScareEventType тип пугающего события
type ScareEventType int

const (
	EventAmbientSound ScareEventType = iota
	EventSuddenNoise
	EventCreatureAppearance
	EventEnvironmentChange
	EventHallucination
	EventWhisper
)

// BehaviorPattern описывает паттерн поведения игрока
type BehaviorPattern struct {
	MovementPreference    float64  // Предпочтение движения (положительное - активное, отрицательное - пассивное)
	ExplorationPreference float64  // Предпочтение исследования (высокое - исследователь, низкое - линейный игрок)
	RiskTolerance         float64  // Толерантность к риску (высокая - смелый, низкая - осторожный)
	ReactivityToScares    float64  // Реакция на испуг (высокая - сильно реагирует, низкая - слабо реагирует)
	PreferredInteractions []string // Типы взаимодействий, которые игрок предпочитает
}

// Director представляет ИИ-директора, управляющего игровыми событиями
type Director struct {
	player             *entity.Player
	world              *world.World
	playerBehavior     BehaviorPattern
	scareHistory       []ScareEvent
	scareEffectiveness map[ScareEventType]float64 // Эффективность различных видов испуга
	lastAnalysisTime   time.Time
	mood               float64 // Общее "настроение" директора от 0 (спокойное) до 1 (агрессивное)
	tension            float64 // Текущий уровень напряжения от 0 до 1
}

// NewDirector создает нового ИИ-директора
func NewDirector(player *entity.Player, world *world.World) *Director {
	return &Director{
		player: player,
		world:  world,
		playerBehavior: BehaviorPattern{
			MovementPreference:    0.5,
			ExplorationPreference: 0.5,
			RiskTolerance:         0.5,
			ReactivityToScares:    0.5,
			PreferredInteractions: []string{},
		},
		scareHistory:       []ScareEvent{},
		scareEffectiveness: make(map[ScareEventType]float64),
		lastAnalysisTime:   time.Now(),
		mood:               0.3, // Начальное настроение
		tension:            0.1, // Начальное напряжение
	}
}

// AnalyzePlayerBehavior анализирует поведение игрока
func (d *Director) AnalyzePlayerBehavior() {
	// Если логов действий игрока нет, ничего не делаем
	if len(d.player.ActionLog) == 0 {
		return
	}

	// Анализируем только логи с последнего анализа
	recentLogs := []entity.PlayerActionRecord{}
	for _, log := range d.player.ActionLog {
		if log.Timestamp.After(d.lastAnalysisTime) {
			recentLogs = append(recentLogs, log)
		}
	}

	d.lastAnalysisTime = time.Now()

	// Если нет новых логов, ничего не делаем
	if len(recentLogs) == 0 {
		return
	}

	// Анализ движения
	moveCount := 0
	for _, log := range recentLogs {
		if log.Action == entity.ActionMove {
			moveCount++
		}
	}

	movementRatio := float64(moveCount) / float64(len(recentLogs))

	// Обновляем предпочтения движения (с инерцией)
	d.playerBehavior.MovementPreference = d.playerBehavior.MovementPreference*0.8 + movementRatio*0.2

	// Анализ исследования (насколько игрок отклоняется от прямого пути)
	// Это более сложный анализ, который мы упростим для примера

	// Обновляем другие аспекты поведения
	// ...

	// Анализ эффективности прошлых попыток испугать
	d.analyzeScareEffectiveness()

	// Обновляем настроение директора и уровень напряжения
	d.updateMoodAndTension()
}

// AdjustWorld изменяет мир на основе анализа поведения игрока
func (d *Director) AdjustWorld() {
	// Решаем, нужно ли создать пугающее событие
	if d.shouldCreateScareEvent() {
		event := d.createScareEvent()
		d.executeScareEvent(event)
	}

	// Модификация окружающего мира
	d.modifyEnvironment()

	// Управление существами
	d.manageCreatures()
}

// shouldCreateScareEvent определяет, нужно ли создать пугающее событие
func (d *Director) shouldCreateScareEvent() bool {
	// Базовый шанс на основе напряжения
	baseChance := d.tension * 0.1

	// Увеличиваем шанс, если игрок долго не испытывал испуга
	if len(d.scareHistory) > 0 {
		lastScare := d.scareHistory[len(d.scareHistory)-1]
		timeSinceLast := time.Since(lastScare.Timestamp)
		if timeSinceLast > 30*time.Second {
			baseChance += 0.1
		}
		if timeSinceLast > 60*time.Second {
			baseChance += 0.2
		}
	} else {
		// Если еще не было испуга, увеличиваем шанс
		baseChance += 0.3
	}

	// Добавляем случайность
	return rand.Float64() < baseChance
}

// createScareEvent создает пугающее событие на основе поведения игрока
func (d *Director) createScareEvent() ScareEvent {
	// Выбираем тип события с учетом эффективности
	eventType := d.chooseEventType()

	// Определяем интенсивность на основе настроения и анализа игрока
	intensity := d.mood * (0.7 + rand.Float64()*0.3)

	// Если игрок слабо реагирует на испуг, увеличиваем интенсивность
	if d.playerBehavior.ReactivityToScares < 0.3 {
		intensity *= 1.5
	}

	// Ограничиваем интенсивность
	if intensity > 1.0 {
		intensity = 1.0
	}

	// Создаем событие
	event := ScareEvent{
		Type:      eventType,
		Intensity: intensity,
		Position:  d.player.Position, // По умолчанию рядом с игроком
		Duration:  time.Duration(2+rand.Intn(5)) * time.Second,
	}

	// Для некоторых типов событий нужна дополнительная настройка
	if eventType == EventCreatureAppearance {
		// Выбираем тип существа
		event.CreatureType = d.chooseCreatureType()

		// Устанавливаем позицию появления существа
		angle := rand.Float64() * 2 * math.Pi
		distance := 10.0 + rand.Float64()*20.0
		event.Position = entity.Vector2D{
			X: d.player.Position.X + math.Cos(angle)*distance,
			Y: d.player.Position.Y + math.Sin(angle)*distance,
		}
	}

	return event
}

// chooseEventType выбирает тип события с учетом эффективности
func (d *Director) chooseEventType() ScareEventType {
	// Если у нас нет данных об эффективности, выбираем случайный тип
	if len(d.scareEffectiveness) == 0 {
		return ScareEventType(rand.Intn(6))
	}

	// Выбираем более эффективные типы с большей вероятностью
	// ...

	// Упрощенный вариант - случайный выбор
	return ScareEventType(rand.Intn(6))
}

// chooseCreatureType выбирает тип существа
func (d *Director) chooseCreatureType() string {
	creatureTypes := []string{
		"shadow", "spider", "phantom", "doppelganger", "wendigo", "faceless",
	}

	return creatureTypes[rand.Intn(len(creatureTypes))]
}

// executeScareEvent выполняет пугающее событие
func (d *Director) executeScareEvent(event ScareEvent) {
	// Добавляем событие в историю
	d.scareHistory = append(d.scareHistory, event)

	// Выполняем действия в зависимости от типа события
	switch event.Type {
	case EventAmbientSound:
		// Воспроизведение звука
		// ...

	case EventSuddenNoise:
		// Внезапный громкий звук
		// ...

	case EventCreatureAppearance:
		// Создание существа
		worldPos := world.Vector2D{X: event.Position.X, Y: event.Position.Y}
		d.world.SpawnCreature(event.CreatureType, worldPos)

	case EventEnvironmentChange:
		// Изменение окружения
		worldPos := world.Vector2D{X: event.Position.X, Y: event.Position.Y}
		d.world.ModifyEnvironment(worldPos, event.Intensity)

	case EventHallucination:
		// Создание галлюцинации
		// ...

	case EventWhisper:
		// Шепот
		// ...
	}

	// Снижаем рассудок игрока в зависимости от интенсивности события
	d.player.ReduceSanity(event.Intensity * 5)
}

// analyzeScareEffectiveness анализирует эффективность прошлых попыток испугать
func (d *Director) analyzeScareEffectiveness() {
	// Этот метод будет анализировать, как сильно изменился уровень рассудка
	// после каждого пугающего события
	// ...
}

// updateMoodAndTension обновляет настроение директора и уровень напряжения
func (d *Director) updateMoodAndTension() {
	// Настроение зависит от того, насколько эффективно мы пугаем игрока
	// ...

	// Напряжение растет со временем и падает после успешного испуга
	d.tension += 0.01

	// Ограничиваем напряжение
	if d.tension > 1.0 {
		d.tension = 1.0
	}
}

// modifyEnvironment изменяет окружающий мир
func (d *Director) modifyEnvironment() {
	// Здесь будет логика изменения окружающего мира
	// ...
}

// manageCreatures управляет существами в мире
func (d *Director) manageCreatures() {
	// Здесь будет логика управления существами
	// ...
}
