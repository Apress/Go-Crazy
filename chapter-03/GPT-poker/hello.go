package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Card struct {
	Value int
	Suit  int
}

const (
	screenWidth  = 800
	screenHeight = 450

	cardWidth  = 70
	cardHeight = 100
)

var (
	deck    []Card
	players [4][]Card
)

func init() {
	for i := 1; i <= 13; i++ {
		for j := 1; j <= 4; j++ {
			deck = append(deck, Card{Value: i, Suit: j})
		}
	}

	// Shuffle the deck
	for i := len(deck) - 1; i > 0; i-- {
		j := rl.GetRandomValue(0, int32(i))
		deck[i], deck[j] = deck[j], deck[i]
	}

	// Deal the cards to the players
	for i := 0; i < 4; i++ {
		players[i] = make([]Card, 5)
		for j := 0; j < 5; j++ {
			players[i][j] = deck[i*5+j]
		}
	}
}

func main() {
	rl.InitWindow(screenWidth, screenHeight, "Poker")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		drawCards()

		rl.EndDrawing()
	}
}

func drawCards() {
	for i := 0; i < 4; i++ {
		for j := 0; j < 5; j++ {
			x := float32(50 + (cardWidth+10)*j)
			y := float32(50 + (cardHeight+10)*i)

			drawCardBack(x, y)

			if i == 0 {
				drawCard(players[i][j], x, y)
			}
		}
	}
}

func drawCardBack(x, y float32) {
	rl.DrawRectangleRec(rl.NewRectangle(x, y, cardWidth, cardHeight), rl.Black)
}

func drawCard(card Card, x, y float32) {
	suitColor := rl.Red
	if card.Suit%2 == 0 {
		suitColor = rl.Black
	}

	rl.DrawRectangleRec(rl.NewRectangle(x, y, cardWidth, cardHeight), rl.White)
	rl.DrawText(fmt.Sprintf("%d", card.Value), int32(x+10), int32(y+10), 28, rl.Black)

	if card.Value == 1 {
		rl.DrawCircle(int32(x+cardWidth/2), int32(y+cardHeight/2), 10, suitColor)
	} else {
		suitX := x + cardWidth/2
		suitY := y + cardHeight/2

		rl.DrawText(string(suitSymbol(card.Suit)), int32(suitX-10), int32(suitY-15), 28, suitColor)
		rl.DrawText(string(suitSymbol(card.Suit)), int32(suitX+10), int32(suitY+10), 28, suitColor)
	}
}

func suitSymbol(suit int) rune {
	switch suit {
	case 1:
		return '♥'
	case 2:
		return '♦'
	case 3:
		return '♣'
	case 4:
		return '♠'
	default:
		return '?'
	}

}
