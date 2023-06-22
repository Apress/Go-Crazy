package main

import (
	"github.com/gen2brain/raylib-go/raylib"
	"time"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

var (
	startTime time.Time
	stopTime  time.Time
)

func main() {

	// Initialize the window and the timer
	rl.InitWindow(screenWidth, screenHeight, "Start-Stop Timer")

	// Set the start and stop buttons
	startButton := rl.NewRectangle(screenWidth/2-50, screenHeight/2-25, 100, 50)
	stopButton := rl.NewRectangle(screenWidth/2-50, screenHeight/2+25, 100, 50)

	// Set the initial state
	isRunning := false
	var elapsedTime time.Duration

	for !rl.WindowShouldClose() {
		// Update the elapsed time if the timer is running
		if isRunning {
			//elapsedTime = time.Since(startTime)
			elapsedTime = time.Since(startTime).Round(time.Second)
		}

		// Check for button clicks
		if rl.CheckCollisionPointRec(rl.GetMousePosition(), startButton) && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			isRunning = true
			startTime = time.Now()
		}
		if rl.CheckCollisionPointRec(rl.GetMousePosition(), stopButton) && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			isRunning = false
			elapsedTime = 0
			stopTime = time.Time{}
		}

		// Draw the buttons and the timer
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		rl.DrawRectangleRec(startButton, rl.Green)
		rl.DrawRectangleRec(stopButton, rl.Red)
		rl.DrawText("START", int32(startButton.X)+20, int32(startButton.Y)+15, 20, rl.Black)
		rl.DrawText("STOP", int32(stopButton.X)+25, int32(stopButton.Y)+15, 20, rl.Black)
		rl.DrawText(elapsedTime.String(), screenWidth/2-50, screenHeight/2-100, 40, rl.Black)
		rl.EndDrawing()
	}

	// Close the window
	rl.CloseWindow()
}
