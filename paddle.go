package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	paddleMargin float64 = 100
	paddleWidth          = screenWidth / 50
	paddleHeight         = screenHeight / 4
	paddleSpeed          = 15
)

type Paddle struct {
	Rect
	Score int
}

func (p *Paddle) Draw(screen *ebiten.Image, color color.Color) {
	ebitenutil.DrawRect(screen, p.X, p.Y, p.W, p.H, color)
}

func (p *Paddle) Move(dy float64) {
	if inScreenBounds(p.X, p.Y+dy, p.W, p.H) {
		p.Y += dy
		return
	}

	if dy > 0 {
		p.Y = screenHeight - paddleHeight
	} else {
		p.Y = 0
	}
}

func NewPaddle(x float64) *Paddle {
	return &Paddle{Rect{x, screenHeight/2 - paddleHeight/2, paddleWidth, paddleHeight}, 0}
}

func movePaddle(key ebiten.Key, p *Paddle, dy float64) {
	if ebiten.IsKeyPressed(key) {
		p.Move(dy)
	}
}
