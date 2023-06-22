package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 800
	screenHeight = 450
)

var (
	running         = true
	backgroundColor = rl.NewColor(147, 211, 196, 255)

	grassSprite  rl.Texture2D
	playerSprite rl.Texture2D

	playerSrc  rl.Rectangle
	playerDest rl.Rectangle

	playerMoving bool
	playerDir    PlayerDirection
	playerFrame  int
	frameCount   int

	playerSpeed float32
	musicPaused = false
	music       rl.Music
	cam         rl.Camera2D
)

type PlayerDirection int

const (
	Down PlayerDirection = iota
	Up
	Left
	Right
)

func init() {
	rl.SetConfigFlags(rl.FlagVsyncHint)
	rl.InitWindow(screenWidth, screenHeight, "Time for a coffee")
	rl.SetExitKey(0)
	rl.SetTargetFPS(60)

	grassSprite = rl.LoadTexture("assets/Tilesets/Grass.png")
	playerSprite = rl.LoadTexture("assets/Characters/Spritesheet.png")

	playerSrc = rl.NewRectangle(0, 0, 48, 48)
	playerDest = rl.NewRectangle(200, 200, 150, 150)

	playerSpeed = 3

	rl.InitAudioDevice()
	music = rl.LoadMusicStream("assets/music/Peanut_Plains_acoustic.mp3")
	musicPaused = false
	rl.PlayMusicStream(music)

	cam = rl.NewCamera2D(rl.NewVector2(screenWidth/2.0, screenHeight/2.0), rl.NewVector2(playerDest.X-playerDest.Width/2, playerDest.Y-playerDest.Height/2), 0.0, 1.0)
}

func update() {
	running = !rl.WindowShouldClose()

	if playerMoving {
		if playerDir == Up {
			playerDest.Y -= playerSpeed
		}
		if playerDir == Down {
			playerDest.Y += playerSpeed
		}
		if playerDir == Left {
			playerDest.X -= playerSpeed
		}
		if playerDir == Right {
			playerDest.X += playerSpeed
		}
		if frameCount%6 == 1 {
			playerFrame++
		}

	}

	frameCount++
	if playerFrame > 3 {
		playerFrame = 0
	}
	playerSrc.X = playerSrc.Width * float32(playerFrame)
	playerSrc.Y = playerSrc.Height * float32(playerDir)

	rl.UpdateMusicStream(music)
	if musicPaused {
		rl.PauseMusicStream(music)
	} else {
		rl.ResumeMusicStream(music)
	}

	cam.Target = rl.NewVector2(playerDest.X-playerDest.Width/2, playerDest.Y-playerDest.Height/2)

	playerMoving = false
}

func input() {
	if rl.IsKeyDown(rl.KeyW) || rl.IsKeyDown(rl.KeyUp) {
		playerMoving = true
		playerDir = Up
	}
	if rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyDown) {
		playerMoving = true
		playerDir = Down
	}
	if rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft) {
		playerMoving = true
		playerDir = Left
	}
	if rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyRight) {
		playerMoving = true
		playerDir = Right
	}

	if rl.IsKeyDown(rl.KeyQ) {
		musicPaused = !musicPaused
	}
}

func quit() {
	rl.UnloadTexture(grassSprite)
	rl.UnloadTexture(playerSprite)
	rl.CloseWindow()
}
func render() {
	rl.BeginDrawing()

	rl.ClearBackground(backgroundColor)

	rl.BeginMode2D(cam)
	drawScene()
	rl.EndMode2D()

	rl.EndDrawing()
}

func drawScene() {
	rl.DrawTexture(grassSprite, 100, 50, rl.White)
	rl.DrawTexturePro(playerSprite, playerSrc, playerDest, rl.NewVector2(playerDest.Width, playerDest.Height), 0, rl.White)
}

func main() {

	for running {
		input()
		update()
		render()
	}

	quit()

}