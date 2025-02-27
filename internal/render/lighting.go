package render

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"

	"nightmare/internal/common"
	"nightmare/internal/entity"
	"nightmare/internal/world"
)

// LightType представляет тип источника света
type LightType int

const (
	LightPoint       LightType = iota // Точечный свет
	LightDirectional                  // Направленный свет
	LightSpot                         // Прожектор
	LightAmbient                      // Рассеянный свет
)

// Light представляет источник света
type Light struct {
	Type         LightType
	Position     entity.Vector2D
	Direction    entity.Vector2D
	Color        color.RGBA
	Intensity    float64
	Radius       float64
	Angle        float64 // Для прожектора - угол конуса в радианах
	Flicker      float64 // Интенсивность мерцания (0 - не мерцает)
	FlickerSpeed float64 // Скорость мерцания
	CastShadows  bool
	IsActive     bool
	Falloff      float64 // Коэффициент затухания света
	TimeOffset   float64 // Смещение времени для мерцания
}

// LightingSystem управляет освещением и тенями в игре
type LightingSystem struct {
	lights          []*Light
	ambientLight    color.RGBA
	lightMap        *ebiten.Image // Карта освещения
	shadowMap       *ebiten.Image // Карта теней
	occlusionMap    [][]bool      // Карта преград для света
	screenWidth     int
	screenHeight    int
	lightMapScale   float64 // Масштаб карты освещения относительно экрана
	shadowQuality   int     // Качество теней (количество лучей)
	globalTime      float64 // Глобальное время для анимаций
	collisionSystem *world.CollisionSystem
}

// NewLightingSystem создает новую систему освещения
func NewLightingSystem(width, height int, collisionSystem *world.CollisionSystem) *LightingSystem {
	// Создаем карты освещения и теней
	lightMapScale := 0.5 // По умолчанию карта освещения в два раза меньше экрана
	lightMapWidth := int(float64(width) * lightMapScale)
	lightMapHeight := int(float64(height) * lightMapScale)

	lightMap := ebiten.NewImage(lightMapWidth, lightMapHeight)
	shadowMap := ebiten.NewImage(lightMapWidth, lightMapHeight)

	// Инициализируем карту преград
	occlusionMap := make([][]bool, lightMapHeight)
	for y := range occlusionMap {
		occlusionMap[y] = make([]bool, lightMapWidth)
	}

	return &LightingSystem{
		lights:          make([]*Light, 0),
		ambientLight:    color.RGBA{20, 20, 30, 255}, // Темно-синий фоновый свет
		lightMap:        lightMap,
		shadowMap:       shadowMap,
		occlusionMap:    occlusionMap,
		screenWidth:     width,
		screenHeight:    height,
		lightMapScale:   lightMapScale,
		shadowQuality:   64, // По умолчанию 64 луча для теней
		globalTime:      0,
		collisionSystem: collisionSystem,
	}
}

// AddLight добавляет источник света
func (ls *LightingSystem) AddLight(light *Light) int {
	ls.lights = append(ls.lights, light)
	return len(ls.lights) - 1 // Возвращаем индекс добавленного света
}

// RemoveLight удаляет источник света
func (ls *LightingSystem) RemoveLight(index int) {
	if index >= 0 && index < len(ls.lights) {
		ls.lights = append(ls.lights[:index], ls.lights[index+1:]...)
	}
}

// GetLight возвращает источник света по индексу
func (ls *LightingSystem) GetLight(index int) *Light {
	if index >= 0 && index < len(ls.lights) {
		return ls.lights[index]
	}
	return nil
}

// SetAmbientLight устанавливает фоновое освещение
func (ls *LightingSystem) SetAmbientLight(ambient color.RGBA) {
	ls.ambientLight = ambient
}

// Update обновляет состояние освещения
func (ls *LightingSystem) Update(deltaTime float64) {
	// Обновляем глобальное время
	ls.globalTime += deltaTime

	// Обновляем источники света (мерцание и т.д.)
	for _, light := range ls.lights {
		if light.Flicker > 0 {
			// Вычисляем мерцание
			flickerTime := ls.globalTime*light.FlickerSpeed + light.TimeOffset
			flickerValue := math.Sin(flickerTime) * light.Flicker

			// Применяем мерцание к интенсивности
			light.Intensity = clampFloat64(light.Intensity+flickerValue, 0, 1)
		}
	}

	// Обновляем карту преград
	ls.updateOcclusionMap()
}

// updateOcclusionMap обновляет карту преград для света
func (ls *LightingSystem) updateOcclusionMap() {
	// Если нет системы коллизий, просто выходим
	if ls.collisionSystem == nil {
		return
	}

	// Обновляем карту преград на основе карты коллизий
	for y := 0; y < len(ls.occlusionMap); y++ {
		for x := 0; x < len(ls.occlusionMap[y]); x++ {
			// Преобразуем координаты из карты освещения в координаты мира
			worldX := float64(x) / ls.lightMapScale
			worldY := float64(y) / ls.lightMapScale

			// Проверяем коллизию
			ls.occlusionMap[y][x] = ls.collisionSystem.CheckCollision(common.Vector2D{X: worldX, Y: worldY})
		}
	}
}

// RenderLighting отрисовывает освещение
func (ls *LightingSystem) RenderLighting(screen *ebiten.Image, viewX, viewY float64) *ebiten.Image {
	// Очищаем карту освещения и заполняем ее фоновым светом
	ls.lightMap.Fill(ls.ambientLight)

	// Очищаем карту теней
	ls.shadowMap.Fill(color.RGBA{0, 0, 0, 0})

	// Отрисовываем каждый источник света
	for _, light := range ls.lights {
		if !light.IsActive {
			continue
		}

		// Проверяем, находится ли свет в пределах видимости
		screenX := (light.Position.X-viewX)*ls.lightMapScale + float64(ls.lightMap.Bounds().Dx())/2
		screenY := (light.Position.Y-viewY)*ls.lightMapScale + float64(ls.lightMap.Bounds().Dy())/2

		// Если свет далеко за пределами экрана, пропускаем его
		if screenX < -light.Radius*2 || screenX > float64(ls.lightMap.Bounds().Dx())+light.Radius*2 ||
			screenY < -light.Radius*2 || screenY > float64(ls.lightMap.Bounds().Dy())+light.Radius*2 {
			continue
		}

		// Отрисовываем свет в зависимости от его типа
		switch light.Type {
		case LightPoint:
			ls.renderPointLight(light, screenX, screenY)
		case LightDirectional:
			ls.renderDirectionalLight(light)
		case LightSpot:
			ls.renderSpotLight(light, screenX, screenY)
		}

		// Отрисовываем тени, если источник отбрасывает их
		if light.CastShadows {
			ls.renderShadows(light, screenX, screenY)
		}
	}

	// Создаем итоговое изображение с освещением
	result := ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())

	// Копируем фоновый свет
	op := &ebiten.DrawImageOptions{}
	result.DrawImage(ls.lightMap, op)

	// Накладываем тени
	op = &ebiten.DrawImageOptions{}
	op.CompositeMode = ebiten.CompositeModeLighter
	result.DrawImage(ls.shadowMap, op)

	return result
}

// renderPointLight отрисовывает точечный источник света
func (ls *LightingSystem) renderPointLight(light *Light, screenX, screenY float64) {
	// Создаем градиент для точечного света
	lightImg := ebiten.NewImage(int(light.Radius*2*ls.lightMapScale), int(light.Radius*2*ls.lightMapScale))

	// Параметры света
	centerX := float64(lightImg.Bounds().Dx()) / 2
	centerY := float64(lightImg.Bounds().Dy()) / 2
	radius := light.Radius * ls.lightMapScale

	// Рисуем градиент света
	for y := 0; y < lightImg.Bounds().Dy(); y++ {
		for x := 0; x < lightImg.Bounds().Dx(); x++ {
			// Вычисляем расстояние до центра
			dx := float64(x) - centerX
			dy := float64(y) - centerY
			distance := math.Sqrt(dx*dx + dy*dy)

			// Вычисляем интенсивность на основе расстояния и затухания
			intensity := 1.0 - math.Pow(distance/radius, light.Falloff)
			intensity = clampFloat64(intensity, 0, 1) * light.Intensity

			// Устанавливаем цвет пикселя
			r := uint8(float64(light.Color.R) * intensity)
			g := uint8(float64(light.Color.G) * intensity)
			b := uint8(float64(light.Color.B) * intensity)
			a := uint8(float64(light.Color.A) * intensity)

			lightImg.Set(x, y, color.RGBA{r, g, b, a})
		}
	}

	// Отрисовываем свет на карту освещения
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(screenX-centerX, screenY-centerY)
	op.CompositeMode = ebiten.CompositeModeLighter
	ls.lightMap.DrawImage(lightImg, op)
}

// renderDirectionalLight отрисовывает направленный источник света
func (ls *LightingSystem) renderDirectionalLight(light *Light) {
	// Создаем изображение для направленного света
	lightImg := ebiten.NewImage(ls.lightMap.Bounds().Dx(), ls.lightMap.Bounds().Dy())

	// Вычисляем вектор направления
	dirX := math.Cos(math.Atan2(light.Direction.Y, light.Direction.X))
	dirY := math.Sin(math.Atan2(light.Direction.Y, light.Direction.X))

	// Рисуем направленный свет
	for y := 0; y < lightImg.Bounds().Dy(); y++ {
		for x := 0; x < lightImg.Bounds().Dx(); x++ {
			// Проецируем точку на направление света
			projection := float64(x)*dirX + float64(y)*dirY

			// Вычисляем интенсивность
			intensity := clampFloat64(1.0-projection/float64(lightImg.Bounds().Dx()), 0, 1)
			intensity = math.Pow(intensity, 2) * light.Intensity

			// Устанавливаем цвет пикселя
			r := uint8(float64(light.Color.R) * intensity)
			g := uint8(float64(light.Color.G) * intensity)
			b := uint8(float64(light.Color.B) * intensity)
			a := uint8(float64(light.Color.A) * intensity)

			lightImg.Set(x, y, color.RGBA{r, g, b, a})
		}
	}

	// Отрисовываем свет на карту освещения
	op := &ebiten.DrawImageOptions{}
	op.CompositeMode = ebiten.CompositeModeLighter
	ls.lightMap.DrawImage(lightImg, op)
}

// renderSpotLight отрисовывает прожектор
func (ls *LightingSystem) renderSpotLight(light *Light, screenX, screenY float64) {
	// Создаем изображение для прожектора
	lightImg := ebiten.NewImage(int(light.Radius*2*ls.lightMapScale), int(light.Radius*2*ls.lightMapScale))

	// Параметры света
	centerX := float64(lightImg.Bounds().Dx()) / 2
	centerY := float64(lightImg.Bounds().Dy()) / 2
	radius := light.Radius * ls.lightMapScale

	// Вычисляем направление прожектора
	angle := math.Atan2(light.Direction.Y, light.Direction.X)

	// Рисуем прожектор
	for y := 0; y < lightImg.Bounds().Dy(); y++ {
		for x := 0; x < lightImg.Bounds().Dx(); x++ {
			// Вычисляем расстояние и угол до центра
			dx := float64(x) - centerX
			dy := float64(y) - centerY
			distance := math.Sqrt(dx*dx + dy*dy)
			pixelAngle := math.Atan2(dy, dx)

			// Вычисляем разницу углов
			angleDiff := math.Abs(normalizeAngle(pixelAngle - angle))

			// Если точка вне конуса прожектора, пропускаем
			if angleDiff > light.Angle/2 {
				continue
			}

			// Вычисляем интенсивность на основе расстояния, угла и затухания
			distanceFactor := 1.0 - math.Pow(distance/radius, light.Falloff)
			angleFactor := 1.0 - angleDiff/(light.Angle/2)

			intensity := distanceFactor * angleFactor * light.Intensity
			intensity = clampFloat64(intensity, 0, 1)

			// Устанавливаем цвет пикселя
			r := uint8(float64(light.Color.R) * intensity)
			g := uint8(float64(light.Color.G) * intensity)
			b := uint8(float64(light.Color.B) * intensity)
			a := uint8(float64(light.Color.A) * intensity)

			lightImg.Set(x, y, color.RGBA{r, g, b, a})
		}
	}

	// Отрисовываем свет на карту освещения
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(screenX-centerX, screenY-centerY)
	op.CompositeMode = ebiten.CompositeModeLighter
	ls.lightMap.DrawImage(lightImg, op)
}

// renderShadows отрисовывает тени от источника света
func (ls *LightingSystem) renderShadows(light *Light, screenX, screenY float64) {
	// Если нет системы коллизий, не отрисовываем тени
	if ls.collisionSystem == nil {
		return
	}

	// Создаем изображение для теней
	shadowImg := ebiten.NewImage(ls.shadowMap.Bounds().Dx(), ls.shadowMap.Bounds().Dy())
	shadowImg.Fill(color.RGBA{0, 0, 0, 0})

	// Параметры для расчета теней
	rayCount := ls.shadowQuality
	maxRayLength := light.Radius * ls.lightMapScale

	// Отрисовываем тени для каждого луча
	for i := 0; i < rayCount; i++ {
		// Вычисляем угол луча
		angle := 2 * math.Pi * float64(i) / float64(rayCount)

		// Для прожектора учитываем только лучи в пределах конуса
		if light.Type == LightSpot {
			lightAngle := math.Atan2(light.Direction.Y, light.Direction.X)
			angleDiff := math.Abs(normalizeAngle(angle - lightAngle))
			if angleDiff > light.Angle/2 {
				continue
			}
		}

		// Направление луча
		dirX := math.Cos(angle)
		dirY := math.Sin(angle)

		// Поиск первого препятствия на пути луча
		var hitDistance float64 = maxRayLength

		// Шаг для поиска препятствий
		rayStep := 1.0

		for dist := 0.0; dist < maxRayLength; dist += rayStep {
			// Координаты точки луча
			rayX := screenX + dirX*dist
			rayY := screenY + dirY*dist

			// Проверяем, находится ли точка в пределах карты
			if rayX < 0 || rayX >= float64(ls.lightMap.Bounds().Dx()) ||
				rayY < 0 || rayY >= float64(ls.lightMap.Bounds().Dy()) {
				hitDistance = dist
				break
			}

			// Проверяем, есть ли препятствие в этой точке
			occlusionX := int(rayX)
			occlusionY := int(rayY)

			// Проверяем, что индексы в пределах карты
			if occlusionX >= 0 && occlusionX < len(ls.occlusionMap[0]) &&
				occlusionY >= 0 && occlusionY < len(ls.occlusionMap) {
				if ls.occlusionMap[occlusionY][occlusionX] {
					hitDistance = dist
					break
				}
			}
		}

		// Отрисовываем тень от точки препятствия до края света
		if hitDistance < maxRayLength {

			// Отрисовываем полигон тени
			// В реальной игре здесь был бы код для отрисовки полигона
			// Для простоты просто отрисуем линию
			for d := hitDistance; d < maxRayLength; d += 0.5 {
				x := screenX + dirX*d
				y := screenY + dirY*d

				// Проверяем, что координаты в пределах изображения
				if x >= 0 && x < float64(shadowImg.Bounds().Dx()) &&
					y >= 0 && y < float64(shadowImg.Bounds().Dy()) {
					// Установка цвета с постепенным уменьшением альфы
					alpha := uint8(255 * (1.0 - (d-hitDistance)/(maxRayLength-hitDistance)))
					shadowImg.Set(int(x), int(y), color.RGBA{0, 0, 0, alpha})
				}
			}
		}
	}

	// Отрисовываем тени на карту теней
	op := &ebiten.DrawImageOptions{}
	op.CompositeMode = ebiten.CompositeModeSourceOver
	ls.shadowMap.DrawImage(shadowImg, op)
}

// CreateLight создает новый источник света
func (ls *LightingSystem) CreateLight(lightType LightType, position entity.Vector2D, color color.RGBA, intensity, radius float64) *Light {
	return &Light{
		Type:         lightType,
		Position:     position,
		Direction:    entity.Vector2D{X: 1, Y: 0}, // По умолчанию направлен вправо
		Color:        color,
		Intensity:    intensity,
		Radius:       radius,
		Angle:        math.Pi / 4, // 45 градусов по умолчанию для прожектора
		Flicker:      0,
		FlickerSpeed: 5.0,
		CastShadows:  true,
		IsActive:     true,
		Falloff:      1.5,                  // Квадратичное затухание по умолчанию
		TimeOffset:   rand.Float64() * 100, // Случайное смещение для мерцания
	}
}

// CreateFlashlight создает источник света для фонарика
func (ls *LightingSystem) CreateFlashlight(position entity.Vector2D) *Light {
	return &Light{
		Type:         LightSpot,
		Position:     position,
		Direction:    entity.Vector2D{X: 1, Y: 0},
		Color:        color.RGBA{255, 240, 200, 255}, // Желтоватый свет
		Intensity:    0.8,
		Radius:       15.0,
		Angle:        math.Pi / 3, // 60 градусов
		Flicker:      0.05,        // Небольшое мерцание
		FlickerSpeed: 10.0,
		CastShadows:  true,
		IsActive:     true,
		Falloff:      1.2,
		TimeOffset:   rand.Float64() * 100,
	}
}

// normalizeAngle нормализует угол в пределах [-π, π]
func normalizeAngle(angle float64) float64 {
	for angle > math.Pi {
		angle -= 2 * math.Pi
	}
	for angle < -math.Pi {
		angle += 2 * math.Pi
	}
	return angle
}

// clampFloat64 ограничивает значение в заданном диапазоне
func clampFloat64(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
