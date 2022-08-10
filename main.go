package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 1920
	screenHeight = 1080
)

var game *Game

type Rect struct {
	X, Y, W, H float64
}

type Vector struct {
	VX, VY float64
}

func main() {
	ebiten.SetWindowTitle("pongo")

	c := NewConfig().Load()
	ebiten.SetWindowSize(c.Width, c.Height)
	ebiten.SetFullscreen(c.Fullscreen)

	game = NewGame(c)

	if !c.Offline {
		go game.net.Announce()
		game.net.Listen()
		go game.net.SendState(game)
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
	// select {}
}
