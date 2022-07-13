package main

import (
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
}

func inBounds(y float64) bool {
	return y >= 0 && y <= screenHeight-paddleHeight
}

func (g *Game) Update() error {
	movePaddle(ebiten.KeyW, paddle1, -paddleSpeed)
	movePaddle(ebiten.KeyS, paddle1, paddleSpeed)

	movePaddle(ebiten.KeyUp, paddle2, -paddleSpeed)
	movePaddle(ebiten.KeyDown, paddle2, paddleSpeed)

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}

	ball.Move()

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
