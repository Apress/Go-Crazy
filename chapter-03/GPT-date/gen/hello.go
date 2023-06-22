package main

import (
	"fmt"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := 800
	screenHeight := 450

	rl.InitWindow(screenWidth, screenHeight, "Real-Time Date Display")

	rl.SetTargetFPS(60)

	font := rl.LoadFont("arial.ttf")
	if font.Size == 0 {
		fmt.Println("Failed to load font")
		return
	}

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		// Get current date and time
		now := time.Now()
		dateStr := now.Format("January 02, 2006 15:04:05")

		// Draw date and time
		rl.DrawText(dateStr, screenWidth/2-rl.MeasureText(dateStr, 20)/2, screenHeight/2-20, 20, rl.Maroon)

		rl.EndDrawing()

		time.Sleep(time.Second)
	}

	rl.UnloadFont(font)
	rl.CloseWindow()
}
