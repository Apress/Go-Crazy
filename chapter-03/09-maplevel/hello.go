package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const (
	screenWidth  = 800
	screenHeight = 450
)

var (
	running         = true
	backgroundColor = rl.NewColor(147, 211, 196, 255)

	fencedSprite rl.Texture2D
	grassSprite  rl.Texture2D
	hillSprite   rl.Texture2D
	houseSprite  rl.Texture2D
	tilledSprite rl.Texture2D
	waterSprite  rl.Texture2D

	tex          rl.Texture2D
	playerSprite rl.Texture2D

	playerSrc  rl.Rectangle
	playerDest rl.Rectangle

	playerMoving                                  bool
	playerDir                                     int
	playerUp, playerDown, playerLeft, playerRight = false, false, false, false
	playerFrame                                   int
	frameCount                                    int

	tileDest   rl.Rectangle
	tileSrc    rl.Rectangle
	tileMap    []int
	srcMap     []string
	mapW, mapH int

	playerSpeed float32
	musicPaused = false
	music       rl.Music

	cam rl.Camera2D
)

func loadMap(mapFile string) {
	fmt.Printf("Loading map: %s\n", mapFile)
	file, err := ioutil.ReadFile(mapFile)
	if err != nil {
		fmt.Printf("Error reading map file: %s: %s\n", mapFile, err)
		os.Exit(1)
	}

	sliced := strings.Split(strings.ReplaceAll(string(file), "\n", " "), " ")

	mapW, mapH = -1, -1
	tileMap = make([]int, mapW*mapH)
	srcMap = make([]string, mapW*mapH)

	for i := 0; i < len(sliced); i++ {
		m, _ := strconv.Atoi(sliced[i])

		if mapW == -1 {
			mapW = m
		} else if mapH == -1 {
			mapH = m
		} else if i < mapW*mapH+2 {
			tileMap = append(tileMap, m)
		} else {
			srcMap = append(srcMap, sliced[i])
		}
	}

}

func init() {
	rl.SetConfigFlags(rl.FlagVsyncHint)
	rl.InitWindow(screenWidth, screenHeight, "Time for a coffee")
	rl.SetExitKey(0)
	rl.SetTargetFPS(60)

	fencedSprite = rl.LoadTexture("assets/Tilesets/Fences.png")
	grassSprite = rl.LoadTexture("assets/Tilesets/Grass.png")
	hillSprite = rl.LoadTexture("assets/Tilesets/Hills.png")
	houseSprite = rl.LoadTexture("assets/Tilesets/House.png")
	tilledSprite = rl.LoadTexture("assets/Tilesets/Tilled.png")
	waterSprite = rl.LoadTexture("assets/Tilesets/Water.png")

	tileSrc = rl.NewRectangle(0, 0, 16, 16)
	tileDest = rl.NewRectangle(0, 0, 16, 16)

	loadMap("world.map")

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
		if playerUp {
			playerDest.Y -= playerSpeed
		}
		if playerDown {
			playerDest.Y += playerSpeed
		}
		if playerLeft {
			playerDest.X -= playerSpeed
		}
		if playerRight {
			playerDest.X += playerSpeed
		}
		if frameCount%8 == 1 {
			playerFrame++
		}

	}
	if frameCount%45 == 1 {
		playerFrame++
	}
	if !playerMoving && playerFrame > 1 {
		playerFrame = 0
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
	playerDown, playerUp, playerRight, playerLeft = false, false, false, false
}

func input() {
	if rl.IsKeyDown(rl.KeyW) || rl.IsKeyDown(rl.KeyUp) {
		playerMoving = true
		playerDir = 1
		playerUp = true
	}
	if rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyDown) {
		playerMoving = true
		playerDir = 0
		playerDown = true
	}
	if rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft) {
		playerMoving = true
		playerDir = 2
		playerLeft = true
	}
	if rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyRight) {
		playerMoving = true
		playerDir = 3
		playerRight = true
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

	for i := 0; i < len(tileMap); i++ {
		tileDest.X = tileDest.Width * float32(i%mapW)
		tileDest.Y = tileDest.Height * float32(i/mapW)

		switch srcMap[i] {
		case "g":
			tex = grassSprite
		case "l":
			tex = hillSprite
		case "f":
			tex = fencedSprite
		case "h":
			tex = houseSprite
		case "w":
			tex = waterSprite
		case "t":
			tex = tilledSprite
		default:
			tex = grassSprite
		}

		tileSrc.X = tileSrc.Width * float32((tileMap[i]-1)%int(tex.Width/int32(tileSrc.Width)))
		tileSrc.Y = tileSrc.Height * float32((tileMap[i]-1)/int(tex.Width/int32(tileSrc.Height)))

		rl.DrawTexturePro(tex, tileSrc, tileDest, rl.NewVector2(tileDest.Width, tileDest.Height), 0, rl.White)

	}

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
