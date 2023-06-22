package main

import (
	"math/rand"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 800
	screenHeight = 450
)

type playerStruct struct {
	position rl.Vector2
	speed    float32
}

type bulletStruct struct {
	texture     rl.Texture2D
	position    rl.Vector2
	speed       float32
	active      bool
	frameWidth  int32
	frameHeight int32
	currentAnim int32
	frameCount  int32
	animSpeed   float32
}

type enemyStruct struct {
	texture     rl.Texture2D
	position    rl.Vector2
	speed       float32
	active      bool
	frameWidth  int32
	frameHeight int32
	currentAnim int32
	frameCount  int32
	animSpeed   float32
}

func main() {
	rl.InitWindow(screenWidth, screenHeight, "Space Invaders")

	playerTexture := rl.LoadTexture("player.png")
	player := playerStruct{rl.NewVector2(screenWidth/2-float32(playerTexture.Width/8)/2, screenHeight-float32(playerTexture.Height)), 5}

	playerBulletTexture := rl.LoadTexture("player_bullet.png")
	playerBullet := bulletStruct{playerBulletTexture, rl.Vector2{}, 10, false, playerBulletTexture.Width, playerBulletTexture.Height, 0, 1, 0}

	enemyBulletTexture := rl.LoadTexture("enemy_bullet.png")
	enemyBullet := bulletStruct{enemyBulletTexture, rl.Vector2{}, 5, false, enemyBulletTexture.Width, enemyBulletTexture.Height, 0, 1, 0}

	enemyTexture := rl.LoadTexture("enemy.png")
	enemies := make([]enemyStruct, 0)
	for x := 0; x < 10; x++ {
		for y := 0; y < 4; y++ {
			enemies = append(enemies, enemyStruct{enemyTexture, rl.NewVector2(float32(x*enemyTexture.Width/10+enemyTexture.Width/20), float32(y*enemyTexture.Height/4+enemyTexture.Height/4)), 1, true, enemyTexture.Width / 4, enemyTexture.Height, 0, 4, 0.1})
		}
	}

	score := 0

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		// Update
		if rl.IsKeyDown(rl.KeyRight) {
			if player.position.X < screenWidth-float32(playerTexture.Width/8) {
				player.position.X += player.speed
			}
		} else if rl.IsKeyDown(rl.KeyLeft) {
			if player.position.X > 0 {
				player.position.X -= player.speed
			}
		}

		if rl.IsKeyPressed(rl.KeySpace) {
			if !playerBullet.active {
				playerBullet.active = true
				playerBullet.position = rl.NewVector2(player.position.X+float32(playerTexture.Width/16), player.position.Y-float32(playerBulletTexture.Height))
			}
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

				// Reset the enemy bullet position
				enemyBullet.position = rl.NewVector2(enemies[rand.Intn(len(enemies))].position.X+float32(enemyTexture.Width/8), enemies[rand.Intn(len(enemies))].position.Y+float32(enemyTexture.Height))
			}
		} else {
			// Fire enemy bullet
			enemyBullet.active = true
			enemyBullet.position = rl.NewVector2(enemies[rand.Intn(len(enemies))].position.X+float32(enemyTexture.Width/8), enemies[rand.Intn(len(enemies))].position.Y+float32(enemyTexture.Height))
		}

		for i := range enemies {
			if enemies[i].active {
				enemies[i].position.X += enemies[i].speed
				if enemies[i].position.X > screenWidth-float32(enemies[i].frameWidth) || enemies[i].position.X < 0 {
					enemies[i].speed *= -1
					enemies[i].position.Y += float32(enemies[i].frameHeight)
				}

				if enemyBullet.active && rl.CheckCollisionPointRec(enemyBullet.position, rl.NewRectangle(enemies[i].position.X, enemies[i].position.Y, float32(enemies[i].frameWidth), float32(enemies[i].frameHeight))) {
					enemies[i].active = false
					enemyBullet.active = false
					score++
				}
			}
		}

		// Draw
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		// Draw player
		rl.DrawTextureRec(playerTexture, rl.NewRectangle(float32(playerTexture.Width/8)*float32(player.currentAnim), 0, float32(playerTexture.Width/8), float32(playerTexture.Height)), player.position, rl.White)

		// Draw player bullet
		if playerBullet.active {
			rl.DrawTexture(playerBullet.texture, int32(playerBullet.position.X), int32(playerBullet.position.Y), rl.White)
		}

		// Draw enemy bullet
		if enemyBullet.active {
			rl.DrawTexture(enemyBullet.texture, int32(enemyBullet.position.X), int32(enemyBullet.position.Y), rl.White)
		}

		// Draw enemies
		for i := range enemies {
			if enemies[i].active {
				rl.DrawTextureRec(enemies[i].texture, rl.NewRectangle(float32(enemies[i].frameWidth)*float32(enemies[i].currentAnim), 0, float32(enemies[i].frameWidth), float32(enemies[i].frameHeight)), enemies[i].position, rl.White)
			}
		}

		// Draw score
		rl.DrawText("Score: "+strconv.Itoa(score), 10, 10, 20, rl.White)

		rl.EndDrawing()
	}

	rl.UnloadTexture(playerTexture)
	rl.UnloadTexture(playerBulletTexture)
	rl.UnloadTexture(enemyBulletTexture)
	rl.UnloadTexture(enemyTexture)
}
