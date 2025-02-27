package util

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// RandomGenerator предоставляет методы для генерации случайных значений
type RandomGenerator struct {
	rand *rand.Rand
}

// NewRandomGenerator создает новый генератор случайных чисел
func NewRandomGenerator(seed int64) *RandomGenerator {
	if seed == 0 {
		seed = time.Now().UnixNano()
	}

	return &RandomGenerator{
		rand: rand.New(rand.NewSource(seed)),
	}
}

// Float64 возвращает случайное число с плавающей точкой в диапазоне [0.0, 1.0)
func (r *RandomGenerator) Float64() float64 {
	return r.rand.Float64()
}

// Range возвращает случайное число с плавающей точкой в диапазоне [min, max)
func (r *RandomGenerator) Range(min, max float64) float64 {
	return min + r.rand.Float64()*(max-min)
}

// RangeInt возвращает случайное целое число в диапазоне [min, max)
func (r *RandomGenerator) RangeInt(min, max int) int {
	return min + r.rand.Intn(max-min)
}

// Chance возвращает true с вероятностью probability
func (r *RandomGenerator) Chance(probability float64) bool {
	return r.rand.Float64() < probability
}

// Choose выбирает случайный элемент из среза
func (r *RandomGenerator) Choose(items []interface{}) interface{} {
	if len(items) == 0 {
		return nil
	}
	return items[r.rand.Intn(len(items))]
}

// ChooseString выбирает случайную строку из среза
func (r *RandomGenerator) ChooseString(items []string) string {
	if len(items) == 0 {
		return ""
	}
	return items[r.rand.Intn(len(items))]
}

// ChooseInt выбирает случайное целое число из среза
func (r *RandomGenerator) ChooseInt(items []int) int {
	if len(items) == 0 {
		return 0
	}
	return items[r.rand.Intn(len(items))]
}

// ChooseFloat64 выбирает случайное число с плавающей точкой из среза
func (r *RandomGenerator) ChooseFloat64(items []float64) float64 {
	if len(items) == 0 {
		return 0
	}
	return items[r.rand.Intn(len(items))]
}

// WeightedChoiceIndex выбирает случайный индекс с учетом весов
func (r *RandomGenerator) WeightedChoiceIndex(weights []float64) int {
	if len(weights) == 0 {
		return -1
	}

	// Вычисляем сумму весов
	totalWeight := 0.0
	for _, weight := range weights {
		totalWeight += weight
	}

	// Генерируем случайное число в диапазоне [0, totalWeight)
	choice := r.rand.Float64() * totalWeight

	// Находим соответствующий индекс
	currentWeight := 0.0
	for i, weight := range weights {
		currentWeight += weight
		if choice < currentWeight {
			return i
		}
	}

	// На случай ошибок округления
	return len(weights) - 1
}

// WeightedChoiceString выбирает случайную строку с учетом весов
func (r *RandomGenerator) WeightedChoiceString(items []string, weights []float64) string {
	if len(items) == 0 || len(items) != len(weights) {
		return ""
	}

	index := r.WeightedChoiceIndex(weights)
	if index >= 0 && index < len(items) {
		return items[index]
	}

	return ""
}

// WeightedChoiceInt выбирает случайное целое число с учетом весов
func (r *RandomGenerator) WeightedChoiceInt(items []int, weights []float64) int {
	if len(items) == 0 || len(items) != len(weights) {
		return 0
	}

	index := r.WeightedChoiceIndex(weights)
	if index >= 0 && index < len(items) {
		return items[index]
	}

	return 0
}

// Shuffle перемешивает элементы среза
func (r *RandomGenerator) Shuffle(items []interface{}) {
	for i := len(items) - 1; i > 0; i-- {
		j := r.rand.Intn(i + 1)
		items[i], items[j] = items[j], items[i]
	}
}

// ShuffleStrings перемешивает элементы среза строк
func (r *RandomGenerator) ShuffleStrings(items []string) {
	for i := len(items) - 1; i > 0; i-- {
		j := r.rand.Intn(i + 1)
		items[i], items[j] = items[j], items[i]
	}
}

// ShuffleInts перемешивает элементы среза целых чисел
func (r *RandomGenerator) ShuffleInts(items []int) {
	for i := len(items) - 1; i > 0; i-- {
		j := r.rand.Intn(i + 1)
		items[i], items[j] = items[j], items[i]
	}
}

// ShuffleFloat64s перемешивает элементы среза чисел с плавающей точкой
func (r *RandomGenerator) ShuffleFloat64s(items []float64) {
	for i := len(items) - 1; i > 0; i-- {
		j := r.rand.Intn(i + 1)
		items[i], items[j] = items[j], items[i]
	}
}

// NormalDistribution возвращает случайное число из нормального распределения
func (r *RandomGenerator) NormalDistribution(mean, stdDev float64) float64 {
	// Используем преобразование Бокса-Мюллера
	u1 := r.rand.Float64()
	u2 := r.rand.Float64()

	// Преобразуем равномерное распределение в нормальное
	z := math.Sqrt(-2.0*math.Log(u1)) * math.Cos(2.0*math.Pi*u2)

	// Масштабируем и смещаем
	return mean + z*stdDev
}

// ExponentialDistribution возвращает случайное число из экспоненциального распределения
func (r *RandomGenerator) ExponentialDistribution(lambda float64) float64 {
	return -math.Log(1.0-r.rand.Float64()) / lambda
}

// PoissonDistribution возвращает случайное число из распределения Пуассона
func (r *RandomGenerator) PoissonDistribution(lambda float64) int {
	if lambda <= 0 {
		return 0
	}

	// Используем алгоритм Кнута
	L := math.Exp(-lambda)
	k := 0
	p := 1.0

	for p > L {
		k++
		p *= r.rand.Float64()
	}

	return k - 1
}

// GenerateUUID создает случайный UUID
func (r *RandomGenerator) GenerateUUID() string {
	uuid := make([]byte, 16)
	r.rand.Read(uuid)

	// Устанавливаем версию UUID (версия 4 - случайный UUID)
	uuid[6] = (uuid[6] & 0x0F) | 0x40
	uuid[8] = (uuid[8] & 0x3F) | 0x80

	return fmt.Sprintf("%x-%x-%x-%x-%x",
		uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

// RandomName генерирует случайное имя
func (r *RandomGenerator) RandomName() string {
	prefixes := []string{
		"dark", "shadow", "night", "blood", "death", "fear", "dread", "grim",
		"fell", "doom", "black", "horror", "terror", "spectral", "phantom",
	}

	roots := []string{
		"wood", "forest", "marsh", "mist", "fog", "grave", "crypt", "tomb",
		"hell", "void", "abyss", "haunt", "hollow", "canyon", "cave", "nightmare",
	}

	suffixes := []string{
		"walker", "stalker", "watcher", "hunter", "keeper", "lurker", "dweller",
		"crawler", "spirit", "wraith", "spectre", "shade", "beast", "fiend",
	}

	// Выбираем компоненты имени
	var parts []string

	if r.Chance(0.7) {
		parts = append(parts, r.ChooseString(prefixes))
	}

	parts = append(parts, r.ChooseString(roots))

	if r.Chance(0.6) {
		parts = append(parts, r.ChooseString(suffixes))
	}

	// Соединяем части
	name := ""
	for i, part := range parts {
		if i == 0 {
			name = part
		} else {
			name = name + r.ChooseString([]string{"", " "}) + part
		}
	}

	// Иногда добавляем числовой суффикс
	if r.Chance(0.2) {
		name = name + r.ChooseString([]string{"", " "}) + fmt.Sprintf("%d", r.RangeInt(1, 100))
	}

	return name
}
