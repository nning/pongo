package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 1920
	screenHeight = 1080
)

var paddle1 = NewPaddle(paddleMargin)
var paddle2 = NewPaddle(screenWidth - paddleMargin - paddleWidth)
var ball = NewBall()

func main() {
	ebiten.SetWindowTitle("pongo")

	c := NewConfig().Load()
	ebiten.SetWindowSize(c.Width, c.Height)
	ebiten.SetFullscreen(c.Fullscreen)

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
