package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"math/rand"
	"time"
)

const (
	screenWidth  = 800
	screenHeight = 600
	cardWidth    = 71
	cardHeight   = 96
)

type Card struct {
	Suit  int
	Value int
	Rect  rl.Rectangle
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Create window and set fps
	rl.InitWindow(screenWidth, screenHeight, "Solitaire")
	rl.SetTargetFPS(60)

	// Create deck of cards
	deck := make([]Card, 52)
	i := 0
	for s := 0; s < 4; s++ {
		for v := 1; v <= 13; v++ {
			deck[i] = Card{
				Suit:  s,
				Value: v,
				Rect: rl.Rectangle{
					X:      float32(i%7)*(cardWidth+10) + 50,
					Y:      float32(i/7)*(cardHeight+10) + 50,
					Width:  cardWidth,
					Height: cardHeight,
				},
			}
			i++
		}
	}

	// Shuffle deck
	for i := range deck {
		j := rand.Intn(i + 1)
		deck[i], deck[j] = deck[j], deck[i]
	}

	// Main game loop
	for !rl.WindowShouldClose() {
		// Update
		for i := range deck {
			if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				mousePos := rl.GetMousePosition()
				if rl.CheckCollisionPointRec(mousePos, deck[i].Rect) {
					deck[i].Rect.Y -= 10
				}
			} else if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
				deck[i].Rect.Y += 10
			}
		}

		// Draw
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		for _, card := range deck {
			rl.DrawRectangleRec(card.Rect, rl.Red)
			rl.DrawText(fmt.Sprintf("%d", card.Value), int32(int(card.Rect.X+5)), int32(int(card.Rect.Y+5)), 20, rl.White)
		}
		rl.EndDrawing()
	}

	// Clean up
	rl.CloseWindow()
}
