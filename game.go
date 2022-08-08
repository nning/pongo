package main

import (
	"fmt"
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	Solo GameMode = iota
	Left
	Right
)

type GameMode int

func (g GameMode) String() string {
	switch g {
	case Solo:
		return "Solo"
	case Left:
		return "Left"
	case Right:
		return "Right"
	}

	return "Unknown"
}

type Game struct {
	config *Config
	mode   GameMode
	net    *Net

	Ball             *Ball
	Paddle1, Paddle2 *Paddle

	gamepadIDsBuf []ebiten.GamepadID
	gamepadIDs    map[ebiten.GamepadID]struct{}
}

func inScreenBounds(x, y, w, h float64) bool {
	return checkBounds(x, y, w, h) == 0
}

func checkBounds(x, y, w, h float64) int {
	if x <= 0 {
		return 4 // left
	}

	if x+w >= screenWidth {
		return 2 // right
	}

	if y <= 0 {
		return 1 // top
	}

	if y+h >= screenHeight {
		return 3 // bottom
	}

	return 0
}

func rectCollision(r1, r2 *Rect) bool {
	return r1.X+r1.W >= r2.X && r1.X <= r2.X+r2.W && r1.Y+r1.H >= r2.Y && r1.Y <= r2.Y+r2.H
}

func (g *Game) handleGamepadConnections() {
	if g.gamepadIDs == nil {
		g.gamepadIDs = map[ebiten.GamepadID]struct{}{}
	}

	g.gamepadIDsBuf = inpututil.AppendJustConnectedGamepadIDs(g.gamepadIDsBuf[:0])
	for _, id := range g.gamepadIDsBuf {
		g.gamepadIDs[id] = struct{}{}
	}

	for id := range g.gamepadIDs {
		if inpututil.IsGamepadJustDisconnected(id) {
			delete(g.gamepadIDs, id)
		}
	}
}

func (g *Game) handleGamepadInput() {
	for id := range g.gamepadIDs {
		if g.mode == Solo || g.mode == Left {
			leftStickY := ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickVertical)
			if leftStickY > 0.1 || leftStickY < -0.1 {
				g.Paddle1.Move(leftStickY * paddleSpeed)
			}
		}

		if g.mode == Solo || g.mode == Right {
			rightStickY := ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisRightStickVertical)
			if rightStickY > 0.1 || rightStickY < -0.1 {
				g.Paddle2.Move(rightStickY * paddleSpeed)
			}
		}

		if inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonCenterLeft) {
			os.Exit(0)
		}

		if inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonCenterRight) {
			os.Exit(0)
		}
	}
}

func (g *Game) Update() error {
	g.handleGamepadConnections()
	g.handleGamepadInput()

	if g.mode == Solo || g.mode == Left {
		movePaddle(ebiten.KeyW, g.Paddle1, -paddleSpeed)
		movePaddle(ebiten.KeyS, g.Paddle1, paddleSpeed)
	}

	if g.mode == Solo || g.mode == Right {
		movePaddle(ebiten.KeyUp, g.Paddle2, -paddleSpeed)
		movePaddle(ebiten.KeyDown, g.Paddle2, paddleSpeed)
	}

	if g.mode == Solo || g.mode == Left {
		screenCollision := g.Ball.Move(g)
		if screenCollision == 2 {
			g.Paddle1.Score++
		} else if screenCollision == 4 {
			g.Paddle2.Score++
		}
	}

	if g.net.LastState != nil {
		if g.mode == Left {
			g.Paddle2.Y = g.net.LastState.Paddle2.Y
		} else {
			g.Ball = g.net.LastState.Ball
			g.Paddle1 = g.net.LastState.Paddle1
			g.Paddle2.Score = g.net.LastState.Paddle2.Score
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}

	if g.config.Debug && ebiten.IsKeyPressed(ebiten.KeyQ) {
		os.Exit(0)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	ebiten.SetWindowTitle(fmt.Sprintf("%d : %d", g.Paddle1.Score, g.Paddle2.Score))

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	var c1 color.Color = color.White
	var c2 color.Color = color.White

	if g.mode == Left {
		c1 = color.RGBA{0, 255, 0, 255}
	} else if g.mode == Right {
		c2 = color.RGBA{0, 255, 0, 255}
	}

	g.Paddle1.Draw(screen, c1)
	g.Paddle2.Draw(screen, c2)

	g.Ball.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func NewGame(c *Config) *Game {
	return &Game{
		config:  c,
		Ball:    NewBall(c.BallSpeed, c.BallAcceleration),
		Paddle1: NewPaddle(paddleMargin),
		Paddle2: NewPaddle(screenWidth - paddleMargin - paddleWidth),
		mode:    Solo,
		net:     NewNet(c),
	}
}
