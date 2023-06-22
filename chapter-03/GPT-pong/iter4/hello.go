package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"math/rand"
)

const (
	screenWidth  = 640
	screenHeight = 480
	paddleWidth  = 20
	paddleHeight = 80
	ballRadius   = 10
	ballVelocity = 5
	maxBalls     = 3 // maximum number of balls allowed on the screen at once
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
	defer rl.CloseWindow()

	leftPaddle := Paddle{X: 0, Y: float32(screenHeight/2 - paddleHeight/2)}
	rightPaddle := Paddle{X: float32(screenWidth - paddleWidth), Y: float32(screenHeight/2 - paddleHeight/2)}

	balls := make([]Ball, 0, maxBalls) // create an empty slice of balls

	for len(balls) < maxBalls {
		ball := Ball{X: float32(rand.Intn(screenWidth-2*ballRadius) + ballRadius), Y: float32(rand.Intn(screenHeight-2*ballRadius) + ballRadius)}
		ball.Velocity = rl.NewVector2(float32(rand.Intn(ballVelocity*2)-ballVelocity), float32(rand.Intn(ballVelocity*2)-ballVelocity))
		balls = append(balls, ball) // add the ball to the slice
	}

	var leftScore, rightScore int // keep track of the scores for each player

	var gameOver bool // flag to indicate if the game is over

	for !rl.WindowShouldClose() {
		if !gameOver {
			// Move the paddles
			if rl.IsKeyDown(rl.KeyW) && leftPaddle.Y > 0 {
				leftPaddle.Y -= 5
			}
			if rl.IsKeyDown(rl.KeyS) && leftPaddle.Y < screenHeight-paddleHeight {
				leftPaddle.Y += 5
			}
			if rl.IsKeyDown(rl.KeyUp) && rightPaddle.Y > 0 {
				rightPaddle.Y -= 5
			}
			if rl.IsKeyDown(rl.KeyDown) && rightPaddle.Y < screenHeight-paddleHeight {
				rightPaddle.Y += 5
			}

			// Update the balls
			for i := 0; i < len(balls); i++ {
				balls[i].X += balls[i].Velocity.X
				balls[i].Y += balls[i].Velocity.Y

				// Check for collisions with the walls and paddles
				if balls[i].X < ballRadius+paddleWidth && balls[i].Y > leftPaddle.Y && balls[i].Y < leftPaddle.Y+paddleHeight {
					balls[i].Velocity.X = -balls[i].Velocity.X
				}
				if balls[i].X > screenWidth-ballRadius-paddleWidth && balls[i].Y > rightPaddle.Y && balls[i].Y < rightPaddle.Y+paddleHeight {
					balls[i].Velocity.X = -balls[i].Velocity.X
				}
				if balls[i].X < ballRadius {
					rightScore++
					balls[i].Reset()
				}
				if balls[i].X > screenWidth-ballRadius {
					leftScore++
					balls[i].Reset()
				}
				if balls[i].Y < ballRadius || balls[i].Y > screenHeight-ballRadius {
					balls[i].Velocity.Y = -balls[i].Velocity.Y
				}
			} // Check if all the balls have been dropped
			if len(balls) > 0 {
				allBallsDropped := true
				for i := 0; i < len(balls); i++ {
					if balls[i].Y > ballRadius && balls[i].Y < screenHeight-ballRadius {
						allBallsDropped = false
						break
					}
				}
				if allBallsDropped {
					// Game over
					gameOver = true
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
		for i := 0; i < len(balls); i++ {
			rl.DrawCircle(int32(balls[i].X), int32(balls[i].Y), ballRadius, rl.White)
		}

		// Draw the scores
		rl.DrawText(fmt.Sprintf("%d", leftScore), screenWidth/4, 10, 30, rl.White)
		rl.DrawText(fmt.Sprintf("%d", rightScore), screenWidth*3/4, 10, 30, rl.White)

		// Draw the game over screen
		if gameOver {
			rl.DrawText("GAME OVER", screenWidth/2-100, screenHeight/2-30, 30, rl.White)
			rl.DrawText(fmt.Sprintf("Left: %d    Right: %d", leftScore, rightScore), screenWidth/2-140, screenHeight/2, 20, rl.White)
			rl.DrawText("Press SPACE to restart", screenWidth/2-140, screenHeight/2+30, 20, rl.White)

			if rl.IsKeyPressed(rl.KeySpace) {
				// Reset the game
				leftScore = 0
				rightScore = 0
				gameOver = false
				for i := 0; i < len(balls); i++ {
					balls[i].Reset()
				}
			}
		}

		rl.EndDrawing()
	}

}

// Reset the ball to its initial position and velocity
func (b Ball) Reset() {
	b.X = float32(rand.Intn(screenWidth-2*ballRadius) + ballRadius)
	b.Y = float32(rand.Intn(screenHeight-2*ballRadius) + ballRadius)
	b.Velocity = rl.NewVector2(float32(rand.Intn(ballVelocity*2)-ballVelocity), float32(rand.Intn(ballVelocity*2)-ballVelocity))
}
