package ui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"nightmare/internal/entity"
	"nightmare/internal/event"
	"nightmare/internal/item"
)

// UIElement базовый интерфейс для элементов UI
type UIElement interface {
	Update() error
	Draw(screen *ebiten.Image)
	GetRect() (x, y, width, height int)
	SetPosition(x, y int)
	IsVisible() bool
	SetVisible(visible bool)
	IsEnabled() bool
	SetEnabled(enabled bool)
	HandleInput() bool
}

// BaseElement базовая структура для элементов UI
type BaseElement struct {
	X           int
	Y           int
	Width       int
	Height      int
	Visible     bool
	Enabled     bool
	FontSize    int
	TextColor   color.RGBA
	BgColor     color.RGBA
	BorderColor color.RGBA
	BorderWidth int
	Padding     int
	OnClick     func()
}

// NewBaseElement создает новый базовый элемент UI
func NewBaseElement(x, y, width, height int) BaseElement {
	return BaseElement{
		X:           x,
		Y:           y,
		Width:       width,
		Height:      height,
		Visible:     true,
		Enabled:     true,
		FontSize:    16,
		TextColor:   color.RGBA{255, 255, 255, 255},
		BgColor:     color.RGBA{50, 50, 50, 200},
		BorderColor: color.RGBA{100, 100, 100, 255},
		BorderWidth: 1,
		Padding:     5,
	}
}

// GetRect возвращает прямоугольник элемента
func (e *BaseElement) GetRect() (x, y, width, height int) {
	return e.X, e.Y, e.Width, e.Height
}

// SetPosition устанавливает позицию элемента
func (e *BaseElement) SetPosition(x, y int) {
	e.X = x
	e.Y = y
}

// IsVisible возвращает видимость элемента
func (e *BaseElement) IsVisible() bool {
	return e.Visible
}

// SetVisible устанавливает видимость элемента
func (e *BaseElement) SetVisible(visible bool) {
	e.Visible = visible
}

// IsEnabled возвращает активность элемента
func (e *BaseElement) IsEnabled() bool {
	return e.Enabled
}

// SetEnabled устанавливает активность элемента
func (e *BaseElement) SetEnabled(enabled bool) {
	e.Enabled = enabled
}

// ContainsPoint проверяет, содержит ли элемент указанную точку
func (e *BaseElement) ContainsPoint(x, y int) bool {
	return x >= e.X && x < e.X+e.Width && y >= e.Y && y < e.Y+e.Height
}

// DrawBackground отрисовывает фон элемента
func (e *BaseElement) DrawBackground(screen *ebiten.Image) {
	// Отрисовываем фон
	ebitenutil.DrawRect(screen, float64(e.X), float64(e.Y), float64(e.Width), float64(e.Height), e.BgColor)

	// Отрисовываем рамку, если нужно
	if e.BorderWidth > 0 {
		for i := 0; i < e.BorderWidth; i++ {
			ebitenutil.DrawRect(screen, float64(e.X-i), float64(e.Y-i), float64(e.Width+i*2), 1, e.BorderColor)
			ebitenutil.DrawRect(screen, float64(e.X-i), float64(e.Y+e.Height+i-1), float64(e.Width+i*2), 1, e.BorderColor)
			ebitenutil.DrawRect(screen, float64(e.X-i), float64(e.Y-i), 1, float64(e.Height+i*2), e.BorderColor)
			ebitenutil.DrawRect(screen, float64(e.X+e.Width+i-1), float64(e.Y-i), 1, float64(e.Height+i*2), e.BorderColor)
		}
	}
}

// Label элемент UI для отображения текста
type Label struct {
	BaseElement
	Text      string
	Alignment int // 0 = left, 1 = center, 2 = right
}

// NewLabel создает новую метку
func NewLabel(x, y, width, height int, text string) *Label {
	base := NewBaseElement(x, y, width, height)
	base.BgColor = color.RGBA{0, 0, 0, 0} // Прозрачный фон по умолчанию
	base.BorderWidth = 0                  // Без рамки по умолчанию

	return &Label{
		BaseElement: base,
		Text:        text,
		Alignment:   0, // По левому краю по умолчанию
	}
}

// Update обновляет состояние метки
func (l *Label) Update() error {
	return nil
}

// Draw отрисовывает метку
func (l *Label) Draw(screen *ebiten.Image) {
	if !l.Visible {
		return
	}

	// Отрисовываем фон
	l.DrawBackground(screen)

	// Отрисовываем текст с учетом выравнивания
	textX := l.X + l.Padding
	if l.Alignment == 1 {
		// По центру
		textWidth := len(l.Text) * l.FontSize / 2 // Приблизительно
		textX = l.X + (l.Width-textWidth)/2
	} else if l.Alignment == 2 {
		// По правому краю
		textWidth := len(l.Text) * l.FontSize / 2 // Приблизительно
		textX = l.X + l.Width - textWidth - l.Padding
	}

	// В реальной игре здесь был бы код для отрисовки текста с учетом font face
	// Для примера используем упрощенный вариант
	ebitenutil.DebugPrintAt(screen, l.Text, textX, l.Y+l.Height/2-l.FontSize/2)
}

// HandleInput обрабатывает ввод для метки
func (l *Label) HandleInput() bool {
	return false // Метка не обрабатывает ввод
}

// Button элемент UI для кнопки
type Button struct {
	BaseElement
	Text         string
	HoverColor   color.RGBA
	PressedColor color.RGBA
	IsHovered    bool
	IsPressed    bool
}

// NewButton создает новую кнопку
func NewButton(x, y, width, height int, text string, onClick func()) *Button {
	base := NewBaseElement(x, y, width, height)
	base.OnClick = onClick

	return &Button{
		BaseElement:  base,
		Text:         text,
		HoverColor:   color.RGBA{70, 70, 70, 200},
		PressedColor: color.RGBA{30, 30, 30, 200},
	}
}

// Update обновляет состояние кнопки
func (b *Button) Update() error {
	if !b.Visible || !b.Enabled {
		return nil
	}

	// Получаем позицию курсора
	x, y := ebiten.CursorPosition()

	// Проверяем, находится ли курсор над кнопкой
	b.IsHovered = b.ContainsPoint(x, y)

	// Проверяем, нажата ли кнопка
	if b.IsHovered {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			b.IsPressed = true
		} else if b.IsPressed && inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			// Кнопка была отпущена, вызываем обработчик нажатия
			b.IsPressed = false
			if b.OnClick != nil {
				b.OnClick()
			}
		} else {
			b.IsPressed = false
		}
	} else {
		b.IsPressed = false
	}

	return nil
}

// Draw отрисовывает кнопку
func (b *Button) Draw(screen *ebiten.Image) {
	if !b.Visible {
		return
	}

	// Выбираем цвет в зависимости от состояния
	bgColor := b.BgColor
	if b.IsPressed {
		bgColor = b.PressedColor
	} else if b.IsHovered {
		bgColor = b.HoverColor
	}

	// Сохраняем оригинальный цвет
	origColor := b.BgColor
	b.BgColor = bgColor

	// Отрисовываем фон
	b.DrawBackground(screen)

	// Восстанавливаем оригинальный цвет
	b.BgColor = origColor

	// Отрисовываем текст по центру
	textX := b.X + (b.Width-len(b.Text)*b.FontSize/2)/2
	textY := b.Y + b.Height/2 - b.FontSize/2

	// В реальной игре здесь был бы код для отрисовки текста с учетом font face
	ebitenutil.DebugPrintAt(screen, b.Text, textX, textY)
}

// HandleInput обрабатывает ввод для кнопки
func (b *Button) HandleInput() bool {
	if !b.Visible || !b.Enabled {
		return false
	}

	// Получаем позицию курсора
	x, y := ebiten.CursorPosition()

	// Проверяем, находится ли курсор над кнопкой
	if b.ContainsPoint(x, y) {
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			// Вызываем обработчик нажатия
			if b.OnClick != nil {
				b.OnClick()
			}
			return true
		}
	}

	return false
}

// Panel элемент UI для панели
type Panel struct {
	BaseElement
	Children []UIElement
}

// NewPanel создает новую панель
func NewPanel(x, y, width, height int) *Panel {
	base := NewBaseElement(x, y, width, height)

	return &Panel{
		BaseElement: base,
		Children:    []UIElement{},
	}
}

// Update обновляет состояние панели и ее дочерних элементов
func (p *Panel) Update() error {
	if !p.Visible {
		return nil
	}

	for _, child := range p.Children {
		if err := child.Update(); err != nil {
			return err
		}
	}

	return nil
}

// Draw отрисовывает панель и ее дочерние элементы
func (p *Panel) Draw(screen *ebiten.Image) {
	if !p.Visible {
		return
	}

	// Отрисовываем фон
	p.DrawBackground(screen)

	// Отрисовываем дочерние элементы
	for _, child := range p.Children {
		child.Draw(screen)
	}
}

// HandleInput обрабатывает ввод для панели и ее дочерних элементов
func (p *Panel) HandleInput() bool {
	if !p.Visible || !p.Enabled {
		return false
	}

	// Обрабатываем ввод для дочерних элементов
	for _, child := range p.Children {
		if child.HandleInput() {
			return true
		}
	}

	return false
}

// AddChild добавляет дочерний элемент
func (p *Panel) AddChild(child UIElement) {
	p.Children = append(p.Children, child)
}

// RemoveChild удаляет дочерний элемент
func (p *Panel) RemoveChild(child UIElement) {
	for i, c := range p.Children {
		if c == child {
			p.Children = append(p.Children[:i], p.Children[i+1:]...)
			return
		}
	}
}

// ProgressBar элемент UI для индикатора прогресса
type ProgressBar struct {
	BaseElement
	Value         float64 // от 0 до 1
	ProgressColor color.RGBA
	ShowText      bool
	Text          string
}

// NewProgressBar создает новый индикатор прогресса
func NewProgressBar(x, y, width, height int) *ProgressBar {
	base := NewBaseElement(x, y, width, height)

	return &ProgressBar{
		BaseElement:   base,
		Value:         0,
		ProgressColor: color.RGBA{0, 200, 0, 255},
		ShowText:      true,
		Text:          "",
	}
}

// Update обновляет состояние индикатора прогресса
func (p *ProgressBar) Update() error {
	return nil
}

// Draw отрисовывает индикатор прогресса
func (p *ProgressBar) Draw(screen *ebiten.Image) {
	if !p.Visible {
		return
	}

	// Отрисовываем фон
	p.DrawBackground(screen)

	// Отрисовываем прогресс
	progressWidth := int(float64(p.Width-2*p.BorderWidth) * p.Value)
	if progressWidth > 0 {
		ebitenutil.DrawRect(
			screen,
			float64(p.X+p.BorderWidth),
			float64(p.Y+p.BorderWidth),
			float64(progressWidth),
			float64(p.Height-2*p.BorderWidth),
			p.ProgressColor,
		)
	}

	// Отрисовываем текст, если нужно
	if p.ShowText {
		text := p.Text
		if text == "" {
			text = fmt.Sprintf("%d%%", int(p.Value*100))
		}

		textX := p.X + (p.Width-len(text)*p.FontSize/2)/2
		textY := p.Y + p.Height/2 - p.FontSize/2

		// В реальной игре здесь был бы код для отрисовки текста с учетом font face
		ebitenutil.DebugPrintAt(screen, text, textX, textY)
	}
}

// HandleInput обрабатывает ввод для индикатора прогресса
func (p *ProgressBar) HandleInput() bool {
	return false // Индикатор прогресса не обрабатывает ввод
}

// SetValue устанавливает значение индикатора прогресса
func (p *ProgressBar) SetValue(value float64) {
	// Ограничиваем значение в диапазоне [0, 1]
	if value < 0 {
		value = 0
	} else if value > 1 {
		value = 1
	}
	p.Value = value
}

// UIManager управляет пользовательским интерфейсом
type UIManager struct {
	elements           []UIElement
	player             *entity.Player
	inventory          *item.Inventory
	eventManager       *event.EventManager
	screenWidth        int
	screenHeight       int
	isMenuVisible      bool
	isInventoryVisible bool
	currentScreen      string

	// Часто используемые элементы UI
	healthBar      *ProgressBar
	sanityBar      *ProgressBar
	messageLine    *Label
	menuPanel      *Panel
	inventoryPanel *Panel
}

// NewUIManager создает новый менеджер UI
func NewUIManager(player *entity.Player, inventory *item.Inventory, eventManager *event.EventManager, screenWidth, screenHeight int) *UIManager {
	ui := &UIManager{
		elements:           []UIElement{},
		player:             player,
		inventory:          inventory,
		eventManager:       eventManager,
		screenWidth:        screenWidth,
		screenHeight:       screenHeight,
		isMenuVisible:      false,
		isInventoryVisible: false,
		currentScreen:      "game",
	}

	// Создаем базовые элементы UI
	ui.createUI()

	// Подписываемся на события
	if eventManager != nil {
		ui.subscribeToEvents()
	}

	return ui
}

// createUI создает основные элементы UI
func (ui *UIManager) createUI() {
	// Создаем индикатор здоровья
	ui.healthBar = NewProgressBar(20, 20, 200, 20)
	ui.healthBar.ProgressColor = color.RGBA{200, 0, 0, 255}
	ui.healthBar.Text = "Health"
	ui.AddElement(ui.healthBar)

	// Создаем индикатор рассудка
	ui.sanityBar = NewProgressBar(20, 50, 200, 20)
	ui.sanityBar.ProgressColor = color.RGBA{0, 0, 200, 255}
	ui.sanityBar.Text = "Sanity"
	ui.AddElement(ui.sanityBar)

	// Создаем строку сообщений
	ui.messageLine = NewLabel(20, ui.screenHeight-40, ui.screenWidth-40, 20, "")
	ui.messageLine.Alignment = 1 // По центру
	ui.AddElement(ui.messageLine)

	// Создаем панель меню
	ui.createMenuPanel()

	// Создаем панель инвентаря
	ui.createInventoryPanel()
}

// createMenuPanel создает панель меню
func (ui *UIManager) createMenuPanel() {
	ui.menuPanel = NewPanel(ui.screenWidth/2-150, ui.screenHeight/2-200, 300, 400)
	ui.menuPanel.SetVisible(false)

	// Заголовок меню
	titleLabel := NewLabel(0, 20, 300, 40, "MENU")
	titleLabel.Alignment = 1 // По центру
	titleLabel.FontSize = 24
	ui.menuPanel.AddChild(titleLabel)

	// Кнопка продолжить
	continueButton := NewButton(50, 100, 200, 40, "Continue", func() {
		ui.ToggleMenu()
	})
	ui.menuPanel.AddChild(continueButton)

	// Кнопка настройки
	settingsButton := NewButton(50, 160, 200, 40, "Settings", func() {
		// Здесь будет код для открытия настроек
	})
	ui.menuPanel.AddChild(settingsButton)

	// Кнопка выход в главное меню
	mainMenuButton := NewButton(50, 220, 200, 40, "Main Menu", func() {
		ui.SetCurrentScreen("main_menu")
	})
	ui.menuPanel.AddChild(mainMenuButton)

	// Кнопка выход из игры
	exitButton := NewButton(50, 280, 200, 40, "Exit Game", func() {
		// Здесь будет код для выхода из игры
	})
	ui.menuPanel.AddChild(exitButton)

	ui.AddElement(ui.menuPanel)
}

// createInventoryPanel создает панель инвентаря
func (ui *UIManager) createInventoryPanel() {
	ui.inventoryPanel = NewPanel(ui.screenWidth-330, 80, 300, 400)
	ui.inventoryPanel.SetVisible(false)

	// Заголовок инвентаря
	titleLabel := NewLabel(0, 10, 300, 30, "INVENTORY")
	titleLabel.Alignment = 1 // По центру
	ui.inventoryPanel.AddChild(titleLabel)

	// Здесь будет код для отображения предметов в инвентаре
	// В реальной игре здесь был бы более сложный код для создания сетки предметов

	// Кнопка закрыть
	closeButton := NewButton(100, 360, 100, 30, "Close", func() {
		ui.ToggleInventory()
	})
	ui.inventoryPanel.AddChild(closeButton)

	ui.AddElement(ui.inventoryPanel)
}

// subscribeToEvents подписывается на события
func (ui *UIManager) subscribeToEvents() {
	// Подписываемся на событие изменения здоровья
	ui.eventManager.AddListener(event.EventPlayerDamaged, func(data event.EventData) {
		// Обновляем индикатор здоровья
		ui.updateHealthBar()
	})

	// Подписываемся на событие изменения рассудка
	ui.eventManager.AddListener(event.EventPlayerSanityChanged, func(data event.EventData) {
		// Обновляем индикатор рассудка
		ui.updateSanityBar()
	})

	// Подписываемся на событие обновления инвентаря
	ui.eventManager.AddCustomListener("inventory_updated", func(data event.EventData) {
		// Обновляем инвентарь
		ui.updateInventoryPanel()
	})
}

// AddElement добавляет элемент UI
func (ui *UIManager) AddElement(element UIElement) {
	ui.elements = append(ui.elements, element)
}

// RemoveElement удаляет элемент UI
func (ui *UIManager) RemoveElement(element UIElement) {
	for i, e := range ui.elements {
		if e == element {
			ui.elements = append(ui.elements[:i], ui.elements[i+1:]...)
			return
		}
	}
}

// Update обновляет состояние UI
func (ui *UIManager) Update() error {
	// Обновляем индикаторы здоровья и рассудка
	ui.updateHealthBar()
	ui.updateSanityBar()

	// Обновляем все элементы UI
	for _, element := range ui.elements {
		if err := element.Update(); err != nil {
			return err
		}
	}

	// Проверяем ввод для всего UI
	ui.handleInput()

	return nil
}

// Draw отрисовывает UI
func (ui *UIManager) Draw(screen *ebiten.Image) {
	// Отрисовываем все элементы UI
	for _, element := range ui.elements {
		element.Draw(screen)
	}
}

// handleInput обрабатывает ввод
func (ui *UIManager) handleInput() {
	// Обрабатываем ввод для всех элементов UI
	for _, element := range ui.elements {
		element.HandleInput()
	}

	// Проверяем нажатие клавиш для открытия/закрытия меню и инвентаря
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		ui.ToggleMenu()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyI) {
		ui.ToggleInventory()
	}
}

// ToggleMenu переключает видимость меню
func (ui *UIManager) ToggleMenu() {
	ui.isMenuVisible = !ui.isMenuVisible
	ui.menuPanel.SetVisible(ui.isMenuVisible)
}

// ToggleInventory переключает видимость инвентаря
func (ui *UIManager) ToggleInventory() {
	ui.isInventoryVisible = !ui.isInventoryVisible
	ui.inventoryPanel.SetVisible(ui.isInventoryVisible)

	// Если инвентарь стал видимым, обновляем его
	if ui.isInventoryVisible {
		ui.updateInventoryPanel()
	}
}

// updateHealthBar обновляет индикатор здоровья
func (ui *UIManager) updateHealthBar() {
	if ui.player != nil {
		ui.healthBar.SetValue(ui.player.Health / entity.MaxHealth)
	}
}

// updateSanityBar обновляет индикатор рассудка
func (ui *UIManager) updateSanityBar() {
	if ui.player != nil {
		ui.sanityBar.SetValue(ui.player.Sanity / entity.MaxSanity)
	}
}

// updateInventoryPanel обновляет панель инвентаря
func (ui *UIManager) updateInventoryPanel() {
	// Очищаем существующие элементы инвентаря
	// В реальной игре здесь был бы код для обновления отображения предметов
}

// ShowMessage показывает сообщение
func (ui *UIManager) ShowMessage(message string) {
	ui.messageLine.Text = message

	// В реальной игре здесь был бы код для автоматического скрытия сообщения через некоторое время
}

// SetCurrentScreen устанавливает текущий экран
func (ui *UIManager) SetCurrentScreen(screen string) {
	ui.currentScreen = screen

	// Скрываем все элементы UI
	for _, element := range ui.elements {
		element.SetVisible(false)
	}

	// Показываем только нужные элементы в зависимости от экрана
	switch screen {
	case "game":
		ui.healthBar.SetVisible(true)
		ui.sanityBar.SetVisible(true)
		ui.messageLine.SetVisible(true)

	case "main_menu":
		// Здесь будет код для отображения главного меню

	case "game_over":
		// Здесь будет код для отображения экрана окончания игры
	}
}

// IsGamePaused возвращает, приостановлена ли игра
func (ui *UIManager) IsGamePaused() bool {
	return ui.isMenuVisible || ui.isInventoryVisible
}
