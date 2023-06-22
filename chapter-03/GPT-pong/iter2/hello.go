package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth   = 800
	screenHeight  = 450
	ballRadius    = 10
	paddleWidth   = 10
	paddleHeight  = 80
	ballVelocityX = 0.1
	ballVelocityY = 0.1
	paddleSpeed   = 2
)

func main() {
	rl.InitWindow(screenWidth, screenHeight, "Pong")

	ballX := float32(screenWidth) / 2
	ballY := float32(screenHeight) / 2
	ballVelocity := rl.NewVector2(ballVelocityX, ballVelocityY)

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

			ballX += ballVelocity.X
			ballY += ballVelocity.Y

			if ballX < ballRadius || ballX > float32(screenWidth)-ballRadius {
				ballVelocity.X = -ballVelocity.X
			}
			if ballY < ballRadius || ballY > float32(screenHeight)-ballRadius {
				ballVelocity.Y = -ballVelocity.Y
			}

			if ballX < float32(paddleWidth)+ballRadius {
				if ballY > leftPaddleY && ballY < leftPaddleY+float32(paddleHeight) {
					ballVelocity.X = -ballVelocity.X
				} else {
					gameOver = true
				}
			}
			if ballX > float32(screenWidth)-float32(paddleWidth)-ballRadius {
				if ballY > rightPaddleY && ballY < rightPaddleY+float32(paddleHeight) {
					ballVelocity.X = -ballVelocity.X
				} else {
					gameOver = true
				}
			}
		} else {
			if rl.IsKeyPressed(rl.KeyEnter) {
				ballX = float32(screenWidth) / 2
				ballY = float32(screenHeight) / 2
				leftPaddleY = float32(screenHeight)/2 - float32(paddleHeight)/2
				rightPaddleY = float32(screenHeight)/2 - float32(paddleHeight)/2
				gameOver = false
			}
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)

		if !gameOver {
			rl.DrawCircle(int32(ballX), int32(ballY), ballRadius, rl.White)
			rl.DrawRectangle(0, int32(leftPaddleY), paddleWidth, paddleHeight, rl.White)
			rl.DrawRectangle(screenWidth-paddleWidth, int32(rightPaddleY), paddleWidth, paddleHeight, rl.White)
		} else {
			rl.DrawText("Game Over", screenWidth/2-80, screenHeight/2-10, 20, rl.White)
			rl.DrawText("Press Enter to restart", screenWidth/
				2-80, screenHeight/2+10, 20, rl.White)
		}

		rl.EndDrawing()
	}

	rl.CloseWindow()

}
