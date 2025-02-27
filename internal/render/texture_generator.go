package render

import (
	"image"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"

	"nightmare/internal/util"
)

// TextureType представляет тип текстуры
type TextureType int

const (
	TextureGrass TextureType = iota
	TextureForest
	TextureDenseForest
	TexturePath
	TextureRocks
	TextureWater
	TextureSwamp
	TextureCorrupted
	TextureWall
	TextureFloor
	TextureCreature
	TextureEffect
)

// TextureGenerator создает процедурно генерируемые текстуры
type TextureGenerator struct {
	noise        *util.NoiseGenerator
	rand         *rand.Rand
	baseTextures map[TextureType]*ebiten.Image
	tileSize     int
}

// NewTextureGenerator создает новый генератор текстур
func NewTextureGenerator(seed int64, tileSize int) *TextureGenerator {
	if seed == 0 {
		seed = rand.Int63()
	}

	return &TextureGenerator{
		noise:        util.NewNoiseGenerator(seed),
		rand:         rand.New(rand.NewSource(seed)),
		baseTextures: make(map[TextureType]*ebiten.Image),
		tileSize:     tileSize,
	}
}

// Initialize инициализирует генератор текстур
func (tg *TextureGenerator) Initialize() {
	// Создаем базовые текстуры для всех типов
	for i := TextureGrass; i <= TextureEffect; i++ {
		tg.generateBaseTexture(i)
	}
}

// generateBaseTexture создает базовую текстуру указанного типа
func (tg *TextureGenerator) generateBaseTexture(textureType TextureType) {
	// Создаем новую текстуру
	texture := ebiten.NewImage(tg.tileSize, tg.tileSize)

	// Заполняем текстуру в зависимости от типа
	switch textureType {
	case TextureGrass:
		tg.generateGrassTexture(texture)
	case TextureForest:
		tg.generateForestTexture(texture)
	case TextureDenseForest:
		tg.generateDenseForestTexture(texture)
	case TexturePath:
		tg.generatePathTexture(texture)
	case TextureRocks:
		tg.generateRocksTexture(texture)
	case TextureWater:
		tg.generateWaterTexture(texture)
	case TextureSwamp:
		tg.generateSwampTexture(texture)
	case TextureCorrupted:
		tg.generateCorruptedTexture(texture)
	case TextureWall:
		tg.generateWallTexture(texture)
	case TextureFloor:
		tg.generateFloorTexture(texture)
	case TextureCreature:
		tg.generateCreatureTexture(texture)
	case TextureEffect:
		tg.generateEffectTexture(texture)
	}

	// Сохраняем текстуру
	tg.baseTextures[textureType] = texture
}

// generateGrassTexture создает текстуру травы
func (tg *TextureGenerator) generateGrassTexture(texture *ebiten.Image) {
	// Создаем базовый цвет травы
	baseColor := color.RGBA{0, 180, 0, 255}

	// Заполняем текстуру шумом
	tg.fillWithColorNoise(texture, baseColor, 20)

	// Добавляем детали травы
	tg.addGrassDetails(texture, 30)
}

// generateForestTexture создает текстуру леса
func (tg *TextureGenerator) generateForestTexture(texture *ebiten.Image) {
	// Создаем базовый цвет леса
	baseColor := color.RGBA{0, 120, 0, 255}

	// Заполняем текстуру шумом
	tg.fillWithColorNoise(texture, baseColor, 30)

	// Добавляем детали деревьев
	tg.addTreeDetails(texture, 5)
}

// generateDenseForestTexture создает текстуру густого леса
func (tg *TextureGenerator) generateDenseForestTexture(texture *ebiten.Image) {
	// Создаем базовый цвет густого леса
	baseColor := color.RGBA{0, 80, 0, 255}

	// Заполняем текстуру шумом
	tg.fillWithColorNoise(texture, baseColor, 40)

	// Добавляем детали деревьев
	tg.addTreeDetails(texture, 10)
}

// generatePathTexture создает текстуру тропинки
func (tg *TextureGenerator) generatePathTexture(texture *ebiten.Image) {
	// Создаем базовый цвет тропинки
	baseColor := color.RGBA{200, 190, 140, 255}

	// Заполняем текстуру шумом
	tg.fillWithColorNoise(texture, baseColor, 30)

	// Добавляем детали тропинки (камни, следы)
	tg.addPathDetails(texture, 10)
}

// generateRocksTexture создает текстуру камней
func (tg *TextureGenerator) generateRocksTexture(texture *ebiten.Image) {
	// Создаем базовый цвет камней
	baseColor := color.RGBA{120, 120, 120, 255}

	// Заполняем текстуру шумом
	tg.fillWithColorNoise(texture, baseColor, 50)

	// Добавляем детали камней
	tg.addRockDetails(texture, 20)
}

// generateWaterTexture создает текстуру воды
func (tg *TextureGenerator) generateWaterTexture(texture *ebiten.Image) {
	// Создаем базовый цвет воды
	baseColor := color.RGBA{0, 0, 180, 255}

	// Заполняем текстуру шумом
	tg.fillWithColorNoise(texture, baseColor, 20)

	// Добавляем детали воды (волны, отражение)
	tg.addWaterDetails(texture, 40)
}

// generateSwampTexture создает текстуру болота
func (tg *TextureGenerator) generateSwampTexture(texture *ebiten.Image) {
	// Создаем базовый цвет болота
	baseColor := color.RGBA{70, 90, 70, 255}

	// Заполняем текстуру шумом
	tg.fillWithColorNoise(texture, baseColor, 50)

	// Добавляем детали болота
	tg.addSwampDetails(texture, 30)
}

// generateCorruptedTexture создает текстуру искаженной местности
func (tg *TextureGenerator) generateCorruptedTexture(texture *ebiten.Image) {
	// Создаем базовый цвет искаженной местности
	baseColor := color.RGBA{80, 0, 80, 255}

	// Заполняем текстуру искаженным шумом
	tg.fillWithDistortedNoise(texture, baseColor, 60)

	// Добавляем детали искажения
	tg.addCorruptionDetails(texture, 50)
}

// generateWallTexture создает текстуру стены
func (tg *TextureGenerator) generateWallTexture(texture *ebiten.Image) {
	// Создаем базовый цвет стены
	baseColor := color.RGBA{100, 100, 100, 255}

	// Заполняем текстуру шумом
	tg.fillWithColorNoise(texture, baseColor, 30)

	// Добавляем детали стены (трещины, кирпичи)
	tg.addWallDetails(texture, 40)
}

// generateFloorTexture создает текстуру пола
func (tg *TextureGenerator) generateFloorTexture(texture *ebiten.Image) {
	// Создаем базовый цвет пола
	baseColor := color.RGBA{130, 110, 80, 255}

	// Заполняем текстуру шумом
	tg.fillWithColorNoise(texture, baseColor, 20)

	// Добавляем детали пола (доски, плитка)
	tg.addFloorDetails(texture, 30)
}

// generateCreatureTexture создает текстуру существа
func (tg *TextureGenerator) generateCreatureTexture(texture *ebiten.Image) {
	// Создаем базовый цвет существа
	baseColor := color.RGBA{50, 50, 50, 255}

	// Заполняем текстуру органическим шумом
	tg.fillWithOrganicNoise(texture, baseColor, 50)

	// Добавляем детали существа
	tg.addCreatureDetails(texture, 60)
}

// generateEffectTexture создает текстуру эффекта
func (tg *TextureGenerator) generateEffectTexture(texture *ebiten.Image) {
	// Создаем базовый цвет эффекта
	baseColor := color.RGBA{200, 200, 200, 150}

	// Заполняем текстуру эфирным шумом
	tg.fillWithEtherealNoise(texture, baseColor, 70)

	// Добавляем детали эффекта
	tg.addEffectDetails(texture, 80)
}

// GetBaseTexture возвращает базовую текстуру указанного типа
func (tg *TextureGenerator) GetBaseTexture(textureType TextureType) *ebiten.Image {
	return tg.baseTextures[textureType]
}

// GenerateUniqueTexture создает уникальную текстуру на основе базовой
func (tg *TextureGenerator) GenerateUniqueTexture(textureType TextureType, seed int64) *ebiten.Image {
	// Получаем базовую текстуру
	baseTexture := tg.GetBaseTexture(textureType)
	if baseTexture == nil {
		return nil
	}

	// Создаем новую текстуру
	texture := ebiten.NewImage(tg.tileSize, tg.tileSize)

	// Копируем базовую текстуру
	texture.DrawImage(baseTexture, nil)

	// Модифицируем текстуру с использованием seed
	tg.modifyTexture(texture, textureType, seed)

	return texture
}

// GenerateCorruptedTexture создает искаженную версию текстуры
func (tg *TextureGenerator) GenerateCorruptedTexture(baseTexture *ebiten.Image, corruptionLevel float64) *ebiten.Image {
	if baseTexture == nil {
		return nil
	}

	// Создаем новую текстуру
	texture := ebiten.NewImage(tg.tileSize, tg.tileSize)

	// Копируем базовую текстуру
	texture.DrawImage(baseTexture, nil)

	// Применяем искажение
	tg.applyCorruption(texture, corruptionLevel)

	return texture
}

// fillWithColorNoise заполняет текстуру цветовым шумом
func (tg *TextureGenerator) fillWithColorNoise(texture *ebiten.Image, baseColor color.RGBA, variance int) {
	bounds := texture.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Генерируем шум для каждого компонента цвета
			noise := tg.noise.Perlin2D(float64(x)/float64(width)*3, float64(y)/float64(height)*3, 0.5)

			// Создаем вариацию цвета на основе шума
			r := clampUint8(int(baseColor.R)+int(noise*float64(variance)), 0, 255)
			g := clampUint8(int(baseColor.G)+int(noise*float64(variance)), 0, 255)
			b := clampUint8(int(baseColor.B)+int(noise*float64(variance)), 0, 255)

			// Устанавливаем цвет пикселя
			texture.Set(x, y, color.RGBA{r, g, b, baseColor.A})
		}
	}
}

// fillWithDistortedNoise заполняет текстуру искаженным шумом
func (tg *TextureGenerator) fillWithDistortedNoise(texture *ebiten.Image, baseColor color.RGBA, variance int) {
	bounds := texture.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Генерируем искаженные координаты
			distX, distY := tg.noise.DomainWarp(float64(x)/float64(width)*3, float64(y)/float64(height)*3, 2.0, 0.5)

			// Генерируем шум на основе искаженных координат
			noise := tg.noise.Perlin2D(distX*3, distY*3, 0.5)

			// Создаем вариацию цвета на основе шума
			r := clampUint8(int(baseColor.R)+int(noise*float64(variance)), 0, 255)
			g := clampUint8(int(baseColor.G)+int(noise*float64(variance)), 0, 255)
			b := clampUint8(int(baseColor.B)+int(noise*float64(variance)), 0, 255)

			// Устанавливаем цвет пикселя
			texture.Set(x, y, color.RGBA{r, g, b, baseColor.A})
		}
	}
}

// fillWithOrganicNoise заполняет текстуру органическим шумом
func (tg *TextureGenerator) fillWithOrganicNoise(texture *ebiten.Image, baseColor color.RGBA, variance int) {
	bounds := texture.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Используем смешанные типы шума для органического эффекта
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Ворлей (ячеистый) шум для органического вида
			cellNoise := tg.noise.WorleyNoise(float64(x)/float64(width)*5, float64(y)/float64(height)*5, 10, 1.0)

			// Перлин шум для плавности
			perlinNoise := tg.noise.Perlin2D(float64(x)/float64(width)*3, float64(y)/float64(height)*3, 0.5)

			// Смешиваем шумы
			mixedNoise := cellNoise*0.6 + perlinNoise*0.4

			// Создаем вариацию цвета на основе шума
			r := clampUint8(int(baseColor.R)+int(mixedNoise*float64(variance)), 0, 255)
			g := clampUint8(int(baseColor.G)+int(mixedNoise*float64(variance)), 0, 255)
			b := clampUint8(int(baseColor.B)+int(mixedNoise*float64(variance)), 0, 255)

			// Устанавливаем цвет пикселя
			texture.Set(x, y, color.RGBA{r, g, b, baseColor.A})
		}
	}
}

// fillWithEtherealNoise заполняет текстуру эфирным шумом
func (tg *TextureGenerator) fillWithEtherealNoise(texture *ebiten.Image, baseColor color.RGBA, variance int) {
	bounds := texture.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Генерируем несколько октав шума
			noise1 := tg.noise.Perlin2D(float64(x)/float64(width)*2, float64(y)/float64(height)*2, 0.5)
			noise2 := tg.noise.Perlin2D(float64(x)/float64(width)*4, float64(y)/float64(height)*4, 0.5)
			noise3 := tg.noise.Perlin2D(float64(x)/float64(width)*8, float64(y)/float64(height)*8, 0.5)

			// Смешиваем шумы с разными весами
			mixedNoise := noise1*0.5 + noise2*0.3 + noise3*0.2

			// Создаем вариацию цвета на основе шума
			r := clampUint8(int(baseColor.R)+int(mixedNoise*float64(variance)), 0, 255)
			g := clampUint8(int(baseColor.G)+int(mixedNoise*float64(variance)), 0, 255)
			b := clampUint8(int(baseColor.B)+int(mixedNoise*float64(variance)), 0, 255)

			// Вычисляем альфа-канал на основе шума для эфирного эффекта
			a := clampUint8(int(baseColor.A)+int((mixedNoise-0.5)*100), 0, 255)

			// Устанавливаем цвет пикселя
			texture.Set(x, y, color.RGBA{r, g, b, a})
		}
	}
}

// addGrassDetails добавляет детали травы
func (tg *TextureGenerator) addGrassDetails(texture *ebiten.Image, count int) {
	bounds := texture.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Добавляем случайные травинки
	for i := 0; i < count; i++ {
		x := tg.rand.Intn(width)
		y := tg.rand.Intn(height)
		length := 1 + tg.rand.Intn(3)
		angle := tg.rand.Float64() * math.Pi

		// Цвет травинки
		r := uint8(0 + tg.rand.Intn(40))
		g := uint8(160 + tg.rand.Intn(70))
		b := uint8(0 + tg.rand.Intn(30))

		// Рисуем травинку
		for j := 0; j < length; j++ {
			dx := int(math.Cos(angle) * float64(j))
			dy := int(math.Sin(angle) * float64(j))

			nx, ny := x+dx, y+dy
			if nx >= 0 && nx < width && ny >= 0 && ny < height {
				texture.Set(nx, ny, color.RGBA{r, g, b, 255})
			}
		}
	}
}

// addTreeDetails добавляет детали деревьев
func (tg *TextureGenerator) addTreeDetails(texture *ebiten.Image, count int) {
	bounds := texture.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Добавляем маленькие деревья
	for i := 0; i < count; i++ {
		x := tg.rand.Intn(width)
		y := tg.rand.Intn(height)
		size := 1 + tg.rand.Intn(3)

		// Цвет ствола
		trunkR := uint8(80 + tg.rand.Intn(40))
		trunkG := uint8(40 + tg.rand.Intn(30))
		trunkB := uint8(0 + tg.rand.Intn(20))

		// Цвет кроны
		crownR := uint8(0 + tg.rand.Intn(30))
		crownG := uint8(60 + tg.rand.Intn(40))
		crownB := uint8(0 + tg.rand.Intn(20))

		// Рисуем ствол
		for j := 0; j < size; j++ {
			nx, ny := x, y+j
			if nx >= 0 && nx < width && ny >= 0 && ny < height {
				texture.Set(nx, ny, color.RGBA{trunkR, trunkG, trunkB, 255})
			}
		}

		// Рисуем крону
		for j := -size; j <= size; j++ {
			for k := -size; k <= size; k++ {
				if j*j+k*k <= size*size {
					nx, ny := x+j, y-size+k
					if nx >= 0 && nx < width && ny >= 0 && ny < height {
						texture.Set(nx, ny, color.RGBA{crownR, crownG, crownB, 255})
					}
				}
			}
		}
	}
}

// addPathDetails добавляет детали тропинки
func (tg *TextureGenerator) addPathDetails(texture *ebiten.Image, count int) {
	bounds := texture.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Добавляем маленькие камни и следы
	for i := 0; i < count; i++ {
		x := tg.rand.Intn(width)
		y := tg.rand.Intn(height)
		size := 1 + tg.rand.Intn(2)

		// Камень или след
		if tg.rand.Float64() < 0.7 {
			// Цвет камня
			r := uint8(160 + tg.rand.Intn(40))
			g := uint8(160 + tg.rand.Intn(40))
			b := uint8(160 + tg.rand.Intn(40))

			// Рисуем камень
			for j := -size; j <= size; j++ {
				for k := -size; k <= size; k++ {
					if j*j+k*k <= size*size {
						nx, ny := x+j, y+k
						if nx >= 0 && nx < width && ny >= 0 && ny < height {
							texture.Set(nx, ny, color.RGBA{r, g, b, 255})
						}
					}
				}
			}
		} else {
			// Цвет следа (темнее основного цвета)
			r := uint8(140 + tg.rand.Intn(30))
			g := uint8(130 + tg.rand.Intn(30))
			b := uint8(100 + tg.rand.Intn(20))

			// Рисуем след (овал)
			angle := tg.rand.Float64() * math.Pi
			for j := -size; j <= size; j++ {
				for k := -size / 2; k <= size/2; k++ {
					// Поворачиваем точку
					rotJ := float64(j)*math.Cos(angle) - float64(k)*math.Sin(angle)
					rotK := float64(j)*math.Sin(angle) + float64(k)*math.Cos(angle)

					nx, ny := x+int(rotJ), y+int(rotK)
					if nx >= 0 && nx < width && ny >= 0 && ny < height {
						if j*j/(size*size)+k*k/((size/2)*(size/2)) <= 1 {
							texture.Set(nx, ny, color.RGBA{r, g, b, 255})
						}
					}
				}
			}
		}
	}
}

// addRockDetails добавляет детали камней
func (tg *TextureGenerator) addRockDetails(texture *ebiten.Image, count int) {
	bounds := texture.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Добавляем трещины и выпуклости
	for i := 0; i < count; i++ {
		x := tg.rand.Intn(width)
		y := tg.rand.Intn(height)
		length := 2 + tg.rand.Intn(5)
		angle := tg.rand.Float64() * math.Pi

		// Трещина или выпуклость
		if tg.rand.Float64() < 0.6 {
			// Цвет трещины (темнее)
			r := uint8(60 + tg.rand.Intn(30))
			g := uint8(60 + tg.rand.Intn(30))
			b := uint8(60 + tg.rand.Intn(30))

			// Рисуем трещину
			for j := 0; j < length; j++ {
				// Добавляем небольшую случайность для реалистичности
				jitter := tg.rand.Float64()*0.8 - 0.4
				dx := int(math.Cos(angle+jitter) * float64(j))
				dy := int(math.Sin(angle+jitter) * float64(j))

				nx, ny := x+dx, y+dy
				if nx >= 0 && nx < width && ny >= 0 && ny < height {
					texture.Set(nx, ny, color.RGBA{r, g, b, 255})
				}
			}
		} else {
			// Цвет выпуклости (светлее)
			r := uint8(150 + tg.rand.Intn(50))
			g := uint8(150 + tg.rand.Intn(50))
			b := uint8(150 + tg.rand.Intn(50))

			// Размер выпуклости
			size := 1 + tg.rand.Intn(2)

			// Рисуем выпуклость
			for j := -size; j <= size; j++ {
				for k := -size; k <= size; k++ {
					if j*j+k*k <= size*size {
						nx, ny := x+j, y+k
						if nx >= 0 && nx < width && ny >= 0 && ny < height {
							texture.Set(nx, ny, color.RGBA{r, g, b, 255})
						}
					}
				}
			}
		}
	}
}

// addWaterDetails добавляет детали воды
func (tg *TextureGenerator) addWaterDetails(texture *ebiten.Image, count int) {
	bounds := texture.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Добавляем блики и волны
	for i := 0; i < count; i++ {
		x := tg.rand.Intn(width)
		y := tg.rand.Intn(height)

		// Блик или волна
		if tg.rand.Float64() < 0.3 {
			// Цвет блика (светлее)
			r := uint8(180 + tg.rand.Intn(75))
			g := uint8(180 + tg.rand.Intn(75))
			b := uint8(220 + tg.rand.Intn(35))

			// Размер блика
			size := 1 + tg.rand.Intn(2)

			// Рисуем блик
			for j := -size; j <= size; j++ {
				for k := -size; k <= size; k++ {
					if j*j+k*k <= size*size {
						nx, ny := x+j, y+k
						if nx >= 0 && nx < width && ny >= 0 && ny < height {
							alpha := 255 - uint8((j*j+k*k)*255/(size*size*2)) // Прозрачнее к краям
							texture.Set(nx, ny, color.RGBA{r, g, b, alpha})
						}
					}
				}
			}
		} else {
			// Цвет волны
			r := uint8(0 + tg.rand.Intn(20))
			g := uint8(0 + tg.rand.Intn(20))
			b := uint8(150 + tg.rand.Intn(60))

			// Длина волны
			length := 3 + tg.rand.Intn(5)

			// Рисуем волну (синусоидальная линия)
			amplitude := 0.5 + tg.rand.Float64()*1.5
			frequency := 0.3 + tg.rand.Float64()*0.7
			angle := tg.rand.Float64() * math.Pi

			for j := 0; j < length; j++ {
				// Синусоидальное смещение
				offset := amplitude * math.Sin(float64(j)*frequency)

				dx := int(math.Cos(angle)*float64(j) - math.Sin(angle)*offset)
				dy := int(math.Sin(angle)*float64(j) + math.Cos(angle)*offset)

				nx, ny := x+dx, y+dy
				if nx >= 0 && nx < width && ny >= 0 && ny < height {
					texture.Set(nx, ny, color.RGBA{r, g, b, 255})
				}
			}
		}
	}
}

// addSwampDetails добавляет детали болота
func (tg *TextureGenerator) addSwampDetails(texture *ebiten.Image, count int) {
	bounds := texture.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Добавляем пузыри и водоросли
	for i := 0; i < count; i++ {
		x := tg.rand.Intn(width)
		y := tg.rand.Intn(height)

		// Пузырь или водоросль
		if tg.rand.Float64() < 0.4 {
			// Цвет пузыря
			r := uint8(150 + tg.rand.Intn(50))
			g := uint8(150 + tg.rand.Intn(50))
			b := uint8(150 + tg.rand.Intn(50))

			// Размер пузыря
			size := 1 + tg.rand.Intn(2)

			// Рисуем пузырь
			for j := -size; j <= size; j++ {
				for k := -size; k <= size; k++ {
					if j*j+k*k <= size*size {
						nx, ny := x+j, y+k
						if nx >= 0 && nx < width && ny >= 0 && ny < height {
							// Край пузыря
							if float64(size-1) <= math.Sqrt(float64(j*j+k*k)) && math.Sqrt(float64(j*j+k*k)) <= float64(size) {
								texture.Set(nx, ny, color.RGBA{r, g, b, 255})
							}
						}
					}
				}
			}
		} else {
			// Цвет водоросли
			r := uint8(20 + tg.rand.Intn(30))
			g := uint8(100 + tg.rand.Intn(50))
			b := uint8(20 + tg.rand.Intn(30))

			// Длина водоросли
			length := 2 + tg.rand.Intn(4)

			// Рисуем водоросль (волнистая линия)
			amplitude := 0.5 + tg.rand.Float64()
			frequency := 0.5 + tg.rand.Float64()

			for j := 0; j < length; j++ {
				// Синусоидальное смещение
				dx := int(amplitude * math.Sin(float64(j)*frequency))

				nx, ny := x+dx, y-j
				if nx >= 0 && nx < width && ny >= 0 && ny < height {
					texture.Set(nx, ny, color.RGBA{r, g, b, 255})
				}
			}
		}
	}
}

// addCorruptionDetails добавляет детали искажения
func (tg *TextureGenerator) addCorruptionDetails(texture *ebiten.Image, count int) {
	bounds := texture.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Добавляем странные формы и "вены"
	for i := 0; i < count; i++ {
		x := tg.rand.Intn(width)
		y := tg.rand.Intn(height)

		// Форма или "вена"
		if tg.rand.Float64() < 0.3 {
			// Цвет формы
			r := uint8(100 + tg.rand.Intn(100))
			g := uint8(0 + tg.rand.Intn(50))
			b := uint8(100 + tg.rand.Intn(100))

			// Размер и форма
			size := 2 + tg.rand.Intn(3)
			deformFactor := 0.3 + tg.rand.Float64()*0.7

			// Рисуем деформированную форму
			for j := -size; j <= size; j++ {
				for k := -size; k <= size; k++ {
					// Деформируем круг
					deform := 1.0 + deformFactor*tg.noise.Perlin2D(float64(j+x)*0.2, float64(k+y)*0.2, 0.5)
					if (j*j + k*k) <= int(float64(size*size)*deform) {
						nx, ny := x+j, y+k
						if nx >= 0 && nx < width && ny >= 0 && ny < height {
							texture.Set(nx, ny, color.RGBA{r, g, b, 255})
						}
					}
				}
			}
		} else {
			// Цвет "вены"
			r := uint8(150 + tg.rand.Intn(105))
			g := uint8(0 + tg.rand.Intn(30))
			b := uint8(150 + tg.rand.Intn(105))

			// Длина "вены"
			length := 5 + tg.rand.Intn(10)
			angle := tg.rand.Float64() * 2 * math.Pi

			// Рисуем "вену" (извилистая линия)
			prevX, prevY := x, y
			for j := 0; j < length; j++ {
				// Случайное изменение направления
				angle += (tg.rand.Float64() - 0.5) * math.Pi / 2

				// Вычисляем новую позицию
				dx := int(math.Cos(angle))
				dy := int(math.Sin(angle))

				nx, ny := prevX+dx, prevY+dy
				if nx >= 0 && nx < width && ny >= 0 && ny < height {
					texture.Set(nx, ny, color.RGBA{r, g, b, 255})

					// Иногда добавляем небольшие ответвления
					if tg.rand.Float64() < 0.2 {
						branchAngle := angle + math.Pi/2*(tg.rand.Float64()*2-1)
						branchLength := 1 + tg.rand.Intn(3)

						for k := 0; k < branchLength; k++ {
							bx := int(math.Cos(branchAngle) * float64(k))
							by := int(math.Sin(branchAngle) * float64(k))

							bnx, bny := nx+bx, ny+by
							if bnx >= 0 && bnx < width && bny >= 0 && bny < height {
								texture.Set(bnx, bny, color.RGBA{r, g, b, 255})
							}
						}
					}

					prevX, prevY = nx, ny
				}
			}
		}
	}
}

// addWallDetails добавляет детали стены
func (tg *TextureGenerator) addWallDetails(texture *ebiten.Image, count int) {
	bounds := texture.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Рисуем кирпичи или трещины
	if tg.rand.Float64() < 0.5 {
		// Кирпичи
		brickWidth := 4 + tg.rand.Intn(3)
		brickHeight := 2 + tg.rand.Intn(2)

		// Цвет кирпича
		r := uint8(120 + tg.rand.Intn(40))
		g := uint8(70 + tg.rand.Intn(30))
		b := uint8(50 + tg.rand.Intn(30))

		// Цвет шва
		mortarR := uint8(160 + tg.rand.Intn(40))
		mortarG := uint8(160 + tg.rand.Intn(40))
		mortarB := uint8(160 + tg.rand.Intn(40))

		// Рисуем кирпичную стену
		for y := 0; y < height; y += brickHeight {
			// Смещение для чередования кирпичей
			offset := 0
			if (y/brickHeight)%2 == 1 {
				offset = brickWidth / 2
			}

			for x := -offset; x < width; x += brickWidth {
				// Рисуем кирпич
				for j := 0; j < brickHeight-1; j++ {
					for k := 0; k < brickWidth-1; k++ {
						nx, ny := x+k, y+j
						if nx >= 0 && nx < width && ny >= 0 && ny < height {
							// Вариация цвета кирпича
							variation := tg.rand.Intn(20) - 10
							brickR := clampUint8(int(r)+variation, 0, 255)
							brickG := clampUint8(int(g)+variation, 0, 255)
							brickB := clampUint8(int(b)+variation, 0, 255)

							texture.Set(nx, ny, color.RGBA{brickR, brickG, brickB, 255})
						}
					}
				}

				// Рисуем швы
				for j := brickHeight - 1; j < brickHeight; j++ {
					for k := 0; k < brickWidth; k++ {
						nx, ny := x+k, y+j
						if nx >= 0 && nx < width && ny >= 0 && ny < height {
							texture.Set(nx, ny, color.RGBA{mortarR, mortarG, mortarB, 255})
						}
					}
				}

				for j := 0; j < brickHeight; j++ {
					for k := brickWidth - 1; k < brickWidth; k++ {
						nx, ny := x+k, y+j
						if nx >= 0 && nx < width && ny >= 0 && ny < height {
							texture.Set(nx, ny, color.RGBA{mortarR, mortarG, mortarB, 255})
						}
					}
				}
			}
		}
	} else {
		// Трещины
		for i := 0; i < count; i++ {
			x := tg.rand.Intn(width)
			y := tg.rand.Intn(height)
			length := 3 + tg.rand.Intn(7)
			angle := tg.rand.Float64() * math.Pi

			// Цвет трещины (темнее)
			r := uint8(60 + tg.rand.Intn(20))
			g := uint8(60 + tg.rand.Intn(20))
			b := uint8(60 + tg.rand.Intn(20))

			// Рисуем трещину с ответвлениями
			tg.drawCrack(texture, x, y, length, angle, r, g, b, 1)
		}
	}
}

// drawCrack рисует трещину с возможными ответвлениями
func (tg *TextureGenerator) drawCrack(texture *ebiten.Image, x, y, length int, angle float64, r, g, b uint8, depth int) {
	bounds := texture.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	prevX, prevY := x, y

	for i := 0; i < length; i++ {
		// Случайное изменение направления
		angle += (tg.rand.Float64() - 0.5) * math.Pi / 4

		// Вычисляем новую позицию
		dx := int(math.Cos(angle))
		dy := int(math.Sin(angle))

		nx, ny := prevX+dx, prevY+dy
		if nx >= 0 && nx < width && ny >= 0 && ny < height {
			texture.Set(nx, ny, color.RGBA{r, g, b, 255})

			// Создаем ответвления с уменьшающейся вероятностью
			if depth < 3 && tg.rand.Float64() < 0.3/float64(depth) {
				branchAngle := angle + math.Pi/2*(tg.rand.Float64()*2-1)
				branchLength := length / 2

				// Рекурсивно рисуем ответвление
				tg.drawCrack(texture, nx, ny, branchLength, branchAngle, r, g, b, depth+1)
			}

			prevX, prevY = nx, ny
		}
	}
}

// addFloorDetails добавляет детали пола
func (tg *TextureGenerator) addFloorDetails(texture *ebiten.Image, count int) {
	bounds := texture.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Рисуем доски или плитку
	if tg.rand.Float64() < 0.5 {
		// Доски
		boardWidth := 4 + tg.rand.Intn(4)

		// Цвет досок
		baseR := uint8(130 + tg.rand.Intn(40))
		baseG := uint8(110 + tg.rand.Intn(30))
		baseB := uint8(80 + tg.rand.Intn(20))

		// Рисуем деревянный пол
		for x := 0; x < width; x += boardWidth {
			// Вариация цвета доски
			variation := tg.rand.Intn(30) - 15
			boardR := clampUint8(int(baseR)+variation, 0, 255)
			boardG := clampUint8(int(baseG)+variation, 0, 255)
			boardB := clampUint8(int(baseB)+variation, 0, 255)

			// Ширина текущей доски
			currentWidth := boardWidth
			if tg.rand.Float64() < 0.3 {
				currentWidth -= 1
			}

			// Рисуем доску
			for j := 0; j < height; j++ {
				for k := 0; k < currentWidth-1; k++ {
					nx := x + k
					if nx < width {
						// Добавляем текстуру дерева
						woodNoise := tg.noise.Perlin2D(float64(nx)*0.2, float64(j)*0.05, 0.5)
						woodVariation := int(woodNoise * 15)

						r := clampUint8(int(boardR)+woodVariation, 0, 255)
						g := clampUint8(int(boardG)+woodVariation, 0, 255)
						b := clampUint8(int(boardB)+woodVariation, 0, 255)

						texture.Set(nx, j, color.RGBA{r, g, b, 255})
					}
				}

				// Рисуем шов между досками
				nx := x + currentWidth - 1
				if nx < width {
					texture.Set(nx, j, color.RGBA{boardR - 30, boardG - 30, boardB - 30, 255})
				}
			}
		}
	} else {
		// Плитка
		tileSize := 8 + tg.rand.Intn(4)

		// Цвет плитки 1
		tile1R := uint8(180 + tg.rand.Intn(40))
		tile1G := uint8(180 + tg.rand.Intn(40))
		tile1B := uint8(180 + tg.rand.Intn(40))

		// Цвет плитки 2
		tile2R := uint8(140 + tg.rand.Intn(40))
		tile2G := uint8(140 + tg.rand.Intn(40))
		tile2B := uint8(140 + tg.rand.Intn(40))

		// Рисуем плиточный пол
		for y := 0; y < height; y += tileSize {
			for x := 0; x < width; x += tileSize {
				// Выбираем цвет плитки (шахматный узор)
				var tileR, tileG, tileB uint8
				if (x/tileSize+y/tileSize)%2 == 0 {
					tileR, tileG, tileB = tile1R, tile1G, tile1B
				} else {
					tileR, tileG, tileB = tile2R, tile2G, tile2B
				}

				// Рисуем плитку
				for j := 0; j < tileSize-1; j++ {
					for k := 0; k < tileSize-1; k++ {
						nx, ny := x+k, y+j
						if nx < width && ny < height {
							// Вариация цвета плитки
							tileNoise := tg.noise.Perlin2D(float64(nx)*0.2, float64(ny)*0.2, 0.5)
							tileVariation := int(tileNoise * 10)

							r := clampUint8(int(tileR)+tileVariation, 0, 255)
							g := clampUint8(int(tileG)+tileVariation, 0, 255)
							b := clampUint8(int(tileB)+tileVariation, 0, 255)

							texture.Set(nx, ny, color.RGBA{r, g, b, 255})
						}
					}
				}

				// Рисуем швы
				for j := 0; j < tileSize; j++ {
					nx, ny := x+tileSize-1, y+j
					if nx < width && ny < height {
						texture.Set(nx, ny, color.RGBA{50, 50, 50, 255})
					}

					nx, ny = x+j, y+tileSize-1
					if nx < width && ny < height {
						texture.Set(nx, ny, color.RGBA{50, 50, 50, 255})
					}
				}
			}
		}
	}
}

// addCreatureDetails добавляет детали существа
func (tg *TextureGenerator) addCreatureDetails(texture *ebiten.Image, count int) {
	bounds := texture.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Добавляем органические детали (вены, глаза, текстуру кожи)
	creatureType := tg.rand.Intn(3)

	switch creatureType {
	case 0: // Кожистое
		// Базовая текстура кожи
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				// Получаем текущий цвет
				r, g, b, a := texture.At(x, y).RGBA()

				// Добавляем текстуру кожи
				skinNoise := tg.noise.Perlin2D(float64(x)*0.2, float64(y)*0.2, 0.5)
				variation := int(skinNoise * 30)

				newR := clampUint8(int(r>>8)+variation, 0, 255)
				newG := clampUint8(int(g>>8)+variation, 0, 255)
				newB := clampUint8(int(b>>8)+variation, 0, 255)

				texture.Set(x, y, color.RGBA{newR, newG, newB, uint8(a >> 8)})
			}
		}

		// Добавляем складки
		for i := 0; i < count/3; i++ {
			x := tg.rand.Intn(width)
			y := tg.rand.Intn(height)
			length := 3 + tg.rand.Intn(7)
			angle := tg.rand.Float64() * math.Pi

			// Цвет складки (темнее)
			r, g, b, _ := texture.At(x, y).RGBA()
			foldR := clampUint8(int(r>>8)-30, 0, 255)
			foldG := clampUint8(int(g>>8)-30, 0, 255)
			foldB := clampUint8(int(b>>8)-30, 0, 255)

			// Рисуем складку
			for j := 0; j < length; j++ {
				// Слегка искривляем линию
				angle += (tg.rand.Float64() - 0.5) * 0.2

				dx := int(math.Cos(angle) * float64(j))
				dy := int(math.Sin(angle) * float64(j))

				nx, ny := x+dx, y+dy
				if nx >= 0 && nx < width && ny >= 0 && ny < height {
					texture.Set(nx, ny, color.RGBA{foldR, foldG, foldB, 255})

					// Добавляем тень рядом со складкой
					for k := -1; k <= 1; k++ {
						for l := -1; l <= 1; l++ {
							if k == 0 && l == 0 {
								continue
							}

							shadedX, shadedY := nx+k, ny+l
							if shadedX >= 0 && shadedX < width && shadedY >= 0 && shadedY < height {
								sr, sg, sb, sa := texture.At(shadedX, shadedY).RGBA()

								// Затемняем цвет для тени
								shadeR := clampUint8(int(sr>>8)-15, 0, 255)
								shadeG := clampUint8(int(sg>>8)-15, 0, 255)
								shadeB := clampUint8(int(sb>>8)-15, 0, 255)

								texture.Set(shadedX, shadedY, color.RGBA{shadeR, shadeG, shadeB, uint8(sa >> 8)})
							}
						}
					}
				}
			}
		}

	case 1: // Слизистое
		// Делаем текстуру более гладкой и блестящей
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				// Получаем текущий цвет
				r, g, b, a := texture.At(x, y).RGBA()

				// Добавляем слизистую текстуру
				slimeNoise := tg.noise.Perlin2D(float64(x)*0.1, float64(y)*0.1, 0.5)
				variation := int(slimeNoise * 40)

				// Делаем более зеленоватым
				newR := clampUint8(int(r>>8)-20+variation/3, 0, 255)
				newG := clampUint8(int(g>>8)+10+variation, 0, 255)
				newB := clampUint8(int(b>>8)-10+variation/2, 0, 255)

				texture.Set(x, y, color.RGBA{newR, newG, newB, uint8(a >> 8)})
			}
		}

		// Добавляем пузыри
		for i := 0; i < count/2; i++ {
			x := tg.rand.Intn(width)
			y := tg.rand.Intn(height)
			size := 1 + tg.rand.Intn(3)

			// Цвет пузыря (светлее и более прозрачный)
			r, g, b, _ := texture.At(x, y).RGBA()
			bubbleR := clampUint8(int(r>>8)+50, 0, 255)
			bubbleG := clampUint8(int(g>>8)+50, 0, 255)
			bubbleB := clampUint8(int(b>>8)+50, 0, 255)

			// Рисуем пузырь
			for j := -size; j <= size; j++ {
				for k := -size; k <= size; k++ {
					dist := j*j + k*k
					if dist <= size*size {
						nx, ny := x+j, y+k
						if nx >= 0 && nx < width && ny >= 0 && ny < height {
							// Центр пузыря светлее
							highlight := 1.0 - float64(dist)/(float64(size*size))

							r := clampUint8(int(bubbleR)+int(highlight*80), 0, 255)
							g := clampUint8(int(bubbleG)+int(highlight*80), 0, 255)
							b := clampUint8(int(bubbleB)+int(highlight*80), 0, 255)

							texture.Set(nx, ny, color.RGBA{r, g, b, 255})
						}
					}
				}
			}
		}

	case 2: // Чешуйчатое/насекомое
		// Добавляем чешуйки или пластины
		scaleSize := 2 + tg.rand.Intn(3)

		for y := 0; y < height; y += scaleSize {
			for x := 0; x < width; x += scaleSize {
				// Получаем базовый цвет области
				baseX, baseY := clamp(x, 0, width-1), clamp(y, 0, height-1)
				r, g, b, _ := texture.At(baseX, baseY).RGBA()

				// Вариация цвета
				variation := tg.rand.Intn(30) - 15

				scaleR := clampUint8(int(r>>8)+variation, 0, 255)
				scaleG := clampUint8(int(g>>8)+variation, 0, 255)
				scaleB := clampUint8(int(b>>8)+variation, 0, 255)

				// Рисуем чешуйку/пластину
				for j := 0; j < scaleSize-1; j++ {
					for k := 0; k < scaleSize-1; k++ {
						nx, ny := x+k, y+j
						if nx < width && ny < height {
							texture.Set(nx, ny, color.RGBA{scaleR, scaleG, scaleB, 255})
						}
					}
				}

				// Рисуем границу чешуйки
				for j := 0; j < scaleSize; j++ {
					nx, ny := x+scaleSize-1, y+j
					if nx < width && ny < height {
						texture.Set(nx, ny, color.RGBA{
							clampUint8(int(scaleR)-40, 0, 255),
							clampUint8(int(scaleG)-40, 0, 255),
							clampUint8(int(scaleB)-40, 0, 255),
							255,
						})
					}

					nx, ny = x+j, y+scaleSize-1
					if nx < width && ny < height {
						texture.Set(nx, ny, color.RGBA{
							clampUint8(int(scaleR)-40, 0, 255),
							clampUint8(int(scaleG)-40, 0, 255),
							clampUint8(int(scaleB)-40, 0, 255),
							255,
						})
					}
				}
			}
		}
	}

	// Добавляем глаза или другие органические особенности
	if tg.rand.Float64() < 0.7 {
		eyeCount := 1 + tg.rand.Intn(5)

		for i := 0; i < eyeCount; i++ {
			x := tg.rand.Intn(width)
			y := tg.rand.Intn(height)
			size := 1 + tg.rand.Intn(3)

			// Цвет глаза (случайный)
			eyeColors := []color.RGBA{
				{255, 0, 0, 255},     // Красный
				{255, 255, 0, 255},   // Желтый
				{0, 255, 0, 255},     // Зеленый
				{0, 0, 255, 255},     // Синий
				{255, 0, 255, 255},   // Фиолетовый
				{0, 0, 0, 255},       // Черный
				{255, 255, 255, 255}, // Белый
			}

			eyeColor := eyeColors[tg.rand.Intn(len(eyeColors))]

			// Рисуем глаз
			for j := -size; j <= size; j++ {
				for k := -size; k <= size; k++ {
					dist := j*j + k*k
					if dist <= size*size {
						nx, ny := x+j, y+k
						if nx >= 0 && nx < width && ny >= 0 && ny < height {
							if dist <= (size/2)*(size/2) {
								// Зрачок (черный)
								texture.Set(nx, ny, color.RGBA{0, 0, 0, 255})
							} else {
								// Радужка
								texture.Set(nx, ny, eyeColor)
							}
						}
					}
				}
			}

			// Добавляем блик в глазу
			blinkX, blinkY := x-size/2, y-size/2
			if blinkX >= 0 && blinkX < width && blinkY >= 0 && blinkY < height {
				texture.Set(blinkX, blinkY, color.RGBA{255, 255, 255, 255})
			}
		}
	}
}

// addEffectDetails добавляет детали эффекта
func (tg *TextureGenerator) addEffectDetails(texture *ebiten.Image, count int) {
	bounds := texture.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Выбираем тип эффекта
	effectType := tg.rand.Intn(3)

	switch effectType {
	case 0: // Дым/туман
		// Добавляем клубы дыма/тумана
		for i := 0; i < count; i++ {
			x := tg.rand.Intn(width)
			y := tg.rand.Intn(height)
			size := 3 + tg.rand.Intn(5)

			// Цвет дыма (сероватый с вариациями)
			smokeR := uint8(200 + tg.rand.Intn(55))
			smokeG := uint8(200 + tg.rand.Intn(55))
			smokeB := uint8(200 + tg.rand.Intn(55))

			// Рисуем клуб дыма
			for j := -size; j <= size; j++ {
				for k := -size; k <= size; k++ {
					// Деформируем круг для более органичной формы
					deform := 1.0 + 0.3*tg.noise.Perlin2D(float64(j+x)*0.2, float64(k+y)*0.2, 0.5)
					if (j*j + k*k) <= int(float64(size*size)*deform) {
						nx, ny := x+j, y+k
						if nx >= 0 && nx < width && ny >= 0 && ny < height {
							// Прозрачность зависит от расстояния до центра
							dist := math.Sqrt(float64(j*j + k*k))
							alpha := uint8(float64(255) * (float64(1.0) - float64(dist)/float64(float64(size)*float64(deform))))

							// Добавляем небольшую вариацию цвета
							variation := tg.rand.Intn(20) - 10
							r := clampUint8(int(smokeR)+variation, 0, 255)
							g := clampUint8(int(smokeG)+variation, 0, 255)
							b := clampUint8(int(smokeB)+variation, 0, 255)

							texture.Set(nx, ny, color.RGBA{r, g, b, alpha})
						}
					}
				}
			}
		}

	case 1: // Энергия/магия
		// Базовый цвет энергии
		baseR := uint8(tg.rand.Intn(100))
		baseG := uint8(100 + tg.rand.Intn(155))
		baseB := uint8(200 + tg.rand.Intn(55))

		// Добавляем энергетические линии
		for i := 0; i < count/2; i++ {
			x := tg.rand.Intn(width)
			y := tg.rand.Intn(height)
			length := 5 + tg.rand.Intn(10)

			// Рисуем извилистую линию
			prevX, prevY := x, y
			angle := tg.rand.Float64() * 2 * math.Pi

			for j := 0; j < length; j++ {
				// Слегка меняем направление
				angle += (tg.rand.Float64() - 0.5) * math.Pi / 2

				// Вычисляем новую позицию
				dx := int(math.Cos(angle) * 2)
				dy := int(math.Sin(angle) * 2)

				nx, ny := prevX+dx, prevY+dy
				if nx >= 0 && nx < width && ny >= 0 && ny < height {
					// Вариация цвета вдоль линии
					progress := float64(j) / float64(length)
					r := clampUint8(int(baseR)+int(progress*100), 0, 255)
					g := clampUint8(int(baseG)-int(progress*50), 0, 255)
					b := clampUint8(int(baseB), 0, 255)

					// Рисуем линию с ореолом
					for k := -1; k <= 1; k++ {
						for l := -1; l <= 1; l++ {
							px, py := nx+k, ny+l
							if px >= 0 && px < width && py >= 0 && py < height {
								dist := k*k + l*l

								if dist == 0 {
									// Центр линии
									texture.Set(px, py, color.RGBA{r, g, b, 255})
								} else {
									// Ореол
									texture.Set(px, py, color.RGBA{r / 2, g / 2, b / 2, 128})
								}
							}
						}
					}

					prevX, prevY = nx, ny
				}
			}
		}

		// Добавляем энергетические искры
		for i := 0; i < count/2; i++ {
			x := tg.rand.Intn(width)
			y := tg.rand.Intn(height)
			size := 1 + tg.rand.Intn(2)

			// Рисуем искру
			for j := -size; j <= size; j++ {
				for k := -size; k <= size; k++ {
					if j*j+k*k <= size*size {
						nx, ny := x+j, y+k
						if nx >= 0 && nx < width && ny >= 0 && ny < height {
							// Яркий центр
							r := uint8(200 + tg.rand.Intn(55))
							g := uint8(200 + tg.rand.Intn(55))
							b := uint8(255)

							texture.Set(nx, ny, color.RGBA{r, g, b, 255})
						}
					}
				}
			}
		}

	case 2: // Огонь/пламя
		// Базовые цвета огня
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				// Градиент от низа к верху
				gradientY := 1.0 - float64(y)/float64(height)

				// Шум для придания формы пламени
				flameNoise := tg.noise.Perlin2D(float64(x)*0.1, float64(y)*0.1, 0.5)

				// Цвет зависит от высоты и шума
				var r, g, b uint8

				if gradientY < 0.3 {
					// Нижняя часть пламени (оранжевая/красная)
					intensity := gradientY / 0.3 * (0.7 + 0.3*flameNoise)
					r = uint8(255 * intensity)
					g = uint8(100 * intensity)
					b = uint8(0)
				} else if gradientY < 0.7 {
					// Средняя часть пламени (оранжевая/желтая)
					normalizedY := (gradientY - 0.3) / 0.4
					intensity := 0.7 + 0.3*flameNoise
					r = 255
					g = uint8(100 + 155*normalizedY*intensity)
					b = uint8(normalizedY * 50 * intensity)
				} else {
					// Верхняя часть пламени (желтая/белая)
					normalizedY := (gradientY - 0.7) / 0.3
					intensity := 0.7 + 0.3*flameNoise
					r = 255
					g = uint8(200 + 55*normalizedY*intensity)
					b = uint8(50 + 205*normalizedY*intensity)
				}

				// Альфа-канал зависит от шума и высоты
				alpha := uint8(255 * math.Max(0, math.Min(1, gradientY*(0.5+0.5*flameNoise))))

				texture.Set(x, y, color.RGBA{r, g, b, alpha})
			}
		}
	}
}

// modifyTexture модифицирует текстуру на основе seed
func (tg *TextureGenerator) modifyTexture(texture *ebiten.Image, textureType TextureType, seed int64) {
	// Используем переданный seed для детерминированной модификации
	r := rand.New(rand.NewSource(seed))

	bounds := texture.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Модифицируем текстуру в зависимости от типа
	switch textureType {
	case TextureGrass, TextureForest, TextureDenseForest:
		// Меняем оттенок зеленого
		hueShift := r.Float64()*0.2 - 0.1
		saturationShift := r.Float64()*0.3 - 0.1
		brightnessShift := r.Float64()*0.2 - 0.1

		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				r, g, b, a := texture.At(x, y).RGBA()
				h, s, v := rgbToHsv(float64(r>>8)/255, float64(g>>8)/255, float64(b>>8)/255)

				// Модифицируем HSV
				h = math.Mod(h+hueShift, 1.0)
				s = clampFloat64(s+saturationShift, 0, 1)
				v = clampFloat64(v+brightnessShift, 0, 1)

				newR, newG, newB := hsvToRgb(h, s, v)

				texture.Set(x, y, color.RGBA{
					uint8(newR * 255),
					uint8(newG * 255),
					uint8(newB * 255),
					uint8(a >> 8),
				})
			}
		}

	case TexturePath, TextureRocks, TextureWall, TextureFloor:
		// Меняем яркость и контраст
		brightnessShift := r.Float64()*0.3 - 0.15
		contrastFactor := 0.8 + r.Float64()*0.4

		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				r, g, b, a := texture.At(x, y).RGBA()

				// Нормализуем к диапазону [0, 1]
				rf := float64(r>>8) / 255
				gf := float64(g>>8) / 255
				bf := float64(b>>8) / 255

				// Применяем яркость
				rf += brightnessShift
				gf += brightnessShift
				bf += brightnessShift

				// Применяем контраст
				rf = (rf-0.5)*contrastFactor + 0.5
				gf = (gf-0.5)*contrastFactor + 0.5
				bf = (bf-0.5)*contrastFactor + 0.5

				// Ограничиваем к диапазону [0, 1]
				rf = clampFloat64(rf, 0, 1)
				gf = clampFloat64(gf, 0, 1)
				bf = clampFloat64(bf, 0, 1)

				texture.Set(x, y, color.RGBA{
					uint8(rf * 255),
					uint8(gf * 255),
					uint8(bf * 255),
					uint8(a >> 8),
				})
			}
		}

	case TextureWater, TextureSwamp:
		// Меняем оттенок синего/зеленого
		hueShift := r.Float64()*0.2 - 0.1
		saturationShift := r.Float64() * 0.3

		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				r, g, b, a := texture.At(x, y).RGBA()
				h, s, v := rgbToHsv(float64(r>>8)/255, float64(g>>8)/255, float64(b>>8)/255)

				// Модифицируем HSV
				h = math.Mod(h+hueShift, 1.0)
				s = clampFloat64(s+saturationShift, 0, 1)

				newR, newG, newB := hsvToRgb(h, s, v)

				texture.Set(x, y, color.RGBA{
					uint8(newR * 255),
					uint8(newG * 255),
					uint8(newB * 255),
					uint8(a >> 8),
				})
			}
		}

	case TextureCorrupted:
		// Усиливаем эффект искажения
		intensityFactor := 0.7 + r.Float64()*0.6

		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				r, g, b, a := texture.At(x, y).RGBA()

				// Усиливаем красный и уменьшаем синий
				newR := clampUint8(int(r>>8)+int(float64(r>>8)*0.2*intensityFactor), 0, 255)
				newG := clampUint8(int(g>>8)-int(float64(g>>8)*0.1*intensityFactor), 0, 255)
				newB := clampUint8(int(b>>8)+int(float64(b>>8)*0.1*intensityFactor), 0, 255)

				texture.Set(x, y, color.RGBA{newR, newG, newB, uint8(a >> 8)})
			}
		}

	case TextureCreature:
		// Случайно выбираем цветовую схему
		colorScheme := r.Intn(5)

		var targetR, targetG, targetB float64
		switch colorScheme {
		case 0: // Красноватый
			targetR, targetG, targetB = 0.8, 0.2, 0.2
		case 1: // Зеленоватый
			targetR, targetG, targetB = 0.2, 0.8, 0.2
		case 2: // Синеватый
			targetR, targetG, targetB = 0.2, 0.2, 0.8
		case 3: // Желтоватый
			targetR, targetG, targetB = 0.8, 0.8, 0.2
		case 4: // Фиолетовый
			targetR, targetG, targetB = 0.8, 0.2, 0.8
		}

		// Интенсивность смещения к целевому цвету
		intensity := 0.3 + r.Float64()*0.4

		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				r, g, b, a := texture.At(x, y).RGBA()

				// Нормализуем к диапазону [0, 1]
				rf := float64(r>>8) / 255
				gf := float64(g>>8) / 255
				bf := float64(b>>8) / 255

				// Смещаем к целевому цвету
				rf = rf*(1-intensity) + targetR*intensity
				gf = gf*(1-intensity) + targetG*intensity
				bf = bf*(1-intensity) + targetB*intensity

				texture.Set(x, y, color.RGBA{
					uint8(rf * 255),
					uint8(gf * 255),
					uint8(bf * 255),
					uint8(a >> 8),
				})
			}
		}

	case TextureEffect:
		// Меняем цвет и прозрачность
		hueShift := r.Float64()
		alphaFactor := 0.7 + r.Float64()*0.6

		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				r, g, b, a := texture.At(x, y).RGBA()
				h, s, v := rgbToHsv(float64(r>>8)/255, float64(g>>8)/255, float64(b>>8)/255)

				// Модифицируем HSV
				h = math.Mod(h+hueShift, 1.0)

				newR, newG, newB := hsvToRgb(h, s, v)
				newAlpha := clampUint8(int(float64(a>>8)*alphaFactor), 0, 255)

				texture.Set(x, y, color.RGBA{
					uint8(newR * 255),
					uint8(newG * 255),
					uint8(newB * 255),
					newAlpha,
				})
			}
		}
	}
}

// applyCorruption применяет эффект искажения к текстуре
func (tg *TextureGenerator) applyCorruption(texture *ebiten.Image, corruptionLevel float64) {
	bounds := texture.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Ограничиваем уровень коррупции
	corruptionLevel = clampFloat64(corruptionLevel, 0, 1)

	// Получаем базовую текстуру искажения
	corruptedTexture := tg.GetBaseTexture(TextureCorrupted)

	// Создаем шум искажения
	distortionNoise := make([][]float64, height)
	for y := range distortionNoise {
		distortionNoise[y] = make([]float64, width)
		for x := range distortionNoise[y] {
			// Генерируем шум для искажения
			distortionNoise[y][x] = tg.noise.Perlin2D(float64(x)*0.05, float64(y)*0.05, 0.5)
		}
	}

	// Применяем искажение
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Локальный уровень искажения
			localCorruption := corruptionLevel * (0.5 + 0.5*distortionNoise[y][x])

			if localCorruption > 0 {
				// Получаем цвета из обеих текстур
				origR, origG, origB, origA := texture.At(x, y).RGBA()
				corrR, corrG, corrB, corrA := corruptedTexture.At(x, y).RGBA()

				// Смешиваем цвета
				r := uint8((float64(origR>>8)*(1-localCorruption) + float64(corrR>>8)*localCorruption))
				g := uint8((float64(origG>>8)*(1-localCorruption) + float64(corrG>>8)*localCorruption))
				b := uint8((float64(origB>>8)*(1-localCorruption) + float64(corrB>>8)*localCorruption))
				a := uint8((float64(origA>>8)*(1-localCorruption) + float64(corrA>>8)*localCorruption))

				// Добавляем пурпурный оттенок
				if localCorruption > 0.5 {
					purpleFactor := (localCorruption - 0.5) * 2
					r = clampUint8(int(r)+int(purpleFactor*100), 0, 255)
					b = clampUint8(int(b)+int(purpleFactor*100), 0, 255)
				}

				// Применяем искажение координат
				if localCorruption > 0.3 {
					// Сильное искажение для высокого уровня коррупции
					distortionFactor := (localCorruption - 0.3) * 10

					// Искажаем координаты
					distX := float64(x) + (distortionNoise[y][x]-0.5)*distortionFactor
					distY := float64(y) + (distortionNoise[y][x]-0.5)*distortionFactor

					// Убеждаемся, что координаты в пределах текстуры
					distX = math.Max(0, math.Min(float64(width-1), distX))
					distY = math.Max(0, math.Min(float64(height-1), distY))

					// Получаем цвет из искаженных координат
					distR, distG, distB, distA := texture.At(int(distX), int(distY)).RGBA()

					// Смешиваем с искаженным цветом
					distortionStrength := math.Min(1, (localCorruption-0.3)*2)
					r = uint8((float64(r)*(1-distortionStrength) + float64(distR>>8)*distortionStrength))
					g = uint8((float64(g)*(1-distortionStrength) + float64(distG>>8)*distortionStrength))
					b = uint8((float64(b)*(1-distortionStrength) + float64(distB>>8)*distortionStrength))
					a = uint8((float64(a)*(1-distortionStrength) + float64(distA>>8)*distortionStrength))
				}

				// Устанавливаем итоговый цвет
				texture.Set(x, y, color.RGBA{r, g, b, a})
			}
		}
	}

	// Добавляем "вены" искажения для сильной коррупции
	if corruptionLevel > 0.7 {
		veinCount := int(corruptionLevel * 10)

		for i := 0; i < veinCount; i++ {
			x := tg.rand.Intn(width)
			y := tg.rand.Intn(height)
			length := 5 + tg.rand.Intn(int(corruptionLevel*20))

			// Цвет "вены"
			veinR := uint8(150 + tg.rand.Intn(105))
			veinG := uint8(tg.rand.Intn(50))
			veinB := uint8(150 + tg.rand.Intn(105))

			// Рисуем "вену"
			tg.drawCorruptionVein(texture, x, y, length, veinR, veinG, veinB)
		}
	}
}

// drawCorruptionVein рисует "вену" искажения
func (tg *TextureGenerator) drawCorruptionVein(texture *ebiten.Image, x, y, length int, r, g, b uint8) {
	bounds := texture.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Начальное направление
	angle := tg.rand.Float64() * 2 * math.Pi

	prevX, prevY := x, y

	for i := 0; i < length; i++ {
		// Случайное изменение направления
		angle += (tg.rand.Float64() - 0.5) * math.Pi / 2

		// Вычисляем новую позицию
		dx := int(math.Cos(angle))
		dy := int(math.Sin(angle))

		nx, ny := prevX+dx, prevY+dy
		if nx >= 0 && nx < width && ny >= 0 && ny < height {
			// Рисуем основную "вену"
			texture.Set(nx, ny, color.RGBA{r, g, b, 255})

			// Рисуем "ореол" вокруг вены
			for j := -1; j <= 1; j++ {
				for k := -1; k <= 1; k++ {
					if j == 0 && k == 0 {
						continue
					}

					haloX, haloY := nx+j, ny+k
					if haloX >= 0 && haloX < width && haloY >= 0 && haloY < height {
						// Получаем текущий цвет
						origR, origG, origB, origA := texture.At(haloX, haloY).RGBA()

						// Смешиваем с цветом вены
						blendFactor := 0.3
						newR := uint8(float64(origR>>8)*(1-blendFactor) + float64(r)*blendFactor)
						newG := uint8(float64(origG>>8)*(1-blendFactor) + float64(g)*blendFactor)
						newB := uint8(float64(origB>>8)*(1-blendFactor) + float64(b)*blendFactor)

						texture.Set(haloX, haloY, color.RGBA{newR, newG, newB, uint8(origA >> 8)})
					}
				}
			}

			// Иногда добавляем ответвления
			if tg.rand.Float64() < 0.2 {
				branchAngle := angle + math.Pi/2*(tg.rand.Float64()*2-1)
				branchLength := 1 + tg.rand.Intn(3)

				bprevX, bprevY := nx, ny

				for j := 0; j < branchLength; j++ {
					// Слегка меняем направление ответвления
					branchAngle += (tg.rand.Float64() - 0.5) * math.Pi / 4

					// Вычисляем новую позицию ответвления
					bdx := int(math.Cos(branchAngle))
					bdy := int(math.Sin(branchAngle))

					bnx, bny := bprevX+bdx, bprevY+bdy
					if bnx >= 0 && bnx < width && bny >= 0 && bny < height {
						// Рисуем ответвление
						texture.Set(bnx, bny, color.RGBA{r, g, b, 255})

						bprevX, bprevY = bnx, bny
					}
				}
			}

			prevX, prevY = nx, ny
		}
	}
}

// rgbToHsv преобразует RGB в HSV
func rgbToHsv(r, g, b float64) (float64, float64, float64) {
	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))

	h, s, v := 0.0, 0.0, max

	// Вычисляем насыщенность
	if max > 0 {
		s = (max - min) / max
	} else {
		return 0, 0, 0
	}

	// Вычисляем оттенок
	if max == min {
		h = 0 // Ахроматический (серый)
	} else {
		d := max - min
		if max == r {
			h = (g - b) / d
			if g < b {
				h += 6
			}
		} else if max == g {
			h = (b-r)/d + 2
		} else {
			h = (r-g)/d + 4
		}
		h /= 6
	}

	return h, s, v
}

// hsvToRgb преобразует HSV в RGB
func hsvToRgb(h, s, v float64) (float64, float64, float64) {
	if s == 0 {
		return v, v, v
	}

	h *= 6
	i := math.Floor(h)
	f := h - i
	p := v * (1 - s)
	q := v * (1 - s*f)
	t := v * (1 - s*(1-f))

	switch int(i) % 6 {
	case 0:
		return v, t, p
	case 1:
		return q, v, p
	case 2:
		return p, v, t
	case 3:
		return p, q, v
	case 4:
		return t, p, v
	default:
		return v, p, q
	}
}

// clampUint8 ограничивает значение для uint8
func clampUint8(value, min, max int) uint8 {
	if value < min {
		value = min
	}
	if value > max {
		value = max
	}
	return uint8(value)
}

// clamp ограничивает значение для int
func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// GenerateTextureAtlas создает атлас текстур
func (tg *TextureGenerator) GenerateTextureAtlas() *ebiten.Image {
	// Определяем количество текстур каждого типа
	typeCounts := map[TextureType]int{
		TextureGrass:       4,
		TextureForest:      4,
		TextureDenseForest: 4,
		TexturePath:        4,
		TextureRocks:       4,
		TextureWater:       4,
		TextureSwamp:       4,
		TextureCorrupted:   4,
		TextureWall:        4,
		TextureFloor:       4,
		TextureCreature:    8,
		TextureEffect:      8,
	}

	// Подсчитываем общее количество текстур
	totalTextures := 0
	for _, count := range typeCounts {
		totalTextures += count
	}

	// Определяем размеры атласа
	atlasWidth := 4
	atlasHeight := (totalTextures + atlasWidth - 1) / atlasWidth

	// Создаем атлас
	atlas := ebiten.NewImage(atlasWidth*tg.tileSize, atlasHeight*tg.tileSize)

	// Генерируем текстуры и размещаем их в атласе
	index := 0
	for textureType := TextureGrass; textureType <= TextureEffect; textureType++ {
		count := typeCounts[textureType]

		for i := 0; i < count; i++ {
			// Генерируем уникальную текстуру
			texture := tg.GenerateUniqueTexture(textureType, int64(i))

			// Определяем позицию в атласе
			x := (index % atlasWidth) * tg.tileSize
			y := (index / atlasWidth) * tg.tileSize

			// Размещаем текстуру в атласе
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x), float64(y))
			atlas.DrawImage(texture, op)

			index++
		}
	}

	return atlas
}

// GetTextureFromAtlas получает текстуру из атласа
func (tg *TextureGenerator) GetTextureFromAtlas(atlas *ebiten.Image, index int) *ebiten.Image {
	atlasWidth := atlas.Bounds().Dx() / tg.tileSize

	// Определяем позицию в атласе
	x := (index % atlasWidth) * tg.tileSize
	y := (index / atlasWidth) * tg.tileSize

	// Создаем новую текстуру
	texture := ebiten.NewImage(tg.tileSize, tg.tileSize)

	// Копируем часть атласа
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(-x), float64(-y))
	texture.DrawImage(atlas.SubImage(image.Rect(x, y, x+tg.tileSize, y+tg.tileSize)).(*ebiten.Image), op)

	return texture
}

// GenerateNormalMap создает карту нормалей для текстуры
func (tg *TextureGenerator) GenerateNormalMap(texture *ebiten.Image) *ebiten.Image {
	bounds := texture.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Создаем карту нормалей
	normalMap := ebiten.NewImage(width, height)

	// Вычисляем нормали
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Получаем высоту в текущей точке
			_, _, v, _ := texture.At(x, y).RGBA()
			height := float64(v) / 65535.0

			// Получаем высоты в соседних точках
			x1, y1 := clamp(x-1, 0, width-1), y
			x2, y2 := clamp(x+1, 0, width-1), y
			x3, y3 := x, clamp(y-1, 0, int(height)-1)
			x4, y4 := x, clamp(y+1, 0, int(height)-1)

			_, _, v1, _ := texture.At(x1, y1).RGBA()
			_, _, v2, _ := texture.At(x2, y2).RGBA()
			_, _, v3, _ := texture.At(x3, y3).RGBA()
			_, _, v4, _ := texture.At(x4, y4).RGBA()

			height1 := float64(v1) / 65535.0
			height2 := float64(v2) / 65535.0
			height3 := float64(v3) / 65535.0
			height4 := float64(v4) / 65535.0

			// Вычисляем градиент
			dx := (height2 - height1) * 2.0
			dy := (height4 - height3) * 2.0

			// Нормализуем вектор
			length := math.Sqrt(dx*dx + dy*dy + 1.0)
			nx := -dx / length
			ny := -dy / length
			nz := 1.0 / length

			// Преобразуем в диапазон [0, 1]
			nx = nx*0.5 + 0.5
			ny = ny*0.5 + 0.5
			nz = nz*0.5 + 0.5

			// Устанавливаем цвет
			normalMap.Set(x, y, color.RGBA{
				uint8(nx * 255),
				uint8(ny * 255),
				uint8(nz * 255),
				255,
			})
		}
	}

	return normalMap
}

// GenerateHeightMap создает карту высот для текстуры
func (tg *TextureGenerator) GenerateHeightMap(textureType TextureType, seed int64) *ebiten.Image {
	// Создаем новую текстуру
	heightMap := ebiten.NewImage(tg.tileSize, tg.tileSize)

	// Используем переданный seed для детерминированной генерации
	r := rand.New(rand.NewSource(seed))

	// Определяем параметры шума в зависимости от типа текстуры
	scale := 0.1
	roughness := 0.0
	baseHeight := 0.5

	switch textureType {
	case TextureGrass:
		scale = 0.1
		roughness = 0.1
		baseHeight = 0.4
	case TextureForest:
		scale = 0.15
		roughness = 0.2
		baseHeight = 0.5
	case TextureDenseForest:
		scale = 0.2
		roughness = 0.3
		baseHeight = 0.6
	case TexturePath:
		scale = 0.05
		roughness = 0.05
		baseHeight = 0.3
	case TextureRocks:
		scale = 0.3
		roughness = 0.5
		baseHeight = 0.7
	case TextureWater:
		scale = 0.05
		roughness = 0.05
		baseHeight = 0.2
	case TextureSwamp:
		scale = 0.1
		roughness = 0.2
		baseHeight = 0.3
	case TextureCorrupted:
		scale = 0.2
		roughness = 0.4
		baseHeight = 0.6
	case TextureWall:
		scale = 0.1
		roughness = 0.3
		baseHeight = 0.8
	case TextureFloor:
		scale = 0.05
		roughness = 0.1
		baseHeight = 0.4
	default:
		scale = 0.1
		roughness = 0.2
		baseHeight = 0.5
	}

	// Небольшая случайная вариация параметров
	scale += r.Float64()*0.05 - 0.025
	roughness += r.Float64()*0.1 - 0.05
	baseHeight += r.Float64()*0.1 - 0.05

	// Генерируем карту высот
	for y := 0; y < tg.tileSize; y++ {
		for x := 0; x < tg.tileSize; x++ {
			// Нормализованные координаты
			nx := float64(x) / float64(tg.tileSize)
			ny := float64(y) / float64(tg.tileSize)

			// Генерируем шум
			noise := tg.noise.Perlin2D(nx*scale*10, ny*scale*10, 0.5) * roughness

			// Добавляем базовую высоту
			height := baseHeight + noise

			// Ограничиваем значение
			height = clampFloat64(height, 0, 1)

			// Устанавливаем цвет (оттенки серого)
			value := uint8(height * 255)
			heightMap.Set(x, y, color.RGBA{value, value, value, 255})
		}
	}

	// Добавляем детали в зависимости от типа текстуры
	switch textureType {
	case TextureRocks:
		// Добавляем выступы и впадины
		for i := 0; i < 20; i++ {
			x := r.Intn(tg.tileSize)
			y := r.Intn(tg.tileSize)
			size := 1 + r.Intn(3)

			// Выступ или впадина
			value := uint8(0)
			if r.Float64() < 0.5 {
				value = 255 // Выступ
			}

			// Рисуем деталь
			for j := -size; j <= size; j++ {
				for k := -size; k <= size; k++ {
					if j*j+k*k <= size*size {
						nx, ny := x+j, y+k
						if nx >= 0 && nx < tg.tileSize && ny >= 0 && ny < tg.tileSize {
							heightMap.Set(nx, ny, color.RGBA{value, value, value, 255})
						}
					}
				}
			}
		}

	case TextureWall:
		// Добавляем кирпичи
		brickWidth := 4 + r.Intn(3)
		brickHeight := 2 + r.Intn(2)

		for y := 0; y < tg.tileSize; y += brickHeight {
			// Смещение для чередования кирпичей
			offset := 0
			if (y/brickHeight)%2 == 1 {
				offset = brickWidth / 2
			}

			for x := -offset; x < tg.tileSize; x += brickWidth {
				// Рисуем кирпич
				for j := 0; j < brickHeight-1; j++ {
					for k := 0; k < brickWidth-1; k++ {
						nx, ny := x+k, y+j
						if nx >= 0 && nx < tg.tileSize && ny >= 0 && ny < tg.tileSize {
							// Варьируем высоту кирпича
							variation := r.Intn(20) - 10
							value := clampUint8(200+variation, 0, 255)
							heightMap.Set(nx, ny, color.RGBA{value, value, value, 255})
						}
					}
				}

				// Рисуем швы
				for j := brickHeight - 1; j < brickHeight; j++ {
					for k := 0; k < brickWidth; k++ {
						nx, ny := x+k, y+j
						if nx >= 0 && nx < tg.tileSize && ny >= 0 && ny < tg.tileSize {
							heightMap.Set(nx, ny, color.RGBA{100, 100, 100, 255})
						}
					}
				}

				for j := 0; j < brickHeight; j++ {
					for k := brickWidth - 1; k < brickWidth; k++ {
						nx, ny := x+k, y+j
						if nx >= 0 && nx < tg.tileSize && ny >= 0 && ny < tg.tileSize {
							heightMap.Set(nx, ny, color.RGBA{100, 100, 100, 255})
						}
					}
				}
			}
		}

	case TextureCorrupted:
		// Добавляем "вены"
		for i := 0; i < 10; i++ {
			x := r.Intn(tg.tileSize)
			y := r.Intn(tg.tileSize)
			length := 5 + r.Intn(10)

			// Рисуем вену
			angle := r.Float64() * 2 * math.Pi
			prevX, prevY := x, y

			for j := 0; j < length; j++ {
				// Случайное изменение направления
				angle += (r.Float64() - 0.5) * math.Pi / 2

				// Вычисляем новую позицию
				dx := int(math.Cos(angle))
				dy := int(math.Sin(angle))

				nx, ny := prevX+dx, prevY+dy
				if nx >= 0 && nx < tg.tileSize && ny >= 0 && ny < tg.tileSize {
					heightMap.Set(nx, ny, color.RGBA{255, 255, 255, 255})

					// Добавляем ореол
					for k := -1; k <= 1; k++ {
						for l := -1; l <= 1; l++ {
							if k == 0 && l == 0 {
								continue
							}

							hx, hy := nx+k, ny+l
							if hx >= 0 && hx < tg.tileSize && hy >= 0 && hy < tg.tileSize {
								r, _, _, _ := heightMap.At(hx, hy).RGBA()
								value := clampUint8(int(r>>8)+50, 0, 255)
								heightMap.Set(hx, hy, color.RGBA{value, value, value, 255})
							}
						}
					}

					prevX, prevY = nx, ny
				}
			}
		}
	}

	return heightMap
}

// BlendTextures смешивает две текстуры с указанным коэффициентом
func (tg *TextureGenerator) BlendTextures(texture1, texture2 *ebiten.Image, blendFactor float64) *ebiten.Image {
	bounds1 := texture1.Bounds()
	bounds2 := texture2.Bounds()

	// Проверяем, что текстуры одинакового размера
	if bounds1.Dx() != bounds2.Dx() || bounds1.Dy() != bounds2.Dy() {
		return texture1
	}

	width, height := bounds1.Dx(), bounds1.Dy()

	// Создаем новую текстуру
	result := ebiten.NewImage(width, height)

	// Смешиваем текстуры
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r1, g1, b1, a1 := texture1.At(x, y).RGBA()
			r2, g2, b2, a2 := texture2.At(x, y).RGBA()

			r := uint8((float64(r1>>8)*(1-blendFactor) + float64(r2>>8)*blendFactor))
			g := uint8((float64(g1>>8)*(1-blendFactor) + float64(g2>>8)*blendFactor))
			b := uint8((float64(b1>>8)*(1-blendFactor) + float64(b2>>8)*blendFactor))
			a := uint8((float64(a1>>8)*(1-blendFactor) + float64(a2>>8)*blendFactor))

			result.Set(x, y, color.RGBA{r, g, b, a})
		}
	}

	return result
}

// CreateTransitionTexture создает текстуру перехода между двумя типами местности
func (tg *TextureGenerator) CreateTransitionTexture(type1, type2 TextureType, seed int64) *ebiten.Image {
	// Получаем базовые текстуры
	texture1 := tg.GetBaseTexture(type1)
	texture2 := tg.GetBaseTexture(type2)

	// Создаем новую текстуру
	result := ebiten.NewImage(tg.tileSize, tg.tileSize)

	// Используем переданный seed для детерминированной генерации
	r := rand.New(rand.NewSource(seed))

	// Генерируем шум для карты смешивания
	noiseScale := 0.1 + r.Float64()*0.1
	blendMap := make([][]float64, tg.tileSize)
	for y := range blendMap {
		blendMap[y] = make([]float64, tg.tileSize)
		for x := range blendMap[y] {
			// Нормализованные координаты
			nx := float64(x) / float64(tg.tileSize)
			ny := float64(y) / float64(tg.tileSize)

			// Генерируем шум
			noise := tg.noise.Perlin2D(nx*noiseScale*10, ny*noiseScale*10, 0.5)

			// Преобразуем в диапазон [0, 1]
			blendMap[y][x] = noise*0.5 + 0.5
		}
	}

	// Смешиваем текстуры
	for y := 0; y < tg.tileSize; y++ {
		for x := 0; x < tg.tileSize; x++ {
			r1, g1, b1, a1 := texture1.At(x, y).RGBA()
			r2, g2, b2, a2 := texture2.At(x, y).RGBA()

			// Коэффициент смешивания
			blend := blendMap[y][x]

			r := uint8((float64(r1>>8)*(1-blend) + float64(r2>>8)*blend))
			g := uint8((float64(g1>>8)*(1-blend) + float64(g2>>8)*blend))
			b := uint8((float64(b1>>8)*(1-blend) + float64(b2>>8)*blend))
			a := uint8((float64(a1>>8)*(1-blend) + float64(a2>>8)*blend))

			result.Set(x, y, color.RGBA{r, g, b, a})
		}
	}

	return result
}

// GenerateRandomSplatter создает случайное брызги на текстуре
func (tg *TextureGenerator) GenerateRandomSplatter(baseTexture *ebiten.Image, splatterColor color.RGBA, density, size float64, seed int64) *ebiten.Image {
	bounds := baseTexture.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Создаем новую текстуру
	result := ebiten.NewImage(width, height)
	result.DrawImage(baseTexture, nil)

	// Используем переданный seed для детерминированной генерации
	r := rand.New(rand.NewSource(seed))

	// Количество брызг
	splatterCount := int(float64(width*height) * density / 100)

	// Добавляем брызги
	for i := 0; i < splatterCount; i++ {
		x := r.Intn(width)
		y := r.Intn(height)
		splatterSize := int(float64(1+r.Intn(3)) * size)

		// Рисуем брызги
		for j := -splatterSize; j <= splatterSize; j++ {
			for k := -splatterSize; k <= splatterSize; k++ {
				// Деформируем круг для более органичной формы
				deform := 1.0 + 0.3*tg.noise.Perlin2D(float64(j+x)*0.2, float64(k+y)*0.2, 0.5)
				if (j*j + k*k) <= int(float64(splatterSize*splatterSize)*deform) {
					nx, ny := x+j, y+k
					if nx >= 0 && nx < width && ny >= 0 && ny < height {
						// Получаем текущий цвет
						r0, g0, b0, a0 := result.At(nx, ny).RGBA()

						// Смешиваем с цветом брызг
						blendFactor := 0.7
						r := uint8(float64(r0>>8)*(1-blendFactor) + float64(splatterColor.R)*blendFactor)
						g := uint8(float64(g0>>8)*(1-blendFactor) + float64(splatterColor.G)*blendFactor)
						b := uint8(float64(b0>>8)*(1-blendFactor) + float64(splatterColor.B)*blendFactor)

						result.Set(nx, ny, color.RGBA{r, g, b, uint8(a0 >> 8)})
					}
				}
			}
		}
	}

	return result
}
