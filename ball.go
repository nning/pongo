package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	ballSize  = screenWidth / 50
	ballSpeed = 5
)

type vector struct {
	vx, vy float64
}

type Ball struct {
	x, y float64
	vector
}

func (b *Ball) Draw(screen *ebiten.Image) {
	ebitenutil.DrawRect(screen, b.x, b.y, ballSize, ballSize, color.White)
}

func (b *Ball) Move() {
	dx := b.vx * ballSpeed
	dy := b.vy * ballSpeed

	if inScreenBounds(b.x+dx, b.y+dy, ballSize, ballSize) {
		b.x += dx
		b.y += dy
	}

}

func NewBall() *Ball {
	return &Ball{screenWidth/2 - ballSize/2, screenHeight/2 - ballSize/2, vector{1, 1}}
}
