package main

import (
	"context"
	"fmt"
	"github.com/PullRequestInc/go-gpt3"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	godotenv.Load()

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatalln("Missing API KEY")
	}

	ctx := context.Background()
	client := gpt3.NewClient(apiKey)

	request := gpt3.CompletionRequest{
		Prompt:    []string{"How many cups of coffee should I drink per day?"},
		MaxTokens: gpt3.IntPtr(100),
	}

	client.CompletionStream(ctx, request, func(resp *gpt3.CompletionResponse) {
		fmt.Print(resp.Choices[0].Text)
	})

}
