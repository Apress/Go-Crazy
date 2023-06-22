package main

import (
	"math/rand"
	"time"

	"github.com/gen2brain/raylib-go/physics"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth   = 800
	screenHeight  = 450
	numDominos    = 20
	dominoWidth   = 40
	dominoHeight  = 80
	floorHeight   = 20
	floorFriction = 1
)

type Domino struct {
	body    *physics.Body
	falling bool
}

func main() {
	rl.InitWindow(screenWidth, screenHeight, "Dominoes")

	physics.Init()
	defer physics.Close()

	rand.Seed(time.Now().UnixNano())

	// Create a static rectangle body for the floor
	floor := physics.NewBodyRectangle(rl.NewVector2(screenWidth/2, screenHeight-floorHeight/2), screenWidth, floorHeight, 1)
	floor.IsGrounded = true

	//physics.SetPhysicsProperties(0, 0, 0, floorFriction)
	//physics.Set

	dominos := make([]Domino, numDominos)
	for i := 0; i < numDominos; i++ {
		x := float32(rand.Intn(screenWidth - dominoWidth))
		y := float32(rand.Intn(screenHeight/2 - dominoHeight))
		body := physics.NewBodyRectangle(rl.NewVector2(x+dominoWidth/2, y+dominoHeight/2), dominoWidth, dominoHeight, 1)
		body.Enabled = false
		dominos[i] = Domino{
			body: body,
		}
	}

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		physics.Update()

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		// Draw the floor
		rl.DrawRectangle(0, screenHeight-floorHeight, screenWidth, floorHeight, rl.Gray)

		for i := 0; i < numDominos; i++ {
			domino := &dominos[i]

			if !domino.falling {
				rl.DrawRectangle(int32(domino.body.Position.X-dominoWidth/2), int32(domino.body.Position.Y-dominoHeight/2), dominoWidth, dominoHeight, rl.Black)
			} else {
				rl.DrawRectangle(int32(domino.body.Position.X-dominoHeight/2), int32(domino.body.Position.Y-dominoWidth/2), dominoHeight, dominoWidth, rl.Black)
			}

			if !domino.falling {
				if i == 0 || dominos[i-1].falling {
					domino.body.Enabled = true
					domino.body.Velocity = rl.NewVector2(0, 5)
					domino.falling = true
				}
			}
		}

		rl.EndDrawing()
	}

	physics.Close()
	rl.CloseWindow()
}