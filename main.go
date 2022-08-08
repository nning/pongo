package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 1920
	screenHeight = 1080
)

type rect struct {
	x, y, w, h float64
}

type vector struct {
	vx, vy float64
}

func main() {
	ebiten.SetWindowTitle("pongo")

	c := NewConfig().Load()
	ebiten.SetWindowSize(c.Width, c.Height)
	ebiten.SetFullscreen(c.Fullscreen)

	g := NewGame(c)

	go g.net.Announce()
	g.net.Listen()
	go g.net.SendState(g.ball, g.paddle1, g.paddle2)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
	// select {}
}
