package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	ballSize = screenWidth / 50
)

type Ball struct {
	x, y float64
}

func (b *Ball) Draw(screen *ebiten.Image) {
	ebitenutil.DrawRect(screen, b.x, b.y, ballSize, ballSize, color.White)
}

func NewBall() *Ball {
	return &Ball{screenWidth/2 - ballSize/2, screenHeight/2 - ballSize/2}
}
