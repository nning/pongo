package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	paddleMargin float64 = 150
	paddleWidth          = screenWidth / 50
	paddleHeight         = screenHeight / 4
	paddleSpeed          = 15
)

type Paddle struct {
	x, y, width, height float64
}

func (p *Paddle) Draw(screen *ebiten.Image) {
	ebitenutil.DrawRect(screen, p.x, p.y, p.width, p.height, color.White)
}

func (p *Paddle) Move(dy float64) {
	if inBounds(p.y + dy) {
		p.y += dy
		return
	}

	if dy > 0 {
		p.y = screenHeight - paddleHeight
	} else {
		p.y = 0
	}
}

func NewPaddle(x float64) *Paddle {
	return &Paddle{x, screenHeight/2 - paddleHeight/2, paddleWidth, paddleHeight}
}

func movePaddle(key ebiten.Key, p *Paddle, dy float64) {
	if ebiten.IsKeyPressed(key) {
		p.Move(dy)
	}
}
