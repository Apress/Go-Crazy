package main

import (
	"fmt"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 800
	screenHeight = 480
	fontSize     = 36
)

func main() {

	rl.InitWindow(screenWidth, screenHeight, "Real-Time Date Display")

	rl.SetTargetFPS(60)

	font := rl.LoadFont("font.ttf")
	if font.BaseSize == 0 {
		fmt.Println("Failed to load font")
		return
	}

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.LightGray)

		// Get current date and time
		now := time.Now()
		dateStr := now.Format("January 02, 2006 15:04:05")

		// Draw date and time
		position := rl.Vector2{X: float32(screenWidth/2 - rl.MeasureTextEx(font, dateStr, fontSize, 0).X/2), Y: screenHeight/2 - 20}
		rl.DrawTextEx(font, dateStr, position, fontSize, 0, rl.Black)

		rl.EndDrawing()

		time.Sleep(time.Second)
	}

	rl.UnloadFont(font)
	rl.CloseWindow()
}
