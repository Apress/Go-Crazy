package main

import (
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 800
	screenHeight = 450
	numDominos   = 20
	dominoWidth  = 40
	dominoHeight = 80
)

type Domino struct {
	pos     rl.Vector2
	vel     rl.Vector2
	falling bool
}

func main() {
	rl.InitWindow(screenWidth, screenHeight, "Dominoes")

	rand.Seed(time.Now().UnixNano())

	dominos := make([]Domino, numDominos)
	for i := 0; i < numDominos; i++ {
		x := float32(rand.Intn(screenWidth - dominoWidth))
		y := float32(rand.Intn(screenHeight/2 - dominoHeight))
		dominos[i] = Domino{
			pos: rl.NewVector2(x, y),
			vel: rl.NewVector2(0, 0),
		}
	}

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		for i := 0; i < numDominos; i++ {
			domino := &dominos[i]

			if !domino.falling {
				rl.DrawRectangle(int32(domino.pos.X), int32(domino.pos.Y), dominoWidth, dominoHeight, rl.Black)
			} else {
				rl.DrawRectangle(int32(domino.pos.X), int32(domino.pos.Y), dominoHeight, dominoWidth, rl.Black)
			}

			if !domino.falling {
				if i == 0 || dominos[i-1].falling {
					domino.vel = rl.NewVector2(0, 5)
					domino.falling = true
				}
			} else {
				if domino.pos.Y > screenHeight-dominoHeight {
					domino.vel = rl.NewVector2(0, 0)
				} else {
					domino.vel.Y += 0.2
				}
			}

			domino.pos = rl.Vector2Add(domino.pos, domino.vel)
		}

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
