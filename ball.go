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
	rect
	vector
	speed        float64
	acceleration float64
}

func (b *Ball) Draw(screen *ebiten.Image) {
	ebitenutil.DrawRect(screen, b.x, b.y, ballSize, ballSize, color.White)
}

func (b *Ball) Move(g *Game) {
	dx := b.vx * b.speed
	dy := b.vy * b.speed

	screenCollision := checkBounds(b.x+dx, b.y+dy, ballSize, ballSize)

	switch screenCollision {
	case 0:
		b.x += dx
		b.y += dy
	case 1: // top
		b.vy = -b.vy
	case 2: // right
		b.vx = -b.vx
	case 3: // bottom
		b.vy = -b.vy
	case 4: // left
		b.vx = -b.vx
	}

	if rectCollision(&b.rect, &g.paddle1.rect) || rectCollision(&b.rect, &g.paddle2.rect) {
		b.vx = -b.vx
		b.speed *= b.acceleration
	}
}

func NewBall(speed, acceleration float64) *Ball {
	return &Ball{
		rect:         rect{screenWidth/2 - ballSize/2, screenHeight/2 - ballSize/2, ballSize, ballSize},
		vector:       vector{1, 1},
		speed:        speed,
		acceleration: acceleration,
	}
}
