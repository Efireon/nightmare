package render

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/font"

	"nightmare/internal/entity"
	"nightmare/internal/world"
)

// Константы рендеринга
const (
	TileSize   = 32 // Размер тайла в пикселях
	ViewRadius = 20 // Радиус видимости в тайлах
)

// Renderer представляет систему рендеринга
type Renderer struct {
	tileImages    map[world.TileType]*ebiten.Image
	entityImages  map[string]*ebiten.Image
	effectImages  map[string]*ebiten.Image
	itemImages    map[string]*ebiten.Image
	uiElements    map[string]*ebiten.Image
	fontNormal    font.Face
	fontLarge     font.Face
	screenWidth   int
	screenHeight  int
	viewOffsetX   float64
	viewOffsetY   float64
	shaderTime    float64
	screenEffects []ScreenEffect
}

// ScreenEffect представляет эффект на экране
type ScreenEffect struct {
	Type      string
	Duration  int
	Intensity float64
	Timer     int
}

// NewRenderer создает новый рендерер
func NewRenderer() (*Renderer, error) {
	r := &Renderer{
		tileImages:    make(map[world.TileType]*ebiten.Image),
		entityImages:  make(map[string]*ebiten.Image),
		effectImages:  make(map[string]*ebiten.Image),
		itemImages:    make(map[string]*ebiten.Image),
		uiElements:    make(map[string]*ebiten.Image),
		screenWidth:   800,
		screenHeight:  600,
		viewOffsetX:   0,
		viewOffsetY:   0,
		shaderTime:    0,
		screenEffects: []ScreenEffect{},
	}

	// Инициализируем изображения тайлов
	err := r.initTileImages()
	if err != nil {
		return nil, err
	}

	// Инициализируем шрифты
	err = r.initFonts()
	if err != nil {
		return nil, err
	}

	return r, nil
}

// initTileImages инициализирует изображения тайлов
func (r *Renderer) initTileImages() error {
	// В реальном проекте здесь будет загрузка тайлов из файлов или генерация
	// В этом примере мы просто создадим цветные прямоугольники

	// Создаем изображения тайлов разных типов
	r.tileImages[world.TileGrass] = createColoredTile(color.RGBA{0, 180, 0, 255})
	r.tileImages[world.TileForest] = createColoredTile(color.RGBA{0, 120, 0, 255})
	r.tileImages[world.TileDenseForest] = createColoredTile(color.RGBA{0, 80, 0, 255})
	r.tileImages[world.TilePath] = createColoredTile(color.RGBA{200, 190, 140, 255})
	r.tileImages[world.TileRocks] = createColoredTile(color.RGBA{120, 120, 120, 255})
	r.tileImages[world.TileWater] = createColoredTile(color.RGBA{0, 0, 180, 255})
	r.tileImages[world.TileSwamp] = createColoredTile(color.RGBA{70, 90, 70, 255})
	r.tileImages[world.TileCorrupted] = createColoredTile(color.RGBA{80, 0, 80, 255})

	// В будущем здесь будет процедурная генерация текстур на основе алгоритмов

	return nil
}

// createColoredTile создает тайл указанного цвета
func createColoredTile(clr color.Color) *ebiten.Image {
	img := ebiten.NewImage(TileSize, TileSize)
	img.Fill(clr)
	return img
}

// initFonts инициализирует шрифты
func (r *Renderer) initFonts() error {
	// В реальном проекте здесь будет загрузка шрифтов
	// В этом примере мы используем стандартный шрифт

	return nil
}

// DrawWorld отрисовывает мир
func (r *Renderer) DrawWorld(screen *ebiten.Image, w *world.World, player *entity.Player) {
	// Обновляем смещение вида
	r.viewOffsetX = player.Position.X
	r.viewOffsetY = player.Position.Y

	// Определяем видимый диапазон тайлов
	startX := int(r.viewOffsetX) - ViewRadius
	endX := int(r.viewOffsetX) + ViewRadius
	startY := int(r.viewOffsetY) - ViewRadius
	endY := int(r.viewOffsetY) + ViewRadius

	// Ограничиваем диапазон тайлов размерами мира
	if startX < 0 {
		startX = 0
	}
	if startY < 0 {
		startY = 0
	}
	if endX >= w.Width {
		endX = w.Width - 1
	}
	if endY >= w.Height {
		endY = w.Height - 1
	}

	// Отрисовываем видимые тайлы
	for y := startY; y <= endY; y++ {
		for x := startX; x <= endX; x++ {
			// Получаем тайл
			tile := w.GetTileAt(x, y)
			if tile == nil {
				continue
			}

			// Вычисляем координаты на экране
			screenX := int((float64(x)-r.viewOffsetX)*TileSize + float64(r.screenWidth)/2)
			screenY := int((float64(y)-r.viewOffsetY)*TileSize + float64(r.screenHeight)/2)

			// Проверяем, что тайл в пределах экрана
			if screenX+TileSize < 0 || screenY+TileSize < 0 ||
				screenX >= r.screenWidth || screenY >= r.screenHeight {
				continue
			}

			// Отрисовываем тайл
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(screenX), float64(screenY))

			// Применяем эффект коррупции
			op.ColorM.Scale(1-tile.Corruption*0.5, 1-tile.Corruption*0.7, 1-tile.Corruption*0.7, 1)

			screen.DrawImage(r.tileImages[tile.Type], op)

			// Отрисовываем объекты на тайле
			for _, obj := range tile.Objects {
				r.drawObject(screen, obj, screenX, screenY)
			}
		}
	}
}

// drawObject отрисовывает объект мира
func (r *Renderer) drawObject(screen *ebiten.Image, obj world.WorldObject, x, y int) {
	// В реальном проекте здесь будет отрисовка объектов
	// В этом примере мы просто нарисуем цветные прямоугольники

	switch obj.Type {
	case "tree":
		ebitenutil.DrawRect(screen, float64(x+TileSize/4), float64(y+TileSize/4),
			TileSize/2, TileSize/2, color.RGBA{0, 100, 0, 255})
	case "dense_tree":
		ebitenutil.DrawRect(screen, float64(x+TileSize/4), float64(y+TileSize/4),
			TileSize/2, TileSize/2, color.RGBA{0, 70, 0, 255})
	case "rock":
		ebitenutil.DrawRect(screen, float64(x+TileSize/3), float64(y+TileSize/3),
			TileSize/3, TileSize/3, color.RGBA{100, 100, 100, 255})
	}
}

// DrawEntities отрисовывает сущности
func (r *Renderer) DrawEntities(screen *ebiten.Image, entities []*world.Entity, player *entity.Player) {
	for _, e := range entities {
		// Вычисляем координаты на экране
		screenX := int((e.Position.X-r.viewOffsetX)*TileSize + float64(r.screenWidth)/2)
		screenY := int((e.Position.Y-r.viewOffsetY)*TileSize + float64(r.screenHeight)/2)

		// Проверяем, что сущность в пределах экрана
		if screenX+TileSize < 0 || screenY+TileSize < 0 ||
			screenX >= r.screenWidth || screenY >= r.screenHeight {
			continue
		}

		// Отрисовываем сущность
		r.drawEntity(screen, e, screenX, screenY)
	}
}

// drawEntity отрисовывает сущность
func (r *Renderer) drawEntity(screen *ebiten.Image, entity *world.Entity, x, y int) {
	// В реальном проекте здесь будет отрисовка сущности
	// В этом примере мы просто нарисуем цветной прямоугольник

	switch entity.Type {
	case "shadow":
		ebitenutil.DrawRect(screen, float64(x), float64(y),
			TileSize, TileSize, color.RGBA{20, 20, 20, 200})
	case "spider":
		ebitenutil.DrawRect(screen, float64(x+TileSize/4), float64(y+TileSize/4),
			TileSize/2, TileSize/2, color.RGBA{0, 0, 0, 255})
	default:
		ebitenutil.DrawRect(screen, float64(x+TileSize/4), float64(y+TileSize/4),
			TileSize/2, TileSize/2, color.RGBA{255, 0, 0, 255})
	}
}

// DrawPlayer отрисовывает игрока
func (r *Renderer) DrawPlayer(screen *ebiten.Image, player *entity.Player) {
	// Игрок всегда в центре экрана
	x := r.screenWidth / 2
	y := r.screenHeight / 2

	// Отрисовываем игрока
	ebitenutil.DrawRect(screen, float64(x-TileSize/2), float64(y-TileSize/2),
		TileSize, TileSize, color.RGBA{255, 255, 0, 255})

	// Отрисовываем направление игрока
	endX := float64(x) + math.Cos(player.Direction)*TileSize
	endY := float64(y) + math.Sin(player.Direction)*TileSize
	ebitenutil.DrawLine(screen, float64(x), float64(y), endX, endY, color.RGBA{255, 0, 0, 255})
}

// DrawUI отрисовывает пользовательский интерфейс
func (r *Renderer) DrawUI(screen *ebiten.Image, player *entity.Player) {
	// Отрисовываем здоровье
	ebitenutil.DrawRect(screen, 20, 20, 200, 20, color.RGBA{100, 100, 100, 255})
	healthWidth := int(player.Health / entity.MaxHealth * 200)
	ebitenutil.DrawRect(screen, 20, 20, float64(healthWidth), 20, color.RGBA{255, 0, 0, 255})
	ebitenutil.DebugPrintAt(screen, "Health", 25, 22)

	// Отрисовываем рассудок
	ebitenutil.DrawRect(screen, 20, 50, 200, 20, color.RGBA{100, 100, 100, 255})
	sanityWidth := int(player.Sanity / entity.MaxSanity * 200)
	ebitenutil.DrawRect(screen, 20, 50, float64(sanityWidth), 20, color.RGBA{0, 0, 255, 255})
	ebitenutil.DebugPrintAt(screen, "Sanity", 25, 52)
}

// DrawMainMenu отрисовывает главное меню
func (r *Renderer) DrawMainMenu(screen *ebiten.Image) {
	// Отрисовываем фон
	screen.Fill(color.RGBA{0, 0, 0, 255})

	// Отрисовываем заголовок
	ebitenutil.DebugPrintAt(screen, "NIGHTMARE FOREST", r.screenWidth/2-70, r.screenHeight/3)

	// Отрисовываем инструкции
	ebitenutil.DebugPrintAt(screen, "Press ENTER to start", r.screenWidth/2-70, r.screenHeight/2)
	ebitenutil.DebugPrintAt(screen, "WASD - move, ESC - pause", r.screenWidth/2-90, r.screenHeight/2+30)
}

// DrawPauseMenu отрисовывает меню паузы
func (r *Renderer) DrawPauseMenu(screen *ebiten.Image) {
	// Затемняем экран
	pauseOverlay := ebiten.NewImage(r.screenWidth, r.screenHeight)
	pauseOverlay.Fill(color.RGBA{0, 0, 0, 128})
	screen.DrawImage(pauseOverlay, nil)

	// Отрисовываем заголовок
	ebitenutil.DebugPrintAt(screen, "PAUSED", r.screenWidth/2-30, r.screenHeight/3)

	// Отрисовываем инструкции
	ebitenutil.DebugPrintAt(screen, "Press ESC to continue", r.screenWidth/2-70, r.screenHeight/2)
}

// DrawGameOver отрисовывает экран окончания игры
func (r *Renderer) DrawGameOver(screen *ebiten.Image) {
	// Отрисовываем фон
	screen.Fill(color.RGBA{0, 0, 0, 255})

	// Отрисовываем заголовок
	ebitenutil.DebugPrintAt(screen, "GAME OVER", r.screenWidth/2-40, r.screenHeight/3)

	// Отрисовываем инструкции
	ebitenutil.DebugPrintAt(screen, "Press ENTER to restart", r.screenWidth/2-70, r.screenHeight/2)
}

// AddScreenEffect добавляет эффект на экран
func (r *Renderer) AddScreenEffect(effectType string, duration int, intensity float64) {
	effect := ScreenEffect{
		Type:      effectType,
		Duration:  duration,
		Intensity: intensity,
		Timer:     duration,
	}
	r.screenEffects = append(r.screenEffects, effect)
}

// UpdateScreenEffects обновляет экранные эффекты
func (r *Renderer) UpdateScreenEffects() {
	// Обновляем таймеры эффектов
	for i := 0; i < len(r.screenEffects); i++ {
		r.screenEffects[i].Timer--
		if r.screenEffects[i].Timer <= 0 {
			// Удаляем эффект, если время истекло
			r.screenEffects = append(r.screenEffects[:i], r.screenEffects[i+1:]...)
			i--
		}
	}

	// Увеличиваем время для шейдеров
	r.shaderTime += 0.016
}

// ApplyScreenEffects применяет экранные эффекты
func (r *Renderer) ApplyScreenEffects(screen *ebiten.Image) {
	// Применяем эффекты
	for _, effect := range r.screenEffects {
		switch effect.Type {
		case "flash":
			// Создаем вспышку
			flash := ebiten.NewImage(r.screenWidth, r.screenHeight)
			flash.Fill(color.RGBA{255, 255, 255, uint8(effect.Intensity * 255)})
			screen.DrawImage(flash, nil)

		case "shake":
			// В реальном проекте здесь будет эффект тряски
			// В этом примере просто сдвигаем экран

		case "vignette":
			// В реальном проекте здесь будет эффект виньетки
			// В этом примере просто затемняем края

		case "distortion":
			// В реальном проекте здесь будет эффект искажения
			// В этом примере ничего не делаем
		}
	}
}
