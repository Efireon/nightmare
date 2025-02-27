package ai

import (
	"math"
	"sort"
	"time"

	"nightmare/internal/entity"
	"nightmare/internal/event"
	"nightmare/internal/util"
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

// ReactorType представляет тип реакции игрока
type ReactorType int

const (
	ReactorCautious   ReactorType = iota // Осторожный игрок
	ReactorBold                          // Смелый игрок
	ReactorPanic                         // Паникующий игрок
	ReactorMethodical                    // Методичный игрок
	ReactorReckless                      // Безрассудный игрок
	ReactorHesitant                      // Нерешительный игрок
)

// ActionType представляет тип действия игрока
type ActionType int

const (
	ActionMove ActionType = iota
	ActionRun
	ActionHide
	ActionInteract
	ActionAttack
	ActionInvestigate
	ActionRetreat
	ActionFreeze
)

// PlayerAction представляет действие игрока
type PlayerAction struct {
	Type      ActionType
	Position  entity.Vector2D
	Direction entity.Vector2D
	Target    interface{}
	Timestamp time.Time
	Context   map[string]interface{}
}

// FearResponse представляет реакцию на страх
type FearResponse struct {
	FearType        FearType
	StrengthOfFear  float64 // От 0 до 1
	ActionsTaken    []ActionType
	SanityLoss      float64
	HeartRateChange float64
	ReactionTime    time.Duration
	DistanceMoved   float64
}

// ObservationContext представляет контекст наблюдения
type ObservationContext struct {
	NearbyCreatures    []*entity.Creature
	LightLevel         float64 // От 0 до 1
	OpenSpace          float64 // От 0 до 1 (1 = открытое пространство)
	NearbyExits        int
	TimeSinceLastScare time.Duration
	CurrentSanity      float64
	RecentSanityLoss   float64
}

// ObserverSystem отвечает за наблюдение и анализ действий игрока
type ObserverSystem struct {
	player       *entity.Player
	eventManager *event.EventManager
	analyzer     *Analyzer
	director     *Director
	random       *util.RandomGenerator

	playerActions  []PlayerAction
	fearResponses  map[FearType][]FearResponse
	reactorProfile map[ReactorType]float64
	fearProfile    map[FearType]float64

	lastObservationTime time.Time
	observationInterval time.Duration

	context ObservationContext

	predictedActions  map[ActionType]float64
	recommendedScares []ScareRecommendation
}

// ScareRecommendation представляет рекомендацию для испуга
type ScareRecommendation struct {
	ScareType  ScareEventType
	FearTarget FearType
	Intensity  float64
	Position   entity.Vector2D
	Timing     time.Duration
	Priority   float64
}

// NewObserverSystem создает новую систему наблюдения
func NewObserverSystem(player *entity.Player, eventManager *event.EventManager, analyzer *Analyzer, director *Director) *ObserverSystem {
	return &ObserverSystem{
		player:       player,
		eventManager: eventManager,
		analyzer:     analyzer,
		director:     director,
		random:       util.NewRandomGenerator(time.Now().UnixNano()),

		playerActions:  []PlayerAction{},
		fearResponses:  make(map[FearType][]FearResponse),
		reactorProfile: make(map[ReactorType]float64),
		fearProfile:    make(map[FearType]float64),

		lastObservationTime: time.Now(),
		observationInterval: 5 * time.Second, // Обновлять анализ каждые 5 секунд

		context: ObservationContext{
			LightLevel:         0.5,
			OpenSpace:          0.5,
			NearbyExits:        2,
			TimeSinceLastScare: 60 * time.Second,
			CurrentSanity:      player.Sanity,
			RecentSanityLoss:   0,
		},

		predictedActions:  make(map[ActionType]float64),
		recommendedScares: []ScareRecommendation{},
	}
}

// Initialize инициализирует систему наблюдения
func (o *ObserverSystem) Initialize() {
	// Инициализируем профили с нейтральными значениями
	for fearType := FearDarkness; fearType <= FearUnknown; fearType++ {
		o.fearProfile[fearType] = 0.5
	}

	for reactorType := ReactorCautious; reactorType <= ReactorHesitant; reactorType++ {
		o.reactorProfile[reactorType] = 0.5
	}

	// Подписываемся на события
	o.subscribeToEvents()
}

// subscribeToEvents подписывается на события
func (o *ObserverSystem) subscribeToEvents() {
	if o.eventManager == nil {
		return
	}

	// Движение игрока
	o.eventManager.AddListener(event.EventPlayerMoved, func(data event.EventData) {
		o.recordMovement(data)
	})

	// Получение урона
	o.eventManager.AddListener(event.EventPlayerDamaged, func(data event.EventData) {
		o.recordDamage(data)
	})

	// Изменение рассудка
	o.eventManager.AddListener(event.EventPlayerSanityChanged, func(data event.EventData) {
		o.recordSanityChange(data)
	})

	// Взаимодействие с объектами
	o.eventManager.AddListener(event.EventPlayerInteracted, func(data event.EventData) {
		o.recordInteraction(data)
	})

	// Пугающие события
	o.eventManager.AddListener(event.EventScareTriggered, func(data event.EventData) {
		o.recordScareResponse(data)
	})
}

// Update обновляет состояние системы наблюдения
func (o *ObserverSystem) Update() {
	currentTime := time.Now()

	// Проверяем, прошел ли достаточный интервал для анализа
	if currentTime.Sub(o.lastObservationTime) >= o.observationInterval {
		o.AnalyzePlayerBehavior()
		o.GenerateScareRecommendations()
		o.lastObservationTime = currentTime
	}

	// Обновляем контекст
	o.updateContext()
}

// updateContext обновляет контекст наблюдения
func (o *ObserverSystem) updateContext() {
	// Обновляем текущий уровень рассудка
	if o.player != nil {
		oldSanity := o.context.CurrentSanity
		o.context.CurrentSanity = o.player.Sanity
		o.context.RecentSanityLoss = oldSanity - o.player.Sanity
		if o.context.RecentSanityLoss < 0 {
			o.context.RecentSanityLoss = 0
		}
	}

	// Увеличиваем время с последнего испуга
	o.context.TimeSinceLastScare += time.Second
}

// recordMovement записывает движение игрока
func (o *ObserverSystem) recordMovement(data event.EventData) {
	// Проверяем, что данные содержат позицию
	if data.Position == nil {
		return
	}

	// Получаем позицию
	position, ok := data.Position.(entity.Vector2D)
	if !ok {
		return
	}

	// Получаем предыдущую позицию из кастомных данных
	var oldPosition entity.Vector2D
	if prevPos, ok := data.Custom["oldPosition"]; ok {
		oldPosition, _ = prevPos.(entity.Vector2D)
	}

	// Вычисляем направление
	direction := entity.Vector2D{
		X: position.X - oldPosition.X,
		Y: position.Y - oldPosition.Y,
	}

	// Определяем тип действия (движение или бег)
	actionType := ActionMove
	if data.Value != nil {
		if speed, ok := data.Value.(float64); ok {
			if speed > 1.5 {
				actionType = ActionRun
			}
		}
	}

	// Записываем действие
	o.addPlayerAction(PlayerAction{
		Type:      actionType,
		Position:  position,
		Direction: direction,
		Timestamp: data.Timestamp,
		Context:   make(map[string]interface{}),
	})
}

// recordDamage записывает получение урона
func (o *ObserverSystem) recordDamage(data event.EventData) {
	// Проверяем, что данные содержат источник урона
	if data.Source == nil {
		return
	}

	// Записываем действие
	o.addPlayerAction(PlayerAction{
		Type:      ActionFreeze, // Предполагаем, что при получении урона игрок на мгновение замирает
		Position:  o.player.Position,
		Timestamp: data.Timestamp,
		Target:    data.Source,
		Context: map[string]interface{}{
			"damageAmount": data.Value,
		},
	})
}

// recordSanityChange записывает изменение рассудка
func (o *ObserverSystem) recordSanityChange(data event.EventData) {
	// Проверяем, что данные содержат значение изменения рассудка
	if data.Value == nil {
		return
	}

	// Получаем значение изменения
	sanityChange, ok := data.Value.(float64)
	if !ok {
		return
	}

	// Если рассудок уменьшился, записываем реакцию на страх
	if sanityChange < 0 {
		oldValue, _ := data.Custom["oldValue"].(float64)
		newValue, _ := data.Custom["newValue"].(float64)

		// Определяем тип страха на основе источника
		fearType := FearUnknown
		if data.Source != nil {
			fearType = o.determineFearType(data.Source)
		}

		// Записываем реакцию на страх
		o.recordFearResponse(FearResponse{
			FearType:       fearType,
			StrengthOfFear: math.Abs(sanityChange) / 20.0, // Нормализуем к диапазону [0, 1]
			SanityLoss:     oldValue - newValue,
		})
	}
}

// recordInteraction записывает взаимодействие игрока
func (o *ObserverSystem) recordInteraction(data event.EventData) {
	// Проверяем, что данные содержат цель взаимодействия
	if data.Target == nil {
		return
	}

	// Записываем действие
	o.addPlayerAction(PlayerAction{
		Type:      ActionInteract,
		Position:  o.player.Position,
		Timestamp: data.Timestamp,
		Target:    data.Target,
		Context:   make(map[string]interface{}),
	})
}

// recordScareResponse записывает реакцию на пугающее событие
func (o *ObserverSystem) recordScareResponse(data event.EventData) {
	// Проверяем, что данные содержат тип пугающего события
	if data.Custom == nil || data.Custom["scareType"] == nil {
		return
	}

	// Получаем тип пугающего события
	scareType, ok := data.Custom["scareType"].(ScareEventType)
	if !ok {
		return
	}

	// Определяем тип страха на основе типа пугающего события
	fearType := o.mapScareEventToFearType(scareType)

	// Получаем интенсивность страха
	intensity := 0.5
	if data.Value != nil {
		if val, ok := data.Value.(float64); ok {
			intensity = val
		}
	}

	// Записываем реакцию на страх
	o.recordFearResponse(FearResponse{
		FearType:       fearType,
		StrengthOfFear: intensity,
		SanityLoss:     o.context.RecentSanityLoss,
	})

	// Сбрасываем время с последнего испуга
	o.context.TimeSinceLastScare = 0
}

// addPlayerAction добавляет действие игрока в историю
func (o *ObserverSystem) addPlayerAction(action PlayerAction) {
	o.playerActions = append(o.playerActions, action)

	// Ограничиваем размер истории
	maxActions := 1000
	if len(o.playerActions) > maxActions {
		o.playerActions = o.playerActions[len(o.playerActions)-maxActions:]
	}
}

// recordFearResponse записывает реакцию на страх
func (o *ObserverSystem) recordFearResponse(response FearResponse) {
	if _, ok := o.fearResponses[response.FearType]; !ok {
		o.fearResponses[response.FearType] = []FearResponse{}
	}

	o.fearResponses[response.FearType] = append(o.fearResponses[response.FearType], response)

	// Ограничиваем размер истории для каждого типа страха
	maxResponses := 20
	if len(o.fearResponses[response.FearType]) > maxResponses {
		o.fearResponses[response.FearType] = o.fearResponses[response.FearType][len(o.fearResponses[response.FearType])-maxResponses:]
	}
}

// determineFearType определяет тип страха на основе источника
func (o *ObserverSystem) determineFearType(source interface{}) FearType {
	switch source.(type) {
	case *entity.Creature:
		return FearCreatures
	default:
		return FearUnknown
	}
}

// mapScareEventToFearType сопоставляет тип пугающего события с типом страха
func (o *ObserverSystem) mapScareEventToFearType(scareType ScareEventType) FearType {
	switch scareType {
	case EventAmbientSound:
		return FearIsolation
	case EventSuddenNoise:
		return FearSuddenNoises
	case EventCreatureAppearance:
		return FearCreatures
	case EventEnvironmentChange:
		return FearUnknown
	case EventHallucination:
		return FearIsolation
	case EventWhisper:
		return FearIsolation
	default:
		return FearUnknown
	}
}

// AnalyzePlayerBehavior анализирует поведение игрока
func (o *ObserverSystem) AnalyzePlayerBehavior() {
	// Анализируем профиль реактора
	o.analyzeReactorProfile()

	// Анализируем профиль страхов
	o.analyzeFearProfile()

	// Предсказываем будущие действия
	o.predictActions()
}

// analyzeReactorProfile анализирует профиль реактора игрока
func (o *ObserverSystem) analyzeReactorProfile() {
	if len(o.playerActions) < 10 {
		return
	}

	// Счетчики для разных типов действий
	moveCount := 0
	runCount := 0
	hideCount := 0
	interactCount := 0
	attackCount := 0
	investigateCount := 0
	retreatCount := 0
	freezeCount := 0

	// Подсчитываем действия
	for _, action := range o.playerActions {
		switch action.Type {
		case ActionMove:
			moveCount++
		case ActionRun:
			runCount++
		case ActionHide:
			hideCount++
		case ActionInteract:
			interactCount++
		case ActionAttack:
			attackCount++
		case ActionInvestigate:
			investigateCount++
		case ActionRetreat:
			retreatCount++
		case ActionFreeze:
			freezeCount++
		}
	}

	totalActions := len(o.playerActions)

	// Вычисляем профиль реактора

	// Осторожный игрок - много скрывается, мало атакует, много отступает
	o.reactorProfile[ReactorCautious] = float64(hideCount+retreatCount-attackCount) / float64(totalActions)

	// Смелый игрок - мало скрывается, много атакует и исследует
	o.reactorProfile[ReactorBold] = float64(attackCount+investigateCount-hideCount) / float64(totalActions)

	// Паникующий игрок - много бегает и замирает
	o.reactorProfile[ReactorPanic] = float64(runCount+freezeCount) / float64(totalActions)

	// Методичный игрок - много исследует и взаимодействует
	o.reactorProfile[ReactorMethodical] = float64(investigateCount+interactCount) / float64(totalActions)

	// Безрассудный игрок - много бегает и атакует
	o.reactorProfile[ReactorReckless] = float64(runCount+attackCount-hideCount) / float64(totalActions)

	// Нерешительный игрок - много замирает и мало взаимодействует
	o.reactorProfile[ReactorHesitant] = float64(freezeCount-interactCount) / float64(totalActions)

	// Нормализуем значения к диапазону [0, 1]
	normalizeProfile(o.reactorProfile)
}

// analyzeFearProfile анализирует профиль страхов игрока
func (o *ObserverSystem) analyzeFearProfile() {
	// Если нет данных о реакциях на страх, выходим
	if len(o.fearResponses) == 0 {
		return
	}

	// Вычисляем среднюю силу страха для каждого типа
	for fearType := FearDarkness; fearType <= FearUnknown; fearType++ {
		responses, ok := o.fearResponses[fearType]
		if !ok || len(responses) == 0 {
			// Если нет данных, оставляем текущее значение
			continue
		}

		// Суммируем силу страха
		totalStrength := 0.0
		for _, response := range responses {
			totalStrength += response.StrengthOfFear
		}

		// Вычисляем среднее значение
		avgStrength := totalStrength / float64(len(responses))

		// Обновляем профиль с инерцией (70% старое значение, 30% новое)
		o.fearProfile[fearType] = o.fearProfile[fearType]*0.7 + avgStrength*0.3
	}
}

// predictActions предсказывает будущие действия игрока
func (o *ObserverSystem) predictActions() {
	// Сбрасываем предыдущие предсказания
	for actionType := ActionMove; actionType <= ActionFreeze; actionType++ {
		o.predictedActions[actionType] = 0
	}

	// Определяем доминирующий тип реактора
	dominantReactor := o.getDominantReactorType()

	// Предсказываем действия на основе типа реактора
	switch dominantReactor {
	case ReactorCautious:
		o.predictedActions[ActionHide] = 0.4
		o.predictedActions[ActionRetreat] = 0.3
		o.predictedActions[ActionMove] = 0.2
		o.predictedActions[ActionInvestigate] = 0.1

	case ReactorBold:
		o.predictedActions[ActionInvestigate] = 0.4
		o.predictedActions[ActionAttack] = 0.3
		o.predictedActions[ActionMove] = 0.2
		o.predictedActions[ActionInteract] = 0.1

	case ReactorPanic:
		o.predictedActions[ActionRun] = 0.5
		o.predictedActions[ActionFreeze] = 0.3
		o.predictedActions[ActionRetreat] = 0.2

	case ReactorMethodical:
		o.predictedActions[ActionInvestigate] = 0.4
		o.predictedActions[ActionInteract] = 0.3
		o.predictedActions[ActionMove] = 0.2
		o.predictedActions[ActionHide] = 0.1

	case ReactorReckless:
		o.predictedActions[ActionRun] = 0.4
		o.predictedActions[ActionAttack] = 0.3
		o.predictedActions[ActionMove] = 0.2
		o.predictedActions[ActionInvestigate] = 0.1

	case ReactorHesitant:
		o.predictedActions[ActionFreeze] = 0.4
		o.predictedActions[ActionMove] = 0.3
		o.predictedActions[ActionRetreat] = 0.2
		o.predictedActions[ActionHide] = 0.1
	}
}

// getDominantReactorType возвращает доминирующий тип реактора
func (o *ObserverSystem) getDominantReactorType() ReactorType {
	// Находим тип реактора с наибольшим значением
	maxValue := -1.0
	dominantType := ReactorCautious

	for reactorType, value := range o.reactorProfile {
		if value > maxValue {
			maxValue = value
			dominantType = reactorType
		}
	}

	return dominantType
}

// GenerateScareRecommendations генерирует рекомендации для испуга
func (o *ObserverSystem) GenerateScareRecommendations() {
	// Очищаем предыдущие рекомендации
	o.recommendedScares = []ScareRecommendation{}

	// Находим наиболее эффективные типы страха
	effectiveFears := o.getEffectiveFearTypes()

	// Генерируем рекомендации для каждого эффективного типа страха
	for _, fearType := range effectiveFears {
		recommendation := o.generateRecommendationForFearType(fearType)
		o.recommendedScares = append(o.recommendedScares, recommendation)
	}

	// Сортируем рекомендации по приоритету
	sort.Slice(o.recommendedScares, func(i, j int) bool {
		return o.recommendedScares[i].Priority > o.recommendedScares[j].Priority
	})
}

// getEffectiveFearTypes возвращает наиболее эффективные типы страха
func (o *ObserverSystem) getEffectiveFearTypes() []FearType {
	// Создаем карту эффективности для каждого типа страха
	effectiveness := make(map[FearType]float64)

	for fearType, value := range o.fearProfile {
		// Более высокий профиль страха означает большую эффективность
		effectiveness[fearType] = value
	}

	// Сортируем типы страха по эффективности
	var fearTypes []FearType
	for fearType := range effectiveness {
		fearTypes = append(fearTypes, fearType)
	}

	sort.Slice(fearTypes, func(i, j int) bool {
		return effectiveness[fearTypes[i]] > effectiveness[fearTypes[j]]
	})

	// Возвращаем до 3 наиболее эффективных типов страха
	count := 3
	if len(fearTypes) < count {
		count = len(fearTypes)
	}

	return fearTypes[:count]
}

// generateRecommendationForFearType генерирует рекомендацию для типа страха
func (o *ObserverSystem) generateRecommendationForFearType(fearType FearType) ScareRecommendation {
	// Выбираем тип пугающего события на основе типа страха
	scareType := o.chooseScareTypeForFearType(fearType)

	// Вычисляем интенсивность на основе профиля страха и текущего состояния
	intensity := o.fearProfile[fearType] * (1.0 + (1.0 - o.context.CurrentSanity/100.0))

	// Ограничиваем интенсивность в диапазоне [0.3, 1.0]
	intensity = math.Max(0.3, math.Min(1.0, intensity))

	// Вычисляем приоритет на основе эффективности и времени с последнего испуга
	timeFactor := 1.0 - math.Exp(-o.context.TimeSinceLastScare.Seconds()/60.0)
	priority := o.fearProfile[fearType] * (0.5 + timeFactor*0.5)

	// Определяем время ожидания перед испугом
	timing := time.Duration(10 * time.Second)
	if o.context.TimeSinceLastScare < 30*time.Second {
		timing = time.Duration(30 * time.Second)
	}

	// Создаем рекомендацию
	return ScareRecommendation{
		ScareType:  scareType,
		FearTarget: fearType,
		Intensity:  intensity,
		Position:   o.player.Position, // По умолчанию рядом с игроком
		Timing:     timing,
		Priority:   priority,
	}
}

// chooseScareTypeForFearType выбирает тип пугающего события для типа страха
func (o *ObserverSystem) chooseScareTypeForFearType(fearType FearType) ScareEventType {
	switch fearType {
	case FearDarkness:
		return EventEnvironmentChange
	case FearCreatures:
		return EventCreatureAppearance
	case FearSuddenNoises:
		return EventSuddenNoise
	case FearIsolation:
		return EventWhisper
	case FearChasing:
		return EventCreatureAppearance
	case FearGore:
		return EventHallucination
	case FearClaustrophobia:
		return EventEnvironmentChange
	case FearOpenSpaces:
		return EventCreatureAppearance
	case FearUnknown:
		// Выбираем случайный тип события
		eventTypes := []ScareEventType{
			EventAmbientSound,
			EventSuddenNoise,
			EventCreatureAppearance,
			EventEnvironmentChange,
			EventHallucination,
			EventWhisper,
		}
		return eventTypes[o.random.RangeInt(0, len(eventTypes))]
	default:
		return EventAmbientSound
	}
}

// GetScareRecommendation возвращает наилучшую рекомендацию для испуга
func (o *ObserverSystem) GetScareRecommendation() *ScareRecommendation {
	if len(o.recommendedScares) == 0 {
		return nil
	}

	return &o.recommendedScares[0]
}

// GetPlayerFearProfile возвращает профиль страхов игрока
func (o *ObserverSystem) GetPlayerFearProfile() map[FearType]float64 {
	return o.fearProfile
}

// GetPlayerReactorProfile возвращает профиль реактора игрока
func (o *ObserverSystem) GetPlayerReactorProfile() map[ReactorType]float64 {
	return o.reactorProfile
}

// GetDominantFear возвращает доминирующий страх
func (o *ObserverSystem) GetDominantFear() FearType {
	maxValue := -1.0
	dominantFear := FearUnknown

	for fearType, value := range o.fearProfile {
		if value > maxValue {
			maxValue = value
			dominantFear = fearType
		}
	}

	return dominantFear
}

// normalizeProfile нормализует профиль к диапазону [0, 1]
func normalizeProfile(profile map[ReactorType]float64) {
	// Find minimum and maximum values
	min := math.MaxFloat64
	max := -math.MaxFloat64

	for _, value := range profile {
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
	}

	// Normalize values
	for key, value := range profile {
		if max > min {
			profile[key] = (value - min) / (max - min)
		} else {
			profile[key] = 0.5 // If all values are the same
		}
	}
}

// GetFearTypeName возвращает название типа страха
func GetFearTypeName(fearType FearType) string {
	switch fearType {
	case FearDarkness:
		return "Darkness"
	case FearCreatures:
		return "Creatures"
	case FearSuddenNoises:
		return "Sudden Noises"
	case FearIsolation:
		return "Isolation"
	case FearChasing:
		return "Chasing"
	case FearGore:
		return "Gore"
	case FearClaustrophobia:
		return "Claustrophobia"
	case FearOpenSpaces:
		return "Open Spaces"
	case FearUnknown:
		return "Unknown"
	default:
		return "Unknown"
	}
}

// GetReactorTypeName возвращает название типа реактора
func GetReactorTypeName(reactorType ReactorType) string {
	switch reactorType {
	case ReactorCautious:
		return "Cautious"
	case ReactorBold:
		return "Bold"
	case ReactorPanic:
		return "Panic"
	case ReactorMethodical:
		return "Methodical"
	case ReactorReckless:
		return "Reckless"
	case ReactorHesitant:
		return "Hesitant"
	default:
		return "Unknown"
	}
}
