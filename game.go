package main

import (
	"fmt"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	config           *Config
	ball             *Ball
	paddle1, paddle2 *Paddle

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

func rectCollision(r1, r2 *rect) bool {
	return r1.x+r1.w >= r2.x && r1.x <= r2.x+r2.w && r1.y+r1.h >= r2.y && r1.y <= r2.y+r2.h
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
		leftStickY := ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickVertical)
		if leftStickY > 0.1 || leftStickY < -0.1 {
			g.paddle1.Move(leftStickY * paddleSpeed)
		}

		rightStickY := ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisRightStickVertical)
		if rightStickY > 0.1 || rightStickY < -0.1 {
			g.paddle2.Move(rightStickY * paddleSpeed)
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

	movePaddle(ebiten.KeyW, g.paddle1, -paddleSpeed)
	movePaddle(ebiten.KeyS, g.paddle1, paddleSpeed)

	movePaddle(ebiten.KeyUp, g.paddle2, -paddleSpeed)
	movePaddle(ebiten.KeyDown, g.paddle2, paddleSpeed)

	screenCollision := g.ball.Move(g)
	if screenCollision == 2 {
		g.paddle1.score++
	} else if screenCollision == 4 {
		g.paddle2.score++
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

	ebiten.SetWindowTitle(fmt.Sprintf("%d : %d", g.paddle1.score, g.paddle2.score))

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, p := range []*Paddle{g.paddle1, g.paddle2} {
		p.Draw(screen)
	}

	g.ball.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func NewGame(c *Config) *Game {
	return &Game{
		config:  c,
		ball:    NewBall(c.BallSpeed, c.BallAcceleration),
		paddle1: NewPaddle(paddleMargin),
		paddle2: NewPaddle(screenWidth - paddleMargin - paddleWidth),
	}
}
