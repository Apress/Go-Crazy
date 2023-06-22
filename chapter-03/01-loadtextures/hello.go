package main

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	screenWidth  = 800
	screenHeight = 450
)

var (
	running         = true
	backgroundColor = rl.NewColor(147, 211, 196, 255)

	grassSprite rl.Texture2D
)

func init() {
	rl.SetConfigFlags(rl.FlagVsyncHint)
	rl.InitWindow(screenWidth, screenHeight, "Time for a coffee")
	rl.SetExitKey(0)
	rl.SetTargetFPS(60)

	grassSprite = rl.LoadTexture("assets/Tilesets/Grass.png")
}

func update() {
	running = !rl.WindowShouldClose()
}

func input() {

}

func quit() {
	rl.UnloadTexture(grassSprite)
	rl.CloseWindow()
}
func render() {
	rl.BeginDrawing()

	rl.ClearBackground(backgroundColor)

	drawScene()

	rl.EndDrawing()
}

func drawScene() {
	rl.DrawTexture(grassSprite, 100, 50, rl.White)
}

func main() {

	for running {
		input()
		update()
		render()
	}

	quit()

}
