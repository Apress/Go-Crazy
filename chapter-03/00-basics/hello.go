package main

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	screenWidth  = 800
	screenHeight = 450
)

var (
	running         = true
	backgroundColor = rl.RayWhite
)

func init() {
	rl.SetConfigFlags(rl.FlagVsyncHint)
	rl.InitWindow(screenWidth, screenHeight, "First Game")
	rl.SetExitKey(0)
	rl.SetTargetFPS(60)
}

func update() {
	running = !rl.WindowShouldClose()
}

func input() {

}

func quit() {
	rl.CloseWindow()
}
func render() {
	rl.BeginDrawing()

	rl.ClearBackground(backgroundColor)

	drawScene()

	rl.EndDrawing()
}

func drawScene() {
	rl.DrawText("Moyashi", 190, 200, 20, rl.Black)
}

func main() {

	for running {
		input()
		update()
		render()
	}

	quit()

}
