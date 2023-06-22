package main

import (
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 800
	screenHeight = 450
	ballRadius   = 10
	paddleWidth  = 10
	paddleHeight = 80
	ballVelocity = 1
	maxBalls     = 5
	paddleSpeed  = 2
)

func main() {
	rl.InitWindow(screenWidth, screenHeight, "Pong")

	balls := make([]Ball, maxBalls)
	for i := range balls {
		balls[i] = NewBall()
	}

	leftPaddleY := float32(screenHeight)/2 - float32(paddleHeight)/2
	rightPaddleY := float32(screenHeight)/2 - float32(paddleHeight)/2

	gameOver := false

	for !rl.WindowShouldClose() {
		if !gameOver {
			if rl.IsKeyDown(rl.KeyW) {
				leftPaddleY -= paddleSpeed
			}
			if rl.IsKeyDown(rl.KeyS) {
				leftPaddleY += paddleSpeed
			}
			if rl.IsKeyDown(rl.KeyUp) {
				rightPaddleY -= paddleSpeed
			}
			if rl.IsKeyDown(rl.KeyDown) {
				rightPaddleY += paddleSpeed
			}

			for i := range balls {
				balls[i].Update(leftPaddleY, rightPaddleY)
			}

			if len(balls) == 0 {
				gameOver = true
			}
		} else {
			if rl.IsKeyPressed(rl.KeyEnter) {
				balls = make([]Ball, maxBalls)
				for i := range balls {
					balls[i] = NewBall()
				}
				leftPaddleY = float32(screenHeight)/2 - float32(paddleHeight)/2
				rightPaddleY = float32(screenHeight)/2 - float32(paddleHeight)/2
				gameOver = false
			}
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)

		if !gameOver {
			for i := range balls {
				balls[i].Draw()
			}
			rl.DrawRectangle(0, int32(leftPaddleY), paddleWidth, paddleHeight, rl.White)
			rl.DrawRectangle(screenWidth-paddleWidth, int32(rightPaddleY), paddleWidth, paddleHeight, rl.White)
		} else {
			rl.DrawText("Game Over", screenWidth/2-80, screenHeight/2-10, 20, rl.White)
			rl.DrawText("Press Enter to restart", screenWidth/2-80, screenHeight/2+10, 20, rl.White)
		}

		rl.EndDrawing()
	}

	rl.CloseWindow()
}

type Ball struct {
	X, Y     float32
	Velocity rl.Vector2
}

func NewBall() Ball {
	return Ball{
		X:        float32(rand.Intn(screenWidth-2*ballRadius) + ballRadius),
		Y:        float32(rand.Intn(screenHeight-2*ballRadius) + ballRadius),
		Velocity: rl.NewVector2(float32(rand.Intn(ballVelocity*2)-ballVelocity), float32(rand.Intn(ballVelocity*2)-ballVelocity)),
	}
}

func (b *Ball) Update(leftPaddleY, rightPaddleY float32) {
	b.X += b.Velocity.X
	b.Y += b.Velocity.Y

	if b.X < ballRadius || b.X > screenWidth-ballRadius {
		b.Velocity.X = -b.Velocity.X
	}

	if b.Y < ballRadius || b.Y > screenHeight-ballRadius {
		b.Velocity.Y = -b.Velocity.Y
	}

	if b.X < paddleWidth+ballRadius && b.Y > leftPaddleY && b.Y < leftPaddleY+paddleHeight {
		b.Velocity.X = -b.Velocity.X
	}

	if b.X > screenWidth-paddleWidth-ballRadius && b.Y > rightPaddleY && b.Y < rightPaddleY+paddleHeight {
		b.Velocity.X = -b.Velocity.X
	}

	if b.X < -ballRadius || b.X > screenWidth+ballRadius {
		b.Reset()
	}

}

func (b *Ball) Draw() {
	rl.DrawCircle(int32(b.X), int32(b.Y), ballRadius, rl.White)
}

func (b Ball) Reset() {
	b.X = float32(rand.Intn(screenWidth-ballRadius) + ballRadius)
	b.Y = float32(rand.Intn(screenHeight-ballRadius) + ballRadius)
	b.Velocity = rl.NewVector2(float32(rand.Intn(ballVelocity)-ballVelocity), float32(rand.Intn(ballVelocity*2)-ballVelocity))
}
