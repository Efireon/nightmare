package sound

import (
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

const (
	sampleRate = 44100 // частота дискретизации звука
)

// SoundType представляет тип звука
type SoundType int

const (
	SoundAmbient SoundType = iota
	SoundEffect
	SoundCreature
	SoundPlayer
	SoundMusic
)

// SoundID представляет идентификатор звука
type SoundID string

// Встроенные звуки
const (
	SoundFootstep     SoundID = "footstep"
	SoundBreath       SoundID = "breath"
	SoundHeartbeat    SoundID = "heartbeat"
	SoundCreak        SoundID = "creak"
	SoundWhisper      SoundID = "whisper"
	SoundScream       SoundID = "scream"
	SoundGrowl        SoundID = "growl"
	SoundRustling     SoundID = "rustling"
	SoundWind         SoundID = "wind"
	SoundThunder      SoundID = "thunder"
	SoundDoor         SoundID = "door"
	SoundCrickets     SoundID = "crickets"
	SoundStaticNoise  SoundID = "static"
	SoundDistantHowl  SoundID = "distant_howl"
	SoundChimes       SoundID = "chimes"
	SoundLaughter     SoundID = "laughter"
	SoundChildrenSong SoundID = "children_song"
	SoundDripping     SoundID = "dripping"
)

// Sound представляет звук
type Sound struct {
	ID       SoundID
	Type     SoundType
	Data     []byte
	Duration time.Duration
	Loop     bool
	Volume   float64
	Pan      float64
}

// SoundInstance представляет экземпляр воспроизводимого звука
type SoundInstance struct {
	Sound      *Sound
	Player     *audio.Player
	StartTime  time.Time
	Position   Vector3D
	Priority   int
	Spatial    bool
	MinDist    float64
	MaxDist    float64
	InstanceID int
}

// Vector3D представляет 3D вектор для пространственного звука
type Vector3D struct {
	X, Y, Z float64
}

// SoundManager управляет звуками в игре
type SoundManager struct {
	audioContext     *audio.Context
	sounds           map[SoundID]*Sound
	activeSounds     map[int]*SoundInstance
	nextInstanceID   int
	listenerPosition Vector3D
	masterVolume     float64

	categoriesVolume map[SoundType]float64

	ambientSounds  []*Sound
	currentAmbient *SoundInstance

	lastHeartbeat time.Time
	heartbeatRate float64 // удары в минуту

	random *rand.Rand
}

// NewSoundManager создает новый менеджер звуков
func NewSoundManager() *SoundManager {
	audioContext := audio.NewContext(sampleRate)

	return &SoundManager{
		audioContext:     audioContext,
		sounds:           make(map[SoundID]*Sound),
		activeSounds:     make(map[int]*SoundInstance),
		nextInstanceID:   1,
		listenerPosition: Vector3D{X: 0, Y: 0, Z: 0},
		masterVolume:     1.0,
		categoriesVolume: map[SoundType]float64{
			SoundAmbient:  1.0,
			SoundEffect:   1.0,
			SoundCreature: 1.0,
			SoundPlayer:   1.0,
			SoundMusic:    0.7,
		},
		ambientSounds: []*Sound{},
		heartbeatRate: 60.0,
		random:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// LoadSound загружает звук из файла
func (sm *SoundManager) LoadSound(id SoundID, filePath string, soundType SoundType, loop bool) error {
	// Здесь мы бы загружали звук из файла
	// В демо-версии просто создаем пустые звуки для примера
	sound := &Sound{
		ID:       id,
		Type:     soundType,
		Data:     make([]byte, 0),
		Duration: 2 * time.Second, // Заглушка
		Loop:     loop,
		Volume:   1.0,
		Pan:      0.0,
	}

	sm.sounds[id] = sound

	// Если это фоновый звук, добавляем его в список фоновых звуков
	if soundType == SoundAmbient {
		sm.ambientSounds = append(sm.ambientSounds, sound)
	}

	return nil
}

// LoadAllSounds загружает все звуки
func (sm *SoundManager) LoadAllSounds() {
	// В реальной игре здесь был бы код загрузки звуков из файлов
	// В демо-версии просто регистрируем звуки

	// Фоновые звуки
	sm.LoadSound(SoundWind, "sound/wind.ogg", SoundAmbient, true)
	sm.LoadSound(SoundCrickets, "sound/crickets.ogg", SoundAmbient, true)
	sm.LoadSound(SoundStaticNoise, "sound/static.ogg", SoundAmbient, true)
	sm.LoadSound(SoundDripping, "sound/dripping.ogg", SoundAmbient, true)

	// Звуки игрока
	sm.LoadSound(SoundFootstep, "sound/footstep.wav", SoundPlayer, false)
	sm.LoadSound(SoundBreath, "sound/breath.wav", SoundPlayer, false)
	sm.LoadSound(SoundHeartbeat, "sound/heartbeat.wav", SoundPlayer, false)

	// Звуки окружения
	sm.LoadSound(SoundCreak, "sound/creak.wav", SoundEffect, false)
	sm.LoadSound(SoundRustling, "sound/rustling.wav", SoundEffect, false)
	sm.LoadSound(SoundThunder, "sound/thunder.wav", SoundEffect, false)
	sm.LoadSound(SoundDoor, "sound/door.wav", SoundEffect, false)
	sm.LoadSound(SoundChimes, "sound/chimes.wav", SoundEffect, false)

	// Звуки существ
	sm.LoadSound(SoundWhisper, "sound/whisper.wav", SoundCreature, false)
	sm.LoadSound(SoundScream, "sound/scream.wav", SoundCreature, false)
	sm.LoadSound(SoundGrowl, "sound/growl.wav", SoundCreature, false)
	sm.LoadSound(SoundDistantHowl, "sound/distant_howl.wav", SoundCreature, false)
	sm.LoadSound(SoundLaughter, "sound/laughter.wav", SoundCreature, false)
	sm.LoadSound(SoundChildrenSong, "sound/children_song.wav", SoundCreature, false)
}

// PlaySound воспроизводит звук
func (sm *SoundManager) PlaySound(id SoundID) int {
	sound, ok := sm.sounds[id]
	if !ok {
		log.Printf("Звук %s не найден", id)
		return 0
	}

	// В реальной игре здесь был бы код создания аудио-плеера
	// В демо-версии просто логируем
	log.Printf("Воспроизведение звука: %s", id)

	// Создаем экземпляр звука
	instance := &SoundInstance{
		Sound:      sound,
		Player:     nil, // В реальной игре здесь был бы аудио-плеер
		StartTime:  time.Now(),
		Position:   Vector3D{X: 0, Y: 0, Z: 0},
		Priority:   1,
		Spatial:    false,
		MinDist:    1.0,
		MaxDist:    50.0,
		InstanceID: sm.nextInstanceID,
	}

	sm.activeSounds[sm.nextInstanceID] = instance
	sm.nextInstanceID++

	return instance.InstanceID
}

// PlaySoundAt воспроизводит пространственный звук
func (sm *SoundManager) PlaySoundAt(id SoundID, position Vector3D, minDist, maxDist float64) int {
	instanceID := sm.PlaySound(id)

	if instanceID > 0 {
		instance := sm.activeSounds[instanceID]
		instance.Position = position
		instance.Spatial = true
		instance.MinDist = minDist
		instance.MaxDist = maxDist

		// Обновляем громкость и панораму
		sm.updateSpatialSound(instance)
	}

	return instanceID
}

// StopSound останавливает звук
func (sm *SoundManager) StopSound(instanceID int) {
	instance, ok := sm.activeSounds[instanceID]
	if !ok {
		return
	}

	// В реальной игре здесь был бы код остановки аудио-плеера
	log.Printf("Остановка звука: %s", instance.Sound.ID)

	delete(sm.activeSounds, instanceID)
}

// StopAllSounds останавливает все звуки
func (sm *SoundManager) StopAllSounds() {
	for id := range sm.activeSounds {
		sm.StopSound(id)
	}
}

// SetListenerPosition устанавливает позицию слушателя
func (sm *SoundManager) SetListenerPosition(position Vector3D) {
	sm.listenerPosition = position

	// Обновляем все пространственные звуки
	for _, instance := range sm.activeSounds {
		if instance.Spatial {
			sm.updateSpatialSound(instance)
		}
	}
}

// updateSpatialSound обновляет пространственный звук
func (sm *SoundManager) updateSpatialSound(instance *SoundInstance) {
	if instance.Player == nil {
		return
	}

	// Вычисляем расстояние до слушателя
	dx := instance.Position.X - sm.listenerPosition.X
	dy := instance.Position.Y - sm.listenerPosition.Y
	dz := instance.Position.Z - sm.listenerPosition.Z

	distance := math.Sqrt(dx*dx + dy*dy + dz*dz)

	// Вычисляем громкость на основе расстояния
	volume := 1.0
	if distance < instance.MinDist {
		volume = 1.0
	} else if distance > instance.MaxDist {
		volume = 0.0
	} else {
		volume = 1.0 - (distance-instance.MinDist)/(instance.MaxDist-instance.MinDist)
	}

	// Учитываем общую громкость
	volume *= sm.masterVolume * sm.categoriesVolume[instance.Sound.Type] * instance.Sound.Volume

	// Вычисляем панораму
	pan := 0.0
	if distance > 0 {
		// Нормализуем вектор направления
		dx /= distance

		// Панорама от -1 до 1 на основе угла
		pan = dx

		// Ограничиваем значение
		if pan < -1.0 {
			pan = -1.0
		} else if pan > 1.0 {
			pan = 1.0
		}
	}

	// Устанавливаем громкость и панораму
	instance.Player.SetVolume(volume)
	instance.Sound.Pan = pan
}

// Update обновляет звуки
func (sm *SoundManager) Update() {
	currentTime := time.Now()

	// Обновляем активные звуки
	for id, instance := range sm.activeSounds {
		// Проверяем, не закончился ли звук
		if !instance.Sound.Loop && currentTime.Sub(instance.StartTime) >= instance.Sound.Duration {
			sm.StopSound(id)
			continue
		}

		// Обновляем пространственные звуки
		if instance.Spatial {
			sm.updateSpatialSound(instance)
		}
	}

	// Управляем фоновыми звуками
	sm.updateAmbientSounds()

	// Управляем сердцебиением
	sm.updateHeartbeat()
}

// updateAmbientSounds обновляет фоновые звуки
func (sm *SoundManager) updateAmbientSounds() {
	// Если нет активного фонового звука, запускаем новый
	if sm.currentAmbient == nil || !sm.isActive(sm.currentAmbient.InstanceID) {
		// Случайно выбираем следующий фоновый звук
		if len(sm.ambientSounds) > 0 {
			sound := sm.ambientSounds[sm.random.Intn(len(sm.ambientSounds))]
			instanceID := sm.PlaySound(sound.ID)
			sm.currentAmbient = sm.activeSounds[instanceID]
		}
	}
}

// updateHeartbeat обновляет звук сердцебиения
func (sm *SoundManager) updateHeartbeat() {
	currentTime := time.Now()

	// Вычисляем интервал между ударами
	interval := time.Duration(60.0/sm.heartbeatRate*1000) * time.Millisecond

	// Если прошло достаточно времени с последнего удара, воспроизводим новый
	if currentTime.Sub(sm.lastHeartbeat) >= interval && sm.heartbeatRate > 0 {
		sm.PlaySound(SoundHeartbeat)
		sm.lastHeartbeat = currentTime
	}
}

// SetHeartbeatRate устанавливает частоту сердцебиения
func (sm *SoundManager) SetHeartbeatRate(bpm float64) {
	sm.heartbeatRate = bpm
}

// SetMasterVolume устанавливает общую громкость
func (sm *SoundManager) SetMasterVolume(volume float64) {
	sm.masterVolume = volume
}

// SetCategoryVolume устанавливает громкость категории звуков
func (sm *SoundManager) SetCategoryVolume(category SoundType, volume float64) {
	sm.categoriesVolume[category] = volume
}

// isActive проверяет, активен ли звук
func (sm *SoundManager) isActive(instanceID int) bool {
	_, ok := sm.activeSounds[instanceID]
	return ok
}

// GenerateRandomSound генерирует случайный звук окружения
func (sm *SoundManager) GenerateRandomSound(position Vector3D) {
	// Список звуков окружения
	environmentSounds := []SoundID{
		SoundCreak, SoundRustling, SoundThunder,
		SoundDoor, SoundChimes, SoundDistantHowl,
	}

	// Выбираем случайный звук
	soundID := environmentSounds[sm.random.Intn(len(environmentSounds))]

	// Воспроизводим звук с некоторой вероятностью
	if sm.random.Float64() < 0.3 {
		sm.PlaySoundAt(soundID, position, 5.0, 50.0)
	}
}

// GenerateScareSound генерирует пугающий звук
func (sm *SoundManager) GenerateScareSound(intensity float64, position Vector3D) {
	// Список пугающих звуков
	scareSounds := []SoundID{
		SoundWhisper, SoundScream, SoundGrowl,
		SoundLaughter, SoundChildrenSong,
	}

	// Выбираем звук в зависимости от интенсивности
	soundIndex := int(float64(len(scareSounds)-1) * intensity)
	soundIndex = clamp(soundIndex, 0, len(scareSounds)-1)

	soundID := scareSounds[soundIndex]

	// Воспроизводим звук
	sm.PlaySoundAt(soundID, position, 10.0, 100.0)
}

// clamp ограничивает значение в диапазоне
func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
