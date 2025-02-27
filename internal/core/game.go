package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"nightmare/internal/ai"
	"nightmare/internal/entity"
	"nightmare/internal/render"
	"nightmare/internal/world"
)

// GameState представляет состояние игры
type GameState int

const (
	StateMainMenu GameState = iota
	StatePlaying
	StatePaused
	StateGameOver
)

// Game реализует интерфейс ebiten.Game
type Game struct {
	state      GameState
	player     *entity.Player
	world      *world.World
	renderer   *render.Renderer
	director   *ai.Director
	frameCount int
}

// NewGame создает новую игру
func NewGame() (*Game, error) {
	// Создаем игрока
	player := entity.NewPlayer()

	// Создаем мир
	world, err := world.NewWorld(256, 256)
	if err != nil {
		return nil, err
	}

	// Создаем рендерер
	renderer, err := render.NewRenderer()
	if err != nil {
		return nil, err
	}

	// Создаем ИИ-директора
	director := ai.NewDirector(player, world)

	return &Game{
		state:      StateMainMenu,
		player:     player,
		world:      world,
		renderer:   renderer,
		director:   director,
		frameCount: 0,
	}, nil
}

// Update обновляет состояние игры
func (g *Game) Update() error {
	g.frameCount++

	switch g.state {
	case StateMainMenu:
		// Обработка ввода в главном меню
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.state = StatePlaying
		}

	case StatePlaying:
		// Обработка ввода игрока
		g.handlePlayerInput()

		// Обновление мира
		g.world.Update()

		// Обновление игрока
		g.player.Update()

		// Обновление ИИ-директора каждые 30 кадров (примерно 0.5 сек)
		if g.frameCount%30 == 0 {
			g.director.AnalyzePlayerBehavior()
			g.director.AdjustWorld()
		}

		// Проверка условий окончания игры
		if g.player.Health <= 0 || g.player.Sanity <= 0 {
			g.state = StateGameOver
		}

	case StatePaused:
		// Обработка ввода в меню паузы
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.state = StatePlaying
		}

	case StateGameOver:
		// Обработка ввода в меню окончания игры
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			// Сбросить игру
			g.resetGame()
		}
	}

	return nil
}

// Draw отрисовывает игру
func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {
	case StateMainMenu:
		g.renderer.DrawMainMenu(screen)

	case StatePlaying, StatePaused:
		// Отрисовка мира
		g.renderer.DrawWorld(screen, g.world, g.player)

		// Отрисовка существ
		g.renderer.DrawEntities(screen, g.world.Entities, g.player)

		// Отрисовка игрока
		g.renderer.DrawPlayer(screen, g.player)

		// Отрисовка UI
		g.renderer.DrawUI(screen, g.player)

		if g.state == StatePaused {
			g.renderer.DrawPauseMenu(screen)
		}

	case StateGameOver:
		g.renderer.DrawGameOver(screen)
	}
}

// Layout возвращает размер игры
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 800, 600
}

// handlePlayerInput обрабатывает ввод игрока
func (g *Game) handlePlayerInput() {
	// Движение
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.player.MoveForward()
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.player.MoveBackward()
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.player.TurnLeft()
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.player.TurnRight()
	}

	// Взаимодействие
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		g.player.Interact(g.world)
	}

	// Пауза
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		if g.state == StatePlaying {
			g.state = StatePaused
		} else if g.state == StatePaused {
			g.state = StatePlaying
		}
	}
}

// resetGame сбрасывает игру
func (g *Game) resetGame() {
	g.player = entity.NewPlayer()

	var err error
	g.world, err = world.NewWorld(256, 256)
	if err != nil {
		panic(err) // В реальной игре нужно обработать ошибку более изящно
	}

	g.director = ai.NewDirector(g.player, g.world)
	g.state = StateMainMenu
	g.frameCount = 0
}
