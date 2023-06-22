package main

import (
	"fmt"
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 800
	screenHeight = 450
	paddleWidth  = 20
	paddleHeight = 100
	ballRadius   = 10
	ballVelocity = 1
)

type Paddle struct {
	X, Y float32
}

type Ball struct {
	X, Y     float32
	Velocity rl.Vector2
}

func main() {
	rl.InitWindow(screenWidth, screenHeight, "Pong")

	// Set the random seed
	rand.Seed(time.Now().UnixNano())

	// Create the left and right paddles
	leftPaddle := Paddle{X: 50, Y: float32(screenHeight/2 - paddleHeight/2)}
	rightPaddle := Paddle{X: float32(screenWidth - 50 - paddleWidth), Y: float32(screenHeight/2 - paddleHeight/2)}

	// Create the balls
	var balls []Ball
	for i := 0; i < 3; i++ {
		balls = append(balls, Ball{X: float32(rand.Intn(screenWidth-2*ballRadius) + ballRadius), Y: float32(rand.Intn(screenHeight-2*ballRadius) + ballRadius), Velocity: rl.NewVector2(float32(rand.Intn(ballVelocity*2)-ballVelocity), float32(rand.Intn(ballVelocity*2)-ballVelocity))})
	}

	// Set the initial score
	leftScore := 0
	rightScore := 0

	// Set the game over flag
	gameOver := false

	for !rl.WindowShouldClose() {
		// Update the game
		if !gameOver {
			// Move the left paddle
			if rl.IsKeyDown(rl.KeyW) && leftPaddle.Y > 0 {
				leftPaddle.Y -= 2
			}
			if rl.IsKeyDown(rl.KeyS) && leftPaddle.Y < screenHeight-paddleHeight {
				leftPaddle.Y += 2
			}

			// Move the right paddle
			if rl.IsKeyDown(rl.KeyUp) && rightPaddle.Y > 0 {
				rightPaddle.Y -= 2
			}
			if rl.IsKeyDown(rl.KeyDown) && rightPaddle.Y < screenHeight-paddleHeight {
				rightPaddle.Y += 2
			}

			// Move the balls
			for i := 0; i < len(balls); i++ {
				balls[i].X += balls[i].Velocity.X
				balls[i].Y += balls[i].Velocity.Y

				// Check for collision with the left paddle
				if balls[i].X-ballRadius < leftPaddle.X+paddleWidth && balls[i].Y > leftPaddle.Y && balls[i].Y < leftPaddle.Y+paddleHeight {
					balls[i].Velocity.X = -balls[i].Velocity.X
					balls[i].Velocity.Y += rl.NewVector2(0, (balls[i].Y-(leftPaddle.Y+paddleHeight/2))/paddleHeight*5).Y
				}

				// Check for collision with the right paddle
				if balls[i].X+ballRadius > rightPaddle.X && balls[i].Y > rightPaddle.Y && balls[i].Y < rightPaddle.Y+paddleHeight {
					balls[i].Velocity.X = -balls[i].Velocity.X
					balls[i].Velocity.Y += rl.NewVector2(0, (balls[i].Y-(rightPaddle.Y+paddleHeight/2))/paddleHeight*5).Y
				}

				// Check for collision with the top and bottom walls
				if balls[i].Y-ballRadius < 0 || balls[i].Y+ballRadius > screenHeight {
					balls[i].Velocity.Y = -balls[i].Velocity.Y
				}

				// Check if a ball is dropped
				if balls[i].X+ballRadius < 0 || balls[i].X-ballRadius > screenWidth {
					if balls[i].X+ballRadius < 0 {
						rightScore++
					} else {
						leftScore++
					}

					// Remove the dropped ball and check if there are any balls remaining
					balls = append(balls[:i], balls[i+1:]...)
					i--
					if len(balls) == 0 {
						gameOver = true
					}
				}
			}
		}

		// Draw the game
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		// Draw the paddles
		rl.DrawRectangle(int32(leftPaddle.X), int32(leftPaddle.Y), paddleWidth, paddleHeight, rl.White)
		rl.DrawRectangle(int32(rightPaddle.X), int32(rightPaddle.Y), paddleWidth, paddleHeight, rl.White)

		// Draw the balls
		for _, ball := range balls {
			rl.DrawCircle(int32(ball.X), int32(ball.Y), ballRadius, rl.White)
		}

		// Draw the score
		rl.DrawText(fmt.Sprintf("%d", leftScore), 100, 10, 50, rl.White)
		rl.DrawText(fmt.Sprintf("%d", rightScore), screenWidth-100, 10, 50, rl.White)

		// Draw the game over screen
		if gameOver {
			rl.DrawText("Game Over", screenWidth/2-100, screenHeight/2-50, 50, rl.White)
			rl.DrawText(fmt.Sprintf("Left: %d - Right: %d", leftScore, rightScore), screenWidth/2-150, screenHeight/2, 30, rl.White)
			rl.DrawText("Press R to restart", screenWidth/2-125, screenHeight/2+50, 30, rl.White)
			if rl.IsKeyPressed(rl.KeyR) {
				// Reset the game
				leftPaddle.Y = float32(screenHeight/2 - paddleHeight/2)
				rightPaddle.Y = float32(screenHeight/2 - paddleHeight/2)
				balls = nil
				for i := 0; i < 3; i++ {
					balls = append(balls, Ball{X: float32(rand.Intn(screenWidth-2*ballRadius) + ballRadius), Y: float32(rand.Intn(screenHeight-2*ballRadius) + ballRadius), Velocity: rl.NewVector2(float32(rand.Intn(ballVelocity*2)-ballVelocity), float32(rand.Intn(ballVelocity*2)-ballVelocity))})
				}
				leftScore = 0
				rightScore = 0
				gameOver = false
			}
		}

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
