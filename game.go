package main

import (
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
}

func inScreenBounds(x, y, w, h float64) bool {
	return x >= 0 && x+w <= screenWidth && y >= 0 && y+h <= screenHeight
}

func (g *Game) Update() error {
	movePaddle(ebiten.KeyW, paddle1, -paddleSpeed)
	movePaddle(ebiten.KeyS, paddle1, paddleSpeed)

	movePaddle(ebiten.KeyUp, paddle2, -paddleSpeed)
	movePaddle(ebiten.KeyDown, paddle2, paddleSpeed)

	ball.Move()

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, p := range []*Paddle{paddle1, paddle2} {
		p.Draw(screen)
	}

	ball.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
