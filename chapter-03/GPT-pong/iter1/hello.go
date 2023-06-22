package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 800
	screenHeight = 450
	ballRadius   = 10
	paddleWidth  = 10
	paddleHeight = 80
	paddleSpeed  = 2
)

func main() {
	rl.InitWindow(screenWidth, screenHeight, "Pong")

	// Initialize ball position and velocity
	ballX := float32(screenWidth) / 2
	ballY := float32(screenHeight) / 2
	ballVelocityX := 0.1
	ballVelocityY := 0.1

	// Initialize paddle positions
	leftPaddleY := float32(screenHeight)/2 - float32(paddleHeight)/2
	rightPaddleY := float32(screenHeight)/2 - float32(paddleHeight)/2

	for !rl.WindowShouldClose() {
		// Move paddles based on input
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

		// Move ball
		ballX += float32(ballVelocityX)
		ballY += float32(ballVelocityY)

		// Bounce ball off walls
		if ballX < ballRadius || ballX > float32(screenWidth)-ballRadius {
			ballVelocityX = -ballVelocityX
		}
		if ballY < ballRadius || ballY > float32(screenHeight)-ballRadius {
			ballVelocityY = -ballVelocityY
		}

		// Bounce ball off paddles
		if ballX < float32(paddleWidth)+ballRadius {
			if ballY > leftPaddleY && ballY < leftPaddleY+float32(paddleHeight) {
				ballVelocityX = -ballVelocityX
			} else {
				// Player on the right wins
				break
			}
		}
		if ballX > float32(screenWidth)-float32(paddleWidth)-ballRadius {
			if ballY > rightPaddleY && ballY < rightPaddleY+float32(paddleHeight) {
				ballVelocityX = -ballVelocityX
			} else {
				// Player on the left wins
				break
			}
		}

		// Draw everything
		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)

		// Draw ball
		rl.DrawCircle(int32(ballX), int32(ballY), ballRadius, rl.White)

		// Draw paddles
		rl.DrawRectangle(0, int32(leftPaddleY), paddleWidth, paddleHeight, rl.White)
		rl.DrawRectangle(screenWidth-paddleWidth, int32(rightPaddleY), paddleWidth, paddleHeight, rl.White)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
