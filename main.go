package main

import (
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth          = 1920
	screenHeight         = 1080
	margin       float64 = 150
	paddleWidth          = screenWidth / 50
	paddleHeight         = screenHeight / 4
	speed                = 15
)

type Paddle struct {
	x      float64
	y      float64
	width  float64
	height float64
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

type Game struct {
}

func inBounds(y float64) bool {
	return y >= 0 && y <= screenHeight-paddleHeight
}

func movePaddle(key ebiten.Key, p *Paddle, dy float64) {
	if ebiten.IsKeyPressed(key) {
		p.Move(dy)
	}
}

func (g *Game) Update() error {
	movePaddle(ebiten.KeyW, p1, -speed)
	movePaddle(ebiten.KeyS, p1, speed)

	movePaddle(ebiten.KeyUp, p2, -speed)
	movePaddle(ebiten.KeyDown, p2, speed)

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, p := range []*Paddle{p1, p2} {
		p.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

var p1 = NewPaddle(margin)
var p2 = NewPaddle(screenWidth - margin - paddleWidth)

func main() {
	ebiten.SetWindowSize(1920, 1080)
	ebiten.SetWindowTitle("pongo")

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
