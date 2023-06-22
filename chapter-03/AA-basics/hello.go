package main

// 1. Library import
import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	// rl is the alias for raylib.
	// here we create a new window
	rl.InitWindow(800, 450, "game window")

	// this is the target frame rate of your game
	rl.SetTargetFPS(60)

	// until the window is being asked to close
	// run the game
	for !rl.WindowShouldClose() {
		// being drawing
		rl.BeginDrawing()

		// prepare the color for the background
		bgColor := rl.RayWhite
		// set the background color
		rl.ClearBackground(bgColor)

		// draw some text on top
		rl.DrawText("Your game!", 190, 200, 20, rl.SkyBlue)

		// finish drawing on the canvas
		rl.EndDrawing()
	}

	// close the window
	// this is called by pressing ESC
	rl.CloseWindow()
}
