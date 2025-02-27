package util

import (
	"math"
	"math/rand"

	"github.com/ojrac/opensimplex-go"
)

// NoiseGenerator генерирует различные виды шума для процедурной генерации
type NoiseGenerator struct {
	simplex opensimplex.Noise
	seed    int64
}

// NewNoiseGenerator создает новый генератор шума
func NewNoiseGenerator(seed int64) *NoiseGenerator {
	return &NoiseGenerator{
		simplex: opensimplex.New(seed),
		seed:    seed,
	}
}

// Perlin2D генерирует 2D шум Перлина
func (n *NoiseGenerator) Perlin2D(x, y float64, scale float64) float64 {
	return n.fbm2D(x, y, 6, 2.0, 0.5, scale)
}

// Perlin3D генерирует 3D шум Перлина
func (n *NoiseGenerator) Perlin3D(x, y, z float64, scale float64) float64 {
	return n.fbm3D(x, y, z, 6, 2.0, 0.5, scale)
}

// fbm2D генерирует дробное броуновское движение в 2D
func (n *NoiseGenerator) fbm2D(x, y float64, octaves int, lacunarity, persistence, scale float64) float64 {
	// Используем Simplex шум как основу
	value := 0.0
	amplitude := 1.0
	frequency := scale
	maxValue := 0.0

	for i := 0; i < octaves; i++ {
		// Добавляем шум с текущей частотой и амплитудой
		value += n.simplex.Eval2(x*frequency, y*frequency) * amplitude

		// Отслеживаем максимальное значение для нормализации
		maxValue += amplitude

		// Увеличиваем частоту и уменьшаем амплитуду для следующей октавы
		amplitude *= persistence
		frequency *= lacunarity
	}

	// Нормализуем в диапазон [0, 1]
	return (value/maxValue + 1) / 2
}

// fbm3D генерирует дробное броуновское движение в 3D
func (n *NoiseGenerator) fbm3D(x, y, z float64, octaves int, lacunarity, persistence, scale float64) float64 {
	// Используем Simplex шум как основу
	value := 0.0
	amplitude := 1.0
	frequency := scale
	maxValue := 0.0

	for i := 0; i < octaves; i++ {
		// Добавляем шум с текущей частотой и амплитудой
		value += n.simplex.Eval3(x*frequency, y*frequency, z*frequency) * amplitude

		// Отслеживаем максимальное значение для нормализации
		maxValue += amplitude

		// Увеличиваем частоту и уменьшаем амплитуду для следующей октавы
		amplitude *= persistence
		frequency *= lacunarity
	}

	// Нормализуем в диапазон [0, 1]
	return (value/maxValue + 1) / 2
}

// WorleyNoise генерирует шум Ворлея (ячеистый шум)
func (n *NoiseGenerator) WorleyNoise(x, y float64, numPoints int, scale float64) float64 {
	// Создаем случайные точки
	r := rand.New(rand.NewSource(n.seed))
	points := make([][2]float64, numPoints)

	for i := 0; i < numPoints; i++ {
		points[i][0] = r.Float64() * scale
		points[i][1] = r.Float64() * scale
	}

	// Находим ближайшую точку
	minDist := math.MaxFloat64

	// Масштабируем входные координаты
	x = math.Mod(x*scale, scale)
	y = math.Mod(y*scale, scale)

	for i := 0; i < numPoints; i++ {
		dx := x - points[i][0]
		dy := y - points[i][1]

		// Используем евклидово расстояние
		dist := math.Sqrt(dx*dx + dy*dy)

		if dist < minDist {
			minDist = dist
		}
	}

	// Нормализуем в диапазон [0, 1]
	return math.Min(minDist/(scale/2), 1.0)
}

// RidgedNoise генерирует ребристый шум
func (n *NoiseGenerator) RidgedNoise(x, y float64, scale float64) float64 {
	// Получаем базовый шум
	value := n.Perlin2D(x, y, scale)

	// Преобразуем в ребристый шум
	value = 1.0 - math.Abs(value-0.5)*2

	return value
}

// DomainWarp применяет искажение домена
func (n *NoiseGenerator) DomainWarp(x, y float64, strength, scale float64) (float64, float64) {
	// Искажаем координаты с помощью шума
	warpX := n.Perlin2D(x, y, scale) * strength
	warpY := n.Perlin2D(x+100, y+100, scale) * strength

	return x + warpX, y + warpY
}

// CreateNoiseTexture создает текстуру на основе шума
func (n *NoiseGenerator) CreateNoiseTexture(width, height int, scale float64,
	noiseFunc func(x, y float64, scale float64) float64) [][]float64 {

	// Создаем двумерный массив для текстуры
	texture := make([][]float64, height)
	for y := range texture {
		texture[y] = make([]float64, width)
	}

	// Заполняем массив значениями шума
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Нормализуем координаты в диапазон [0, 1]
			nx := float64(x) / float64(width)
			ny := float64(y) / float64(height)

			// Получаем значение шума
			texture[y][x] = noiseFunc(nx, ny, scale)
		}
	}

	return texture
}

// CreateHeightmap создает карту высот
func (n *NoiseGenerator) CreateHeightmap(width, height int, scale, mountainScale float64) [][]float64 {
	// Создаем базовую карту высот с Perlin шумом
	heightmap := n.CreateNoiseTexture(width, height, scale, n.Perlin2D)

	// Добавляем горы с ребристым шумом
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			nx := float64(x) / float64(width)
			ny := float64(y) / float64(height)

			// Получаем ребристый шум
			mountain := n.RidgedNoise(nx, ny, mountainScale)

			// Смешиваем с базовой картой высот
			heightmap[y][x] = heightmap[y][x]*0.7 + mountain*0.3
		}
	}

	return heightmap
}

// ApplyCorruption применяет эффект "порчи" к текстуре
func (n *NoiseGenerator) ApplyCorruption(texture [][]float64, x, y, radius, intensity float64) {
	width := len(texture[0])
	height := len(texture)

	// Применяем эффект порчи в указанном радиусе
	radiusSq := radius * radius

	for j := int(y - radius); j <= int(y+radius); j++ {
		for i := int(x - radius); i <= int(x+radius); i++ {
			// Проверяем, что координаты в пределах текстуры
			if i < 0 || j < 0 || i >= width || j >= height {
				continue
			}

			// Вычисляем расстояние до центра
			dx := float64(i) - x
			dy := float64(j) - y
			distSq := dx*dx + dy*dy

			// Если точка внутри радиуса
			if distSq <= radiusSq {
				// Вычисляем силу воздействия (ближе к центру - сильнее)
				strength := (1.0 - distSq/radiusSq) * intensity

				// Добавляем "порчу"
				corruption := n.Perlin2D(float64(i)*0.1, float64(j)*0.1, 0.5) * strength
				texture[j][i] = (texture[j][i] + corruption) / (1.0 + strength)
			}
		}
	}
}

// CreateTextureFromCorruption создает текстуру на основе уровня "порчи"
func (n *NoiseGenerator) CreateTextureFromCorruption(width, height int, corruptionMap [][]float64) [][]float64 {
	// Создаем текстуру той же размерности
	texture := make([][]float64, height)
	for y := range texture {
		texture[y] = make([]float64, width)
	}

	// Заполняем текстуру на основе уровня "порчи"
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Получаем базовый шум
			baseNoise := n.Perlin2D(float64(x)*0.1, float64(y)*0.1, 0.5)

			// Искажаем координаты в зависимости от уровня "порчи"
			warpStrength := corruptionMap[y][x] * 10.0
			warpX, warpY := n.DomainWarp(float64(x)*0.05, float64(y)*0.05, warpStrength, 0.5)

			// Получаем искаженный шум
			warpedNoise := n.Perlin2D(warpX, warpY, 0.5)

			// Смешиваем шумы в зависимости от уровня "порчи"
			texture[y][x] = baseNoise*(1.0-corruptionMap[y][x]) + warpedNoise*corruptionMap[y][x]
		}
	}

	return texture
}
