package render

import (
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// EffectType представляет тип визуального эффекта
type EffectType int

const (
	EffectVignette EffectType = iota
	EffectFilmGrain
	EffectChromaticAberration
	EffectDistortion
	EffectBlur
	EffectGlitch
	EffectPulse
	EffectFlash
	EffectFog
	EffectShadow
)

// EffectManager управляет визуальными эффектами
type EffectManager struct {
	effects       map[EffectType]*Effect
	noiseTexture  *ebiten.Image
	glitchTexture *ebiten.Image
	startTime     time.Time
}

// Effect представляет визуальный эффект
type Effect struct {
	Type        EffectType
	Intensity   float64
	Duration    float64
	ElapsedTime float64
	Active      bool
	Params      map[string]float64
}

// NewEffectManager создает новый менеджер эффектов
func NewEffectManager() *EffectManager {
	manager := &EffectManager{
		effects:   make(map[EffectType]*Effect),
		startTime: time.Now(),
	}

	// Инициализируем все типы эффектов
	for i := EffectVignette; i <= EffectShadow; i++ {
		manager.effects[i] = &Effect{
			Type:      i,
			Intensity: 0,
			Duration:  0,
			Active:    false,
			Params:    make(map[string]float64),
		}
	}

	// Генерируем текстуру шума для различных эффектов
	manager.noiseTexture = generateNoiseTexture(256, 256)

	// Генерируем текстуру глюка
	manager.glitchTexture = generateGlitchTexture(32, 32)

	return manager
}

// AddEffect добавляет временный эффект
func (em *EffectManager) AddEffect(effectType EffectType, intensity, duration float64) {
	effect := em.effects[effectType]
	effect.Intensity = intensity
	effect.Duration = duration
	effect.ElapsedTime = 0
	effect.Active = true
}

// SetBaseEffect устанавливает постоянный базовый эффект
func (em *EffectManager) SetBaseEffect(effectType EffectType, intensity float64) {
	effect := em.effects[effectType]
	effect.Intensity = intensity
	effect.Duration = -1 // -1 означает бесконечность
	effect.Active = true
}

// SetEffectParam устанавливает параметр эффекта
func (em *EffectManager) SetEffectParam(effectType EffectType, param string, value float64) {
	effect := em.effects[effectType]
	effect.Params[param] = value
}

// Update обновляет состояние эффектов
func (em *EffectManager) Update(deltaTime float64) {
	// Обновляем таймеры для всех активных эффектов
	for _, effect := range em.effects {
		if effect.Active && effect.Duration > 0 {
			effect.ElapsedTime += deltaTime

			// Если время эффекта истекло
			if effect.ElapsedTime >= effect.Duration {
				effect.Active = false
				effect.Intensity = 0
			}
		}
	}
}

// Apply применяет эффекты к изображению
func (em *EffectManager) Apply(screen *ebiten.Image) {
	// Копируем изображение, чтобы применять эффекты последовательно
	temp := ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())
	temp.DrawImage(screen, nil)

	// Применяем каждый активный эффект
	for _, effect := range em.effects {
		if effect.Active && effect.Intensity > 0 {
			switch effect.Type {
			case EffectVignette:
				em.applyVignette(temp, effect.Intensity)

			case EffectFilmGrain:
				em.applyFilmGrain(temp, effect.Intensity)

			case EffectChromaticAberration:
				em.applyChromaticAberration(temp, effect.Intensity)

			case EffectDistortion:
				em.applyDistortion(temp, effect.Intensity)

			case EffectBlur:
				em.applyBlur(temp, effect.Intensity)

			case EffectGlitch:
				em.applyGlitch(temp, effect.Intensity)

			case EffectPulse:
				em.applyPulse(temp, effect.Intensity)

			case EffectFlash:
				em.applyFlash(temp, effect.Intensity)

			case EffectFog:
				em.applyFog(temp, effect.Intensity)

			case EffectShadow:
				em.applyShadow(temp, effect.Intensity)
			}
		}
	}

	// Очищаем экран и копируем результат
	screen.Clear()
	screen.DrawImage(temp, nil)
}

// applyVignette добавляет эффект виньетки
func (em *EffectManager) applyVignette(img *ebiten.Image, intensity float64) {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	vignette := ebiten.NewImage(width, height)

	// Создаем градиент от центра к краям
	centerX := float64(width) / 2
	centerY := float64(height) / 2
	maxDist := math.Sqrt(centerX*centerX + centerY*centerY)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Вычисляем расстояние от текущей точки до центра
			dx := float64(x) - centerX
			dy := float64(y) - centerY
			dist := math.Sqrt(dx*dx + dy*dy)

			// Нормализуем расстояние и применяем нелинейную функцию
			normalized := math.Pow(dist/maxDist, 2) * intensity

			// Ограничиваем значение
			if normalized > 1 {
				normalized = 1
			}

			// Устанавливаем цвет пикселя
			alpha := uint8(normalized * 255)
			vignette.Set(x, y, color.RGBA{0, 0, 0, alpha})
		}
	}

	// Накладываем виньетку на изображение
	img.DrawImage(vignette, nil)
}

// applyFilmGrain добавляет эффект зернистости
func (em *EffectManager) applyFilmGrain(img *ebiten.Image, intensity float64) {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Создаем шум
	noise := ebiten.NewImage(width, height)

	// Используем предгенерированную текстуру шума
	op := &ebiten.DrawImageOptions{}

	// Меняем текстуру шума каждый кадр
	op.GeoM.Translate(float64(time.Now().UnixNano()%256), float64(time.Now().UnixNano()/256%256))

	// Масштабируем шум
	scale := 4.0
	op.GeoM.Scale(float64(width)/(256*scale), float64(height)/(256*scale))

	// Устанавливаем интенсивность
	op.ColorM.Scale(1, 1, 1, intensity)

	noise.DrawImage(em.noiseTexture, op)

	// Накладываем шум на изображение
	img.DrawImage(noise, nil)
}

// applyChromaticAberration добавляет эффект хроматической аберрации
func (em *EffectManager) applyChromaticAberration(img *ebiten.Image, intensity float64) {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Создаем изображения для каждого цветового канала
	red := ebiten.NewImage(width, height)
	green := ebiten.NewImage(width, height)
	blue := ebiten.NewImage(width, height)

	// Копируем соответствующие каналы
	redOp := &ebiten.DrawImageOptions{}
	redOp.ColorM.Scale(1, 0, 0, 1)
	red.DrawImage(img, redOp)

	greenOp := &ebiten.DrawImageOptions{}
	greenOp.ColorM.Scale(0, 1, 0, 1)
	green.DrawImage(img, greenOp)

	blueOp := &ebiten.DrawImageOptions{}
	blueOp.ColorM.Scale(0, 0, 1, 1)
	blue.DrawImage(img, blueOp)

	// Очищаем изображение
	img.Clear()

	// Накладываем каналы со смещением
	offset := intensity * 5.0

	// Красный канал
	redDrawOp := &ebiten.DrawImageOptions{}
	redDrawOp.GeoM.Translate(-offset, 0)
	img.DrawImage(red, redDrawOp)

	// Зеленый канал - без смещения
	img.DrawImage(green, nil)

	// Синий канал
	blueDrawOp := &ebiten.DrawImageOptions{}
	blueDrawOp.GeoM.Translate(offset, 0)
	img.DrawImage(blue, blueDrawOp)
}

// applyDistortion добавляет эффект искажения
func (em *EffectManager) applyDistortion(img *ebiten.Image, intensity float64) {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Создаем временное изображение
	temp := ebiten.NewImage(width, height)
	temp.DrawImage(img, nil)

	// Очищаем изображение
	img.Clear()

	// Текущее время для анимации
	currentTime := time.Since(em.startTime).Seconds()

	// Применяем искажение
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Вычисляем смещение на основе синусоидальной волны
			angle := float64(x+y)/20.0 + currentTime
			offsetX := math.Sin(angle) * intensity * 10.0
			offsetY := math.Cos(angle) * intensity * 5.0

			// Получаем пиксель из исходного изображения со смещением
			srcX := x + int(offsetX)
			srcY := y + int(offsetY)

			// Проверяем границы
			if srcX >= 0 && srcX < width && srcY >= 0 && srcY < height {
				r, g, b, a := temp.At(srcX, srcY).RGBA()
				img.Set(x, y, color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)})
			}
		}
	}
}

// applyBlur добавляет эффект размытия
func (em *EffectManager) applyBlur(img *ebiten.Image, intensity float64) {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Создаем временное изображение
	temp := ebiten.NewImage(width, height)
	temp.DrawImage(img, nil)

	// Очищаем изображение
	img.Clear()

	// Размер ядра размытия
	kernelSize := int(intensity * 10)
	if kernelSize < 1 {
		kernelSize = 1
	}

	// Применяем размытие (упрощенная реализация)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var rSum, gSum, bSum, aSum, count uint32

			// Суммируем пиксели в окрестности
			for ky := -kernelSize; ky <= kernelSize; ky++ {
				for kx := -kernelSize; kx <= kernelSize; kx++ {
					nx, ny := x+kx, y+ky

					if nx >= 0 && nx < width && ny >= 0 && ny < height {
						r, g, b, a := temp.At(nx, ny).RGBA()
						rSum += r
						gSum += g
						bSum += b
						aSum += a
						count++
					}
				}
			}

			// Вычисляем среднее значение
			if count > 0 {
				img.Set(x, y, color.RGBA{
					uint8((rSum / count) >> 8),
					uint8((gSum / count) >> 8),
					uint8((bSum / count) >> 8),
					uint8((aSum / count) >> 8),
				})
			}
		}
	}
}

// applyGlitch добавляет эффект глюка
func (em *EffectManager) applyGlitch(img *ebiten.Image, intensity float64) {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Создаем временное изображение
	temp := ebiten.NewImage(width, height)
	temp.DrawImage(img, nil)

	// Очищаем изображение
	img.Clear()

	// Генерируем случайные полосы
	numStripes := int(intensity * 10)

	for i := 0; i < numStripes; i++ {
		// Случайная позиция и высота полосы
		y := rand.Intn(height)
		h := 1 + rand.Intn(10)
		xOffset := (rand.Float64()*2 - 1) * intensity * 20

		// Рисуем полосу с глюком
		for j := 0; j < h; j++ {
			if y+j < height {
				for x := 0; x < width; x++ {
					// Смещенная позиция
					srcX := x + int(xOffset)

					// Проверяем границы
					if srcX >= 0 && srcX < width {
						r, g, b, a := temp.At(srcX, y+j).RGBA()
						img.Set(x, y+j, color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)})
					}
				}
			}
		}
	}

	// Восстанавливаем остальные пиксели
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if img.At(x, y) == (color.RGBA{0, 0, 0, 0}) {
				r, g, b, a := temp.At(x, y).RGBA()
				img.Set(x, y, color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)})
			}
		}
	}

	// Накладываем текстуру глюка
	if intensity > 0.5 {
		glitchOp := &ebiten.DrawImageOptions{}
		glitchOp.ColorM.Scale(1, 1, 1, (intensity-0.5)*2*0.3)
		glitchOp.GeoM.Scale(float64(width)/32, float64(height)/32)
		img.DrawImage(em.glitchTexture, glitchOp)
	}
}

// applyPulse добавляет эффект пульсации
func (em *EffectManager) applyPulse(img *ebiten.Image, intensity float64) {
	// Текущее время для анимации
	currentTime := time.Since(em.startTime).Seconds()

	// Вычисляем пульсацию на основе синуса
	pulse := 1.0 + math.Sin(currentTime*5)*intensity*0.2

	// Применяем масштабирование
	op := &ebiten.DrawImageOptions{}

	// Масштабируем относительно центра
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	centerX, centerY := float64(width)/2, float64(height)/2

	op.GeoM.Translate(-centerX, -centerY)
	op.GeoM.Scale(pulse, pulse)
	op.GeoM.Translate(centerX, centerY)

	// Очищаем изображение и рисуем масштабированную версию
	temp := ebiten.NewImage(width, height)
	temp.DrawImage(img, op)
	img.Clear()
	img.DrawImage(temp, nil)
}

// applyFlash добавляет эффект вспышки
func (em *EffectManager) applyFlash(img *ebiten.Image, intensity float64) {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Создаем белое изображение для вспышки
	flash := ebiten.NewImage(width, height)
	flash.Fill(color.RGBA{255, 255, 255, uint8(intensity * 255)})

	// Накладываем вспышку
	img.DrawImage(flash, nil)
}

// applyFog добавляет эффект тумана
func (em *EffectManager) applyFog(img *ebiten.Image, intensity float64) {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Текущее время для анимации
	currentTime := time.Since(em.startTime).Seconds()

	// Создаем туман
	fog := ebiten.NewImage(width, height)

	// Используем предгенерированную текстуру шума
	op := &ebiten.DrawImageOptions{}

	// Анимируем туман
	op.GeoM.Translate(math.Sin(currentTime*0.5)*20, math.Cos(currentTime*0.5)*20)

	// Масштабируем шум
	op.GeoM.Scale(float64(width)/256*2, float64(height)/256*2)

	// Настраиваем цвет тумана
	op.ColorM.Scale(0.8, 0.8, 0.8, intensity)

	fog.DrawImage(em.noiseTexture, op)

	// Накладываем туман
	img.DrawImage(fog, nil)
}

// applyShadow добавляет эффект тени
func (em *EffectManager) applyShadow(img *ebiten.Image, intensity float64) {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Задаем параметры тени
	shadowOffset := intensity * 10
	shadowAlpha := intensity * 0.5

	// Создаем тень
	shadow := ebiten.NewImage(width, height)
	shadow.Fill(color.RGBA{0, 0, 0, uint8(shadowAlpha * 255)})

	// Применяем тень
	shadowOp := &ebiten.DrawImageOptions{}
	shadowOp.GeoM.Translate(shadowOffset, shadowOffset)

	// Рисуем тень под изображением
	temp := ebiten.NewImage(width, height)
	temp.DrawImage(shadow, shadowOp)
	temp.DrawImage(img, nil)

	// Обновляем изображение
	img.Clear()
	img.DrawImage(temp, nil)
}

// generateNoiseTexture создает текстуру шума
func generateNoiseTexture(width, height int) *ebiten.Image {
	img := ebiten.NewImage(width, height)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Генерируем случайное значение для каждого пикселя
			value := uint8(rand.Intn(256))
			img.Set(x, y, color.RGBA{value, value, value, 255})
		}
	}

	return img
}

// generateGlitchTexture создает текстуру "глюка"
func generateGlitchTexture(width, height int) *ebiten.Image {
	img := ebiten.NewImage(width, height)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Генерируем случайные RGB значения
			r := uint8(rand.Intn(256))
			g := uint8(rand.Intn(256))
			b := uint8(rand.Intn(256))

			// 50% вероятность прозрачности
			a := uint8(0)
			if rand.Float64() < 0.5 {
				a = 255
			}

			img.Set(x, y, color.RGBA{r, g, b, a})
		}
	}

	return img
}
