package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 800
	screenHeight = 450
	gravity      = 500
)

var (
	mushrooms []mushroom
)

type player struct {
	position rl.Vector2
	velocity rl.Vector2
	size     rl.Vector2
}

type mushroom struct {
	position rl.Vector2
	size     rl.Vector2
}

func main() {
	rl.InitWindow(screenWidth, screenHeight, "Simple Platformer")

	p := player{
		position: rl.NewVector2(50, 50),
		size:     rl.NewVector2(50, 50),
	}

	mushrooms = []mushroom{
		{position: rl.NewVector2(200, 350), size: rl.NewVector2(30, 30)},
		{position: rl.NewVector2(400, 300), size: rl.NewVector2(30, 30)},
		{position: rl.NewVector2(600, 250), size: rl.NewVector2(30, 30)},
	}

	for !rl.WindowShouldClose() {
		p.update()

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		rl.DrawRectangleV(p.position, p.size, rl.Red)
		for _, m := range mushrooms {
			rl.DrawRectangleV(m.position, m.size, rl.Green)
		}
		rl.EndDrawing()
	}

	rl.CloseWindow()
}

func (p *player) update() {
	// Apply gravity to player velocity
	p.velocity.Y += gravity * rl.GetFrameTime()

	// Move player horizontally
	if rl.IsKeyDown(rl.KeyRight) {
		p.velocity.X = 300 * rl.GetFrameTime()
	} else if rl.IsKeyDown(rl.KeyLeft) {
		p.velocity.X = -300 * rl.GetFrameTime()
	} else {
		p.velocity.X = 0
	}

	// Check for collisions with ground
	if p.position.Y+p.size.Y >= screenHeight {
		p.velocity.Y = 0
		p.position.Y = screenHeight - p.size.Y
	}

	// Check for collisions with mushrooms
	for _, m := range mushrooms {
		if rl.CheckCollisionRecs(rl.NewRectangle(p.position.X, p.position.Y, p.size.X, p.size.Y), rl.NewRectangle(m.position.X, m.position.Y, m.size.X, m.size.Y)) {
			// Player collides with mushroom
			p.velocity.Y = -300
		}
	}

	// Apply player velocity to position
	p.position.X += p.velocity.X
	p.position.Y += p.velocity.Y
}
