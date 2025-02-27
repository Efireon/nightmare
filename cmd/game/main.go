package main

import (
	"log"
	"math/rand"
	"time"

	"nightmare/internal/core"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// Инициализация генератора случайных чисел
	rand.Seed(time.Now().UnixNano())

	// Создание игры
	game, err := core.NewGame()
	if err != nil {
		log.Fatalf("Не удалось создать игру: %v", err)
	}

	// Настройка окна
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Nightmare Forest")

	// Запуск игрового цикла
	if err := ebiten.RunGame(game); err != nil {
		log.Fatalf("Игра завершилась с ошибкой: %v", err)
	}
}
