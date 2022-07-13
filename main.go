package main

import (
	"log"

	ebiten "github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 320
	screenHeight = 240
)

type Game struct {
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(1920, 1080)
	ebiten.SetWindowTitle("pongo")

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
