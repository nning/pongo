package main

import (
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	log "github.com/sirupsen/logrus"
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

	if len(os.Args) > 1 && os.Args[1] == "-d" {
		c.Debug = true
	}

	if c.Debug {
		log.SetLevel(log.DebugLevel)
	}

	game = NewGame(c)

	if !c.Offline {
		go game.net.Announce()
		game.net.Listen()
		go game.net.SendState(game)
	}

	go game.RecordFPS()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
	// select {}
}
