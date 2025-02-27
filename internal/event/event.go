package event

import (
	"sync"
	"time"
)

// EventType представляет тип события
type EventType int

// Определение типов событий
const (
	EventPlayerMoved EventType = iota
	EventPlayerDamaged
	EventPlayerInteracted
	EventPlayerSanityChanged
	EventCreatureSpawned
	EventCreatureKilled
	EventCreatureDetected
	EventWorldChanged
	EventScareTriggered
	EventItemPickedUp
	EventItemUsed
	EventAmbientChanged
	EventGameStateChanged
	EventCustom // Для пользовательских событий
)

// EventData представляет данные события
type EventData struct {
	Type      EventType
	Source    interface{}
	Target    interface{}
	Position  interface{} // Позиция, если применимо
	Value     interface{} // Числовое значение, если применимо
	Timestamp time.Time
	Custom    map[string]interface{} // Дополнительные данные
}

// EventCallback представляет функцию обратного вызова для события
type EventCallback func(data EventData)

// EventManager управляет событиями в игре
type EventManager struct {
	listeners       map[EventType][]EventCallback
	customListeners map[string][]EventCallback
	queuedEvents    []EventData
	mutex           sync.RWMutex
}

// NewEventManager создает новый менеджер событий
func NewEventManager() *EventManager {
	return &EventManager{
		listeners:       make(map[EventType][]EventCallback),
		customListeners: make(map[string][]EventCallback),
		queuedEvents:    []EventData{},
	}
}

// AddListener добавляет слушателя события
func (em *EventManager) AddListener(eventType EventType, callback EventCallback) {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	if _, ok := em.listeners[eventType]; !ok {
		em.listeners[eventType] = []EventCallback{}
	}

	em.listeners[eventType] = append(em.listeners[eventType], callback)
}

// AddCustomListener добавляет слушателя пользовательского события
func (em *EventManager) AddCustomListener(eventName string, callback EventCallback) {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	if _, ok := em.customListeners[eventName]; !ok {
		em.customListeners[eventName] = []EventCallback{}
	}

	em.customListeners[eventName] = append(em.customListeners[eventName], callback)
}

// RemoveListener удаляет слушателя события
func (em *EventManager) RemoveListener(eventType EventType, callback EventCallback) {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	if listeners, ok := em.listeners[eventType]; ok {
		for i, cb := range listeners {
			if &cb == &callback {
				// Удаляем элемент
				em.listeners[eventType] = append(listeners[:i], listeners[i+1:]...)
				break
			}
		}
	}
}

// RemoveCustomListener удаляет слушателя пользовательского события
func (em *EventManager) RemoveCustomListener(eventName string, callback EventCallback) {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	if listeners, ok := em.customListeners[eventName]; ok {
		for i, cb := range listeners {
			if &cb == &callback {
				// Удаляем элемент
				em.customListeners[eventName] = append(listeners[:i], listeners[i+1:]...)
				break
			}
		}
	}
}

// Trigger запускает событие
func (em *EventManager) Trigger(eventType EventType, source, target, position, value interface{}) {
	data := EventData{
		Type:      eventType,
		Source:    source,
		Target:    target,
		Position:  position,
		Value:     value,
		Timestamp: time.Now(),
		Custom:    make(map[string]interface{}),
	}

	em.TriggerWithData(data)
}

// TriggerCustom запускает пользовательское событие
func (em *EventManager) TriggerCustom(eventName string, source, target, position, value interface{}, customData map[string]interface{}) {
	data := EventData{
		Type:      EventCustom,
		Source:    source,
		Target:    target,
		Position:  position,
		Value:     value,
		Timestamp: time.Now(),
		Custom:    customData,
	}

	data.Custom["name"] = eventName

	em.TriggerWithData(data)
}

// TriggerWithData запускает событие с данными
func (em *EventManager) TriggerWithData(data EventData) {
	em.mutex.Lock()
	em.queuedEvents = append(em.queuedEvents, data)
	em.mutex.Unlock()
}

// ProcessEvents обрабатывает очередь событий
func (em *EventManager) ProcessEvents() {
	em.mutex.Lock()
	events := em.queuedEvents
	em.queuedEvents = []EventData{}
	em.mutex.Unlock()

	for _, eventData := range events {
		em.dispatchEvent(eventData)
	}
}

// dispatchEvent отправляет событие слушателям
func (em *EventManager) dispatchEvent(data EventData) {
	em.mutex.RLock()
	defer em.mutex.RUnlock()

	// Вызываем слушателей для типа события
	if listeners, ok := em.listeners[data.Type]; ok {
		for _, callback := range listeners {
			callback(data)
		}
	}

	// Если это пользовательское событие, вызываем соответствующих слушателей
	if data.Type == EventCustom {
		if eventName, ok := data.Custom["name"].(string); ok {
			if listeners, ok := em.customListeners[eventName]; ok {
				for _, callback := range listeners {
					callback(data)
				}
			}
		}
	}
}

// NewScareEvent создает новое событие испуга
func NewScareEvent(scareType, source, position interface{}, intensity float64) EventData {
	return EventData{
		Type:      EventScareTriggered,
		Source:    source,
		Position:  position,
		Value:     intensity,
		Timestamp: time.Now(),
		Custom: map[string]interface{}{
			"scareType": scareType,
		},
	}
}

// NewPlayerMovedEvent создает новое событие перемещения игрока
func NewPlayerMovedEvent(player, oldPosition, newPosition interface{}) EventData {
	return EventData{
		Type:      EventPlayerMoved,
		Source:    player,
		Position:  newPosition,
		Timestamp: time.Now(),
		Custom: map[string]interface{}{
			"oldPosition": oldPosition,
		},
	}
}

// NewPlayerDamagedEvent создает новое событие получения урона игроком
func NewPlayerDamagedEvent(player, source interface{}, amount float64) EventData {
	return EventData{
		Type:      EventPlayerDamaged,
		Source:    source,
		Target:    player,
		Value:     amount,
		Timestamp: time.Now(),
	}
}

// NewPlayerInteractedEvent создает новое событие взаимодействия игрока
func NewPlayerInteractedEvent(player, target, position interface{}) EventData {
	return EventData{
		Type:      EventPlayerInteracted,
		Source:    player,
		Target:    target,
		Position:  position,
		Timestamp: time.Now(),
	}
}

// NewPlayerSanityChangedEvent создает новое событие изменения рассудка игрока
func NewPlayerSanityChangedEvent(player interface{}, oldValue, newValue float64) EventData {
	return EventData{
		Type:      EventPlayerSanityChanged,
		Source:    player,
		Value:     newValue - oldValue,
		Timestamp: time.Now(),
		Custom: map[string]interface{}{
			"oldValue": oldValue,
			"newValue": newValue,
		},
	}
}

// NewCreatureSpawnedEvent создает новое событие появления существа
func NewCreatureSpawnedEvent(creature, position interface{}) EventData {
	return EventData{
		Type:      EventCreatureSpawned,
		Source:    creature,
		Position:  position,
		Timestamp: time.Now(),
	}
}

// NewCreatureKilledEvent создает новое событие уничтожения существа
func NewCreatureKilledEvent(creature, killer, position interface{}) EventData {
	return EventData{
		Type:      EventCreatureKilled,
		Source:    killer,
		Target:    creature,
		Position:  position,
		Timestamp: time.Now(),
	}
}

// NewCreatureDetectedEvent создает новое событие обнаружения существа
func NewCreatureDetectedEvent(creature, detector, position interface{}) EventData {
	return EventData{
		Type:      EventCreatureDetected,
		Source:    detector,
		Target:    creature,
		Position:  position,
		Timestamp: time.Now(),
	}
}

// NewWorldChangedEvent создает новое событие изменения мира
func NewWorldChangedEvent(source, position interface{}, radius float64) EventData {
	return EventData{
		Type:      EventWorldChanged,
		Source:    source,
		Position:  position,
		Value:     radius,
		Timestamp: time.Now(),
	}
}

// NewItemPickedUpEvent создает новое событие подбора предмета
func NewItemPickedUpEvent(player, item, position interface{}) EventData {
	return EventData{
		Type:      EventItemPickedUp,
		Source:    player,
		Target:    item,
		Position:  position,
		Timestamp: time.Now(),
	}
}

// NewItemUsedEvent создает новое событие использования предмета
func NewItemUsedEvent(player, item, target, position interface{}) EventData {
	return EventData{
		Type:      EventItemUsed,
		Source:    player,
		Target:    item,
		Position:  position,
		Timestamp: time.Now(),
		Custom: map[string]interface{}{
			"itemTarget": target,
		},
	}
}

// NewAmbientChangedEvent создает новое событие изменения окружения
func NewAmbientChangedEvent(source, oldAmbient, newAmbient interface{}) EventData {
	return EventData{
		Type:      EventAmbientChanged,
		Source:    source,
		Target:    newAmbient,
		Timestamp: time.Now(),
		Custom: map[string]interface{}{
			"oldAmbient": oldAmbient,
		},
	}
}

// NewGameStateChangedEvent создает новое событие изменения состояния игры
func NewGameStateChangedEvent(source, oldState, newState interface{}) EventData {
	return EventData{
		Type:      EventGameStateChanged,
		Source:    source,
		Value:     newState,
		Timestamp: time.Now(),
		Custom: map[string]interface{}{
			"oldState": oldState,
		},
	}
}
