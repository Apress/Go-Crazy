package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
)

const (
	screenWidth  = 640
	screenHeight = 480
	maxGuesses   = 6
)

var (
	secretWord     string
	guessedLetters []string
	wrongGuesses   int
	gameOver       bool
	font           rl.Font
)

func main() {
	rand.Seed(time.Now().UnixNano())
	rl.InitWindow(screenWidth, screenHeight, "Hangman")
	rl.SetTargetFPS(60)

	font = rl.LoadFont("font.ttf")

	// Load word list
	words := loadWordList("words.txt")
	if len(words) == 0 {
		fmt.Println("Failed to load word list")
		return
	}

	// Select a random word
	secretWord = strings.ToUpper(words[rand.Intn(len(words))])

	// Initialize guessed letters array
	guessedLetters = make([]string, len(secretWord))
	for i := range guessedLetters {
		guessedLetters[i] = "_"
	}

	// Main game loop
	for !rl.WindowShouldClose() {
		// Handle input
		if rl.IsKeyPressed(rl.KeyR) {
			restartGame(words)
		} else if !gameOver {
			handleInput()
		}

		// Draw graphics
		rl.BeginDrawing()
		rl.ClearBackground(rl.White)
		drawHangman()
		drawWord()
		drawGuessedLetters()
		if gameOver {
			if gameWin() {
				drawGameWin()
			} else {
				drawGameOver()
			}

		}
		rl.EndDrawing()
	}

	rl.CloseWindow()
}

func loadWordList(filename string) []string {
	words := make([]string, 0)

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Failed to open file %s: %v\n", filename, err)
		return words
	}

	lines := strings.Split(string(file), "\n")
	for _, line := range lines {
		//line := rl.FgetS(file, rl.MaxStringLength)
		if line == "" {
			break
		}
		words = append(words, strings.TrimSpace(line))
	}

	return words
}

func restartGame(words []string) {
	secretWord = strings.ToUpper(words[rand.Intn(len(words))])
	guessedLetters = make([]string, len(secretWord))
	for i := range guessedLetters {
		guessedLetters[i] = "_"
	}
	wrongGuesses = 0
	gameOver = false
}

func handleInput() {
	// Check for new letter input
	for i := rl.KeyA; i <= rl.KeyZ; i++ {
		if rl.IsKeyPressed(int32(i)) {
			letter := string(i)
			if strings.Contains(secretWord, letter) {
				// Letter is in secret word, update guessed letters
				for j, c := range secretWord {
					if letter == string(c) {
						guessedLetters[j] = letter
					}
				}
				if gameWin() {
					gameOver = true
				}

			} else {
				// Letter is not in secret word, increment wrong guesses
				wrongGuesses++
				if wrongGuesses >= maxGuesses {
					// Too many wrong guesses, game over
					gameOver = true
				}
			}
		}
	}
}

func gameWin() bool {
	return !contains(guessedLetters, "_")
}

func drawHangman() {
	// Draw gallows
	rl.DrawLine(screenWidth/2-100, screenHeight/2-50, screenWidth/2-100, screenHeight/2+150, rl.Black)
	rl.DrawLine(screenWidth/2-130, screenHeight/2+130, screenWidth/2-70, screenHeight/2+170, rl.Black)

	switch wrongGuesses {
	case 1:
		// Draw head
		rl.DrawCircle(screenWidth/2-100, screenHeight/2+60, 20, rl.Black)
	case 2:
		// Draw body
		rl.DrawLine(screenWidth/2-100, screenHeight/2+80, screenWidth/2-100, screenHeight/2+130, rl.Black)
	case 3:
		// Draw left arm
		rl.DrawLine(screenWidth/2-100, screenHeight/2+80, screenWidth/2-130, screenHeight/2+100, rl.Black)
	case 4:
		// Draw right arm
		rl.DrawLine(screenWidth/2-100, screenHeight/2+80, screenWidth/2-70, screenHeight/2+100, rl.Black)
	case 5:
		// Draw left leg
		rl.DrawLine(screenWidth/2-100, screenHeight/2+130, screenWidth/2-130, screenHeight/2+170, rl.Black)
	case 6:
		// Draw right leg
		rl.DrawLine(screenWidth/2-100, screenHeight/2+130, screenWidth/2-70, screenHeight/2+170, rl.Black)
	}
}

func drawWord() {
	// Draw secret word with underscores for unknown letters
	x := screenWidth/2 - (len(secretWord) * 20 / 2)
	y := screenHeight/2 + 200
	for _, c := range secretWord {
		letter := string(c)
		if contains(guessedLetters, letter) {
			rl.DrawText(letter, int32(x), int32(y), 32, rl.Black)
		} else {
			rl.DrawText("", int32(x), int32(y), 32, rl.Black)
		}
		x += 20
	}
}

func drawGuessedLetters() {
	// Draw guessed letters
	letters := strings.Join(guessedLetters, " ")
	rl.DrawText("Guessed letters: "+letters, 10, 10, 20, rl.Black)
}

func drawGameOver() {
	// Draw game over message
	message := "Game over! The word was " + secretWord + ".\nPress R to restart."
	rl.DrawTextEx(font, message, rl.Vector2{X: 50, Y: screenHeight / 3}, 20, 1, rl.Black)
}

func drawGameWin() {
	// Draw game over message
	message := "You Win! The word was " + secretWord + ".\nPress R to restart."
	rl.DrawTextEx(font, message, rl.Vector2{X: 50, Y: screenHeight / 3}, 20, 1, rl.Black)
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
