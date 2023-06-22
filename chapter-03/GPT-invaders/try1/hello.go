package main

import (
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 800
	screenHeight = 450
)

type bullet struct {
	position    rl.Vector2
	speed       float32
	active      bool
	texture     rl.Texture2D
	frameWidth  int32
	frameHeight int32
	currentAnim int32
	frameCount  int32
	animSpeed   float32
}

type enemy struct {
	position    rl.Vector2
	speed       float32
	active      bool
	texture     rl.Texture2D
	frameWidth  int32
	frameHeight int32
	currentAnim int32
	frameCount  int32
	animSpeed   float32
}

func main() {
	rand.Seed(time.Now().UnixNano())

	rl.InitWindow(screenWidth, screenHeight, "Space Invaders")

	player := rl.NewRectangle(float32(screenWidth/2-20), float32(screenHeight-50), 40, 40)
	playerSpeed := float32(5)
	playerTexture := rl.LoadTexture("player.png")
	playerBullet := bullet{}
	playerBullet.speed = 10
	playerBullet.texture = rl.LoadTexture("player_bullet.png")
	playerBullet.frameWidth = playerBullet.texture.Width / 4
	playerBullet.frameHeight = playerBullet.texture.Height
	playerBullet.frameCount = 4
	playerBullet.animSpeed = 0.2

	enemies := make([]enemy, 0)
	for i := 0; i < 10; i++ {
		for j := 0; j < 5; j++ {
			enemy := enemy{}
			enemy.position = rl.NewVector2(float32(50+i*75), float32(50+j*50))
			enemy.speed = 1
			enemy.active = true
			enemy.texture = rl.LoadTexture("enemy.png")
			enemy.frameWidth = enemy.texture.Width / 4
			enemy.frameHeight = enemy.texture.Height
			enemy.frameCount = 4
			enemy.animSpeed = 0.2
			enemies = append(enemies, enemy)
		}
	}

	enemyBullet := bullet{}
	enemyBullet.speed = 5
	enemyBullet.texture = rl.LoadTexture("enemy_bullet.png")
	enemyBullet.frameWidth = enemyBullet.texture.Width / 4
	enemyBullet.frameHeight = enemyBullet.texture.Height
	enemyBullet.frameCount = 4
	enemyBullet.animSpeed = 0.2

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		// Update
		if rl.IsKeyDown(rl.KeyRight) {
			player.X += playerSpeed
		}
		if rl.IsKeyDown(rl.KeyLeft) {
			player.X -= playerSpeed
		}

		if playerBullet.active {
			playerBullet.position.Y -= playerBullet.speed

			if playerBullet.position.Y < 0 {
				playerBullet.active = false
			}
		}

		if enemyBullet.active {
			enemyBullet.position.Y += enemyBullet.speed

			if enemyBullet.position.Y > screenHeight {
				enemyBullet.active = false
			}
		}

		for i := range enemies {
			if enemies[i].active {
				enemies[i].position.X += enemies[i].speed

				if enemies[i].position.X > screenWidth-50 || enemies[i].position.X < 0 {
					for j := range enemies {
						enemies[j].speed = -enemies[j].speed
						enemies[j].position.Y += 10
					}

					if rand.Intn(1000) < 5 {
						if !enemyBullet.active {
							enemyBullet.active = true
							enemyBullet.position = rl.NewVector2(enemies[i].position.X+float32(enemies[i].frameWidth/2-enemyBullet.frameWidth/2), enemies[i].position.Y+float32(enemies[i].frameHeight))
						}
					}
				}
			}
		}

		// Collisions
		if rl.CheckCollisionRecs(player, rl.NewRectangle(enemyBullet.position.X, enemyBullet.position.Y, float32(enemyBullet.frameWidth), float32(enemyBullet.frameHeight))) {
			rl.CloseWindow()
		}

		for i := range enemies {
			if enemies[i].active {
				if rl.CheckCollisionRecs(rl.NewRectangle(playerBullet.position.X, playerBullet.position.Y, float32(playerBullet.texture.Width), float32(playerBullet.texture.Height)), rl.NewRectangle(enemies[i].position.X, enemies[i].position.Y, float32(enemies[i].frameWidth), float32(enemies[i].frameHeight))) {
					enemies[i].active = false
					playerBullet.active = false
				}
			}
		}

		// Draw
		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)

		if playerBullet.active {
			drawAnimatedTexture(playerBullet.texture, playerBullet.position, playerBullet.frameWidth, playerBullet.frameHeight, playerBullet.currentAnim, playerBullet.frameCount, playerBullet.animSpeed)
		}

		if enemyBullet.active {
			drawAnimatedTexture(enemyBullet.texture, enemyBullet.position, enemyBullet.frameWidth, enemyBullet.frameHeight, enemyBullet.currentAnim, enemyBullet.frameCount, enemyBullet.animSpeed)
		}

		for i := range enemies {
			if enemies[i].active {
				drawAnimatedTexture(enemies[i].texture, enemies[i].position, enemies[i].frameWidth, enemies[i].frameHeight, enemies[i].currentAnim, enemies[i].frameCount, enemies[i].animSpeed)
			}
		}

		drawAnimatedTexture(playerTexture, rl.NewVector2(player.X, player.Y), playerTexture.Width/4, playerTexture.Height, 0, 4, 0.1)

		rl.EndDrawing()
	}

	rl.UnloadTexture(playerTexture)
	rl.UnloadTexture(playerBullet.texture)
	rl.UnloadTexture(enemyBullet.texture)
	for i := range enemies {
		rl.UnloadTexture(enemies[i].texture)
	}

	rl.CloseWindow()
}

func drawAnimatedTexture(texture rl.Texture2D, position rl.Vector2, frameWidth, frameHeight, currentAnim, frameCount int32, animSpeed float32) {
	sourceRec := rl.NewRectangle(float32(currentAnim)*float32(frameWidth), 0, float32(frameWidth), float32(frameHeight))
	destRec := rl.NewRectangle(position.X, position.Y, float32(frameWidth), float32(frameHeight))
	rl.DrawTexturePro(texture, sourceRec, destRec, rl.NewVector2(0, 0), 0, rl.White)

	if int(rl.GetFrameTime()*60)%int(1/animSpeed*60) == 0 {
		currentAnim++
	}
	if currentAnim >= frameCount {
		currentAnim = 0
	}
}
