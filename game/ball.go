package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	ballSize = screenWidth / 50
)

type Ball struct {
	Rect
	Vector
	Speed        float64
	Acceleration float64
}

func (b *Ball) Draw(screen *ebiten.Image) {
	ebitenutil.DrawRect(screen, b.X, b.Y, ballSize, ballSize, color.White)
}

func (b *Ball) Move(g *Game) int {
	dx := b.VX * b.Speed
	dy := b.VY * b.Speed

	screenCollision := checkBounds(b.X+dx, b.Y+dy, ballSize, ballSize)

	switch screenCollision {
	case 0:
		b.X += dx
		b.Y += dy
	case 1: // top
		b.VY = -b.VY
	case 2: // right
		b.VX = -b.VX
	case 3: // bottom
		b.VY = -b.VY
	case 4: // left
		b.VX = -b.VX
	}

	if rectCollision(&b.Rect, &g.Paddle1.Rect) || rectCollision(&b.Rect, &g.Paddle2.Rect) {
		b.VX = -b.VX
		b.Speed *= b.Acceleration
	}

	return screenCollision
}

func NewBall(speed, acceleration float64) *Ball {
	return &Ball{
		Rect:         Rect{screenWidth/2 - ballSize/2, screenHeight/2 - ballSize/2, ballSize, ballSize},
		Vector:       Vector{1, 1},
		Speed:        speed,
		Acceleration: acceleration,
	}
}
