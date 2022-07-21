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
	rect
	score int
}

func (p *Paddle) Draw(screen *ebiten.Image) {
	ebitenutil.DrawRect(screen, p.x, p.y, p.w, p.h, color.White)
}

func (p *Paddle) Move(dy float64) {
	if inScreenBounds(p.x, p.y+dy, p.w, p.h) {
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
	return &Paddle{rect{x, screenHeight/2 - paddleHeight/2, paddleWidth, paddleHeight}, 0}
}

func movePaddle(key ebiten.Key, p *Paddle, dy float64) {
	if ebiten.IsKeyPressed(key) {
		p.Move(dy)
	}
}
